// ecma-query is a CLI for querying the ECMAScript specification.
package main

import (
	"os"

	"github.com/podikoglou/ecma-query/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
