package cli

import (
	"github.com/shinzonetwork/view-creator/core/service"
	"github.com/shinzonetwork/view-creator/core/util"
	"github.com/spf13/cobra"
)

func MakeAddQueryCommand(viewName *string) *cobra.Command {
	return &cobra.Command{
		Use:   "query '<query>'",
		Short: "Add or update the query of the view",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := util.EnsureSchemaFileExists(); err != nil {
				return err
			}
			store := mustGetContextStore(cmd)
			view, err := service.UpdateQuery(*viewName, args[0], store)
			if err != nil {
				return err
			}

			printViewPretty(cmd, view, false, false)

			return nil
		},
	}
}
