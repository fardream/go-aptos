package aptos_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-aptos/aptos"
)

const network = aptos.Devnet

// Example Submitting a transaction to the devnet.
// But first, we need to create a new account, and obtain some gas from faucet.
func ExampleClient_SubmitTransaction() {
	// rest url and faucet url
	restUrl, faucetUrl, _ := aptos.GetDefaultEndpoint(network)

	// create new account, and fund it with faucet
	account, _ := aptos.NewLocalAccountWithRandomKey()

	fmt.Fprintf(os.Stderr, "private key: 0x%s\n", hex.EncodeToString(account.PrivateKey.Seed()))
	fmt.Fprintf(os.Stderr, "address: %s\n", account.Address.String())

	// aptos client
	client := aptos.MustNewClient(network, restUrl)

	// do faucet. note the faucet transaction is still inflight when returned, we need to check for the completion of the transaction.
	faucetTxes, err := aptos.RequestFromFaucet(context.Background(), faucetUrl, &account.Address, 1000000000)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	for _, faucetTx := range faucetTxes {
		txType, err := client.WaitForTransaction(ctx, faucetTx)
		if err != nil {
			spew.Dump(err)
		}
		fmt.Fprintf(os.Stderr, "fauct tx type: %s\n", txType.Type)
	}

	// construct a transaction from aux suite of transactions.
	auxConfig, _ := aptos.GetAuxClientConfig(network)

	tx := aptos.Transaction{
		Sender:                  account.Address,
		ExpirationTimestampSecs: aptos.JsonUint64(time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC).Unix()),
		Payload: aptos.NewEntryFunctionPayload(
			aptos.MustNewMoveFunctionTag(auxConfig.Address, "fake_coin", "mint"),
			[]*aptos.MoveStructTag{aptos.MustNewMoveStructTag(auxConfig.Address, "fake_coin", "USDC", nil)},
			[]*aptos.EntryFunctionArg{aptos.EntryFunctionArg_Uint64(1000000000000)}),
	}

	// fill the missing information (
	if err := client.FillTransactionData(context.Background(), &tx, false); err != nil {
		panic(err)
	}

	// sign the transaction.
	signature, err := account.Sign(&tx)
	if err != nil {
		panic(err)
	}

	// submit the transaction.
	request := &aptos.SubmitTransactionRequest{
		Transaction: &tx,
		Signature:   *signature,
	}
	// we can take a look at the json body.
	body, _ := request.Body()
	fmt.Fprintln(os.Stderr, string(body))

	// submit the transaction.
	resp, err := client.SubmitTransaction(context.Background(), request)
	if err != nil {
		panic(err)
	} else {
		spew.Fdump(os.Stderr, resp)
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		txInfo, err := client.WaitForTransaction(ctx, resp.Parsed.Hash)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(txInfo.Type)
		}
	}
}
