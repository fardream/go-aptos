package aptos

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type Uint128 struct {
	underlying big.Int
}

var (
	_ json.Marshaler   = (*Uint128)(nil)
	_ json.Unmarshaler = (*Uint128)(nil)
	_ EntryFunctionArg = (*Uint128)(nil)
)

func (i Uint128) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.underlying.String())
}

var maxU128 = (&big.Int{}).Lsh(big.NewInt(1), 64)

func (i *Uint128) check() error {
	if i.underlying.Sign() < 0 {
		return fmt.Errorf("%s is negative", i.underlying.String())
	}

	if i.underlying.Cmp(maxU128) >= 0 {
		return fmt.Errorf("%s is greater than Max Uint 128", i.underlying.String())
	}

	return nil
}

func (i *Uint128) UnmarshalJSON(data []byte) error {
	if err := i.underlying.UnmarshalJSON(data); err != nil {
		return err
	}

	return i.check()
}

func NewEntryFunctionArg_Uint128(s string) (*Uint128, error) {
	r := &big.Int{}
	r, ok := r.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse %s as an integer", s)
	}

	i := &Uint128{underlying: *r}
	if err := i.check(); err != nil {
		return nil, err
	}

	return i, nil
}

func (i Uint128) ToBCS() []byte {
	r := make([]byte, 16)
	bigEndianBytes := i.underlying.Bytes()
	copy(r, bigEndianBytes)
	n := len(bigEndianBytes)
	for j := 0; j < n/2; j++ {
		r[j], r[n-1-j] = r[n-1-j], r[j]
	}

	return r
}
