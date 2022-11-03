package aptos

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"
)

func mustEd25519PublicKey(s string) ed25519.PublicKey {
	return must(hex.DecodeString(strings.TrimPrefix(s, "0x")))
}

// GetAuxClientConfig gets the aux exchange configuration for an aptos network.
// Only mainnet, devnet, and testnet are supported.
func GetAuxClientConfig(chain Network) (*AuxClientConfig, error) {
	switch chain {
	case Mainnet:
		return &AuxClientConfig{
			Address:           MustParseAddress("0xbd35135844473187163ca197ca93b2ab014370587bb0ed3befff9e902d6bb541"),
			Deployer:          MustParseAddress("0x5a5e124ea1f3fc5fcfae3c198765c3b4c8d72c7236ae97ef6e5a9bc7cfda549c"),
			DataFeedAddress:   MustParseAddress("0x73daac91bd205cec351524974cfae156985f947e07d55f2acfcb38981fdb8898"),
			DataFeedPublicKey: mustEd25519PublicKey("0xa257c3a9f8c0316326681fc525c038886e39b3495c99bb28e1bca01ff6216634"),
		}, nil
	case Testnet:
		return &AuxClientConfig{
			Address:           MustParseAddress("0x8b7311d78d47e37d09435b8dc37c14afd977c5cfa74f974d45f0258d986eef53"),
			Deployer:          MustParseAddress("0x27a5ed998335d3b74ee2329bdb803f25095ca1137015a115e748b366c44f73be"),
			DataFeedAddress:   MustParseAddress("0x490d9592c7f246ecd5eef80e0e5592fef813d0adb43b26dbedc0d045282c36b8"),
			DataFeedPublicKey: mustEd25519PublicKey("0x5252282e6fd74873a1a777e707496919cb118fb65ba46e5271ebd4c2af716a28"),
		}, nil
	case Devnet:
		return &AuxClientConfig{
			Address:           MustParseAddress("0xea383dc2819210e6e427e66b2b6aa064435bf672dc4bdc55018049f0c361d01a"),
			Deployer:          MustParseAddress("0x52746eee4d2ecc79f940f617d1e98f885467c185e93a444bc6231a8b1960c48a"),
			DataFeedAddress:   MustParseAddress("0x84f372536c73df84327d2af63992f4443e2bd1aec8695fa85693e256fc1f904f"),
			DataFeedPublicKey: mustEd25519PublicKey("0x2a27ecf198ff20db2634c43177e0d492df63105fa7106706b91a22dc42797d88"),
		}, nil
	default:
		return nil, fmt.Errorf("unknown network: %s. if you are looking to use localnet, consider GetAuxInfoFromLocalAccount", chain)
	}
}

// MustGetAuxClientConfig returns the aux exchange configuration for an aptos network.
// Panic if fails.
func MustGetAuxClientConfig(chain Network) *AuxClientConfig {
	return must(GetAuxClientConfig(chain))
}

// GetAuxClientConfigFromLocalAccount returns the aux configuration based on a local account.
// The input local account must be the account that deploys aux to the chain.
func GetAuxClientConfigFromLocalAccount(localAccount *LocalAccount) *AuxClientConfig {
	return &AuxClientConfig{
		Address:           CalculateResourceAddress(localAccount.Address, []byte("aux")),
		Deployer:          localAccount.Address,
		DataFeedAddress:   localAccount.Address,
		DataFeedPublicKey: localAccount.PublicKey,
	}
}
