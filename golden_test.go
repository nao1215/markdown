package markdown_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/internal/golden"
)

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
