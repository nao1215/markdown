//go:build linux || darwin

// Package main is generating mermaid sankey diagram.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/sankey"
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

	diagram := sankey.NewDiagram(io.Discard).
		Link("Agricultural 'waste'", "Bio-conversion", 124.729). //nolint:mnd
		Link("Bio-conversion", "Liquid", 0.597).                 //nolint:mnd
		Link("Bio-conversion", "Losses, and more", 26.862).      //nolint:mnd
		Link("Bio-conversion", "Solid", 280.322).                //nolint:mnd
		Link("Bio-conversion", "Gas", 81.144).                   //nolint:mnd
		String()

	err = markdown.NewMarkdown(f, markdown.WithBlockSpacing()).
		H2("Sankey").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()
	if err != nil {
		panic(err)
	}
}
