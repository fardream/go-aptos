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

const AuxClobMarketModuleName = "clob_market"

func (info *AuxClientConfig) ClobMarket_PlaceOrder(
	sender Address,
	isBid bool,
	baseCoin,
	quoteCoin *MoveTypeTag,
	limitPrice uint64,
	quantity uint64,
	auxToBurnPerLot uint64,
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
		Uint128{},                                  // client_order_id: u128,
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

func (info *AuxClientConfig) ClobMarket_CancelAll(sender Address, baseCoin, quoteCoin *MoveTypeTag, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxClobMarketModuleName, "cancel_all")
	payload := NewEntryFunctionPayload(
		function,
		[]*MoveTypeTag{baseCoin, quoteCoin},
		nil,
	)

	tx := &Transaction{Payload: payload}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}
