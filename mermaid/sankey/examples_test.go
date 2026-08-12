//go:build linux || darwin

package sankey_test

import (
	"fmt"
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

// ExampleNewDiagram shows the shape every sankey diagram has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = sankey.NewDiagram(os.Stdout).
		Link("Bill", "Electricity", 120).
		Build()

	// Output:
	// sankey-beta
	//
	// Bill,Electricity,120
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := sankey.NewDiagram(io.Discard).
		Link("Bill", "Electricity", 120).
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// sankey-beta
	//
	// Bill,Electricity,120
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := sankey.NewDiagram(nil).
		Link("Bill", "Electricity", 120).
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := sankey.NewDiagram(io.Discard).
		Link("", "Electricity", 120)
	fmt.Println("error:", d.Error())

	// Output:
	// error: source must not be empty
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = sankey.NewDiagram(os.Stdout).
		Link("Bill", "Electricity", 120).
		LF().
		Link("Bill", "Electricity", 120).
		Build()

	// Output:
	// sankey-beta
	//
	// Bill,Electricity,120
	//
	// Bill,Electricity,120
}

// ExampleDiagram_full shows a sankey diagram built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := sankey.NewDiagram(io.Discard).
		Link("Bill", "Electricity", 120).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// sankey-beta
	//
	// Bill,Electricity,120
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []sankey.Option{sankey.WithTitle("Overview")}

	_ = sankey.NewDiagram(os.Stdout, options...).
		Link("Bill", "Electricity", 120).
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// sankey-beta
	//
	// Bill,Electricity,120
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = sankey.NewDiagram(os.Stdout, sankey.WithTitle("Overview")).
		Link("Bill", "Electricity", 120).
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// sankey-beta
	//
	// Bill,Electricity,120
}

// ExampleDiagram_Link adds one flow. The value is the width of the band drawn
// between the two nodes, and a node appears in the diagram because a link names
// it: there is nothing else to declare.
func ExampleDiagram_Link() {
	_ = sankey.NewDiagram(os.Stdout).
		Link("Bill", "Electricity", 120).
		Link("Bill", "Gas", 80).
		Link("Electricity", "Lighting", 45).
		Build()

	// Output:
	// sankey-beta
	//
	// Bill,Electricity,120
	// Bill,Gas,80
	// Electricity,Lighting,45
}
