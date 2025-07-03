package cli

import (
	"github.com/shinzonetwork/view-creator/core/service"
	"github.com/spf13/cobra"
)

func makeAddSdlCommand(viewName *string) *cobra.Command {
	return &cobra.Command{
		Use:   "sdl '<sdl>'",
		Short: "Add or update the sdl of the view",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := mustGetContextStore(cmd)
			view, err := service.UpdateSDL(*viewName, args[0], store)
			if err != nil {
				return err
			}

			printViewPretty(cmd, view, false, false)

			return nil
		},
	}
}
