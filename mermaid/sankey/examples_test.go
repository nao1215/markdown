//go:build linux || darwin

package sankey_test

import (
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/sankey"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := sankey.NewDiagram(io.Discard).
		Link("Coal", "Electricity", 100).
		Link("Gas", "Electricity", 50).
		Link("Electricity", "Homes", 120).
		Link("Electricity", "Losses, and more", 30).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Sankey").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Sankey
	// ```mermaid
	// sankey-beta
	//
	// Coal,Electricity,100
	// Gas,Electricity,50
	// Electricity,Homes,120
	// Electricity,"Losses, and more",30
	// ```
}
