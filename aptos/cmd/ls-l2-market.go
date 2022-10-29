package cmd

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos"
)

func GetListL2MarketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls-l2",
		Short: "list level 2 market data",
		Long:  "List level 2 market data (price and quantities).\n" + commonLongDescription,
		Args:  cobra.NoArgs,
	}

	sharedArgs := NewSharedArgsWithBaseQuoteCoins()

	sharedArgs.SetCmd(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if sharedArgs.endpoint == "" {
			var err error
			sharedArgs.endpoint, _, err = aptos.GetDefaultEndpoint(sharedArgs.network)
			orPanic(err)
		}

		auxConfig := getOrPanic(aptos.GetAuxClientConfig(sharedArgs.network))

		client := aptos.MustNewClient(sharedArgs.network, sharedArgs.endpoint)

		baseCoin := getOrPanic(parseCoinType(sharedArgs.network, sharedArgs.baseCoinStr))
		quoteCoin := getOrPanic(parseCoinType(sharedArgs.network, sharedArgs.quoteCoinStr))

		L2 := getOrPanic(aptos.NewAuxClient(client, auxConfig, nil).ListLevel2(context.Background(), baseCoin, quoteCoin, aptos.TransactionOption_MaxGasAmount(sharedArgs.maxGasAmount)))

		spew.Dump(L2)
	}

	return cmd
}
