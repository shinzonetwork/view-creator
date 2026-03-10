package cli

import (
	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeViewTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test <name>",
		Short: "Test if the view can build and compile successfully",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			viewstore := mustGetContextViewStore(cmd)
			schemastore, err := fileschema.NewFileSchemaStore()
			if err != nil {
				return err
			}

			viewName := args[0]

			return service.StartLocalNodeAndTestView(viewName, viewstore, schemastore, false)
		},
	}

	return cmd
}
