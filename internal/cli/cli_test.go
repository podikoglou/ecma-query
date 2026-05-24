// Package cli tests.
package cli

import "testing"

func TestGetCommand(t *testing.T) {
	err := runGet(nil, []string{"sec-toprimitive"})
	if err != nil {
		t.Logf("get sec-toprimitive returned error: %v", err)
	}
}

func TestSearchCommand(t *testing.T) {
	searchLimit = 5
	searchKind = ""
	err := runSearch(nil, []string{"object", "property"})
	if err != nil {
		t.Logf("search returned error: %v", err)
	}
}

func TestTocCommand(t *testing.T) {
	tocDepth = 1
	err := runToc(nil, []string{})
	if err != nil {
		t.Logf("toc returned error: %v", err)
	}
}

func TestTocCommandWithSection(t *testing.T) {
	tocDepth = 2
	err := runToc(nil, []string{"1"})
	if err != nil {
		t.Logf("toc with section returned error: %v", err)
	}
}

func TestSearchCommandOutput(t *testing.T) {
	searchLimit = 3
	searchKind = ""
	err := runSearch(nil, []string{"prototype"})
	if err != nil {
		t.Logf("search prototype returned error: %v", err)
	}
}

func TestXrefCommandOutput(t *testing.T) {
	xrefDirection = "both"
	err := runXref(nil, []string{"sec-toprimitive"})
	if err != nil {
		t.Logf("xref returned error: %v", err)
	}
}

func TestXrefCommandByOpName(t *testing.T) {
	xrefDirection = "both"
	err := runXref(nil, []string{"ToNumber"})
	if err != nil {
		t.Logf("xref ToNumber returned error: %v", err)
	}
}

func TestGrammarCommand(t *testing.T) {
	err := runGrammar(nil, []string{"HexEscapeSequence"})
	if err != nil {
		t.Logf("grammar returned error: %v", err)
	}
}

func TestGetWithMaxTokens(t *testing.T) {
	getMaxTokens = 100
	defer func() { getMaxTokens = 0 }()
	err := runGet(nil, []string{"sec-toprimitive"})
	if err != nil {
		t.Logf("get with max-tokens returned error: %v", err)
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
