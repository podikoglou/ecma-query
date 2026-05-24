package cli

import (
	"fmt"
	"strings"

	"github.com/podikoglou/ecma-query/internal/spec"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <QUERY>",
	Short: "Full-text search across the spec",
	Long:  "Returns ranked results with snippets for exploratory queries.",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runSearch,
}

var (
	searchLimit int
	searchKind  string
)

func init() {
	searchCmd.Flags().IntVar(&searchLimit, "limit", 10, "maximum number of results")
	searchCmd.Flags().StringVar(&searchKind, "kind", "", "filter by entity kind: operation, method, section, type, grammar, slot")
}

func runSearch(_ *cobra.Command, args []string) error {
	query := strings.Join(args, " ")
	s := getSpec()

	tokens := tokenizeSearch(query)
	if len(tokens) == 0 {
		return nil
	}

	scores := map[string]int{}
	for _, tok := range tokens {
		clauseMap, ok := s.TextIndex[tok]
		if !ok {
			continue
		}
		for clauseID, freq := range clauseMap {
			scores[clauseID] += freq
		}
	}

	type result struct {
		cn    *spec.ClauseNode
		score int
	}
	var results []result
	for clauseID, score := range scores {
		cn, ok := s.ByID[clauseID]
		if !ok {
			continue
		}
		results = append(results, result{cn, score})
	}

	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].score < results[j].score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	var filtered []result
	for _, r := range results {
		if searchKind != "" && !matchKind(r.cn, searchKind) {
			continue
		}
		filtered = append(filtered, r)
	}
	if searchKind != "" {
		results = filtered
	}

	limit := searchLimit
	if limit <= 0 {
		limit = 10
	}
	if limit > len(results) {
		limit = len(results)
	}

	for i := 0; i < limit; i++ {
		r := results[i]
		snippet := extractSnippet(r.cn, tokens[0], 60)
		md := "## " + r.cn.H1Text + "\n"
		if r.cn.Section != "" {
			md += "Section: " + r.cn.Section + "\n"
		}
		md += "Kind: " + kindLabel(r.cn) + "\n"
		if snippet != "" {
			md += "> " + snippet + "\n"
		}
		url := clauseURL(r.cn.ID)
		if url != "" {
			md += url + "\n"
		}
		md += "\n"
		fmt.Print(md)
	}
	return nil
}

func tokenizeSearch(query string) []string {
	query = strings.ToLower(query)
	var tokens []string
	var buf strings.Builder
	for _, r := range query {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			buf.WriteRune(r)
		} else {
			if buf.Len() >= 3 {
				tokens = append(tokens, buf.String())
			}
			buf.Reset()
		}
	}
	if buf.Len() >= 3 {
		tokens = append(tokens, buf.String())
	}
	return tokens
}

func matchKind(cn *spec.ClauseNode, kind string) bool {
	switch strings.ToLower(kind) {
	case "operation":
		return cn.Type == "abstract operation" || cn.Type == "syntax-directed operation" || cn.Type == "concrete method" || cn.Type == "host-defined abstract operation"
	case "method":
		return cn.Type == "built-in function" || cn.Type == "built-in method"
	case "section":
		return cn.Type == "" && cn.Kind == "clause"
	case "type":
		return cn.Type == "spec type"
	case "grammar":
		return cn.Type == "grammar"
	case "slot":
		return cn.Type == "internal slot"
	}
	return true
}
