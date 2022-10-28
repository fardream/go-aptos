package cmd

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos"
)

func GetListAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls-account",
		Short: "listing the account resources",
		Args:  cobra.NoArgs,
	}

	network := aptos.Mainnet
	address := aptos.Address{}
	cmd.Flags().VarP(&network, "network", "c", "network")
	cmd.Flags().VarP(&address, "address", "a", "address")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		r, _, err := aptos.GetDefaultEndpoint(network)
		if err != nil {
			orPanic(err)
		}

		client := aptos.MustNewClient(network, r)
		resp := getOrPanic(client.GetAccountResources(context.Background(), &aptos.GetAccountResourcesRequest{
			Address: address,
		}))

		spew.Dump(resp)
	}
	return cmd
}
