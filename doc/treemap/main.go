//go:build linux || darwin

// Package main is generating mermaid treemap diagram.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/treemap"
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

	diagram := treemap.NewDiagram(io.Discard, treemap.WithTitle("Budget")).
		Section("Ops").
		Leaf("Salaries", 1200). //nolint:mnd
		Section("Cloud").
		Leaf("Compute", 400). //nolint:mnd
		Parent().
		Leaf("Travel", 300). //nolint:mnd
		Parent().
		Section("Marketing").
		Leaf("Ads", 800). //nolint:mnd
		String()

	err = markdown.NewMarkdown(f, markdown.WithBlockSpacing()).
		H2("Treemap").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()
	if err != nil {
		panic(err)
	}
}
