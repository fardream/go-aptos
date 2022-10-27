package aptos

import (
	"context"
)

var AptosCoin = MoveTypeTag{
	MoveModuleTag: MoveModuleTag{
		Address: AptosStdAddress,
		Module:  "aptos_coin",
	},
	Name: "AptosCoin",
}

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

type Coin struct {
	Value JsonUint64 `json:"value"`
}

type CoinStore struct {
	Coin           Coin         `json:"coin"`
	Frozen         bool         `json:"frozen"`
	DepositEvents  EventHandler `json:"deposit_events"`
	WithdrawEvents EventHandler `json:"withdraw_events"`
}

func (client *Client) GetCoinBalance(ctx context.Context, address Address, coinType *MoveTypeTag) (uint64, error) {
	coinStore, err := GetAccountResourceWithType[CoinStore](ctx, client, address, GetCoinStoreType(coinType), 0)
	if err != nil {
		return 0, err
	}

	return uint64(coinStore.Coin.Value), nil
}
