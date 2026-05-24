package spec

import "testing"

func TestLoad(t *testing.T) {
	s, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if s.Root == nil {
		t.Fatal("Root is nil")
	}
	if s.ClauseCount() == 0 {
		t.Fatal("zero clauses parsed")
	}
	if len(s.ByID) == 0 {
		t.Fatal("ByID map is empty")
	}

	t.Logf("Parsed %d clauses, %d IDs", s.ClauseCount(), len(s.ByID))
}

func TestParse_empty(t *testing.T) {
	s, err := Parse([]byte("<html></html>"))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if s.ClauseCount() != 0 {
		t.Fatalf("expected 0 clauses, got %d", s.ClauseCount())
	}
}

func TestParse_singleClause(t *testing.T) {
	html := `<html><body><emu-clause id="sec-test"><h1>Test Clause</h1></emu-clause></body></html>`
	s, err := Parse([]byte(html))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if s.ClauseCount() != 1 {
		t.Fatalf("expected 1 clause, got %d", s.ClauseCount())
	}
	cn := s.Nodes[0]
	if cn.ID != "sec-test" {
		t.Errorf("expected ID 'sec-test', got %q", cn.ID)
	}
	if cn.H1Text != "Test Clause" {
		t.Errorf("expected H1 'Test Clause', got %q", cn.H1Text)
	}
	if cn.Kind != "clause" {
		t.Errorf("expected kind 'clause', got %q", cn.Kind)
	}
	if cn.Depth != 0 {
		t.Errorf("expected depth 0, got %d", cn.Depth)
	}
	if c, ok := s.ByID["sec-test"]; !ok {
		t.Error("ByID key 'sec-test' missing")
	} else if c != cn {
		t.Error("ByID returned wrong pointer")
	}
}

func TestParse_nested(t *testing.T) {
	html := `<html><body><emu-clause id="sec-a"><h1>A</h1><emu-clause id="sec-b"><h1>B</h1></emu-clause></emu-clause></body></html>`
	s, err := Parse([]byte(html))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if s.ClauseCount() != 2 {
		t.Fatalf("expected 2 clauses, got %d", s.ClauseCount())
	}
	if s.Nodes[0].Children[0] != s.Nodes[1] {
		t.Error("nested parent/child relationship wrong")
	}
	if s.Nodes[1].Depth != 1 {
		t.Errorf("expected nested depth 1, got %d", s.Nodes[1].Depth)
	}
}

func TestParse_annexAndIntro(t *testing.T) {
	html := `<html><body><emu-intro id="sec-intro"><h1>Intro</h1></emu-intro><emu-annex id="sec-annex"><h1>Annex</h1></emu-annex></body></html>`
	s, err := Parse([]byte(html))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if s.ClauseCount() != 2 {
		t.Fatalf("expected 2 clauses, got %d", s.ClauseCount())
	}
	if s.Nodes[0].Kind != "intro" {
		t.Errorf("expected kind 'intro', got %q", s.Nodes[0].Kind)
	}
	if s.Nodes[1].Kind != "annex" {
		t.Errorf("expected kind 'annex', got %q", s.Nodes[1].Kind)
	}
}

func TestParse_typeAttribute(t *testing.T) {
	html := `<html><body><emu-clause id="sec-ao" type="abstract operation"><h1>ToNumber ( x )</h1></emu-clause></body></html>`
	s, err := Parse([]byte(html))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	if s.Nodes[0].Type != "abstract operation" {
		t.Errorf("expected type 'abstract operation', got %q", s.Nodes[0].Type)
	}
}
