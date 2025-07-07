package cli

import (
	"github.com/spf13/cobra"
)

func MakeViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "Commands for working with views",
		Long:  "Use this command group to create, delete, and manage views in Viewkit.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if err := setContextViewStore(cmd); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.AddCommand(MakeViewInitCommand())
	cmd.AddCommand(MakeViewRollbackCommand())
	cmd.AddCommand(MakeViewDeleteCommand())
	cmd.AddCommand(MakeViewInspectCommand())
	cmd.AddCommand(MakeViewAddCommand())
	cmd.AddCommand(MakeViewRemoveCommand())

	return cmd
}
