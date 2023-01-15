package aptos

import (
	"context"
	"encoding/json"
	"net/http"
)

// [View] is a move function on chain that is marked with "[view]". Those functions can be called from off-chain and the return value of the function will be returned off-chain (hence the name view).
//
// [View]: https://fullnode.devnet.aptoslabs.com/v1/spec#/operations/view
func (client *Client) View(ctx context.Context, request *ViewRequest) (*AptosResponse[ViewResponse], error) {
	return doRequestForType[ViewResponse](ctx, client, request)
}

// ViewRequest is similar to entry function payload
type ViewRequest struct {
	Function      *MoveFunctionTag    `json:"function" url:"-"`
	TypeArguments []*MoveTypeTag      `json:"type_arguments" url:"-"`
	Arguments     []*EntryFunctionArg `json:"arguments" url:"-"`

	LedgerVersion *uint64 `json:"-" bcs:"-" url:"ledger_version,omitempty"`
}

var _ AptosRequest = (*ViewRequest)(nil)

func (r *ViewRequest) Body() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func (r *ViewRequest) HttpMethod() string {
	return http.MethodPost
}

func (r *ViewRequest) PathSegments() ([]string, error) {
	return []string{"view"}, nil
}

type ViewResponse = json.RawMessage
