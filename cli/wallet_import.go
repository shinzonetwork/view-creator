package cli

import (
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeWalletImportCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "import [mnemonic]",
		Short: "Import wallet from mnemonic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := service.ImportMnemonic(args[0]); err != nil {
				return err
			}

			wallet, err := service.LoadWallet()
			if err != nil {
				return err
			}

			fmt.Println("✅ Wallet imported successfully")
			fmt.Println("Address: ", wallet.Address)
			return nil
		},
	}
}
