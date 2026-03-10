package cli

import (
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeWalletGenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate and save a new wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			mnemonic, address, err := service.GenerateWallet()
			if err != nil {
				return err
			}

			fmt.Println("✅ Wallet generated")
			fmt.Println("Mnemonic:", mnemonic)
			fmt.Println("Address:", address)
			return nil
		},
	}
}
