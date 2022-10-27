package aptos_test

import (
	"context"
	"testing"

	"github.com/fardream/go-aptos/aptos"
)

func TestClientGetCoinBalance(t *testing.T) {
	client := aptos.NewClient("https://fullnode.devnet.aptoslabs.com/v1")

	address := aptos.MustParseAddress("0xea383dc2819210e6e427e66b2b6aa064435bf672dc4bdc55018049f0c361d01a")

	balance, err := client.GetCoinBalance(context.Background(), address, &aptos.AptosCoin)
	if err != nil {
		t.Fatalf("failed to get balance: %v", err)
	}

	if balance <= 0 {
		t.Fatalf("balance is zero!")
	}
}
