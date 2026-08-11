//go:build linux || darwin

// Package main is generating mermaid timeline diagram.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/timeline"
)

// This file is gated by //go:build linux || darwin, so //go:generate is skipped
// on Windows. To regenerate generated.md on Windows, run under WSL or via CI.
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

	diagram := timeline.NewDiagram(
		io.Discard,
		timeline.WithTitle("History of Social Media"),
	).
		Period("2002", "LinkedIn").
		Section("Second wave").
		Period("2004", "Facebook", "Google").
		Period("2005", "YouTube").
		Section("Third wave").
		Period("2006", "Twitter").
		Event("Reddit").
		String()

	err = markdown.NewMarkdown(f, markdown.WithBlockSpacing()).
		H2("Timeline").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()
	if err != nil {
		panic(err)
	}
}
