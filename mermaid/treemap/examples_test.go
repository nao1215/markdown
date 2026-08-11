//go:build linux || darwin

package treemap_test

import (
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/treemap"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := treemap.NewDiagram(io.Discard, treemap.WithTitle("Budget")).
		Section("Ops").
		Leaf("Salaries", 1200).
		Section("Cloud").
		Leaf("Compute", 400).
		Parent().
		Leaf("Travel", 300).
		Parent().
		Section("Marketing").
		Leaf("Ads", 800).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Treemap").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Treemap
	// ```mermaid
	// ---
	// title: "Budget"
	// ---
	// treemap-beta
	// "Ops"
	//     "Salaries": 1200
	//     "Cloud"
	//         "Compute": 400
	//     "Travel": 300
	// "Marketing"
	//     "Ads": 800
	// ```
}
