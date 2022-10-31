package aptos

import (
	"context"
	"fmt"
)

//go:generate stringer -type AuxClobMarketOrderStatus -linecomment

// AuxClobMarketOrderStatus contains the onchain status of the order
//   - Placed: the order is placed into the limit order book and is open.
//   - Cancelled: the order is cancelled. Order can be cancelled due to user sending a cancel transaction,
//     self trade handling of a new order, time out, or fok/ioc/post only/passive join failed to meet the condition.
//   - Filled: the order is fully filled.
type AuxClobMarketOrderStatus uint64

const (
	AuxClobMarketOrderStatus_Placed    AuxClobMarketOrderStatus = iota // Placed
	AuxClobMarketOrderStatus_Filled                                    // Filled
	AuxClobMarketOrderSTatus_Cancelled                                 // Cancelled
)

// AuxClobMarketTrader contains the market state, a client to aptos/aux,
type AuxClobMarketTrader struct {
	baseCoin  *MoveTypeTag
	quoteCoin *MoveTypeTag

	// OriginalState contains the state of the market when the AuxClobMarketTrader first start trading.
	// This is useful to monitor the event queues to check if an order is filled or cancelled.
	OriginalState *AuxClobMarket

	originalStart AuxClobMarketEventCounter
	current       AuxClobMarketEventCounter

	auxClient *AuxClient
}

type AuxClobMarketEventCounter struct {
	fillEventsStart   uint64
	placedEventsStart uint64
	cancelEventsStart uint64
}

func (counter *AuxClobMarketEventCounter) FillFromAuxClobMarket(market *AuxClobMarket) {
	counter.fillEventsStart = uint64(market.FillEvents.Counter)
	counter.cancelEventsStart = uint64(market.CancelEvents.Counter)
	counter.placedEventsStart = uint64(market.PlacedEvents.Counter)
}

// NewAuxClobMarketTrader creates a new trader.
func NewAuxClobMarketTrader(ctx context.Context, auxClient *AuxClient, baseCoin, quoteCoin *MoveTypeTag) (*AuxClobMarketTrader, error) {
	originalState, err := auxClient.GetClobMarket(ctx, baseCoin, quoteCoin, 0)
	if err != nil {
		return nil, err
	}

	r := &AuxClobMarketTrader{
		OriginalState: originalState,
		auxClient:     auxClient,
		baseCoin:      baseCoin,
		quoteCoin:     quoteCoin,
	}

	r.current.FillFromAuxClobMarket(originalState)
	r.originalStart.FillFromAuxClobMarket(originalState)

	return r, nil
}

// AuxClobMarketPlaceOrderResult contains the results from a [AuxClientConfig.ClobMarket_PlaceOrder] transaction.
//
// If the transaction is successfully committed to the blockchain, the order will get an order id even if it never goes onto the
// order book. Note there is no way to differentiate between an empty client order id or a 0 client order id.
//
// Status of the order indicates the status right after the transaction is committed.
type AuxClobMarketPlaceOrderResult struct {
	RawTransaction *TransactionWithInfo

	// OrderId for this order
	OrderId *Uint128
	// ClientOrderId if there is any. Otherwise this will be 0.
	ClientOrderId *Uint128
	// Status of the order
	OrderStatus AuxClobMarketOrderStatus

	Events []*AuxClobMarketOrderEvent

	FillEvents []*AuxClobMarket_OrderFillEvent
}

// AuxClobMarketPlaceOrderError is the error returned if the transaction is successfully submitted to the chain and executed,
// but post processing somehow failed. It will contain a copy of the [AuxClobMarketPlaceOrderResult] that is processed uptil failure.
type AuxClobMarketPlaceOrderError struct {
	ErroredResult        *AuxClobMarketPlaceOrderResult
	InnerError           error
	IsTransactionFailure bool
}

var _ error = (*AuxClobMarketPlaceOrderError)(nil)

func (err *AuxClobMarketPlaceOrderError) Error() string {
	if err.IsTransactionFailure {
		return fmt.Sprintf("tx hash: %s - failed, vm status: %s", err.ErroredResult.RawTransaction.Hash, err.ErroredResult.RawTransaction.VmStatus)
	} else {
		return fmt.Sprintf("tx hash: %s - error: %v", err.ErroredResult.RawTransaction.Hash, err.InnerError)
	}
}

// IsAuxClobMarketPlaceOrderError checks if the error is [AuxClobMarketPlaceOrderError],
// returns the casted [AuxClobMarketPlaceOrderError] and a bool indicate if it is [AuxClobMarketPlaceOrderError]
func IsAuxClobMarketPlaceOrderError(err error) (*AuxClobMarketPlaceOrderError, bool) {
	placeOrderErr, ok := err.(*AuxClobMarketPlaceOrderError)

	return placeOrderErr, ok
}

// PlaceOrder places a new order on clob, check [AuxClientConfig.ClobMarket_PlaceOrder] for more information on the parameters.
func (trader *AuxClobMarketTrader) PlaceOrder(
	ctx context.Context,
	isBid bool,
	limitPrice uint64,
	quantity uint64,
	auxToBurnPerLot uint64,
	clientOrderId Uint128,
	orderType AuxClobMarketOrderType,
	ticksToSlide uint64,
	directionAggressive bool,
	timeoutTimestamp uint64,
	selfTradeType AuxClobMarketSelfTradeType,
	options ...TransactionOption,
) (*AuxClobMarketPlaceOrderResult, error) {
	// Make sure in the input client order id is not zero.
	if clientOrderId.Cmp(&zeroUint128) != 0 {
		return nil, fmt.Errorf("must pass in non zero client order id to identify returning orders")
	}

	tx := trader.auxClient.config.ClobMarket_PlaceOrder(
		trader.auxClient.userAddress,
		isBid, trader.baseCoin,
		trader.quoteCoin,
		limitPrice, quantity,
		auxToBurnPerLot,
		clientOrderId,
		orderType,
		ticksToSlide,
		directionAggressive,
		timeoutTimestamp,
		selfTradeType,
		options...)

	if err := trader.auxClient.client.FillTransactionData(ctx, tx, false); err != nil {
		return nil, err
	}

	txInfo, err := trader.auxClient.client.SignSubmitTransactionWait(ctx, trader.auxClient.signer, tx, false)
	if err != nil {
		return nil, err
	}
	result := &AuxClobMarketPlaceOrderResult{
		RawTransaction: txInfo,
	}

	if !txInfo.Success {
		return nil, &AuxClobMarketPlaceOrderError{
			InnerError:           nil,
			ErroredResult:        result,
			IsTransactionFailure: true,
		}
	}

	result.Events = FilterAuxClobMarketOrderEvent(txInfo.Events, trader.auxClient.config.Address, false, true)
	result.ClientOrderId = &Uint128{}
	*result.ClientOrderId = clientOrderId

findOrderIdLoop:
	for _, ev := range result.Events {
		switch {
		case ev.AuxClobMarket_OrderPlacedEvent != nil:
			placedEvent := ev.AuxClobMarket_OrderPlacedEvent
			// if the order placed event is not nil, we got the order id, since this is the only order in this transaction.
			if clientOrderId.Cmp(&placedEvent.ClientOrderId) != 0 {
				continue findOrderIdLoop
			}
			result.OrderId = &Uint128{}
			*result.OrderId = placedEvent.OrderId

			break findOrderIdLoop
		case ev.AuxClobMarket_OrderFillEvent != nil:
			fill := ev.AuxClobMarket_OrderFillEvent
			if fill.ClientOrderId.Cmp(&clientOrderId) != 0 {
				continue findOrderIdLoop
			}
			if result.OrderId == nil {
				result.OrderId = &Uint128{}
				*result.OrderId = fill.OrderId
			}
			if fill.RemainingQuantity == 0 {
				result.OrderStatus = AuxClobMarketOrderStatus_Filled
				break findOrderIdLoop
			}
			result.FillEvents = append(result.FillEvents, fill)
		case ev.AuxClobMarket_OrderCancelEvent != nil:
			cancel := ev.AuxClobMarket_OrderCancelEvent
			if cancel.ClientOrderId.Cmp(&clientOrderId) != 0 {
				continue findOrderIdLoop
			}
			if result.OrderId == nil {
				result.OrderId = &Uint128{}
				*result.OrderId = cancel.OrderId
			}
			result.OrderStatus = AuxClobMarketOrderSTatus_Cancelled
			break findOrderIdLoop
		}
	}

	if result.OrderId == nil {
		return nil, &AuxClobMarketPlaceOrderError{
			ErroredResult: result,
			InnerError:    fmt.Errorf("tx hash: %s, failed to find the client order id in transaction events", txInfo.Hash),
		}
	}

	return result, nil
}

// AuxClobMarketCancelOrderResult contains the results from a [AuxClientConfig.ClobMarket_Cancel] transaction.
type AuxClobMarketCancelOrderResult struct {
	RawTransation *TransactionWithInfo
	IsCancelled   bool
}

// Cancel an order
func (trader *AuxClobMarketTrader) CancelOrder(
	ctx context.Context,
	orderId Uint128,
	options ...TransactionOption,
) (*AuxClobMarketCancelOrderResult, error) {
	tx := trader.auxClient.config.ClobMarket_CancelOrder(trader.auxClient.userAddress, trader.baseCoin, trader.quoteCoin, orderId, options...)

	if err := trader.auxClient.client.FillTransactionData(ctx, tx, false); err != nil {
		return nil, err
	}

	txInfo, err := trader.auxClient.client.SignSubmitTransactionWait(ctx, trader.auxClient.signer, tx, false)
	if err != nil {
		return nil, err
	}

	if txInfo.Success {
		return &AuxClobMarketCancelOrderResult{
			IsCancelled:   true,
			RawTransation: txInfo,
		}, nil
	} else {
		return &AuxClobMarketCancelOrderResult{
			IsCancelled:   false,
			RawTransation: txInfo,
		}, nil
	}
}

type AuxClobMarketCancelAllResult struct {
	RawTransation     *TransactionWithInfo
	CancelledOrderIds []Uint128
}

func (trader *AuxClobMarketTrader) CancelAll(ctx context.Context, options ...TransactionOption) (*AuxClobMarketCancelAllResult, error) {
	tx := trader.auxClient.config.ClobMarket_CancelAll(trader.auxClient.userAddress, trader.baseCoin, trader.quoteCoin, options...)

	if err := trader.auxClient.client.FillTransactionData(ctx, tx, false); err != nil {
		return nil, err
	}

	txInfo, err := trader.auxClient.client.SignSubmitTransactionWait(ctx, trader.auxClient.signer, tx, false)
	if err != nil {
		return nil, err
	}

	if txInfo.Success {
		result := &AuxClobMarketCancelAllResult{
			RawTransation: txInfo,
		}

		events := FilterAuxClobMarketOrderEvent(result.RawTransation.Events, trader.auxClient.config.Address, false, true)
		for _, ev := range events {
			if ev.AuxClobMarket_OrderCancelEvent != nil {
				result.CancelledOrderIds = append(result.CancelledOrderIds, ev.AuxClobMarket_OrderCancelEvent.OrderId)
			}
		}
		return result, nil
	} else {
		return &AuxClobMarketCancelAllResult{
			RawTransation: txInfo,
		}, nil
	}
}
