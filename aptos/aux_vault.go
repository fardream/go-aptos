package aptos

import "context"

// GetAuxOnChainSignerAddress calculates the onchain account holding assets for a given address.
// Assets for an aux user is held in a separate resource account, which is derived from the aux module
// address and seed "aux-user".
func GetAuxOnChainSignerAddress(auxModuleAddress, userAddress Address) Address {
	seeds := []byte{}
	seeds = append(seeds, auxModuleAddress[:]...)
	seeds = append(seeds, []byte("aux-user")...)

	return CalculateResourceAddress(userAddress, seeds)
}

// AuxUserAccount
type AuxUserAccount struct {
	AuthorizedTraders Table `json:"authorized_traders"`
}

// AuxCoinBalance
type AuxCoinBalance struct {
	Balance          JsonUint64 `json:"balance"`           // note on aux this is uint128.
	AvailableBalance JsonUint64 `json:"available_balance"` // note on aux this is uint128.
}

// GetCoinBalanceType get the coin balance and available balance in vault for a user.
func (info *AuxClientConfig) GetCoinBalanceType(coinType *MoveTypeTag) *MoveTypeTag {
	return &MoveTypeTag{
		MoveModuleTag: MoveModuleTag{
			Address: info.Address,
			Module:  "vault",
		},
		Name:                  "CoinBalance",
		GenericTypeParameters: []*MoveTypeTag{coinType},
	}
}

// GetAuxCoinBalance retrieves the balance for the user.
func (client *Client) GetAuxCoinBalance(ctx context.Context, auxInfo *AuxClientConfig, user Address, coinType *MoveTypeTag) (*AuxCoinBalance, error) {
	onchainAddress := GetAuxOnChainSignerAddress(auxInfo.Address, user)
	return GetAccountResourceWithType[AuxCoinBalance](ctx, client, onchainAddress, auxInfo.GetCoinBalanceType(coinType), 0)
}

// AuxVaultModuleName is the module name for vault.
const AuxVaultModuleName = "vault"

func (info *AuxClientConfig) Vault_CreateAuxAccount(sender Address, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxVaultModuleName, "create_aux_account")

	payload := NewEntryFunctionPayload(function, nil, nil)

	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Vault_Deposit deposits into vault.
func (info *AuxClientConfig) Vault_Deposit(sender Address, to Address, coinType *MoveTypeTag, amount uint64, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxVaultModuleName, "deposit")

	if to.IsZero() {
		to = sender
	}

	tx := &Transaction{
		Payload: NewEntryFunctionPayload(function, []*MoveTypeTag{coinType}, []EntryFunctionArg{
			to,
			JsonUint64(amount),
		}),
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// Vault_Withdraw withdraw from the vault.
func (info *AuxClientConfig) Vault_Withdraw(sender Address, coinType *MoveTypeTag, amount uint64, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxVaultModuleName, "withdraw")

	tx := &Transaction{
		Payload: NewEntryFunctionPayload(function, []*MoveTypeTag{coinType}, []EntryFunctionArg{JsonUint64(amount)}),
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}
