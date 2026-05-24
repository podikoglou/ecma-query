package cli

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/podikoglou/ecma-query/internal/spec"
	"golang.org/x/net/html"
)

var specOnce sync.Once
var theSpec *spec.Spec

func getSpec() *spec.Spec {
	specOnce.Do(func() {
		var err error
		theSpec, err = spec.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to load spec:", err)
			os.Exit(1)
		}
	})
	return theSpec
}

func resolve(s *spec.Spec, query string) (*spec.ClauseNode, []*spec.ClauseNode, []MatchResult) {
	if cn, ok := s.SectionNums[query]; ok {
		return cn, nil, nil
	}

	if cn, ok := s.ByID[query]; ok {
		return cn, nil, nil
	}

	norm := normalizeQuery(query)
	if nodes, ok := s.OpsByName[norm]; ok {
		if len(nodes) == 1 {
			return nodes[0], nil, nil
		}
		matches := make([]MatchResult, 0, len(nodes))
		for _, n := range nodes {
			matches = append(matches, MatchResult{
				Kind:    n.Type,
				Name:    n.Name,
				Section: n.Section,
			})
		}
		return nil, nodes, matches
	}

	bracketed := query
	if !strings.HasPrefix(query, "[[") {
		bracketed = "[[" + query + "]]"
	}
	bracketedLower := strings.ToLower(bracketed)
	for _, cn := range s.Nodes {
		if cn.H1Text != "" && strings.Contains(strings.ToLower(cn.H1Text), bracketedLower) {
			return cn, nil, nil
		}
	}

	queryLower := strings.ToLower(query)
	for _, cn := range s.Nodes {
		if cn.Type != "" && strings.ToLower(cn.Type) == queryLower {
			return cn, nil, nil
		}
	}

	return nil, nil, nil
}

func normalizeQuery(q string) string {
	q = strings.ToLower(q)
	q = strings.ReplaceAll(q, " ", "")
	q = strings.ReplaceAll(q, "(", "")
	q = strings.ReplaceAll(q, ")", "")
	return q
}

func suggestions(query string, s *spec.Spec) []string {
	type candidate struct {
		name  string
		score int
	}

	queryLower := strings.ToLower(strings.TrimSpace(query))
	var candidates []candidate

	for key := range s.OpsByName {
		if len(key) < 2 {
			continue
		}
		dist := levenshtein(queryLower, key)
		candidates = append(candidates, candidate{name: key, score: dist})
	}

	for sec := range s.SectionNums {
		if len(sec) < 1 {
			continue
		}
		dist := levenshtein(queryLower, sec)
		candidates = append(candidates, candidate{name: sec, score: dist})
	}

	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].score > candidates[j].score {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	seen := map[string]bool{}
	var result []string
	for _, c := range candidates {
		if len(result) >= 3 {
			break
		}
		if seen[c.name] {
			continue
		}
		seen[c.name] = true
		result = append(result, c.name)
	}
	return result
}

func levenshtein(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	d := make([][]int, la+1)
	for i := range d {
		d[i] = make([]int, lb+1)
		d[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		d[0][j] = j
	}

	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			d[i][j] = min3(d[i-1][j]+1, d[i][j-1]+1, d[i-1][j-1]+cost)
		}
	}
	return d[la][lb]
}

func min3(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= a && b <= c {
		return b
	}
	return c
}

func extractSteps(node *spec.ClauseNode) []string {
	if node.HTMLNode == nil {
		return nil
	}
	var steps []string
	var findAlg func(*html.Node)
	findAlg = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "emu-alg" {
			collectListItems(n, &steps)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findAlg(c)
		}
	}
	findAlg(node.HTMLNode)
	return steps
}

func collectListItems(n *html.Node, steps *[]string) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "ol" {
			collectListItems(c, steps)
		}
		if c.Type == html.ElementNode && c.Data == "li" {
			text := strings.TrimSpace(extractText(c))
			if text != "" {
				*steps = append(*steps, text)
			}
			collectListItems(c, steps)
		}
	}
}

func extractSummary(node *spec.ClauseNode) string {
	if node.HTMLNode == nil {
		return ""
	}
	n := node.HTMLNode
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "p" {
			return strings.TrimSpace(extractText(c))
		}
	}
	return ""
}

func buildBreadcrumb(node *spec.ClauseNode) []string {
	if node == nil {
		return nil
	}
	var parts []string
	for p := node.Parent; p != nil; p = p.Parent {
		t := p.H1Text
		if t == "" {
			t = p.Section
		}
		if t != "" {
			parts = append([]string{t}, parts...)
		}
	}
	return parts
}

func clauseURL(id string) string {
	if id == "" {
		return ""
	}
	return "https://tc39.es/ecma262/#" + id
}

func estimateTokens(text string) int {
	words := strings.Fields(text)
	return int(float64(len(words)) / 0.75)
}

func kindLabel(cn *spec.ClauseNode) string {
	if cn.Type != "" {
		return cn.Type
	}
	if cn.Kind == "annex" {
		return "annex"
	}
	if cn.Kind == "intro" {
		return "intro"
	}
	return "section"
}

func extractText(n *html.Node) string {
	var buf strings.Builder
	var collect func(*html.Node)
	collect = func(nn *html.Node) {
		if nn.Type == html.TextNode {
			buf.WriteString(nn.Data)
		}
		for c := nn.FirstChild; c != nil; c = c.NextSibling {
			collect(c)
		}
	}
	collect(n)
	return collapseWhitespace(buf.String())
}

func collapseWhitespace(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	prevSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace {
				buf.WriteByte(' ')
				prevSpace = true
			}
		} else {
			buf.WriteRune(r)
			prevSpace = false
		}
	}
	return strings.TrimSpace(buf.String())
}

func extractSnippet(node *spec.ClauseNode, token string, window int) string {
	if node.HTMLNode == nil {
		return ""
	}
	raw := extractText(node.HTMLNode)
	lower := strings.ToLower(raw)
	idx := strings.Index(lower, strings.ToLower(token))
	if idx < 0 {
		if len(raw) > window*2 {
			return raw[:window*2]
		}
		return raw
	}

	start := idx - window
	if start < 0 {
		start = 0
	}
	end := idx + len(token) + window
	if end > len(raw) {
		end = len(raw)
	}

	snippet := raw[start:end]
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(raw) {
		snippet = snippet + "..."
	}
	return snippet
}
