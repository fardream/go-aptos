package aptos

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MoveFunctionTag identifies a function on chain. Note the function's generic parameters or parameters are not part of the tag.
type MoveFunctionTag struct {
	MoveModuleTag

	Name string
}

func NewMoveFunctionTag(address Address, module, name string) (*MoveFunctionTag, error) {
	moduleTag, err := NewMoveModuleTag(address, module)
	if err != nil {
		return nil, err
	}
	if !identifierRegex.MatchString(name) {
		return nil, fmt.Errorf("%s is not type name", name)
	}

	return &MoveFunctionTag{MoveModuleTag: *moduleTag, Name: name}, nil
}

func MustNewMoveFunctionTag(address Address, module, name string) *MoveFunctionTag {
	return must(NewMoveFunctionTag(address, module, name))
}

var _ json.Marshaler = (*MoveFunctionTag)(nil)

func (f MoveFunctionTag) String() string {
	return fmt.Sprintf("%s::%s::%s", f.Address.String(), f.Module, f.Name)
}

func (f MoveFunctionTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

var _ json.Unmarshaler = (*MoveFunctionTag)(nil)

func internalParseModuleFunctionTag(str string, f *MoveFunctionTag) error {
	segments := strings.Split(str, "::")
	if len(segments) != 3 {
		return fmt.Errorf("%s is not in the format of address::module", str)
	}

	addressStr := segments[0]

	address, err := ParseAddress(addressStr)
	if err != nil {
		return fmt.Errorf("%s doesn't contain a valid address: %w", str, err)
	}

	moduleNameStr := segments[1]
	if !identifierRegex.MatchString(moduleNameStr) {
		return fmt.Errorf("module name %s in %s is invalid", moduleNameStr, str)
	}

	nameStr := segments[2]
	if !identifierRegex.MatchString(nameStr) {
		return fmt.Errorf("function name %s in %s is invalid", nameStr, str)
	}
	f.Address = address
	f.Module = moduleNameStr
	f.Name = nameStr

	return nil
}

func ParseModuleFunctionTag(str string) (*MoveFunctionTag, error) {
	r := &MoveFunctionTag{}

	if err := internalParseModuleFunctionTag(str, r); err != nil {
		return nil, err
	}

	return r, nil
}

func (f *MoveFunctionTag) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	return internalParseModuleFunctionTag(str, f)
}

func (f MoveFunctionTag) ToBCS() []byte {
	return append(f.MoveModuleTag.ToBCS(), StringToBCS(f.Name)...)
}
