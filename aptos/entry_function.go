package aptos

import (
	"encoding/json"
	"fmt"
	"strings"
)

type EntryFunctionPayload struct {
	Type          string                `json:"type"`
	Function      string                `json:"function"`
	TypeArguments []*MoveTypeTag        `json:"type_arguments"`
	Arguments     EntryFunctionArgSlice `json:"arguments"`
}

type EntryFunctionArg interface {
	ToBCS() []byte
}

func NewEntryFunctionPayload(functionName string, typeArguments []*MoveTypeTag, arguments []EntryFunctionArg) *EntryFunctionPayload {
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

type EntryFunctionArg_Uint64 = JsonUint64

type EntryFunctionArg_String string

var _ EntryFunctionArg = (*EntryFunctionArg_String)(nil)

func (s EntryFunctionArg_String) ToBCS() []byte {
	stringBytes := []byte(string(s))

	prefix := ULEB128Encode(len(s))
	return append(prefix, stringBytes...)
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
