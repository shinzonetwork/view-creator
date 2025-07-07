package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MakeViewRollbackCommand() *cobra.Command {
	var targetVersion int

	cmd := &cobra.Command{
		Use:   "rollback <viewName>",
		Short: "Rollback a view to a specific version (or previous if not specified)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			viewName := args[0]
			store := mustGetContextViewStore(cmd)

			currentView, err := store.Load(viewName)
			if err != nil {
				return fmt.Errorf("failed to load view '%s': %w", viewName, err)
			}

			if len(currentView.Metadata.Revisions) == 0 {
				return fmt.Errorf("no revisions found for view '%s'", viewName)
			}

			if targetVersion < 0 {
				targetVersion = currentView.Metadata.Revisions[len(currentView.Metadata.Revisions)-1].Version
			}

			rolledBackView, err := store.Rollback(viewName, targetVersion)
			if err != nil {
				return fmt.Errorf("rollback failed: %w", err)
			}

			printViewPretty(cmd, rolledBackView, false, false)
			return nil
		},
	}

	cmd.Flags().IntVar(&targetVersion, "version", -1, "Target version to rollback to")
	return cmd
}
