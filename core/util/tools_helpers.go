package util

import (
	"os"
	"path/filepath"

	"github.com/shinzonetwork/view-creator/tools"
)

func EnsureModelsFile() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	target := filepath.Join(home, ".shinzo", "tools", "schema.graphql")

	if _, err := os.Stat(target); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return os.WriteFile(target, tools.DefaultSchema, 0644)
	}

	return nil
}
