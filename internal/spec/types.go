// Package spec parses the ECMAScript specification HTML into an in-memory clause tree.
package spec

import "golang.org/x/net/html"

// Spec holds the parsed ECMAScript specification as a clause tree with all indexes.
type Spec struct {
	Root  *ClauseNode            // Top-level root clause node.
	Nodes []*ClauseNode          // Flat list of all clause nodes in document order.
	ByID  map[string]*ClauseNode // Lookup by HTML id attribute (e.g. "sec-...").

	SectionNums  map[string]*ClauseNode         // Lookup by section number string (e.g. "1.2.3" or "A").
	OpsByName    map[string][]*ClauseNode       // Lookup by normalized operation name.
	Outgoing     map[string]map[string]struct{} // sourceID → set of targetIDs (outgoing xrefs).
	Incoming     map[string]map[string]struct{} // targetID → set of sourceIDs (incoming xrefs).
	XRefs        []XRef                         // Flat list of all cross-references with context.
	TextIndex    map[string]map[string]int      // token → clauseID → count (full-text search index).
	ClauseTokens map[string]int                 // clauseID → total indexed token count.
	Grammars     map[string][]GrammarProd       // Grammar production name → list of productions.
}

// ClauseNode represents a section of the spec (clause, annex, or intro).
type ClauseNode struct {
	ID       string        // HTML id attribute (e.g. "sec-intro", "sec-ecmascript-data-types-and-values").
	Kind     string        // One of "clause", "annex", or "intro".
	Type     string        // The "type" attribute on <emu-clause>, used for operation lookup.
	Section  string        // Assigned section number (e.g. "1", "1.2.3", "A", "B.1").
	Name     string        // Extracted operation name (from H1Text when Type is set).
	Depth    int           // Nesting depth, 0 for top-level clauses.
	HTMLNode *html.Node    // Underlying DOM node (<emu-clause>, <emu-annex>, or <emu-intro>).
	H1Text   string        // The heading text (from the first <h1> or <h2> child).
	Parent   *ClauseNode   // Parent clause in the tree; nil for root-level clauses.
	Children []*ClauseNode // Nested child clauses.
}

// GrammarProd represents a grammar production found in <emu-grammar> elements.
type GrammarProd struct {
	Name    string // Production name (text before "::").
	Section string // Section number where this production was found.
	HTML    string // Raw text content of the <emu-grammar> element.
}

// XRef represents a cross-reference from one clause to another via <emu-xref href="...">.
type XRef struct {
	Source  string // ID of the clause containing the cross-reference.
	Target  string // ID of the referenced clause.
	Context string // Nearest heading text near the xref, for display context.
}
