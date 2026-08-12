//go:build linux || darwin

// Package main is generating mermaid Venn diagram.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/venn"
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

	diagram := venn.NewDiagram(io.Discard, venn.WithTitle("What the languages share")).
		SetWithLabel("go", "Go").
		SetWithLabel("rust", "Rust").
		SetWithLabel("compiled", "Compiled and statically typed").
		String()

	err = markdown.NewMarkdown(f, markdown.WithBlockSpacing()).
		H2("Venn").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()
	if err != nil {
		panic(err)
	}
}
