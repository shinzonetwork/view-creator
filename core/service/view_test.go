package service_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store/local"
)

func TestViewService_CRUD(t *testing.T) {
	tempDir := t.TempDir()

	viewStore, err := local.NewLocalStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create view store: %v", err)
	}

	schemaStore, err := fileschema.NewFileSchemaStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create schema store: %v", err)
	}

	name := "testview"
	_, err = service.InitView(name, viewStore)
	if err != nil {
		t.Fatalf("InitView failed: %v", err)
	}

	view, err := service.InspectView(name, viewStore)
	if err != nil {
		t.Fatalf("InspectView failed: %v", err)
	}
	if view.Name != name {
		t.Errorf("expected view name %q, got %q", name, view.Name)
	}

	query := "TempLog { address }"
	if err := os.WriteFile(filepath.Join(tempDir, ".shinzo", "schema", "custom_schema.graphql"), []byte("type TempLog { address: String }"), 0644); err != nil {
		t.Fatalf("failed to write test schema: %v", err)
	}

	view, err = service.UpdateQuery(name, query, viewStore, schemaStore)
	if err != nil {
		t.Fatalf("UpdateQuery failed: %v", err)
	}
	if view.Query == nil || !strings.Contains(*view.Query, "address") {
		t.Errorf("expected query to contain 'address', got %v", view.Query)
	}

	sdl := "type Something @materialized(if: false) { x: String }"
	view, err = service.UpdateSDL(name, sdl, viewStore)
	if err != nil {
		t.Fatalf("UpdateSDL failed: %v", err)
	}
	if view.Sdl == nil || !strings.Contains(*view.Sdl, "Something") {
		t.Errorf("expected SDL to contain 'Something', got %v", view.Sdl)
	}

	view, err = service.ClearQuery(name, viewStore)
	if err != nil {
		t.Fatalf("ClearQuery failed: %v", err)
	}
	if view.Query != nil {
		t.Error("expected query to be cleared")
	}

	view, err = service.ClearSDL(name, viewStore)
	if err != nil {
		t.Fatalf("ClearSDL failed: %v", err)
	}
	if view.Sdl != nil {
		t.Error("expected SDL to be cleared")
	}

	err = service.DeleteView(name, viewStore)
	if err != nil {
		t.Fatalf("DeleteView failed: %v", err)
	}
}

func TestViewService_LensLifecycle(t *testing.T) {
	tempDir := t.TempDir()

	viewStore, err := local.NewLocalStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create view store: %v", err)
	}

	name := "lensview"
	_, err = service.InitView(name, viewStore)
	if err != nil {
		t.Fatalf("InitView failed: %v", err)
	}

	wasmPath := filepath.Join(tempDir, "test.wasm")

	// a minimal valid WASM binary header
	wasmBytes := []byte{0x00, 0x61, 0x73, 0x6D, 0x01, 0x00, 0x00, 0x00}
	if err := os.WriteFile(wasmPath, wasmBytes, 0644); err != nil {
		t.Fatalf("failed to write valid wasm file: %v", err)
	}

	view, err := service.InitLens(name, "testlens", wasmPath, map[string]any{"arg": "val"}, viewStore)
	if err != nil {
		t.Fatalf("InitLens failed: %v", err)
	}
	if len(view.Transform.Lenses) != 1 || view.Transform.Lenses[0].Label != "testlens" {
		t.Errorf("lens not properly added: %+v", view.Transform.Lenses)
	}

	view, err = service.RemoveLens(name, "testlens", viewStore)
	if err != nil {
		t.Fatalf("RemoveLens failed: %v", err)
	}
	if len(view.Transform.Lenses) != 0 {
		t.Error("lens was not removed properly")
	}
}
