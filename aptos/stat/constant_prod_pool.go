package stat

import (
	"context"
	"fmt"
	"math"

	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-aptos/aptos/known"
)

// ConstantProductPool contains stats of the constant product pool.
type ConstantProductPool struct {
	Coin0        *aptos.MoveStructTag
	Coin0Reserve uint64
	Coin0Value   float64

	Coin1        *aptos.MoveStructTag
	Coin1Reserve uint64
	Coin1Value   float64

	TotalValueLocked float64
}

// ConstantProductPoolProtocol contains stat for all pools on a protocol
type ConstantProductPoolProtocol struct {
	Coins       map[string]*CoinStat
	StableCoins []*aptos.MoveStructTag
	Pools       map[string]*ConstantProductPool

	TotalValueLocked float64
}

// NewStatForConstantProductPool creates a new struct to hold the stat.
func NewStatForConstantProductPool() *ConstantProductPoolProtocol {
	return &ConstantProductPoolProtocol{
		Coins:       make(map[string]*CoinStat),
		StableCoins: make([]*aptos.MoveStructTag, 0),
		Pools:       make(map[string]*ConstantProductPool),
	}
}

// AddStableCoins add stable coins to the protocol stat
func (p *ConstantProductPoolProtocol) AddStableCoins(stableCoins ...*aptos.MoveStructTag) {
	p.StableCoins = append(p.StableCoins, stableCoins...)
	p.AddCoins(stableCoins...)
}

// AddPools add pools to the stat
func (p *ConstantProductPoolProtocol) AddPools(pools ...*ConstantProductPool) error {
	for _, pool := range pools {
		if pool.Coin0 == nil {
			return fmt.Errorf("coin 0 is nil")
		}
		if pool.Coin1 == nil {
			return fmt.Errorf("coin 1 is nil")
		}
		poolName := fmt.Sprintf("%s-%s", pool.Coin0.String(), pool.Coin1.String())
		if _, found := p.Pools[poolName]; found {
			return fmt.Errorf("%s already in the pools", poolName)
		}

		p.Pools[poolName] = pool
		p.AddCoins(pool.Coin0, pool.Coin1)
	}

	return nil
}

// AddSinglePool add a single pool to the stat of the pool
func (p *ConstantProductPoolProtocol) AddSinglePool(coin0 *aptos.MoveStructTag, reserve0 uint64, coin1 *aptos.MoveStructTag, reserve1 uint64) error {
	if coin0 == nil {
		return fmt.Errorf("coin 0 is nil")
	}
	if coin1 == nil {
		return fmt.Errorf("coin 0 is nil")
	}
	poolName := fmt.Sprintf("%s-%s", coin0.String(), coin1.String())
	if _, found := p.Pools[poolName]; found {
		return fmt.Errorf("%s already in the pools", poolName)
	}

	p.Pools[poolName] = &ConstantProductPool{
		Coin0:        coin0,
		Coin0Reserve: reserve0,
		Coin1:        coin1,
		Coin1Reserve: reserve1,
	}
	p.AddCoins(coin0, coin1)

	return nil
}

// AddCoins add coins to the pool.
func (p *ConstantProductPoolProtocol) AddCoins(coins ...*aptos.MoveStructTag) {
	for _, coin := range coins {
		coinName := coin.String()
		if _, found := p.Coins[coinName]; found {
			continue
		}

		p.Coins[coinName] = &CoinStat{
			MoveTypeTag:   coin,
			Decimals:      0,
			CoinRegistry:  nil,
			Price:         0,
			TotalValue:    0,
			TotalQuantity: 0,
		}
	}
}

// FillCoinInfo fills the missing coins.
// Consider calling [known.ReloadHippoCoinRegistry] before hand.
// If all coin infos are available from Hippo Registry, no query will be made to the aptos network.
func (p *ConstantProductPoolProtocol) FillCoinInfo(ctx context.Context, network aptos.Network, aptosClient *aptos.Client) error {
	for _, coinInfo := range p.Coins {
		if coinInfo.CoinRegistry != nil {
			continue
		}
		registry := known.GetCoinInfo(network, coinInfo.MoveTypeTag)
		if registry != nil {
			coinInfo.CoinRegistry = registry
			coinInfo.Name = registry.Name
			coinInfo.Symbol = registry.Symbol
			coinInfo.Decimals = registry.Decimals
			coinInfo.IsHippo = true
		} else {
			info, err := aptosClient.GetCoinInfo(ctx, coinInfo.MoveTypeTag)
			if err != nil {
				return err
			}

			coinInfo.Decimals = info.Decimals
			coinInfo.Symbol = info.Symbol + "(Non Hippo)"
			coinInfo.Name = info.Name + "(Non Hippo)"
			coinInfo.IsHippo = false
		}
	}

	return nil
}

// FillStat fills the stat.
// It assumes the stat of the coins have already been filled.
//
//  1. Set all stable coins' price to 1.
//  2. Walk through the pool list. If one of the coin price is known, the unknown coin price will be known_price * known_coin_quantity / unknown_coin_quantity.
//  3. Rewalk the pool until all coin price are filled, or 256 iterations are reached.
//  4. Calculate TVL for pools and totals.
func (p *ConstantProductPoolProtocol) FillStat() {
	for _, stableCoin := range p.StableCoins {
		name := stableCoin.String()
		if coin, found := p.Coins[name]; found {
			coin.IsStable = true
		}
	}

	missing := 0
	// first, fill reserve
	for _, coin := range p.Coins {
		coin.TotalQuantity = 0
		coin.TotalValue = 0

		if coin.IsStable {
			coin.Price = 1
		} else {
			coin.Price = 0
			missing += 1
		}
	}

	for i := 0; i < 256; i++ {
		for _, pool := range p.Pools {
			name0 := pool.Coin0.String()
			coin0Info := p.Coins[name0]
			name1 := pool.Coin1.String()
			coin1Info := p.Coins[name1]
			if coin0Info.Price == 0 {
				if coin1Info.Price != 0 && (i >= 1 || coin1Info.IsStable) {
					value0 := float64(pool.Coin0Reserve) / math.Pow10(int(coin0Info.Decimals))
					value1 := float64(pool.Coin1Reserve) / math.Pow10(int(coin1Info.Decimals))
					if value0 > 0 {
						coin0Info.Price = value1 / value0 * coin1Info.Price
						missing -= 1
					}
				}
			}

			if coin1Info.Price == 0 {
				if coin0Info.Price != 0 && (i >= 1 || coin0Info.IsStable) {
					value0 := float64(pool.Coin0Reserve) / math.Pow10(int(coin0Info.Decimals))
					value1 := float64(pool.Coin1Reserve) / math.Pow10(int(coin1Info.Decimals))
					if value1 > 0 {
						coin1Info.Price = value0 / value1 * coin0Info.Price
						missing -= 1
					}
				}
			}
		}
		if missing <= 0 {
			break
		}
	}

	for _, pool := range p.Pools {
		name0 := pool.Coin0.String()
		coin0Info := p.Coins[name0]
		d0 := math.Pow10(int(coin0Info.Decimals))
		coin0Info.TotalQuantity += pool.Coin0Reserve
		coin0Info.TotalValue = float64(coin0Info.TotalQuantity) * coin0Info.Price / d0
		pool.Coin0Value = float64(pool.Coin0Reserve) * coin0Info.Price / d0

		name1 := pool.Coin1.String()
		coin1Info := p.Coins[name1]
		d1 := math.Pow10(int(coin1Info.Decimals))
		coin1Info.TotalQuantity += pool.Coin1Reserve
		coin1Info.TotalValue = float64(coin1Info.TotalQuantity) * coin1Info.Price / d1
		pool.Coin1Value = float64(pool.Coin1Reserve) * coin1Info.Price / d1

		pool.TotalValueLocked = pool.Coin0Value + pool.Coin1Value

		p.TotalValueLocked += pool.TotalValueLocked
	}
}
