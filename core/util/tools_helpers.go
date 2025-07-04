package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shinzonetwork/view-creator/tools"
)

func EnsureSchemaFilesExist() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	schemaDir := filepath.Join(home, ".shinzo", "tools", "schema")
	if err := os.MkdirAll(schemaDir, 0755); err != nil {
		return err
	}

	defaultPath := filepath.Join(schemaDir, "default_schema.graphql")
	customPath := filepath.Join(schemaDir, "custom_schema.graphql")

	if _, err := os.Stat(defaultPath); os.IsNotExist(err) || isFileEmpty(defaultPath) {
		if err := os.WriteFile(defaultPath, []byte(strings.TrimSpace(tools.DefaultSchema)), 0644); err != nil {
			return fmt.Errorf("failed to write default schema: %w", err)
		}
	}

	if _, err := os.Stat(customPath); os.IsNotExist(err) || isFileEmpty(customPath) {
		if err := os.WriteFile(customPath, []byte(""), 0644); err != nil {
			return fmt.Errorf("failed to write custom schema: %w", err)
		}
	}

	return nil
}

func isFileEmpty(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Size() == 0
}
