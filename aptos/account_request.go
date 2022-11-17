package aptos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// AccountModule contains the byte codes and the abi of the module.
type AccountModule struct {
	Bytecode MoveBytecode   `json:"bytecode"`
	Abi      *MoveModuleABI `json:"abi"`
}

// AccountResource includes the type and json encoding of the data.
type AccountResource struct {
	Type *MoveStructTag  `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// [GetAccount]
//
// [GetAccount]: https://fullnode.mainnet.aptoslabs.com/v1/spec#/operations/get_account
func (client *Client) GetAccount(ctx context.Context, request *GetAccountRequest) (*AptosResponse[GetAccountResponse], error) {
	return doRequestForType[GetAccountResponse](ctx, client, request)
}

type GetAccountRequest struct {
	GetRequest

	Address       Address `url:"-"`
	LedgerVersion *uint64 `url:"ledger_version,omitempty"`
}

var _ AptosRequest = (*GetAccountRequest)(nil)

func (r *GetAccountRequest) PathSegments() ([]string, error) {
	if r.Address.IsZero() {
		return nil, fmt.Errorf("empty address for account request")
	}

	return []string{"accounts", r.Address.String()}, nil
}

type GetAccountResponse struct {
	SequenceNumber    JsonUint64 `json:"sequence_number"`
	AuthenticationKey string     `json:"authentication_key"`
}

// [GetAccountResources]
//
// [GetAccountResources]: https://fullnode.mainnet.aptoslabs.com/v1/accounts/{address}/resources
func (client *Client) GetAccountResources(ctx context.Context, request *GetAccountResourcesRequest) (*AptosResponse[GetAccountResourcesResponse], error) {
	return doRequestForType[GetAccountResourcesResponse](ctx, client, request)
}

var _ AptosRequest = (*GetAccountResourcesRequest)(nil)

type GetAccountResourcesRequest struct {
	GetRequest

	Address       Address `url:"-"`
	LedgerVersion *uint64 `url:"ledger_version,omitempty"`
}

func (r *GetAccountResourcesRequest) PathSegments() ([]string, error) {
	if r.Address.IsZero() {
		return nil, fmt.Errorf("empty address for account resources request")
	}

	return []string{"accounts", r.Address.String(), "resources"}, nil
}

// TypedAccountResource
type TypedAccountResource[T any] struct {
	AccountResource
	ParsedData *T
}

type GetAccountResourcesResponse []AccountResource

// [GetAccountResource]
//
// [GetAccountResource]: https://fullnode.mainnet.aptoslabs.com/v1/spec#/operations/get_account_resource
func (client *Client) GetAccountResource(ctx context.Context, request *GetAccountResourceRequest) (*AptosResponse[GetAccountResourceResponse], error) {
	return doRequestForType[GetAccountResourceResponse](ctx, client, request)
}

type GetAccountResourceRequest struct {
	GetRequest

	Address       Address        `url:"-"`
	LedgerVersion *JsonUint64    `url:"ledger_version,omitempty"`
	Type          *MoveStructTag `url:"-"`
}

func (r *GetAccountResourceRequest) PathSegments() ([]string, error) {
	if r.Address.IsZero() {
		return nil, fmt.Errorf("empty address for account resource request")
	}

	if r.Type == nil {
		return nil, fmt.Errorf("missing type information")
	}

	return []string{"accounts", r.Address.String(), "resource", url.PathEscape(r.Type.String())}, nil
}

type GetAccountResourceResponse struct {
	*AccountResource `json:",inline"`
}

// GetAccountResourceWithType get the resource of specified move type, then marshal it into requested type T.
//
// This is equivalent of calling [Client.GetAccountResource], then marshal the response into the type.
//
// This is a function since golang doesn't support generic method.
func GetAccountResourceWithType[T any](ctx context.Context, client *Client, address Address, moveType *MoveStructTag, ledgerVersion uint64) (*T, error) {
	request := &GetAccountResourceRequest{
		Address: address,
		Type:    moveType,
	}
	if ledgerVersion > 0 {
		request.LedgerVersion = new(JsonUint64)
		*(request.LedgerVersion) = JsonUint64(ledgerVersion)
	}

	resp, err := client.GetAccountResource(ctx, request)
	if err != nil {
		return nil, err
	}

	result := new(T)

	if err := json.Unmarshal(resp.Parsed.Data, result); err != nil {
		return nil, err
	}

	return result, nil
}

// [GetAccountModules]
//
// [GetAccountModules]: https://fullnode.mainnet.aptoslabs.com/v1/spec#/operations/get_account_modules
func (client *Client) GetAccountModules(ctx context.Context, request *GetAccountModulesRequest) (*AptosResponse[GetAccountModulesResponse], error) {
	return doRequestForType[GetAccountModulesResponse](ctx, client, request)
}

var _ AptosRequest = (*GetAccountModulesRequest)(nil)

type GetAccountModulesRequest struct {
	GetRequest

	Address       Address `url:"-"`
	LedgerVersion *uint64 `url:"ledger_version,omitempty"`
}

func (r *GetAccountModulesRequest) PathSegments() ([]string, error) {
	if r.Address.IsZero() {
		return nil, fmt.Errorf("empty address for account modules request")
	}

	return []string{"accounts", r.Address.String(), "modules"}, nil
}

type GetAccountModulesResponse = []AccountModule

// [GetAccountModule]
//
// [GetAccountModule]: https://fullnode.mainnet.aptoslabs.com/v1/spec#/operations/get_account_modules
func (client *Client) GetAccountModule(ctx context.Context, request *GetAccountModuleRequest) (*AptosResponse[GetAccountModuleResponse], error) {
	return doRequestForType[GetAccountModuleResponse](ctx, client, request)
}

var _ AptosRequest = (*GetAccountModuleRequest)(nil)

type GetAccountModuleRequest struct {
	GetRequest

	Address       Address `url:"-"`
	LedgerVersion *uint64 `url:"ledger_version,omitempty"`
	ModuleName    string  `url:"-"`
}

func (r *GetAccountModuleRequest) PathSegments() ([]string, error) {
	if r.Address.IsZero() {
		return nil, fmt.Errorf("empty address for account modules request")
	}

	if !identifierRegex.MatchString(r.ModuleName) {
		return nil, fmt.Errorf("%s is not a valid module name", r.ModuleName)
	}

	return []string{"accounts", r.Address.String(), "module", r.ModuleName}, nil
}

type GetAccountModuleResponse = AccountModule
