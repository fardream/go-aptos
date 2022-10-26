package aptos

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type EntryFunctionArg_Uint128 struct {
	underlying big.Int
}

var (
	_ json.Marshaler   = (*EntryFunctionArg_Uint128)(nil)
	_ json.Unmarshaler = (*EntryFunctionArg_Uint128)(nil)
	_ EntryFunctionArg = (*EntryFunctionArg_Uint128)(nil)
)

func (i EntryFunctionArg_Uint128) MarshalJSON() ([]byte, error) {
	return i.underlying.MarshalJSON()
}

var maxU128 = (&big.Int{}).Lsh(big.NewInt(1), 64)

func (i *EntryFunctionArg_Uint128) check() error {
	if i.underlying.Sign() < 0 {
		return fmt.Errorf("%s is negative", i.underlying.String())
	}

	if i.underlying.Cmp(maxU128) >= 0 {
		return fmt.Errorf("%s is greater than Max Uint 128", i.underlying.String())
	}

	return nil
}

func (i *EntryFunctionArg_Uint128) UnmarshalJSON(data []byte) error {
	if err := i.underlying.UnmarshalJSON(data); err != nil {
		return err
	}

	return i.check()
}

func NewEntryFunctionArg_Uint128(s string) (*EntryFunctionArg_Uint128, error) {
	r := &big.Int{}
	r, ok := r.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse %s as an integer", s)
	}

	i := &EntryFunctionArg_Uint128{underlying: *r}
	if err := i.check(); err != nil {
		return nil, err
	}

	return i, nil
}

func (i EntryFunctionArg_Uint128) ToBCS() []byte {
	bigEndianBytes := i.underlying.Bytes()
	bigEndianBytes = bigEndianBytes[:8]
	for j := 0; j < 4; j++ {
		bigEndianBytes[j], bigEndianBytes[7-j] = bigEndianBytes[7-j], bigEndianBytes[j]
	}

	return bigEndianBytes
}
