//go:build linux || darwin

package state_test

import (
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/state"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := state.NewDiagram(
		os.Stdout,
		state.WithTitle("Simple State Diagram"),
	).
		StartTransition("Still").
		Transition("Still", "Moving").
		TransitionWithNote("Moving", "Crash", "sudden stop").
		EndTransition("Crash").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("State Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## State Diagram
	// ```mermaid
	// ---
	// title: Simple State Diagram
	// ---
	// stateDiagram-v2
	//     [*] --> Still
	//     Still --> Moving
	//     Moving --> Crash : sudden stop
	//     Crash --> [*]
	// ```
}
