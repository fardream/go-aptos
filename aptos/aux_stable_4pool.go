// Code generated, DO NOT EDIT by hand
// from github.com/fardream/go-aptos/aptos/gen-aux-stable-pool

package aptos

import "github.com/fardream/go-bcs/bcs"

// AuxStable4Pool is a pool with 4 coins that priced at parity.
type AuxStable4Pool struct {
	/// FeeNumerator, denominator is 10^10
	FeeNumerator *bcs.Uint128 `json:"fee_numerator"`
	/// BalancedReserve
	BalancedReserve *bcs.Uint128 `json:"balanced_reserve"`
	/// Amp
	Amp *bcs.Uint128 `json:"amp"`

	Reserve0 *Coin        `json:"reserve_0"`
	Fee0     *Coin        `json:"fee_0"`
	Scaler0  *bcs.Uint128 `json:"scaler_0"`

	Reserve1 *Coin        `json:"reserve_1"`
	Fee1     *Coin        `json:"fee_1"`
	Scaler1  *bcs.Uint128 `json:"scaler_1"`

	Reserve2 *Coin        `json:"reserve_2"`
	Fee2     *Coin        `json:"fee_2"`
	Scaler2  *bcs.Uint128 `json:"scaler_2"`

	Reserve3 *Coin        `json:"reserve_3"`
	Fee3     *Coin        `json:"fee_3"`
	Scaler3  *bcs.Uint128 `json:"scaler_3"`
}

// AuxStable4PoolModuleName aux::stable_4pool
const AuxStable4PoolModuleName = "stable_4pool"

// AuxRouter4PoolModuleName aux::router_4pool
const AuxRouter4PoolModuleName = "router_4pool"

// Stable4PoolType returns the move type tag ([MoveTypeTag]) for a stable 4 pool
func (info *AuxClientConfig) Stable4PoolType(coin0, coin1, coin2, coin3 *MoveTypeTag) (*MoveTypeTag, error) {
	return NewMoveTypeTag(
		info.Address,
		AuxStable4PoolModuleName,
		"Pool",
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
			coin3,
		})
}

// Stable4Pool_CreatePool construct the transaction to create a new 4pool
func (info *AuxClientConfig) Stable4Pool_CreatePool(
	sender Address,
	coin0 *MoveTypeTag,
	coin1 *MoveTypeTag,
	coin2 *MoveTypeTag,
	coin3 *MoveTypeTag,
	feeNumerator uint64,
	amp uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxStable4PoolModuleName, "create_pool")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
			coin3,
		},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint128(feeNumerator, 0),
			EntryFunctionArg_Uint128(amp, 0),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Router4Pool_AddLiquidity constructs the transaction to add liquidity to a pool
func (info *AuxClientConfig) Router4Pool_AddLiquidity(
	sender Address,
	coin0 *MoveTypeTag,
	amount0 uint64,
	coin1 *MoveTypeTag,
	amount1 uint64,
	coin2 *MoveTypeTag,
	amount2 uint64,
	coin3 *MoveTypeTag,
	amount3 uint64,
	minLpAmount uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter4PoolModuleName, "add_liquidity")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
			coin3,
		},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint64(amount0),
			EntryFunctionArg_Uint64(amount1),
			EntryFunctionArg_Uint64(amount2),
			EntryFunctionArg_Uint64(amount3),
			EntryFunctionArg_Uint64(minLpAmount),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Router4Pool_RemoveLiquidityForCoin constructs the transaction to remove liquidity from the pool
// by specifying the coin amount to withdraw.
func (info *AuxClientConfig) Router4Pool_RemoveLiquidityForCoin(
	sender Address,
	coin0 *MoveTypeTag,
	amount0ToWithdraw uint64,
	coin1 *MoveTypeTag,
	amount1ToWithdraw uint64,
	coin2 *MoveTypeTag,
	amount2ToWithdraw uint64,
	coin3 *MoveTypeTag,
	amount3ToWithdraw uint64,
	maxLpAmount uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter4PoolModuleName, "remove_liquidity_for_coin")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
			coin3,
		},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint64(amount0ToWithdraw),
			EntryFunctionArg_Uint64(amount1ToWithdraw),
			EntryFunctionArg_Uint64(amount2ToWithdraw),
			EntryFunctionArg_Uint64(amount3ToWithdraw),
			EntryFunctionArg_Uint64(maxLpAmount),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Router4Pool_RemoveLiquidity constructs the transaction to remove liquidity from the pool
// by specifying the amount of lp coins to burn.
func (info *AuxClientConfig) Router4Pool_RemoveLiquidity(
	sender Address,
	coin0 *MoveTypeTag,
	coin1 *MoveTypeTag,
	coin2 *MoveTypeTag,
	coin3 *MoveTypeTag,
	lpAmount uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter4PoolModuleName, "remove_liquidity")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
			coin3,
		},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint64(lpAmount),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Router4Pool_SwapExactCoinForCoin constructs the transaction to swap coins by
// by specifying the input amount of coins to swap.
func (info *AuxClientConfig) Router4Pool_SwapExactCoinForCoin(
	sender Address,
	coin0 *MoveTypeTag,
	amount0 uint64,
	coin1 *MoveTypeTag,
	amount1 uint64,
	coin2 *MoveTypeTag,
	amount2 uint64,
	coin3 *MoveTypeTag,
	amount3 uint64,
	outCoinIndex int,
	minQuantityOut uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter4PoolModuleName, "swap_exact_coin_for_coin")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
			coin3,
		},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint64(amount0),
			EntryFunctionArg_Uint64(amount1),
			EntryFunctionArg_Uint64(amount2),
			EntryFunctionArg_Uint64(amount3),
			EntryFunctionArg_Uint8(uint8(outCoinIndex)),
			EntryFunctionArg_Uint64(minQuantityOut),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Router4Pool_SwapCoinForExactCoin constructs the transaction to swap coins by
// by specifying the input amount of coins to swap.
func (info *AuxClientConfig) Router4Pool_SwapCoinForExactCoin(
	sender Address,
	coin0 *MoveTypeTag,
	requestAmount0 uint64,
	coin1 *MoveTypeTag,
	requestAmount1 uint64,
	coin2 *MoveTypeTag,
	requestAmount2 uint64,
	coin3 *MoveTypeTag,
	requestAmount3 uint64,
	inCoinIndex int,
	maxQuantityIn uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter4PoolModuleName, "swap_coin_for_exact_coin")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{
			coin0,
			coin1,
			coin2,
			coin3,
		},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint64(requestAmount0),
			EntryFunctionArg_Uint64(requestAmount1),
			EntryFunctionArg_Uint64(requestAmount2),
			EntryFunctionArg_Uint64(requestAmount3),
			EntryFunctionArg_Uint8(uint8(inCoinIndex)),
			EntryFunctionArg_Uint64(maxQuantityIn),
		},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}
