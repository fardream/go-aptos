package aptos

import "net/http"

// AptosRequest defines a helper interface for quickly construct aptos request.
// The URL query string is generated from [query.Values], and all the fields that should not be included in the query string should be ignored with tag url:"-".
type AptosRequest interface {
	// PathSegments doesn't include the query string.
	PathSegments() ([]string, error)
	Body() ([]byte, error)
	HttpMethod() string
}

// GetRequest embed this struct for a get request where only path segments are necessary.
type GetRequest struct{}

func (*GetRequest) Body() ([]byte, error) {
	return nil, nil
}

func (*GetRequest) HttpMethod() string {
	return http.MethodGet
}

// pathSegmentHolder is a get request where all requests share the same path segments
type pathSegmentHolder struct {
	Segments []string `json:"-" url:"-"`

	GetRequest
}

func (p *pathSegmentHolder) PathSegments() ([]string, error) {
	return p.Segments, nil
}

func newPathSegmentHolder(segments ...string) *pathSegmentHolder {
	return &pathSegmentHolder{
		Segments: segments,
	}
}

// AptosReponseHeader contains the header information on a successful aptos response
type AptosReponseHeader struct {
	AptosBlockHeight         string
	AptosChainId             string
	AptosEpoch               string
	AptosLedgerOldestVersion string
	AptosLedgerTimestampUsec string
	AptosLedgerVersion       string
	AptosOldestBlockHeight   string
}

type AptosResponse[T any] struct {
	RawData []byte
	Parsed  *T
	Headers *AptosReponseHeader
}
