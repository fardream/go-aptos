package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-aptos/aptos"
	"github.com/spf13/cobra"
)

func GetMintFakeCoinCmd() *cobra.Command {
	const longDescription = `Mint fake coin to test trades on https://aux.exchagne.

Only support devnet and testnet.

Following fake coins are available:
- ETH
- BTC
- USDC
- USDT
- SOL
- AUX

Minting those coins doesnt need authority, anyone can mint as much as they want.

` + commonLongDescription

	cmd := &cobra.Command{
		Use:   "mint-fake-coin",
		Short: "mint fake coin to use on testnet/devnet",
		Long:  longDescription,
	}

	args := NewSharedArgs()
	args.SetCmd(cmd)

	var fakeCoin string
	var amount uint64

	cmd.PersistentFlags().StringVarP(&fakeCoin, "coin", "x", fakeCoin, "fake coin to mint")
	cmd.MarkPersistentFlagRequired("coin")

	cmd.PersistentFlags().Uint64VarP(&amount, "amount", "a", amount, "amount to mint")
	cmd.MarkPersistentFlagRequired("amount")

	cmd.Run = func(*cobra.Command, []string) {
		args.UpdateProfileForCmd(cmd)

		if args.network != aptos.Devnet && args.network != aptos.Testnet {
			orPanic(fmt.Errorf("unsupported network: %s", args.network))
		}

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

		coin := getOrPanic(aptos.ParseAuxFakeCoin(fakeCoin))

		tx := auxConfig.FakeCoin_RegisterAndMint(account.Address, coin, amount, aptos.TransactionOption_MaxGasAmount(args.maxGasAmount))

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
