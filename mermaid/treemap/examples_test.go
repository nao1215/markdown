//go:build linux || darwin

package treemap_test

import (
	"fmt"
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

// ExampleNewDiagram shows the shape every treemap has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = treemap.NewDiagram(os.Stdout).
		Section("Ops").Leaf("Salaries", 1200).
		Build()

	// Output:
	// treemap-beta
	// "Ops"
	//     "Salaries": 1200
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := treemap.NewDiagram(io.Discard).
		Section("Ops").Leaf("Salaries", 1200).
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// treemap-beta
	// "Ops"
	//     "Salaries": 1200
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := treemap.NewDiagram(nil).
		Section("Ops").Leaf("Salaries", 1200).
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := treemap.NewDiagram(io.Discard).
		Parent()
	fmt.Println("error:", d.Error())

	// Output:
	// error: Parent was called at the top level; there is nothing to go up to
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = treemap.NewDiagram(os.Stdout).
		Section("Ops").Leaf("Salaries", 1200).
		LF().
		Section("Ops").Leaf("Salaries", 1200).
		Build()

	// Output:
	// treemap-beta
	// "Ops"
	//     "Salaries": 1200
	//
	//     "Ops"
	//         "Salaries": 1200
}

// ExampleDiagram_full shows a treemap built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := treemap.NewDiagram(io.Discard).
		Section("Ops").Leaf("Salaries", 1200).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// treemap-beta
	// "Ops"
	//     "Salaries": 1200
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []treemap.Option{treemap.WithTitle("Overview")}

	_ = treemap.NewDiagram(os.Stdout, options...).
		Section("Ops").Leaf("Salaries", 1200).
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// treemap-beta
	// "Ops"
	//     "Salaries": 1200
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = treemap.NewDiagram(os.Stdout, treemap.WithTitle("Overview")).
		Section("Ops").Leaf("Salaries", 1200).
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// treemap-beta
	// "Ops"
	//     "Salaries": 1200
}

// ExampleDiagram_Section opens a level of the hierarchy. Everything added until
// Parent belongs to it, and a section carries no value of its own: mermaid
// gives it the sum of what it holds.
func ExampleDiagram_Section() {
	_ = treemap.NewDiagram(os.Stdout).
		Section("Ops").
		Leaf("Salaries", 1200).
		Parent().
		Section("Marketing").
		Leaf("Ads", 800).
		Build()

	// Output:
	// treemap-beta
	// "Ops"
	//     "Salaries": 1200
	// "Marketing"
	//     "Ads": 800
}

// ExampleDiagram_Leaf puts a value in the current level. A leaf is what a
// treemap actually draws: the area it gets is its value against the whole.
func ExampleDiagram_Leaf() {
	_ = treemap.NewDiagram(os.Stdout).
		Leaf("Salaries", 1200).
		Leaf("Travel", 300).
		Build()

	// Output:
	// treemap-beta
	// "Salaries": 1200
	// "Travel": 300
}

// ExampleDiagram_Parent leaves the section opened last. Calling it at the top
// level is an error rather than a silent no-op, because a chain that has gone
// up more often than down is not the diagram its author meant.
func ExampleDiagram_Parent() {
	_ = treemap.NewDiagram(os.Stdout).
		Section("Ops").
		Section("Cloud").
		Leaf("Compute", 400).
		Parent().
		Leaf("Travel", 300).
		Parent().
		Leaf("Everything else", 100).
		Build()

	// Output:
	// treemap-beta
	// "Ops"
	//     "Cloud"
	//         "Compute": 400
	//     "Travel": 300
	// "Everything else": 100
}
