package aptos

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/sha3"
)

// SignatureLength aptos uses ed25519 and signature is 64 bytes.
const SignatureLength = 64

// GenerateAuthenticationKey calculates the authentication key for a scheme.
//
// The information is based on documentation on [aptos.dev]
//
// Account in aptos is presented by SHA3-256 of
//   - a public key of ed25519 public key (pub_key|0x00)
//   - a series of ed25519 public keys, the number of signature required (pub_key_1 | pub_key_2 ... | pub_key_n | K | 0x01)
//   - an address and some seed (address | seed| 0xFF) if on chain.
//
// [aptos.dev]: https://aptos.dev/concepts/basics-accounts/#signature-schemes
func GenerateAuthenticationKey(
	totalSignerCount int,
	requiredSignerCount int,
	signerPublicKeys ...ed25519.PublicKey,
) (Address, error) {
	if len(signerPublicKeys) != totalSignerCount {
		return Address{}, fmt.Errorf(
			"require %d public keys, but only %d are present",
			totalSignerCount,
			len(signerPublicKeys),
		)
	}

	if requiredSignerCount > totalSignerCount {
		return Address{}, fmt.Errorf("required signature count is %d, less than total public key count: %d", requiredSignerCount, totalSignerCount)
	}

	var allBytes []byte
	for _, pubKey := range signerPublicKeys {
		allBytes = append(allBytes, pubKey...)
	}

	if totalSignerCount == 1 {
		allBytes = append(allBytes, 0)
	} else {
		allBytes = append(allBytes, byte(requiredSignerCount))
		allBytes = append(allBytes, 1)
	}

	return sha3.Sum256(allBytes), nil
}

// LocalAccount contains the private key, public key, and the address.
type LocalAccount struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	Address    Address
}

// NewLocalAccountFromPrivateKey creates a local account based on the private key.
// the authentication key will be the one calculated from public key.
func NewLocalAccountFromPrivateKey(privateKey *ed25519.PrivateKey) (*LocalAccount, error) {
	publicKey, ok := privateKey.Public().(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot get public key from private key")
	}
	accountAddress, err := GenerateAuthenticationKey(1, 1, publicKey)
	if err != nil {
		return nil, err
	}
	return &LocalAccount{
		PrivateKey: *privateKey,
		PublicKey:  publicKey,
		Address:    accountAddress,
	}, nil
}

func parseHexString(hexString string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(hexString, "0x"))
}

func prefixedHexString(input []byte) string {
	return "0x" + hex.EncodeToString(input)
}

// NewPrivateKeyFromHexString generates a private key from hex string.
func NewPrivateKeyFromHexString(hexString string) (*ed25519.PrivateKey, error) {
	pk, err := parseHexString(hexString)
	if err != nil {
		return nil, err
	}

	privateKey := ed25519.NewKeyFromSeed(pk)
	return &privateKey, nil
}

// IsOriginalAuthenticationKey checks if the authentication key
// is the authentication key generated from the single public key.
func (account *LocalAccount) IsOriginalAuthenticationKey() bool {
	authKey, _ := GenerateAuthenticationKey(1, 1, account.PublicKey)
	return authKey == account.Address
}

// NewLocalAccountWithRandomKey creates a new account with random key.
func NewLocalAccountWithRandomKey() (*LocalAccount, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	address, err := GenerateAuthenticationKey(1, 1, pub)
	if err != nil {
		return nil, err
	}

	return &LocalAccount{
		PublicKey:  pub,
		PrivateKey: priv,
		Address:    address,
	}, nil
}

var (
	_ Signer        = (*LocalAccount)(nil)
	_ RawDataSigner = (*LocalAccount)(nil)
)

func (account *LocalAccount) Sign(tx *Transaction) (*SingleSignature, error) {
	if !cmp.Equal(tx.Sender, account.Address) {
		return nil, fmt.Errorf("can only sign for self")
	}

	return account.SignRawData(EncodeTransaction(tx))
}

func (account *LocalAccount) SignForSimulation(tx *Transaction) (*SingleSignature, error) {
	if !cmp.Equal(tx.Sender, account.Address) {
		return nil, fmt.Errorf("can only sign for self")
	}

	return account.SignRawDataForSimulation(EncodeTransaction(tx))
}

func (account *LocalAccount) SignerAddress() Address {
	return account.Address
}

// SignRawData
func (account *LocalAccount) SignRawData(message []byte) (*SingleSignature, error) {
	signature := ed25519.Sign(account.PrivateKey, message)
	return NewSingleSignature(&account.PublicKey, signature), nil
}

// SignRawDataForSimulation
func (account *LocalAccount) SignRawDataForSimulation(message []byte) (*SingleSignature, error) {
	return NewSingleSignatureForSimulation(&account.PublicKey), nil
}
