package markdown

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/nao1215/markdown/internal"
)

// Every defect found in the recent downstream survey was a string-handling
// defect: a pipe in a cell dropped a column, a newline in alert text escaped the
// callout, a heading with punctuation produced an anchor nobody checked. The
// table-driven tests only cover the inputs someone thought of, so these targets
// state the properties instead and let the fuzzer look for the counterexample.

// FuzzEscapeTableCell asserts the two properties the escaper promises: the
// result never contains a character that ends a cell or a row, and escaping
// twice changes nothing.
func FuzzEscapeTableCell(f *testing.F) {
	for _, seed := range []string{
		"", "plain", "a|b", "a|b|c", `a\|b`, `a\\|b`, "a\nb", "a\r\nb", "a\rb",
		`a\`, "[Go](https://go.dev)", "**bold**", "|", "||", "\n", "\r",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, in string) {
		got := EscapeTableCell(in)

		if strings.ContainsAny(got, "\r\n") {
			t.Errorf("EscapeTableCell(%q) = %q, which still ends the row", in, got)
		}
		if unescapedPipeCount(got) != 0 {
			t.Errorf("EscapeTableCell(%q) = %q, which still ends the cell", in, got)
		}
		if again := EscapeTableCell(got); again != got {
			t.Errorf("EscapeTableCell is not idempotent: %q became %q", got, again)
		}
	})
}

// FuzzTableRowKeepsColumnCount asserts the property the escaper exists for: for
// any cell content, an escaped row still has exactly as many cells as the
// header. This is the property the unescaped path violates.
func FuzzTableRowKeepsColumnCount(f *testing.F) {
	for _, seed := range []string{"a", "a|b", "x\ny", "|", `\|`, ""} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, cell string) {
		out := NewMarkdown(nil).Table(TableSet{
			Header:      []string{"one", "two"},
			Rows:        [][]string{{cell, "second"}},
			EscapeCells: true,
		}).String()

		lines := strings.Split(strings.TrimRight(out, internal.LineFeed()), internal.LineFeed())
		row := lines[len(lines)-1]

		// Two edges plus one separator between the cells.
		if got := unescapedPipeCount(row); got != 3 {
			t.Errorf("row %q built from cell %q has %d delimiters, want 3", row, cell, got)
		}
	})
}

// FuzzAlertQuotesEveryLine asserts that no line of an alert can escape the
// callout, whatever the caller puts in the text.
func FuzzAlertQuotesEveryLine(f *testing.F) {
	for _, seed := range []string{
		"", "one line", "first\nsecond", "intro\n- item", "a\n\nb", "a\r\nb",
		"> already quoted", "\n", "   ", "```\ncode\n```",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		out := NewMarkdown(nil).Note(text).String()

		for i, line := range strings.Split(out, internal.LineFeed()) {
			if !strings.HasPrefix(line, ">") {
				t.Errorf("line %d of the alert is outside the callout: %q\nfull output:\n%q", i, line, out)
			}
		}
	})
}

// FuzzGitHubAnchor asserts that a generated anchor holds only the characters
// GitHub keeps, and that generating it twice gives the same answer. A heading
// is arbitrary user text, so this is where punctuation and non-ASCII land.
func FuzzGitHubAnchor(f *testing.F) {
	for _, seed := range []string{
		"", "Simple", "With Space", "punctuation!?", "`code`", "日本語の見出し",
		"a--b", "  leading", "MiXeD CaSe", "1. numbered", "emoji 🎉",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, heading string) {
		anchor := generateGitHubAnchor(heading)

		for _, r := range anchor {
			if r == '-' || unicode.IsLetter(r) || unicode.IsNumber(r) {
				continue
			}
			t.Errorf("anchor %q from heading %q contains %q", anchor, heading, r)
		}
		if again := generateGitHubAnchor(heading); again != anchor {
			t.Errorf("anchor for %q is not stable: %q then %q", heading, anchor, again)
		}
	})
}

// FuzzTableOfContentsAnchorsMatchHeadings asserts that every link the table of
// contents emits points at an anchor derived from a heading in the document.
//
// Uniqueness is deliberately not asserted: GitHub resolves a duplicate anchor to
// the first heading that produces it, so a document containing both "a" twice
// and a literal "a-1" genuinely has a collision. Renaming our anchor would point
// the link at something GitHub does not have.
func FuzzTableOfContentsAnchorsMatchHeadings(f *testing.F) {
	for _, seed := range []string{"", "Same", "Same|Same", "A|B|A", "a-1|a|a", "`code`|日本語"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, joined string) {
		headings := strings.Split(joined, "|")
		if len(headings) > 32 { //nolint:mnd // keep the generated document small
			headings = headings[:32]
		}

		valid := map[string]bool{}
		m := NewMarkdown(nil).TableOfContents(TableOfContentsDepthH2)
		for _, heading := range headings {
			m = m.H2(heading)
			base := generateGitHubAnchor(heading)
			valid[base] = true
			for i := 1; i <= len(headings); i++ {
				valid[fmt.Sprintf("%s-%d", base, i)] = true
			}
		}

		// Scan only between the markers. A heading is arbitrary text and can
		// itself contain "](#", so scanning the whole document would pick up the
		// heading rather than the generated link and report a false failure.
		// The link is the last "](#" on the line for the same reason.
		entries := 0
		inTOC := false
		for _, line := range strings.Split(m.String(), internal.LineFeed()) {
			switch line {
			case TableOfContentsMarkerBegin:
				inTOC = true
				continue
			case TableOfContentsMarkerEnd:
				inTOC = false
				continue
			}
			if !inTOC {
				continue
			}

			open := strings.LastIndex(line, "](#")
			if open == -1 || !strings.HasSuffix(line, ")") {
				continue
			}
			entries++
			anchor := line[open+len("](#") : len(line)-1]
			if !valid[anchor] {
				t.Errorf("anchor %q does not come from any heading in %q", anchor, headings)
			}
		}

		if entries != len(headings) {
			t.Errorf("table of contents has %d entries for %d headings", entries, len(headings))
		}
	})
}
