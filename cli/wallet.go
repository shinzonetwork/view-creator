package cli

import (
	"github.com/spf13/cobra"
)

func MakeWalletCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wallet",
		Short: "Manage your Shinzo wallet",
	}

	cmd.AddCommand(MakeWalletGenerateCmd())
	cmd.AddCommand(MakeWalletImportCommand())
	cmd.AddCommand(MakeWalletShowCommand())

	return cmd
}
