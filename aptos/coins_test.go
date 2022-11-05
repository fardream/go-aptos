package aptos_test

import (
	"context"
	"testing"

	"github.com/fardream/go-aptos/aptos"
)

func TestClientGetCoinBalance(t *testing.T) {
	client := aptos.MustNewClient(aptos.Testnet, "")

	balance, err := client.GetCoinBalance(context.Background(), aptos.MustGetAuxClientConfig(aptos.Testnet).Address, &aptos.AptosCoin)
	if err != nil {
		t.Fatalf("failed to get balance: %v", err)
	}

	if balance <= 0 {
		t.Fatalf("balance is zero!")
	}
}
