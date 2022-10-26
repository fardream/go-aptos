package aptos

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

// AuxAmmPool is a constant product amm
type AuxAmmPool struct {
	FeeBps   JsonUint64 `json:"feebps"`
	Frozen   bool       `json:"frozen"`
	XReserve Coin       `json:"x_reserve"`
	YReserve Coin       `json:"y_reserve"`

	AddLiquidityEvents    *EventHandler `json:"add_liquidity_events"`
	RemoveLiquidityEvents *EventHandler `json:"remove_liquidity_events"`
	SwapEvents            *EventHandler `json:"swap_events"`
}

// MarketType provides the market for a pair of currencies
func (info *AuxClientConfig) MarketType(baseCoin *MoveTypeTag, quoteCoin *MoveTypeTag) (*MoveTypeTag, error) {
	return NewMoveTypeTag(info.Address, "clob_market", "Market", []*MoveTypeTag{baseCoin, quoteCoin})
}

// AuxLevel2Event_Level is price/quantity in an aux level 2 event
type AuxLevel2Event_Level struct {
	Price    JsonUint64 `json:"price"`
	Quantity JsonUint64 `json:"quantity"`
}

// AuxLevel2Event contains the bids and asks from a `load_market_into_event` call.
// Since tranversing the orderbook from off-chain is difficult, we run those entry functions to emit data into event queues.
type AuxLevel2Event struct {
	Bids []*AuxLevel2Event_Level `json:"bids"`
	Asks []*AuxLevel2Event_Level `json:"asks"`
}

type AuxOpenOrderEventInfo struct {
	Id                string     `json:"id"`
	CilentOrderId     string     `json:"client_order_id"`
	Price             JsonUint64 `json:"price"`
	Quantity          JsonUint64 `json:"quantity"`
	AuxAuToBurnPerLot JsonUint64 `json:"aux_au_to_burn_per_lot"`
	IsBid             bool       `json:"is_bid"`
	OwnerId           Address    `json:"order_id"`
	TimeoutTimestamp  JsonUint64 `json:"timeout_timestsamp"`
	OrderType         JsonUint64 `json:"order_type"`
	Timestamp         JsonUint64 `json:"timestamp"`
}

type AuxAllOrdersEvent struct {
	Bids [][]*AuxOpenOrderEventInfo `json:"bids"`
	Asks [][]*AuxOpenOrderEventInfo `json:"asks"`
}

//go:generate stringer -type AuxClobMarketSelfTradeType -linecomment

// AuxClobMarketSelfTradeType gives instruction on how to handle self trade
type AuxClobMarketSelfTradeType uint64

const (
	AuxClobMarketSelfTradeType_CancelPassive    AuxClobMarketSelfTradeType = iota + 200 // CANCEL_PASSIVE
	AuxClobMarketSelfTradeType_CancelAggressive                                         // CANCEL_AGGRESSIVE
	AuxClobMarketSelfTradeType_CancelBoth                                               // CANCEL_BOTH
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
