package aptos_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-aptos/aptos/known"
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

func ExampleClient_LoadEvents_auxFillEvents() {
	client := aptos.MustNewClient(aptos.Mainnet, "")
	auxConfig, _ := aptos.GetAuxClientConfig(aptos.Mainnet)

	auxClient := aptos.NewAuxClient(client, auxConfig, nil)

	usdc := known.GetCoinInfoBySymbol(aptos.Mainnet, "USDC").TokenType.Type
	apt := known.GetCoinInfoBySymbol(aptos.Mainnet, "APT").TokenType.Type

	// get the market information
	market, err := auxClient.GetClobMarket(context.Background(), apt, usdc, 0)
	if err != nil {
		panic(err)
	}

	// get the creation number of the event handler for fills, can do the same for OrderPlaced and OrderCancel
	creationNumber := uint64(market.FillEvents.GUID.Id.CreationNumber)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	rawevents, err := client.LoadEvents(ctx, auxConfig.Address, creationNumber, 0, 312, 50)
	if err != nil {
		spew.Dump(err)
		panic(err)
	}

	// converts the raw events into proper fill events
	events := aptos.FilterAuxClobMarketOrderEvent(rawevents, auxConfig.Address, false, true)

	spew.Sdump(os.Stderr, events)
	fmt.Println(len(events))

	// Output: 312
}
