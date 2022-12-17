package aptos

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var (
	identifierRegex  = regexp.MustCompile("^[A-z_][A-z0-9_]*$")
	whiteSpaceRegex  = regexp.MustCompile(`\s+`)
	genericTypeRegex = regexp.MustCompile(`^([A-z0-9_:]+)+<(.+)>$`)
)

// MoveStructTag represents the type of a move struct in the format of
// address::module_name::TypeName
// off-chain, the address is a 0x prefixed hex encoded string, but during move development there can be named addresses.
type MoveStructTag struct {
	MoveModuleTag
	// Name of the type
	Name string

	// GenericTypeParameters
	GenericTypeParameters []*MoveTypeTag
}

// Get the string presentation of the type, in the form of
// 0xaddresshex::module_name::TypeName
// or
// 0xaddresshex::module_name::TypeName<T1, T2>
func (t *MoveStructTag) String() string {
	genericListStr := ""
	if len(t.GenericTypeParameters) > 0 {
		genericListStr = fmt.Sprintf(
			"<%s>",
			strings.Join(
				mapSlices(
					t.GenericTypeParameters,
					func(t *MoveTypeTag) string {
						return t.String()
					},
				),
				",",
			),
		)
	}

	return fmt.Sprintf("%s::%s::%s%s", t.Address.String(), t.Module, t.Name, genericListStr)
}

// NewMoveStructTag
func NewMoveStructTag(address Address, module string, name string, genericTypeParameters []*MoveStructTag) (*MoveStructTag, error) {
	moduleTag, err := NewMoveModuleTag(address, module)
	if err != nil {
		return nil, err
	}
	if !identifierRegex.MatchString(name) {
		return nil, fmt.Errorf("%s is not type name", name)
	}

	return &MoveStructTag{
		MoveModuleTag: *moduleTag,
		Name:          name,

		GenericTypeParameters: mapSlices(
			genericTypeParameters,
			func(t *MoveStructTag) *MoveTypeTag {
				return &MoveTypeTag{
					Struct: t,
				}
			}),
	}, nil
}

func MustNewMoveStructTag(address Address, module, name string, genericTypeParameters []*MoveStructTag) *MoveStructTag {
	return must(NewMoveStructTag(address, module, name, genericTypeParameters))
}

func makeCanonicalSegment(input string) string {
	return whiteSpaceRegex.ReplaceAllString(input, "")
}

func parseMoveStructTagInternal(fullName string, moveTypeTag *MoveStructTag) error {
	name := makeCanonicalSegment(fullName)

	var segments []string

	genericMatches := genericTypeRegex.FindStringSubmatch(name)
	var genericParameters []*MoveTypeTag
	if len(genericMatches) == 3 {
		name = genericMatches[1]
		var err error
		genericParameters, err = parseGenericTypeListString(genericMatches[2])
		if err != nil {
			return err
		}
	}

	segments = strings.Split(name, "::")
	if len(segments) != 3 {
		return fmt.Errorf("%s is not in the format of address::module::Name", fullName)
	}

	addressStr := segments[0]

	address, err := ParseAddress(addressStr)
	if err != nil {
		return fmt.Errorf("%s doesn't contain a valid address: %w", fullName, err)
	}

	moduleNameStr := segments[1]
	if !identifierRegex.MatchString(moduleNameStr) {
		return fmt.Errorf("module name %s in %s is invalid", moduleNameStr, fullName)
	}
	typeNameStr := segments[2]

	moveTypeTag.Address = address
	moveTypeTag.Module = moduleNameStr
	moveTypeTag.Name = typeNameStr
	moveTypeTag.GenericTypeParameters = genericParameters

	return nil
}

// ParseMoveStructTag takes the full name of the move type tag
func ParseMoveStructTag(fullName string) (*MoveStructTag, error) {
	r := &MoveStructTag{}
	if err := parseMoveStructTagInternal(fullName, r); err != nil {
		return nil, err
	}

	return r, nil
}

func parseGenericTypeListString(genericTypeListString string) ([]*MoveTypeTag, error) {
	leftBracketCount := 0
	var parsedTypes []*MoveTypeTag
	start := 0
	end := 0
	l := len(genericTypeListString)
	if l == 0 {
		return nil, nil
	}

	for idx := 0; idx < l; idx++ {
		switch genericTypeListString[idx] {
		case '<':
			leftBracketCount += 1
		case '>':
			leftBracketCount -= 1
		case ',':
			if leftBracketCount == 0 {
				end = idx
				aTypeStr := genericTypeListString[start:end]
				aType, err := ParseMoveTypeTag(aTypeStr)
				if err != nil {
					return nil, err
				}
				parsedTypes = append(parsedTypes, aType)
				start = idx + 1
				end = idx + 1
			}
		}
	}

	if end < l-1 {
		aType, err := ParseMoveTypeTag(genericTypeListString[start:l])
		if err != nil {
			return nil, err
		}
		parsedTypes = append(parsedTypes, aType)
	}

	return parsedTypes, nil
}

func mapSlices[E ~[]TIn, TIn any, TOut any](input E, mapper func(TIn) TOut) []TOut {
	var r []TOut
	for _, e := range input {
		r = append(r, mapper(e))
	}

	return r
}

var _ json.Marshaler = (*MoveStructTag)(nil)

func (t *MoveStructTag) MarshalJSON() ([]byte, error) {
	typeName := t.String()
	return json.Marshal(typeName)
}

var _ json.Unmarshaler = (*MoveStructTag)(nil)

func (t *MoveStructTag) UnmarshalJSON(data []byte) error {
	var dataStr string
	err := json.Unmarshal(data, &dataStr)
	if err != nil {
		return err
	}
	return parseMoveStructTagInternal(dataStr, t)
}

// Type is to support cobra value
func (t *MoveStructTag) Type() string {
	return "move-type-tag"
}

// Set is to support cobra value
func (t *MoveStructTag) Set(data string) error {
	return parseMoveStructTagInternal(data, t)
}
