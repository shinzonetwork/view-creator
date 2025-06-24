package cli

import (
	"github.com/spf13/cobra"
)

func MakeRootCommand() *cobra.Command {
	var cmd = &cobra.Command{
		SilenceUsage: true,
		Use:          "viewkit",
		Short:        "A CLI tool to manage shinzo views",
		Long:         "Viewkit helps you initialize, manage, and publish Shinzo views through a simple CLI interface.",
	}

	return cmd
}
