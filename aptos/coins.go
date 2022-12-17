package aptos

import (
	"context"
)

// AptosCoin is the type for aptos coin
var AptosCoin = MoveStructTag{
	MoveModuleTag: MoveModuleTag{
		Address: AptosStdAddress,
		Module:  "aptos_coin",
	},
	Name: "AptosCoin",
}

// GetCoinStoreType returns the 0x1::coin::CoinStore<T>
func GetCoinStoreType(coin *MoveStructTag) *MoveStructTag {
	return &MoveStructTag{
		MoveModuleTag: MoveModuleTag{
			Address: AptosStdAddress,
			Module:  "coin",
		},
		Name: "CoinStore",

		GenericTypeParameters: []*MoveTypeTag{
			{Struct: coin},
		},
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
func (client *Client) GetCoinBalance(ctx context.Context, address Address, coinType *MoveStructTag) (uint64, error) {
	coinStore, err := GetAccountResourceWithType[CoinStore](ctx, client, address, GetCoinStoreType(coinType), 0)
	if err != nil {
		return 0, err
	}

	return uint64(coinStore.Coin.Value), nil
}

// CoinInfo, this is golang equivalentg of 0x1::coin::CoinInfo.
type CoinInfo struct {
	Decimals uint8  `json:"decimals"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
}

// GetCoinInfoType returns the CoinInfo<coinType> for coin.
func GetCoinInfoType(coinType *MoveStructTag) *MoveStructTag {
	return &MoveStructTag{
		MoveModuleTag: MoveModuleTag{
			Address: AptosStdAddress,
			Module:  "coin",
		},
		GenericTypeParameters: []*MoveTypeTag{{Struct: coinType}},
		Name:                  "CoinInfo",
	}
}

// GetCoinInfo retrieves [CoinInfo]
func (client *Client) GetCoinInfo(ctx context.Context, coinType *MoveStructTag) (*CoinInfo, error) {
	return GetAccountResourceWithType[CoinInfo](ctx, client, coinType.Address, GetCoinInfoType(coinType), 0)
}
