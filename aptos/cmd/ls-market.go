package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-aptos/aptos/known"
)

func GetListMarketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls-market",
		Short: "get all markets for aux",
		Args:  cobra.NoArgs,
	}

	network := aptos.Devnet
	endpoint := ""
	cmd.Flags().VarP(&network, "network", "n", "network for the market.")
	cmd.Flags().StringVarP(&endpoint, "endpoint", "u", endpoint, "endpoint for the rest api, default to the one provided by aptos labs.")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if endpoint == "" {
			var err error
			endpoint, _, err = aptos.GetDefaultEndpoint(network)
			orPanic(err)
		}

		auxConfig := getOrPanic(aptos.GetAuxClientConfig(network))

		moduleAddress := auxConfig.Address

		client := aptos.MustNewClient(network, endpoint)

		resources := getOrPanic(client.GetAccountResources(context.Background(), &aptos.GetAccountResourcesRequest{
			Address: moduleAddress,
		}))

		for _, resource := range *resources.Parsed {
			resourceType := resource.Type

			if resourceType.Module == "clob_market" && resourceType.Name == "Market" && cmp.Equal(auxConfig.Address, resourceType.Address) {
				var market aptos.AuxClobMarket
				if err := json.Unmarshal(resource.Data, &market); err != nil {
					fmt.Printf("failed to parse %s due to %v\n", string(resource.Data), err)
					continue

				}
				baseName := resourceType.GenericTypeParameters[0].String()
				if baseCoin := known.GetCoinInfo(network, resourceType.GenericTypeParameters[0].Struct); baseCoin != nil {
					baseName = baseCoin.Name
				}
				quoteName := resourceType.GenericTypeParameters[1].String()
				if quoteCoin := known.GetCoinInfo(network, resourceType.GenericTypeParameters[1].Struct); quoteCoin != nil {
					quoteName = quoteCoin.Name
				}
				fmt.Printf("found market for %s (decimal: %d) <-> %s (decimal: %d) lot size: %d, tick size: %d\n",
					baseName, market.BaseDecimals, quoteName, market.QuoteDecimals, market.LotSize, market.TickSize,
				)
			}
		}
	}

	return cmd
}
