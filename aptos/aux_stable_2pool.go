// Code generated, DO NOT EDIT by hand
// from github.com/fardream/go-aptos/aptos/gen-aux-stable-pool

package aptos

import "github.com/fardream/go-bcs/bcs"

// AuxStable2Pool is a pool with 2 coins that priced at parity.
type AuxStable2Pool struct {
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
}

// AuxStable2PoolModuleName aux::stable_2pool
const AuxStable2PoolModuleName = "stable_2pool"

// AuxRouter2PoolModuleName aux::router_2pool
const AuxRouter2PoolModuleName = "router_2pool"

// Stable2PoolType returns the move struct tag ([MoveStructTag]) for a stable 2 pool
func (info *AuxClientConfig) Stable2PoolType(coin0, coin1 *MoveStructTag) (*MoveStructTag, error) {
	return NewMoveStructTag(
		info.Address,
		AuxStable2PoolModuleName,
		"Pool",
		[]*MoveStructTag{
			coin0,
			coin1,
		})
}

// Stable2Pool_CreatePool construct the transaction to create a new 2pool
func (info *AuxClientConfig) Stable2Pool_CreatePool(
	sender Address,
	coin0 *MoveStructTag,
	coin1 *MoveStructTag,
	feeNumerator uint64,
	amp uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxStable2PoolModuleName, "create_pool")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveStructTag{
			coin0,
			coin1,
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

// Router2Pool_AddLiquidity constructs the transaction to add liquidity to a pool
func (info *AuxClientConfig) Router2Pool_AddLiquidity(
	sender Address,
	coin0 *MoveStructTag,
	amount0 uint64,
	coin1 *MoveStructTag,
	amount1 uint64,
	minLpAmount uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter2PoolModuleName, "add_liquidity")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveStructTag{
			coin0,
			coin1,
		},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint64(amount0),
			EntryFunctionArg_Uint64(amount1),
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

// Router2Pool_RemoveLiquidityForCoin constructs the transaction to remove liquidity from the pool
// by specifying the coin amount to withdraw.
func (info *AuxClientConfig) Router2Pool_RemoveLiquidityForCoin(
	sender Address,
	coin0 *MoveStructTag,
	amount0ToWithdraw uint64,
	coin1 *MoveStructTag,
	amount1ToWithdraw uint64,
	maxLpAmount uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter2PoolModuleName, "remove_liquidity_for_coin")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveStructTag{
			coin0,
			coin1,
		},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint64(amount0ToWithdraw),
			EntryFunctionArg_Uint64(amount1ToWithdraw),
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

// Router2Pool_RemoveLiquidity constructs the transaction to remove liquidity from the pool
// by specifying the amount of lp coins to burn.
func (info *AuxClientConfig) Router2Pool_RemoveLiquidity(
	sender Address,
	coin0 *MoveStructTag,
	coin1 *MoveStructTag,
	lpAmount uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter2PoolModuleName, "remove_liquidity")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveStructTag{
			coin0,
			coin1,
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

// Router2Pool_SwapExactCoinForCoin constructs the transaction to swap coins by
// by specifying the input amount of coins to swap.
func (info *AuxClientConfig) Router2Pool_SwapExactCoinForCoin(
	sender Address,
	coin0 *MoveStructTag,
	amount0 uint64,
	coin1 *MoveStructTag,
	amount1 uint64,
	outCoinIndex int,
	minQuantityOut uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter2PoolModuleName, "swap_exact_coin_for_coin")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveStructTag{
			coin0,
			coin1,
		},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint64(amount0),
			EntryFunctionArg_Uint64(amount1),
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

// Router2Pool_SwapCoinForExactCoin constructs the transaction to swap coins by
// by specifying the input amount of coins to swap.
func (info *AuxClientConfig) Router2Pool_SwapCoinForExactCoin(
	sender Address,
	coin0 *MoveStructTag,
	requestAmount0 uint64,
	coin1 *MoveStructTag,
	requestAmount1 uint64,
	inCoinIndex int,
	maxQuantityIn uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxRouter2PoolModuleName, "swap_coin_for_exact_coin")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveStructTag{
			coin0,
			coin1,
		},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint64(requestAmount0),
			EntryFunctionArg_Uint64(requestAmount1),
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
