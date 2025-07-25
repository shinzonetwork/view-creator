package cli

import (
	"github.com/spf13/cobra"
)

func MakeSchemaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Manage custom models in the Viewkit schema",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setContextSchemaStore(cmd); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(MakeSchemaAddCommand())
	cmd.AddCommand(MakeSchemaListCommand())
	cmd.AddCommand(MakeSchemaUpdateCommand())
	cmd.AddCommand(MakeSchemaRemoveCommand())
	cmd.AddCommand(MakeSchemaInspectCommand())
	cmd.AddCommand(MakeSchemaResetCommand())

	return cmd
}
