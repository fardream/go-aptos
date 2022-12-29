package stat

import (
	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-aptos/aptos/known"
)

// CoinStat contains the stat for a coin locked into a protocol.
type CoinStat struct {
	// MoveTypeTag is the canonical identifier for the coin.
	MoveTypeTag *aptos.MoveStructTag
	// Decimals of the coin.
	Decimals uint8
	// CoinRegistry is the Hippo Coin Registry
	CoinRegistry *known.HippoCoinRegistryEntry
	// IsStable tells if the coin is stable or not. This is provided by the user.
	IsStable bool
	// Name of the coin. If CoinRegistry is not nil, this is from the registry, otherwise it's from the name on chain with (Not Hippo) appended.
	Name string
	// Symbol of the coin. If CoinRegistry is not nil, this is from the registry, otherwise it's from the symbol on chain with (Not Hippo) appended.
	Symbol string
	// If the information is from hippo or not
	IsHippo bool

	// Price of the coin
	Price float64
	// TotalQuantity is the total quantity of the coin.
	TotalQuantity uint64
	// TotalValue of the coin.
	TotalValue float64
}
