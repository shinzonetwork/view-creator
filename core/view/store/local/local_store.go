package local

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/shinzonetwork/shinzo-view-creator/core/models"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store"
)

type LocalStore struct {
	BasePath string
}

func NewLocalStore(path ...string) (*LocalStore, error) {
	var base string

	// check if a custom path was provided
	if len(path) == 0 || path[0] == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("unable to get home directory: %w", err)
		}
		base = filepath.Join(home, ".shinzo", "views")
	} else {
		base = filepath.Join(path[0], ".shinzo", "views")
	}

	if err := os.MkdirAll(base, 0755); err != nil {
		return nil, fmt.Errorf("unable to create base directory: %w", err)
	}

	return &LocalStore{BasePath: base}, nil
}

func (s *LocalStore) Create(name string, timestamp string) (models.View, error) {
	// Create a new view folder with the name
	folderBasePath := filepath.Join(s.BasePath, name)

	// Check if the folder already exists
	if _, err := os.Stat(folderBasePath); err == nil {
		// Folder exists
		return models.View{}, store.ErrViewAlreadyExist
	} else if !os.IsNotExist(err) {
		// Some other unexpected error
		return models.View{}, fmt.Errorf("failed to check if view exists: %w", err)
	}

	if err := os.MkdirAll(folderBasePath, 0755); err != nil {
		return models.View{}, fmt.Errorf("failed to create view dir: %w", err)
	}

	// Create assets subdirectory in the new view folder
	assetsDir := filepath.Join(folderBasePath, "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		return models.View{}, fmt.Errorf("failed to create assets dir: %w", err)
	}

	view := models.View{
		Name:  name,
		Query: nil,
		Sdl:   nil,
		Transform: models.Transform{
			Lenses: []models.Lens{},
		},
		Metadata: models.Metadata{
			Version:   0,
			Total:     0,
			Revisions: []models.Revision{},
			CreatedAt: timestamp,
			UpdatedAt: timestamp,
		},
	}

	// create view file in the new folder dir
	viewFilePath := filepath.Join(folderBasePath, "view.json")

	file, err := os.Create(viewFilePath)
	if err != nil {
		return models.View{}, fmt.Errorf("failed to create view.json: %w", err)
	}
	defer file.Close()

	// insert view json into view file
	if err := json.NewEncoder(file).Encode(view); err != nil {
		return models.View{}, fmt.Errorf("failed to write JSON: %w", err)
	}

	return view, nil
}

func (s *LocalStore) Load(name string) (models.View, error) {
	folderBasePath := filepath.Join(s.BasePath, name)

	if _, err := os.Stat(folderBasePath); os.IsNotExist(err) {
		return models.View{}, store.ErrViewDoesNotExist
	} else if err != nil {
		return models.View{}, fmt.Errorf("failed to check if view exists: %w", err)
	}

	viewFilePath := filepath.Join(folderBasePath, "view.json")

	data, err := os.ReadFile(viewFilePath)
	if err != nil {
		return models.View{}, fmt.Errorf("failed to retrieve view.json file: %w", err)
	}

	var view models.View
	if err := json.Unmarshal([]byte(data), &view); err != nil {
		return models.View{}, fmt.Errorf("failed to unmarshall view.json file: %w", err)
	}

	return view, nil
}

func (s *LocalStore) List() ([]models.View, error) {
	entries, err := os.ReadDir(s.BasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read views directory: %w", err)
	}

	var views []models.View

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		viewFilePath := filepath.Join(s.BasePath, entry.Name(), "view.json")

		file, err := os.Open(viewFilePath)
		if err != nil {
			continue
		}
		defer file.Close()

		var view models.View
		if err := json.NewDecoder(file).Decode(&view); err != nil {
			continue
		}

		views = append(views, view)
	}

	return views, nil
}

func (s *LocalStore) Save(name string, view models.View) (models.View, error) {
	folderBasePath := filepath.Join(s.BasePath, name)

	// Check if view folder exists
	if _, err := os.Stat(folderBasePath); os.IsNotExist(err) {
		return models.View{}, store.ErrViewDoesNotExist
	} else if err != nil {
		return models.View{}, fmt.Errorf("failed to check if view exists: %w", err)
	}

	// Load the current state (before mutation)
	current, err := s.Load(name)
	if err != nil {
		return models.View{}, fmt.Errorf("failed to load current view before saving: %w", err)
	}

	// Generate revision snapshot and update metadata
	updatedMeta, err := MakeRevisionSnapshot(current.Metadata, current, view)
	if err != nil {
		return models.View{}, fmt.Errorf("failed to generate revision: %w", err)
	}
	view.Metadata = updatedMeta

	viewFilePath := filepath.Join(folderBasePath, "view.json")
	tempFilePath := filepath.Join(folderBasePath, "view.tmp.json")

	// Create temp file
	file, err := os.OpenFile(tempFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return models.View{}, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	// Write to temp file
	if err := json.NewEncoder(file).Encode(view); err != nil {
		return models.View{}, fmt.Errorf("failed to encode view to temp file: %w", err)
	}

	// Rename temp file to view.json (atomic move)
	if err := os.Rename(tempFilePath, viewFilePath); err != nil {
		return models.View{}, fmt.Errorf("failed to replace view.json: %w", err)
	}

	return view, nil
}

func (s *LocalStore) Delete(name string) error {
	folderBasePath := filepath.Join(s.BasePath, name)

	if _, err := os.Stat(folderBasePath); os.IsNotExist(err) {
		return store.ErrViewDoesNotExist
	} else if err != nil {
		return fmt.Errorf("failed to check if view exists: %w", err)
	}

	if err := os.RemoveAll(folderBasePath); err != nil {
		return fmt.Errorf("failed to delete view: %w", err)
	}

	return nil
}

func (s *LocalStore) UploadAsset(name string, label string, file io.Reader) (string, error) {
	folderBasePath := filepath.Join(s.BasePath, name)
	assetFolderPath := filepath.Join(folderBasePath, "assets")

	assetPath := filepath.Join(assetFolderPath, fmt.Sprintf("%s.wasm", label))
	outFile, err := os.Create(assetPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		return "", err
	}

	return assetPath, nil
}

func (s *LocalStore) DeleteAsset(viewName string, label string) error {
	folderBasePath := filepath.Join(s.BasePath, viewName)
	assetFolderPath := filepath.Join(folderBasePath, "assets")
	assetPath := filepath.Join(assetFolderPath, fmt.Sprintf("%s.wasm", label))

	if err := os.Remove(assetPath); err != nil {
		if os.IsNotExist(err) {
			return nil // already deleted or never existed
		}
		return fmt.Errorf("failed to delete asset: %w", err)
	}
	return nil
}

func (s *LocalStore) Rollback(viewName string, targetVersion int) (models.View, error) {
	view, err := s.Load(viewName)
	if err != nil {
		return models.View{}, err
	}

	var patchJSON []byte
	found := false
	for _, rev := range view.Metadata.Revisions {
		if rev.Version == targetVersion {
			patchJSON = []byte(rev.Diff)
			found = true
			break
		}
	}
	if !found {
		return models.View{}, fmt.Errorf("version %d not found", targetVersion)
	}

	currentJSON, err := json.Marshal(view)
	if err != nil {
		return models.View{}, err
	}
	meta := view.Metadata
	meta, err = MakeRevisionSnapshot(meta, view, view)
	if err != nil {
		return models.View{}, err
	}

	rolledBackJSON, err := jsonpatch.MergePatch(currentJSON, patchJSON)
	if err != nil {
		return models.View{}, fmt.Errorf("failed to apply patch: %w", err)
	}

	var rolledBackView models.View
	if err := json.Unmarshal(rolledBackJSON, &rolledBackView); err != nil {
		return models.View{}, fmt.Errorf("failed to unmarshal rolled back view: %w", err)
	}

	rolledBackView.Metadata = meta

	return s.Save(viewName, rolledBackView)
}

// GetAssetBlob finds the lens with the given label and returns its wasm blob as a base64 string.
func (s *LocalStore) GetAssetBlob(viewName string, lensLabel string) (string, error) {
	// Load the view
	view, err := s.Load(viewName)
	if err != nil {
		return "", fmt.Errorf("failed to load view: %w", err)
	}

	// Find the lens by label
	var lens *models.Lens
	for i := range view.Transform.Lenses {
		if view.Transform.Lenses[i].Label == lensLabel {
			lens = &view.Transform.Lenses[i]
			break
		}
	}
	if lens == nil {
		return "", fmt.Errorf("lens with label %q not found", lensLabel)
	}

	// Resolve path to wasm file
	assetPath := lens.Path
	if !filepath.IsAbs(assetPath) {
		assetPath = filepath.Join(s.BasePath, viewName, "assets", fmt.Sprintf("%s.wasm", lensLabel))
	}

	// Read the file
	data, err := os.ReadFile(assetPath)
	if err != nil {
		return "", fmt.Errorf("failed to read asset file %q: %w", assetPath, err)
	}

	// Encode to base64 and return
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}

func MakeRevisionSnapshot(meta models.Metadata, oldView any, newView any) (models.Metadata, error) {
	oldJSON, err := json.Marshal(oldView)
	if err != nil {
		return meta, fmt.Errorf("failed to marshal old view: %w", err)
	}

	newJSON, err := json.Marshal(newView)
	if err != nil {
		return meta, fmt.Errorf("failed to marshal new view: %w", err)
	}

	patch, err := jsonpatch.CreateMergePatch(newJSON, oldJSON)
	if err != nil {
		return meta, fmt.Errorf("failed to create patch: %w", err)
	}

	if string(patch) == "{}" {
		return meta, nil
	}

	revision := models.Revision{
		Version:   meta.Version,
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		Diff:      string(patch),
	}

	meta.Revisions = append(meta.Revisions, revision)
	meta.Version++
	meta.Total++
	meta.UpdatedAt = strconv.FormatInt(time.Now().Unix(), 10)

	return meta, nil
}

func ApplyPatch(original any, patchStr string) (any, error) {
	originalJSON, _ := json.Marshal(original)
	patch, err := jsonpatch.DecodePatch([]byte(patchStr))
	if err != nil {
		return nil, err
	}
	modified, err := patch.Apply(originalJSON)
	if err != nil {
		return nil, err
	}
	var result any
	if err := json.Unmarshal(modified, &result); err != nil {
		return nil, err
	}
	return result, nil
}
