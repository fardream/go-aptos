package aptos

import (
	"crypto/ed25519"
	"encoding/json"

	"github.com/fardream/go-bcs/bcs"
	"golang.org/x/crypto/sha3"
)

// Transaction doesn't have signatures attached to them.
type Transaction struct {
	// Sender.
	// Note if the signer is the first parameter, it doesn't need to be included in the payload parameter list.
	Sender Address `json:"sender"`
	// SequenceNumber of the transaction for the sender.
	// transactions with sequence number less than the curernt on chain sequence number for address will be rejected.
	SequenceNumber JsonUint64 `json:"sequence_number"`
	// Payload
	Payload *TransactionPayload `json:"payload"`
	// MaxGasAmount
	MaxGasAmount JsonUint64 `json:"max_gas_amount"`
	// UnitGasPrice
	GasUnitPrice JsonUint64 `json:"gas_unit_price"`
	// ExpirationTimestampSecs
	ExpirationTimestampSecs JsonUint64 `json:"expiration_timestamp_secs"`

	// chain id - this is not serialized into json for payload
	ChainId uint8 `json:"-"`
}

// rawTransactionPrefix is sha3-256 of "APTOS::RawTransaction"
var rawTransactionPrefix []byte = sha3_Sum256Slice([]byte("APTOS::RawTransaction"))

// ToBCS get the signing bytes of the transaction.
// This is calling [EncodeTransaction] under the hood.
func (tx *Transaction) ToBCS() []byte {
	return EncodeTransaction(tx)
}

// sha3_Sum256Slice returns a the digest as a slice instead of [32]byte.
// [sha3.Sum256] returns a byte array of length 32 ([32]byte), whereas a slice is expected
// in many of use cases.
func sha3_Sum256Slice(data []byte) []byte {
	hashed := sha3.Sum256(data)
	return hashed[:]
}

// EncodeTransaction for signing.
// See here: [doc on aptos.dev], also see the [implementation in typescript]
//
// The process is follows:
//
//   - generate sha3_256 of "APTOS::RawTransaction"
//
// Then bcs serialize in the following order:
//
//   - sender
//   - sequence_number
//   - payload
//   - max_gas_amount
//   - gas_unit_price
//   - expiration_timestamp_secs
//   - chain_id
//
// for entry function payload, see [EntryFunctionPayload.ToBCS].
//
// [doc on aptos.dev]: https://aptos.dev/guides/creating-a-signed-transaction#signing-message
// [implementation in typescript]: https://github.com/aptos-labs/aptos-core/blob/ef6d3f45dfaeafcd76aa189b855d37a408a8e85e/ecosystem/typescript/sdk/src/transaction_builder/builder.ts#L69-L89
func EncodeTransaction(tx *Transaction) []byte {
	txBytes, err := bcs.Marshal(tx)
	if err != nil {
		panic(err)
	}

	return append(rawTransactionPrefix, txBytes...)
}

// GetHash get the hash of the transaction that can be used to look up the transaction on chain.
//
// Hash of the transaction is sha3-256 of ("RawTransaction" | bcs encoded transaction).
// BCS encoded transaction can be obtained by [Transaction.ToBCS] method.
//
// See [here].
//
// [here]: https://fullnode.mainnet.aptoslabs.com/v1/spec#/operations/get_transaction_by_hash
func (tx *Transaction) GetHash() []byte {
	signingBytes := EncodeTransaction(tx)
	// prefixBytes := []byte("RawTransaction")
	// hashed := sha3.Sum256(append(prefixBytes, signingBytes...))
	hashed := sha3.Sum256(signingBytes)
	return hashed[:]
}

// SingleSignature
type SingleSignature struct {
	Type      string `json:"type"`
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
}

// Ed25519SinatureType is the signature type for single signer based on a public/private key of ed25519 type.
const Ed25519SignatureType = "ed25519_signature"

// 64 zero bytes
const simulationSignature = "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

// NewSingleSignature creates a new signature
func NewSingleSignature(publicKey *ed25519.PublicKey, signature []byte) *SingleSignature {
	return &SingleSignature{
		Signature: prefixedHexString(signature),
		PublicKey: prefixedHexString(*publicKey),
		Type:      Ed25519SignatureType,
	}
}

// NewSingleSignatureForSimulation creates a new signature
func NewSingleSignatureForSimulation(publicKey *ed25519.PublicKey) *SingleSignature {
	return &SingleSignature{
		Type:      Ed25519SignatureType,
		PublicKey: prefixedHexString(*publicKey),
		Signature: simulationSignature,
	}
}

// TransactionInfo contains the information about the transaction that has been submitted to the blockchain.
type TransactionInfo struct {
	// Hash of the transaction.
	Hash string `json:"hash"`

	StateChangeHash     string `json:"state_change_hash"`
	EventRootHash       string `json:"event_root_hash"`
	StateCheckPointHash string `json:"state_checkpoint_hash"`

	GasUsed JsonUint64 `json:"gas_used"`

	// If the transaction is successful or not.
	Success bool `json:"success"`

	// VmStatus is useful for debug if the transaction failed.
	VmStatus            string `json:"vm_status"`
	AccumulatorRootHash string `json:"accumulator_root_hash"`

	Changes []json.RawMessage `json:"changes"`
	Events  []*RawEvent       `json:"events"`

	Timestamp JsonUint64 `json:"timestamp"`
}

// TransactionWithInfo is contains the transaction itself and the results of the transaciton execution.
type TransactionWithInfo struct {
	*Transaction     `json:",inline"`
	Type             string `json:"type"`
	*TransactionInfo `json:",inline"`
}
