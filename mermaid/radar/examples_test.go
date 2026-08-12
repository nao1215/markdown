//go:build linux || darwin

package radar_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/radar"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	chart := radar.NewDiagram(io.Discard, radar.WithTitle("Grades")).
		Axis("Math", "Science", "English").
		Axis("History", "Art").
		Curve("Alice", 85, 90, 80, 70, 75).
		Curve("Bob", 70, 75, 85, 80, 90).
		Max(100).
		Min(0).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Radar").
		CodeBlocks(md.SyntaxHighlightMermaid, chart).
		Build()

	// Output:
	// ## Radar
	// ```mermaid
	// ---
	// title: "Grades"
	// ---
	// radar-beta
	//   axis a1["Math"], a2["Science"], a3["English"]
	//   axis a4["History"], a5["Art"]
	//   curve c1["Alice"]{85, 90, 80, 70, 75}
	//   curve c2["Bob"]{70, 75, 85, 80, 90}
	//   max 100
	//   min 0
	// ```
}

// ExampleNewDiagram shows the shape every radar chart has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = radar.NewDiagram(os.Stdout).
		Axis("Math", "Science").Curve("Alice", 85, 90).
		Build()

	// Output:
	// radar-beta
	//   axis a1["Math"], a2["Science"]
	//   curve c1["Alice"]{85, 90}
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := radar.NewDiagram(io.Discard).
		Axis("Math", "Science").Curve("Alice", 85, 90).
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// radar-beta
	//   axis a1["Math"], a2["Science"]
	//   curve c1["Alice"]{85, 90}
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := radar.NewDiagram(nil).
		Axis("Math", "Science").Curve("Alice", 85, 90).
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := radar.NewDiagram(io.Discard).
		Curve("Alice", 85, 90)
	fmt.Println("error:", d.Error())

	// Output:
	// error: <nil>
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = radar.NewDiagram(os.Stdout).
		Axis("Math", "Science").Curve("Alice", 85, 90).
		LF().
		Axis("Math", "Science").Curve("Alice", 85, 90).
		Build()

	// Output:
	// radar-beta
	//   axis a1["Math"], a2["Science"]
	//   curve c1["Alice"]{85, 90}
	//
	//   axis a3["Math"], a4["Science"]
	//   curve c2["Alice"]{85, 90}
}

// ExampleDiagram_full shows a radar chart built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := radar.NewDiagram(io.Discard).
		Axis("Math", "Science").Curve("Alice", 85, 90).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// radar-beta
	//   axis a1["Math"], a2["Science"]
	//   curve c1["Alice"]{85, 90}
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []radar.Option{radar.WithTitle("Overview")}

	_ = radar.NewDiagram(os.Stdout, options...).
		Axis("Math", "Science").Curve("Alice", 85, 90).
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// radar-beta
	//   axis a1["Math"], a2["Science"]
	//   curve c1["Alice"]{85, 90}
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = radar.NewDiagram(os.Stdout, radar.WithTitle("Overview")).
		Axis("Math", "Science").Curve("Alice", 85, 90).
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// radar-beta
	//   axis a1["Math"], a2["Science"]
	//   curve c1["Alice"]{85, 90}
}

// ExampleDiagram_Axis declares the axes, in order. Every curve then gives its
// values in that same order. mermaid wants an identifier in front of each
// label; nothing in a radar chart refers to one, so the package numbers them
// and the caller passes only the labels.
func ExampleDiagram_Axis() {
	_ = radar.NewDiagram(os.Stdout).
		Axis("Math", "Science", "English").
		Curve("Alice", 85, 90, 80).
		Build()

	// Output:
	// radar-beta
	//   axis a1["Math"], a2["Science"], a3["English"]
	//   curve c1["Alice"]{85, 90, 80}
}

// ExampleDiagram_Curve adds one subject's values, in the order the axes were
// declared.
func ExampleDiagram_Curve() {
	_ = radar.NewDiagram(os.Stdout).
		Axis("Math", "Science").
		Curve("Alice", 85, 90).
		Curve("Bob", 70, 75).
		Build()

	// Output:
	// radar-beta
	//   axis a1["Math"], a2["Science"]
	//   curve c1["Alice"]{85, 90}
	//   curve c2["Bob"]{70, 75}
}

// ExampleDiagram_Max sets the outer edge of the chart. Without it mermaid picks
// a bound from the values, so two charts drawn from different data cannot be
// compared by eye.
func ExampleDiagram_Max() {
	_ = radar.NewDiagram(os.Stdout).
		Axis("Math", "Science").
		Curve("Alice", 85, 90).
		Max(100).
		Build()

	// Output:
	// radar-beta
	//   axis a1["Math"], a2["Science"]
	//   curve c1["Alice"]{85, 90}
	//   max 100
}

// ExampleDiagram_Min sets the center of the chart.
func ExampleDiagram_Min() {
	_ = radar.NewDiagram(os.Stdout).
		Axis("Math", "Science").
		Curve("Alice", 85, 90).
		Min(0).
		Max(100).
		Build()

	// Output:
	// radar-beta
	//   axis a1["Math"], a2["Science"]
	//   curve c1["Alice"]{85, 90}
	//   min 0
	//   max 100
}
