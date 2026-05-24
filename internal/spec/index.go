package spec

// This file builds all query indexes: section numbers, operation names,
// cross-reference graphs, full-text search, and grammar productions.

import (
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

// Index builds all query indexes (sections, operations, xrefs, text, grammars).
func (s *Spec) Index() {
	s.SectionNums = make(map[string]*ClauseNode)
	s.OpsByName = make(map[string][]*ClauseNode)
	s.Outgoing = make(map[string]map[string]struct{})
	s.Incoming = make(map[string]map[string]struct{})
	s.TextIndex = make(map[string]map[string]int)
	s.ClauseTokens = make(map[string]int)
	s.Grammars = make(map[string][]GrammarProd)

	s.assignSectionNums()
	s.indexOperations()
	s.indexXRefs()
	s.indexText()
	s.indexGrammars()
}

// stopWords are common English words excluded from the full-text search index.
var stopWords = map[string]bool{
	"the": true, "a": true, "an": true, "is": true, "are": true,
	"was": true, "were": true, "be": true, "been": true, "being": true,
	"have": true, "has": true, "had": true, "do": true, "does": true,
	"did": true, "will": true, "would": true, "shall": true, "should": true,
	"may": true, "might": true, "must": true, "can": true, "could": true,
	"not": true, "and": true, "but": true, "or": true, "nor": true,
	"for": true, "with": true, "from": true, "into": true, "onto": true,
	"upon": true, "this": true, "that": true, "these": true, "those": true,
	"its": true, "then": true, "than": true, "also": true, "each": true,
	"all": true, "any": true, "some": true, "such": true, "both": true,
	"if": true,
}

// assignSectionNums assigns hierarchical section numbers to top-level clauses
// (numeric for clauses, letter-based for annexes) and recurses into children.
func (s *Spec) assignSectionNums() {
	var clauseNum, annexNum int
	for _, cn := range s.Nodes {
		// why: only number top-level clauses to avoid double-numbering;
		// children are numbered by assignChildSections when a top-level parent is found.
		if cn.Parent != nil {
			continue
		}
		switch cn.Kind {
		case "intro":
			cn.Section = ""
		case "clause":
			clauseNum++
			cn.Section = strconv.Itoa(clauseNum)
			s.SectionNums[cn.Section] = cn
			assignChildSections(cn, s.SectionNums)
		case "annex":
			annexNum++
			cn.Section = annexLetter(annexNum)
			s.SectionNums[cn.Section] = cn
			assignChildSections(cn, s.SectionNums)
		}
	}
}

// assignChildSections recursively assigns dotted section numbers to child clauses.
func assignChildSections(node *ClauseNode, sectionNums map[string]*ClauseNode) {
	if node.Section == "" {
		return
	}
	for i, child := range node.Children {
		child.Section = node.Section + "." + strconv.Itoa(i+1)
		sectionNums[child.Section] = child
		assignChildSections(child, sectionNums)
	}
}

// annexLetter converts a 1-based index to an uppercase annex letter (A-Z).
func annexLetter(n int) string {
	if n <= 26 {
		return string(rune('A' + n - 1))
	}
	return ""
}

// indexOperations builds the OpsByName index from clauses that have a type attribute.
func (s *Spec) indexOperations() {
	for _, cn := range s.Nodes {
		if cn.Type == "" {
			continue
		}
		name := extractOpName(cn.H1Text)
		if name == "" {
			continue
		}
		cn.Name = name
		key := normalizeOpName(name)
		s.OpsByName[key] = append(s.OpsByName[key], cn)
	}
}

// extractOpName derives an operation name from heading text by taking the portion
// before the first '(' character.
func extractOpName(h1Text string) string {
	text := strings.TrimSpace(h1Text)
	if idx := strings.IndexByte(text, '('); idx >= 0 {
		return strings.TrimSpace(text[:idx])
	}
	return text
}

// normalizeOpName lowercases the operation name and strips spaces for case-insensitive lookup.
func normalizeOpName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "")
	return name
}

// indexXRefs builds the Outgoing, Incoming, and XRefs indexes by walking
// every <emu-xref> element in every clause.
func (s *Spec) indexXRefs() {
	for _, cn := range s.Nodes {
		if cn.ID == "" {
			continue
		}
		walkXRefs(cn.HTMLNode, cn.ID, s, cn.HTMLNode)
	}
}

// walkXRefs recursively finds <emu-xref> elements, records outgoing/incoming links,
// and captures the surrounding heading context for each xref.
func walkXRefs(node *html.Node, sourceID string, spec *Spec, clauseRoot *html.Node) {
	if node.Type == html.ElementNode && node.Data == "emu-xref" {
		href := ""
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				href = attr.Val
				break
			}
		}
		if href != "" && strings.HasPrefix(href, "#") {
			targetID := href[1:]
			if spec.Outgoing[sourceID] == nil {
				spec.Outgoing[sourceID] = make(map[string]struct{})
			}
			spec.Outgoing[sourceID][targetID] = struct{}{}

			if spec.Incoming[targetID] == nil {
				spec.Incoming[targetID] = make(map[string]struct{})
			}
			spec.Incoming[targetID][sourceID] = struct{}{}

			context := findContext(node, clauseRoot)
			spec.XRefs = append(spec.XRefs, XRef{
				Source:  sourceID,
				Target:  targetID,
				Context: context,
			})
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		walkXRefs(c, sourceID, spec, clauseRoot)
	}
}

// findContext walks up from an xref node to find the nearest heading element
// within the same clause, for display context.
func findContext(node *html.Node, clauseRoot *html.Node) string {
	for n := node.Parent; n != nil && n != clauseRoot; n = n.Parent {
		if n.Type == html.ElementNode && isHeading(n.Data) {
			return textContent(n)
		}
	}
	return ""
}

// isHeading returns true if tag is an HTML heading element (h1-h6).
func isHeading(tag string) bool {
	switch tag {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return true
	}
	return false
}

// indexText builds the full-text search index (TextIndex) and per-clause token counts (ClauseTokens).
func (s *Spec) indexText() {
	for _, cn := range s.Nodes {
		if cn.ID == "" {
			continue
		}
		rawText := textContent(cn.HTMLNode)
		tokens := tokenize(rawText)
		var valid []string
		for _, tok := range tokens {
			// why: skip 1-2 char tokens (mostly noise: "it", "to", "of", "in", "on" etc.)
			// — they pollute the full-text index without adding meaningful search signal.
			if len(tok) < 3 {
				continue
			}
			if stopWords[tok] {
				continue
			}
			valid = append(valid, tok)
		}
		s.ClauseTokens[cn.ID] = len(valid)
		for _, tok := range valid {
			if s.TextIndex[tok] == nil {
				s.TextIndex[tok] = make(map[string]int)
			}
			s.TextIndex[tok][cn.ID]++
		}
	}
}

// tokenize splits text into lowercase alphanumeric tokens.
func tokenize(text string) []string {
	text = strings.ToLower(text)
	var tokens []string
	var buf strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			buf.WriteRune(r)
		} else {
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}
		}
	}
	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}
	return tokens
}

// indexGrammars finds all <emu-grammar> elements in every clause and indexes them by production name.
func (s *Spec) indexGrammars() {
	for _, cn := range s.Nodes {
		walkGrammar(cn.HTMLNode, cn.Section, s)
	}
}

// walkGrammar recursively finds <emu-grammar> elements and indexes them by production name.
func walkGrammar(node *html.Node, section string, spec *Spec) {
	if node.Type == html.ElementNode && node.Data == "emu-grammar" {
		rawText := textContent(node)
		name := extractGrammarName(rawText)
		if name != "" {
			gp := GrammarProd{
				Name:    name,
				Section: section,
				HTML:    rawText,
			}
			spec.Grammars[name] = append(spec.Grammars[name], gp)
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		walkGrammar(c, section, spec)
	}
}

// extractGrammarName returns the production name (text before "::") from grammar text.
func extractGrammarName(text string) string {
	text = strings.TrimSpace(text)
	if idx := strings.Index(text, "::"); idx >= 0 {
		return strings.TrimSpace(text[:idx])
	}
	return ""
}
