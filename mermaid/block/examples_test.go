//go:build linux || darwin

package block_test

import (
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/block"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := block.NewDiagram(
		io.Discard,
		block.WithTitle("Checkout Architecture"),
	).
		Columns(3).
		Row(
			block.Node("Frontend"),
			block.ArrowRight("toBackend", block.WithArrowLabel("calls")),
			block.Node("Backend"),
		).
		Row(
			block.Space(2),
			block.ArrowDown("toDB"),
		).
		Row(
			block.Node("Database", block.WithNodeLabel("Primary DB"), block.WithNodeShape(block.ShapeCylinder)),
			block.Space(),
			block.Node("Cache", block.WithNodeLabel("Cache"), block.WithNodeShape(block.ShapeRound)),
		).
		Link("Backend", "Database").
		LinkWithLabel("Backend", "reads from", "Cache").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Block Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Block Diagram
	// ```mermaid
	// block
	//     title Checkout Architecture
	//     columns 3
	//     Frontend toBackend<["calls"]>(right) Backend
	//     space:2 toDB<["&nbsp;"]>(down)
	//     Database[("Primary DB")] space Cache("Cache")
	//     Backend --> Database
	//     Backend -- "reads from" --> Cache
	// ```
}
