package cli

import (
	"github.com/shinzonetwork/view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeViewInitCommand() *cobra.Command {
	var verbose bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "init [name]",
		Short: "Initialize a new view",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			storeImpl := mustGetContextStore(cmd)
			view, err := service.InitView(args[0], storeImpl)
			if err != nil {
				return err
			}

			printViewPretty(cmd, view, verbose, jsonOutput)
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output the view in raw JSON format")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Show full output including revision history")

	return cmd
}
