package xychart_test

import (
	"strings"
	"testing"

	"github.com/nao1215/markdown/mermaid/xychart"
)

// TestXAxisRange covers the numeric x-axis, including the rejected range. A
// chart whose axis silently fails to render is worse than one that reports the
// problem, so the error path matters as much as the happy one.
func TestXAxisRange(t *testing.T) {
	t.Parallel()

	t.Run("with a title", func(t *testing.T) {
		t.Parallel()

		d := xychart.NewDiagram(nil).XAxisRangeWithTitle("Month", 1, 12)
		if err := d.Error(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `    x-axis Month 1 --> 12`
		if !strings.Contains(d.String(), want) {
			t.Errorf("axis missing from output:\n got: %q\nwant to contain: %q", d.String(), want)
		}
	})

	t.Run("min not below max is rejected", func(t *testing.T) {
		t.Parallel()

		d := xychart.NewDiagram(nil).XAxisRangeWithTitle("Month", 12, 12)
		if d.Error() == nil {
			t.Fatal("an inverted range must be reported")
		}
	})

	t.Run("a later call is a no-op once the diagram holds an error", func(t *testing.T) {
		t.Parallel()

		d := xychart.NewDiagram(nil).
			XAxisRangeWithTitle("Month", 12, 12).
			XAxisRangeWithTitle("Month", 1, 12)

		if strings.Contains(d.String(), "1 --> 12") {
			t.Errorf("statement was appended after an error:\n%s", d.String())
		}
	})
}

// TestLineFeedStopsOnError pins the same rule for LF: once the diagram is in an
// error state, nothing more is appended.
func TestLineFeedStopsOnError(t *testing.T) {
	t.Parallel()

	clean := xychart.NewDiagram(nil).LF()
	if !strings.HasSuffix(clean.String(), "\n") {
		t.Errorf("LF did not append a blank line: %q", clean.String())
	}

	broken := xychart.NewDiagram(nil).XAxisRangeWithTitle("Month", 12, 12)
	before := broken.String()
	if after := broken.LF().String(); after != before {
		t.Errorf("LF appended after an error:\n before: %q\n  after: %q", before, after)
	}
}
