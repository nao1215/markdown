package markdown

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal"
)

// lf is the platform line feed, spelled short because these tests are mostly
// exact-output comparisons.
func lf() string { return internal.LineFeed() }

// TestAlertSingleLineIsUnchanged is the compatibility guard for the alert
// rewrite: single-line text, which is what nearly every caller passes, has to
// produce exactly what it produced before.
func TestAlertSingleLineIsUnchanged(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		call func(*Markdown) *Markdown
		want string
	}{
		"note":      {func(m *Markdown) *Markdown { return m.Note("hello") }, "> [!NOTE]  " + lf() + "> hello"},
		"tip":       {func(m *Markdown) *Markdown { return m.Tip("hello") }, "> [!TIP]  " + lf() + "> hello"},
		"important": {func(m *Markdown) *Markdown { return m.Important("hello") }, "> [!IMPORTANT]  " + lf() + "> hello"},
		"warning":   {func(m *Markdown) *Markdown { return m.Warning("hello") }, "> [!WARNING]  " + lf() + "> hello"},
		"caution":   {func(m *Markdown) *Markdown { return m.Caution("hello") }, "> [!CAUTION]  " + lf() + "> hello"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := tt.call(NewMarkdown(nil)).String(); got != tt.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}

// TestAlertQuotesEveryLine covers the defect: only the first line carried the
// quote marker, so a list or a blank line in the text escaped the callout.
func TestAlertQuotesEveryLine(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		text string
		want []string
	}{
		"prose over two lines": {
			text: "first\nsecond",
			want: []string{"> [!NOTE]  ", "> first", "> second"},
		},
		"list in the body": {
			text: "intro\n- alpha\n- beta",
			want: []string{"> [!NOTE]  ", "> intro", "> - alpha", "> - beta"},
		},
		"blank line stays inside the callout": {
			text: "first\n\nsecond",
			want: []string{"> [!NOTE]  ", "> first", ">", "> second"},
		},
		"windows line endings in the text": {
			text: "first\r\nsecond",
			want: []string{"> [!NOTE]  ", "> first", "> second"},
		},
		"lines the caller already quoted are left alone": {
			// Callers hand-wrote these continuations to work around the bug.
			// Prefixing them again would nest a quote inside the alert.
			text: "first  \n> second  \n> * third",
			want: []string{"> [!NOTE]  ", "> first  ", "> second  ", "> * third"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := strings.Split(NewMarkdown(nil).Note(tt.text).String(), lf())
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("alert body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestBlockquoteSplitsOnPlainNewline covers the Windows case: splitting on the
// platform line feed meant a plain Go literal was never split at all.
func TestBlockquoteSplitsOnPlainNewline(t *testing.T) {
	t.Parallel()

	for name, text := range map[string]string{
		"unix":    "first\nsecond",
		"windows": "first\r\nsecond",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			want := "> first" + lf() + "> second"
			if got := NewMarkdown(nil).Blockquote(text).String(); got != want {
				t.Errorf("blockquote mismatch:\n got: %q\nwant: %q", got, want)
			}
		})
	}
}

// TestDetailsSurroundsBodyWithBlankLines covers the reason nested markdown used
// to render as literal text: an HTML block runs to the next blank line.
func TestDetailsSurroundsBodyWithBlankLines(t *testing.T) {
	t.Parallel()

	// The body is written by the caller and is passed through untouched, so it
	// has to be built with the platform line feed for this comparison.
	body := "- alpha" + lf() + "- beta"
	got := NewMarkdown(nil).Details("summary", body).String()
	want := strings.Join([]string{
		"<details>",
		"<summary>summary</summary>",
		"",
		"- alpha",
		"- beta",
		"",
		"</details>",
		"",
	}, lf())

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("details mismatch (-want +got):\n%s", diff)
	}
}

// TestDetailsDoesNotSwallowTheNextBlock is the other half of the same defect:
// without the trailing blank line the block after </details> stayed inside the
// HTML block and rendered as text.
func TestDetailsDoesNotSwallowTheNextBlock(t *testing.T) {
	t.Parallel()

	got := NewMarkdown(nil).
		Details("summary", "body").
		CodeBlocks(SyntaxHighlightGo, "x := 1").
		String()

	if !strings.Contains(got, "</details>"+lf()+lf()+"```go") {
		t.Errorf("no blank line between </details> and the next block:\n%q", got)
	}
}

// TestBuildTerminatesWithNewline pins MD047 and the append case: writing a
// second document to the same writer used to splice it onto the last line.
func TestBuildTerminatesWithNewline(t *testing.T) {
	t.Parallel()

	t.Run("single document", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		if err := NewMarkdown(&buf).H1("Title").Build(); err != nil {
			t.Fatalf("failed to build: %v", err)
		}
		if want := "# Title" + lf(); buf.String() != want {
			t.Errorf("got %q, want %q", buf.String(), want)
		}
	})

	t.Run("content already ending in a line feed gets no second one", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		if err := NewMarkdown(&buf).Table(TableSet{
			Header: []string{"a"},
			Rows:   [][]string{{"1"}},
		}).Build(); err != nil {
			t.Fatalf("failed to build: %v", err)
		}
		if strings.HasSuffix(buf.String(), lf()+lf()+lf()) {
			t.Errorf("trailing newlines piled up: %q", buf.String())
		}
	})

	t.Run("two documents to one writer stay on separate lines", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		if err := NewMarkdown(&buf).H1("First").Build(); err != nil {
			t.Fatalf("failed to build: %v", err)
		}
		if err := NewMarkdown(&buf).H1("Second").Build(); err != nil {
			t.Fatalf("failed to build: %v", err)
		}
		want := "# First" + lf() + "# Second" + lf()
		if buf.String() != want {
			t.Errorf("got %q, want %q", buf.String(), want)
		}
	})

	t.Run("String is left alone", func(t *testing.T) {
		t.Parallel()

		if got := NewMarkdown(nil).H1("Title").String(); got != "# Title" {
			t.Errorf("String() gained a line feed: %q", got)
		}
	})
}

// TestBlankLineBetweenBlocks covers the separations markdown requires. The
// cases that gain a blank line render wrongly without one; the rest must stay
// byte for byte as they were.
func TestBlankLineBetweenBlocks(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(*Markdown) *Markdown
		want  []string
	}{
		"list then table": {
			build: func(m *Markdown) *Markdown {
				return m.BulletList("alpha").Table(TableSet{Header: []string{"a"}, Rows: [][]string{{"1"}}})
			},
			want: []string{"- alpha", "", "| a |", "|---------|", "| 1 |", ""},
		},
		"list then paragraph": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha").PlainText("after") },
			want:  []string{"- alpha", "", "after"},
		},
		"list then heading": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha").H2("Next") },
			want:  []string{"- alpha", "", "## Next"},
		},
		"consecutive list items stay tight": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha", "beta") },
			want:  []string{"- alpha", "- beta"},
		},
		"ordered list items stay tight": {
			build: func(m *Markdown) *Markdown { return m.OrderedList("alpha", "beta") },
			want:  []string{"1. alpha", "2. beta"},
		},
		"a different list kind is separated": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha").OrderedList("one") },
			want:  []string{"- alpha", "", "1. one"},
		},
		"checkbox items stay tight": {
			build: func(m *Markdown) *Markdown {
				return m.CheckBox([]CheckBoxSet{{Text: "alpha"}, {Checked: true, Text: "beta"}})
			},
			want: []string{"- [ ] alpha", "- [x] beta"},
		},
		"alert then paragraph": {
			build: func(m *Markdown) *Markdown { return m.Note("careful").PlainText("after") },
			want:  []string{"> [!NOTE]  ", "> careful", "", "after"},
		},
		"two alerts do not merge": {
			build: func(m *Markdown) *Markdown { return m.Note("first").Warning("second") },
			want:  []string{"> [!NOTE]  ", "> first", "", "> [!WARNING]  ", "> second"},
		},
		"blockquote then paragraph": {
			build: func(m *Markdown) *Markdown { return m.Blockquote("quoted").PlainText("after") },
			want:  []string{"> quoted", "", "after"},
		},
		"an existing LF spacer is not doubled": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha").LF().PlainText("after") },
			want:  []string{"- alpha", "  ", "after"},
		},
		"an existing empty paragraph is not doubled": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha").PlainText("").PlainText("after") },
			want:  []string{"- alpha", "", "after"},
		},
		"headings and paragraphs are untouched": {
			build: func(m *Markdown) *Markdown { return m.H1("Title").PlainText("text").H2("Section") },
			want:  []string{"# Title", "text", "## Section"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := strings.Split(tt.build(NewMarkdown(nil)).String(), lf())
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("block layout mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestTableOfContentsMarkersStayTight makes sure the block separation does not
// push a blank line into the marker pair, and that the entries still land
// between the markers. Matching the marker pair as one string used to be the
// mechanism, and it broke silently the moment anything separated them.
func TestTableOfContentsMarkersStayTight(t *testing.T) {
	t.Parallel()

	got := NewMarkdown(nil).
		H1("Doc").
		TableOfContents(TableOfContentsDepthH2).
		H2("Section").
		String()

	want := strings.Join([]string{
		"# Doc",
		TableOfContentsMarkerBegin,
		"- [Doc](#doc)",
		"  - [Section](#section)",
		TableOfContentsMarkerEnd,
		"",
		"## Section",
	}, lf())

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("table of contents mismatch (-want +got):\n%s", diff)
	}
}

// TestInsertTableOfContentsWithMissingMarkers guards the splice against a body
// that does not hold both markers.
func TestInsertTableOfContentsWithMissingMarkers(t *testing.T) {
	t.Parallel()

	toc := []string{"- [A](#a)"}

	tests := map[string][]string{
		"no markers at all": {"# Title"},
		"only the opening":  {TableOfContentsMarkerBegin, "# Title"},
		"only the closing":  {"# Title", TableOfContentsMarkerEnd},
		"reversed order":    {TableOfContentsMarkerEnd, TableOfContentsMarkerBegin},
	}

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if diff := cmp.Diff(body, insertTableOfContents(body, toc)); diff != "" {
				t.Errorf("body was modified (-want +got):\n%s", diff)
			}
		})
	}
}

// TestWithBlockSpacing covers the opt-in that separates every block, which is
// what markdownlint and mkdocs want. The default is deliberately tighter, so
// both shapes are pinned here.
func TestWithBlockSpacing(t *testing.T) {
	t.Parallel()

	chain := func(m *Markdown) *Markdown {
		return m.H1("Title").
			PlainText("paragraph").
			BulletList("alpha", "beta").
			CodeBlocks(SyntaxHighlightGo, "x := 1")
	}

	compact := chain(NewMarkdown(nil)).String()
	spaced := chain(NewMarkdown(nil, WithBlockSpacing())).String()

	wantCompact := strings.Join([]string{
		"# Title",
		"paragraph",
		"- alpha",
		"- beta",
		"",
		"```go",
		"x := 1",
		"```",
	}, lf())
	if diff := cmp.Diff(wantCompact, compact); diff != "" {
		t.Errorf("default output changed (-want +got):\n%s", diff)
	}

	wantSpaced := strings.Join([]string{
		"# Title",
		"",
		"paragraph",
		"",
		"- alpha",
		"- beta",
		"",
		"```go",
		"x := 1",
		"```",
	}, lf())
	if diff := cmp.Diff(wantSpaced, spaced); diff != "" {
		t.Errorf("spaced output mismatch (-want +got):\n%s", diff)
	}
}

// TestBlankLine covers the explicit spacer, and that it is not doubled by the
// automatic separation.
func TestBlankLine(t *testing.T) {
	t.Parallel()

	got := NewMarkdown(nil).BulletList("alpha").BlankLine().PlainText("after").String()
	want := strings.Join([]string{"- alpha", "", "after"}, lf())

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("blank line mismatch (-want +got):\n%s", diff)
	}
}
