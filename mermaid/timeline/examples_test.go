//go:build linux || darwin

package timeline_test

import (
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
