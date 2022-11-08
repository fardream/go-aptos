// Code generated, DO NOT EDIT by hand
// from github.com/fardream/go-aptos/aptos/gen-aux-stable-pool

package aptos

// AuxStable3Pool is a pool with 3 coins that priced at parity.
type AuxStable3Pool struct {
	/// FeeNumerator, denominator is 10^10
	FeeNumerator *Uint128 `json:"fee_numerator"`
	/// BalancedReserve
	BalancedReserve *Uint128 `json:"balanced_reserve"`
	/// Amp
	Amp *Uint128 `json:"amp"`

	Reserve0 *Coin    `json:"reserve_0"`
	Fee0     *Coin    `json:"fee_0"`
	Scaler0  *Uint128 `json:"scaler_0"`

	Reserve1 *Coin    `json:"reserve_1"`
	Fee1     *Coin    `json:"fee_1"`
	Scaler1  *Uint128 `json:"scaler_1"`

	Reserve2 *Coin    `json:"reserve_2"`
	Fee2     *Coin    `json:"fee_2"`
	Scaler2  *Uint128 `json:"scaler_2"`
}

// AuxStable3PoolModuleName aux::stable_3pool
const AuxStable3PoolModuleName = "stable_3pool"

// AuxRouter3PoolModuleName aux::router_3pool
const AuxRouter3PoolModuleName = "router_3pool"

// Stable3PoolType returns the move type tag ([MoveTypeTag]) for a stable 3 pool
func (info *AuxClientConfig) Stable3PoolType(coin0, coin1, coin2 *MoveTypeTag) (*MoveTypeTag, error) {
	return NewMoveTypeTag(
		info.Address,
		AuxStable3PoolModuleName,
		"Pool",
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
		})
}

// Stable3Pool_CreatePool construct the transaction to create a new 3pool
func (info *AuxClientConfig) Stable3Pool_CreatePool(
	sender Address,
	coin0 *MoveTypeTag,
	coin1 *MoveTypeTag,
	coin2 *MoveTypeTag,
	feeNumerator uint64,
	amp uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxStable3PoolModuleName, "create_pool")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
		},
		[]EntryFunctionArg{
			NewUint128FromUint64(feeNumerator, 0),
			NewUint128FromUint64(amp, 0),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Router3Pool_AddLiquidity constructs the transaction to add liquidity to a pool
func (info *AuxClientConfig) Router3Pool_AddLiquidity(
	sender Address,
	coin0 *MoveTypeTag,
	amount0 uint64,
	coin1 *MoveTypeTag,
	amount1 uint64,
	coin2 *MoveTypeTag,
	amount2 uint64,
	minLpAmount uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter3PoolModuleName, "add_liquidity")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
		},
		[]EntryFunctionArg{
			JsonUint64(amount0),
			JsonUint64(amount1),
			JsonUint64(amount2),
			JsonUint64(minLpAmount),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Router3Pool_RemoveLiquidityForCoin constructs the transaction to remove liquidity from the pool
// by specifying the coin amount to withdraw.
func (info *AuxClientConfig) Router3Pool_RemoveLiquidityForCoin(
	sender Address,
	coin0 *MoveTypeTag,
	amount0ToWithdraw uint64,
	coin1 *MoveTypeTag,
	amount1ToWithdraw uint64,
	coin2 *MoveTypeTag,
	amount2ToWithdraw uint64,
	maxLpAmount uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter3PoolModuleName, "remove_liquidity_for_coin")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
		},
		[]EntryFunctionArg{
			JsonUint64(amount0ToWithdraw),
			JsonUint64(amount1ToWithdraw),
			JsonUint64(amount2ToWithdraw),
			JsonUint64(maxLpAmount),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Router3Pool_RemoveLiquidity constructs the transaction to remove liquidity from the pool
// by specifying the amount of lp coins to burn.
func (info *AuxClientConfig) Router3Pool_RemoveLiquidity(
	sender Address,
	coin0 *MoveTypeTag,
	coin1 *MoveTypeTag,
	coin2 *MoveTypeTag,
	lpAmount uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter3PoolModuleName, "remove_liquidity")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
		},
		[]EntryFunctionArg{
			JsonUint64(lpAmount),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Router3Pool_SwapExactCoinForCoin constructs the transaction to swap coins by
// by specifying the input amount of coins to swap.
func (info *AuxClientConfig) Router3Pool_SwapExactCoinForCoin(
	sender Address,
	coin0 *MoveTypeTag,
	amount0 uint64,
	coin1 *MoveTypeTag,
	amount1 uint64,
	coin2 *MoveTypeTag,
	amount2 uint64,
	outCoinIndex int,
	minQuantityOut uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter3PoolModuleName, "swap_exact_coin_for_coin")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
		},
		[]EntryFunctionArg{
			JsonUint64(amount0),
			JsonUint64(amount1),
			JsonUint64(amount2),
			EntryFunctionArg_Uint8(outCoinIndex),
			JsonUint64(minQuantityOut),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Router3Pool_SwapCoinForExactCoin constructs the transaction to swap coins by
// by specifying the input amount of coins to swap.
func (info *AuxClientConfig) Router3Pool_SwapCoinForExactCoin(
	sender Address,
	coin0 *MoveTypeTag,
	requestAmount0 uint64,
	coin1 *MoveTypeTag,
	requestAmount1 uint64,
	coin2 *MoveTypeTag,
	requestAmount2 uint64,
	inCoinIndex int,
	maxQuantityIn uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter3PoolModuleName, "swap_coin_for_exact_coin")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
		},
		[]EntryFunctionArg{
			JsonUint64(requestAmount0),
			JsonUint64(requestAmount1),
			JsonUint64(requestAmount2),
			EntryFunctionArg_Uint8(inCoinIndex),
			JsonUint64(maxQuantityIn),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}
