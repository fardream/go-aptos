package aptos

import (
	"encoding/json"
	"fmt"
	"strings"
)

// EntryFunctionPayload
type EntryFunctionPayload struct {
	Type          string                `json:"type"`
	Function      *MoveFunctionTag      `json:"function"`
	TypeArguments []*MoveTypeTag        `json:"type_arguments"`
	Arguments     EntryFunctionArgSlice `json:"arguments"`
}

// EntryFunctionArg is argument to entry function
type EntryFunctionArg interface {
	ToBCS() []byte
}

// NewEntryFunctionPayload
func NewEntryFunctionPayload(functionName *MoveFunctionTag, typeArguments []*MoveTypeTag, arguments []EntryFunctionArg) *EntryFunctionPayload {
	r := &EntryFunctionPayload{
		Type:          "entry_function_payload",
		Function:      functionName,
		TypeArguments: typeArguments,
		Arguments:     arguments,
	}

	if r.TypeArguments == nil {
		r.TypeArguments = make([]*MoveTypeTag, 0)
	}
	if r.Arguments == nil {
		r.Arguments = make([]EntryFunctionArg, 0)
	}

	return r
}

// EntryFunctionArg_Uint8 represents u8 in move, equivalent to a byte
type EntryFunctionArg_Uint8 uint8

var _ EntryFunctionArg = (*EntryFunctionArg_Uint8)(nil)

func (u *EntryFunctionArg_Uint8) UnmarshalJSON(data []byte) error {
	var v uint8
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*u = EntryFunctionArg_Uint8(v)

	return nil
}

func (u EntryFunctionArg_Uint8) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint8(u))
}

func (u EntryFunctionArg_Uint8) ToBCS() []byte {
	return []byte{byte(u)}
}

// EntryFunctionArg_Uint64 is equivalent to uint64, or u64 in move.
type EntryFunctionArg_Uint64 = JsonUint64

type EntryFunctionArg_String string

var _ EntryFunctionArg = (*EntryFunctionArg_String)(nil)

func (s EntryFunctionArg_String) ToBCS() []byte {
	return StringToBCS(string(s))
}

func (s EntryFunctionArg_String) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *EntryFunctionArg_String) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return nil
	}

	*s = EntryFunctionArg_String(str)

	return nil
}

// EntryFunctionArg_Bool
type EntryFunctionArg_Bool bool

var _ EntryFunctionArg = (*EntryFunctionArg_Bool)(nil)

func (b EntryFunctionArg_Bool) ToBCS() []byte {
	if bool(b) {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

func (b EntryFunctionArg_Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(bool(b))
}

func (b *EntryFunctionArg_Bool) UnmarshalJSON(data []byte) error {
	var bt bool
	if err := json.Unmarshal(data, &bt); err != nil {
		return err
	}

	*b = EntryFunctionArg_Bool(bt)

	return nil
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
type EntryFunctionArgSlice []EntryFunctionArg

var _ json.Unmarshaler = (*EntryFunctionArgSlice)(nil)

func (s *EntryFunctionArgSlice) UnmarshalJSON(data []byte) error {
	var objects []json.RawMessage
	if err := json.Unmarshal(data, &objects); err != nil {
		return err
	}

	result := []EntryFunctionArg{}

	for _, msg := range objects {
		var j JsonUint64
		if err := j.UnmarshalJSON(msg); err == nil {
			result = append(result, j)
			continue
		}
		var b EntryFunctionArg_Bool
		if err := b.UnmarshalJSON(msg); err == nil {
			result = append(result, b)
			continue
		}
		var str string
		if err := json.Unmarshal(msg, &str); err == nil {
			if strings.HasPrefix(str, "0x") {
				addr, err := ParseAddress(str)
				if err == nil {
					result = append(result, addr)
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

func (s EntryFunctionArgSlice) ToBCS() []byte {
	r := ULEB128Encode(len(s))
	for _, a := range s {
		ab := a.ToBCS()
		r = append(r, ULEB128Encode(len(ab))...)
		r = append(r, ab...)
	}

	return r
}

type EntryFunctionArgVector[T EntryFunctionArg] []T

func (v EntryFunctionArgVector[T]) ToBCS() []byte {
	r := ULEB128Encode(len(v))
	for _, x := range v {
		r = append(r, x.ToBCS()...)
	}

	return r
}

// ToBCS encodes EntryFunctionPayload to bytes.
//   - first byte is 2, indicating [EntryFunctionPayload]
//   - serialize function name
//   - serialize generic type arguments
//   - serialize the arguments. Note arguments are serialized first of number of arguments, then each argument needs the length of their serialized bytes as prefix.
func (f EntryFunctionPayload) ToBCS() []byte {
	r := []byte{2}
	r = append(r, f.Function.ToBCS()...)
	r = append(r, ULEB128Encode(len(f.TypeArguments))...)
	for _, t := range f.TypeArguments {
		r = append(r, t.ToBCS()...)
	}

	r = append(r, f.Arguments.ToBCS()...)

	return r
}
