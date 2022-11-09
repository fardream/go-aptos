package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"text/template"

	"mvdan.cc/gofumpt/format"
)

//go:generate go run ./

//go:embed aux_stable_pool.go.template
var auxStablePoolGoTemplate string

type coinData struct {
	I       int
	NotLast bool
}
type genData struct {
	N   int
	Xis []*coinData
}

func newGenData(n int) *genData {
	r := &genData{
		N: n,
	}

	for i := 0; i < n; i++ {
		r.Xis = append(r.Xis, &coinData{
			I:       i,
			NotLast: i < n-1,
		})
	}

	return r
}

func main() {
	tmpl, err := template.New("temp").Parse(auxStablePoolGoTemplate)
	if err != nil {
		panic(err)
	}

	for _, n := range []int{2, 3, 4} {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, newGenData(n)); err != nil {
			panic(err)
		}

		formatted, err := format.Source(buf.Bytes(), format.Options{
			LangVersion: "v1.19",
			ModulePath:  "github.com/fardream/go-aptos/aptos",
			ExtraRules:  true,
		})
		if err != nil {
			panic(err)
		}

		if err := os.WriteFile(path.Join("..", fmt.Sprintf("aux_stable_%dpool.go", n)), formatted, 0o666); err != nil {
			panic(err)
		}
	}
}
