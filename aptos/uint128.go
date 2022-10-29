package aptos

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// Uint128, for convenience of implementation, this uses [big.Int]
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

var maxU128 = (&big.Int{}).Lsh(big.NewInt(1), 128)

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
	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err != nil {
		return err
	}

	_, ok := i.underlying.SetString(dataStr, 10)
	if !ok {
		return fmt.Errorf("failed to parse %s as an integer", dataStr)
	}

	return i.check()
}

func NewUint128(s string) (*Uint128, error) {
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

func (i *Uint128) Cmp(j *Uint128) int {
	return i.underlying.Cmp(&j.underlying)
}

var zeroUint128 = Uint128{underlying: *big.NewInt(0)}

// 63 ones
const ones63 uint64 = (1 << 63) - 1

// 1 << 63
var oneLsh63 = big.NewInt(0).Lsh(big.NewInt(1), 63)

func NewBigIntFromUint64(i uint64) *big.Int {
	r := big.NewInt(int64(i & ones63))
	if i > ones63 {
		r = r.Add(r, oneLsh63)
	}
	return r
}

func NewUint128FromUint64(lo, hi uint64) *Uint128 {
	loBig := NewBigIntFromUint64(lo)
	hiBig := NewBigIntFromUint64(hi)
	hiBig = hiBig.Lsh(hiBig, 64)
	return &Uint128{
		underlying: *hiBig.Add(hiBig, loBig),
	}
}

func (u *Uint128) Big() *big.Int {
	return (&big.Int{}).Set(&u.underlying)
}

func (u Uint128) String() string {
	return u.underlying.String()
}
