package cli

import "github.com/spf13/cobra"

func NewViewCreatorCommand() *cobra.Command {
	view := MakeViewCommand()
	tool := MakeToolsCommand()

	root := MakeRootCommand()
	root.AddCommand(
		view,
		tool,
	)

	return root
}
