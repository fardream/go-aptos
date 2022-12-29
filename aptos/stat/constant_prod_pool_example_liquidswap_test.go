package stat_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-aptos/aptos/known"
	"github.com/fardream/go-aptos/aptos/stat"
)

var liquidSwapUsdSymbols = []string{
	"USDC",
	"USDCso",
	"USDT",
	"ceUSDC",
	"ceDAI",
	"ceUSDT",
	"zUSDC",
	"zUSDT",
	"BUSD",
	"ceBUSD",
	"zBUSD",
}

type LiquidSwapLiquidityPool struct {
	CoinXReserve aptos.Coin `json:"coin_x_reserve"`
	CoinYReserve aptos.Coin `json:"coin_y_reserve"`
}

// ExampleConstantProductPool shows an example of how to query the constant product pool of a protocol on chain (for aux.exchange)
func ExampleConstantProductPool_liquidswap() {
	client, _ := aptos.NewClient(aptos.Mainnet, "")

	resp, err := client.GetAccountResources(context.Background(), &aptos.GetAccountResourcesRequest{
		Address: aptos.MustParseAddress("0x05a97986a9d031c4567e15b797be516910cfcb4156312482efc6a19c0a30c948"),
	})
	if err != nil {
		panic(err)
	}
	protocol := stat.NewStatForConstantProductPool()

	known.ReloadHippoCoinRegistry(known.HippoCoinRegistryUrl)

	for _, usdSymbol := range liquidSwapUsdSymbols {
		stable := known.GetCoinInfoBySymbol(aptos.Mainnet, usdSymbol)
		if stable != nil {
			protocol.AddStableCoins(stable.TokenType.Type)
		}
	}

	for _, resource := range *resp.Parsed {
		resourceType := resource.Type

		if resourceType.Module == "liquidity_pool" && resourceType.Name == "LiquidityPool" {
			var amm LiquidSwapLiquidityPool
			if err := json.Unmarshal(resource.Data, &amm); err != nil {
				fmt.Printf("failed to parse %s: %s\n", string(resource.Data), err.Error())
				continue
			}
			protocol.AddSinglePool(resourceType.GenericTypeParameters[0].Struct, uint64(amm.CoinXReserve.Value), resourceType.GenericTypeParameters[1].Struct, uint64(amm.CoinYReserve.Value))
		}
	}

	protocol.FillCoinInfo(context.Background(), aptos.Mainnet, client)

	protocol.FillStat()

	var coinBuf bytes.Buffer

	fmt.Fprintln(&coinBuf, "Coin Type, Coin Symbol, Coin Name, Coin Decimal, Total Reserve, Price, Total Value, IsHippo")

	for _, coin := range protocol.Coins {
		isHippo := 0
		if coin.IsHippo {
			isHippo = 1
		}

		fmt.Fprintf(&coinBuf, "%s,%s,%s,%d,%d,%f,%f,%d\n", coin.MoveTypeTag.String(), coin.Symbol, coin.Name, coin.Decimals, coin.TotalQuantity, coin.Price, coin.TotalValue, isHippo)
	}

	fmt.Fprint(os.Stderr, coinBuf.String())

	var poolBuf bytes.Buffer

	fmt.Fprint(&poolBuf, "Coin 0 Type, Coin 0 Symbol, Coin 0 Is Hippo, Coin 0 Decimal, Coin 0 Reserve, Coin 0 Price, Coin 0 Value,")
	fmt.Fprint(&poolBuf, "Coin 1 Type, Coin 1 Symbol, Coin 1 Is Hippo, Coin 1 Decimal, Coin 1 Reserve, Coin 1 Price, Coin 1 Value,")
	fmt.Fprintln(&poolBuf, "Total Value")

	for _, pool := range protocol.Pools {
		coin0Name := pool.Coin0.String()
		coin0Info := protocol.Coins[coin0Name]
		coin0IsHippo := 0
		if coin0Info.IsHippo {
			coin0IsHippo = 1
		}
		fmt.Fprintf(&poolBuf, "%s,%s,%d,%d,%d,%f,%f,", coin0Name, coin0Info.Symbol, coin0IsHippo, coin0Info.Decimals, pool.Coin0Reserve, coin0Info.Price, pool.Coin0Value)

		coin1Name := pool.Coin1.String()
		coin1Info := protocol.Coins[coin1Name]
		coin1IsHippo := 0
		if coin1Info.IsHippo {
			coin1IsHippo = 1
		}
		fmt.Fprintf(&poolBuf, "%s,%s,%d,%d,%d,%f,%f,", coin1Name, coin1Info.Symbol, coin1IsHippo, coin1Info.Decimals, pool.Coin1Reserve, coin1Info.Price, pool.Coin1Value)

		fmt.Fprintf(&poolBuf, "%f\n", pool.TotalValueLocked)
	}

	fmt.Fprint(os.Stderr, poolBuf.String())

	// Output:
}
