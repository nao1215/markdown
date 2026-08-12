//go:build linux || darwin

// Package main is generating flowchart.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/flowchart"
)

//go:generate go run main.go

func main() {
	f, err := os.Create("generated.md")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	fc := flowchart.NewFlowchart(
		io.Discard,
		flowchart.WithTitle("mermaid flowchart builder"),
		flowchart.WithOrientalTopToBottom(),
	).
		Subgraph("ingest", "Ingest").
		SubgraphDirection(flowchart.DirectionLR).
		NodeWithText("A", "Node A").
		StadiumNode("B", "Node B").
		LinkWithArrowHead("A", "B").
		SubgraphEnd().
		SubroutineNode("C", "Node C").
		DatabaseNode("D", "Database").
		LinkWithArrowHeadAndText("B", "D", "send original data").
		LinkWithArrowHead("B", "C").
		DottedLinkWithText("C", "D", "send filtered data").
		ClassDef("stored", "fill:#d4f7d4,stroke:#2b8a3e").
		Class("D", "stored").
		Style("C", "fill:#fff3bf,stroke:#e67700").
		ClickHref("D", "https://example.com/database", "The database").
		String()

	err = markdown.NewMarkdown(f, markdown.WithBlockSpacing()).
		H2("Flowchart").
		CodeBlocks(markdown.SyntaxHighlightMermaid, fc).
		Build()

	if err != nil {
		panic(err)
	}
}
