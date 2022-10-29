package aptos

import (
	"context"
	"fmt"
)

// [GetEventsByCreationNumber]
// [GetEventsByCreationNumber]: https://fullnode.mainnet.aptoslabs.com/v1/spec#/operations/get_events_by_creation_number
func (client *Client) GetEventsByCreationNumber(ctx context.Context, request *GetEventsByCreationNumberRequest) (*AptosResponse[GetEventsByCreationNumberResponse], error) {
	return doRequestForType[GetEventsByCreationNumberResponse](ctx, client, request)
}

type GetEventsByCreationNumberRequest struct {
	GetRequest

	Limit *JsonUint64 `url:"limit,omitempty"`
	Start *JsonUint64 `url:"start,omitempty"`

	Address        Address    `url:"-"`
	CreationNumber JsonUint64 `url:"-"`
}

func (request *GetEventsByCreationNumberRequest) PathSegments() ([]string, error) {
	return []string{"accounts", request.Address.String(), "events", fmt.Sprintf("%d", request.CreationNumber)}, nil
}

type GetEventsByCreationNumberResponse []*RawEvent

// [GetEventsByEventHandler]
// [GetEventsByEventHandler]: https://fullnode.mainnet.aptoslabs.com/v1/spec#/operations/get_events_by_event_handle
func (client *Client) GetEventsByEventHandler(ctx context.Context, request *GetEventsByEventHandlerRequest) (*AptosResponse[GetEventsByEventHandlerResponse], error) {
	return doRequestForType[GetEventsByEventHandlerResponse](ctx, client, request)
}

type GetEventsByEventHandlerRequest struct {
	GetRequest

	Limit *JsonUint64 `url:"limit,omitempty"`
	Start *JsonUint64 `url:"start,omitempty"`

	Address      Address      `url:"-"`
	EventHandler *MoveTypeTag `url:"-"`
	FieldName    string       `url:"-"`
}

func (request *GetEventsByEventHandlerRequest) PathSegments() ([]string, error) {
	if !identifierRegex.MatchString(request.FieldName) {
		return nil, fmt.Errorf("%s is not proper field name", request.FieldName)
	}
	return []string{"accounts", request.Address.String(), "events", request.EventHandler.String(), request.FieldName}, nil
}

type GetEventsByEventHandlerResponse []*RawEvent

// LoadEvents loads the events as defined by the creation number and address.
// Load sliceSize events at one request.
func (client *Client) LoadEvents(ctx context.Context, address Address, creationNumber uint64, start, end, sliceSize uint64) ([]*RawEvent, error) {
	result := make([]*RawEvent, 0, end-start)
	for ; start < end; start += sliceSize {
		if start+sliceSize > end {
			sliceSize = end - start
		}
		startJsonUint64 := JsonUint64(start)
		limit := JsonUint64(sliceSize)
		resp, err := client.GetEventsByCreationNumber(ctx, &GetEventsByCreationNumberRequest{
			Address:        address,
			CreationNumber: JsonUint64(creationNumber),
			Start:          &startJsonUint64,
			Limit:          &limit,
		})
		if err != nil {
			return result, err
		}

		result = append(result, (*resp.Parsed)...)
	}

	return result, nil
}
