package aptos

import (
	"encoding/json"
	"fmt"

	"github.com/fardream/go-bcs/bcs"
)

// TransactionPayload is payload field of a [transaction].
// Although there are three types, script, module bundle, and entry function, only
// entry function is supported here.
//
// [transaction]: https://fullnode.mainnet.aptoslabs.com/v1/spec#/schemas/Transaction
type TransactionPayload struct {
	ScriptPayload       *json.RawMessage `bcs:"-"`
	ModuleBundlePayload *json.RawMessage `bcs:"-"`
	*EntryFunctionPayload
}

// transactionPayload_EntryFunction is a helper class for json de/serialization.
type transactionPayload_EntryFunction struct {
	*EntryFunctionPayload `json:",inline"`
	Type                  string `json:"type"`
}

var (
	_ bcs.Enum         = (*TransactionPayload)(nil)
	_ json.Marshaler   = (*TransactionPayload)(nil)
	_ json.Unmarshaler = (*TransactionPayload)(nil)
)

func (p TransactionPayload) IsBcsEnum() {
}

const entryFunctionPayloadTypeStr = "entry_function_payload"

func (p TransactionPayload) MarshalJSON() ([]byte, error) {
	if p.EntryFunctionPayload == nil {
		return nil, fmt.Errorf("only entry function payload is supported")
	}

	return json.Marshal(transactionPayload_EntryFunction{
		EntryFunctionPayload: p.EntryFunctionPayload,
		Type:                 entryFunctionPayloadTypeStr,
	})
}

func (p *TransactionPayload) UnmarshalJSON(data []byte) error {
	var tmp transactionPayload_EntryFunction
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	switch tmp.Type {
	case entryFunctionPayloadTypeStr:
		p.EntryFunctionPayload = tmp.EntryFunctionPayload
	case "script_payload":
		p.ScriptPayload = (*json.RawMessage)(&data)
	case "module_bundel_payloaad":
		p.ModuleBundlePayload = (*json.RawMessage)(&data)
	}

	return nil
}
