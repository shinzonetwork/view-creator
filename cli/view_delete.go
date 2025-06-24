package cli

import (
	"github.com/shinzonetwork/view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeViewDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delet a saved view",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			storeImpl := mustGetContextStore(cmd)
			err := service.DeleteView(args[0], storeImpl)
			if err != nil {
				return err
			}

			cmd.Printf("deleted view %s Successfully\n", args[0])
			return nil
		},
	}

	return cmd
}
