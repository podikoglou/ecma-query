package spec

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func parseHTMLStr(t *testing.T, htmlStr string) *html.Node {
	t.Helper()
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		t.Fatalf("html.Parse error: %v", err)
	}
	return doc
}

func extractBody(doc *html.Node) *html.Node {
	var body *html.Node
	var find func(*html.Node)
	find = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			find(c)
		}
	}
	find(doc)
	return body
}

func TestRenderHTML_simpleParagraph(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><p>Hello <b>world</b></p></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "Hello **world**") {
		t.Errorf("expected 'Hello **world**', got %q", result)
	}
}

func TestRenderHTML_inlineFormatting(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><p><b>bold</b> <em>italic</em> <code>code</code> <var>variable</var></p></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "**bold** *italic* `code` _variable_") {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestRenderHTML_links(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><p><a href="https://example.com">Example</a></p></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "[Example](https://example.com)") {
		t.Errorf("expected [Example](https://example.com), got %q", result)
	}
}

func TestRenderHTML_headings(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><h1>Heading 1</h1><h2>Heading 2</h2><h3>Heading 3</h3></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "# Heading 1") {
		t.Errorf("expected '# Heading 1', got %q", result)
	}
	if !strings.Contains(result, "## Heading 2") {
		t.Errorf("expected '## Heading 2', got %q", result)
	}
	if !strings.Contains(result, "### Heading 3") {
		t.Errorf("expected '### Heading 3', got %q", result)
	}
}

func TestRenderHTML_unorderedList(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><ul><li>Item 1</li><li>Item 2</li></ul></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "- Item 1") {
		t.Errorf("expected '- Item 1', got %q", result)
	}
	if !strings.Contains(result, "- Item 2") {
		t.Errorf("expected '- Item 2', got %q", result)
	}
}

func TestRenderHTML_orderedList(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><ol><li>First</li><li>Second</li></ol></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "1. First") {
		t.Errorf("expected '1. First', got %q", result)
	}
	if !strings.Contains(result, "2. Second") {
		t.Errorf("expected '2. Second', got %q", result)
	}
}

func TestRenderHTML_emuAlg(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><emu-alg>
  1. If _input_ is an Object, then
    1. Let _exoticToPrim_ be ? GetMethod(_input_, %Symbol.toPrimitive%).
    1. If _exoticToPrim_ is not *undefined*, then
      1. Let _result_ be ? Call(_exoticToPrim_, _input_, "preferredType").
  2. Return _input_.
</emu-alg></body></html>`)
	result := RenderHTML(extractBody(doc))

	t.Logf("emu-alg result:\n%s", result)

	if !strings.Contains(result, "1. If _input_ is an Object,") {
		t.Errorf("expected top-level step 1, got %q", result)
	}
	if !strings.Contains(result, "  1. Let _exoticToPrim_") {
		t.Errorf("expected nested step at indent, got %q", result)
	}
	if !strings.Contains(result, "    1. Let _result_") {
		t.Errorf("expected double-nested step, got %q", result)
	}
	if !strings.Contains(result, "2. Return _input_") {
		t.Errorf("expected top-level step 2, got %q", result)
	}
}

func TestRenderHTML_emuXRef(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><p><emu-xref href="#sec-toprimitive">ToPrimitive</emu-xref></p></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "[ToPrimitive](#sec-toprimitive)") {
		t.Errorf("expected cross-reference, got %q", result)
	}
}

func TestRenderHTML_emuNote(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><emu-note><span class="note">Note 1</span>This is a note.</emu-note></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "> **") {
		t.Errorf("expected blockquote, got %q", result)
	}
	if !strings.Contains(result, "Note 1") {
		t.Errorf("expected note text, got %q", result)
	}
}

func TestRenderHTML_emuGrammar(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><emu-grammar>Expression :
  AssignmentExpression</emu-grammar></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "```") {
		t.Errorf("expected fenced code block, got %q", result)
	}
	if !strings.Contains(result, "Expression") {
		t.Errorf("expected grammar content, got %q", result)
	}
}

func TestRenderHTML_table(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><emu-table>
  <emu-caption>Test Table</emu-caption>
  <table>
    <thead><tr><th>A</th><th>B</th></tr></thead>
    <tbody><tr><td>1</td><td>2</td></tr></tbody>
  </table>
</emu-table></body></html>`)
	result := RenderHTML(extractBody(doc))

	t.Logf("table result:\n%s", result)

	if !strings.Contains(result, "Test Table") {
		t.Errorf("expected caption, got %q", result)
	}
	if !strings.Contains(result, "| A | B |") {
		t.Errorf("expected header row, got %q", result)
	}
	if !strings.Contains(result, "| --- | --- |") {
		t.Errorf("expected separator, got %q", result)
	}
	if !strings.Contains(result, "| 1 | 2 |") {
		t.Errorf("expected data row, got %q", result)
	}
}

func TestRenderHTML_emuVal(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><p><emu-val>*undefined*</emu-val></p></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "*undefined*") {
		t.Errorf("expected *undefined*, got %q", result)
	}
}

func TestRenderHTML_emuProdRef(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><p><emu-prodref>IdentifierReference</emu-prodref></p></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "_IdentifierReference_") {
		t.Errorf("expected _IdentifierReference_, got %q", result)
	}
}

func TestRenderHTML_emuEqn(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><p><emu-eqn>1 + 1 = 2</emu-eqn></p></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "`1 + 1 = 2`") {
		t.Errorf("expected inline code, got %q", result)
	}
}

func TestRenderHTML_emuFigure(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><emu-figure><emu-caption>Figure caption</emu-caption><img src="foo.png"></emu-figure></body></html>`)
	result := RenderHTML(extractBody(doc))
	if strings.Contains(result, "![]") {
		t.Errorf("expected no image, got %q", result)
	}
	if !strings.Contains(result, "Figure caption") {
		t.Errorf("expected caption, got %q", result)
	}
}

func TestRenderHTML_skipElements(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><div><span><ins>content</ins></span></div></body></html>`)
	result := RenderHTML(extractBody(doc))
	if strings.Contains(result, "<div>") || strings.Contains(result, "<span>") || strings.Contains(result, "<ins>") {
		t.Errorf("expected no HTML wrapper tags, got %q", result)
	}
	if !strings.Contains(result, "content") {
		t.Errorf("expected inner content, got %q", result)
	}
}

func TestRenderHTML_emuMetaSkip(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><div><emu-meta>metadata</emu-meta>visible</div></body></html>`)
	result := RenderHTML(extractBody(doc))
	if strings.Contains(result, "metadata") {
		t.Errorf("expected emu-meta skipped, got %q", result)
	}
	if !strings.Contains(result, "visible") {
		t.Errorf("expected visible text, got %q", result)
	}
}

func TestRenderHTML_pre(t *testing.T) {
	doc := parseHTMLStr(t, "<html><body><pre>line1\nline2</pre></body></html>")
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "```") {
		t.Errorf("expected fenced code block, got %q", result)
	}
	if !strings.Contains(result, "line1") {
		t.Errorf("expected content, got %q", result)
	}
}

func TestRenderHTML_dfn(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><p><dfn>defined term</dfn></p></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "*defined term*") {
		t.Errorf("expected *defined term*, got %q", result)
	}
}

func TestRenderHTML_subSup(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><p>x<sub>1</sub> y<sup>2</sup></p></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "x1 y2") {
		t.Errorf("expected x1 y2, got %q", result)
	}
}

func TestRenderHTML_br(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><p>line1<br>line2</p></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "line1") || !strings.Contains(result, "line2") {
		t.Errorf("expected line1 and line2, got %q", result)
	}
}

func TestRenderHTML_defList(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><dl><dt>Term</dt><dd>Definition</dd></dl></body></html>`)
	result := RenderHTML(extractBody(doc))
	if !strings.Contains(result, "**Term**: Definition") {
		t.Errorf("expected '**Term**: Definition', got %q", result)
	}
}

func TestRenderHTML_empty(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body></body></html>`)
	result := RenderHTML(extractBody(doc))
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestRenderHTML_emuTableNoTable(t *testing.T) {
	doc := parseHTMLStr(t, `<html><body><emu-table>No table here</emu-table></body></html>`)
	result := RenderHTML(extractBody(doc))
	if strings.Contains(result, "|") {
		t.Errorf("expected no table output, got %q", result)
	}
}

func TestRenderNode_fullClause(t *testing.T) {
	htmlStr := `<html><body><emu-clause id="sec-test">
  <h1>1.0 Test Clause</h1>
  <p>First paragraph.</p>
  <emu-clause id="sec-child">
    <h2>1.0.1 Child</h2>
    <p>Child content <emu-xref href="#sec-target">Ref</emu-xref>.</p>
  </emu-clause>
  <p>More content.</p>
</emu-clause></body></html>`
	s, err := Parse([]byte(htmlStr))
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}

	if len(s.Nodes) > 0 {
		s.Nodes[0].Section = "1.0"
	}
	if len(s.Nodes) > 1 {
		s.Nodes[1].Section = "1.0.1"
		s.Nodes[1].Parent = s.Nodes[0]
		s.Nodes[0].Children = []*ClauseNode{s.Nodes[1]}
	}

	result, err := s.RenderNode(s.Nodes[0])
	if err != nil {
		t.Fatalf("RenderNode() error: %v", err)
	}

	t.Logf("RenderNode result:\n%s", result)

	if !strings.Contains(result, "## 1.0 Test Clause") {
		t.Errorf("expected heading, got %q", result)
	}
	if !strings.Contains(result, "First paragraph") {
		t.Errorf("expected first paragraph, got %q", result)
	}
	if !strings.Contains(result, "#sec-target") {
		t.Errorf("expected xref target, got %q", result)
	}
}

func TestRenderNode_toprimitive(t *testing.T) {
	s, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	cn, ok := s.ByID["sec-toprimitive"]
	if !ok {
		t.Skip("sec-toprimitive not found in spec")
	}

	result, err := s.RenderNode(cn)
	if err != nil {
		t.Fatalf("RenderNode() error: %v", err)
	}

	limit := len(result)
	if limit > 2000 {
		limit = 2000
	}
	t.Logf("ToPrimitive render:\n%s\n---END---", result[:limit])

	if !strings.Contains(result, "##") {
		t.Error("expected markdown heading")
	}
	if !strings.Contains(result, "ToPrimitive") {
		t.Error("expected ToPrimitive in output")
	}
	if !strings.Contains(result, "See also") {
		t.Error("expected 'See also' section")
	}

	if strings.Contains(result, "<emu-") {
		t.Error("output contains raw emu- tags")
	}
	if strings.Contains(result, "<p>") {
		t.Error("output contains raw <p> tags")
	}
	if strings.Contains(result, "<ol>") || strings.Contains(result, "<li>") {
		t.Error("output contains raw list tags")
	}
}
