package aptos

import "encoding/json"

// Event emitted from aptos transactions
type Event[T any] struct {
	// GUID is the identifier of the event handler
	GUID GUID_ID `json:"guid"`
	// SequenceNumber of the event.
	// This is monotonically increasing without any gaps.
	SequenceNumber JsonUint64 `json:"sequence_number"`
	// Type of the event
	Type MoveStructTag `json:"type"`
	// Data of the event
	Data *T `json:"data"`
}

// RawEvent stores the data as [json.RawMessage]/byte slice
type RawEvent = Event[json.RawMessage]
