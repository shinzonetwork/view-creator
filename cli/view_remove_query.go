package cli

import (
	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeRemoveQueryCommand(viewName *string) *cobra.Command {
	return &cobra.Command{
		Use:   "query",
		Short: "Remove the query from the view",
		RunE: func(cmd *cobra.Command, args []string) error {
			store := mustGetContextViewStore(cmd)

			view, err := service.ClearQuery(*viewName, store)
			if err != nil {
				return err
			}

			printViewPretty(cmd, view, false, false)
			return nil
		},
	}
}
