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

func (auxInfo *AuxClientConfig) GetCoinBalanceType(coinType *MoveTypeTag) *MoveTypeTag {
	return &MoveTypeTag{
		MoveModuleTag: MoveModuleTag{
			Address: auxInfo.Address,
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
