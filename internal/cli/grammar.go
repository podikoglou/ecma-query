package cli

import (
	"fmt"
	"strings"

	"github.com/podikoglou/ecma-query/internal/spec"
	"github.com/spf13/cobra"
)

var grammarCmd = &cobra.Command{
	Use:   "grammar <PRODUCTION-NAME>",
	Short: "Grammar production query",
	Long:  "Retrieve grammar productions.",
	Args:  cobra.ExactArgs(1),
	RunE:  runGrammar,
}

func runGrammar(_ *cobra.Command, args []string) error {
	query := args[0]
	s := getSpec()

	prods, ok := s.Grammars[query]
	if !ok {
		exitNotFound(query)
		return nil
	}

	renderGrammarMD(query, prods)
	return nil
}

func renderGrammarMD(name string, prods []spec.GrammarProd) {
	fmt.Printf("## %s\n\n", name)
	for _, gp := range prods {
		fmt.Println("```")
		fmt.Println(strings.TrimSpace(gp.HTML))
		fmt.Println("```")
		fmt.Println()
	}
}
