package cli

import (
	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeAddSdlCommand(viewName *string) *cobra.Command {
	return &cobra.Command{
		Use:   "sdl '<sdl>'",
		Short: "Add or update the sdl of the view",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := mustGetContextViewStore(cmd)
			view, err := service.UpdateSDL(*viewName, args[0], store)
			if err != nil {
				return err
			}

			printViewPretty(cmd, view, false, false)

			return nil
		},
	}
}
