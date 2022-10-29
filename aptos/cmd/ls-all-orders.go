package cmd

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos"
)

func GetListAllOrdersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls-all-orders",
		Short: "list all orders for clob",
		Long: `List all orders for the clob.
The base/quote coins can either be fully qualified types, or a short hand name like USDC.
To see a list of all coins that are known, check ls-known command.` + commonLongDescription,
		Args: cobra.NoArgs,
	}

	args := NewSharedArgsWithBaseQuoteCoins()

	args.SetCmd(cmd)

	cmd.Run = func(*cobra.Command, []string) {
		if args.endpoint == "" {
			var err error
			args.endpoint, _, err = aptos.GetDefaultEndpoint(args.network)
			orPanic(err)
		}

		auxConfig := getOrPanic(aptos.GetAuxClientConfig(args.network))

		client := aptos.MustNewClient(args.network, args.endpoint)

		baseCoin := getOrPanic(parseCoinType(args.network, args.baseCoinStr))
		quoteCoin := getOrPanic(parseCoinType(args.network, args.quoteCoinStr))

		allOrders := getOrPanic(aptos.NewAuxClient(client, auxConfig, nil).ListAllOrders(context.Background(), baseCoin, quoteCoin, aptos.TransactionOption_MaxGasAmount(args.maxGasAmount)))

		spew.Dump(allOrders)
	}

	return cmd
}
