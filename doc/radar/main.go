//go:build linux || darwin

// Package main is generating mermaid radar chart.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/radar"
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

	chart := radar.NewDiagram(io.Discard, radar.WithTitle("Grades")).
		Axis("Math", "Science", "English").
		Axis("History", "Art").
		Curve("Alice", 85, 90, 80, 70, 75). //nolint:mnd
		Curve("Bob", 70, 75, 85, 80, 90).   //nolint:mnd
		Max(100).                           //nolint:mnd
		Min(0).
		String()

	err = markdown.NewMarkdown(f, markdown.WithBlockSpacing()).
		H2("Radar").
		CodeBlocks(markdown.SyntaxHighlightMermaid, chart).
		Build()
	if err != nil {
		panic(err)
	}
}
