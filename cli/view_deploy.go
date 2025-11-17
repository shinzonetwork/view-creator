package cli

import (
	"fmt"

	"github.com/shinzonetwork/view-creator/core/schema/store/fileschema"
	"github.com/shinzonetwork/view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeViewDeployCommand() *cobra.Command {
	var target string

	cmd := &cobra.Command{
		Use:   "deploy <name>",
		Short: "Deploy a view to local, devnet, or mainnet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			viewstore := mustGetContextViewStore(cmd)
			schemastore, err := fileschema.NewFileSchemaStore()
			if err != nil {
				return err
			}

			viewName := args[0]

			switch target {
			case "local":
				return service.StartLocalNodeAndDeployView(viewName, viewstore, schemastore)
			case "devnet":
				wallet, err := service.LoadWallet()
				if err != nil {
					return err
				}
				return service.StartLocalNodeTestAndDeploy(viewName, viewstore, schemastore, wallet)
			case "mainnet":
				return fmt.Errorf("target '%s' not yet supported", target)
			default:
				return fmt.Errorf("invalid target '%s'. Must be one of: local, devnet, mainnet", target)
			}
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "Where to deploy the view: local, devnet, or mainnet (required)")

	cmd.MarkFlagRequired("target")
	return cmd
}
