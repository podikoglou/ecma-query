package spec

import (
	"testing"
)

func TestSectionNumbers(t *testing.T) {
	s, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(s.SectionNums) == 0 {
		t.Fatal("SectionNums is empty")
	}
	for _, cn := range s.Nodes {
		if cn.Parent == nil && cn.Kind == "clause" {
			if cn.Section == "" {
				t.Errorf("top-level clause %q has empty section", cn.ID)
			}
		}
	}
	if _, ok := s.SectionNums["1"]; !ok {
		t.Error("section '1' not found in SectionNums")
	}
	if cn, ok := s.SectionNums["1"]; ok {
		if cn.Kind != "clause" {
			t.Errorf("section '1' has wrong kind: %q", cn.Kind)
		}
	}
	t.Logf("section count: %d", len(s.SectionNums))

	childFound := false
	for _, cn := range s.Nodes {
		for _, child := range cn.Children {
			if child.Section != "" {
				childFound = true
				if child.Section == cn.Section+"." {
					t.Logf("child %q section %q from parent %q", child.ID, child.Section, cn.ID)
				}
				break
			}
		}
		if childFound {
			break
		}
	}
	if !childFound && len(s.Nodes) > 1 {
		t.Log("no children found with section numbers (might be flat)")
	}
}

func TestOperationNames(t *testing.T) {
	s, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	found := false
	for _, cn := range s.Nodes {
		if cn.Type == "abstract operation" {
			found = true
			if cn.Name == "" {
				t.Error("abstract operation has empty name")
			} else {
				t.Logf("abstract operation: %q → %q", cn.H1Text, cn.Name)
			}
			break
		}
	}
	if !found {
		t.Skip("no abstract operations in spec")
	}

	key := normalizeOpName("ToNumber")
	if nodes, ok := s.OpsByName[key]; ok {
		t.Logf("ToNumber found (%d nodes)", len(nodes))
		if len(nodes) > 0 && nodes[0].Name != "ToNumber" {
			t.Errorf("expected Name 'ToNumber', got %q", nodes[0].Name)
		}
	} else {
		t.Error("ToNumber not found in OpsByName")
	}
}

func TestCrossReferences(t *testing.T) {
	s, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(s.Outgoing) == 0 {
		t.Error("Outgoing is empty")
	}
	if len(s.Incoming) == 0 {
		t.Error("Incoming is empty")
	}
	if len(s.XRefs) == 0 {
		t.Error("XRefs is empty")
	}

	t.Logf("outgoing edges: %d, incoming edges: %d, total xrefs: %d",
		len(s.Outgoing), len(s.Incoming), len(s.XRefs))

	count := 0
	for _, targets := range s.Outgoing {
		for target := range targets {
			if _, ok := s.Incoming[target]; !ok {
				count++
				if count <= 5 {
					t.Logf("target %q has no inbound entry", target)
				}
			}
		}
	}
	if count > 0 {
		t.Logf("%d outgoing targets have no incoming entry (may be undefined IDs)", count)
	}

	for _, xr := range s.XRefs[:min(3, len(s.XRefs))] {
		if xr.Source == "" || xr.Target == "" {
			t.Error("XRef has empty Source or Target")
		}
	}
}

func TestFullTextIndex(t *testing.T) {
	s, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(s.TextIndex) == 0 {
		t.Error("TextIndex is empty")
	}
	if len(s.ClauseTokens) == 0 {
		t.Error("ClauseTokens is empty")
	}

	t.Logf("text index has %d unique tokens across %d clauses",
		len(s.TextIndex), len(s.ClauseTokens))

	totalTokens := 0
	for _, count := range s.ClauseTokens {
		totalTokens += count
	}
	if totalTokens == 0 {
		t.Error("zero total tokens")
	}
	t.Logf("total tokens: %d", totalTokens)

	for tok, clauseMap := range s.TextIndex {
		if len(clauseMap) > 50 {
			t.Logf("common token %q appears in %d clauses", tok, len(clauseMap))
			break
		}
	}
}

func TestGrammarIndex(t *testing.T) {
	s, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if len(s.Grammars) == 0 {
		t.Error("Grammars is empty")
	}

	total := 0
	for name, prods := range s.Grammars {
		total += len(prods)
		if len(prods) > 3 {
			t.Logf("production %q appears %d times", name, len(prods))
		}
	}
	t.Logf("total grammar productions: %d, unique names: %d", total, len(s.Grammars))

	count := 0
	for _, prods := range s.Grammars {
		for _, gp := range prods {
			if gp.Name == "" {
				t.Error("GrammarProd has empty Name")
			}
			if gp.HTML == "" {
				t.Error("GrammarProd has empty HTML")
			}
			count++
			if count >= 3 {
				return
			}
		}
	}
}

func TestIndex_empty(t *testing.T) {
	s := &Spec{
		ByID: make(map[string]*ClauseNode),
	}
	s.Index()

	if len(s.SectionNums) != 0 {
		t.Error("empty spec should have zero section nums")
	}
	if len(s.OpsByName) != 0 {
		t.Error("empty spec should have zero ops")
	}
	if len(s.Outgoing) != 0 {
		t.Error("empty spec should have zero outgoing xrefs")
	}
	if len(s.TextIndex) != 0 {
		t.Error("empty spec should have zero text tokens")
	}
	if len(s.Grammars) != 0 {
		t.Error("empty spec should have zero grammars")
	}
}
