package cmd

import "github.com/spf13/cobra"

func GetRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aptos-aux",
		Short: "aux exchange utility cli - aptos version",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		GetKeyAndMnemonicCmd(),
		GetCalculateResourceAddressCmd(),
		GetLaunchAptosNodeCmd(),
		GetListAccountCmd(),
		GetListMarketCmd(),
		GetListChainCmd(),
		GetListPoolCmd(),
		GetListL2MarketCmd(),
		GetListAllOrdersCmd(),
		GetPlaceClobOrderCmd(),
	)

	return cmd
}
