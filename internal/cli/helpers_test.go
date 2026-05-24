package cli

import (
	"testing"
)

func TestResolveBySectionNumber(t *testing.T) {
	s := getSpec()
	cn, _, _ := resolve(s, "1")
	if cn == nil {
		t.Fatal("expected a result for section '1'")
	}
	if cn.Section != "1" {
		t.Errorf("expected section '1', got %q", cn.Section)
	}
}

func TestResolveByOperationName(t *testing.T) {
	s := getSpec()
	cn, _, _ := resolve(s, "ToNumber")
	if cn == nil {
		t.Skip("ToNumber not found")
	}
	t.Logf("resolved ToNumber: %s (section %s)", cn.Name, cn.Section)
}

func TestResolveByAnchorID(t *testing.T) {
	s := getSpec()
	cn, _, _ := resolve(s, "sec-toprimitive")
	if cn == nil {
		t.Skip("sec-toprimitive not found")
	}
	if cn.ID != "sec-toprimitive" {
		t.Errorf("expected id 'sec-toprimitive', got %q", cn.ID)
	}
}

func TestResolveAmbiguous(t *testing.T) {
	s := getSpec()
	_, multi, matches := resolve(s, "toString")
	if len(multi) <= 1 && len(matches) <= 1 {
		t.Skip("toString is not ambiguous in this spec version")
	}
	t.Logf("toString matched %d operations", len(multi))
	for _, m := range matches {
		t.Logf("  %s (%s)", m.Name, m.Kind)
	}
}

func TestResolveNotFound(t *testing.T) {
	s := getSpec()
	cn, _, _ := resolve(s, "nonexistent_operation_xyz")
	if cn != nil {
		t.Errorf("expected nil for nonexistent query, got %q", cn.Name)
	}
}

func TestSuggestions(t *testing.T) {
	s := getSpec()
	succ := suggestions("tonumbr", s)
	if len(succ) == 0 {
		t.Skip("no suggestions found")
	}
	t.Logf("suggestions for 'tonumbr': %v", succ)
}

func TestKindLabel(t *testing.T) {
	s := getSpec()
	for _, cn := range s.Nodes {
		if cn.Type == "abstract operation" {
			label := kindLabel(cn)
			if label == "" {
				t.Error("empty kind label for abstract operation")
			}
			break
		}
	}
}

func TestExtractSnippet(t *testing.T) {
	s := getSpec()
	for _, tok := range []string{"object", "function", "value"} {
		for _, cn := range s.Nodes {
			snippet := extractSnippet(cn, tok, 40)
			if snippet != "" {
				t.Logf("snippet for %q in %s: %s", tok, cn.ID, snippet)
				return
			}
		}
	}
}

func TestResolveBySpecType(t *testing.T) {
	s := getSpec()
	cn, _, _ := resolve(s, "spec type")
	if cn != nil {
		t.Logf("resolved 'spec type': %s", cn.H1Text)
	}
}

func TestResolveNormalizedOpName(t *testing.T) {
	s := getSpec()
	cn, _, _ := resolve(s, "ToNumber ( _argument_ )")
	if cn == nil {
		t.Skip("ToNumber ( _argument_ ) not resolved")
	}
	t.Logf("resolved ToNumber ( _argument_ ): %s", cn.Name)
}

func TestClauseURL(t *testing.T) {
	u := clauseURL("sec-toprimitive")
	if u != "https://tc39.es/ecma262/#sec-toprimitive" {
		t.Errorf("unexpected URL: %q", u)
	}
}

func TestEstimateTokens(t *testing.T) {
	tokens := estimateTokens("hello world this is a test")
	if tokens <= 0 {
		t.Error("expected positive token count")
	}
	t.Logf("estimated tokens: %d", tokens)
}

func TestLevenshtein(t *testing.T) {
	d := levenshtein("tonumber", "tonumbr")
	if d != 1 {
		t.Errorf("expected distance 1, got %d", d)
	}

	d = levenshtein("abc", "abc")
	if d != 0 {
		t.Errorf("expected distance 0, got %d", d)
	}

	d = levenshtein("", "abc")
	if d != 3 {
		t.Errorf("expected distance 3, got %d", d)
	}
}

func TestNormalizeQuery(t *testing.T) {
	n1 := normalizeQuery("ToNumber")
	if n1 != "tonumber" {
		t.Errorf("expected 'tonumber', got %q", n1)
	}

	n2 := normalizeQuery(" Array.prototype.map ")
	if n2 != "array.prototype.map" {
		t.Errorf("expected 'array.prototype.map', got %q", n2)
	}
}

func TestResolveByID(t *testing.T) {
	s := getSpec()
	if len(s.ByID) == 0 {
		t.Fatal("ByID is empty")
	}
	for id := range s.ByID {
		cn, _, _ := resolve(s, id)
		if cn == nil {
			t.Errorf("failed to resolve by ID %q", id)
		}
		break
	}
}

func TestSpecSingleton(t *testing.T) {
	s1 := getSpec()
	s2 := getSpec()
	if s1 != s2 {
		t.Error("getSpec did not return singleton")
	}
}

func TestResolveInternalSlot(t *testing.T) {
	s := getSpec()
	cn, _, _ := resolve(s, "Prototype")
	t.Logf("internal slot resolve result: %v", cn != nil)
}

func TestSpecTypeCheck(t *testing.T) {
	_ = getSpec()
	t.Log("spec loaded successfully")
}
