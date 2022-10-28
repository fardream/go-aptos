package aptos_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-aptos/aptos"
)

const network = aptos.Devnet

func ExampleClient_SubmitTransaction() {
	restUrl, faucetUrl, _ := aptos.GetDefaultEndpoint(network)

	account, _ := aptos.NewLocalAccountWithRandomKey()

	var err error

	fmt.Fprintf(os.Stderr, "private key: 0x%s\n", hex.EncodeToString(account.PrivateKey.Seed()))
	fmt.Fprintf(os.Stderr, "address: %s\n", account.Address.String())
	auxConfig, _ := aptos.GetAuxClientConfig(network)

	faucetTxes, err := aptos.RequestFromFaucet(context.Background(), faucetUrl, &account.Address, 1000000000)
	if err != nil {
		panic(err)
	}

	client := aptos.MustNewClient(network, restUrl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	for _, faucetTx := range faucetTxes {
		txType, err := client.WaitForTransaction(ctx, faucetTx)
		if err != nil {
			spew.Dump(err)
		}
		fmt.Fprintf(os.Stderr, "fauct tx type: %s\n", txType.Type)
	}

	tx := aptos.Transaction{
		Sender:                  account.Address,
		ExpirationTimestampSecs: aptos.JsonUint64(time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC).Unix()),
		Payload: aptos.NewEntryFunctionPayload(
			aptos.MustNewMoveFunctionTag(auxConfig.Address, "fake_coin", "mint"),
			[]*aptos.MoveTypeTag{aptos.MustNewMoveTypeTag(auxConfig.Address, "fake_coin", "USDC", nil)},
			[]aptos.EntryFunctionArg{aptos.JsonUint64(1000000000000)}),
	}

	if err := client.FillTransactionData(context.Background(), &tx, false); err != nil {
		panic(err)
	}

	encode, err := client.EncodeSubmission(context.Background(),
		&aptos.EncodeSubmissionRequest{
			Transaction: &tx,
		},
	)
	if err != nil {
		panic(err)
	}

	data, err := hex.DecodeString(strings.TrimPrefix(string(*encode.Parsed), "0x"))
	if err != nil {
		panic(err)
	}

	signature, _ := account.SignRawData(data)

	request := &aptos.SubmitTransactionRequest{
		Transaction: &tx,
		Signature:   *signature,
	}
	body, _ := request.Body()
	fmt.Fprintln(os.Stderr, string(body))
	resp, err := client.SubmitTransaction(context.Background(), request)
	if err != nil {
		panic(err)
	} else {
		spew.Fdump(os.Stderr, resp)
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		status, err := client.WaitForTransaction(ctx, resp.Parsed.Hash)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(status)
		}
	}
}
