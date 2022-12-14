package aptos_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-aptos/aptos"
)

func ExampleAuxClientConfig_Router3Pool_SwapExactCoinForCoin() {
	client := aptos.MustNewClient(aptos.Devnet, "")

	userHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	configFileBytes, err := os.ReadFile(path.Join(userHome, ".aptos", "config.yaml"))
	if err != nil {
		panic(err)
	}

	configFile, err := aptos.ParseAptosConfigFile(configFileBytes)
	if err != nil {
		panic(err)
	}

	localProfile, ok := configFile.Profiles[string(aptos.Devnet)]
	if !ok {
		panic(fmt.Errorf("profile %s is not in config file", aptos.Devnet))
	}

	localAccount, err := localProfile.GetLocalAccount()
	if err != nil {
		panic(err)
	}

	auxConfig, err := aptos.GetAuxClientConfig(aptos.Devnet)
	if err != nil {
		panic(err)
	}

	usdc, _ := aptos.GetAuxFakeCoinCoinType(auxConfig.Address, aptos.AuxFakeCoin_USDC)
	usdt, _ := aptos.GetAuxFakeCoinCoinType(auxConfig.Address, aptos.AuxFakeCoin_USDT)
	usdcd8, _ := aptos.GetAuxFakeCoinCoinType(auxConfig.Address, aptos.AuxFakeCoin_USDCD8)

	auxCilent := aptos.NewAuxClient(client, auxConfig, localAccount)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*120)
	defer cancel()

	if _, err := auxCilent.GetStable3Pool(ctx, usdc, usdt, usdcd8, 0); err != nil {
		creatPoolTx := auxConfig.Stable3Pool_CreatePool(localAccount.Address, usdc, usdt, usdcd8, 15_000_000, 85)
		if err := client.FillTransactionData(context.Background(), creatPoolTx, false); err != nil {
			panic(err)
		}

		if txInfo, err := client.SignSubmitTransactionWait(ctx, localAccount, creatPoolTx, false); err != nil {
			spew.Fdump(os.Stderr, err)
		} else {
			spew.Fdump(os.Stderr, txInfo)
		}
	}

	for _, coin := range []aptos.AuxFakeCoin{aptos.AuxFakeCoin_USDC, aptos.AuxFakeCoin_USDT, aptos.AuxFakeCoin_USDCD8} {
		tx := auxConfig.FakeCoin_RegisterAndMint(localAccount.Address, coin, 1_000_000_000_0000)
		if err := client.FillTransactionData(ctx, tx, false); err != nil {
			panic(err)
		}
		if _, err := client.SignSubmitTransactionWait(ctx, localAccount, tx, false); err != nil {
			panic(err)
		}
	}

	addLiquidityTx := auxConfig.Router3Pool_AddLiquidity(localAccount.Address, usdc, 1_000_000_0000, usdt, 1_000_000_0000, usdcd8, 100_000_000_0000, 0)
	if err := client.FillTransactionData(ctx, addLiquidityTx, false); err != nil {
		panic(err)
	}
	if txInfo, err := client.SignSubmitTransactionWait(ctx, localAccount, addLiquidityTx, false); err != nil {
		panic(err)
	} else {
		spew.Fdump(os.Stderr, txInfo)
	}

	swapTx := auxConfig.Router3Pool_SwapExactCoinForCoin(localAccount.Address, usdc, 1_000, usdt, 1_000, usdcd8, 0, 2, 190_000, aptos.TransactionOption_MaxGasAmount(100_000))
	if err := client.FillTransactionData(ctx, swapTx, false); err != nil {
		panic(err)
	}
	if txInfo, err := client.SignSubmitTransactionWait(ctx, localAccount, swapTx, false); err != nil {
		spew.Fdump(os.Stderr, err)
	} else {
		spew.Fdump(os.Stderr, txInfo)
	}

	// Output:
}
