package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shinzonetwork/shinzo-view-creator/core/models"
	"github.com/shinzonetwork/shinzo-view-creator/core/schema"
	schemastore "github.com/shinzonetwork/shinzo-view-creator/core/schema/store"
	"github.com/shinzonetwork/shinzo-view-creator/core/util"
	viewstore "github.com/shinzonetwork/shinzo-view-creator/core/view/store"
)

func InitView(name string, s viewstore.ViewStore) (models.View, error) {
	return s.Create(name, strconv.FormatInt(time.Now().Unix(), 10))
}

func InspectView(name string, s viewstore.ViewStore) (models.View, error) {
	return s.Load(name)
}

func DeleteView(name string, s viewstore.ViewStore) error {
	return s.Delete(name)
}

func UpdateQuery(name string, query string, viewstore viewstore.ViewStore, schemastore schemastore.SchemaStore) (models.View, error) {
	view, err := viewstore.Load(name)
	if err != nil {
		return models.View{}, err
	}

	if err := schema.ValidateQuery(schemastore, query); err != nil {
		return models.View{}, err
	}

	view.Query = &query

	view, err = viewstore.Save(name, view)
	if err != nil {
		return models.View{}, err
	}

	return view, nil
}

func UpdateSDL(name string, sdl string, s viewstore.ViewStore) (models.View, error) {
	view, err := s.Load(name)
	if err != nil {
		return models.View{}, err
	}

	if err := util.ValidateSDL(sdl); err != nil {
		return models.View{}, err
	}

	view.Sdl = &sdl

	view, err = s.Save(name, view)
	if err != nil {
		return models.View{}, err
	}

	return view, nil
}

func ClearSDL(name string, s viewstore.ViewStore) (models.View, error) {
	view, err := s.Load(name)
	if err != nil {
		return models.View{}, err
	}

	view.Sdl = nil

	return s.Save(name, view)
}

func ClearQuery(name string, s viewstore.ViewStore) (models.View, error) {
	view, err := s.Load(name)
	if err != nil {
		return models.View{}, err
	}

	view.Query = nil

	return s.Save(name, view)
}

func InitLens(name string, label string, path string, args map[string]any, s viewstore.ViewStore) (models.View, error) {
	view, err := s.Load(name)
	if err != nil {
		return models.View{}, err
	}

	for _, lens := range view.Transform.Lenses {
		if lens.Label == label {
			return models.View{}, fmt.Errorf(`lens with label "%s" already exists`, label)
		}
	}

	// Read file into memory
	var wasmBytes []byte
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		resp, err := http.Get(path)
		if err != nil {
			return models.View{}, fmt.Errorf("failed to download from URL: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return models.View{}, fmt.Errorf("download failed with HTTP %d", resp.StatusCode)
		}

		wasmBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return models.View{}, fmt.Errorf("failed to read wasm from response: %w", err)
		}
	} else {
		wasmBytes, err = os.ReadFile(path)
		if err != nil {
			return models.View{}, fmt.Errorf("failed to read local wasm file: %w", err)
		}
	}

	// Validate WASM
	if _, err := util.IsValidWasm(bytes.NewReader(wasmBytes)); err != nil {
		return models.View{}, fmt.Errorf("invalid wasm file: %w", err)
	}

	// Upload asset
	if _, err := s.UploadAsset(name, label, bytes.NewReader(wasmBytes)); err != nil {
		return models.View{}, fmt.Errorf("failed to upload asset: %w", err)
	}

	// Add to view
	newLens := models.Lens{
		Label:     label,
		Arguments: args,
		Path:      fmt.Sprintf("assets/%s.wasm", label),
	}
	view.Transform.Lenses = append(view.Transform.Lenses, newLens)

	return s.Save(name, view)
}

func RemoveLens(name string, label string, s viewstore.ViewStore) (models.View, error) {
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
			continue
		}
		updatedLenses = append(updatedLenses, lens)
	}

	if !found {
		return models.View{}, fmt.Errorf(`lens with label "%s" not found`, label)
	}

	view.Transform.Lenses = updatedLenses

	if err := s.DeleteAsset(name, label); err != nil {
		return models.View{}, fmt.Errorf("failed to delete lens asset: %w", err)
	}

	return s.Save(name, view)
}

func Rollback(name string, version int, s viewstore.ViewStore) (models.View, error) {
	return s.Rollback(name, version)
}
