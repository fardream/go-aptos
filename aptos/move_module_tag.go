package aptos

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MoveModuleTag identifies a move module on the chain
type MoveModuleTag struct {
	// Address of module
	Address Address
	// Module name
	Module string
}

func (m MoveModuleTag) String() string {
	return fmt.Sprintf("%s::%s", m.Address.String(), m.Module)
}

var _ json.Marshaler = (*MoveModuleTag)(nil)

func (m MoveModuleTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

var _ json.Unmarshaler = (*MoveModuleTag)(nil)

func (m *MoveModuleTag) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	return parseMoveModuleTagInternal(str, m)
}

func parseMoveModuleTagInternal(module string, moduleTag *MoveModuleTag) error {
	segments := strings.Split(module, "::")
	if len(segments) != 2 {
		return fmt.Errorf("%s is not in the format of address::module", module)
	}

	addressStr := segments[0]

	address, err := ParseAddress(addressStr)
	if err != nil {
		return fmt.Errorf("%s doesn't contain a valid address: %w", module, err)
	}

	moduleNameStr := segments[1]
	if !identifierRegex.MatchString(moduleNameStr) {
		return fmt.Errorf("module name %s in %s is invalid", moduleNameStr, module)
	}

	moduleTag.Address = address
	moduleTag.Module = moduleNameStr

	return nil
}

func ParseMoveModuleTag(str string) (*MoveModuleTag, error) {
	r := &MoveModuleTag{}
	if err := parseMoveModuleTagInternal(str, r); err != nil {
		return nil, err
	}

	return r, nil
}

// NewMoveModuleTag creates a new module module tag
func NewMoveModuleTag(address Address, module string) (*MoveModuleTag, error) {
	if !identifierRegex.MatchString(module) {
		return nil, fmt.Errorf("%s is not valid module name", module)
	}

	return &MoveModuleTag{
		Address: address,
		Module:  module,
	}, nil
}

func MustNewMoveModuleTag(address Address, module string) *MoveModuleTag {
	return must(NewMoveModuleTag(address, module))
}

func (m *MoveModuleTag) Set(str string) error {
	return parseMoveModuleTagInternal(str, m)
}

func (m MoveModuleTag) Type() string {
	return "move-module"
}
