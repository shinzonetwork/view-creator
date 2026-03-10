package cli

import (
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeSchemaAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <schema>",
		Short: "Add a custom schema type to the viewkit schema",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			schema := args[0]
			schemastore := mustGetContextSchemaStore(cmd)

			if err := service.AddCustomSchema(schemastore, schema); err != nil {
				return fmt.Errorf("failed to add schema: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Schema added successfully.")
			return nil
		},
	}

	return cmd
}
