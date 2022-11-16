package aptos

import (
	"fmt"
	"strings"
)

//go:generate stringer -type AuxFakeCoin -linecomment

// AuxFakeCoin contains some fake coins to use on devnet and testnet.
// They don't have any value and can be freely minted to anyone.
// Simply call mint (or register_and_mint if not signed up for it already).
type AuxFakeCoin int

const (
	AuxFakeCoin_USDC   AuxFakeCoin = iota // USDC
	AuxFakeCoin_ETH                       // ETH
	AuxFakeCoin_BTC                       // BTC
	AuxFakeCoin_SOL                       // SOL
	AuxFakeCoin_AUX                       // AUX
	AuxFakeCoin_USDT                      // USDT
	AuxFakeCoin_USDCD8                    // USDCD8
)

// AuxAllFakeCoins contains all the fake coins provided by aux for testing
var AuxAllFakeCoins []AuxFakeCoin = []AuxFakeCoin{
	AuxFakeCoin_USDC,
	AuxFakeCoin_ETH,
	AuxFakeCoin_BTC,
	AuxFakeCoin_SOL,
	AuxFakeCoin_USDT,
	AuxFakeCoin_AUX,
	AuxFakeCoin_USDCD8,
}

// GetAuxFakeCoinType returns the fake coin type. Note, this is actually not a type of coin as defined by aptos framework.
// Use [GetAuxFakeCoinCoinType] to get the coin type.
func GetAuxFakeCoinType(moduleAddress Address, fakeCoin AuxFakeCoin) (*MoveTypeTag, error) {
	if fakeCoin >= AuxFakeCoin(len(_AuxFakeCoin_index)-1) {
		return nil, fmt.Errorf("unknown fake coin: %d", int(fakeCoin))
	}

	return NewMoveTypeTag(
		moduleAddress,
		AuxFakeCoinModuleName,
		fakeCoin.String(),
		nil)
}

// GetAuxFakeCoinCoinType returns the fake coin **coin** type - **this is a coin as defined by the aptos framework.**
func GetAuxFakeCoinCoinType(moduleAddress Address, fakeCoin AuxFakeCoin) (*MoveTypeTag, error) {
	coinType, err := GetAuxFakeCoinType(moduleAddress, fakeCoin)
	if err != nil {
		return nil, err
	}
	return NewMoveTypeTag(
		moduleAddress,
		AuxFakeCoinModuleName,
		"FakeCoin",
		[]*MoveTypeTag{coinType},
	)
}

// GetAuxFakeCoinDecimal provides the decimals for fake coins
func GetAuxFakeCoinDecimal(fakeCoin AuxFakeCoin) uint8 {
	switch fakeCoin {
	case AuxFakeCoin_BTC, AuxFakeCoin_ETH, AuxFakeCoin_SOL, AuxFakeCoin_USDCD8:
		return 8
	case AuxFakeCoin_USDC, AuxFakeCoin_USDT, AuxFakeCoin_AUX:
		return 6
	}

	return 0
}

// ParseAuxFakeCoin converts a string into fake coin.
func ParseAuxFakeCoin(s string) (AuxFakeCoin, error) {
	switch strings.ToUpper(s) {
	case "ETH":
		return AuxFakeCoin_ETH, nil
	case "BTC":
		return AuxFakeCoin_BTC, nil
	case "AUX":
		return AuxFakeCoin_AUX, nil
	case "USDC":
		return AuxFakeCoin_USDC, nil
	case "USDT":
		return AuxFakeCoin_USDT, nil
	case "SOL":
		return AuxFakeCoin_SOL, nil
	case "USDCD8":
		return AuxFakeCoin_USDCD8, nil
	default:
		return AuxFakeCoin_USDC, fmt.Errorf("failed to recognize the fake coin: %s", s)
	}
}

// AuxFakeCoinModuleName is the module name for fake coin.
const AuxFakeCoinModuleName = "fake_coin"

// FakeCoin_RegisterAndMint register and mint fake coins. Any signer can self sign and get those coins. If the sender is not registered, this operation will
// register user for the coin. If the sender is registered for the coin, it will simply mint.
func (info *AuxClientConfig) FakeCoin_RegisterAndMint(sender Address, fakeCoin AuxFakeCoin, amount uint64, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxFakeCoinModuleName, "register_and_mint")

	tx := &Transaction{
		Payload: NewEntryFunctionPayload(
			function,
			[]*MoveTypeTag{
				must(GetAuxFakeCoinType(info.Address, fakeCoin)),
			},
			[]*EntryFunctionArg{
				EntryFunctionArg_Uint64(amount),
			}),
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// FakeCoin_Register registers the user for the fake coin. No effect if the user is already registered.
func (info *AuxClientConfig) FakeCoin_Register(sender Address, fakeCoin AuxFakeCoin, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxFakeCoinModuleName, "register")

	tx := &Transaction{
		Payload: NewEntryFunctionPayload(
			function,
			[]*MoveTypeTag{
				must(GetAuxFakeCoinType(info.Address, fakeCoin)),
			},
			nil,
		),
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// FakeCoin_Mint mints coins to the user. The user must be registered.
func (info *AuxClientConfig) FakeCoin_Mint(sender Address, fakeCoin AuxFakeCoin, amount uint64, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxFakeCoinModuleName, "mint")

	tx := &Transaction{
		Payload: NewEntryFunctionPayload(
			function,
			[]*MoveTypeTag{
				must(GetAuxFakeCoinType(info.Address, fakeCoin)),
			},
			[]*EntryFunctionArg{
				EntryFunctionArg_Uint64(amount),
			}),
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}

// FakeCoin_Burn burns the fake coins for a user. The is useful when tests require users' balances must start from zero.
func (info *AuxClientConfig) FakeCoin_Burn(sender Address, fakeCoin AuxFakeCoin, amount uint64, options ...TransactionOption) *Transaction {
	function := MustNewMoveFunctionTag(info.Address, AuxFakeCoinModuleName, "burn")

	tx := &Transaction{
		Payload: NewEntryFunctionPayload(
			function,
			[]*MoveTypeTag{
				must(GetAuxFakeCoinType(info.Address, fakeCoin)),
			},
			[]*EntryFunctionArg{
				EntryFunctionArg_Uint64(amount),
			}),
	}

	ApplyTransactionOptions(tx, options...)

	tx.Sender = sender

	return tx
}
