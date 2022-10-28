package aptos

// Signer is an interface to sign a transaction
type Signer interface {
	// Sign transaction
	Sign(tx *Transaction) (*SingleSignature, error)
	// Sign transaction for simulation
	SignForSimulation(tx *Transaction) (*SingleSignature, error)
}

// RawDataSigner signs arbitrary bytes. Prefer to use Signer
type RawDataSigner interface {
	SignRawData([]byte) (*SingleSignature, error)
	SignRawDataForSimulation([]byte) (*SingleSignature, error)
}
