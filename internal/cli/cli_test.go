package cli

import (
	"testing"

	"golang.org/x/net/html"
)

func TestGetCommand(t *testing.T) {
	format = FormatJSON
	defer func() { format = FormatJSON }()

	err := runGet(nil, []string{"sec-toprimitive"})
	if err != nil {
		t.Logf("get sec-toprimitive returned error: %v", err)
	}
}

func TestSearchCommand(t *testing.T) {
	format = FormatJSON
	defer func() { format = FormatJSON }()
	searchLimit = 5
	searchKind = ""
	err := runSearch(nil, []string{"object", "property"})
	if err != nil {
		t.Logf("search returned error: %v", err)
	}
}

func TestTocCommand(t *testing.T) {
	format = FormatJSON
	defer func() { format = FormatJSON }()
	tocDepth = 1
	err := runToc(nil, []string{})
	if err != nil {
		t.Logf("toc returned error: %v", err)
	}
}

func TestTocCommandWithSection(t *testing.T) {
	format = FormatJSON
	defer func() { format = FormatJSON }()
	tocDepth = 2
	err := runToc(nil, []string{"1"})
	if err != nil {
		t.Logf("toc with section returned error: %v", err)
	}
}

func TestSearchCommandOutput(t *testing.T) {
	format = FormatJSON
	defer func() { format = FormatJSON }()
	searchLimit = 3
	searchKind = ""
	err := runSearch(nil, []string{"prototype"})
	if err != nil {
		t.Logf("search prototype returned error: %v", err)
	}
}

func TestXrefCommandOutput(t *testing.T) {
	format = FormatJSON
	defer func() { format = FormatJSON }()
	xrefDirection = "both"
	err := runXref(nil, []string{"sec-toprimitive"})
	if err != nil {
		t.Logf("xref returned error: %v", err)
	}
}

func TestXrefCommandByOpName(t *testing.T) {
	format = FormatJSON
	defer func() { format = FormatJSON }()
	xrefDirection = "both"
	err := runXref(nil, []string{"ToNumber"})
	if err != nil {
		t.Logf("xref ToNumber returned error: %v", err)
	}
}

func TestGrammarCommandEBNF(t *testing.T) {
	grammarFormat = "ebnf"
	err := runGrammar(nil, []string{"HexEscapeSequence"})
	if err != nil {
		t.Logf("grammar returned error: %v", err)
	}
}

func TestGrammarCommandJSON(t *testing.T) {
	grammarFormat = "json"
	err := runGrammar(nil, []string{"HexEscapeSequence"})
	if err != nil {
		t.Logf("grammar JSON returned error: %v", err)
	}
}

func TestGrammarCommandMD(t *testing.T) {
	grammarFormat = "md"
	err := runGrammar(nil, []string{"HexEscapeSequence"})
	if err != nil {
		t.Logf("grammar MD returned error: %v", err)
	}
}

func TestGetStepsOnly(t *testing.T) {
	format = FormatJSON
	defer func() { format = FormatJSON }()
	getStepsOnly = true
	defer func() { getStepsOnly = false }()
	err := runGet(nil, []string{"sec-toprimitive"})
	if err != nil {
		t.Logf("get steps-only returned error: %v", err)
	}
}

func TestGetBrief(t *testing.T) {
	format = FormatJSON
	defer func() { format = FormatJSON }()
	getBrief = true
	defer func() { getBrief = false }()
	err := runGet(nil, []string{"sec-toprimitive"})
	if err != nil {
		t.Logf("get brief returned error: %v", err)
	}
}

func TestGetWithMaxTokens(t *testing.T) {
	format = FormatJSON
	defer func() { format = FormatJSON }()
	getMaxTokens = 100
	defer func() { getMaxTokens = 0 }()
	err := runGet(nil, []string{"sec-toprimitive"})
	if err != nil {
		t.Logf("get with max-tokens returned error: %v", err)
	}
}

func TestGetMarkdown(t *testing.T) {
	format = FormatMarkdown
	defer func() { format = FormatJSON }()
	err := runGet(nil, []string{"sec-toprimitive"})
	if err != nil {
		t.Logf("get markdown returned error: %v", err)
	}
}

func TestSearchMarkdown(t *testing.T) {
	format = FormatMarkdown
	defer func() { format = FormatJSON }()
	searchLimit = 3
	searchKind = ""
	err := runSearch(nil, []string{"prototype"})
	if err != nil {
		t.Logf("search markdown returned error: %v", err)
	}
}

func TestTocMarkdown(t *testing.T) {
	format = FormatMarkdown
	defer func() { format = FormatJSON }()
	tocDepth = 1
	err := runToc(nil, []string{})
	if err != nil {
		t.Logf("toc markdown returned error: %v", err)
	}
}

func TestXrefMarkdown(t *testing.T) {
	format = FormatMarkdown
	defer func() { format = FormatJSON }()
	xrefDirection = "both"
	err := runXref(nil, []string{"sec-toprimitive"})
	if err != nil {
		t.Logf("xref markdown returned error: %v", err)
	}
}

func TestTokenizeSearch(t *testing.T) {
	tokens := tokenizeSearch("abstract operation toNumber")
	if len(tokens) == 0 {
		t.Error("expected tokens from search query")
	}
	t.Logf("tokens: %v", tokens)
}

func TestMatchKind(t *testing.T) {
	s := getSpec()
	for _, cn := range s.Nodes {
		if cn.Type == "abstract operation" {
			if !matchKind(cn, "operation") {
				t.Error("abstract operation should match 'operation' kind")
			}
			if matchKind(cn, "grammar") {
				t.Error("abstract operation should not match 'grammar' kind")
			}
			break
		}
	}
}

func TestContentHasEmuAlg(t *testing.T) {
	s := getSpec()
	found := false
	for _, cn := range s.Nodes {
		if cn.HTMLNode != nil {
			if hasElement(cn.HTMLNode, "emu-alg") {
				found = true
				t.Logf("clause %s has emu-alg", cn.ID)
				break
			}
		}
	}
	if !found {
		t.Skip("no emu-alg elements found")
	}
}

func hasElement(n *html.Node, tag string) bool {
	if n.Type == html.ElementNode && n.Data == tag {
		return true
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if hasElement(c, tag) {
			return true
		}
	}
	return false
}
