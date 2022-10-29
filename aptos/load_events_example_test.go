package aptos_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-aptos/aptos"
)

func ExampleClient_LoadEvents() {
	client := aptos.MustNewClient(aptos.Mainnet, "")
	auxConfig, _ := aptos.GetAuxClientConfig(aptos.Mainnet)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	// pool: Aptos/USDC 0xbd35135844473187163ca197ca93b2ab014370587bb0ed3befff9e902d6bb541::amm::Pool<0x1::aptos_coin::AptosCoin, 0x5e156f1207d0ebfa19a9eeff00d62a282278fb8719f4fab3a586a0a2c0fffbea::coin::T>
	events, err := client.LoadEvents(ctx, auxConfig.Address, 72, 93, 405, 25)
	if err != nil {
		spew.Dump(err)
		panic(err)
	}

	fmt.Println(len(events))

	spew.Sdump(os.Stderr, events)
	// Output: 312
}
