// Package spec parses the ECMAScript specification HTML into an in-memory clause tree.
package spec

import "golang.org/x/net/html"

// Spec holds the parsed ECMAScript specification as a clause tree.
type Spec struct {
	Root  *ClauseNode
	Nodes []*ClauseNode
	ByID  map[string]*ClauseNode

	SectionNums  map[string]*ClauseNode
	OpsByName    map[string][]*ClauseNode
	Outgoing     map[string]map[string]struct{}
	Incoming     map[string]map[string]struct{}
	XRefs        []XRef
	TextIndex    map[string]map[string]int
	ClauseTokens map[string]int
	Grammars     map[string][]GrammarProd
}

// ClauseNode represents a section of the spec (clause, annex, or intro).
type ClauseNode struct {
	ID       string
	Kind     string
	Type     string
	Section  string
	Name     string
	Depth    int
	HTMLNode *html.Node
	H1Text   string
	Parent   *ClauseNode
	Children []*ClauseNode
}

// GrammarProd represents a grammar production found in the spec.
type GrammarProd struct {
	Name    string
	Section string
	HTML    string
}

// XRef represents a cross-reference from one clause to another.
type XRef struct {
	Source  string
	Target  string
	Context string
}
