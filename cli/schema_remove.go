package cli

import (
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeSchemaRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove <type>",
		Short: "Remove a custom schema type from the viewkit schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			schemastore := mustGetContextSchemaStore(cmd)

			if err := service.RemoveCustomSchema(schemastore, name); err != nil {
				return fmt.Errorf("failed to remove schema: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Schema removed.")
			return nil
		},
	}

	return cmd
}
