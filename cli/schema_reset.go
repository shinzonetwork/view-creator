package cli

import (
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeSchemaResetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Clear all custom schema types (does not affect defaults)",
		RunE: func(cmd *cobra.Command, args []string) error {
			schemastore := mustGetContextSchemaStore(cmd)

			if err := service.ResetCustomSchemas(schemastore); err != nil {
				return fmt.Errorf("failed to reset custom schemas: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Custom schema cleared.")
			return nil
		},
	}
	return cmd
}
