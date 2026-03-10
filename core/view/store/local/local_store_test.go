package local_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/shinzonetwork/shinzo-view-creator/core/models"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store"
	"github.com/shinzonetwork/shinzo-view-creator/core/view/store/local"
)

func TestLocalStoreInitCreateDir(t *testing.T) {
	temp := t.TempDir()

	store, err := local.NewLocalStore(temp)
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	expected := filepath.Join(temp, ".shinzo", "views")

	info, err := os.Stat(expected)
	if err != nil {
		t.Fatalf("expected dir %s to exist, got error: %v", expected, err)
	}

	if !info.IsDir() {
		t.Fatalf("expected %s to be a directory", expected)
	}

	_ = store
}

func TestLocalStoreCreateViewDir(t *testing.T) {
	temp := t.TempDir()

	store, err := local.NewLocalStore(temp)
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	name := "testexample"
	timestamp := "1750696562"

	_, err = store.Create(name, timestamp)
	if err != nil {
		t.Fatalf("failed to create view: %v", err)
	}

	base := filepath.Join(temp, ".shinzo", "views", name)

	// 1. Check view folder
	if stat, err := os.Stat(base); err != nil || !stat.IsDir() {
		t.Fatalf("expected view dir at %s, got err: %v", base, err)
	}

	// 2. Check view.json file
	viewJSON := filepath.Join(base, "view.json")
	if stat, err := os.Stat(viewJSON); err != nil || stat.IsDir() {
		t.Fatalf("expected view.json file at %s, got err: %v", viewJSON, err)
	}

	// 3. Check assets directory
	assets := filepath.Join(base, "assets")
	if stat, err := os.Stat(assets); err != nil || !stat.IsDir() {
		t.Fatalf("expected assets dir at %s, got err: %v", assets, err)
	}

	// validate contents of view.json
	data, err := os.ReadFile(viewJSON)
	if err != nil {
		t.Fatalf("failed to read view.json: %v", err)
	}

	var view models.View
	if err := json.Unmarshal(data, &view); err != nil {
		t.Fatalf("failed to unmarshal view.json: %v", err)
	}

	// Validate fields
	if view.Name != name {
		t.Errorf("expected view name %q, got %q", name, view.Name)
	}

	if view.Metadata.CreatedAt != timestamp {
		t.Errorf("expected createdAt %q, got %q", timestamp, view.Metadata.CreatedAt)
	}

	if view.Metadata.UpdatedAt != timestamp {
		t.Errorf("expected updatedAt %q, got %q", timestamp, view.Metadata.UpdatedAt)
	}
}

func TestLocalStoreCreateDuplicate(t *testing.T) {
	temp := t.TempDir()

	localstore, err := local.NewLocalStore(temp)
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	name := "dupe-view"
	timestamp := "1750696562"

	// First creation should succeed
	if _, err := localstore.Create(name, timestamp); err != nil {
		t.Fatalf("failed to create first view: %v", err)
	}

	// Second creation with the same name should fail
	_, err = localstore.Create(name, timestamp)
	if err == nil {
		t.Fatalf("expected error when creating duplicate view, got nil")
	}

	if err != store.ErrViewAlreadyExist {
		t.Errorf("expected ErrViewAlreadyExists, got: %v", err)
	}
}

func TestLocalStoreLoadViewInDir(t *testing.T) {
	temp := t.TempDir()

	localstore, err := local.NewLocalStore(temp)
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	name := "testexample"
	timestamp := "1750696562"

	_, err = localstore.Create(name, timestamp)
	if err != nil {
		t.Fatalf("failed to create view: %v", err)
	}

	view, err := localstore.Load(name)
	if err != nil {
		t.Fatalf("failed to create load view: %v", err)
	}

	// Validate fields
	if view.Name != name {
		t.Errorf("expected view name %q, got %q", name, view.Name)
	}

	if view.Metadata.CreatedAt != timestamp {
		t.Errorf("expected createdAt %q, got %q", timestamp, view.Metadata.CreatedAt)
	}

	if view.Metadata.UpdatedAt != timestamp {
		t.Errorf("expected updatedAt %q, got %q", timestamp, view.Metadata.UpdatedAt)
	}
}

func TestLocalStoreShouldNotLoadViewThatDontExists(t *testing.T) {
	temp := t.TempDir()

	localstore, err := local.NewLocalStore(temp)
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	name := "testexample"

	// load should return error because view does not exist
	_, err = localstore.Load(name)
	if err == nil {
		t.Fatalf("expected error when loading not exisitng view, got nil")
	}

	if err != store.ErrViewDoesNotExist {
		t.Errorf("expected ErrViewDoesNotExist, got: %v", err)
	}
}

func TestLocalStoreDeleteRemovesViewDir(t *testing.T) {
	temp := t.TempDir()

	localstore, err := local.NewLocalStore(temp)
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	name := "testexample"
	timestamp := "1750696562"

	_, err = localstore.Create(name, timestamp)
	if err != nil {
		t.Fatalf("failed to create view: %v", err)
	}

	// Delete it
	if err := localstore.Delete(name); err != nil {
		t.Fatalf("failed to delete view: %v", err)
	}

	// Assert folder is gone
	deletedPath := filepath.Join(temp, ".shinzo", "views", name)
	if _, err := os.Stat(deletedPath); !os.IsNotExist(err) {
		t.Errorf("expected view folder to be deleted, but it still exists: %v", err)
	}
}

func TestLocalStoreDeleteNonExistentView(t *testing.T) {
	temp := t.TempDir()

	localstore, err := local.NewLocalStore(temp)
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	err = localstore.Delete("testexample")
	if err == nil {
		t.Fatal("expected error when deleting non-existent view, got nil")
	}

	if err != store.ErrViewDoesNotExist {
		t.Errorf("expected ErrViewDoesNotExist, got: %v", err)
	}
}

func TestLocalStoreListReturnsCreatedViews(t *testing.T) {
	temp := t.TempDir()

	localstore, err := local.NewLocalStore(temp)
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	viewsToCreate := []string{"alpha", "beta", "gamma"}

	for _, name := range viewsToCreate {
		_, err := localstore.Create(name, "1750696562")
		if err != nil {
			t.Fatalf("failed to create view %s: %v", name, err)
		}
	}

	views, err := localstore.List()
	if err != nil {
		t.Fatalf("failed to list views: %v", err)
	}

	if len(views) != len(viewsToCreate) {
		t.Errorf("expected %d views, got %d", len(viewsToCreate), len(views))
	}

	names := map[string]bool{}
	for _, v := range views {
		names[v.Name] = true
	}

	for _, name := range viewsToCreate {
		if !names[name] {
			t.Errorf("expected view %q to be in list, but it wasn't", name)
		}
	}
}

func TestLocalStoreSaveUpdatesViewJson(t *testing.T) {
	temp := t.TempDir()
	store, err := local.NewLocalStore(temp)
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	name := "testexample"
	timestamp := "1750696562"

	original, err := store.Create(name, timestamp)
	if err != nil {
		t.Fatalf("failed to create view: %v", err)
	}

	// Use the provided query and SDL
	original.Query = String("Log {address topics data transactionHash blockNumber}")
	original.Sdl = String("type FilteredAndDecodedLogs @materialized(if: false) {hash: String block: String address: String signature: String }")

	_, err = store.Save(name, original)
	if err != nil {
		t.Fatalf("failed to save updated view: %v", err)
	}

	// Load it back from disk
	loaded, err := store.Load(name)
	if err != nil {
		t.Fatalf("failed to load saved view: %v", err)
	}

	if loaded.Query == nil || *loaded.Query != *original.Query {
		t.Errorf("expected Query to be %q, got: %v", *original.Query, *loaded.Query)
	}

	if loaded.Sdl == nil || *loaded.Sdl != *original.Sdl {
		t.Errorf("expected Sdl to be %q, got: %v", *original.Sdl, *loaded.Sdl)
	}
}

func TestLocalStoreLifecycle(t *testing.T) {
	temp := t.TempDir()
	localstore, err := local.NewLocalStore(temp)
	if err != nil {
		t.Fatalf("failed to initialize store: %v", err)
	}

	name := "fulltest"
	timestamp := "1750696562"

	// 1. Create
	view, err := localstore.Create(name, timestamp)
	if err != nil {
		t.Fatalf("failed to create view: %v", err)
	}

	// 2. Update and Save
	view.Query = String("query { field }")
	view.Sdl = String("type T { field: String }")
	_, err = localstore.Save(name, view)
	if err != nil {
		t.Fatalf("failed to save view: %v", err)
	}

	// 3. Load
	loaded, err := localstore.Load(name)
	if err != nil {
		t.Fatalf("failed to load view: %v", err)
	}

	if loaded.Query == nil || *loaded.Query != *view.Query {
		t.Errorf("expected Query %q, got %v", *view.Query, loaded.Query)
	}
	if loaded.Sdl == nil || *loaded.Sdl != *view.Sdl {
		t.Errorf("expected Sdl %q, got %v", *view.Sdl, loaded.Sdl)
	}

	// 4. List
	listed, err := localstore.List()
	if err != nil {
		t.Fatalf("failed to list views: %v", err)
	}
	if len(listed) != 1 || listed[0].Name != name {
		t.Errorf("expected list to contain %q, got %+v", name, listed)
	}

	// 5. Delete
	err = localstore.Delete(name)
	if err != nil {
		t.Fatalf("failed to delete view: %v", err)
	}

	// 6. Confirm Deletion
	_, err = localstore.Load(name)
	if err != store.ErrViewDoesNotExist {
		t.Errorf("expected ErrViewDoesNotExist after delete, got %v", err)
	}
}

func String(s string) *string {
	return &s
}
