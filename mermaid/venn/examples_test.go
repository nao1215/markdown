//go:build linux || darwin

package venn_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/venn"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := venn.NewDiagram(io.Discard, venn.WithTitle("What they share")).
		SetWithLabel("go", "Go").
		SetWithLabel("rust", "Rust").
		SetWithLabel("both", "Compiled and typed").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Venn").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Venn
	// ```mermaid
	// venn-beta
	//     title What they share
	//     set go["Go"]
	//     set rust["Rust"]
	//     set both["Compiled and typed"]
	// ```
}

// ExampleNewDiagram shows the shape every Venn diagram has: a writer, the sets,
// and Build.
func ExampleNewDiagram() {
	_ = venn.NewDiagram(os.Stdout).Set("go").Set("rust").Build()

	// Output:
	// venn-beta
	//     set go
	//     set rust
}

// ExampleDiagram_Set adds a set drawn with its own name inside. Where two sets
// overlap is mermaid's business rather than something to declare, which is why
// there is no call for an intersection.
func ExampleDiagram_Set() {
	_ = venn.NewDiagram(os.Stdout).
		Set("go").
		Set("rust").
		Set("both").
		Build()

	// Output:
	// venn-beta
	//     set go
	//     set rust
	//     set both
}

// ExampleDiagram_SetWithLabel adds a set drawn with a label rather than its
// name, for a set whose identifier is not what a reader should see.
func ExampleDiagram_SetWithLabel() {
	_ = venn.NewDiagram(os.Stdout).
		SetWithLabel("go", "Written in Go").
		SetWithLabel("rust", "Written in Rust").
		Build()

	// Output:
	// venn-beta
	//     set go["Written in Go"]
	//     set rust["Written in Rust"]
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := venn.NewDiagram(io.Discard).Set("go").String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// venn-beta
	//     set go
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded.
func ExampleDiagram_Build() {
	err := venn.NewDiagram(nil).Set("go").Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the error the chain recorded, for code that
// wants to look before writing anything. A set name mermaid cannot read is
// reported rather than mangled, because there is nothing to escape to.
func ExampleDiagram_Error() {
	d := venn.NewDiagram(io.Discard).Set("not a name")
	fmt.Println("error:", d.Error())

	// Output:
	// error: set name "not a name" must hold only letters, digits, underscores and hyphens; mermaid reads nothing else there
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = venn.NewDiagram(os.Stdout).Set("go").LF().Set("rust").Build()

	// Output:
	// venn-beta
	//     set go
	//
	//     set rust
}

// ExampleWithTitle sets the title the diagram is drawn with. It is not quoted:
// mermaid reads the rest of the line, so a quoted title would draw its own
// quotation marks.
func ExampleWithTitle() {
	_ = venn.NewDiagram(os.Stdout, venn.WithTitle(`What "they" share`)).
		Set("go").
		Build()

	// Output:
	// venn-beta
	//     title What "they" share
	//     set go
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []venn.Option{venn.WithTitle("What they share")}

	_ = venn.NewDiagram(os.Stdout, options...).Set("go").Build()

	// Output:
	// venn-beta
	//     title What they share
	//     set go
}
