package aptos

import (
	"context"
)

// GetLedgerInfo
// https://fullnode.mainnet.aptoslabs.com/v1/spec#/operations/get_ledger_info
func (client *Client) GetLedgerInfo(ctx context.Context) (*AptosResponse[GetLedgerInfoResponse], error) {
	return doRequestForType[GetLedgerInfoResponse](ctx, client, newPathSegmentHolder())
}

type GetLedgerInfoResponse struct {
	*LedgerInfo `json:",inline"`
}

// Get ChainId
func (client *Client) GetChainId(ctx context.Context) (uint8, error) {
	if client.ledgerInfo == nil {
		if err := client.RefreshData(ctx); err != nil {
			return 0, err
		}
	}

	return uint8(client.ledgerInfo.ChainId), nil
}
