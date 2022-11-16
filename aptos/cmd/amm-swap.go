package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos"
)

func GetAmmSwapCmd() *cobra.Command {
	const longDescription = `Swap two coins. The output amount is decided by the input amount.

All the amounts are specified in atomic unit.

The coin x/y can either be fully qualified types, or a short hand name like USDC.
To see a list of all coins that are known, check "ls-known" command.

` + commonLongDescription

	cmd := &cobra.Command{
		Use:   "swap",
		Short: "swap coins on amm",
		Args:  cobra.NoArgs,
		Long:  longDescription,
	}

	var amountIn uint64
	var minAmountOut uint64

	coinX := ""
	coinY := ""

	args := NewSharedArgs()
	args.SetCmd(cmd)

	cmd.PersistentFlags().StringVarP(&coinX, "coin-in", "x", coinX, "input coin for the amm")
	cmd.MarkPersistentFlagRequired("coin-x")
	cmd.PersistentFlags().StringVarP(&coinY, "coin-out", "y", coinY, "output coin for the amm")
	cmd.MarkPersistentFlagRequired("coin-y")
	cmd.PersistentFlags().Uint64VarP(&amountIn, "amount-in", "i", amountIn, "input amount")
	cmd.MarkPersistentFlagRequired("in")
	cmd.PersistentFlags().Uint64VarP(&minAmountOut, "min-amount-out", "o", minAmountOut, "min amount out")

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

		client := getOrPanic(aptos.NewClient(args.network, args.endpoint))

		inCoin := getOrPanic(parseCoinType(args.network, coinX))
		outCoin := getOrPanic(parseCoinType(args.network, coinY))

		tx := auxConfig.Amm_SwapExactCoinForCoin(account.Address, inCoin, outCoin, amountIn, minAmountOut, aptos.TransactionOption_MaxGasAmount(args.maxGasAmount))

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
	return cmd
}
