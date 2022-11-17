package aptos_test

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-bcs/bcs"
)

const testUserMenmonic = "escape summer cupboard disagree coach mother permit sugar short excite road smoke"

func TestClient_EncodeSubmission(t *testing.T) {
	config, _ := aptos.GetAuxClientConfig(aptos.Testnet)
	moduleAddress := config.Address
	devnetClient := aptos.MustNewClient(aptos.Testnet, "")
	account, _ := aptos.NewLocalAccountFromMnemonic(testUserMenmonic, "")
	t.Logf("sender is %s", account.Address.String())
	expirationSecs := aptos.JsonUint64(time.Date(3000, 12, 31, 0, 0, 0, 0, time.UTC).Unix())
	t.Logf("expiration is %s", hex.EncodeToString(expirationSecs.ToBCS()))
	t.Logf("gas price is %s", hex.EncodeToString(aptos.JsonUint64(100).ToBCS()))
	t.Logf("max gas is %s", hex.EncodeToString(aptos.JsonUint64(20000).ToBCS()))
	t.Logf("sequence is %s", hex.EncodeToString(aptos.JsonUint64(0).ToBCS()))
	t.Logf("argument 1 is %s", hex.EncodeToString(aptos.JsonUint64(10000000000).ToBCS()))
	t.Logf("module name fake_coin is %s", hex.EncodeToString(bcs.MustMarshal("fake_coin")))
	t.Logf("function name mint is %s", hex.EncodeToString(bcs.MustMarshal("mint")))
	t.Logf("type name USDC is %s", hex.EncodeToString(bcs.MustMarshal("USDC")))

	tx := &aptos.Transaction{
		Sender:                  account.Address,
		ExpirationTimestampSecs: aptos.JsonUint64(expirationSecs),
		GasUnitPrice:            aptos.JsonUint64(100),
		MaxGasAmount:            aptos.JsonUint64(20000),
		SequenceNumber:          aptos.JsonUint64(0),
		Payload: aptos.NewEntryFunctionPayload(
			aptos.MustNewMoveFunctionTag(config.Address, "fake_coin", "mint"),
			[]*aptos.MoveStructTag{aptos.MustNewMoveStructTag(moduleAddress, "fake_coin", "USDC", nil)},
			[]*aptos.EntryFunctionArg{aptos.EntryFunctionArg_Uint64(10000000000)},
		),
		ChainId: aptos.GetChainIdForNetwork(aptos.Testnet),
	}

	request := &aptos.EncodeSubmissionRequest{
		Transaction: tx,
	}
	bodyStr, _ := request.Body()
	t.Log(string(bodyStr))
	r, err := devnetClient.EncodeSubmission(context.Background(), request)
	if err != nil {
		e, ok := err.(*aptos.AptosRestError)
		if ok {
			t.Fatalf("failed: %s", string(e.Message))
		} else {
			t.Fatalf("failed to encode: %#v", err)
		}
	}

	expected := "0xb5e97db07fa0bd0e5598aa3643a9bc6f6693bddc1a9fec9e674a461eaa00b1937d928500a7c0176468d16bd391c5b551bcea5c08394b19690ec8233fd464bcef0000000000000000028b7311d78d47e37d09435b8dc37c14afd977c5cfa74f974d45f0258d986eef530966616b655f636f696e046d696e7401078b7311d78d47e37d09435b8dc37c14afd977c5cfa74f974d45f0258d986eef530966616b655f636f696e045553444300010800e40b5402000000204e000000000000640000000000000000ae3e930700000002"
	result := string(*r.Parsed)
	if expected != result {
		t.Fatalf("want: %s\ngot:%s\n", expected, result)
	}

	result = "0x" + hex.EncodeToString(aptos.EncodeTransaction(tx))
	if expected != result {
		t.Fatalf("want: %s\ngot:%s\n", expected, result)
	}
}
