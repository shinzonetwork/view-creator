package cli

import (
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeSchemaInspectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect <type>",
		Short: "Show the full definition of a schema type",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			schemastore := mustGetContextSchemaStore(cmd)

			def, err := service.GetSchemaTypeDefinition(schemastore, args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), def)
			return nil
		},
	}
	return cmd
}
