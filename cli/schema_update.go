package cli

import (
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeSchemaUpdateCommand() *cobra.Command {
	var version string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the default schemas from a remote source",
		RunE: func(cmd *cobra.Command, args []string) error {
			schemaStore := mustGetContextSchemaStore(cmd)

			if err := service.UpdateDefaultSchemas(schemaStore, version); err != nil {
				return fmt.Errorf("failed to update default schemas: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Default schemas updated from remote source.")
			return nil
		},
	}

	cmd.Flags().StringVar(&version, "version", "", "Git branch or tag to fetch the default schema from (default: main)")
	return cmd
}
