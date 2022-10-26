package aptos

import "crypto/ed25519"

type AuxClientConfig struct {
	Address           Address
	Deployer          Address
	DataFeedAddress   Address
	DataFeedPublicKey ed25519.PublicKey
}
