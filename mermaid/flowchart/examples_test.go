//go:build linux || darwin

package flowchart_test

import (
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/flowchart"
)

// ExampleFlowchart skip this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleFlowchart() {
	fc := flowchart.NewFlowchart(
		os.Stdout,
		flowchart.WithTitle("mermaid flowchart builder"),
		flowchart.WithOrientalTopToBottom(),
	).
		NodeWithText("A", "Node A").
		StadiumNode("B", "Node B").
		SubroutineNode("C", "Node C").
		DatabaseNode("D", "Database").
		LinkWithArrowHead("A", "B").
		LinkWithArrowHeadAndText("B", "D", "send original data").
		LinkWithArrowHead("B", "C").
		DottedLinkWithText("C", "D", "send filtered data").
		String()

	markdown.NewMarkdown(os.Stdout).
		H2("Flowchart").
		CodeBlocks(markdown.SyntaxHighlightMermaid, fc).
		Build() //nolint

	// Output:
	//## Flowchart
	//```mermaid
	//---
	//title: mermaid flowchart builder
	//---
	//flowchart TB
	//     A["Node A"]
	//     B(["Node B"])
	//     C[["Node C"]]
	//     D[("Database")]
	//     A-->B
	//     B-->|"send original data"|D
	//     B-->C
	//     C-. "send filtered data" .-> D
	//```
}
