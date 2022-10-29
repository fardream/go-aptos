package aptos

import "encoding/json"

// MoveModuleABI is the smart contract module [abi].
//
// [abi]: https://fullnode.mainnet.aptoslabs.com/v1/spec#/schemas/MoveModule
type MoveModuleABI struct {
	Address Address `json:"address"`
	Name    string  `json:"name"`

	Friends []string `json:"friends"`

	Functions []*MoveModuleABI_Function `json:"exposed_functions"`
	Structs   []*MoveModuleABI_Struct   `json:"structs"`
}

// MoveModuleABI_Function is the function definition contained in a [MoveModuleABI]
type MoveModuleABI_Function struct {
	Name              string            `json:"name"`
	Visibility        string            `json:"visibility"`
	IsEntry           bool              `json:"is_entry"`
	GenericTypeParams []json.RawMessage `json:"generic_type_params"`
	Params            []string          `json:"params"`
	Return            []string          `json:"return"`
}

// MoveModuleABI_Struct is the struct definition contained in a [MoveModuleABI]
type MoveModuleABI_Struct struct {
	Name              string            `json:"name"`
	IsNative          bool              `json:"is_native"`
	Abilities         []string          `json:"abilities"`
	GenericTypeParams []json.RawMessage `json:"generic_type_params"`
	Fields            []json.RawMessage `json:"fields"`
}
