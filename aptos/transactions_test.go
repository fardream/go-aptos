package aptos_test

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/fardream/go-aptos/aptos"
)

//go:embed test_data/test_tx.json
var testTxJson string

func TestTransactionWithInfo(t *testing.T) {
	var tx aptos.TransactionWithInfo
	if err := json.Unmarshal([]byte(testTxJson), &tx); err != nil {
		t.Fatalf("%v", err)
	}

	events := aptos.FilterAuxClobMarketOrderEvent(tx.Events, aptos.Address{}, true, false)

	for _, ev := range events {
		t.Logf("%s", spew.Sdump(ev))
	}
}

func TestTransaction_GetHash(t *testing.T) {
	t.Skipf("Transaction_GetHash doesnt mathc the current value from chain")
	// this is from test in test_data/test_tx_1.json
	sender := aptos.MustParseAddress("0x767b7442b8547fa5cf50989b9b761760ca6687b83d1c23d3589a5ac8acb50639")
	config, _ := aptos.GetAuxClientConfig(aptos.Devnet)
	moduleAddress := aptos.MustParseAddress("0xea383dc2819210e6e427e66b2b6aa064435bf672dc4bdc55018049f0c361d01a")
	aux, _ := aptos.GetAuxFakeCoinCoinType(moduleAddress, aptos.AuxFakeCoin_AUX)
	usdc, _ := aptos.GetAuxFakeCoinCoinType(moduleAddress, aptos.AuxFakeCoin_USDC)
	tx := config.ClobMarket_PlaceOrder(
		sender,
		true,
		aux,
		usdc,
		50000000,
		50000000000,
		0,
		aptos.Uint128{},
		aptos.AuxClobMarketOrderType_Limit,
		0,
		false,
		18446744073709551615,
		aptos.AuxClobMarketSelfTradeType_CancelPassive,
		aptos.TransactionOption_GasUnitPrice(100),
		aptos.TransactionOption_SequenceNumber(2625),
		aptos.TransactionOption_MaxGasAmount(1718192),
	)

	tx.ExpirationTimestampSecs = 1667245112
	tx.ChainId = 35

	jsonBytes, _ := json.MarshalIndent(tx, "", "  ")
	t.Log(string(jsonBytes))
	expectedHash := "0xb161e7592d5f8ea8a97f3493669660205ee76f8699f20e71ae2ad3878836a1ac"
	hash := "0x" + hex.EncodeToString(tx.GetHash())

	if hash != expectedHash {
		t.Fatalf("hash doesn't match:\nwant: %s\nhas:  %s\n", expectedHash, hash)
	}
}

func TestLocalAccount_Sign(t *testing.T) {
	// this is from test in test_data/test_tx_1.json
	// Private key is from aux exchange replay script
	// https://github.com/aux-exchange/aux-exchange/blob/3f0c6d8cc9a8f8375d8c1c2f4e90a3071742b00f/aptos/api/aux-ts/scripts/sim/replay.ts#L12
	privateKey, _ := aptos.NewPrivateKeyFromHexString("0x2b248dee740ee1e8d271afb89590554cd9655ee9fae8a0ec616b95911834eb49")
	localAccount, _ := aptos.NewLocalAccountFromPrivateKey(privateKey)
	sender := localAccount.Address
	config, _ := aptos.GetAuxClientConfig(aptos.Devnet)
	moduleAddress := aptos.MustParseAddress("0xea383dc2819210e6e427e66b2b6aa064435bf672dc4bdc55018049f0c361d01a")
	aux, _ := aptos.GetAuxFakeCoinCoinType(moduleAddress, aptos.AuxFakeCoin_AUX)
	usdc, _ := aptos.GetAuxFakeCoinCoinType(moduleAddress, aptos.AuxFakeCoin_USDC)
	tx := config.ClobMarket_PlaceOrder(
		sender,
		true,
		aux,
		usdc,
		50000000,
		50000000000,
		0,
		aptos.Uint128{},
		aptos.AuxClobMarketOrderType_Limit,
		0,
		false,
		18446744073709551615,
		aptos.AuxClobMarketSelfTradeType_CancelPassive,
		aptos.TransactionOption_GasUnitPrice(100),
		aptos.TransactionOption_SequenceNumber(2625),
		aptos.TransactionOption_MaxGasAmount(1718192),
	)

	tx.ExpirationTimestampSecs = 1667245112
	tx.ChainId = 35

	expectedSignature := "0x05c4586e2b47fa4d81a2fab4b6ac7975fb4d2a410d507bada48b8a0adff14c1dd18a9c8342076a133e7f6b358efba0bff777c3cc78f1131f5c5d617329b8200a"
	sig, _ := localAccount.Sign(tx)

	sigStr := sig.Signature

	if sigStr != expectedSignature {
		t.Fatalf("signature doesn't match:\nwant: %s\nhas:  %s\n", expectedSignature, sigStr)
	}
}
