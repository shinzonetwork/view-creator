package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MakeViewAddCommand() *cobra.Command {
	var viewName string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add components to an existing view",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setContextViewStore(cmd); err != nil {
				return err
			}

			if err := setContextSchemaStore(cmd); err != nil {
				return err
			}

			if viewName == "" {
				return fmt.Errorf("view name is required (use --name)")
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&viewName, "name", "", "Name of the view")

	cmd.AddCommand(MakeAddQueryCommand(&viewName))
	cmd.AddCommand(MakeAddSdlCommand(&viewName))
	cmd.AddCommand(MakeAddLensCommand(&viewName))

	return cmd
}
