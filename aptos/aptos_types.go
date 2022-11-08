package aptos

// [TableWithLength] is a wrapper around [Table], which keeps track of the length of the table.
//
// [TableWithLength]: https://github.com/aptos-labs/aptos-core/blob/ef6d3f45dfaeafcd76aa189b855d37a408a8e85e/aptos-move/framework/aptos-stdlib/sources/table_with_length.move
type TableWithLength struct {
	Inner  *Table     `json:"inner"`
	Length JsonUint64 `json:"length"`
}

// [Table] is a storage class provided by aptos framework where each individual element can be accessed independently.
// Table is unlike the move standard vector where the whole vector needs to be loaded even if only one element is needed
// in the vector. However, Table is actually not stored on chain directly, and needs a separate table api to query off-chain.
//
// [Table]: https://github.com/aptos-labs/aptos-core/blob/ef6d3f45dfaeafcd76aa189b855d37a408a8e85e/aptos-move/framework/aptos-stdlib/sources/table.move
type Table struct {
	Handle Address `json:"handle"`
}

// GUID_ID is the [ID] type of aptos framework.
//
// [ID]: https://github.com/aptos-labs/aptos-core/blob/ef6d3f45dfaeafcd76aa189b855d37a408a8e85e/aptos-move/framework/aptos-framework/sources/guid.move#L10-L16
type GUID_ID struct {
	CreationNumber JsonUint64 `json:"creation_number"`
	AccountAddress Address    `json:"account_address"`
}

// [GUID] contains an [ID]. This is an onchain struct.
//
// [GUID]: https://github.com/aptos-labs/aptos-core/blob/ef6d3f45dfaeafcd76aa189b855d37a408a8e85e/aptos-move/framework/aptos-framework/sources/guid.move#L5-L8
type GUID struct {
	Id GUID_ID `json:"id"`
}

// [EventHandler] contains the information for an event handler.
//
// [EventHandler]: https://github.com/aptos-labs/aptos-core/blob/ef6d3f45dfaeafcd76aa189b855d37a408a8e85e/aptos-move/framework/aptos-framework/sources/event.move
type EventHandler struct {
	Counter JsonUint64 `json:"counter"`
	GUID    struct {
		Id struct {
			CreationNumber JsonUint64 `json:"creation_num"`
			AccountAddress Address    `json:"addr"`
		} `json:"id"`
	} `json:"guid"`
}

// AptosStdAddress is the aptos standard library and aptos framework's address on chain, which is 0x1.
var AptosStdAddress = MustParseAddress("0x1")

// LedgerInfo contains basic information about the chain.
type LedgerInfo struct {
	ChainId             uint8      `json:"chain_id"`
	Epoch               JsonUint64 `json:"epoch"`
	LedgerVersion       JsonUint64 `json:"ledger_version"`
	OldestLedgerVersion JsonUint64 `json:"oldest_ledger_version"`
	LedgerTimestamp     JsonUint64 `json:"ledger_timestamp"`
	NodeRole            string     `json:"node_role"`
	OldestBlockHeight   JsonUint64 `json:"oldest_block_height"`
	BlockHeight         JsonUint64 `json:"block_height"`
	GitHash             string     `json:"git_hash"`
}
