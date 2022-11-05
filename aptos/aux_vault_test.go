package aptos_test

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/fardream/go-aptos/aptos"
)

func ExampleClient_GetAuxCoinBalance() {
	client := aptos.MustNewClient(aptos.Devnet, "")
	eth, _ := aptos.GetAuxFakeCoinCoinType(devnetConfig.Address, aptos.AuxFakeCoin_ETH)
	cb, err := client.GetAuxCoinBalance(context.Background(), devnetConfig, trader.Address, eth)
	if err != nil {
		panic(fmt.Errorf("failed to get balance: %v", err))
	}

	spew.Dump(cb)
}
