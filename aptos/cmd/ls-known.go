package cmd

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/fardream/go-aptos/aptos/known"
)

func GetListKnownCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls-known",
		Short: "list known coins, market, amms",
		Long:  "list known coins market amms etc",
	}

	cmd.Run = func(*cobra.Command, []string) {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Network", "Symbol", "Name", "Decimal", "Move Type"})

		allCoins := known.GetAllCoins()

		for network, coinmap := range allCoins {
			for _, coin := range *coinmap {
				table.Append(
					[]string{
						network.String(),
						coin.Symbol,
						coin.Name,
						strconv.Itoa(int(coin.Decimals)),
						coin.TokenType.Type.String(),
					},
				)
			}
		}

		table.Render()
	}
	return cmd
}
