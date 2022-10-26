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

func GetCreatePoolCmd() *cobra.Command {
	const longDescription = `Create a constant product AMM on https://aux.exchange

The coin x/y can either be fully qualified types, or a short hand name like USDC.
To see a list of all coins that are known, check "ls-known" command.

If no fee bps is specified, it will be set to 30.

` + commonLongDescription
	cmd := &cobra.Command{
		Use:   "create-pool",
		Short: "create a pool",
		Args:  cobra.NoArgs,
		Long:  longDescription,
	}

	coinX := ""
	coinY := ""

	var feeBps uint64 = 30

	args := NewSharedArgs()
	args.SetCmd(cmd)

	cmd.PersistentFlags().StringVarP(&coinX, "coin-x", "x", coinX, "x coin for the amm")
	cmd.MarkPersistentFlagRequired("coin-x")
	cmd.PersistentFlags().StringVarP(&coinY, "coin-y", "y", coinY, "x coin for the amm")
	cmd.MarkPersistentFlagRequired("coin-y")
	cmd.PersistentFlags().Uint64VarP(&feeBps, "fee-bps", "f", feeBps, "fee in bps")

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

		baseCoin := getOrPanic(parseCoinType(args.network, coinX))
		quoteCoin := getOrPanic(parseCoinType(args.network, coinY))

		tx := auxConfig.BuildCreatePool(account.Address, baseCoin, quoteCoin, feeBps, aptos.TransactionOption_MaxGasAmount(args.maxGasAmount))

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
