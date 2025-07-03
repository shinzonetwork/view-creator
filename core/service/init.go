package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shinzonetwork/view-creator/core/models"
	"github.com/shinzonetwork/view-creator/core/store"
	"github.com/shinzonetwork/view-creator/core/util"
)

func InitView(name string, s store.ViewStore) (models.View, error) {
	return s.Create(name, strconv.FormatInt(time.Now().Unix(), 10))
}

func InspectView(name string, s store.ViewStore) (models.View, error) {
	return s.Load(name)
}

func DeleteView(name string, s store.ViewStore) error {
	return s.Delete(name)
}

func UpdateQuery(name string, query string, s store.ViewStore) (models.View, error) {
	view, err := s.Load(name)
	if err != nil {
		return models.View{}, err
	}

	// TODO: validate query

	view.Query = &query

	view, err = s.Save(name, view)
	if err != nil {
		return models.View{}, err
	}

	return view, nil
}

func UpdateSDL(name string, sdl string, s store.ViewStore) (models.View, error) {
	view, err := s.Load(name)
	if err != nil {
		return models.View{}, err
	}

	// TODO: validate sdl

	view.Sdl = &sdl

	view, err = s.Save(name, view)
	if err != nil {
		return models.View{}, err
	}

	return view, nil
}

func ClearSDL(name string, s store.ViewStore) (models.View, error) {
	view, err := s.Load(name)
	if err != nil {
		return models.View{}, err
	}

	view.Sdl = nil

	return s.Save(name, view)
}

func ClearQuery(name string, s store.ViewStore) (models.View, error) {
	view, err := s.Load(name)
	if err != nil {
		return models.View{}, err
	}

	view.Query = nil

	return s.Save(name, view)
}

func InitLens(name string, label string, path string, args map[string]any, s store.ViewStore) (models.View, error) {
	view, err := s.Load(name)
	if err != nil {
		return models.View{}, err
	}

	// Check if lens already exists
	for _, lens := range view.Transform.Lenses {
		if lens.Label == label {
			return models.View{}, fmt.Errorf(`lens with label "%s" already exists`, label)
		}
	}

	var file io.ReadCloser

	// Check if it's a URL
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			return models.View{}, fmt.Errorf("failed to download from URL: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return models.View{}, fmt.Errorf("download failed with HTTP %d", resp.StatusCode)
		}
		file = resp.Body
	} else {
		localFile, err := os.Open(path)
		if err != nil {
			return models.View{}, fmt.Errorf("failed to open local file: %w", err)
		}
		file = localFile
	}
	defer file.Close()

	_, err = util.IsValidWasm(file)
	if err != nil {
		return models.View{}, fmt.Errorf("invalid wasm file: %w", err)
	}

	// Upload the asset
	if _, err := s.UploadAsset(name, label, file); err != nil {
		return models.View{}, fmt.Errorf("failed to upload asset: %w", err)
	}

	// Add the new lens with label and args
	newLens := models.Lens{
		Label:     label,
		Arguments: args,
		Path:      fmt.Sprintf("assets/%s.wasm", label),
	}
	view.Transform.Lenses = append(view.Transform.Lenses, newLens)

	// Save the updated view
	return s.Save(name, view)
}

func RemoveLens(name string, label string, s store.ViewStore) (models.View, error) {
	view, err := s.Load(name)
	if err != nil {
		return models.View{}, err
	}

	var (
		updatedLenses []models.Lens
		found         bool
	)

	for _, lens := range view.Transform.Lenses {
		if lens.Label == label {
			found = true
			continue // skip the one we want to remove
		}
		updatedLenses = append(updatedLenses, lens)
	}

	if !found {
		return models.View{}, fmt.Errorf(`lens with label "%s" not found`, label)
	}

	// Update view with remaining lenses
	view.Transform.Lenses = updatedLenses

	// Delete the wasm file associated with the lens
	if err := s.DeleteAsset(name, label); err != nil {
		return models.View{}, fmt.Errorf("failed to delete lens asset: %w", err)
	}

	// Save and return the updated view
	return s.Save(name, view)
}
