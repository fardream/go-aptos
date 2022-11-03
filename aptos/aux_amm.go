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

// AmmPoolType returns the move type ([MoveTypeTag]) for a pool
func (info *AuxClientConfig) AmmPoolType(coinX *MoveTypeTag, coinY *MoveTypeTag) (*MoveTypeTag, error) {
	return NewMoveTypeTag(info.Address, "amm", "Pool", []*MoveTypeTag{coinX, coinY})
}

// AuxAmmModuleName aux::amm
const AuxAmmModuleName = "amm"

// Amm_CreatePool creates a new pool with the give coin x and coin y.
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

	tx.Sender = sender

	return tx
}

// Amm_UpdateFee updates the fee of the amm pool.
// the pool is identified by the coin types.
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

	tx.Sender = sender

	return tx
}
