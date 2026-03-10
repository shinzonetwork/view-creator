package cli

import (
	"fmt"

	"github.com/shinzonetwork/shinzo-view-creator/core/service"
	"github.com/spf13/cobra"
)

func MakeWalletShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "inspect",
		Short: "Show the saved wallet address",
		RunE: func(cmd *cobra.Command, args []string) error {
			wallet, err := service.LoadWallet()
			if err != nil {
				return err
			}
			fmt.Println("Address: ", wallet.Address)
			return nil
		},
	}
}
