package aptos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

// Client for aptos
type Client struct {
	restUrl     string
	gasEstimate *EstimateGasPriceResponse
	ledgerInfo  *LedgerInfo
	chainId     uint8

	defaultTransactionOptions TransactionOptions
}

func GetChainIdForNetwork(network Network) uint8 {
	switch network {
	case Mainnet:
		return 1
	case Testnet:
		return 2
	default:
		return 0
	}
}

// NewClient creates a new client for the given network.
// Values will be taken from the default of the network.
// URL can be left empty.
// Client's default option includes expire after 5 minutes and max gas of 20,000.
func NewClient(network Network, restUrl string, transactionOptions ...TransactionOption) (*Client, error) {
	url := restUrl
	var err error
	if url == "" {
		url, _, err = GetDefaultEndpoint(network)
		if err != nil {
			return nil, err
		}
	}

	client := &Client{
		restUrl: url,
	}

	for _, opt := range transactionOptions {
		client.defaultTransactionOptions.SetOption(opt)
	}

	if len(transactionOptions) == 0 {
		client.defaultTransactionOptions.SetOption(TransactionOption_ExpireAfter(5 * time.Minute))
		client.defaultTransactionOptions.SetOption(TransactionOption_MaxGasAmount(20000))
	}

	client.SetChainId(GetChainIdForNetwork(network))

	return client, nil
}

// MustNewClient creates a new client, panic if error happens.
func MustNewClient(network Network, restUrl string, transactionOptions ...TransactionOption) *Client {
	return must(NewClient(network, restUrl, transactionOptions...))
}

// SetChainId after client is created.
func (client *Client) SetChainId(chainId uint8) {
	client.chainId = chainId
}

// RefreshData updates gas price estimates and ledger info.
func (client *Client) RefreshData(ctx context.Context) error {
	if est, err := client.EstimateGasPrice(ctx); err != nil {
		return err
	} else {
		client.gasEstimate = est.Parsed
	}

	if info, err := client.GetLedgerInfo(ctx); err != nil {
		return err
	} else {
		client.ledgerInfo = info.Parsed.LedgerInfo
		client.chainId = info.Parsed.ChainId
	}

	return nil
}

// AptosRestError contains the http status code, message body and message of the response.
// This is returned when status code >= 400 is returned.
type AptosRestError struct {
	// HttpStatusCode returned
	HttpStatusCode int

	// Body of the response
	Body []byte

	// Message
	Message string
}

var _ error = (*AptosRestError)(nil)

func (e *AptosRestError) Error() string {
	return fmt.Sprintf("http failed: %d %s %s", e.HttpStatusCode, e.Message, e.Body)
}

func doRequest[TResponse any](ctx context.Context, client *Client, method string, pathSegments []string, queryString string, body []byte) (*AptosResponse[TResponse], error) {
	fullUrl, err := url.JoinPath(client.restUrl, pathSegments...)
	if err != nil {
		return nil, err
	}

	if queryString != "" {
		fullUrl = fmt.Sprintf("%s?%s", fullUrl, queryString)
	}

	r, err := http.NewRequestWithContext(ctx, method, fullUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	r.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, &AptosRestError{HttpStatusCode: resp.StatusCode, Message: resp.Status, Body: msg}
	}

	res := &AptosResponse[TResponse]{
		RawData: msg,
		Parsed:  new(TResponse),
		Headers: &AptosReponseHeader{},
	}

	// headers
	res.Headers.AptosBlockHeight = resp.Header.Get("X-APTOS-BLOCK-HEIGHT")
	res.Headers.AptosChainId = resp.Header.Get("X-APTOS-CHAIN-ID")
	res.Headers.AptosEpoch = resp.Header.Get("X-APTOS-EPOCH")
	res.Headers.AptosLedgerOldestVersion = resp.Header.Get("X-APTOS-LEDGER-OLDEST-VERSION")
	res.Headers.AptosLedgerTimestampUsec = resp.Header.Get("X-APTOS-LEDGER-TIMESTAMPUSEC")
	res.Headers.AptosLedgerVersion = resp.Header.Get("X-APTOS-LEDGER-VERSION")
	res.Headers.AptosOldestBlockHeight = resp.Header.Get("X-APTOS-OLDEST-BLOCK-HEIGHT")

	err = json.Unmarshal(msg, res.Parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w\n, response msg: %s", err, string(msg))
	}

	return res, nil
}

// doRequestForType takes an AptosRequest, construct the request and pass on to [doRequest].
func doRequestForType[TResponse any](ctx context.Context, client *Client, request AptosRequest) (*AptosResponse[TResponse], error) {
	pathSegments, err := request.PathSegments()
	if err != nil {
		return nil, fmt.Errorf("failed to construct path: %w", err)
	}

	queryV, err := query.Values(request)
	if err != nil {
		return nil, err
	}
	queryString := queryV.Encode()

	body, err := request.Body()
	if err != nil {
		return nil, err
	}

	return doRequest[TResponse](ctx, client, request.HttpMethod(), pathSegments, queryString, body)
}
