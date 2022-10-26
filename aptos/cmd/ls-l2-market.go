package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos"
)

func GetListL2MarketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls-l2",
		Short: "list level 2 market data",
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

		client := aptos.NewClient(sharedArgs.endpoint)

		baseCoin := getOrPanic(parseCoinType(sharedArgs.network, sharedArgs.baseCoinStr))
		quoteCoin := getOrPanic(parseCoinType(sharedArgs.network, sharedArgs.quoteCoinStr))

		tx := auxConfig.BuildLoadMarketIntoEvent(baseCoin, quoteCoin, aptos.TransactionOption_MaxGasAmount(sharedArgs.maxGasAmount))

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

		var L2 aptos.AuxLevel2Event
		orPanic(json.Unmarshal(*events[0].Data, &L2))

		spew.Dump(L2)
	}

	return cmd
}
