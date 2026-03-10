package cli

import (
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeRemoveLensCommand(viewName *string) *cobra.Command {
	var label string

	cmd := &cobra.Command{
		Use:   "lens",
		Short: "Remove a lens from the view",
		RunE: func(cmd *cobra.Command, args []string) error {
			if label == "" {
				return fmt.Errorf("--label is required")
			}

			store := mustGetContextViewStore(cmd)

			view, err := service.RemoveLens(*viewName, label, store)
			if err != nil {
				return err
			}

			printViewPretty(cmd, view, false, false)
			return nil
		},
	}

	cmd.Flags().StringVar(&label, "label", "", "Label of the lens to remove")

	return cmd
}
