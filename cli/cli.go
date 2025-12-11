package cli

import "github.com/spf13/cobra"

func NewViewCreatorCommand() *cobra.Command {
	view := MakeViewCommand()
	tool := MakeToolsCommand()
	wallet := MakeWalletCommand()

	root := MakeRootCommand()
	root.AddCommand(
		view,
		tool,
		wallet,
	)

	return root
}
