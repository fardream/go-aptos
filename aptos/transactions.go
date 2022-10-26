package aptos

import (
	"crypto/ed25519"
	"encoding/json"

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
	// MaxGasAmount
	MaxGasAmount JsonUint64 `json:"max_gas_amount"`
	// UnitGasPrice
	GasUnitPrice JsonUint64 `json:"gas_unit_price"`
	// ExpirationTimestampSecs
	ExpirationTimestampSecs JsonUint64            `json:"expiration_timestamp_secs"`
	Payload                 *EntryFunctionPayload `json:"payload"`
}

var rawTransactionPrefix []byte

func init() {
	hash := sha3.Sum256([]byte("APTOS::RawTransaction"))
	rawTransactionPrefix = hash[:]
}

// EncodeTransaction for signing
// See here: https://aptos.dev/guides/creating-a-signed-transaction#signing-message
// also see here
//
// The process is follows
// - generate sha3_256 of "APTOS::RawTransaction"
// - bcs serialize in the following order:
//   - sender
//   - sequence_number
//   - payload
//   - max_gas_amount
//   - gas_unit_price
//   - expiration_timestamp_secs
//   - chain_id
//
// for entry function payload, the serialization is as follows
func EncodeTransaction(tx *Transaction, chainId uint8) []byte {
	var encoded []byte
	encoded = append(encoded, rawTransactionPrefix...)
	encoded = append(encoded, tx.Sender.ToBCS()...)
	encoded = append(encoded, tx.SequenceNumber.ToBCS()...)
	encoded = append(encoded, tx.Payload.ToBCS()...)
	encoded = append(encoded, tx.MaxGasAmount.ToBCS()...)
	encoded = append(encoded, tx.GasUnitPrice.ToBCS()...)
	encoded = append(encoded, tx.ExpirationTimestampSecs.ToBCS()...)
	encoded = append(encoded, chainId)

	return encoded
}

// SingleSignature
type SingleSignature struct {
	Type      string `json:"type"`
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
}

const (
	Ed25519SignatureType = "ed25519_signature"
	// 64 zero bytes
	simulationSignature = "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
)

func NewSingleSignature(publicKey *ed25519.PublicKey, signature []byte) *SingleSignature {
	return &SingleSignature{
		Signature: prefixedHexString(signature),
		PublicKey: prefixedHexString(*publicKey),
		Type:      Ed25519SignatureType,
	}
}

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
