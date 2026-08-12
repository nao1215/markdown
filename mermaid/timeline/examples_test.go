//go:build linux || darwin

package timeline_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/timeline"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := timeline.NewDiagram(
		io.Discard,
		timeline.WithTitle("History of Social Media"),
	).
		Period("2002", "LinkedIn").
		Section("Second wave").
		Period("2004", "Facebook", "Google").
		Period("2005", "YouTube").
		Event("Reddit").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Timeline").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Timeline
	// ```mermaid
	// timeline
	//     title History of Social Media
	//     2002 : LinkedIn
	//     section Second wave
	//         2004 : Facebook : Google
	//         2005 : YouTube : Reddit
	// ```
}

// ExampleNewDiagram shows the shape every timeline has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = timeline.NewDiagram(os.Stdout).
		Section("2024").Period("Q1", "Kickoff").
		Build()

	// Output:
	// timeline
	//     section 2024
	//         Q1 : Kickoff
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := timeline.NewDiagram(io.Discard).
		Section("2024").Period("Q1", "Kickoff").
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// timeline
	//     section 2024
	//         Q1 : Kickoff
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := timeline.NewDiagram(nil).
		Section("2024").Period("Q1", "Kickoff").
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := timeline.NewDiagram(io.Discard).
		Period("Q1", "Kickoff")
	fmt.Println("error:", d.Error())

	// Output:
	// error: <nil>
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = timeline.NewDiagram(os.Stdout).
		Section("2024").Period("Q1", "Kickoff").
		LF().
		Section("2024").Period("Q1", "Kickoff").
		Build()

	// Output:
	// timeline
	//     section 2024
	//         Q1 : Kickoff
	//
	//     section 2024
	//         Q1 : Kickoff
}

// ExampleDiagram_full shows a timeline built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := timeline.NewDiagram(io.Discard).
		Section("2024").Period("Q1", "Kickoff").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// timeline
	//     section 2024
	//         Q1 : Kickoff
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []timeline.Option{timeline.WithTitle("Overview")}

	_ = timeline.NewDiagram(os.Stdout, options...).
		Section("2024").Period("Q1", "Kickoff").
		Build()

	// Output:
	// timeline
	//     title Overview
	//     section 2024
	//         Q1 : Kickoff
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = timeline.NewDiagram(os.Stdout, timeline.WithTitle("Overview")).
		Section("2024").Period("Q1", "Kickoff").
		Build()

	// Output:
	// timeline
	//     title Overview
	//     section 2024
	//         Q1 : Kickoff
}

// ExampleDiagram_Section groups the periods that follow it. A timeline needs a
// section before anything else, and a period added without one is reported
// from Build rather than drawn in the wrong place.
func ExampleDiagram_Section() {
	_ = timeline.NewDiagram(os.Stdout).
		Section("2024").
		Period("Q1", "Kickoff").
		Section("2025").
		Period("Q1", "Launch").
		Build()

	// Output:
	// timeline
	//     section 2024
	//         Q1 : Kickoff
	//     section 2025
	//         Q1 : Launch
}

// ExampleDiagram_Period adds a point on the timeline with the events that
// happened in it. The events are optional, and each is drawn beside the period
// rather than under it.
func ExampleDiagram_Period() {
	_ = timeline.NewDiagram(os.Stdout).
		Section("2024").
		Period("Q1", "Kickoff", "Design review").
		Period("Q2").
		Build()

	// Output:
	// timeline
	//     section 2024
	//         Q1 : Kickoff : Design review
	//         Q2
}

// ExampleDiagram_Event adds an event to the period added last, which is how a
// long list of them stays readable in the chain.
func ExampleDiagram_Event() {
	_ = timeline.NewDiagram(os.Stdout).
		Section("2024").
		Period("Q1").
		Event("Kickoff").
		Event("Design review").
		Build()

	// Output:
	// timeline
	//     section 2024
	//         Q1 : Kickoff : Design review
}
