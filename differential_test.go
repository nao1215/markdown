package markdown_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/nao1215/markdown"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
)

// The golden files pin the bytes this library emits. They cannot say whether
// those bytes mean what the call meant, because a document can be stable and
// wrong at the same time: an unescaped pipe silently drops a table column, and
// a heading built from text containing a newline silently becomes two.
//
// The tests here close that gap by reading the output back with goldmark in GFM
// mode, which is the parser GitHub's own renderer is closest to, and asserting
// the structure rather than the bytes. Nothing here compares rendered HTML: that
// would pin goldmark's version rather than this library's behavior.

// parse reads src back with the extensions this library targets and returns the
// document node.
func parse(t *testing.T, src string) (ast.Node, []byte) {
	t.Helper()

	md := goldmark.New(goldmark.WithExtensions(extension.GFM, extension.Footnote))
	source := []byte(src)
	return md.Parser().Parse(text.NewReader(source)), source
}

// build runs fn against a fresh builder and returns the document it wrote.
func build(t *testing.T, fn func(*markdown.Markdown) *markdown.Markdown) string {
	t.Helper()

	buf := &bytes.Buffer{}
	if err := fn(markdown.NewMarkdown(buf)).Build(); err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}
	return buf.String()
}

// nodesOfKind collects every node of the given kind, in document order.
func nodesOfKind(root ast.Node, kind ast.NodeKind) []ast.Node {
	var found []ast.Node
	_ = ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering && n.Kind() == kind {
			found = append(found, n)
		}
		return ast.WalkContinue, nil
	})
	return found
}

// cellText returns the text a parsed table cell holds, with the markup its
// inline nodes carry left out.
func cellText(cell ast.Node, source []byte) string {
	var b strings.Builder
	_ = ast.Walk(cell, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if text, ok := n.(*ast.Text); ok {
			b.Write(text.Segment.Value(source))
		}
		return ast.WalkContinue, nil
	})
	return b.String()
}

// firstOfKind returns the first node of the given kind, failing the test when
// the document has none.
func firstOfKind(t *testing.T, root ast.Node, kind ast.NodeKind) ast.Node {
	t.Helper()

	found := nodesOfKind(root, kind)
	if len(found) == 0 {
		t.Fatalf("the document has no %s node", kind)
	}
	return found[0]
}

func TestHeadingsParseAtTheLevelTheyWereBuiltAt(t *testing.T) {
	t.Parallel()

	document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
		return m.H1("One").H2("Two").H3("Three").H4("Four").H5("Five").H6("Six")
	})

	root, _ := parse(t, document)
	headings := nodesOfKind(root, ast.KindHeading)
	if len(headings) != 6 {
		t.Fatalf("parsed %d headings, want 6", len(headings))
	}

	for i, node := range headings {
		heading, ok := node.(*ast.Heading)
		if !ok {
			t.Fatalf("node %d is %T, want *ast.Heading", i, node)
		}
		if want := i + 1; heading.Level != want {
			t.Errorf("heading %d is level %d, want %d", i, heading.Level, want)
		}
	}
}

func TestATableParsesWithTheRowsAndColumnsItWasBuiltWith(t *testing.T) {
	t.Parallel()

	document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
		return m.Table(markdown.TableSet{
			Header: []string{"one", "two", "three"},
			Rows: [][]string{
				{"a", "b", "c"},
				{"d", "e", "f"},
			},
		})
	})

	root, _ := parse(t, document)
	table := firstOfKind(t, root, extast.KindTable)

	rows := 0
	for child := table.FirstChild(); child != nil; child = child.NextSibling() {
		cells := child.ChildCount()
		if cells != 3 {
			t.Errorf("a table row parsed with %d cells, want 3", cells)
		}
		if child.Kind() == extast.KindTableRow {
			rows++
		}
	}
	if rows != 2 {
		t.Errorf("the table parsed with %d body rows, want 2", rows)
	}
}

func TestAlignmentReachesTheParsedTable(t *testing.T) {
	t.Parallel()

	document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
		return m.Table(markdown.TableSet{
			Header: []string{"left", "center", "right"},
			Rows:   [][]string{{"a", "b", "c"}},
			Alignment: []markdown.TableAlignment{
				markdown.AlignLeft,
				markdown.AlignCenter,
				markdown.AlignRight,
			},
		})
	})

	root, _ := parse(t, document)
	table := firstOfKind(t, root, extast.KindTable)

	header := table.FirstChild()
	if header == nil {
		t.Fatal("the table parsed with no header row")
	}

	want := []extast.Alignment{extast.AlignLeft, extast.AlignCenter, extast.AlignRight}
	i := 0
	for cell := header.FirstChild(); cell != nil; cell = cell.NextSibling() {
		parsed, ok := cell.(*extast.TableCell)
		if !ok {
			t.Fatalf("header child %d is %T, want *extast.TableCell", i, cell)
		}
		if i >= len(want) {
			t.Fatalf("the header parsed with more than %d cells", len(want))
		}
		if parsed.Alignment != want[i] {
			t.Errorf("column %d parsed with alignment %v, want %v", i, parsed.Alignment, want[i])
		}
		i++
	}
	if i != len(want) {
		t.Errorf("the header parsed with %d cells, want %d", i, len(want))
	}
}

func TestListsParseWithTheItemsTheyWereBuiltWith(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build     func(*markdown.Markdown) *markdown.Markdown
		items     int
		isOrdered bool
	}{
		"bullet list": {
			build: func(m *markdown.Markdown) *markdown.Markdown {
				return m.BulletList("one", "two", "three")
			},
			items: 3,
		},
		"ordered list": {
			build: func(m *markdown.Markdown) *markdown.Markdown {
				return m.OrderedList("one", "two")
			},
			items:     2,
			isOrdered: true,
		},
		"checkbox list": {
			build: func(m *markdown.Markdown) *markdown.Markdown {
				return m.CheckBox([]markdown.CheckBoxSet{
					{Checked: true, Text: "done"},
					{Checked: false, Text: "not done"},
				})
			},
			items: 2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			root, _ := parse(t, build(t, tt.build))
			list, ok := firstOfKind(t, root, ast.KindList).(*ast.List)
			if !ok {
				t.Fatal("the first list node is not an *ast.List")
			}

			if got := list.ChildCount(); got != tt.items {
				t.Errorf("the list parsed with %d items, want %d", got, tt.items)
			}
			if list.IsOrdered() != tt.isOrdered {
				t.Errorf("the list parsed as ordered=%v, want %v", list.IsOrdered(), tt.isOrdered)
			}
		})
	}
}

func TestCheckboxesParseAsTaskItems(t *testing.T) {
	t.Parallel()

	document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
		return m.CheckBox([]markdown.CheckBoxSet{
			{Checked: true, Text: "done"},
			{Checked: false, Text: "not done"},
		})
	})

	root, _ := parse(t, document)
	tasks := nodesOfKind(root, extast.KindTaskCheckBox)
	if len(tasks) != 2 {
		t.Fatalf("parsed %d task checkboxes, want 2", len(tasks))
	}

	want := []bool{true, false}
	for i, node := range tasks {
		task, ok := node.(*extast.TaskCheckBox)
		if !ok {
			t.Fatalf("task %d is %T, want *extast.TaskCheckBox", i, node)
		}
		if task.IsChecked != want[i] {
			t.Errorf("task %d parsed as checked=%v, want %v", i, task.IsChecked, want[i])
		}
	}
}

func TestACodeBlockParsesWithItsLanguage(t *testing.T) {
	t.Parallel()

	document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
		return m.CodeBlocks(markdown.SyntaxHighlightGo, "package main")
	})

	root, source := parse(t, document)
	block, ok := firstOfKind(t, root, ast.KindFencedCodeBlock).(*ast.FencedCodeBlock)
	if !ok {
		t.Fatal("the first fenced code block node is not an *ast.FencedCodeBlock")
	}

	if got := string(block.Language(source)); got != "go" {
		t.Errorf("the code block parsed with language %q, want %q", got, "go")
	}
}

func TestQuotesAndAlertsParseAsBlockquotes(t *testing.T) {
	t.Parallel()

	tests := map[string]func(*markdown.Markdown) *markdown.Markdown{
		"blockquote": func(m *markdown.Markdown) *markdown.Markdown { return m.Blockquote("quoted") },
		"note":       func(m *markdown.Markdown) *markdown.Markdown { return m.Note("a note") },
		"tip":        func(m *markdown.Markdown) *markdown.Markdown { return m.Tip("a tip") },
		"important":  func(m *markdown.Markdown) *markdown.Markdown { return m.Important("important") },
		"warning":    func(m *markdown.Markdown) *markdown.Markdown { return m.Warning("a warning") },
		"caution":    func(m *markdown.Markdown) *markdown.Markdown { return m.Caution("a caution") },
	}

	for name, fn := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			root, _ := parse(t, build(t, fn))
			if got := len(nodesOfKind(root, ast.KindBlockquote)); got != 1 {
				t.Errorf("parsed %d blockquotes, want 1", got)
			}
		})
	}
}

func TestAMultiLineAlertStaysOneBlockquote(t *testing.T) {
	t.Parallel()

	// Every line of an alert carries the marker, so a body with a list in it
	// belongs to the callout instead of escaping it and rendering as a sibling.
	document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
		return m.Note("first line\n\n- an item\n- another item")
	})

	root, _ := parse(t, document)
	quotes := nodesOfKind(root, ast.KindBlockquote)
	if len(quotes) != 1 {
		t.Fatalf("parsed %d blockquotes, want 1: the alert body escaped the callout", len(quotes))
	}
	if got := len(nodesOfKind(quotes[0], ast.KindList)); got != 1 {
		t.Errorf("the blockquote holds %d lists, want 1", got)
	}
}

func TestLinksAndImagesParseWithTheirDestination(t *testing.T) {
	t.Parallel()

	const url = "https://example.com/page?a=1&b=2"

	t.Run("link", func(t *testing.T) {
		t.Parallel()

		document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
			return m.PlainText(markdown.Link("the text", url))
		})

		root, _ := parse(t, document)
		link, ok := firstOfKind(t, root, ast.KindLink).(*ast.Link)
		if !ok {
			t.Fatal("the first link node is not an *ast.Link")
		}
		if got := string(link.Destination); got != url {
			t.Errorf("the link parsed with destination %q, want %q", got, url)
		}
	})

	t.Run("image", func(t *testing.T) {
		t.Parallel()

		document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
			return m.PlainText(markdown.Image("the alt text", url))
		})

		root, _ := parse(t, document)
		image, ok := firstOfKind(t, root, ast.KindImage).(*ast.Image)
		if !ok {
			t.Fatal("the first image node is not an *ast.Image")
		}
		if got := string(image.Destination); got != url {
			t.Errorf("the image parsed with destination %q, want %q", got, url)
		}
	})
}

func TestAFootnoteReferenceFindsItsDefinition(t *testing.T) {
	t.Parallel()

	document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
		return m.
			PlainTextf("a statement%s", markdown.FootnoteReference("1")).
			BlankLine().
			PlainText(markdown.FootnoteDefinition("1", "the note")).
			BlankLine()
	})

	root, _ := parse(t, document)
	if got := len(nodesOfKind(root, extast.KindFootnoteLink)); got != 1 {
		t.Errorf("parsed %d footnote references, want 1", got)
	}
	if got := len(nodesOfKind(root, extast.KindFootnote)); got != 1 {
		t.Errorf("parsed %d footnote definitions, want 1", got)
	}
}

// TestEscapedCellsStayOneCell is the point of EscapeTableCell. A pipe in the
// data ends the cell and drops the last column of the row; a newline ends the
// row outright. Neither produces an error, and the row still has the right
// length before it is serialized, so nothing but a parser can catch it.
func TestEscapedCellsStayOneCell(t *testing.T) {
	t.Parallel()

	cells := map[string]string{
		"a pipe":               "a|b",
		"two pipes":            "a|b|c",
		"a newline":            "first\nsecond",
		"a carriage return":    "first\r\nsecond",
		"a backtick":           "`code`",
		"emphasis markers":     "*star* _underscore_",
		"brackets":             "[not a link]",
		"a backslash":          `a\b`,
		"an escaped pipe":      `a\|b`,
		"a trailing backslash": `a\`,
	}

	for name, cell := range cells {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
				return m.Table(markdown.TableSet{
					Header:      []string{"one", "two"},
					Rows:        [][]string{{cell, "second column"}},
					EscapeCells: true,
				})
			})

			root, source := parse(t, document)
			table := firstOfKind(t, root, extast.KindTable)

			rows := nodesOfKind(table, extast.KindTableRow)
			if len(rows) != 1 {
				t.Fatalf("the table parsed with %d body rows, want 1: %q", len(rows), document)
			}
			if got := rows[0].ChildCount(); got != 2 {
				t.Errorf("the row parsed with %d cells, want 2: %q", got, document)
			}
			// The column after the offending one is the one that disappears
			// when the cell is not escaped, so it is the one worth asserting.
			last := rows[0].LastChild()
			if last == nil {
				t.Fatalf("the row parsed with no cells: %q", document)
			}
			if got := cellText(last, source); got != "second column" {
				t.Errorf("the last cell holds %q, want %q: %q", got, "second column", document)
			}
		})
	}
}

// TestAnUnescapedPipeIsWhyEscapeCellsExists pins the failure the option
// prevents, so that a change making escaping unnecessary shows up here rather
// than leaving a dead option behind.
//
// The failure is quiet. GFM drops the cells a row has beyond its header, so a
// pipe in the first value pushes the last column off the end of the row and the
// table still parses, one column short, with no error anywhere.
func TestAnUnescapedPipeIsWhyEscapeCellsExists(t *testing.T) {
	t.Parallel()

	document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
		return m.Table(markdown.TableSet{
			Header: []string{"one", "two"},
			Rows:   [][]string{{"a|b", "second column"}},
		})
	})

	root, source := parse(t, document)
	table := firstOfKind(t, root, extast.KindTable)

	rows := nodesOfKind(table, extast.KindTableRow)
	if len(rows) != 1 {
		t.Fatalf("the table parsed with %d body rows, want 1", len(rows))
	}

	last := rows[0].LastChild()
	if last == nil {
		t.Fatal("the row parsed with no cells")
	}
	if got := cellText(last, source); got == "second column" {
		t.Error("the unescaped pipe no longer costs the row its last column; EscapeCells may no longer be needed")
	}
}

// TestHeadingTextDoesNotProduceExtraHeadings pins that text carrying its own
// markdown cannot turn one heading into several.
func TestHeadingTextDoesNotProduceExtraHeadings(t *testing.T) {
	t.Parallel()

	texts := map[string]string{
		"a hash":            "# not a heading",
		"trailing hashes":   "text ###",
		"an underline rule": "text\n---",
	}

	for name, heading := range texts {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
				return m.H2(heading)
			})

			root, _ := parse(t, document)
			if got := len(nodesOfKind(root, ast.KindHeading)); got != 1 {
				t.Errorf("parsed %d headings from one H2 call, want 1: %q", got, document)
			}
		})
	}
}

// TestEveryMermaidSubpackageProducesAMermaidBlock covers the one part of a
// diagram a markdown parser sees. The body is opaque to it, and the renderer
// checks that separately in continuous integration.
func TestEveryMermaidSubpackageProducesAMermaidBlock(t *testing.T) {
	t.Parallel()

	for _, diagram := range mermaidDiagrams(t) {
		t.Run(diagram.name, func(t *testing.T) {
			t.Parallel()

			document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
				return m.CodeBlocks(markdown.SyntaxHighlightMermaid, diagram.body)
			})

			root, source := parse(t, document)
			block, ok := firstOfKind(t, root, ast.KindFencedCodeBlock).(*ast.FencedCodeBlock)
			if !ok {
				t.Fatal("the first fenced code block node is not an *ast.FencedCodeBlock")
			}
			if got := string(block.Language(source)); got != "mermaid" {
				t.Errorf("the block parsed with language %q, want %q", got, "mermaid")
			}

			var body strings.Builder
			for i := 0; i < block.Lines().Len(); i++ {
				line := block.Lines().At(i)
				body.WriteString(string(line.Value(source)))
			}
			if !strings.Contains(body.String(), diagram.body) {
				t.Errorf("the block parsed with body %q, want it to contain the diagram", body.String())
			}
		})
	}
}

// TestGeneratedDocumentsHoldNoStrayFence guards the invariant that keeps a
// mermaid diagram inside its block: a diagram body carrying "```" would end the
// fence early and spill the rest of the document into the page.
func TestGeneratedDocumentsHoldNoStrayFence(t *testing.T) {
	t.Parallel()

	for _, diagram := range mermaidDiagrams(t) {
		t.Run(diagram.name, func(t *testing.T) {
			t.Parallel()

			if strings.Contains(diagram.body, "```") {
				t.Errorf("the diagram body holds a fence, which would end the code block early: %q", diagram.body)
			}
		})
	}
}

// diagram is one mermaid subpackage's output, named for the package.
type diagram struct {
	name string
	body string
}

// mermaidDiagrams builds one diagram per mermaid subpackage.
func mermaidDiagrams(t *testing.T) []diagram {
	t.Helper()

	builders := mermaidBuilders()
	diagrams := make([]diagram, 0, len(builders))
	for _, b := range builders {
		body := b.build()
		if body == "" {
			t.Fatalf("the %s builder produced an empty diagram", b.name)
		}
		diagrams = append(diagrams, diagram{name: b.name, body: body})
	}
	return diagrams
}

// String makes a diagram printable in a failure message.
func (d diagram) String() string {
	return fmt.Sprintf("%s: %q", d.name, d.body)
}
