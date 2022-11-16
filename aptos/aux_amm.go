package aptos

// AuxAmmPool is a constant product amm
type AuxAmmPool struct {
	FeeBps   JsonUint64 `json:"fee_bps"`
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

// AuxAmm_AddLiquidityEvent is emitted when liquidity is added. contract [here]
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/v1.0.4/aptos/contract/aux/sources/amm.move#L86-L93
type AuxAmm_AddLiquidityEvent struct {
	Timestamp  JsonUint64   `json:"timestamp"`
	XCoinType  *MoveTypeTag `json:"x_coin_type"`
	YCoinType  *MoveTypeTag `json:"y_coin_type"`
	XAddedAu   JsonUint64   `json:"x_added_au"`
	YAddedAu   JsonUint64   `json:"y_added_au"`
	LpMintedAu JsonUint64   `json:"lp_minted_au"`
}

// AuxAmm_RemoveLiquidityEvent is emitted when liquidity is removed. contract [here]
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/v1.0.4/aptos/contract/aux/sources/amm.move#L95-L102
type AuxAmm_RemoveLiquidityEvent struct {
	Timestamp  JsonUint64   `json:"timestamp"`
	XCoinType  *MoveTypeTag `json:"x_coin_type"`
	YCoinType  *MoveTypeTag `json:"y_coin_type"`
	XRemovedAu JsonUint64   `json:"x_removed_au"`
	YRemovedAu JsonUint64   `json:"y_removed_au"`
	LpBurnedAu JsonUint64   `json:"lp_burned_au"`
}

// AuxAmm_SwapEvent is emitted when a swap happens on chain. contract [here]
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/v1.0.4/aptos/contract/aux/sources/amm.move#L74-L84
type AuxAmm_SwapEvent struct {
	SenderAddr  Address      `json:"sender_addr"`
	Timestamp   JsonUint64   `json:"timestamp"`
	InCoinType  *MoveTypeTag `json:"in_coin_type"`
	OutCoinType *MoveTypeTag `json:"out_coin_type"`
	InReserve   JsonUint64   `json:"in_reserve"`
	OutReserve  JsonUint64   `json:"out_reserve"`
	InAu        JsonUint64   `json:"in_au"`
	OutAu       JsonUint64   `json:"out_au"`
	FeeBps      JsonUint64   `json:"fee_bps"`
}

// Amm_CreatePool creates a new pool with the give coin x and coin y. contract [here]
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/v1.0.4/aptos/contract/aux/sources/amm.move#L117-L163
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
// contract [here].
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/v1.0.4/aptos/contract/aux/sources/amm.move#L165-L173
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

// Amm_SwapExactCoinForCoin swaps coins, with the output amount decided by the input amount.
// See contract [here]
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/v1.0.4/aptos/contract/aux/sources/amm.move#L176-L199
func (info *AuxClientConfig) Amm_SwapExactCoinForCoin(sender Address, coinX, coinY *MoveTypeTag, amountIn uint64, minAmountOut uint64, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxAmmModuleName, "swap_exact_coin_for_coin_with_signer")

	payload := NewEntryFunctionPayload(function, []*MoveTypeTag{coinX, coinY}, []EntryFunctionArg{
		JsonUint64(amountIn),
		JsonUint64(minAmountOut),
	})

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Amm_AddLiquidity adds liquidity to amm.
// See contract [here]
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/v1.0.4/aptos/contract/aux/sources/amm.move#L287-L306
func (info *AuxClientConfig) Amm_AddLiquidity(sender Address, coinX *MoveTypeTag, xAmount uint64, coinY *MoveTypeTag, yAmount uint64, maxSlippage uint64, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxAmmModuleName, "add_liquidity")

	payload := NewEntryFunctionPayload(function, []*MoveTypeTag{coinX, coinY}, []EntryFunctionArg{
		JsonUint64(xAmount),
		JsonUint64(yAmount),
		JsonUint64(maxSlippage),
	})

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}
