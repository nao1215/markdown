package markdown

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal"
)

// delimiterRow returns the alignment row of the first table in the rendered
// document, which is where markdown records column alignment.
func delimiterRow(t *testing.T, rendered string) string {
	t.Helper()

	lines := strings.Split(rendered, internal.LineFeed())
	if len(lines) < 2 {
		t.Fatalf("rendered table has %d line(s), want at least 2:\n%s", len(lines), rendered)
	}
	return lines[1]
}

// TestCustomTableHonoursAlignment pins the behaviour that CustomTable used to
// drop on the floor: TableSet.Alignment is a documented public field, but it
// never reached the rendered output, so choosing CustomTable silently cost you
// alignment.
func TestCustomTableHonoursAlignment(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		alignment []TableAlignment
		want      string
	}{
		"left, centre, right": {
			alignment: []TableAlignment{AlignLeft, AlignCenter, AlignRight},
			want:      "|:----|:------:|-----:|",
		},
		"default keeps a plain rule": {
			alignment: []TableAlignment{AlignDefault, AlignDefault, AlignDefault},
			want:      "|-----|--------|------|",
		},
		"shorter than the header falls back to default": {
			alignment: []TableAlignment{AlignRight},
			want:      "|----:|--------|------|",
		},
		"longer than the header ignores the extra entries": {
			alignment: []TableAlignment{AlignLeft, AlignLeft, AlignLeft, AlignRight, AlignCenter},
			want:      "|:----|:-------|:-----|",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			if err := NewMarkdown(&buf).CustomTable(TableSet{
				Header:    []string{"key", "middle", "last"},
				Rows:      [][]string{{"a", "b", "c"}},
				Alignment: tt.alignment,
			}, TableOptions{}).Build(); err != nil {
				t.Fatalf("failed to build: %v", err)
			}

			if got := delimiterRow(t, buf.String()); got != tt.want {
				t.Errorf("delimiter row mismatch:\n got: %q\nwant: %q\nfull output:\n%s", got, tt.want, buf.String())
			}
		})
	}
}

// TestCustomTableWithoutAlignmentIsUnchanged is the compatibility guard: a
// TableSet with no Alignment must render exactly as it did before alignment
// support was wired in, so existing golden files keep passing.
func TestCustomTableWithoutAlignmentIsUnchanged(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := NewMarkdown(&buf).CustomTable(TableSet{
		Header: []string{"key", "middle", "last"},
		Rows:   [][]string{{"a", "b", "c"}, {"longer", "value", "here"}},
	}, TableOptions{}).Build(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	want := strings.Join([]string{
		"|  key   | middle | last |",
		"|--------|--------|------|",
		"| a      | b      | c    |",
		"| longer | value  | here |",
		"",
	}, internal.LineFeed())

	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("output changed (-want +got):\n%s", diff)
	}
}

// TestTableAndCustomTableAgreeOnAlignment keeps the two table renderers from
// drifting apart again: whatever alignment a caller asks for, both must place
// the colons on the same sides.
func TestTableAndCustomTableAgreeOnAlignment(t *testing.T) {
	t.Parallel()

	set := TableSet{
		Header:    []string{"key", "middle", "last"},
		Rows:      [][]string{{"a", "b", "c"}},
		Alignment: []TableAlignment{AlignRight, AlignCenter, AlignLeft},
	}

	var plain, custom bytes.Buffer
	if err := NewMarkdown(&plain).Table(set).Build(); err != nil {
		t.Fatalf("failed to build Table: %v", err)
	}
	if err := NewMarkdown(&custom).CustomTable(set, TableOptions{}).Build(); err != nil {
		t.Fatalf("failed to build CustomTable: %v", err)
	}

	sides := func(row string) []string {
		cells := strings.Split(strings.Trim(row, "|"), "|")
		out := make([]string, 0, len(cells))
		for _, cell := range cells {
			switch {
			case strings.HasPrefix(cell, ":") && strings.HasSuffix(cell, ":"):
				out = append(out, "center")
			case strings.HasPrefix(cell, ":"):
				out = append(out, "left")
			case strings.HasSuffix(cell, ":"):
				out = append(out, "right")
			default:
				out = append(out, "default")
			}
		}
		return out
	}

	if diff := cmp.Diff(sides(delimiterRow(t, plain.String())), sides(delimiterRow(t, custom.String()))); diff != "" {
		t.Errorf("Table and CustomTable disagree on alignment (-Table +CustomTable):\n%s", diff)
	}
}

// TestDelimiterCell covers the widths where a marker does not fit, which the
// table renderers can reach for a very narrow column.
func TestDelimiterCell(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		width int
		align TableAlignment
		want  string
	}{
		"left needs two characters":    {width: 1, align: AlignLeft, want: "-"},
		"left at the minimum":          {width: 2, align: AlignLeft, want: ":-"},
		"right needs two characters":   {width: 1, align: AlignRight, want: "-"},
		"right at the minimum":         {width: 2, align: AlignRight, want: "-:"},
		"centre needs three":           {width: 2, align: AlignCenter, want: "--"},
		"centre at the minimum":        {width: 3, align: AlignCenter, want: ":-:"},
		"default is always all dashes": {width: 4, align: AlignDefault, want: "----"},
		"zero width":                   {width: 0, align: AlignCenter, want: ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := delimiterCell(tt.width, tt.align); got != tt.want {
				t.Errorf("delimiterCell(%d, %v) = %q, want %q", tt.width, tt.align, got, tt.want)
			}
		})
	}
}

// TestApplyAlignmentToDelimiterRowLeavesOddInputAlone makes sure the rewrite
// never corrupts something that is not a delimiter row.
func TestApplyAlignmentToDelimiterRowLeavesOddInputAlone(t *testing.T) {
	t.Parallel()

	alignment := []TableAlignment{AlignCenter}

	tests := map[string]string{
		"no alignment requested":    "| a |" + internal.LineFeed() + "|---|",
		"single line":               "| a |",
		"second line is not a rule": "| a |" + internal.LineFeed() + "| b |",
		"empty":                     "",
	}

	for name, rendered := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			use := alignment
			if name == "no alignment requested" {
				use = nil
			}
			if got := applyAlignmentToDelimiterRow(rendered, use); got != rendered {
				t.Errorf("input was rewritten:\n got: %q\nwant: %q", got, rendered)
			}
		})
	}
}
