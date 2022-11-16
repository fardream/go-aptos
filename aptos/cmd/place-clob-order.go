package cmd

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-bcs/bcs"
)

func GetPlaceClobOrderCmd() *cobra.Command {
	const longDescription = `
The base/quote coins can either be fully qualified types, or a short hand name like USDC.
To see a list of all coins that are known, check "ls-known-coins" command.

Price is quoted per unit of base in quote quantities. For example, if BTC has decimal of 8, and USDC has decimal of 6.
To buy 0.5 BTC at a price of 19000 USDC, price should be 19,000,000,000, and the quantity should be 50,000,000.

` + commonLongDescription
	cmd := &cobra.Command{
		Use:   "place-clob-order",
		Short: "place an order on clob",
		Args:  cobra.NoArgs,
		Long:  "Place a limit order on central limit order book of https://aux.exchange\n" + longDescription,
	}

	args := NewSharedArgsWithBaseQuoteCoins()
	args.SetCmd(cmd)
	var limitPrice uint64 = 0
	var quantity uint64 = 0

	cmd.PersistentFlags().Uint64VarP(&limitPrice, "limit-price", "p", limitPrice, "limit price")
	cmd.MarkPersistentFlagRequired("limit-price")
	cmd.PersistentFlags().Uint64VarP(&quantity, "quantity", "t", quantity, "quantity")
	cmd.MarkPersistentFlagRequired("quantity")

	buyCmd := &cobra.Command{
		Use:   "buy",
		Short: "buy on aux",
		Args:  cobra.NoArgs,
		Long:  "Buy on central limit order book of https://aux.exchange\n" + longDescription,
	}

	sellCmd := &cobra.Command{
		Use:   "sell",
		Short: "sell on aux",
		Long:  "Sell on central limit order book of https://aux.exchange\n" + longDescription,
		Args:  cobra.NoArgs,
	}

	buildCmd := func(isBuy bool) func(*cobra.Command, []string) {
		isBid := isBuy
		return func(c *cobra.Command, s []string) {
			args.UpdateProfileForCmd(c)
			configFile, _ := getConfigFileLocation()
			configs := getOrPanic(aptos.ParseAptosConfigFile(getOrPanic(os.ReadFile(configFile))))
			if configs.Profiles == nil {
				orPanic(fmt.Errorf("empty configuration at %s", configFile))
			}
			config, ok := configs.Profiles[args.profile]
			if !ok {
				orPanic(fmt.Errorf("cannot find profile %s in config file %s", args.profile, configFile))
			}

			if args.endpoint == "" && config.RestUrl != "" {
				args.endpoint = config.RestUrl
			}
			if args.endpoint == "" {
				var err error
				args.endpoint, _, err = aptos.GetDefaultEndpoint(args.network)
				orPanic(err)
			}
			account := getOrPanic(config.GetLocalAccount())

			auxConfig := getOrPanic(aptos.GetAuxClientConfig(args.network))

			client := aptos.MustNewClient(args.network, args.endpoint)

			baseCoin := getOrPanic(parseCoinType(args.network, args.baseCoinStr))
			quoteCoin := getOrPanic(parseCoinType(args.network, args.quoteCoinStr))

			tx := auxConfig.ClobMarket_PlaceOrder(
				account.Address,
				isBid,
				baseCoin,
				quoteCoin,
				limitPrice,
				quantity,
				0,
				bcs.Uint128{},
				aptos.AuxClobMarketOrderType_Limit,
				0,
				false,
				math.MaxUint64,
				aptos.AuxClobMarketSelfTradeType_CancelBoth,
				aptos.TransactionOption_MaxGasAmount(args.maxGasAmount))
			orPanic(client.FillTransactionData(context.Background(), tx, false))

			if args.simulate {
				resp := getOrPanic(
					client.SimulateTransaction(context.Background(), &aptos.SimulateTransactionRequest{
						Transaction: tx,
						Signature:   aptos.NewSingleSignatureForSimulation(&account.PublicKey),
					}),
				)
				fmt.Println(string(resp.RawData))
				return
			}

			spew.Dump(getOrPanic(client.SignSubmitTransactionWait(context.Background(), account, tx, false)))
		}
	}

	buyCmd.Run = buildCmd(true)
	sellCmd.Run = buildCmd(false)

	cmd.AddCommand(buyCmd, sellCmd)

	return cmd
}
