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
	Long:  "Retrieve grammar productions. Defaults to EBNF output (instead of the global JSON default).",
	Args:  cobra.ExactArgs(1),
	RunE:  runGrammar,
}

var grammarFormat string

func init() {
	grammarCmd.Flags().StringVar(&grammarFormat, "format", "ebnf", "json, md, or ebnf")
}

func runGrammar(_ *cobra.Command, args []string) error {
	query := args[0]
	s := getSpec()

	prods, ok := s.Grammars[query]
	if !ok {
		ExitNotFound(query)
		return nil
	}

	switch grammarFormat {
	case "json":
		resp := buildGrammarJSON(query, prods)
		printJSON(resp)
	case "md":
		renderGrammarMD(query, prods)
	default:
		renderGrammarEBNF(prods)
	}
	return nil
}

func buildGrammarJSON(name string, prods []spec.GrammarProd) GrammarResponse {
	resp := GrammarResponse{
		Name: name,
		Kind: "syntactic",
	}
	if len(prods) > 0 {
		resp.Section = prods[0].Section
		for _, gp := range prods {
			item := parseGrammarProd(gp.HTML)
			resp.Productions = append(resp.Productions, item)
		}
	}
	if strings.HasPrefix(name, "::") || strings.Contains(name, "Lexical") {
		resp.Kind = "lexical"
	}
	for _, gp := range prods {
		if strings.Contains(gp.HTML, "lexical") || strings.Contains(gp.Section, "lexical") {
			resp.Kind = "lexical"
			break
		}
	}
	return resp
}

func parseGrammarProd(text string) GrammarItem {
	lines := strings.Split(text, "\n")
	var item GrammarItem
	var currentAlt []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.Contains(trimmed, "::") {
			parts := strings.SplitN(trimmed, "::", 2)
			if item.LHS == "" {
				item.LHS = strings.TrimSpace(parts[0])
			}
			rhs := strings.TrimSpace(parts[1])
			if rhs != "" {
				currentAlt = parseRHSTerms(rhs)
				item.Alternatives = append(item.Alternatives, currentAlt)
			}
		} else {
			currentAlt = parseRHSTerms(trimmed)
			item.Alternatives = append(item.Alternatives, currentAlt)
		}
	}
	return item
}

func parseRHSTerms(rhs string) []string {
	var terms []string
	rhs = strings.TrimSpace(rhs)
	i := 0
	for i < len(rhs) {
		switch rhs[i] {
		case '`':
			j := i + 1
			for j < len(rhs) && rhs[j] != '`' {
				j++
			}
			if j < len(rhs) {
				terms = append(terms, rhs[i+1:j])
				i = j + 1
			} else {
				terms = append(terms, rhs[i:])
				i = len(rhs)
			}
		case ' ', '\t':
			i++
		default:
			j := i
			for j < len(rhs) && rhs[j] != ' ' && rhs[j] != '\t' && rhs[j] != '`' {
				j++
			}
			term := strings.TrimSpace(rhs[i:j])
			if term != "" {
				terms = append(terms, term)
			}
			i = j
		}
	}
	return terms
}

func renderGrammarEBNF(prods []spec.GrammarProd) {
	for _, gp := range prods {
		fmt.Println(strings.TrimSpace(gp.HTML))
		fmt.Println()
	}
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
