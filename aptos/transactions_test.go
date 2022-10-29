package aptos_test

import (
	_ "embed"
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
