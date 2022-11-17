package aptos_test

import (
	"testing"

	"github.com/fardream/go-aptos/aptos"
	"github.com/google/go-cmp/cmp"
)

var parseMoveTagTestCases = []struct {
	typeStr  string
	expected aptos.MoveTypeTag
}{
	{
		typeStr:  "u8",
		expected: aptos.MoveTypeTag{Uint8: &struct{}{}},
	},
	{
		typeStr:  "vector<u8>",
		expected: aptos.MoveTypeTag{Vector: &aptos.MoveTypeTag{Uint8: &struct{}{}}},
	},
	{
		typeStr: "vector<0x1::coin::CoinInfo<vector<u8>, 0x1::coin::Coin, u128>>",
		expected: aptos.MoveTypeTag{
			Vector: &aptos.MoveTypeTag{
				Struct: &aptos.MoveStructTag{
					MoveModuleTag: aptos.MoveModuleTag{
						Module: "coin", Address: aptos.AptosStdAddress,
					},
					Name: "CoinInfo",
					GenericTypeParameters: []*aptos.MoveTypeTag{
						{
							Vector: &aptos.MoveTypeTag{
								Uint8: &struct{}{},
							},
						},
						{
							Struct: &aptos.MoveStructTag{
								MoveModuleTag: aptos.MoveModuleTag{
									Address: aptos.AptosStdAddress,
									Module:  "coin",
								},
								Name: "Coin",
							},
						},
						{
							Uint128: &struct{}{},
						},
					},
				},
			},
		},
	},
}

func TestParseMoveTypeTag(t *testing.T) {
	for _, aCase := range parseMoveTagTestCases {
		got, err := aptos.ParseMoveTypeTag(aCase.typeStr)
		if err != nil {
			t.Errorf("failed to parse %s: %v", aCase.typeStr, err)
		}
		if !cmp.Equal(*got, aCase.expected) {
			t.Errorf("got:  %v\nwant: %v\n", got, aCase.expected)
		}
	}
}
