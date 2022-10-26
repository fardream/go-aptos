package cmd

import (
	"github.com/fardream/go-aptos/aptos"
	"github.com/spf13/cobra"
)

type SharedArgs struct {
	network      aptos.Network
	profile      string
	endpoint     string
	maxGasAmount uint64
	simulate     bool
}

func NewSharedArgs() *SharedArgs {
	return &SharedArgs{
		network:      aptos.Devnet,
		profile:      "",
		endpoint:     "",
		maxGasAmount: 200000,
		simulate:     false,
	}
}

func (args *SharedArgs) SetCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().VarP(&args.network, "network", "n", "network for the market.")
	cmd.PersistentFlags().StringVarP(&args.endpoint, "endpoint", "u", args.endpoint, "endpoint for the rest api, default to the one provided by aptos labs.")
	cmd.PersistentFlags().Uint64VarP(&args.maxGasAmount, "max-gas-amount", "m", args.maxGasAmount, "max gas amount - make sure the account has enough aptos liquidity.")
	cmd.PersistentFlags().StringVarP(&args.profile, "profile", "k", args.profile, "aptos profile to use. if a network is selected but this is unset, will use profile with that name.")
	cmd.PersistentFlags().BoolVarP(&args.simulate, "simulate", "s", args.simulate, "simulate the transaction")
}

type SharedArgsWithBaseQuoteCoins struct {
	*SharedArgs
	baseCoinStr  string
	quoteCoinStr string
}

func NewSharedArgsWithBaseQuoteCoins() *SharedArgsWithBaseQuoteCoins {
	return &SharedArgsWithBaseQuoteCoins{
		SharedArgs: NewSharedArgs(),
	}
}

func (args *SharedArgsWithBaseQuoteCoins) SetCmd(cmd *cobra.Command) {
	args.SharedArgs.SetCmd(cmd)
	cmd.PersistentFlags().StringVarP(&args.baseCoinStr, "base-coin", "b", args.baseCoinStr, "base coin for the market")
	cmd.MarkPersistentFlagRequired("base-coin")
	cmd.PersistentFlags().StringVarP(&args.quoteCoinStr, "quote-coin", "q", args.quoteCoinStr, "quote coin for the market")
	cmd.MarkPersistentFlagRequired("quote-coin")
}

func (args *SharedArgs) UpdateProfileForCmd(cmd *cobra.Command) {
	if !cmd.PersistentFlags().Changed("profile") {
		args.profile = string(args.network)
	}
}
