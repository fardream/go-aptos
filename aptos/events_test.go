package aptos_test

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/fardream/go-aptos/aptos"
)

// Example: getting clob market events by their creation numbers
func ExampleClient_GetEventsByCreationNumber() {
	client := aptos.MustNewClient(aptos.Devnet, "")

	auxConfig, _ := aptos.GetAuxClientConfig(aptos.Devnet)

	fakeEth, _ := aptos.GetAuxFakeCoinCoinType(auxConfig.Address, aptos.AuxFakeCoin_ETH)
	fakeUsdc, _ := aptos.GetAuxFakeCoinCoinType(auxConfig.Address, aptos.AuxFakeCoin_USDC)
	marketType := aptos.MustNewMoveTypeTag(auxConfig.Address, "clob_market", "Market", []*aptos.MoveTypeTag{fakeEth, fakeUsdc})

	market, err := aptos.GetAccountResourceWithType[aptos.AuxClobMarket](context.Background(), client, auxConfig.Address, marketType, 0)
	if err != nil {
		panic(err)
	}

	if market.PlacedEvents == nil {
		panic(fmt.Errorf("placed_events is nil"))
	}

	creation_number := market.PlacedEvents.GUID.Id.CreationNumber

	resp, err := client.GetEventsByCreationNumber(context.Background(), &aptos.GetEventsByCreationNumberRequest{
		CreationNumber: creation_number,
		Address:        auxConfig.Address,
	})
	if err != nil {
		panic(err)
	}

	spew.Fdump(os.Stderr, resp)

	fmt.Println("done")

	// Output:done
}

func ExampleClient_GetEventsByEventHandler() {
	client := aptos.MustNewClient(aptos.Devnet, "")

	auxConfig, _ := aptos.GetAuxClientConfig(aptos.Devnet)

	fakeEth, _ := aptos.GetAuxFakeCoinCoinType(auxConfig.Address, aptos.AuxFakeCoin_ETH)
	fakeUsdc, _ := aptos.GetAuxFakeCoinCoinType(auxConfig.Address, aptos.AuxFakeCoin_USDC)
	marketType := aptos.MustNewMoveTypeTag(auxConfig.Address, "clob_market", "Market", []*aptos.MoveTypeTag{fakeEth, fakeUsdc})

	resp, err := client.GetEventsByEventHandler(context.Background(), &aptos.GetEventsByEventHandlerRequest{
		EventHandler: marketType,
		Address:      auxConfig.Address,
		FieldName:    "placed_events",
	})
	if err != nil {
		panic(err)
	}

	spew.Fdump(os.Stderr, resp)

	fmt.Println("done")

	// Output:done
}
