package aptos_test

import (
	"math/big"
	"testing"

	"github.com/fardream/go-aptos/aptos"
)

func TestNewUint128FromUint64(t *testing.T) {
	expected := &big.Int{}
	expected.Add(big.NewInt(0).Lsh(big.NewInt(96), 64), big.NewInt(50))
	result := aptos.NewUint128FromUint64(50, 96)
	if result.Big().Cmp(expected) != 0 {
		t.Fatalf("want: %s, got: %s", expected.String(), result.Big().String())
	}
}

func TestNewBigIntFromUint64(t *testing.T) {
	var expected uint64 = (1 << 63) + 12345
	result := aptos.NewBigIntFromUint64(expected)

	if result.Uint64() != expected {
		t.Fatalf("want: %d, got: %s", expected, result.String())
	}
}
