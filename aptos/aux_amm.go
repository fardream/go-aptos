package aptos

// AuxAmmPool is a constant product amm
type AuxAmmPool struct {
	FeeBps   JsonUint64 `json:"feebps"`
	Frozen   bool       `json:"frozen"`
	XReserve Coin       `json:"x_reserve"`
	YReserve Coin       `json:"y_reserve"`

	AddLiquidityEvents    *EventHandler `json:"add_liquidity_events"`
	RemoveLiquidityEvents *EventHandler `json:"remove_liquidity_events"`
	SwapEvents            *EventHandler `json:"swap_events"`
}

const AuxAmmModuleName = "amm"

func (info *AuxClientConfig) Amm_CreatePool(sender Address, coinX, coinY *MoveTypeTag, feeBps uint64, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxAmmModuleName, "create_pool")
	playload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{coinX, coinY},
		[]EntryFunctionArg{
			JsonUint64(feeBps),
		},
	)

	tx := &Transaction{Payload: playload}

	ApplyTransactionOptions(tx, options...)

	if tx.Sender.IsZero() {
		tx.Sender = sender
	}

	return tx
}

func (info *AuxClientConfig) Amm_UpdateFee(sender Address, coinX, coinY *MoveTypeTag, feeBps uint64, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxAmmModuleName, "update_fee")
	playload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{coinX, coinY},
		[]EntryFunctionArg{
			JsonUint64(feeBps),
		},
	)

	tx := &Transaction{Payload: playload}

	ApplyTransactionOptions(tx, options...)

	if tx.Sender.IsZero() {
		tx.Sender = sender
	}
	return tx
}
