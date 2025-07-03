package local

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shinzonetwork/view-creator/core/models"
	"github.com/shinzonetwork/view-creator/core/store"
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
