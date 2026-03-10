package cli

import (
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeSchemaListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all schema types in the Viewkit schema (default and custom)",
		RunE: func(cmd *cobra.Command, args []string) error {
			schemastore := mustGetContextSchemaStore(cmd)

			defaults, customs, err := service.ListSchemas(schemastore)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Schemas:")
			for _, d := range defaults {
				fmt.Fprintf(cmd.OutOrStdout(), "• %s (default)\n", d)
			}
			for _, c := range customs {
				fmt.Fprintf(cmd.OutOrStdout(), "• %s (custom)\n", c)
			}

			return nil
		},
	}

	return cmd
}
