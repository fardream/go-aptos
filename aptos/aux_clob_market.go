package aptos

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

// AuxCritbit is a critbit tree, used to store order books.
// It has better tail behavior when adding/removing large number of elements,
// but on average performs worse than red/black tree.
type AuxCritbit struct {
	Entries  *TableWithLength `json:"entries,omitempty"`
	MaxIndex JsonUint64       `json:"max_index"`
	MinIndex JsonUint64       `json:"min_index"`
	Root     JsonUint64       `json:"root"`
	Tree     *TableWithLength `json:"tree"`
}

// AuxClobMarket contains two sided book of bids and asks.
type AuxClobMarket struct {
	Asks          *AuxCritbit `json:"asks,omitempty"`
	Bids          *AuxCritbit `json:"bids,omitempty"`
	BaseDecimals  uint8       `json:"base_decimals"`
	QuoteDecimals uint8       `json:"quote_decimals"`
	LotSize       JsonUint64  `json:"lot_size"`
	TickSize      JsonUint64  `json:"tick_size"`

	FillEvents   *EventHandler `json:"fill_events"`
	PlacedEvents *EventHandler `json:"placed_events"`
	CancelEvents *EventHandler `json:"cancel_events"`
}

// MarketType provides the market for a pair of currencies
func (info *AuxClientConfig) MarketType(baseCoin *MoveTypeTag, quoteCoin *MoveTypeTag) (*MoveTypeTag, error) {
	return NewMoveTypeTag(info.Address, "clob_market", "Market", []*MoveTypeTag{baseCoin, quoteCoin})
}

// AuxClobMarket_Level2Event_Level is price/quantity in an aux level 2 event
type AuxClobMarket_Level2Event_Level struct {
	Price    JsonUint64 `json:"price"`
	Quantity JsonUint64 `json:"quantity"`
}

// AuxClobMarket_Level2Event contains the bids and asks from a `load_market_into_event` call.
// Since tranversing the orderbook from off-chain is difficult, we run those entry functions to emit data into event queues.
type AuxClobMarket_Level2Event struct {
	Bids []*AuxClobMarket_Level2Event_Level `json:"bids"`
	Asks []*AuxClobMarket_Level2Event_Level `json:"asks"`
}

type AuxClobMarket_OpenOrderEventInfo struct {
	Id                Uint128    `json:"id"`
	CilentOrderId     Uint128    `json:"client_order_id"`
	Price             JsonUint64 `json:"price"`
	Quantity          JsonUint64 `json:"quantity"`
	AuxAuToBurnPerLot JsonUint64 `json:"aux_au_to_burn_per_lot"`
	IsBid             bool       `json:"is_bid"`
	OwnerId           Address    `json:"owner_id"`
	TimeoutTimestamp  JsonUint64 `json:"timeout_timestsamp"`
	OrderType         JsonUint64 `json:"order_type"`
	Timestamp         JsonUint64 `json:"timestamp"`
}

type AuxClobMarket_AllOrdersEvent struct {
	Bids [][]*AuxClobMarket_OpenOrderEventInfo `json:"bids"`
	Asks [][]*AuxClobMarket_OpenOrderEventInfo `json:"asks"`
}

type AuxClobMarket_OrderPlacedEvent struct {
	OrderId       Uint128    `json:"order_id"`
	ClientOrderId Uint128    `json:"client_order_id"`
	Owner         Address    `json:"owner"`
	IsBid         bool       `json:"is_bid"`
	Quantity      JsonUint64 `json:"qty"`
	Price         JsonUint64 `json:"price"`
	Timestamp     JsonUint64 `json:"timestamp"`
}

type AuxClobMarket_OrderCancelEvent struct {
	OrderId        Uint128    `json:"order_id"`
	ClientOrderId  Uint128    `json:"client_order_id"`
	Owner          Address    `json:"owner"`
	CancelQuantity JsonUint64 `json:"cancel_qty"`
	Timestamp      JsonUint64 `json:"timestamp"`
}

type AuxClobMarket_OrderFillEvent struct {
	OrderId           Uint128    `json:"order_id"`
	ClientOrderId     Uint128    `json:"client_order_id"`
	Owner             Address    `json:"owner"`
	IsBid             bool       `json:"is_bid"`
	BaseQuantity      JsonUint64 `json:"base_qty"`
	Price             JsonUint64 `json:"price"`
	Fee               JsonUint64 `json:"fee"`
	Rebate            JsonUint64 `json:"rebate"`
	RemainingQuantity JsonUint64 `json:"remaining_qty"`
	Timestamp         JsonUint64 `json:"timestamp"`
}

// AuxClobMarketOrderEvent is an union of all order events
//   - [AuxClobMarket_OrderFillEvent]
//   - [AuxClobMarket_OrderCancelEvent]
//   - [AuxClobMarket_OrderPlacedEvent]
type AuxClobMarketOrderEvent struct {
	*AuxClobMarket_OrderFillEvent
	*AuxClobMarket_OrderCancelEvent
	*AuxClobMarket_OrderPlacedEvent
	// the event is identified, but failed to parse
	ParsingFailure error
}

// IsOrderEvent checks if this AuxClobMarketOrderEvent is order related or not.
func (ev *AuxClobMarketOrderEvent) IsOrderEvent() bool {
	return ev.AuxClobMarket_OrderCancelEvent != nil ||
		ev.AuxClobMarket_OrderFillEvent != nil ||
		ev.AuxClobMarket_OrderPlacedEvent != nil ||
		ev.ParsingFailure != nil
}

// FilterAuxClobMarketOrderEvent filers out the clob market events emitted during the placing/cancelling process.
// It will parse out all the order related events. However, the output slice will have the same length of the input slice.
// The events unrelated to the orders will have [AuxClobMarketOrderEvent.IsOrderEvent] returns false. This is useful when the events needs
// to be put into context of composition with other protocols.
func FilterAuxClobMarketOrderEvent(events []*RawEvent, moduleAddress Address, ignoreAddress bool, dropOtherEvents bool) []*AuxClobMarketOrderEvent {
	var result []*AuxClobMarketOrderEvent

	for _, rawEvent := range events {
		ev := new(AuxClobMarketOrderEvent)
		switch {
		case rawEvent.Type.Module != AuxClobMarketModuleName:
		case !ignoreAddress && !cmp.Equal(moduleAddress, rawEvent.Type.Address):
		case rawEvent.Type.Name == "OrderPlacedEvent":
			placedEvent := new(AuxClobMarket_OrderPlacedEvent)
			if err := json.Unmarshal(*rawEvent.Data, placedEvent); err != nil {
				ev.ParsingFailure = err
			} else {
				ev.AuxClobMarket_OrderPlacedEvent = placedEvent
			}
		case rawEvent.Type.Name == "OrderCancelEvent":
			cancelEvent := new(AuxClobMarket_OrderCancelEvent)
			if err := json.Unmarshal(*rawEvent.Data, cancelEvent); err != nil {
				ev.ParsingFailure = err
			} else {
				ev.AuxClobMarket_OrderCancelEvent = cancelEvent
			}
		case rawEvent.Type.Name == "OrderFillEvent":
			fillEvent := new(AuxClobMarket_OrderFillEvent)
			if err := json.Unmarshal(*rawEvent.Data, fillEvent); err != nil {
				ev.ParsingFailure = err
			} else {
				ev.AuxClobMarket_OrderFillEvent = fillEvent
			}
		default:
		}

		if !dropOtherEvents || ev.IsOrderEvent() {
			result = append(result, ev)
		}
	}
	return result
}

//go:generate stringer -type AuxClobMarketSelfTradeType -linecomment

// AuxClobMarketSelfTradeType gives instruction on how to handle self trade
type AuxClobMarketSelfTradeType uint64

const (
	// cancel the order that is on the book
	AuxClobMarketSelfTradeType_CancelPassive AuxClobMarketSelfTradeType = iota + 200 // CANCEL_PASSIVE
	// cancel the order that is being placed.
	AuxClobMarketSelfTradeType_CancelAggressive // CANCEL_AGGRESSIVE
	// cancel both.
	AuxClobMarketSelfTradeType_CancelBoth // CANCEL_BOTH
)

//go:generate stringer -type AuxClobMarketOrderType -linecomment

// AuxClobMarketOrderType, can be limit, fok, ioc, post only or passive join
type AuxClobMarketOrderType uint64

const (
	AuxClobMarketOrderType_Limit        AuxClobMarketOrderType = iota + 100 // LIMIT
	AuxClobMarketOrderType_FOK                                              // FOK
	AuxClobMarketOrderType_IOC                                              // IOC
	AuxClobMarketOrderType_POST_ONLY                                        // POST_ONLY
	AuxClobMarketOrderType_PASSIVE_JOIN                                     // PASSIVE_JOIN
)

// AuxClobMarketModuleName is the module name for clob market.
const AuxClobMarketModuleName = "clob_market"

// ClobMarket_PlaceOrder creates a transaction to place an orde on aux.echange.
//
// Each order placed on the clob will receive an order id, even if it is cancelled or filled immediately. The order id is unique to the market,
// which is specified by base coin - quote coin pair.
//
// To link an order id generate from the contract on the client side, user can pass in a clientOrderId, which is
// unsigned int128 (which go doesn't support). However, the contract doesn't check of uniqueness of clientOrderIds.
//
// Limit price of the order must be in the quote coin decimals for one unit of base coin, and quantity is specified
// in base coin quantity. For example, assume coin Base has a decimal of 8, and coin Quote has a decimal of 6. To buy 0.5 unit of
// base at a price of 66.8, the limit price should be 66,800,000, and the quantity should be 50,000,000.
//
// also see contract at [here]
//
// [here]: https://github.com/aux-exchange/aux-exchange/blob/17f1b0ff677e93170610f21c088b347e18c2694d/aptos/contract/aux/sources/clob_market.move#L379-L393
func (info *AuxClientConfig) ClobMarket_PlaceOrder(
	sender Address,
	isBid bool,
	baseCoin,
	quoteCoin *MoveTypeTag,
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
) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxClobMarketModuleName, "place_order")
	payload := NewEntryFunctionPayload(function, []*MoveTypeTag{baseCoin, quoteCoin}, []EntryFunctionArg{
		// sender: &signer, // sender is the user who initiates the trade (can also be the vault_account_owner itself) on behalf of vault_account_owner. Will only succeed if sender is the creator of the account, or on the access control list of the account published under vault_account_owner address
		sender,                                     // vault_account_owner: address, // vault_account_owner is, from the module's internal perspective, the address that actually makes the trade. It will be the actual account that has changes in balance (fee, volume tracker, etc is all associated with vault_account_owner, and independent of sender (i.e. delegatee))
		EntryFunctionArg_Bool(isBid),               // is_bid: bool,
		JsonUint64(limitPrice),                     // limit_price: u64,
		JsonUint64(quantity),                       // quantity: u64,
		JsonUint64(auxToBurnPerLot),                // aux_au_to_burn_per_lot: u64,
		clientOrderId,                              // client_order_id: u128,
		JsonUint64(orderType),                      // order_type: u64,
		JsonUint64(ticksToSlide),                   // ticks_to_slide: u64, // # of ticks to slide for post only
		EntryFunctionArg_Bool(directionAggressive), // direction_aggressive: bool, // only used in passive join order
		JsonUint64(timeoutTimestamp),               // timeout_timestamp: u64, // if by the timeout_timestamp the submitted order is not filled, then it would be cancelled automatically, if the timeout_timestamp <= current_timestamp, the order would not be placed and cancelled immediately
		JsonUint64(selfTradeType),                  // self_trade_action_type: u64 // self_trade_action_type
	})

	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// ClobMarket_LoadMarketIntoEvent constructs a transaction to load price level and total quantities of each price level into an event.
// This is useful if the price/quantity of the market is needed since the market is stored on [TableWithLength] and is cumbersome to query.
func (info *AuxClientConfig) ClobMarket_LoadMarketIntoEvent(baseCoin, quoteCoin *MoveTypeTag, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxClobMarketModuleName, "load_market_into_event")
	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{baseCoin, quoteCoin},
		[]EntryFunctionArg{},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	if tx.Sender.IsZero() {
		tx.Sender = info.DataFeedAddress
	}

	return tx
}

func (info *AuxClientConfig) ClobMarket_LoadAllOrdersIntoEvent(baseCoin, quoteCoin *MoveTypeTag, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxClobMarketModuleName, "load_all_orders_into_event")
	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{baseCoin, quoteCoin},
		[]EntryFunctionArg{},
	)

	tx := &Transaction{
		Payload: payload,
	}

	ApplyTransactionOptions(tx, options...)

	if tx.Sender.IsZero() {
		tx.Sender = info.DataFeedAddress
	}

	return tx
}

// ClobMarket_CreateMarket constructs a transaction to create a market.
//
// lot size and quote size must guarantee that the minimal quote coin quantity is available and no rounding happens.
// This requires (assuming base coin has decimal of b)
//   - lot size * tick size / 10^b > 0  (the minimal quote coin quantity must be greater than zero)
//   - lot size * tick size % 10^b == 0 (the minimal quote coin quantity must be whole integers)
func (info *AuxClientConfig) ClobMarket_CreateMarket(sender Address, baseCoin, quoteCoin *MoveTypeTag, lotSize, tickSize uint64, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxClobMarketModuleName, "create_market")
	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{baseCoin, quoteCoin},
		[]EntryFunctionArg{
			JsonUint64(lotSize),
			JsonUint64(tickSize),
		},
	)

	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// ClobMarket_CancelAll constructs a transaction to cancel all open orders on a given market.
func (info *AuxClientConfig) ClobMarket_CancelAll(sender Address, baseCoin, quoteCoin *MoveTypeTag, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxClobMarketModuleName, "cancel_all")
	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{baseCoin, quoteCoin},
		[]EntryFunctionArg{
			sender,
		},
	)

	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// ClobMarket_CancelOrder constructs a transaction to cancel an open orde on a given market.
func (info *AuxClientConfig) ClobMarket_CancelOrder(sender Address, baseCoin, quoteCoin *MoveTypeTag, orderId Uint128, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxClobMarketModuleName, "cancel_order")
	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{baseCoin, quoteCoin},
		[]EntryFunctionArg{sender, orderId},
	)

	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

func (client *AuxClient) GetClobMarket(ctx context.Context, baseCoin, quoteCoin *MoveTypeTag, ledgerVersion uint64) (*AuxClobMarket, error) {
	marketType, err := client.config.MarketType(baseCoin, quoteCoin)
	if err != nil {
		return nil, err
	}
	return GetAccountResourceWithType[AuxClobMarket](ctx, client.client, client.config.Address, marketType, ledgerVersion)
}

func (client *AuxClient) ListLevel2(ctx context.Context, baseCoin, quoteCoin *MoveTypeTag, options ...TransactionOption) (*AuxClobMarket_Level2Event, error) {
	tx := client.config.ClobMarket_LoadMarketIntoEvent(baseCoin, quoteCoin, options...)
	if err := client.client.FillTransactionData(ctx, tx, false); err != nil {
		return nil, err
	}
	resp, err := client.client.SimulateTransaction(context.Background(), &SimulateTransactionRequest{
		Transaction: tx,
		Signature:   NewSingleSignatureForSimulation(&client.config.DataFeedPublicKey),
	})
	if err != nil {
		return nil, err
	}

	parsed := *resp.Parsed
	if len(parsed) != 1 {
		return nil, fmt.Errorf("more than one transactions in response: %s", string(resp.RawData))
	}
	if !parsed[0].Success {
		return nil, fmt.Errorf("simulation failed: %s, resp is %s", parsed[0].VmStatus, string(resp.RawData))
	}
	events := parsed[0].Events
	if len(events) != 1 {
		return nil, fmt.Errorf("there should only be one events: %#v", events)
	}

	var level2 AuxClobMarket_Level2Event
	if err := json.Unmarshal(*events[0].Data, &level2); err != nil {
		return nil, err
	}

	return &level2, nil
}

func (client *AuxClient) ListAllOrders(ctx context.Context, baseCoin, quoteCoin *MoveTypeTag, options ...TransactionOption) (*AuxClobMarket_AllOrdersEvent, error) {
	tx := client.config.ClobMarket_LoadAllOrdersIntoEvent(baseCoin, quoteCoin, options...)
	if err := client.client.FillTransactionData(ctx, tx, false); err != nil {
		return nil, err
	}
	resp, err := client.client.SimulateTransaction(context.Background(), &SimulateTransactionRequest{
		Transaction: tx,
		Signature:   NewSingleSignatureForSimulation(&client.config.DataFeedPublicKey),
	})
	if err != nil {
		return nil, err
	}

	parsed := *resp.Parsed
	if len(parsed) != 1 {
		return nil, fmt.Errorf("more than one transactions in response: %s", string(resp.RawData))
	}
	if !parsed[0].Success {
		return nil, fmt.Errorf("simulation failed: %s, resp is %s", parsed[0].VmStatus, string(resp.RawData))
	}
	events := parsed[0].Events
	if len(events) != 1 {
		return nil, fmt.Errorf("there should only be one events: %#v", events)
	}

	var allOrders AuxClobMarket_AllOrdersEvent
	if err := json.Unmarshal(*events[0].Data, &allOrders); err != nil {
		return nil, err
	}

	return &allOrders, nil
}
