package markdown

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal"
)

// TestEscapeTableCell covers the characters that end a cell or a row, and the
// idempotence the function promises.
func TestEscapeTableCell(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   string
		want string
	}{
		"plain text is returned unchanged":   {in: "hello", want: "hello"},
		"empty":                              {in: "", want: ""},
		"pipe":                               {in: "a|b", want: `a\|b`},
		"several pipes":                      {in: "a|b|c", want: `a\|b\|c`},
		"pipe at the edges":                  {in: "|a|", want: `\|a\|`},
		"already escaped pipe is left alone": {in: `a\|b`, want: `a\|b`},
		"escaped backslash then pipe":        {in: `a\\|b`, want: `a\\\|b`},
		"newline becomes a break":            {in: "a\nb", want: "a<br>b"},
		"carriage return and newline":        {in: "a\r\nb", want: "a<br>b"},
		"lone carriage return":               {in: "a\rb", want: "a<br>b"},
		"trailing backslash":                 {in: `a\`, want: `a\`},
		"markup the caller built survives":   {in: "[Go](https://go.dev)", want: "[Go](https://go.dev)"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := EscapeTableCell(tt.in)
			if got != tt.want {
				t.Errorf("EscapeTableCell(%q) = %q, want %q", tt.in, got, tt.want)
			}
			if again := EscapeTableCell(got); again != got {
				t.Errorf("not idempotent: %q became %q", got, again)
			}
		})
	}
}

// TestTableEscapeCellsKeepsColumnCount is the point of the option: an unescaped
// pipe splits the cell in two, so GitHub renders the row with the wrong number
// of columns and drops the last value.
func TestTableEscapeCellsKeepsColumnCount(t *testing.T) {
	t.Parallel()

	for name, render := range map[string]func(*Markdown, TableSet) *Markdown{
		"Table":       func(m *Markdown, ts TableSet) *Markdown { return m.Table(ts) },
		"CustomTable": func(m *Markdown, ts TableSet) *Markdown { return m.CustomTable(ts, TableOptions{}) },
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			if err := render(NewMarkdown(&buf), TableSet{
				Header:      []string{"cmd", "desc"},
				Rows:        [][]string{{"a|b", "keeps its value"}},
				EscapeCells: true,
			}).Build(); err != nil {
				t.Fatalf("failed to build: %v", err)
			}

			lines := strings.Split(strings.TrimRight(buf.String(), internal.LineFeed()), internal.LineFeed())
			row := lines[len(lines)-1]
			if cells := unescapedPipeCount(row); cells != 3 {
				t.Errorf("row has %d unescaped pipes, want 3 (one per edge plus one separator): %q", cells, row)
			}
			if !strings.Contains(row, "keeps its value") {
				t.Errorf("the last column was dropped: %q", row)
			}
		})
	}
}

// unescapedPipeCount counts the pipes that actually delimit cells.
func unescapedPipeCount(row string) int {
	count := 0
	for i := 0; i < len(row); i++ {
		if row[i] == '\\' {
			i++
			continue
		}
		if row[i] == '|' {
			count++
		}
	}
	return count
}

// TestTableWithoutEscapeCellsIsUnchanged is the compatibility guard: the option
// is off by default and the rendering must not move.
func TestTableWithoutEscapeCellsIsUnchanged(t *testing.T) {
	t.Parallel()

	set := TableSet{
		Header: []string{"cmd", "desc"},
		Rows:   [][]string{{"a|b", "text"}},
	}

	var buf bytes.Buffer
	if err := NewMarkdown(&buf).Table(set).Build(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	want := strings.Join([]string{
		"| cmd | desc |",
		"|---------|---------|",
		"| a|b | text |",
		"",
	}, internal.LineFeed())

	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("output changed (-want +got):\n%s", diff)
	}
}

// TestEscapeCellsDoesNotMutateTheCallersSlices makes sure enabling the option
// does not write back into the data the caller passed in.
func TestEscapeCellsDoesNotMutateTheCallersSlices(t *testing.T) {
	t.Parallel()

	header := []string{"a|b"}
	rows := [][]string{{"c|d"}}

	var buf bytes.Buffer
	if err := NewMarkdown(&buf).Table(TableSet{
		Header:      header,
		Rows:        rows,
		EscapeCells: true,
	}).Build(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	if header[0] != "a|b" || rows[0][0] != "c|d" {
		t.Errorf("caller data was modified: header=%v rows=%v", header, rows)
	}
}
