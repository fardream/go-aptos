package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos"
)

func GetUpdatePoolFeeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-pool-fee",
		Short: "update pool fee",
		Args:  cobra.NoArgs,
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
	cmd.MarkPersistentFlagRequired("fee-bps")

	cmd.Run = func(*cobra.Command, []string) {
		configFile, _ := getConfigFileLocation()
		configs := getOrPanic(aptos.ParseAptosConfigFile(getOrPanic(os.ReadFile(configFile))))
		if configs.Profiles == nil {
			orPanic(fmt.Errorf("empty configuration at %s", configFile))
		}
		if !cmd.PersistentFlags().Changed("profile") {
			args.profile = string(args.network)
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

		tx := auxConfig.Amm_UpdateFee(account.Address, baseCoin, quoteCoin, feeBps, aptos.TransactionOption_MaxGasAmount(args.maxGasAmount))

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
