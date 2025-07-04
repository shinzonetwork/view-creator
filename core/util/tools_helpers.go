package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shinzonetwork/view-creator/tools"
)

func EnsureSchemaFileExists() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	schemaPath := filepath.Join(home, ".shinzo", "tools", "schema.graphql")

	_, statErr := os.Stat(schemaPath)
	if os.IsNotExist(statErr) || isFileEmpty(schemaPath) {
		content := fmt.Sprintf(
			"# --- DEFAULT SCHEMAS ---\n\n%s\n\n# --- CUSTOM SCHEMAS ---\n",
			strings.TrimSpace(tools.DefaultSchema),
		)
		if err := os.MkdirAll(filepath.Dir(schemaPath), 0755); err != nil {
			return err
		}
		return os.WriteFile(schemaPath, []byte(content), 0644)
	}

	return nil
}

func isFileEmpty(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Size() == 0
}
