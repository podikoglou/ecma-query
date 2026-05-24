// Package spec parses the ECMAScript specification HTML into an in-memory clause tree.
package spec

import "golang.org/x/net/html"

// Spec holds the parsed ECMAScript specification as a clause tree.
type Spec struct {
	Root  *ClauseNode
	Nodes []*ClauseNode
	ByID  map[string]*ClauseNode
}

// ClauseNode represents a section of the spec (clause, annex, or intro).
type ClauseNode struct {
	ID       string
	Kind     string
	Type     string
	Depth    int
	HTMLNode *html.Node
	H1Text   string
	Parent   *ClauseNode
	Children []*ClauseNode
}
