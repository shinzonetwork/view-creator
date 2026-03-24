package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeViewCreateCommand() *cobra.Command {
	var query string
	var sdl string
	var lenses []string
	var verbose bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a view with name, query, SDL, and lenses in one command",
		Long: `Create a complete view in a single command.

Lenses are specified with the --lens flag, which can be repeated.
Each lens value uses the format: label:path (or label:url for remote WASM).
To include arguments, use: label:path:{"key":"value"}

Examples:
  viewkit view create my-view \
    --query '{ ethTransactions { hash from to value } }' \
    --sdl 'type MyView { hash: String, from: String, to: String, value: String }' \
    --lens 'filter:./filter.wasm:{"address":"0xA0b8..."}' \
    --lens 'decode:./decode.wasm'
`,
		Args: cobra.ExactArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setContextViewStore(cmd); err != nil {
				return err
			}
			if err := setContextSchemaStore(cmd); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			viewstore := mustGetContextViewStore(cmd)
			schemastore := mustGetContextSchemaStore(cmd)

			// Init the view
			_, err := service.InitView(name, viewstore)
			if err != nil {
				return fmt.Errorf("failed to init view: %w", err)
			}

			// Cleanup on failure
			cleanup := func(err error) error {
				_ = service.DeleteView(name, viewstore)
				return err
			}

			// Add query if provided
			if query != "" {
				if _, err := service.UpdateQuery(name, query, viewstore, schemastore); err != nil {
					return cleanup(fmt.Errorf("failed to set query: %w", err))
				}
			}

			// Add SDL if provided
			if sdl != "" {
				if _, err := service.UpdateSDL(name, sdl, viewstore); err != nil {
					return cleanup(fmt.Errorf("failed to set SDL: %w", err))
				}
			}

			// Add lenses if provided
			for _, lensSpec := range lenses {
				label, path, lensArgs, err := parseLensSpec(lensSpec)
				if err != nil {
					return cleanup(fmt.Errorf("invalid lens spec %q: %w", lensSpec, err))
				}

				if _, err := service.InitLens(name, label, path, lensArgs, viewstore); err != nil {
					return cleanup(fmt.Errorf("failed to add lens %q: %w", label, err))
				}
			}

			// Load final view and print
			view, err := service.InspectView(name, viewstore)
			if err != nil {
				return cleanup(err)
			}

			printViewPretty(cmd, view, verbose, jsonOutput)
			return nil
		},
	}

	cmd.Flags().StringVar(&query, "query", "", "GraphQL query for the view")
	cmd.Flags().StringVar(&sdl, "sdl", "", "Schema definition for the view")
	cmd.Flags().StringArrayVar(&lenses, "lens", nil, "Lens in format label:path or label:path:{\"args\":\"json\"} (repeatable)")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output the view in raw JSON format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Show full output including revision history")

	return cmd
}

// parseLensSpec parses a lens specification string.
// Format: label:path or label:path:{"args":"json"}
func parseLensSpec(spec string) (label string, path string, args map[string]any, err error) {
	// Find the first colon to get the label
	firstColon := strings.Index(spec, ":")
	if firstColon == -1 {
		return "", "", nil, fmt.Errorf("expected format label:path or label:path:{args}")
	}

	label = spec[:firstColon]
	rest := spec[firstColon+1:]

	if label == "" {
		return "", "", nil, fmt.Errorf("label cannot be empty")
	}

	// Find where JSON args start (look for :{)
	jsonStart := strings.Index(rest, ":{")
	if jsonStart == -1 {
		// No args, rest is the path
		path = rest
		if path == "" {
			return "", "", nil, fmt.Errorf("path cannot be empty")
		}
		return label, path, nil, nil
	}

	path = rest[:jsonStart]
	argsStr := rest[jsonStart+1:]

	if path == "" {
		return "", "", nil, fmt.Errorf("path cannot be empty")
	}

	args = make(map[string]any)
	if err := json.Unmarshal([]byte(argsStr), &args); err != nil {
		return "", "", nil, fmt.Errorf("invalid args JSON: %w", err)
	}

	return label, path, args, nil
}
