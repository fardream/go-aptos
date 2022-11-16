package aptos

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fardream/go-bcs/bcs"
)

// EntryFunctionPayload
type EntryFunctionPayload struct {
	Function      *MoveFunctionTag      `json:"function"`
	TypeArguments []*MoveTypeTag        `json:"type_arguments"`
	Arguments     EntryFunctionArgSlice `json:"arguments"`
}

type EntryFunctionArg struct {
	Bool    *bool
	Uint8   *uint8
	Uint64  *uint64
	Uint128 *bcs.Uint128
	Address *Address
	Signer  *struct{}
	Vector  *[]byte
	Struct  *MoveTypeTag
}

var _ bcs.Enum = (*EntryFunctionArg)(nil)

func (m EntryFunctionArg) IsBcsEnum() {}

// NewEntryFunctionPayload
func NewEntryFunctionPayload(
	functionName *MoveFunctionTag,
	typeArguments []*MoveTypeTag,
	arguments []*EntryFunctionArg,
) *TransactionPayload {
	r := &EntryFunctionPayload{
		Function:      functionName,
		TypeArguments: typeArguments,
		Arguments:     arguments,
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
		Uint64: new(uint64),
	}

	*r.Uint64 = v

	return r
}

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

// EntryFunctionArgSlice
//
// Slices of [EntryFunctionArg] need special handling during serialization and deserialization.
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
type EntryFunctionArgSlice []*EntryFunctionArg

var _ json.Unmarshaler = (*EntryFunctionArgSlice)(nil)

func (s *EntryFunctionArgSlice) UnmarshalJSON(data []byte) error {
	var objects []json.RawMessage
	if err := json.Unmarshal(data, &objects); err != nil {
		return err
	}

	result := []*EntryFunctionArg{}

	for _, msg := range objects {
		var j JsonUint64
		if err := j.UnmarshalJSON(msg); err == nil {
			result = append(result, EntryFunctionArg_Uint64(uint64(j)))
			continue
		}
		var b bool
		if err := json.Unmarshal(msg, &b); err == nil {
			result = append(result, EntryFunctionArg_Bool(b))
			continue
		}
		var str string
		if err := json.Unmarshal(msg, &str); err == nil {
			if strings.HasPrefix(str, "0x") {
				addr, err := ParseAddress(str)
				if err == nil {
					result = append(result, EntryFunctionArg_Address(addr))
				}
				continue
			}

			result = append(result, EntryFunctionArg_String(str))
			continue
		}

		return fmt.Errorf("failed to unmarshal %s", string(msg))
	}

	*s = result

	return nil
}
