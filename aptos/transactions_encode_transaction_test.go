package aptos_test

import (
	"encoding/hex"
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/fardream/go-aptos/aptos"
)

const networkForEncodeTransactionTest = aptos.Devnet

func TestEncodeTransaction(t *testing.T) {
	auxConfig, _ := aptos.GetAuxClientConfig(networkForEncodeTransactionTest)

	eth, _ := aptos.GetAuxFakeCoinCoinType(auxConfig.Address, aptos.AuxFakeCoin_ETH)
	btc, _ := aptos.GetAuxFakeCoinCoinType(auxConfig.Address, aptos.AuxFakeCoin_BTC)

	seqnum := aptos.TransactionOption_SequenceNumber(3)
	expiry := aptos.TransactionOption_ExpireAt(time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC))
	maxGas := aptos.TransactionOption_MaxGasAmount(345678)
	gasUnitPrice := aptos.TransactionOption_GasUnitPrice(111)
	tx := auxConfig.ClobMarket_PlaceOrder(auxConfig.DataFeedAddress, true, eth, btc, 30, 90, 1, aptos.AuxClobMarketOrderType_FOK, 1, false,
		math.MaxUint64-5000000, aptos.AuxClobMarketSelfTradeType_CancelBoth, seqnum, expiry, maxGas, gasUnitPrice)
	jsonPayload, _ := json.MarshalIndent(tx, "", "  ")
	t.Logf("tx:\n%s\n", string(jsonPayload))

	encoded := aptos.EncodeTransaction(tx, 34)
	obtained := "0x" + hex.EncodeToString(encoded)
	expected := "0xb5e97db07fa0bd0e5598aa3643a9bc6f6693bddc1a9fec9e674a461eaa00b19384f372536c73df84327d2af63992f4443e2bd1aec8695fa85693e256fc1f904f030000000000000002ea383dc2819210e6e427e66b2b6aa064435bf672dc4bdc55018049f0c361d01a0b636c6f625f6d61726b65740b706c6163655f6f726465720207ea383dc2819210e6e427e66b2b6aa064435bf672dc4bdc55018049f0c361d01a0966616b655f636f696e0846616b65436f696e0107ea383dc2819210e6e427e66b2b6aa064435bf672dc4bdc55018049f0c361d01a0966616b655f636f696e034554480007ea383dc2819210e6e427e66b2b6aa064435bf672dc4bdc55018049f0c361d01a0966616b655f636f696e0846616b65436f696e0107ea383dc2819210e6e427e66b2b6aa064435bf672dc4bdc55018049f0c361d01a0966616b655f636f696e03425443000b2084f372536c73df84327d2af63992f4443e2bd1aec8695fa85693e256fc1f904f0101081e00000000000000085a000000000000000801000000000000001000000000000000000000000000000000086500000000000000080100000000000000010008bfb4b3ffffffffff08ca000000000000004e460500000000006f0000000000000000cc5e910700000022"

	if obtained != expected {
		t.Fatalf("\nwant: %s\nhas:  %s\n", expected, obtained)
	}
}
