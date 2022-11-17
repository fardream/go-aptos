package aptos_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fardream/go-aptos/aptos"
	"github.com/google/go-cmp/cmp"
)

// How to get account resource
func ExampleClient_GetAccountResource() {
	client := aptos.MustNewClient(aptos.Testnet, "")
	// account resource is identified by a type.
	// AptosCoin is type
	aptosCoin, _ := aptos.NewMoveStructTag(aptos.AptosStdAddress, "aptos_coin", "AptosCoin", nil)
	// The coin value of an account is stored in a coin store
	aptosCoinStore, _ := aptos.NewMoveStructTag(aptos.AptosStdAddress, "coin", "CoinStore", []*aptos.MoveStructTag{aptosCoin})

	// let's check the coin balance of our deployer
	auxConfig, _ := aptos.GetAuxClientConfig(aptos.Testnet)

	// getting the coin store will get us the coin balance
	resp, err := client.GetAccountResource(context.Background(), &aptos.GetAccountResourceRequest{
		Address: auxConfig.Address,
		Type:    aptosCoinStore,
	})
	if err != nil {
		panic(err)
	}

	// check we have the correct type
	if !cmp.Equal(resp.Parsed.Type, aptosCoinStore) {
		panic(fmt.Errorf("differenet types: %s - %s", resp.Parsed.Type.String(), aptosCoinStore.String()))
	}

	// unfortuantely, we still need to parse the store into a proper golang object
	var coinStore aptos.CoinStore
	if err := json.Unmarshal(resp.Parsed.Data, &coinStore); err != nil {
		panic(err)
	}

	if coinStore.Coin.Value == 0 {
		panic("empty store")
	}

	fmt.Fprintln(os.Stderr, coinStore.Coin.Value)
	fmt.Println("we have money")

	// Output:we have money
}

// Listing Account Resources
func ExampleClient_GetAccountResources() {
	client := aptos.MustNewClient(aptos.Devnet, "")
	// let's check the coin balance of our deployer
	auxConfig, _ := aptos.GetAuxClientConfig(aptos.Devnet)

	// getting the coin store will get us the coin balance
	resp, err := client.GetAccountResources(context.Background(), &aptos.GetAccountResourcesRequest{
		Address: auxConfig.Deployer,
	})
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, string(resp.RawData))
	fmt.Println("done")
	// Output:done
}
