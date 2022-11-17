package aptos_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-bcs/bcs"
	"github.com/gorilla/websocket"
)

func must[T any](in T, err error) T {
	orPanic(err)

	return in
}

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}

type Order struct {
	OrderId   string `json:"order_id"`
	Price     string `json:"price"`
	OrderType string `json:"order_type"`
	Type      string `json:"type"`
	Size      string `json:"size"`
	Side      string `json:"side"`
}

type TraderWithOrdreIds struct {
	*aptos.AuxClobMarketTrader

	orderIds []string
}

// Example creating 5 traders, each maintain 3 orders at most
func ExampleAuxClobMarketTrader() {
	const network = aptos.Devnet

	const (
		desiredAptosBalance = 1_000_000_000_000     // 10_000 ATP
		desiredUSDCBalance  = 1_000_000_000_000_000 // 1_000_000_000 USDC
		desiredAuxBalance   = 1_000_000_000_000_000 //
	)

	auxConfig, _ := aptos.GetAuxClientConfig(network)

	restUrl, faucetUrl, _ := aptos.GetDefaultEndpoint(network)
	auxFakeCoinCoin := must(aptos.GetAuxFakeCoinCoinType(auxConfig.Address, aptos.AuxFakeCoin_AUX))
	usdcFakeCoinCoin := must(aptos.GetAuxFakeCoinCoinType(auxConfig.Address, aptos.AuxFakeCoin_USDC))

	// cancel the whole process after 15 minutes
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*15)
	defer cancel()

	client := aptos.MustNewClient(network, restUrl)

	traders := make([]*TraderWithOrdreIds, 0, 5)
	for _, account := range []*aptos.LocalAccount{
		must(aptos.NewLocalAccountWithRandomKey()),
		must(aptos.NewLocalAccountWithRandomKey()),
		must(aptos.NewLocalAccountWithRandomKey()),
		must(aptos.NewLocalAccountWithRandomKey()),
		must(aptos.NewLocalAccountWithRandomKey()),
	} { // get some gas
		txes := must(aptos.RequestFromFaucet(ctx, faucetUrl, &(account.Address), desiredAptosBalance*2))
		for _, txhash := range txes {
			client.WaitForTransaction(ctx, txhash)
		}

		// create user account
		must(client.SignSubmitTransactionWait(ctx, account, auxConfig.Vault_CreateAuxAccount(account.Address), false))
		<-time.After(5 * time.Second)
		// request fake coins
		must(client.SignSubmitTransactionWait(ctx, account, auxConfig.FakeCoin_RegisterAndMint(account.Address, aptos.AuxFakeCoin_USDC, desiredUSDCBalance), false))
		<-time.After(5 * time.Second)
		must(client.SignSubmitTransactionWait(ctx, account, auxConfig.FakeCoin_RegisterAndMint(account.Address, aptos.AuxFakeCoin_AUX, desiredAuxBalance), false))
		<-time.After(5 * time.Second)

		// deposit fake coins
		must(client.SignSubmitTransactionWait(ctx, account, auxConfig.Vault_Deposit(account.Address, account.Address, usdcFakeCoinCoin, desiredUSDCBalance), false))
		<-time.After(5 * time.Second)
		must(client.SignSubmitTransactionWait(ctx, account, auxConfig.Vault_Deposit(account.Address, account.Address, auxFakeCoinCoin, desiredAuxBalance), false))
		<-time.After(5 * time.Second)

		traders = append(traders,
			&TraderWithOrdreIds{
				AuxClobMarketTrader: must(
					aptos.NewAuxClobMarketTrader(
						ctx,
						aptos.NewAuxClient(client, auxConfig, account),
						auxFakeCoinCoin,
						usdcFakeCoinCoin),
				),
			})
	}

	// connect to coinbase
	asset := "APT-USD"
	endpoint := "wss://ws-feed.exchange.coinbase.com"
	dialer := &websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
	}
	conn, rsp, err := dialer.DialContext(ctx, endpoint, nil)
	if err != nil {
		orPanic(fmt.Errorf("failed to open connection: %v %v", err, rsp))
	}
	defer conn.Close()
	orPanic(conn.WriteJSON(map[string]any{
		"type":        "subscribe",
		"product_ids": []string{asset},
		"channels":    []string{"full"},
	}))

	var wg sync.WaitGroup
	defer wg.Wait()

	waitForWs := make(chan struct{})
	orderChannel := make(chan *Order, 100)

	// first goroutine read the data from websocket and pipe it into a channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			// once the websocket is disconnected, we indicate that we are done.
			waitForWs <- struct{}{}
			close(waitForWs)
			close(orderChannel)
		}()
	readLoop:
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("failed to read websocket...: %v", err)
				break
			}

			var order Order
			if err := json.Unmarshal(data, &order); err != nil {
				fmt.Printf("failed to parse: %v\n", err)
				continue
			}

			if !(order.Type == "received" && order.OrderType == "limit") {
				continue
			}
			// stop piping order if cancelled
			select {
			case orderChannel <- &order:
			case <-ctx.Done():
				break readLoop
			}

		}
	}()

	// a second websocket will read from the channel,
	// and select the next trader to trade.
	// each trader will wait 30 seconds to avoid flooding the fullnode.
	wg.Add(1)
	go func() {
		defer wg.Done()
		buyId := 0
		sellId := 1
		wait := time.Second * 5
		var clienOrderId uint64 = 1
	orderLoop:
		for {
			var order *Order
			var ok bool

			// make sure we don't hang on orderChannel if ctx is cancelled
			select {
			case order, ok = <-orderChannel:
			case <-ctx.Done():
				break orderLoop
			}

			if !ok {
				break
			}

			// stop waiting if cancelled
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				break orderLoop
			}

			price, err := strconv.ParseFloat(order.Price, 64)
			if err != nil {
				fmt.Printf("failed to parse price: %s %v\n", order.Price, err)
			}
			size, err := strconv.ParseFloat(order.Size, 64)
			if err != nil {
				fmt.Printf("failed to parse size: %s %v\n", order.Size, err)
			}
			priceInt := uint64(price * 1_000_000)
			priceInt = (priceInt / 10000) * 10000
			sizeInt := uint64(size * 1_000_000)
			sizeInt = (sizeInt / 100000) * 100000

			var trader *TraderWithOrdreIds
			if order.Side == "buy" {
				buyId %= len(traders)
				trader = traders[buyId]
				buyId += 2
				priceInt += 10000
			} else {
				sellId %= len(traders)
				trader = traders[sellId]
				sellId += 2
			}

			fmt.Printf("place %s size %d price %d\n", order.Side, sizeInt, priceInt)

			result, err := trader.PlaceOrder(
				ctx,
				order.Side == "buy",
				priceInt,
				sizeInt,
				0,
				*bcs.NewUint128FromUint64(clienOrderId, 0),
				aptos.AuxClobMarketOrderType_Limit,
				0,
				false,
				math.MaxInt64,
				aptos.AuxClobMarketSelfTradeType_CancelBoth,
				aptos.TransactionOption_MaxGasAmount(30000),
			)
			// print out the order
			if err != nil {
				spew.Dump(err)
			} else {
				fmt.Println(string(result.RawTransaction.Hash))
				if result.OrderStatus == aptos.AuxClobMarketOrderStatus_Placed {
					trader.orderIds = append(trader.orderIds, result.OrderId.Big().String())
					if len(trader.orderIds) >= 3 {
						idToCancel := must(bcs.NewUint128(trader.orderIds[0]))
						trader.orderIds = trader.orderIds[1:]
						result, err := trader.CancelOrder(ctx, *idToCancel)
						if err != nil {
							spew.Dump(err)
						} else {
							spew.Dump(result)
						}
					}
				}
			}
		}
	}()

	select {
	case <-waitForWs:
	case <-ctx.Done():
		conn.Close()
		<-waitForWs
	}
}
