package markdown

import (
	"bytes"
	"strings"
	"testing"

	"github.com/nao1215/markdown/internal"
)

// TestNormalizeLineFeeds covers the conversion applied to text that reaches the
// builder from elsewhere, such as a table rendered by tablewriter.
func TestNormalizeLineFeeds(t *testing.T) {
	t.Parallel()

	lineFeed := internal.LineFeed()

	tests := map[string]struct {
		in   string
		want string
	}{
		"unix input":       {in: "a\nb", want: "a" + lineFeed + "b"},
		"windows input":    {in: "a\r\nb", want: "a" + lineFeed + "b"},
		"mixed input":      {in: "a\r\nb\nc", want: "a" + lineFeed + "b" + lineFeed + "c"},
		"no line feeds":    {in: "abc", want: "abc"},
		"empty":            {in: "", want: ""},
		"trailing newline": {in: "a\n", want: "a" + lineFeed},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := normalizeLineFeeds(tt.in); got != tt.want {
				t.Errorf("normalizeLineFeeds(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestDocumentUsesOneLineEnding builds a document that exercises every
// construct and asserts it contains no line ending other than the platform one.
//
// This is a Windows regression test. It is trivially true on Linux and macOS,
// where the platform line feed is "\n", but on Windows it catches any construct
// that emits a bare "\n" and leaves the document with mixed endings. That is
// exactly how CustomTable shipped alignment that silently did nothing there:
// tablewriter separates rows with "\n" on every platform.
func TestDocumentUsesOneLineEnding(t *testing.T) {
	t.Parallel()

	lineFeed := internal.LineFeed()

	var buf bytes.Buffer
	err := NewMarkdown(&buf).
		H1("Title").
		TableOfContents(TableOfContentsDepthH3).
		PlainText("paragraph").
		H2("Lists").
		BulletList("alpha", "beta").
		OrderedList("one", "two").
		CheckBox([]CheckBoxSet{{Text: "todo"}, {Checked: true, Text: "done"}}).
		H2("Quotes").
		Blockquote("quoted"+lineFeed+"over two lines").
		Note("note"+lineFeed+"over two lines").
		Warning("warning").
		H2("Blocks").
		CodeBlocks(SyntaxHighlightGo, "x := 1").
		Details("summary", "body").
		HorizontalRule().
		H2("Tables").
		Table(TableSet{
			Header:    []string{"name", "value"},
			Rows:      [][]string{{"a", "1"}},
			Alignment: []TableAlignment{AlignLeft, AlignRight},
		}).
		CustomTable(TableSet{
			Header:    []string{"name", "value"},
			Rows:      [][]string{{"a", "1"}},
			Alignment: []TableAlignment{AlignLeft, AlignRight},
		}, TableOptions{}).
		LF().
		Build()
	if err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	if lineFeed == "\n" {
		if strings.Contains(buf.String(), "\r") {
			t.Errorf("document contains a carriage return:\n%q", buf.String())
		}
		return
	}

	// On Windows every "\n" has to be preceded by "\r".
	stripped := strings.ReplaceAll(buf.String(), lineFeed, "")
	if strings.ContainsAny(stripped, "\r\n") {
		t.Errorf("document mixes line endings:\n%q", buf.String())
	}
}

// TestCustomTableAlignmentSurvivesLineFeedConversion pins the specific defect:
// the delimiter row rewrite splits the rendered table on the platform line
// feed, so it found nothing to rewrite while tablewriter was emitting "\n".
func TestCustomTableAlignmentSurvivesLineFeedConversion(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := NewMarkdown(&buf).CustomTable(TableSet{
		Header:    []string{"name", "value"},
		Rows:      [][]string{{"a", "1"}},
		Alignment: []TableAlignment{AlignLeft, AlignRight},
	}, TableOptions{}).Build(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	lines := strings.Split(buf.String(), internal.LineFeed())
	if len(lines) < 2 {
		t.Fatalf("table was not split into rows: %q", buf.String())
	}
	if !strings.HasPrefix(lines[1], "|:") || !strings.HasSuffix(strings.TrimSuffix(lines[1], "|"), ":") {
		t.Errorf("alignment markers missing from the delimiter row: %q", lines[1])
	}
}
