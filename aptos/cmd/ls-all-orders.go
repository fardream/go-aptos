package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos"
)

func GetListAllOrdersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls-all-orders",
		Short: "list all orders for clob",
		Long:  "List all orders for the clob.\nThe base/quote coins can either be fully qualified types, or a short hand name like USDC.\nTo see a list of all coins that are known, check ls-known command." + commonLongDescription,
		Args:  cobra.NoArgs,
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

		client := aptos.NewClient(args.endpoint)

		baseCoin := getOrPanic(parseCoinType(args.network, args.baseCoinStr))
		quoteCoin := getOrPanic(parseCoinType(args.network, args.quoteCoinStr))

		tx := auxConfig.ClobMarket_LoadAllOrdersIntoEvent(baseCoin, quoteCoin, aptos.TransactionOption_MaxGasAmount(args.maxGasAmount))

		orPanic(client.FillTransactionData(context.Background(), tx, false))

		resp := getOrPanic(
			client.SimulateTransaction(context.Background(), &aptos.SimulateTransactionRequest{
				Transaction: tx,
				Signature:   aptos.NewSingleSignatureForSimulation(&auxConfig.DataFeedPublicKey),
			}),
		)
		parsed := *resp.Parsed
		if len(parsed) != 1 {
			orPanic(fmt.Errorf("more than one transactions in response: %s", string(resp.RawData)))
		}
		if !parsed[0].Success {
			orPanic(fmt.Errorf("simulation failed: %s, resp is %s", parsed[0].VmStatus, string(resp.RawData)))
		}
		events := parsed[0].Events
		if len(events) != 1 {
			orPanic(fmt.Errorf("there should only be one events: %#v", events))
		}

		var L2 aptos.AuxAllOrdersEvent
		orPanic(json.Unmarshal(*events[0].Data, &L2))

		spew.Dump(L2)
	}

	return cmd
}
