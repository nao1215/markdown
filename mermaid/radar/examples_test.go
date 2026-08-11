//go:build linux || darwin

package radar_test

import (
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/radar"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	chart := radar.NewDiagram(io.Discard, radar.WithTitle("Grades")).
		Axis("Math", "Science", "English").
		Axis("History", "Art").
		Curve("Alice", 85, 90, 80, 70, 75).
		Curve("Bob", 70, 75, 85, 80, 90).
		Max(100).
		Min(0).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Radar").
		CodeBlocks(md.SyntaxHighlightMermaid, chart).
		Build()

	// Output:
	// ## Radar
	// ```mermaid
	// ---
	// title: "Grades"
	// ---
	// radar-beta
	//   axis a1["Math"], a2["Science"], a3["English"]
	//   axis a4["History"], a5["Art"]
	//   curve c1["Alice"]{85, 90, 80, 70, 75}
	//   curve c2["Bob"]{70, 75, 85, 80, 90}
	//   max 100
	//   min 0
	// ```
}
