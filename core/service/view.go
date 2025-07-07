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
	"github.com/shinzonetwork/view-creator/core/schema"
	schemastore "github.com/shinzonetwork/view-creator/core/schema/store"
	"github.com/shinzonetwork/view-creator/core/util"
	viewstore "github.com/shinzonetwork/view-creator/core/view/store"
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

	var file io.ReadCloser

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

	if _, err := s.UploadAsset(name, label, file); err != nil {
		return models.View{}, fmt.Errorf("failed to upload asset: %w", err)
	}

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
