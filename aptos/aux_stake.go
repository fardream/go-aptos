package aptos

import "github.com/fardream/go-bcs/bcs"

// AuxStakePool provide staking reward for users locking up their coins.
// See code [here].
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/2022-12-13/aptos/contract/aux/sources/stake.move#L158-L178
type AuxStakePool struct {
	Authority         Address      `json:"address"`
	StartTime         JsonUint64   `json:"start_time"`
	EndTime           JsonUint64   `json:"end_time"`
	RewardRemaining   JsonUint64   `json:"reward_remaining"`
	Stake             Coin         `json:"stake"`
	Reward            Coin         `json:"reward"`
	LastUpdateTime    JsonUint64   `json:"last_update_time"`
	AccRewardPerShare *bcs.Uint128 `json:"acc_reward_per_share"`

	CreatePoolEvents EventHandler `json:"create_pool_events"`
	DepositEvents    EventHandler `json:"deposit_events"`
	WithdrawEvents   EventHandler `json:"widthraw_events"`
	ModifyPoolEvents EventHandler `json:"modify_pool_events"`
	ClaimEvents      EventHandler `json:"claim_events"`
}

const AuxStakeModuleName = "stake"

// StakePoolType returns the aux::stake::Pool<Stake, Reward>
func (info *AuxClientConfig) StakePoolType(stake, reward *MoveStructTag) (*MoveStructTag, error) {
	return NewMoveStructTag(
		info.Address,
		AuxStakeModuleName,
		"Pool",
		[]*MoveStructTag{
			stake,
			reward,
		})
}

// Stake_CreatePool creates a new staking pool.
// see [here]
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/2022-12-13/aptos/contract/aux/sources/stake.move#L233-L326
func (info *AuxClientConfig) Stake_CreatePool(
	sender Address,
	stake *MoveStructTag,
	reward *MoveStructTag,
	rewardAmount uint64,
	durationUs uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxStakeModuleName, "create")

	payload := NewEntryFunctionPayload(
		function,
		[]*MoveStructTag{stake, reward},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint64(rewardAmount),
			EntryFunctionArg_Uint64(durationUs),
		})

	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Stake_DeleteEmptyPool deletes a pool when it's empty.
// see [here]
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/2022-12-13/aptos/contract/aux/sources/stake.move#L328-L384
func (info *AuxClientConfig) Stake_DeleteEmptyPool(sender Address, stake *MoveStructTag, reward *MoveStructTag, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxStakeModuleName, "delete_empty_pool")
	payload := NewEntryFunctionPayload(function, []*MoveStructTag{stake, reward}, []*EntryFunctionArg{})

	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Stake_EndRewardEarly ends the reward early for a pool.
// see [here]
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/2022-12-13/aptos/contract/aux/sources/stake.move#L416-L451
func (info *AuxClientConfig) Stake_EndRewardEarly(sender Address, stake *MoveStructTag, reward *MoveStructTag, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxStakeModuleName, "end_reward_early")
	payload := NewEntryFunctionPayload(function, []*MoveStructTag{stake, reward}, []*EntryFunctionArg{})

	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Stake_ModifyPool modifies the time and reward of the pool.
// see [here].
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/2022-12-13/aptos/contract/aux/sources/stake.move#L453-L525
func (info *AuxClientConfig) Stake_ModifyPool(
	sender Address,
	stake *MoveStructTag,
	reward *MoveStructTag,
	rewardAmount uint64,
	rewardIncrease bool,
	timeAmountUs uint64,
	timeIncrease bool,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxStakeModuleName, "modify_pool")
	payload := NewEntryFunctionPayload(
		function,
		[]*MoveStructTag{stake, reward},
		[]*EntryFunctionArg{
			EntryFunctionArg_Uint64(rewardAmount),
			EntryFunctionArg_Bool(rewardIncrease),
			EntryFunctionArg_Uint64(timeAmountUs),
			EntryFunctionArg_Bool(timeIncrease),
		})

	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Stake_Deposit deposit stake coins into the reward pool to earn reward.
// see [here].
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/2022-12-13/aptos/contract/aux/sources/stake.move#L527-L591
func (info *AuxClientConfig) Stake_Deposit(
	sender Address,
	stake *MoveStructTag,
	reward *MoveStructTag,
	amount uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxStakeModuleName, "deposit")
	payload := NewEntryFunctionPayload(function, []*MoveStructTag{stake, reward}, []*EntryFunctionArg{EntryFunctionArg_Uint64(amount)})
	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Stake_Withdraw withdraw staked coins from the pool. Reward will be claimed.
// see [here].
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/2022-12-13/aptos/contract/aux/sources/stake.move#L593-L640
func (info *AuxClientConfig) Stake_Withdraw(
	sender Address,
	stake *MoveStructTag,
	reward *MoveStructTag,
	amount uint64,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxStakeModuleName, "withdraw")
	payload := NewEntryFunctionPayload(function, []*MoveStructTag{stake, reward}, []*EntryFunctionArg{EntryFunctionArg_Uint64(amount)})
	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Stake_Claim claims the reward but doesn't unstake the coins.
// see [here].
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/2022-12-13/aptos/contract/aux/sources/stake.move#L642-L671
func (info *AuxClientConfig) Stake_Claim(
	sender Address,
	stake *MoveStructTag,
	reward *MoveStructTag,
	options ...TransactionOption,
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxStakeModuleName, "claim")
	payload := NewEntryFunctionPayload(function, []*MoveStructTag{stake, reward}, []*EntryFunctionArg{})
	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}
