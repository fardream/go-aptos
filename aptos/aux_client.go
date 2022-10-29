package aptos

// AuxClient combines [AuxClientConfig], [Client], and [Signer] for aptos for convenient access
type AuxClient struct {
	config *AuxClientConfig

	client *Client

	signer Signer

	userAddress Address
}

func NewAuxClient(client *Client, config *AuxClientConfig, signer Signer) *AuxClient {
	r := &AuxClient{
		client: client,
		config: config,
		signer: signer,
	}

	if signer != nil {
		r.userAddress = signer.SignerAddress()
	}

	return r
}
