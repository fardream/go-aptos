package aptos

import (
	"context"
)

// AptosCoin is the type for aptos coin
var AptosCoin = MoveTypeTag{
	MoveModuleTag: MoveModuleTag{
		Address: AptosStdAddress,
		Module:  "aptos_coin",
	},
	Name: "AptosCoin",
}

// GetCoinStoreType returns the 0x1::coin::CoinStore<T>
func GetCoinStoreType(coin *MoveTypeTag) *MoveTypeTag {
	return &MoveTypeTag{
		MoveModuleTag: MoveModuleTag{
			Address: AptosStdAddress,
			Module:  "coin",
		},
		Name: "CoinStore",

		GenericTypeParameters: []*MoveTypeTag{coin},
	}
}

// Coin, this is golang equivalent of 0x1::coin::Coin
type Coin struct {
	Value JsonUint64 `json:"value"`
}

// CoinStore, this is golang equivalent of 0x1::coin::CoinStore
type CoinStore struct {
	Coin           Coin         `json:"coin"`
	Frozen         bool         `json:"frozen"`
	DepositEvents  EventHandler `json:"deposit_events"`
	WithdrawEvents EventHandler `json:"withdraw_events"`
}

// GetCoinBalance
func (client *Client) GetCoinBalance(ctx context.Context, address Address, coinType *MoveTypeTag) (uint64, error) {
	coinStore, err := GetAccountResourceWithType[CoinStore](ctx, client, address, GetCoinStoreType(coinType), 0)
	if err != nil {
		return 0, err
	}

	return uint64(coinStore.Coin.Value), nil
}
