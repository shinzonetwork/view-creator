package cli

import "github.com/spf13/cobra"

func NewViewCreatorCommand() *cobra.Command {
	view := MakeViewCommand()

	root := MakeRootCommand()
	root.AddCommand(
		view,
	)

	return root
}
