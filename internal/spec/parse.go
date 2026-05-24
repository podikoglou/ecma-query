package spec

// This file walks the spec DOM tree and builds a clause tree from <emu-clause> elements.

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// Load parses the embedded spec.html into a Spec tree and builds all indexes.
func Load() (*Spec, error) {
	s, err := Parse(specHTML)
	if err != nil {
		return nil, err
	}
	s.Index()
	return s, nil
}

// Parse parses HTML bytes into a Spec clause tree.
func Parse(data []byte) (*Spec, error) {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("parse spec HTML: %w", err)
	}

	s := &Spec{
		ByID: make(map[string]*ClauseNode),
	}

	s.Root = walk(doc, 0, nil, s)
	return s, nil
}

// walk recursively traverses the HTML tree, building ClauseNodes from emu-* elements.
// It returns the first child clause found (as a convenience for the caller).
func walk(node *html.Node, depth int, parent *ClauseNode, spec *Spec) *ClauseNode {
	var childClause *ClauseNode

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}

		kind := clauseKind(c.Data)
		if kind == "" {
			// why: the spec HTML wraps emu-* elements in intermediate <div> etc.;
			// we must recurse into them to find nested clause/annex/intro elements.
			sub := walk(c, depth, parent, spec)
			if sub != nil && childClause == nil {
				childClause = sub
			}
			continue
		}

		cn := buildClauseNode(c, kind, depth, parent)
		spec.Nodes = append(spec.Nodes, cn)
		if cn.ID != "" {
			spec.ByID[cn.ID] = cn
		}

		if parent != nil {
			parent.Children = append(parent.Children, cn)
		}

		walk(c, depth+1, cn, spec)

		if childClause == nil {
			childClause = cn
		}
	}

	return childClause
}

// clauseKind maps HTML tag names to Spec clause kinds.
// Returns "" for non-clause elements.
func clauseKind(tag string) string {
	switch tag {
	case "emu-clause":
		return "clause"
	case "emu-annex":
		return "annex"
	case "emu-intro":
		return "intro"
	default:
		return ""
	}
}

// buildClauseNode creates a ClauseNode from an HTML element node,
// extracting id and type attributes and the heading text.
func buildClauseNode(node *html.Node, kind string, depth int, parent *ClauseNode) *ClauseNode {
	cn := &ClauseNode{
		Kind:     kind,
		Depth:    depth,
		HTMLNode: node,
		Parent:   parent,
	}

	for _, attr := range node.Attr {
		switch attr.Key {
		case "id":
			cn.ID = attr.Val
		case "type":
			cn.Type = attr.Val
		}
	}

	cn.H1Text = extractHeading(cn)
	return cn
}

// extractHeading finds the heading text from a clause's first <h1> or <h2> child.
func extractHeading(cn *ClauseNode) string {
	for c := cn.HTMLNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		if c.Data == "h1" {
			return textContent(c)
		}
	}

	for c := cn.HTMLNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		if c.Data == "h2" {
			return textContent(c)
		}
	}

	return ""
}

// textContent extracts the concatenated text content of an HTML node tree,
// skipping all tags and returning trimmed whitespace.
func textContent(node *html.Node) string {
	var buf strings.Builder
	var collect func(*html.Node)
	collect = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			collect(c)
		}
	}
	collect(node)
	return strings.TrimSpace(buf.String())
}

// ClauseCount returns the total number of parsed clause nodes.
func (s *Spec) ClauseCount() int {
	return len(s.Nodes)
}
