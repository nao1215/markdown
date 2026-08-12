package markdown_test

import (
	"bytes"
	"errors"
	"fmt"
	goast "go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/arch"
	"github.com/nao1215/markdown/mermaid/block"
	"github.com/nao1215/markdown/mermaid/c4"
	"github.com/nao1215/markdown/mermaid/class"
	"github.com/nao1215/markdown/mermaid/er"
	"github.com/nao1215/markdown/mermaid/flowchart"
	"github.com/nao1215/markdown/mermaid/gantt"
	"github.com/nao1215/markdown/mermaid/gitgraph"
	"github.com/nao1215/markdown/mermaid/kanban"
	"github.com/nao1215/markdown/mermaid/mindmap"
	"github.com/nao1215/markdown/mermaid/packet"
	"github.com/nao1215/markdown/mermaid/piechart"
	"github.com/nao1215/markdown/mermaid/quadrant"
	"github.com/nao1215/markdown/mermaid/requirement"
	"github.com/nao1215/markdown/mermaid/sequence"
	"github.com/nao1215/markdown/mermaid/state"
	"github.com/nao1215/markdown/mermaid/userjourney"
	"github.com/nao1215/markdown/mermaid/venn"
	"github.com/nao1215/markdown/mermaid/xychart"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v3"
)

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return markdown.NewMarkdown(w).H1("Title")
	})
}

// TestRecordedErrorContract asserts that a table whose rows do not match its
// header surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return markdown.NewMarkdown(w).Table(markdown.TableSet{
			Header: []string{"one", "two"},
			Rows:   [][]string{{"only one cell"}},
		})
	})
}

// TestErrorSurfacesFromBothErrorAndBuild pins that the two ways of asking for
// the error agree. Callers use whichever suits the shape of their code, and a
// builder that reported an error from only one of them would be a trap.
func TestErrorSurfacesFromBothErrorAndBuild(t *testing.T) {
	t.Parallel()

	m := markdown.NewMarkdown(io.Discard).Table(markdown.TableSet{
		Header: []string{"one", "two"},
		Rows:   [][]string{{"only one cell"}},
	})

	fromError := m.Error()
	if fromError == nil {
		t.Fatal("Error() = nil, want the error the chain recorded")
	}
	if !errors.Is(fromError, markdown.ErrMismatchColumn) {
		t.Errorf("Error() = %v, want an error wrapping ErrMismatchColumn", fromError)
	}

	fromBuild := m.Build()
	if fromBuild == nil {
		t.Fatal("Build() = nil, want the error the chain recorded")
	}
	if !errors.Is(fromBuild, markdown.ErrMismatchColumn) {
		t.Errorf("Build() = %v, want an error wrapping ErrMismatchColumn", fromBuild)
	}
}

// TestTheChainContinuesAfterAnError pins that a rejected call does not stop the
// document. Callers write one long chain and check it once at the end, so the
// blocks after a bad table still have to appear.
func TestTheChainContinuesAfterAnError(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := markdown.NewMarkdown(buf).
		H1("Before").
		Table(markdown.TableSet{
			Header: []string{"one", "two"},
			Rows:   [][]string{{"only one cell"}},
		}).
		H2("After").
		Build()
	if err == nil {
		t.Fatal("Build() = nil, want the error the chain recorded")
	}

	for _, want := range []string{"# Before", "## After"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("document = %q, want it to contain %q", buf.String(), want)
		}
	}
}

// TestTheFirstErrorIsKept pins which error a caller sees when a chain records
// more than one. The first is the one that explains the rest, so it is the one
// that survives; the later ones are appended to its message rather than
// replacing it.
func TestTheFirstErrorIsKept(t *testing.T) {
	t.Parallel()

	m := markdown.NewMarkdown(io.Discard).
		TableOfContents(markdown.TableOfContentsDepthH3).
		TableOfContents(markdown.TableOfContentsDepthH3)

	err := m.Error()
	if err == nil {
		t.Fatal("Error() = nil, want the error the chain recorded")
	}
	if want := "table of contents has already been generated"; !strings.Contains(err.Error(), want) {
		t.Errorf("Error() = %v, want it to mention %q", err, want)
	}
}

// TestBuildWithNilWriterKeepsTheEarlierError pins that a nil writer does not
// hide what went wrong before it. Both failures are worth knowing about, so the
// message carries the earlier one too.
func TestBuildWithNilWriterKeepsTheEarlierError(t *testing.T) {
	t.Parallel()

	err := markdown.NewMarkdown(nil).
		Table(markdown.TableSet{
			Header: []string{"one", "two"},
			Rows:   [][]string{{"only one cell"}},
		}).
		Build()
	if err == nil {
		t.Fatal("Build() = nil, want an error")
	}

	for _, want := range []string{"destination writer is nil", markdown.ErrMismatchColumn.Error()} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Build() = %v, want it to mention %q", err, want)
		}
	}
}

// TestBuildWithFailingWriterKeepsTheEarlierError pins the same for a writer
// that refuses the document.
func TestBuildWithFailingWriterKeepsTheEarlierError(t *testing.T) {
	t.Parallel()

	err := markdown.NewMarkdown(buildertest.FailingWriter{}).
		Table(markdown.TableSet{
			Header: []string{"one", "two"},
			Rows:   [][]string{{"only one cell"}},
		}).
		Build()
	if err == nil {
		t.Fatal("Build() = nil, want an error")
	}

	for _, want := range []string{"failed to write markdown text", markdown.ErrMismatchColumn.Error()} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Build() = %v, want it to mention %q", err, want)
		}
	}
}

// TestBuildEndsTheDocumentWithALineFeed pins the trailing newline. markdownlint
// MD047 requires it, and without it a second document written to the same
// writer would splice its first line onto the last line of this one.
func TestBuildEndsTheDocumentWithALineFeed(t *testing.T) {
	t.Parallel()

	tests := map[string]func(*markdown.Markdown) *markdown.Markdown{
		"a document ending in a heading": func(m *markdown.Markdown) *markdown.Markdown {
			return m.H1("Title")
		},
		"a document ending in a table": func(m *markdown.Markdown) *markdown.Markdown {
			return m.Table(markdown.TableSet{
				Header: []string{"key", "value"},
				Rows:   [][]string{{"a", "b"}},
			})
		},
		"an empty document": func(m *markdown.Markdown) *markdown.Markdown {
			return m
		},
	}

	for name, build := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			if err := build(markdown.NewMarkdown(buf)).Build(); err != nil {
				t.Fatalf("Build() = %v, want nil", err)
			}
			if !strings.HasSuffix(buf.String(), "\n") {
				t.Errorf("document = %q, want it to end with a line feed", buf.String())
			}
		})
	}
}

// TestGoldenDocument pins the output of every builder method of the root
// package. The document is deliberately one long chain rather than several
// small ones, because the blank lines the builder inserts between blocks depend
// on which blocks are adjacent, and that spacing is part of the contract too.
func TestGoldenDocument(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := markdown.NewMarkdown(buf).
		H1("Golden Document").
		PlainText("A paragraph written with PlainText.").
		PlainTextf("A paragraph written with %s.", "PlainTextf").
		H2("Headings").
		H3("Third level").
		H4("Fourth level").
		H5("Fifth level").
		H6("Sixth level").
		H1f("%s built with H1f", "Heading").
		H2f("%s built with H2f", "Heading").
		H3f("%s built with H3f", "Heading").
		H4f("%s built with H4f", "Heading").
		H5f("%s built with H5f", "Heading").
		H6f("%s built with H6f", "Heading").
		H2("Inline syntax").
		PlainText(markdown.Bold("bold")).
		PlainText(markdown.Italic("italic")).
		PlainText(markdown.BoldItalic("bold italic")).
		PlainText(markdown.Strikethrough("strikethrough")).
		PlainText(markdown.Code("code")).
		PlainText(markdown.Highlight("highlight")).
		PlainText(markdown.Link("nao1215/markdown", "https://github.com/nao1215/markdown")).
		PlainText(markdown.Image("gopher", "https://go.dev/images/gophers/ladder.svg")).
		PlainText(markdown.ReferenceLink("The Go Programming Language", "go-site")).
		PlainText(markdown.ReferenceLinkDefinition("go-site", "https://go.dev")).
		PlainText(markdown.ReferenceLinkDefinition("go-doc", "https://pkg.go.dev", `The "Go" package index`)).
		PlainTextf("A statement with a footnote%s", markdown.FootnoteReference("1")).
		PlainText(markdown.FootnoteDefinition("1", "The footnote body.")).
		H2("Math").
		PlainTextf("Inline math: %s", markdown.InlineMath("E=mc^2")).
		PlainText(markdown.BlockMath("x^2 + y^2 = z^2")).
		H2("Lists").
		BulletList("first bullet", "second bullet").
		OrderedList("first item", "second item").
		CheckBox([]markdown.CheckBoxSet{
			{Checked: true, Text: "done"},
			{Checked: false, Text: "not done"},
		}).
		H2("Quotes and alerts").
		Blockquote("A quote\nspanning two lines.").
		Note("A note.").
		Notef("A note built with %s.", "Notef").
		Tip("A tip.").
		Tipf("A tip built with %s.", "Tipf").
		Important("Something important.").
		Importantf("Something important built with %s.", "Importantf").
		Warning("A warning.").
		Warningf("A warning built with %s.", "Warningf").
		Caution("A caution.").
		Cautionf("A caution built with %s.", "Cautionf").
		H2("Badges").
		RedBadge("red").
		RedBadgef("%s-formatted", "red").
		YellowBadge("yellow").
		YellowBadgef("%s-formatted", "yellow").
		GreenBadge("green").
		GreenBadgef("%s-formatted", "green").
		BlueBadge("blue").
		BlueBadgef("%s-formatted", "blue").
		H2("Code blocks").
		CodeBlocks(markdown.SyntaxHighlightGo, `func main() {
	fmt.Println("hello")
}`).
		CodeBlocks(markdown.SyntaxHighlightNone, "plain text block").
		H2("Tables").
		Table(markdown.TableSet{
			Header: []string{"Name", "Kind", "Note"},
			Rows: [][]string{
				{"markdown", "library", "builder"},
				{"mermaid", "sub package", "diagrams"},
			},
		}).
		H2("Details").
		Details("Summary", "Hidden body.").
		Detailsf("Formatted summary", "Hidden body built with %s.", "Detailsf").
		H2("Separators").
		LF().
		BlankLine().
		HorizontalRule().
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("document.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenDocumentWithBlockSpacing pins the output of the same document
// shape under WithBlockSpacing, which is the only construction time option of
// the root package.
func TestGoldenDocumentWithBlockSpacing(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := markdown.NewMarkdown(buf, markdown.WithBlockSpacing()).
		H1("Block Spacing").
		PlainText("A paragraph.").
		BulletList("first bullet", "second bullet").
		OrderedList("first item", "second item").
		Blockquote("A quote.").
		Note("A note.").
		CodeBlocks(markdown.SyntaxHighlightShell, "echo hello").
		Table(markdown.TableSet{
			Header: []string{"Key", "Value"},
			Rows:   [][]string{{"spacing", "on"}},
		}).
		HorizontalRule().
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("document_block_spacing.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenTables pins table rendering: every alignment, cell escaping, and
// the tablewriter backed CustomTable with both of its options.
func TestGoldenTables(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := markdown.NewMarkdown(buf).
		H2("Alignment").
		Table(markdown.TableSet{
			Header: []string{"default", "left", "center", "right"},
			Rows: [][]string{
				{"a", "b", "c", "d"},
				{"e", "f", "g", "h"},
			},
			Alignment: []markdown.TableAlignment{
				markdown.AlignDefault,
				markdown.AlignLeft,
				markdown.AlignCenter,
				markdown.AlignRight,
			},
		}).
		H2("Escaped cells").
		Table(markdown.TableSet{
			Header: []string{"input", "note"},
			Rows: [][]string{
				{"a|b", "a pipe would split the cell"},
				{"line\nbreak", "a newline would end the row"},
			},
			EscapeCells: true,
		}).
		H2("Escaped by hand").
		PlainText(markdown.EscapeTableCell("a|b")).
		H2("Custom table").
		CustomTable(
			markdown.TableSet{
				Header: []string{"name", "value"},
				Rows: [][]string{
					{"custom", "table"},
					{"auto", "format"},
				},
				Alignment: []markdown.TableAlignment{
					markdown.AlignLeft,
					markdown.AlignRight,
				},
			},
			markdown.TableOptions{AutoWrapText: false, AutoFormatHeaders: true},
		).
		CustomTable(
			markdown.TableSet{
				Header: []string{"name", "value"},
				Rows:   [][]string{{"wrapped", "a longer value that may be wrapped by tablewriter"}},
			},
			markdown.TableOptions{AutoWrapText: true, AutoFormatHeaders: false},
		).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("tables.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenTableOfContents pins both table of contents entry points, including
// the anchor suffix that repeated headings get.
func TestGoldenTableOfContents(t *testing.T) {
	t.Parallel()

	t.Run("full depth", func(t *testing.T) {
		t.Parallel()

		buf := &bytes.Buffer{}
		err := markdown.NewMarkdown(buf).
			H1("Table Of Contents").
			TableOfContents(markdown.TableOfContentsDepthH3).
			H2("Section").
			H3("Subsection").
			H4("Excluded from the table of contents").
			H2("Section").
			H3("Subsection with a Symbol: !").
			Build()
		if err != nil {
			t.Fatalf("Build() = %v, want nil", err)
		}

		if err := golden.Assert("toc.md", buf.String()); err != nil {
			t.Error(err)
		}
	})

	t.Run("depth range", func(t *testing.T) {
		t.Parallel()

		buf := &bytes.Buffer{}
		err := markdown.NewMarkdown(buf).
			H1("Excluded Title").
			H2("Contents").
			TableOfContentsWithRange(markdown.TableOfContentsDepthH2, markdown.TableOfContentsDepthH4).
			H2("Section").
			H3("Subsection").
			H4("Detail").
			H5("Excluded detail").
			Build()
		if err != nil {
			t.Fatalf("Build() = %v, want nil", err)
		}

		if err := golden.Assert("toc_range.md", buf.String()); err != nil {
			t.Error(err)
		}
	})
}

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

// Every defect this package has shipped lived on a boundary between the text a
// caller passes in and the syntax that text is dropped into: a pipe ends a table
// cell, a newline ends a row, a colon ends a YAML key, a backtick run ends a
// fence. The targets in markdown_test.go cover the table and alert boundaries
// from inside the package. The ones here cover the rest through the exported
// API, with a real parser as the oracle rather than a string comparison,
// because the failure mode is always a document that looks fine and means
// something else.

// boundarySeeds are the inputs that have historically broken one boundary or
// another, plus the shapes of text that reach a builder from a database or a
// command's output.
func boundarySeeds() []string {
	return []string{
		"",
		"plain text",
		"a: b",
		"# comment",
		"*ref",
		"&anchor",
		"~",
		"true",
		"123",
		"  leading and trailing  ",
		"---",
		"...",
		"line\nbreak",
		"line\r\nbreak",
		"carriage\rreturn",
		"null\x00byte",
		"bell\x07",
		"tab\tseparated",
		`quote"inside`,
		`back\slash`,
		"emoji 🎉 here",
		"日本語のタイトル",
		"مرحبا بالعالم",
		"```",
		"````",
		"a|b",
		"[link](https://go.dev)",
		strings.Repeat("long ", 2048),
	}
}

// FuzzFrontMatterTitle asserts that a diagram title survives the YAML front
// matter it is written into.
//
// mermaid runs the front matter through a YAML parser before it draws anything,
// so a title is not a string in a template but a value in a document. A bare
// scalar built from arbitrary text can end the value early, turn into a comment,
// resolve to a boolean, or fail to parse at all, and any of those loses the
// whole diagram rather than one line of it.
func FuzzFrontMatterTitle(f *testing.F) {
	for _, seed := range boundarySeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, title string) {
		if strings.TrimSpace(title) == "" {
			// A title that is blank once trimmed writes no front matter at all,
			// which is a different property and is covered by the unit tests.
			return
		}
		diagram := flowchart.NewFlowchart(io.Discard, flowchart.WithTitle(title)).
			NodeWithText("A", "Node A").
			String()

		block, ok := frontMatter(diagram)
		if !ok {
			t.Fatalf("the diagram built with title %q has no front matter block:\n%s", title, diagram)
		}

		var parsed struct {
			Title string `yaml:"title"`
		}
		if err := yaml.Unmarshal([]byte(block), &parsed); err != nil {
			t.Fatalf("the front matter of the diagram built with title %q is not valid YAML: %v\n%s", title, err, block)
		}
		if want := asYAMLReadsIt(title); parsed.Title != want {
			t.Errorf("the front matter title parsed as %q, want %q", parsed.Title, want)
		}
	})
}

// asYAMLReadsIt returns the title a YAML parser reads back out of the front
// matter this library writes for it.
//
// It is the identity for every title that is valid UTF-8, which is every title
// a caller has. Below that the two disagree in one specific way, documented in
// SPEC.md: strconv.Quote writes a byte no valid UTF-8 sequence covers as \xNN,
// meaning that byte, and YAML reads \xNN as the code point U+00NN. Spelling the
// disagreement out here rather than skipping those inputs is what keeps the
// fuzzer able to find a second one.
func asYAMLReadsIt(title string) string {
	if utf8.ValidString(title) {
		return title
	}

	var b strings.Builder
	for i := 0; i < len(title); {
		r, size := utf8.DecodeRuneInString(title[i:])
		if r == utf8.RuneError && size == 1 {
			b.WriteRune(rune(title[i]))
			i++
			continue
		}
		b.WriteString(title[i : i+size])
		i += size
	}
	return b.String()
}

// frontMatter returns the body between the first two "---" lines of a diagram.
func frontMatter(diagram string) (string, bool) {
	lines := strings.Split(strings.ReplaceAll(diagram, "\r\n", "\n"), "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return "", false
	}
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			return strings.Join(lines[1:i], "\n"), true
		}
	}
	return "", false
}

// FuzzCodeBlockContent asserts that content placed in a fenced block comes back
// out of a parser unchanged.
//
// The property only holds for content that does not open a fence of its own:
// three backticks at the start of a line end the block, and the rest of the
// document spills onto the page as prose. That case is a documented limitation
// of CodeBlocks rather than something this target can assert, so it is skipped
// here and pinned separately by TestACodeBlockEndsAtAFenceInItsContent.
func FuzzCodeBlockContent(f *testing.F) {
	for _, seed := range boundarySeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, content string) {
		if holdsFence(content) {
			return
		}

		document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
			return m.CodeBlocks(markdown.SyntaxHighlightGo, content)
		})

		root, source := parse(t, document)
		blocks := nodesOfKind(root, ast.KindFencedCodeBlock)
		if len(blocks) != 1 {
			t.Fatalf("content %q produced %d fenced code blocks, want 1:\n%s", content, len(blocks), document)
		}

		block, ok := blocks[0].(*ast.FencedCodeBlock)
		if !ok {
			t.Fatalf("the parsed node is %T, want *ast.FencedCodeBlock", blocks[0])
		}

		var got strings.Builder
		for i := 0; i < block.Lines().Len(); i++ {
			line := block.Lines().At(i)
			got.Write(line.Value(source))
		}

		// The builder writes a line ending after the content, and a parser
		// reports the block's lines with theirs, so compare without the one at
		// the end. Content whose own last line is empty is indistinguishable
		// from that, which is why both sides are trimmed rather than one.
		if want, have := normalizeContent(content), normalizeContent(got.String()); want != have {
			t.Errorf("content %q came back as %q", want, have)
		}
	})
}

// holdsFence reports whether content has a line that opens or closes a fence.
func holdsFence(content string) bool {
	for _, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		if strings.HasPrefix(strings.TrimLeft(line, " "), "```") {
			return true
		}
	}
	return false
}

// normalizeContent puts code block content in the one shape a round trip can be
// compared in: every line ending resolved to "\n", and the trailing one dropped.
//
// A parser reads "\r", "\n", and "\r\n" as the same thing, so a round trip
// through a document cannot preserve which one the caller wrote. Nor can a
// document tell content ending in line breaks apart from content that does not,
// since the closing fence has to start on its own line either way, which is why
// every trailing one is dropped rather than a single one.
func normalizeContent(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.TrimRight(s, "\n")
}

// TestACodeBlockEndsAtAFenceInItsContent pins the limitation FuzzCodeBlockContent
// skips. CodeBlocks writes a three backtick fence, so content holding a line of
// three backticks closes the block early and everything after it renders as
// prose. Callers whose content can hold a fence have to say so themselves; this
// test exists so the limitation is a stated one rather than a surprise.
func TestACodeBlockEndsAtAFenceInItsContent(t *testing.T) {
	t.Parallel()

	document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
		return m.CodeBlocks(markdown.SyntaxHighlightMermaid, "before\n```\nafter")
	})

	root, _ := parse(t, document)
	if got := len(nodesOfKind(root, ast.KindFencedCodeBlock)); got < 2 {
		t.Errorf("parsed %d fenced code blocks, want at least 2: the limitation this test documents no longer holds", got)
	}
}

// FuzzHeadingText asserts that a heading built from arbitrary text stays one
// heading, for the text a heading can hold.
//
// A heading is a single line by construction, so text carrying a line break is
// out of scope here: markdown has no multi-line heading to put it in. The rest
// of the punctuation that means something at the start of a line, "#" above all,
// has to be harmless.
func FuzzHeadingText(f *testing.F) {
	for _, seed := range boundarySeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		if strings.ContainsAny(text, "\r\n") {
			return
		}

		document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
			return m.H2(text)
		})

		root, _ := parse(t, document)
		if got := len(nodesOfKind(root, ast.KindHeading)); got != 1 {
			t.Errorf("H2(%q) produced %d headings, want 1:\n%s", text, got, document)
		}
	})
}

// FuzzLinkAndImageText asserts that a link or an image stays inline.
//
// The destination is not asserted for arbitrary input, because a URL holding a
// space or an unbalanced parenthesis has no inline form and the caller has to
// encode it. What has to hold for any input is that the helper cannot introduce
// block structure: a link built from a database value must not turn into a
// heading, a list, or a fenced block halfway through a paragraph.
func FuzzLinkAndImageText(f *testing.F) {
	for _, seed := range boundarySeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		if strings.ContainsAny(text, "\r\n") {
			// A line break in the text ends the paragraph, which is a property
			// of markdown rather than of these helpers.
			return
		}

		for name, inline := range map[string]string{
			"link":  markdown.Link(text, "https://example.com"),
			"image": markdown.Image(text, "https://example.com/image.png"),
			"url":   markdown.Link("the text", text),
		} {
			document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
				return m.PlainTextf("before %s after", inline)
			})

			root, _ := parse(t, document)
			for _, kind := range []ast.NodeKind{
				ast.KindHeading,
				ast.KindList,
				ast.KindFencedCodeBlock,
				ast.KindBlockquote,
				ast.KindThematicBreak,
			} {
				if got := len(nodesOfKind(root, kind)); got != 0 {
					t.Errorf("%s built from %q produced %d %s nodes, want 0:\n%s", name, text, got, kind, document)
				}
			}
		}
	})
}

// FuzzSafeLinkDestinationRoundTrips asserts the destination itself, for the
// destinations that have an inline form at all.
func FuzzSafeLinkDestinationRoundTrips(f *testing.F) {
	for _, seed := range []string{
		"https://go.dev",
		"https://example.com/page?a=1&b=2",
		"./relative/path.md",
		"#anchor",
		"mailto:someone@example.com",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, url string) {
		if url == "" || !isInlineDestination(url) {
			return
		}

		document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
			return m.PlainText(markdown.Link("the text", url))
		})

		root, _ := parse(t, document)
		links := nodesOfKind(root, ast.KindLink)
		if len(links) != 1 {
			t.Fatalf("url %q produced %d links, want 1:\n%s", url, len(links), document)
		}

		link, ok := links[0].(*ast.Link)
		if !ok {
			t.Fatalf("the parsed node is %T, want *ast.Link", links[0])
		}
		if got := string(link.Destination); got != url {
			t.Errorf("url %q came back as %q", url, got)
		}
	})
}

// isInlineDestination reports whether url can be written between parentheses
// without being encoded first. Whitespace ends the destination, parentheses and
// angle brackets have to balance, and a control character is not a URL at all.
func isInlineDestination(url string) bool {
	if strings.ContainsAny(url, "()<>\\") {
		return false
	}
	for _, r := range url {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return false
		}
	}
	return true
}

// mermaidBuilder names one mermaid subpackage and builds a small diagram with
// it.
//
// The list lives in the root package's tests because what it is used for is a
// property of the pair: a diagram is only ever seen by a reader after it has
// been handed to CodeBlocks, and it is the combination that has to hold. Seven
// lines per subpackage here beats a near-identical test file in each of the
// seventeen.
type mermaidBuilder struct {
	// name is the subpackage the diagram comes from.
	name string
	// build returns the diagram body, as CodeBlocks would receive it.
	build func() string
}

// mermaidBuilders holds one builder per mermaid subpackage. A new subpackage
// belongs here as soon as it exists.
func mermaidBuilders() []mermaidBuilder {
	return []mermaidBuilder{
		{
			name: "arch",
			build: func() string {
				return arch.NewArchitecture(io.Discard).
					Service("api", arch.IconServer, "API").
					String()
			},
		},
		{
			name: "block",
			build: func() string {
				return block.NewDiagram(io.Discard).
					Row(block.Node("a", block.WithNodeLabel("A"))).
					String()
			},
		},
		{
			name: "c4",
			build: func() string {
				return c4.NewDiagram(io.Discard).
					Person("customer", "Customer").
					System("ledger", "Ledger").
					Rel("customer", "ledger", "Uses").
					String()
			},
		},
		{
			name: "class",
			build: func() string {
				return class.NewDiagram(io.Discard).
					Class("Account", class.WithPublicField("string", "id")).
					String()
			},
		},
		{
			name: "er",
			build: func() string {
				return er.NewDiagram(io.Discard).
					NoRelationship(er.NewEntity("teachers", []*er.Attribute{
						{Type: "int", Name: "id", IsPrimaryKey: true, Comment: "Teacher ID"},
					})).
					String()
			},
		},
		{
			name: "flowchart",
			build: func() string {
				return flowchart.NewFlowchart(io.Discard).
					NodeWithText("A", "Node A").
					String()
			},
		},
		{
			name: "gantt",
			build: func() string {
				return gantt.NewChart(io.Discard).
					Section("Planning").
					Task("Design", "2024-01-01", "2d").
					String()
			},
		},
		{
			name: "gitgraph",
			build: func() string {
				return gitgraph.NewDiagram(io.Discard).
					Commit(gitgraph.WithCommitTag("v1.0.0")).
					String()
			},
		},
		{
			name: "kanban",
			build: func() string {
				return kanban.NewDiagram(io.Discard).
					Column("Todo").
					Task("Define scope").
					String()
			},
		},
		{
			name: "mindmap",
			build: func() string {
				return mindmap.NewDiagram(io.Discard).
					Root("Product").
					Child("Market").
					String()
			},
		},
		{
			name: "packet",
			build: func() string {
				return packet.NewDiagram(io.Discard).
					Field(0, 15, "Source Port").
					String()
			},
		},
		{
			name: "piechart",
			build: func() string {
				return piechart.NewPieChart(io.Discard).
					LabelAndIntValue("Go", 120).
					String()
			},
		},
		{
			name: "quadrant",
			build: func() string {
				return quadrant.NewChart(io.Discard).
					XAxis("Low Reach", "High Reach").
					Point("Campaign A", 0.3, 0.6).
					String()
			},
		},
		{
			name: "requirement",
			build: func() string {
				return requirement.NewDiagram(io.Discard).
					Requirement(
						"a requirement",
						requirement.WithID("1"),
						requirement.WithText("the system shall do the thing"),
						requirement.WithRisk(requirement.RiskLow),
						requirement.WithVerifyMethod(requirement.VerifyMethodTest),
					).
					String()
			},
		},
		{
			name: "sequence",
			build: func() string {
				return sequence.NewDiagram(io.Discard).
					SyncRequest("Client", "Server", "GET /users").
					String()
			},
		},
		{
			name: "state",
			build: func() string {
				return state.NewDiagram(io.Discard).
					State("Draft", "The order is being written").
					Transition("Draft", "Placed").
					String()
			},
		},
		{
			name: "userjourney",
			build: func() string {
				return userjourney.NewDiagram(io.Discard).
					Section("Discovery").
					Task("Find the site", userjourney.ScoreSatisfied, "Visitor").
					String()
			},
		},
		{
			name: "venn",
			build: func() string {
				return venn.NewDiagram(io.Discard).
					SetWithLabel("go", "Go").
					SetWithLabel("rust", "Rust").
					String()
			},
		},
		{
			name: "xychart",
			build: func() string {
				return xychart.NewDiagram(io.Discard).
					XAxisLabels("jan", "feb").
					Bar(5000, 6000).
					String()
			},
		},
	}
}

// TestReadmeShowsWhatTheGeneratorsProduce pins the README against the documents
// it links to.
//
// Every mermaid section shows a builder chain, then the document that chain
// produces, then a link to the committed sample under doc/. The three had drifted
// apart: the chains left out WithBlockSpacing while every generator used it, so a
// reader who copied an example and diffed the result against the linked file
// found a difference that nothing explained. Comparing them here is what keeps
// them together, since the samples are regenerated by "make generate" and the
// README is written by hand.
func TestReadmeShowsWhatTheGeneratorsProduce(t *testing.T) {
	t.Parallel()

	readme, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}

	// Normalized because git checks this file out with CRLF line endings on
	// Windows, where every line would otherwise end in a carriage return that
	// neither the pattern below nor the comparison expects.
	sections := readmeSamples(t, strings.ReplaceAll(string(readme), "\r\n", "\n"))
	if len(sections) == 0 {
		t.Fatal("found no linked samples in README.md; has the link wording changed?")
	}

	for name, shown := range sections {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// The path is built from a directory name matched out of the
			// README by [\w-]+, so it cannot climb out of doc/.
			generated, err := os.ReadFile(filepath.Join("doc", name, "generated.md")) //nolint:gosec
			if err != nil {
				t.Fatalf("read the sample the README links to: %v", err)
			}

			want := strings.TrimRight(strings.ReplaceAll(string(generated), "\r\n", "\n"), "\n")
			if shown != want {
				t.Errorf(
					"the README block for %s is not what doc/%s/generated.md holds.\n"+
						"Run make generate, then copy the file into the README block.\n\ngot:\n%s\n\nwant:\n%s",
					name, name, shown, want,
				)
			}
		})
	}
}

// TestEveryExportedSymbolHasAnExample walks this module's packages and reports
// any exported function, or method on an exported type, that godoc would show
// without a runnable example beside it.
//
// The examples are the documentation this library ships: pkg.go.dev puts one
// under the symbol it is named for, and a reader deciding how to call something
// looks there first. Counting them here is what keeps a new API from arriving
// without one, since nothing else notices.
//
// Names have to match Go's rules exactly or godoc silently attaches the example
// to nothing: ExampleMarkdown_H1 for method H1 on Markdown, ExampleBold for the
// function Bold. The list of examples comes from go/doc rather than from a scan
// for the name, so an example that takes an argument, returns a value or has no
// output to check against does not count as documenting anything, which is what
// godoc does with it too.
func TestEveryExportedSymbolHasAnExample(t *testing.T) {
	t.Parallel()

	packages, err := filepath.Glob("mermaid/*")
	if err != nil {
		t.Fatalf("list the mermaid subpackages: %v", err)
	}
	packages = append(packages, ".")

	for _, dir := range packages {
		t.Run(dir, func(t *testing.T) {
			t.Parallel()

			symbols, examples := exportedSymbols(t, dir)
			for _, symbol := range symbols {
				if !examples[symbol] {
					t.Errorf(
						"%s has no example. Add func %s() to %s/examples_test.go, with an // Output: block.",
						symbol, symbol, dir,
					)
				}
			}
		})
	}
}

// readmeSamples returns the document each "Plain text output" link is followed
// by, keyed by the directory under doc/ that it links to.
func readmeSamples(t *testing.T, readme string) map[string]string {
	t.Helper()

	link := regexp.MustCompile(`^Plain text output: \[markdown is here\]\(\./doc/([\w-]+)/generated\.md\)$`)
	samples := make(map[string]string)

	lines := strings.Split(readme, "\n")
	for i := 0; i < len(lines); i++ {
		m := link.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}

		// The sample follows the link inside a fence of four or more backticks,
		// long enough to hold the ```mermaid fence the document itself opens.
		if i+1 >= len(lines) || !strings.HasPrefix(lines[i+1], "````") {
			t.Errorf("the %s link is not followed by a fenced block", m[1])
			continue
		}
		close := strings.Repeat("`", len(lines[i+1])-len(strings.TrimLeft(lines[i+1], "`")))

		body := []string{}
		i += 2
		for i < len(lines) && lines[i] != close {
			body = append(body, lines[i])
			i++
		}
		if i >= len(lines) {
			t.Errorf("the block after the %s link is never closed", m[1])
			continue
		}
		samples[m[1]] = strings.Join(body, "\n")
	}
	return samples
}

// exportedSymbols returns the example name every exported symbol of the package
// in dir would be documented under, and the examples the package actually has.
//
// Types are in, because pkg.go.dev puts an example under a type and a reader
// looking at TableSet wants to see one filled in. Constants and variables are
// not: there are sixty one of them here, they are the enumerations a method
// takes rather than anything a caller calls, and an example under each would
// bury the page rather than document it.
func exportedSymbols(t *testing.T, dir string) ([]string, map[string]bool) {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("parse %s: %v", dir, err)
	}

	want := []string{}
	examples := map[string]bool{}
	for _, pkg := range pkgs {
		for name, file := range pkg.Files {
			if strings.HasSuffix(name, "_test.go") {
				// go/doc reports the examples godoc will actually show: it
				// leaves out anything that takes an argument, returns a value,
				// or has no output to check against, which a scan for the name
				// alone would count.
				for _, example := range doc.Examples(file) {
					if example.Output != "" {
						examples["Example"+example.Name] = true
					}
				}
				continue
			}
			want = append(want, documentedSymbols(file)...)
		}
	}
	return want, examples
}

// documentedSymbols returns the example name each exported function, method on
// an exported type, and type in file would be documented under.
func documentedSymbols(file *goast.File) []string {
	symbols := []string{}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *goast.FuncDecl:
			if !d.Name.IsExported() {
				continue
			}
			if symbol, ok := documentedAs(d); ok {
				symbols = append(symbols, symbol)
			}
		case *goast.GenDecl:
			if d.Tok != token.TYPE {
				continue
			}
			for _, spec := range d.Specs {
				if ts, ok := spec.(*goast.TypeSpec); ok && ts.Name.IsExported() {
					symbols = append(symbols, "Example"+ts.Name.Name)
				}
			}
		}
	}
	return symbols
}

// documentedAs returns the example name godoc attaches to fn.
func documentedAs(fn *goast.FuncDecl) (string, bool) {
	if fn.Recv == nil {
		return "Example" + fn.Name.Name, true
	}
	if len(fn.Recv.List) != 1 {
		return "", false
	}

	receiver := fn.Recv.List[0].Type
	if star, ok := receiver.(*goast.StarExpr); ok {
		receiver = star.X
	}
	ident, ok := receiver.(*goast.Ident)
	if !ok || !ident.IsExported() {
		return "", false
	}
	return "Example" + ident.Name + "_" + fn.Name.Name, true
}

// TestABuilderBelongsToOneGoroutine pins the concurrency contract SPEC.md
// states: a builder is not safe for concurrent use, and the way to build two
// documents at once is to build two builders.
//
// There is no test here that shares one builder between goroutines, because
// such a test is a data race: it would pass or fail depending on the scheduler,
// and under -race it would fail on purpose, which is not a test but a promise
// to break the build. What is pinned instead is the shape the contract points
// callers at, which has to keep working.
func TestABuilderBelongsToOneGoroutine(t *testing.T) {
	t.Parallel()

	const writers = 8

	documents := make([]string, writers)
	var wg sync.WaitGroup
	for i := range writers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// One builder per goroutine, which is what the contract asks for.
			documents[i] = markdown.NewMarkdown(io.Discard).
				H2(fmt.Sprintf("Section %d", i)).
				PlainText("Built on its own goroutine.").
				String()
		}()
	}
	wg.Wait()

	for i, document := range documents {
		if want := fmt.Sprintf("## Section %d", i); !strings.Contains(document, want) {
			t.Errorf("document %d = %q, want it to contain %q", i, document, want)
		}
	}
}
