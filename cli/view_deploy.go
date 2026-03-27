package cli

import (
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/schema/store/fileschema"
	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeViewDeployCommand() *cobra.Command {
	var target string
	var rpc string
	var url string
	var debug bool

	cmd := &cobra.Command{
		Use:   "deploy <name>",
		Short: "Deploy a view to local, playground, devnet, or mainnet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			viewstore := mustGetContextViewStore(cmd)
			schemastore, err := fileschema.NewFileSchemaStore()
			if err != nil {
				return err
			}

			viewName := args[0]

			if (target == "devnet" || target == "mainnet") && rpc == "" {
				return fmt.Errorf("--rpc is required when --target is %s", target)
			}

			switch target {
			case "local":
				return service.StartLocalNodeAndDeployView(viewName, viewstore, schemastore, debug)

			case "playground":
				return service.DeployViewToPlayground(viewName, viewstore, schemastore, url)

			case "devnet":
				wallet, err := service.LoadWallet()
				if err != nil {
					return err
				}
				return service.StartLocalNodeTestAndDeploy(viewName, viewstore, schemastore, wallet, rpc, debug)

			case "mainnet":
				return fmt.Errorf("target '%s' not yet supported", target)

			default:
				return fmt.Errorf("invalid target '%s'. Must be one of: local, playground, devnet, mainnet", target)
			}
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "Where to deploy the view: local, playground, devnet, or mainnet (required)")
	cmd.Flags().StringVar(&rpc, "rpc", "", "RPC endpoint URL (required for devnet/mainnet)")
	cmd.Flags().StringVar(&url, "url", "http://127.0.0.1:9181", "DefraDB URL for playground target")
	cmd.Flags().BoolVar(&debug, "debug", false, "Enable debug output")

	cmd.MarkFlagRequired("target")
	return cmd
}
