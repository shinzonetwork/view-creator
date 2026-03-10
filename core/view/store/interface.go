package store

import (
	"io"

	"github.com/shinzonetwork/shinzo-view-creator/core/models"
)

// ViewStore defines a contract for persisting and retrieving views across various storage.
//
// This interface abstracts the storage mechanism for views, allowing implementations to use
// local file systems, cloud storage, GitHub, databases, or remote services.
//
// All methods operate on a view identified by its name.
type ViewStore interface {
	// Create initializes and stores a new view with the given name and timestamp.
	// Returns the newly created View.
	Create(name string, timestamp string) (models.View, error)

	// Load retrieves a view by its name.
	// Returns the loaded View or an empty View if not found.
	Load(name string) (models.View, error)

	// List returns all currently stored views.
	List() ([]models.View, error)

	// Save persists updates to a view identified by its name.
	// Returns the updated View.
	Save(name string, view models.View) (models.View, error)

	// Delete removes the view identified by name from the store.
	// Returns the deleted View.
	Delete(name string) error

	// Uploads an asset (e.g. .wasm) to a view, identified by label.
	UploadAsset(viewName string, label string, file io.Reader) (string, error)

	// Deletes an asset by view name and label.
	DeleteAsset(viewName string, label string) error

	GetAssetBlob(viewName string, label string) (string, error)

	// Revert the view back to a previous version
	Rollback(viewName string, version int) (models.View, error)
}
