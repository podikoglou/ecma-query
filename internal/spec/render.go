package spec

// This file converts spec HTML into Markdown, supporting all emu-* custom elements.

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

// mdWriter accumulates Markdown output while tracking rendering state.
type mdWriter struct {
	buf       strings.Builder // Output buffer.
	spec      *Spec           // Spec reference for resolving cross-references.
	inAlg     bool            // Whether we are inside an <emu-alg> block.
	algDepth  int             // Current nesting depth for algorithm ordered lists.
	olCounter []int           // Per-depth counters for algorithm step numbering.
	isBlock   bool            // Whether the previous output was a block element (controls spacing).
}

// RenderNode converts a clause tree to Markdown with breadcrumbs and cross-references.
func (s *Spec) RenderNode(node *ClauseNode) (string, error) {
	w := &mdWriter{spec: s}
	w.renderNode(node)
	return strings.TrimSpace(w.buf.String()), nil
}

// RenderHTML converts an HTML subtree to Markdown (no breadcrumbs or xrefs).
func RenderHTML(root *html.Node) string {
	w := &mdWriter{}
	w.walk(root)
	return strings.TrimSpace(w.buf.String())
}

// renderNode renders a clause node with heading, breadcrumb trail, body content, and "See also" xrefs.
func (w *mdWriter) renderNode(node *ClauseNode) {
	heading := node.H1Text
	if heading != "" {
		w.buf.WriteString("## ")
		w.buf.WriteString(heading)
		w.buf.WriteString("\n\n")
		w.isBlock = true
	}

	crumbs := w.buildBreadcrumb(node)
	if crumbs != "" {
		w.buf.WriteString("> ")
		w.buf.WriteString(crumbs)
		w.buf.WriteString("\n\n")
		w.isBlock = true
	}

	if node.HTMLNode != nil {
		for c := node.HTMLNode.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && (c.Data == "h1" || c.Data == "h2") {
				continue
			}
			w.walk(c)
		}
	}

	if node.ID != "" {
		if targets, ok := w.spec.Outgoing[node.ID]; ok && len(targets) > 0 {
			w.ensureNewline()
			w.buf.WriteString("See also:")
			first := true
			for targetID := range targets {
				if targetNode, ok := w.spec.ByID[targetID]; ok {
					if !first {
						w.buf.WriteString(",")
					}
					first = false
					name := targetNode.Name
					if name == "" {
						name = targetNode.H1Text
					}
					if name == "" {
						name = targetID
					}
					w.buf.WriteString(" [")
					w.buf.WriteString(name)
					w.buf.WriteString("](#")
					w.buf.WriteString(targetID)
					w.buf.WriteString(")")
				}
			}
			w.buf.WriteString("\n")
		}
	}
}

// buildBreadcrumb builds a "Section → parent → child" breadcrumb string from the clause ancestry.
func (w *mdWriter) buildBreadcrumb(node *ClauseNode) string {
	if node == nil {
		return ""
	}
	var parts []string

	for p := node.Parent; p != nil; p = p.Parent {
		text := p.H1Text
		if text == "" {
			text = p.Section
		}
		if p.Section != "" {
			text = p.Section + " " + text
		}
		if text != "" {
			parts = append([]string{text}, parts...)
		}
	}

	if node.Section != "" {
		parts = append(parts, node.Section)
	} else if node.H1Text != "" {
		parts = append(parts, node.H1Text)
	}

	return strings.Join(parts, " → ")
}

// walk dispatches rendering based on the node type (text, element, or document).
func (w *mdWriter) walk(n *html.Node) {
	switch n.Type {
	case html.TextNode:
		w.emitText(n.Data)
	case html.ElementNode:
		w.walkElement(n)
	case html.DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			w.walk(c)
		}
	}
}

// walkElement routes an HTML element to the appropriate render method based on its tag.
func (w *mdWriter) walkElement(n *html.Node) {
	switch n.Data {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		w.renderHeading(n)
	case "p":
		w.renderParagraph(n)
	case "b", "strong":
		w.renderInline(n, "**", "**")
	case "em", "i", "dfn":
		w.renderInline(n, "*", "*")
	case "code":
		w.renderInline(n, "`", "`")
	case "var":
		w.renderInline(n, "_", "_")
	case "sub", "sup":
		// why: subscript and superscript are rendered as plain text
		// since Markdown has no native support for them.
		w.walkChildren(n)
	case "a":
		w.renderLink(n)
	case "ul":
		w.renderUnorderedList(n)
	case "ol":
		w.renderOrderedList(n)
	case "dl":
		w.renderDefList(n)
	case "pre":
		w.renderPre(n)
	case "br":
		w.buf.WriteString("  \n")
	case "emu-alg":
		w.renderEmuAlg(n)
	case "emu-note":
		w.renderEmuNote(n)
	case "emu-xref":
		w.renderEmuXRef(n)
	case "emu-grammar":
		w.renderEmuGrammar(n)
	case "emu-table":
		w.renderEmuTable(n)
	case "emu-eqn":
		w.renderInline(n, "`", "`")
	case "emu-val":
		w.walkChildren(n)
	case "emu-figure":
		w.renderEmuFigure(n)
	case "emu-prodref":
		w.renderEmuProdRef(n)
	case "emu-concrete-method-dfns":
		w.walkChildren(n)
	case "emu-meta", "emu-import":
		// why: metadata and import elements are internal to the spec HTML structure
		// and should not produce visible output.
	case "figure", "img":
		// why: figures are often non-essential diagrams; images can't be rendered in text.
		// <emu-figure> is handled separately via renderEmuFigure.
	case "span", "div", "ins":
		w.walkChildren(n)
	case "style", "link", "meta":
		// why: embedded stylesheets and metadata are not part of the rendered content.
	default:
		w.walkChildren(n)
	}
}

// walkChildren renders all child nodes in order.
func (w *mdWriter) walkChildren(n *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		w.walk(c)
	}
}

// renderHeading renders h1-h6 as Markdown ATX headings at the corresponding level.
func (w *mdWriter) renderHeading(n *html.Node) {
	w.ensureNewline()
	level := n.Data[1] - '0'
	w.buf.WriteString(strings.Repeat("#", int(level)))
	w.buf.WriteString(" ")
	w.walkChildren(n)
	w.buf.WriteString("\n\n")
	w.isBlock = true
}

// renderParagraph renders a <p> element as a Markdown paragraph.
func (w *mdWriter) renderParagraph(n *html.Node) {
	w.ensureNewline()
	w.walkChildren(n)
	w.buf.WriteString("\n\n")
	w.isBlock = true
}

// renderInline wraps child content with prefix/suffix (used for bold, italic, code, etc.).
func (w *mdWriter) renderInline(n *html.Node, prefix, suffix string) {
	w.buf.WriteString(prefix)
	w.walkChildren(n)
	w.buf.WriteString(suffix)
}

// renderLink renders an <a> element as a Markdown link [text](href).
func (w *mdWriter) renderLink(n *html.Node) {
	href := ""
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}
	text := renderInlineText(n)
	w.buf.WriteString("[")
	w.buf.WriteString(text)
	w.buf.WriteString("](")
	w.buf.WriteString(href)
	w.buf.WriteString(")")
}

// renderUnorderedList renders <ul> as a Markdown bullet list.
func (w *mdWriter) renderUnorderedList(n *html.Node) {
	w.ensureNewline()
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "li" {
			w.buf.WriteString("- ")
			w.buf.WriteString(renderInlineText(c))
			w.buf.WriteString("\n")
		}
	}
	w.buf.WriteString("\n")
	w.isBlock = true
}

// renderOrderedList renders <ol> as a Markdown numbered list.
// When inside an <emu-alg> block, it delegates to renderAlgOrderedList for indented step numbering.
func (w *mdWriter) renderOrderedList(n *html.Node) {
	if w.inAlg {
		depth := w.algDepth
		w.algDepth++
		for depth >= len(w.olCounter) {
			w.olCounter = append(w.olCounter, 0)
		}
		w.renderAlgOrderedList(n, depth)
		w.algDepth--
	} else {
		w.ensureNewline()
		counter := 0
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "li" {
				counter++
				fmt.Fprintf(&w.buf, "%d. ", counter)
				w.buf.WriteString(renderInlineText(c))
				w.buf.WriteString("\n")
			}
		}
		w.buf.WriteString("\n")
		w.isBlock = true
	}
}

// renderAlgOrderedList renders nested <ol> steps within an <emu-alg> block with indentation and step counters.
func (w *mdWriter) renderAlgOrderedList(n *html.Node, depth int) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode || c.Data != "li" {
			continue
		}
		w.olCounter[depth]++

		indent := strings.Repeat("  ", depth)
		w.buf.WriteString(indent)

		if depth == 0 {
			fmt.Fprintf(&w.buf, "%d. ", w.olCounter[0])
		} else {
			w.buf.WriteString("1. ")
		}

		content := w.renderAlgLI(c, depth+1)
		w.buf.WriteString(content)
		w.buf.WriteString("\n")
	}
	w.buf.WriteString("\n")
	w.isBlock = true
}

// renderAlgLI renders a single <li> within an algorithm, handling nested <ol> children recursively.
func (w *mdWriter) renderAlgLI(li *html.Node, childDepth int) string {
	var parts []string
	for c := li.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "ol" {
			subWriter := &mdWriter{spec: w.spec, inAlg: true}
			subWriter.algDepth = childDepth
			subWriter.olCounter = append([]int{}, w.olCounter...)
			subWriter.renderOrderedList(c)
			nested := strings.TrimRight(subWriter.buf.String(), "\n")
			if len(parts) > 0 {
				parts = append(parts, "\n")
			}
			parts = append(parts, nested)
		} else {
			parts = append(parts, w.renderInlineNode(c))
		}
	}
	result := strings.TrimSpace(strings.Join(parts, ""))
	result = stripAlgPrefix(result)
	return result
}

// stripAlgPrefix removes numeric algorithm step prefixes like "1. ", "2. " etc.
// from the beginning of a string. This is needed because the step number is
// already rendered from the list counter.
func stripAlgPrefix(s string) string {
	s = strings.TrimSpace(s)
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			continue
		}
		if s[i] == '.' && i+1 < len(s) && s[i+1] == ' ' {
			return s[i+2:]
		}
		break
	}
	return s
}

// renderInlineNode renders an HTML node as inline text, returning the result as a string.
func (w *mdWriter) renderInlineNode(n *html.Node) string {
	childWriter := &mdWriter{spec: w.spec}
	childWriter.walk(n)
	return childWriter.buf.String()
}

// renderDefList renders <dl> as bold term: definition pairs.
func (w *mdWriter) renderDefList(n *html.Node) {
	w.ensureNewline()
	var dt, dd string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		switch c.Data {
		case "dt":
			dt = renderInlineText(c)
		case "dd":
			dd = renderInlineText(c)
			if dt != "" && dd != "" {
				w.buf.WriteString("**")
				w.buf.WriteString(dt)
				w.buf.WriteString("**: ")
				w.buf.WriteString(dd)
				w.buf.WriteString("\n")
				dt = ""
			}
		}
	}
	w.buf.WriteString("\n")
	w.isBlock = true
}

// renderPre renders <pre> as a Markdown fenced code block.
func (w *mdWriter) renderPre(n *html.Node) {
	w.ensureNewline()
	w.buf.WriteString("```\n")
	w.buf.WriteString(strings.TrimRight(textContentRaw(n), "\n"))
	w.buf.WriteString("\n```\n\n")
	w.isBlock = true
}

// renderEmuAlg renders an <emu-alg> block as indented numbered steps.
func (w *mdWriter) renderEmuAlg(n *html.Node) {
	w.ensureNewline()
	lines := w.collectAlgLines(n)
	w.renderParsedAlg(lines)
	w.buf.WriteString("\n")
	w.isBlock = true
}

// collectAlgLines extracts text lines from an <emu-alg> block, preserving inline formatting.
func (w *mdWriter) collectAlgLines(n *html.Node) []string {
	var buf strings.Builder
	w.collectAlgText(n, &buf)
	raw := buf.String()

	var lines []string
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimRight(line, " \t")
		if strings.TrimSpace(trimmed) == "" {
			continue
		}
		lines = append(lines, trimmed)
	}
	return lines
}

// collectAlgText recursively collects text from an <emu-alg> subtree, converting
// inline emu-* elements to Markdown equivalents in the output.
func (w *mdWriter) collectAlgText(n *html.Node, buf *strings.Builder) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			buf.WriteString(c.Data)
		} else if c.Type == html.ElementNode {
			switch c.Data {
			case "emu-xref":
				buf.WriteString("[")
				buf.WriteString(renderInlineText(c))
				buf.WriteString("]")
				for _, attr := range c.Attr {
					if attr.Key == "href" {
						buf.WriteString("(")
						buf.WriteString(attr.Val)
						buf.WriteString(")")
						break
					}
				}
			case "emu-val":
				buf.WriteString(renderInlineText(c))
			case "emu-eqn":
				buf.WriteString("`")
				buf.WriteString(renderInlineText(c))
				buf.WriteString("`")
			case "b", "strong":
				buf.WriteString("**")
				buf.WriteString(renderInlineText(c))
				buf.WriteString("**")
			case "em", "i", "dfn":
				buf.WriteString("*")
				buf.WriteString(renderInlineText(c))
				buf.WriteString("*")
			case "var":
				buf.WriteString("_")
				buf.WriteString(renderInlineText(c))
				buf.WriteString("_")
			case "code":
				buf.WriteString("`")
				buf.WriteString(renderInlineText(c))
				buf.WriteString("`")
			case "sub", "sup":
				buf.WriteString(renderInlineText(c))
			case "emu-grammar":
				buf.WriteString("```\n")
				buf.WriteString(strings.TrimSpace(textContentRaw(c)))
				buf.WriteString("\n```")
			default:
				buf.WriteString(renderInlineText(c))
			}
		}
	}
}

// renderParsedAlg renders algorithm lines with indentation-based nesting and auto-numbering.
func (w *mdWriter) renderParsedAlg(lines []string) {
	if len(lines) == 0 {
		return
	}

	// why: the spec uses indentation (not <ol>) for algorithm nesting in some places.
	// We compute the minimum indent to establish a baseline, then convert relative
	// indentation levels into numbered steps with increasing counters.
	minIndent := -1
	for _, line := range lines {
		n := 0
		for n < len(line) && line[n] == ' ' {
			n++
		}
		if n < len(line) && (minIndent == -1 || n < minIndent) {
			minIndent = n
		}
	}
	if minIndent < 0 {
		minIndent = 0
	}

	type algLine struct {
		level int
		text  string
	}
	var parsed []algLine
	for _, line := range lines {
		relIndent := 0
		if len(line) > minIndent {
			relIndent = countIndent(line[minIndent:])
		}
		content := strings.TrimLeft(line[minIndent:], " ")
		content = stripAlgPrefix(content)
		parsed = append(parsed, algLine{level: relIndent, text: content})
	}

	counters := make([]int, 8)
	for _, pl := range parsed {
		if pl.level >= len(counters) {
			continue
		}
		// why: reset deeper counters when indentation level decreases,
		// so each nested group starts at 1.
		for i := pl.level + 1; i < len(counters); i++ {
			counters[i] = 0
		}
		counters[pl.level]++

		indent := strings.Repeat("  ", pl.level)
		if pl.level == 0 {
			fmt.Fprintf(&w.buf, "%s%d. %s\n", indent, counters[0], pl.text)
		} else {
			fmt.Fprintf(&w.buf, "%s1. %s\n", indent, pl.text)
		}
	}
}

// countIndent returns the number of 2-space indentation levels at the start of s.
func countIndent(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			n++
		} else {
			break
		}
	}
	return n / 2
}

// renderEmuNote renders <emu-note> as a Markdown blockquote with bold "Note:" prefix.
func (w *mdWriter) renderEmuNote(n *html.Node) {
	w.ensureNewline()
	span := findSpanWithClass(n, "note")
	if span != nil {
		w.buf.WriteString("> **")
		w.buf.WriteString(renderInlineText(span))
		w.buf.WriteString("**")
	}

	content := renderInlineText(n)
	if span != nil {
		// why: the <span class="note"> text is rendered as "**Note:**" above,
		// so we strip it from the body content to avoid duplication.
		spanText := renderInlineText(span)
		content = strings.TrimPrefix(content, spanText)
	}

	if strings.TrimSpace(content) != "" {
		w.buf.WriteString(content)
	}
	w.buf.WriteString("\n\n")
	w.isBlock = true
}

// findSpanWithClass finds a direct <span> child with the given class attribute.
func findSpanWithClass(n *html.Node, class string) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "span" {
			for _, attr := range c.Attr {
				if attr.Key == "class" && attr.Val == class {
					return c
				}
			}
		}
	}
	return nil
}

// renderEmuXRef renders an <emu-xref> as a Markdown link [text](href).
func (w *mdWriter) renderEmuXRef(n *html.Node) {
	href := ""
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}
	text := renderInlineText(n)
	w.buf.WriteString("[")
	w.buf.WriteString(text)
	w.buf.WriteString("](")
	w.buf.WriteString(href)
	w.buf.WriteString(")")
}

// renderEmuGrammar renders <emu-grammar> as a Markdown fenced code block.
func (w *mdWriter) renderEmuGrammar(n *html.Node) {
	w.ensureNewline()
	w.buf.WriteString("```\n")
	w.buf.WriteString(strings.TrimSpace(textContentRaw(n)))
	w.buf.WriteString("\n```\n\n")
	w.isBlock = true
}

// renderEmuTable renders <emu-table> as a Markdown table with caption and pipe-formatted rows.
func (w *mdWriter) renderEmuTable(n *html.Node) {
	w.ensureNewline()

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == "emu-caption" {
			w.buf.WriteString(renderInlineText(c))
			w.buf.WriteString("\n\n")
		}
	}

	var tableNode *html.Node
	var findTable func(*html.Node)
	findTable = func(node *html.Node) {
		if tableNode != nil {
			return
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "table" {
				tableNode = c
				return
			}
			findTable(c)
		}
	}
	findTable(n)

	if tableNode == nil {
		return
	}

	w.renderTable(tableNode)
	w.buf.WriteString("\n")
	w.isBlock = true
}

// renderTable renders a <table> as a Markdown pipe table with headers and rows.
func (w *mdWriter) renderTable(tableNode *html.Node) {
	var headers []string
	var rows [][]string
	var colCount int

	for c := tableNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		switch c.Data {
		case "thead":
			headers = w.renderTableRowCells(c, &colCount)
		case "tbody":
			for tr := c.FirstChild; tr != nil; tr = tr.NextSibling {
				if tr.Type == html.ElementNode && tr.Data == "tr" {
					cells := w.renderTableRowCells(tr, &colCount)
					rows = append(rows, cells)
				}
			}
		}
	}

	if len(headers) == 0 && len(rows) == 0 {
		return
	}

	if len(headers) > 0 {
		w.writeTableRow(headers, colCount)
		w.writeTableSeparator(colCount)
	}

	for _, row := range rows {
		w.writeTableRow(row, colCount)
	}
}

// renderTableRowCells extracts cell text from <th> and <td> elements within a table row container.
func (w *mdWriter) renderTableRowCells(container *html.Node, colCount *int) []string {
	var cells []string
	for c := container.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		switch c.Data {
		case "tr":
			return w.renderTableRowCells(c, colCount)
		case "th", "td":
			content := strings.TrimSpace(renderInlineTextNoWrap(c))
			content = strings.ReplaceAll(content, "\n", " ")
			cells = append(cells, content)
		}
	}
	if len(cells) > *colCount {
		*colCount = len(cells)
	}
	return cells
}

// writeTableRow writes a single Markdown table row with the given cells and column count.
func (w *mdWriter) writeTableRow(cells []string, colCount int) {
	for i := 0; i < colCount; i++ {
		w.buf.WriteString("| ")
		if i < len(cells) {
			w.buf.WriteString(cells[i])
		}
		w.buf.WriteString(" ")
	}
	w.buf.WriteString("|\n")
}

// writeTableSeparator writes the Markdown table header separator row (| --- | --- |).
func (w *mdWriter) writeTableSeparator(colCount int) {
	for i := 0; i < colCount; i++ {
		w.buf.WriteString("| --- ")
	}
	w.buf.WriteString("|\n")
}

// renderEmuFigure renders <emu-figure> as its caption text only.
func (w *mdWriter) renderEmuFigure(n *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && (c.Data == "emu-caption" || c.Data == "figcaption") {
			w.buf.WriteString(renderInlineText(c))
			w.buf.WriteString("\n\n")
		}
	}
	w.isBlock = true
}

// renderEmuProdRef renders <emu-prodref> as italic text.
func (w *mdWriter) renderEmuProdRef(n *html.Node) {
	w.buf.WriteString("_")
	w.buf.WriteString(renderInlineText(n))
	w.buf.WriteString("_")
}

// emitText writes normalized text content, suppressing leading whitespace after block elements.
func (w *mdWriter) emitText(text string) {
	normalized := normalizeWS(text)
	if normalized == "" {
		return
	}
	// why: after a block element (<p>, heading, etc.), leading whitespace-only
	// text nodes should be suppressed to avoid stray indentation.
	if w.isBlock && strings.TrimSpace(normalized) == "" {
		return
	}
	w.isBlock = false
	w.buf.WriteString(normalized)
}

// ensureNewline ensures there is a blank line before the next block element,
// inserting one or two newlines as needed.
func (w *mdWriter) ensureNewline() {
	s := w.buf.String()
	if s == "" {
		return
	}
	// why: we maintain at least one blank line between block elements.
	// If the buffer already ends with "\n\n", no action is needed.
	// If it ends with a single "\n", add one more to get a blank line.
	// Otherwise, add a full "\n\n" to create a blank line.
	if strings.HasSuffix(s, "\n\n") {
		return
	}
	if strings.HasSuffix(s, "\n") {
		w.buf.WriteString("\n")
		return
	}
	w.buf.WriteString("\n\n")
}

// textContentRaw extracts raw text content from an HTML node tree,
// preserving whitespace as-is (no trimming or collapse).
func textContentRaw(node *html.Node) string {
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
	return buf.String()
}

// renderInlineText renders an HTML subtree as a single line of inline text with collapsed whitespace.
func renderInlineText(node *html.Node) string {
	w := &mdWriter{}
	w.walkChildren(node)
	result := strings.TrimSpace(w.buf.String())
	return collapseWS(result)
}

// renderInlineTextNoWrap renders an HTML subtree as inline text without whitespace collapse.
// Used for table cells where newlines are replaced with spaces by the caller.
func renderInlineTextNoWrap(node *html.Node) string {
	w := &mdWriter{}
	w.walkChildren(node)
	return strings.TrimSpace(w.buf.String())
}

// normalizeWS collapses consecutive whitespace (spaces, tabs, newlines) into a single space.
func normalizeWS(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	prevSpace := false
	for _, r := range s {
		if r == '\n' || r == '\t' || r == ' ' {
			if !prevSpace {
				buf.WriteByte(' ')
				prevSpace = true
			}
		} else {
			buf.WriteRune(r)
			prevSpace = false
		}
	}
	return buf.String()
}

// collapseWS collapses all whitespace into single spaces, stripping newlines entirely.
func collapseWS(s string) string {
	var buf strings.Builder
	buf.Grow(len(s))
	prevSpace := false
	for _, r := range s {
		if r == '\n' {
			continue
		}
		if unicode.IsSpace(r) {
			if !prevSpace {
				buf.WriteByte(' ')
				prevSpace = true
			}
		} else {
			buf.WriteRune(r)
			prevSpace = false
		}
	}
	return strings.TrimSpace(buf.String())
}
