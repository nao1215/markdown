//go:build linux || darwin

package packet_test

import (
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
