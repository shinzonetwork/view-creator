package cli

import (
	"github.com/spf13/cobra"
)

func MakeToolsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tools",
		Short: "Developer tools for working with the Viewkit",
	}

	cmd.AddCommand(MakeSchemaCommand())

	return cmd
}
