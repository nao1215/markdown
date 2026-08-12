//go:build linux || darwin

package packet_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/packet"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := packet.NewDiagram(
		io.Discard,
		packet.WithTitle("UDP Packet"),
	).
		Next(16, "Source Port").
		Next(16, "Destination Port").
		Field(32, 47, "Length").
		Field(48, 63, "Checksum").
		Field(64, 95, "Data (variable length)").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Packet").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Packet
	// ```mermaid
	// packet
	//     title UDP Packet
	//     +16: "Source Port"
	//     +16: "Destination Port"
	//     32-47: "Length"
	//     48-63: "Checksum"
	//     64-95: "Data (variable length)"
	// ```
}

// ExampleNewDiagram shows the shape every packet diagram has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = packet.NewDiagram(os.Stdout).
		Field(0, 15, "Source Port").
		Build()

	// Output:
	// packet
	//     0-15: "Source Port"
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := packet.NewDiagram(io.Discard).
		Field(0, 15, "Source Port").
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// packet
	//     0-15: "Source Port"
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := packet.NewDiagram(nil).
		Field(0, 15, "Source Port").
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := packet.NewDiagram(io.Discard).
		Field(15, 0, "Backwards")
	fmt.Println("error:", d.Error())

	// Output:
	// error: start bit must be less than or equal to end bit
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = packet.NewDiagram(os.Stdout).
		Field(0, 15, "Source Port").
		LF().
		Field(0, 15, "Source Port").
		Build()

	// Output:
	// packet
	//     0-15: "Source Port"
	//
	//     0-15: "Source Port"
}

// ExampleDiagram_full shows a packet diagram built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := packet.NewDiagram(io.Discard).
		Field(0, 15, "Source Port").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// packet
	//     0-15: "Source Port"
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []packet.Option{packet.WithTitle("Overview")}

	_ = packet.NewDiagram(os.Stdout, options...).
		Field(0, 15, "Source Port").
		Build()

	// Output:
	// packet
	//     title Overview
	//     0-15: "Source Port"
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = packet.NewDiagram(os.Stdout, packet.WithTitle("Overview")).
		Field(0, 15, "Source Port").
		Build()

	// Output:
	// packet
	//     title Overview
	//     0-15: "Source Port"
}

// ExampleDiagram_Field adds a field spanning a range of bits. The range is
// inclusive at both ends, so a sixteen bit field runs from 0 to 15.
func ExampleDiagram_Field() {
	_ = packet.NewDiagram(os.Stdout).
		Field(0, 15, "Source Port").
		Field(16, 31, "Destination Port").
		Build()

	// Output:
	// packet
	//     0-15: "Source Port"
	//     16-31: "Destination Port"
}

// ExampleDiagram_Bit adds a field one bit wide, which is what a flag is.
func ExampleDiagram_Bit() {
	_ = packet.NewDiagram(os.Stdout).
		Bit(0, "SYN").
		Bit(1, "ACK").
		Build()

	// Output:
	// packet
	//     0: "SYN"
	//     1: "ACK"
}

// ExampleDiagram_Next adds a field by its width rather than its position, so a
// packet can be described without counting the bits that came before it.
func ExampleDiagram_Next() {
	_ = packet.NewDiagram(os.Stdout).
		Next(16, "Source Port").
		Next(16, "Destination Port").
		Next(32, "Sequence Number").
		Build()

	// Output:
	// packet
	//     +16: "Source Port"
	//     +16: "Destination Port"
	//     +32: "Sequence Number"
}
