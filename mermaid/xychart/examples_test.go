//go:build linux || darwin

package xychart_test

import (
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/xychart"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := xychart.NewDiagram(
		io.Discard,
		xychart.WithTitle("Sales Revenue"),
	).
		XAxisLabels("Jan", "Feb", "Mar", "Apr", "May", "Jun").
		YAxisRangeWithTitle("Revenue (k$)", 0, 100).
		Bar(25, 40, 60, 80, 70, 90).
		Line(30, 50, 70, 85, 75, 95).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("XY Chart").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## XY Chart
	// ```mermaid
	// xychart
	//     title "Sales Revenue"
	//     x-axis [Jan, Feb, Mar, Apr, May, Jun]
	//     y-axis "Revenue (k$)" 0 --> 100
	//     bar [25, 40, 60, 80, 70, 90]
	//     line [30, 50, 70, 85, 75, 95]
	// ```
}
