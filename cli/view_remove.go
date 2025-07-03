package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MakeViewRemoveCommand() *cobra.Command {
	var viewName string

	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove components to an existing view",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setContextStore(cmd); err != nil {
				return err
			}

			if viewName == "" {
				return fmt.Errorf("view name is required (use --name)")
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&viewName, "name", "", "Name of the view")

	cmd.AddCommand(MakeRemoveQueryCommand(&viewName))
	cmd.AddCommand(MakeRemoveSdlCommand(&viewName))
	cmd.AddCommand(MakeRemoveLensCommand(&viewName))

	return cmd
}
