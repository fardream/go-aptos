package aptos

import "github.com/fardream/go-bcs/bcs"

type TransactionPayload struct {
	ScriptPayload       uint8 `bcs:"-"`
	ModuleBundlePayload uint8 `bcs:"-"`
	*EntryFunctionPayload
}

var _ bcs.Enum = (*TransactionPayload)(nil)

func (p TransactionPayload) IsBcsEnum() {
}
