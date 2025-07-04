package cli

import (
	"fmt"

	"github.com/shinzonetwork/view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeSchemaUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the default schemas from a remote source",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := service.UpdateDefaultSchemas(); err != nil {
				return fmt.Errorf("failed to update default schemas: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Default schemas updated from remote source.")
			return nil
		},
	}
	return cmd
}
