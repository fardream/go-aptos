package aptos

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fardream/go-bcs/bcs"
)

// EntryFunctionPayload
type EntryFunctionPayload struct {
	Function      *MoveFunctionTag    `json:"function"`
	TypeArguments []*MoveTypeTag      `json:"type_arguments"`
	Arguments     []*EntryFunctionArg `json:"arguments"`
}

// EntryFunctionArg is the argument to entry function.
// This is kind of an enum, but on the wire actually
// all bcs byte vectors.
// This means the values are first bcs serialized to vector of bytes.
// Then the vector of bytes are serialized.
type EntryFunctionArg struct {
	Bool    *bool
	Uint8   *uint8
	Uint64  *JsonUint64
	Uint128 *bcs.Uint128
	Address *Address
	Signer  *struct{}
	Vector  *[]byte
}

var (
	_ bcs.Marshaler    = (*EntryFunctionArg)(nil)
	_ json.Marshaler   = (*EntryFunctionArg)(nil)
	_ json.Unmarshaler = (*EntryFunctionArg)(nil)
)

func marshalWithByteLength[T any](v T) ([]byte, error) {
	if firstBytes, err := bcs.Marshal(v); err != nil {
		return nil, err
	} else {
		return bcs.Marshal(firstBytes)
	}
}

// MarshalBCS customizes bcs Marshal for [EntryFunctionArg].
// This method first bcs serialize the non-nil value, then
// serialize the resulted byte vector (simply prepending the length with ULEB128 encoding).
func (m EntryFunctionArg) MarshalBCS() ([]byte, error) {
	switch {
	case m.Bool != nil:
		return marshalWithByteLength(m.Bool)
	case m.Uint8 != nil:
		return marshalWithByteLength(m.Uint8)
	case m.Uint64 != nil:
		return marshalWithByteLength(m.Uint64)
	case m.Uint128 != nil:
		return marshalWithByteLength(m.Uint128)
	case m.Address != nil:
		return marshalWithByteLength(m.Address)
	case m.Vector != nil:
		return bcs.Marshal(m.Vector)
	default:
		return nil, fmt.Errorf("unset arg")
	}
}

// reset sets all values to nil.
func (m *EntryFunctionArg) reset() {
	m.Bool = nil
	m.Uint8 = nil
	m.Uint64 = nil
	m.Uint128 = nil
	m.Address = nil
	m.Signer = nil
	m.Vector = nil
}

func (m EntryFunctionArg) MarshalJSON() ([]byte, error) {
	switch {
	case m.Bool != nil:
		return json.Marshal(*m.Bool)
	case m.Address != nil:
		return json.Marshal(*m.Address)
	case m.Uint8 != nil:
		return json.Marshal(*m.Uint8)
	case m.Uint64 != nil:
		return json.Marshal(*m.Uint64)
	case m.Uint128 != nil:
		return json.Marshal(*m.Uint128)
	case m.Vector != nil:
		return json.Marshal(*m.Vector)
	}

	return nil, fmt.Errorf("unsupported: %v", m)
}

// NewEntryFunctionPayload
func NewEntryFunctionPayload(
	functionName *MoveFunctionTag,
	typeArguments []*MoveStructTag,
	arguments []*EntryFunctionArg,
) *TransactionPayload {
	r := &EntryFunctionPayload{
		Function: functionName,
		TypeArguments: mapSlices(typeArguments, func(t *MoveStructTag) *MoveTypeTag {
			return &MoveTypeTag{
				Struct: t,
			}
		}),
		Arguments: arguments,
	}

	if r.TypeArguments == nil {
		r.TypeArguments = make([]*MoveTypeTag, 0)
	}
	if r.Arguments == nil {
		r.Arguments = make([]*EntryFunctionArg, 0)
	}

	return &TransactionPayload{
		EntryFunctionPayload: r,
	}
}

// EntryFunctionArg_Uint8 represents u8 in move, equivalent to a byte
func EntryFunctionArg_Uint8(v uint8) *EntryFunctionArg {
	r := &EntryFunctionArg{
		Uint8: new(uint8),
	}

	*r.Uint8 = v

	return r
}

// EntryFunctionArg_Uint64 is equivalent to uint64, or u64 in move.
func EntryFunctionArg_Uint64(v uint64) *EntryFunctionArg {
	r := &EntryFunctionArg{
		Uint64: new(JsonUint64),
	}

	*r.Uint64 = JsonUint64(v)

	return r
}

// EntryFunctionArg_Uint128 creates a new [EntryFunctionArg] with Uint128 set from lo/hi [uint64].
func EntryFunctionArg_Uint128(lo uint64, hi uint64) *EntryFunctionArg {
	r := &EntryFunctionArg{
		Uint128: new(bcs.Uint128),
	}

	*r.Uint128 = *bcs.NewUint128FromUint64(lo, hi)

	return r
}

func EntryFunctionArg_String(v string) *EntryFunctionArg {
	r := &EntryFunctionArg{
		Vector: new([]byte),
	}

	*r.Vector = append([]byte{}, []byte(v)...)

	return r
}

// EntryFunctionArg_Bool
func EntryFunctionArg_Bool(v bool) *EntryFunctionArg {
	r := &EntryFunctionArg{
		Bool: new(bool),
	}

	*r.Bool = v

	return r
}

func EntryFunctionArg_Address(v Address) *EntryFunctionArg {
	r := &EntryFunctionArg{
		Address: new(Address),
	}

	*r.Address = v

	return r
}

// UnmarshalJSON for [EntryFunctionArg]. json doesn't have type information, so the process uses a series of heuristics.
//
//   - deserializing response from rest api either in json or bcs is difficult without knowing the types of the elements before
//     hand.
//
//     The following logic is used to deserialize the slices: first, the element of the slice will be first tested if it is an u64 or bool.
//     Then, it is checked to see if it is a string. If it is a string and it has 0x prefix, cast it to address. If casting to address is unsuccessful,
//     keep it as a string.
//
//   - during serialization, the element of entry function argument slice is prefixed with the length of the
//     serialized bytes. For example, instead of serialize true to 01, serialize it to 0101.
func (s *EntryFunctionArg) UnmarshalJSON(data []byte) error {
	var j JsonUint64
	if err := j.UnmarshalJSON(data); err == nil {
		s.reset()
		s.Uint64 = new(JsonUint64)
		*s.Uint64 = j
		return nil
	}
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		s.reset()
		s.Bool = new(bool)
		*s.Bool = b
		return nil
	}
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if strings.HasPrefix(str, "0x") {
			addr, err := ParseAddress(str)
			if err == nil {
				s.reset()
				s.Address = new(Address)
				*s.Address = addr
				return nil
			}
		}
		s.reset()
		s.Vector = new([]byte)
		*s.Vector = append(*s.Vector, []byte(str)...)
		return nil
	}

	return fmt.Errorf("failed to unmarshal %s", string(data))
}
