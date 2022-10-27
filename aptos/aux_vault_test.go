package aptos_test

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-aptos/aptos"
)

func TestClient_GetAuxCoinBalance(t *testing.T) {
	client := aptos.NewClient(devnetUrl)
	eth, _ := aptos.GetAuxFakeCoinCoinType(devnetConfig.Address, aptos.AuxFakeCoin_ETH)
	cb, err := client.GetAuxCoinBalance(context.Background(), devnetConfig, trader.Address, eth)
	if err != nil {
		t.Fatalf("failed to get balance: %v", err)
	}

	spew.Dump(cb)
}
