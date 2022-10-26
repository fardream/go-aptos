package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos"
)

func GetCreateMarketCmd() *cobra.Command {
	const longDescription = `Create a central limit order book on https://aux.exchange.

The base/quote coins can either be fully qualified types, or a short hand name like USDC.
To see a list of all coins that are known, check "ls-known" command.

Tick size is specified in decimals of the quote coin. Lot size is specified in decimals of the base coin.
Price of the market is specified as amount of quote coin per 1 unit of base coin.

For example, on devnet test market, BTC has a decimal of 8 and USDC has a decimal of 6.
- lot size is 0.01 BTC, and should be specified 1000000 (1 million).
- tick size is 0.01 USDC, and should be specified 10000 (10 thousand)

To buy 0.5 BTC at a price of 19000, price should be 19,000,000,000, and the quantity should be 50,000,000.

` + commonLongDescription

	cmd := &cobra.Command{
		Use:   "create-market",
		Short: "create new market on aux exchange",
		Long:  longDescription,
		Args:  cobra.NoArgs,
	}

	args := NewSharedArgsWithBaseQuoteCoins()
	args.SetCmd(cmd)

	var lotSize, tickSize uint64
	cmd.PersistentFlags().Uint64Var(&lotSize, "lot-size", lotSize, "lot size")
	cmd.MarkPersistentFlagRequired("lot-size")
	cmd.PersistentFlags().Uint64Var(&tickSize, "tick-size", tickSize, "tick size")
	cmd.MarkPersistentFlagRequired("tick-size")

	cmd.Run = func(*cobra.Command, []string) {
		args.UpdateProfileForCmd(cmd)
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

		client := aptos.NewClient(args.endpoint)

		baseCoin := getOrPanic(parseCoinType(args.network, args.baseCoinStr))
		quoteCoin := getOrPanic(parseCoinType(args.network, args.quoteCoinStr))

		tx := auxConfig.BuildCreateMarket(account.Address, baseCoin, quoteCoin, lotSize, tickSize, aptos.TransactionOption_MaxGasAmount(args.maxGasAmount))
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

		chainId := getOrPanic(client.GetChainId(context.Background()))
		signingData := aptos.EncodeTransaction(tx, chainId)
		signature := getOrPanic(account.Sign(signingData))

		spew.Dump(getOrPanic(client.SubmitTransaction(context.Background(), &aptos.SubmitTransactionRequest{
			Transaction: tx,
			Signature:   *signature,
		})))
	}

	return cmd
}
