package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-aptos/aptos/known"
	"github.com/mattn/go-isatty"
)

func getOrPanic[T any](v T, err error) T {
	orPanic(err)

	return v
}

func orPanic(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type warnPrinter struct {
	colorRed   string
	colorReset string
}

func newWarnPrinter() *warnPrinter {
	colorRed := ""
	colorReset := ""
	if isatty.IsTerminal(os.Stdin.Fd()) {
		colorReset = "\033[0m"
		colorRed = "\033[31m"
	}

	return &warnPrinter{
		colorRed:   colorRed,
		colorReset: colorReset,
	}
}

var redWarn = newWarnPrinter()

func (rw *warnPrinter) Printf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	fmt.Printf("%s%s%s", redWarn.colorRed, s, redWarn.colorReset)
}

func parseCoinType(network aptos.Network, coinType string) (*aptos.MoveStructTag, error) {
	coinInfo := known.GetCoinInfoBySymbol(network, coinType)
	if coinInfo != nil {
		return coinInfo.TokenType.Type, nil
	}

	coin, err := aptos.ParseMoveStructTag(coinType)
	if err != nil {
		return nil, fmt.Errorf("cannot find %s on aptos %s as a known coin, and failed to parse as a move type: %w", coinType, network, err)
	}

	return coin, nil
}
