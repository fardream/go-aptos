// contains commands that are supported by the aptos-aux cli.
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
		GetCreatePoolCmd(),
		GetUpdatePoolFeeCmd(),
		GetCreateMarketCmd(),
		GetListKnownCmd(),
		GetMintFakeCoinCmd(),
		GetAmmAddLiquidityCmd(),
		GetAmmSwapCmd(),
		GetAmmRemoveLiquidityCmd(),
	)

	return cmd
}
