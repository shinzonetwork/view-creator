package cli

import (
	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeRemoveSdlCommand(viewName *string) *cobra.Command {
	return &cobra.Command{
		Use:   "sdl",
		Short: "Remove the SDL from the view",
		RunE: func(cmd *cobra.Command, args []string) error {
			store := mustGetContextViewStore(cmd)

			view, err := service.ClearSDL(*viewName, store)
			if err != nil {
				return err
			}

			printViewPretty(cmd, view, false, false)
			return nil
		},
	}
}
