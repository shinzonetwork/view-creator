package store

// SchemaStore defines a contract for managing GraphQL schemas used in the developer tool.
//
// This interface abstracts the storage and retrieval of schema files (default and custom),
// allowing implementations to store them locally, in-memory, or from remote sources.
//
// A schema is composed of two distinct parts:
// - Default schema: typically maintained by the CLI and updated from a remote source
// - Custom schema: user-defined types added and managed locally
type SchemaStore interface {
	// LoadD returns the full contents of the default schema and custom schema combined.
	Load() (string, error)

	// LoadDefault returns the full contents of the default schema.
	// This is typically read-only and managed by the CLI.
	LoadDefault() (string, error)

	// LoadCustom returns the full contents of the custom schema.
	// This contains user-defined types and is editable.
	LoadCustom() (string, error)

	// SaveCustom replaces the entire contents of the custom schema file.
	// Useful for batch edits or updates from an external source.
	SaveCustom(schema string) error

	// UpdateDefaultFromRemote fetches and replaces the default schema
	// from a remote location based on the provided version (e.g., branch or tag).
	// If version is empty, it defaults to "main".
	UpdateDefaultFromRemote(version string) error

	// ResetCustom clears the custom schema, removing all user-defined types.
	// Default schema remains untouched.
	ResetCustom() error

	// ListTypes returns all types found in the default and custom schemas.
	// Results are returned as two separate slices:
	//   - defaultTypes: from the default schema
	//   - customTypes: from the custom schema
	ListTypes() (defaultTypes []string, customTypes []string, err error)

	// GetTypeDefinition returns the full definition block of a type by name.
	// Searches both default and custom schemas.
	GetTypeDefinition(typeName string) (string, error)
}
