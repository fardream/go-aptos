package aptos_test

import "github.com/fardream/go-aptos/aptos"

var (
	devnetUrl    string
	devnetFaucet string
	devnetConfig *aptos.AuxClientConfig
	trader       *aptos.LocalAccount
)

const traderMnemonic = "apart canvas monitor nephew certain oxygen happy answer element oven cup ladder"

func init() {
	devnetUrl, devnetFaucet, _ = aptos.GetDefaultEndpoint(aptos.Devnet)
	devnetConfig, _ = aptos.GetAuxClientConfig(aptos.Devnet)
	trader, _ = aptos.NewLocalAccountFromMnemonic(traderMnemonic, "")
}
