package aptos

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/fardream/go-bcs/bcs"
)

// MoveTypeTag is an Enum indicating the type of a move value.
// It is a [bcs.Enum]
// aptos right now supports the following types
// - bool
// - u8
// - u64
// - u128
// - address
// - signer
// - vector
// - struct
type MoveTypeTag struct {
	Bool    *struct{}      // 0
	Uint8   *struct{}      // 1
	Uint64  *struct{}      // 2
	Uint128 *struct{}      // 3
	Address *struct{}      // 4
	Signer  *struct{}      // 5
	Vector  *MoveTypeTag   // 6
	Struct  *MoveStructTag // 7
}

var (
	_ bcs.Enum         = (*MoveTypeTag)(nil)
	_ json.Marshaler   = (*MoveTypeTag)(nil)
	_ json.Unmarshaler = (*MoveTypeTag)(nil)
)

func (m MoveTypeTag) IsBcsEnum() {}

func (m MoveTypeTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m MoveTypeTag) String() string {
	allBytes, err := m.MarshalText()
	if err != nil {
		return err.Error()
	} else {
		return string(allBytes)
	}
}

func (m MoveTypeTag) MarshalText() ([]byte, error) {
	switch {
	case m.Bool != nil:
		return []byte("bool"), nil
	case m.Uint8 != nil:
		return []byte("u8"), nil
	case m.Uint64 != nil:
		return []byte("u64"), nil
	case m.Uint128 != nil:
		return []byte("u128"), nil
	case m.Address != nil:
		return []byte("address"), nil
	case m.Signer != nil:
		return []byte("signer"), nil
	case m.Vector != nil:
		return []byte(fmt.Sprintf("vector<%s>", m.Vector)), nil
	case m.Struct != nil:
		return []byte(m.Struct.String()), nil
	default:
		return nil, fmt.Errorf("-- unset move type tag --")
	}
}

func newEmptyStruct() *struct{} {
	return &struct{}{}
}

func (m *MoveTypeTag) reset() {
	m.Bool = nil
	m.Uint8 = nil
	m.Uint64 = nil
	m.Uint128 = nil
	m.Address = nil
	m.Signer = nil
	m.Vector = nil
	m.Struct = nil
}

func (m *MoveTypeTag) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	return m.unmarshalFromStr(str)
}

var vectorRegex = regexp.MustCompile("^vector<(.+)>$")

func (m *MoveTypeTag) unmarshalFromStr(str string) error {
	m.reset()
	switch str {
	case "bool":
		m.Bool = newEmptyStruct()
		return nil
	case "u8":
		m.Uint8 = newEmptyStruct()
		return nil
	case "u64":
		m.Uint64 = newEmptyStruct()
		return nil
	case "u128":
		m.Uint128 = newEmptyStruct()
		return nil
	case "address":
		m.Address = newEmptyStruct()
		return nil
	case "signer":
		m.Signer = newEmptyStruct()
		return nil
	}

	vectorData := vectorRegex.FindStringSubmatch(str)

	if len(vectorData) == 2 {
		m.Vector = new(MoveTypeTag)
		return m.Vector.unmarshalFromStr(vectorData[1])
	}

	structTag, err := ParseMoveStructTag(str)
	if err != nil {
		return fmt.Errorf("exhausted all options, and parsing as struct tag failed: %w", err)
	}

	m.Struct = structTag

	return nil
}

func (m *MoveTypeTag) UnmarshalTEXT(data []byte) error {
	return m.unmarshalFromStr(string(data))
}

func ParseMoveTypeTag(str string) (*MoveTypeTag, error) {
	r := &MoveTypeTag{}
	if err := r.unmarshalFromStr(str); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}
