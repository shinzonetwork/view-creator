package cli

import "github.com/spf13/cobra"

func MakeNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage Viewkit local node",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setContextSchemaStore(cmd); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
