package aptos_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-aptos/aptos"
)

func ExampleClient_SubmitTransaction_multiple() {
	aptosClient := aptos.MustNewClient(aptos.Devnet, "")
	auxConfig := aptos.MustGetAuxClientConfig(aptos.Devnet)

	localAccount, _ := aptos.NewLocalAccountWithRandomKey()
	fmt.Fprintf(os.Stderr, "0x%s\n", hex.EncodeToString(localAccount.PrivateKey.Seed()))
	_, facuetUrl, _ := aptos.GetDefaultEndpoint(aptos.Devnet)

	txes, err := aptos.RequestFromFaucet(context.Background(), facuetUrl, &(localAccount.Address), 1000000000)
	if err != nil {
		panic(err)
	}

	for _, tx := range txes {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if tx, err := aptosClient.WaitForTransaction(ctx, tx); err != nil {
			panic(err)
		} else {
			if !tx.Success {
				panic(tx)
			}
		}
	}

	// a new account always starts with sequence number 0.
	registerTx := auxConfig.FakeCoin_Register(
		localAccount.Address,
		aptos.AuxFakeCoin_SOL,
		aptos.TransactionOption_SequenceNumber(0), // seqnum 0.
		aptos.TransactionOption_ExpireAfter(5*time.Minute))
	if err := aptosClient.FillTransactionData(context.Background(), registerTx, true); err != nil {
		panic(err)
	}
	// a second transaction that is guaranteed to be executed after the first.
	mintTx := auxConfig.FakeCoin_Mint(
		localAccount.Address,
		aptos.AuxFakeCoin_SOL,
		10000000000,
		aptos.TransactionOption_SequenceNumber(1), // seqnum 1
		aptos.TransactionOption_ExpireAfter(5*time.Minute))
	if err := aptosClient.FillTransactionData(context.Background(), mintTx, true); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	// send the min transaction first, then send the register transaction.
	// mint transaction will wait for register transaction because it has a later sequence number.
	for _, tx := range []*aptos.Transaction{mintTx, registerTx} {
		fmt.Fprintln(os.Stderr, "wait for 5 seconds")
		<-time.After(5 * time.Second)
		wg.Add(1)
		go func(tx *aptos.Transaction) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()

			txInfo, err := aptosClient.SignSubmitTransactionWait(ctx, localAccount, tx, false)
			if err != nil {
				spew.Fdump(os.Stderr, err)
			} else {
				fmt.Fprintf(os.Stderr, "function: %s - hash: %s- status: %s\n", txInfo.Payload.Function.Name, txInfo.Hash, txInfo.VmStatus)
			}
		}(tx)
	}

	fmt.Println("done")
	// Output: done
}
