package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos"
)

func GetCreateMarketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-market",
		Short: "create new market on aux exchange",
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

		resp := getOrPanic(client.EncodeSubmission(context.Background(), &aptos.EncodeSubmissionRequest{
			Transaction: tx,
		}))

		signature := getOrPanic(account.Sign(getOrPanic(hex.DecodeString(strings.TrimPrefix(string(*resp.Parsed), "0x")))))

		spew.Dump(getOrPanic(client.SubmitTransaction(context.Background(), &aptos.SubmitTransactionRequest{
			Transaction: tx,
			Signature:   *signature,
		})))
	}

	return cmd
}
