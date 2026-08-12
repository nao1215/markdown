//go:build linux || darwin

package mindmap_test

import (
	"fmt"
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
	// title: "Product Strategy Mindmap"
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

// ExampleNewDiagram shows the shape every mindmap has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = mindmap.NewDiagram(os.Stdout).
		Root("Product").Child("Market").
		Build()

	// Output:
	// mindmap
	//     Product
	//         Market
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := mindmap.NewDiagram(io.Discard).
		Root("Product").Child("Market").
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// mindmap
	//     Product
	//         Market
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := mindmap.NewDiagram(nil).
		Root("Product").Child("Market").
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := mindmap.NewDiagram(io.Discard).
		Child("a child with no root")
	fmt.Println("error:", d.Error())

	// Output:
	// error: root node must be defined first
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = mindmap.NewDiagram(os.Stdout).
		Root("Product").
		Child("Market").
		LF().
		Sibling("Engineering").
		Build()

	// Output:
	// mindmap
	//     Product
	//         Market
	//
	//         Engineering
}

// ExampleDiagram_full shows a mindmap built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := mindmap.NewDiagram(io.Discard).
		Root("Product").Child("Market").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// mindmap
	//     Product
	//         Market
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []mindmap.Option{mindmap.WithTitle("Overview")}

	_ = mindmap.NewDiagram(os.Stdout, options...).
		Root("Product").Child("Market").
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// mindmap
	//     Product
	//         Market
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = mindmap.NewDiagram(os.Stdout, mindmap.WithTitle("Overview")).
		Root("Product").Child("Market").
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// mindmap
	//     Product
	//         Market
}

// ExampleDiagram_Root adds the node everything else hangs from. A mindmap needs
// one before anything else, and a child added without it is reported from Build
// rather than drawn at the top level.
func ExampleDiagram_Root() {
	_ = mindmap.NewDiagram(os.Stdout).
		Root("Product").
		Child("Market").
		Build()

	// Output:
	// mindmap
	//     Product
	//         Market
}

// ExampleDiagram_Child adds a node one level below the current one and descends
// into it, so a second Child is a grandchild of the first.
func ExampleDiagram_Child() {
	_ = mindmap.NewDiagram(os.Stdout).
		Root("Product").
		Child("Market").
		Child("Segments").
		Build()

	// Output:
	// mindmap
	//     Product
	//         Market
	//             Segments
}

// ExampleDiagram_Sibling adds a node beside the current one rather than under
// it, which is how a level gains a second entry.
func ExampleDiagram_Sibling() {
	_ = mindmap.NewDiagram(os.Stdout).
		Root("Product").
		Child("Market").
		Sibling("Engineering").
		Build()

	// Output:
	// mindmap
	//     Product
	//         Market
	//         Engineering
}

// ExampleDiagram_Parent moves back up a level, so what follows is a sibling of
// the node above rather than of the current one.
func ExampleDiagram_Parent() {
	_ = mindmap.NewDiagram(os.Stdout).
		Root("Product").
		Child("Market").
		Child("Segments").
		Parent().
		Sibling("Pricing").
		Build()

	// Output:
	// mindmap
	//     Product
	//         Market
	//             Segments
	//         Pricing
}

// ExampleDiagram_Node puts a node at a level named outright, for a tree that is
// walked rather than written by hand. The root is level 0, and a level more
// than one deeper than the last is reported rather than drawn.
func ExampleDiagram_Node() {
	_ = mindmap.NewDiagram(os.Stdout).
		Node(0, "Product").
		Node(1, "Market").
		Node(2, "Segments").
		Node(1, "Engineering").
		Build()

	// Output:
	// mindmap
	//     Product
	//         Market
	//             Segments
	//         Engineering
}
