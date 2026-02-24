//go:build linux || darwin

package mindmap_test

import (
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/mindmap"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := mindmap.NewDiagram(
		io.Discard,
		mindmap.WithTitle("Product Strategy Mindmap"),
	).
		Root("Product Strategy").
		Child("Market").
		Child("SMB").
		Sibling("Enterprise").
		Parent().
		Sibling("Execution").
		Child("Q1").
		Sibling("Q2").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Mindmap").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Mindmap
	// ```mermaid
	// ---
	// title: Product Strategy Mindmap
	// ---
	// mindmap
	//     Product Strategy
	//         Market
	//             SMB
	//             Enterprise
	//         Execution
	//             Q1
	//             Q2
	// ```
}
