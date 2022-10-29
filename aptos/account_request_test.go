package aptos_test

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-aptos/aptos"
)

// Example: get the modules of [aux.exchange] from testnet
//
// [aux.exchange]: https://aux.exchange
func ExampleClient_GetAccountModules() {
	client := aptos.MustNewClient(aptos.Testnet, "")
	auxConfig, _ := aptos.GetAuxClientConfig(aptos.Testnet)

	modules, err := client.GetAccountModules(context.Background(), &aptos.GetAccountModulesRequest{
		Address: auxConfig.Address,
	})
	if err != nil {
		panic(err)
	}

	spew.Fdump(os.Stderr, modules.Parsed)
	if len(*modules.Parsed) >= 10 {
		fmt.Println("got all the modules")
	}

	// Output: got all the modules
}

// Example: get the clob_market of [aux.exchange] from testnet
//
// [aux.exchange]: https://aux.exchange
func ExampleClient_GetAccountModule() {
	client := aptos.MustNewClient(aptos.Testnet, "")
	auxConfig, _ := aptos.GetAuxClientConfig(aptos.Testnet)

	modules, err := client.GetAccountModule(context.Background(), &aptos.GetAccountModuleRequest{
		Address:    auxConfig.Address,
		ModuleName: "clob_market",
	})
	if err != nil {
		panic(err)
	}

	spew.Fdump(os.Stderr, modules.Parsed)
	fmt.Println("got clob_market")

	// Output: got clob_market
}
