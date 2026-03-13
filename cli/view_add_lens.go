package cli

import (
	"encoding/json"
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeAddLensCommand(viewName *string) *cobra.Command {
	var wasmPath string
	var wasmURL string
	var argsJson string
	var label string

	cmd := &cobra.Command{
		Use:   "lens",
		Short: "Add lenses in a view",
		RunE: func(cmd *cobra.Command, args []string) error {
			var path string

			if wasmURL != "" {
				path = wasmURL
			} else if wasmPath != "" {
				path = wasmPath
			} else {
				return fmt.Errorf("either --path or --url must be provided")
			}

			if label == "" {
				return fmt.Errorf("--label is required")
			}

			var argsMap map[string]any
			if argsJson != "" {
				if err := json.Unmarshal([]byte(argsJson), &argsMap); err != nil {
					return fmt.Errorf("invalid --args JSON: %w", err)
				}
			}

			store := mustGetContextViewStore(cmd)

			// function to add lens
			view, err := service.InitLens(*viewName, label, path, argsMap, store)
			if err != nil {
				return err
			}

			printViewPretty(cmd, view, false, false)

			return nil
		},
	}

	cmd.Flags().StringVar(&label, "label", "", "name of lens")
	cmd.Flags().StringVar(&wasmPath, "path", "", "Path to the WASM file (local)")
	cmd.Flags().StringVar(&wasmURL, "url", "", "URL to download the WASM file from")
	cmd.Flags().StringVar(&argsJson, "args", "", "arguments of the lens transform")

	return cmd
}
