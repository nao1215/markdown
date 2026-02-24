//go:build linux || darwin

// Package main is generating mermaid XY chart.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/xychart"
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

	diagram := xychart.NewDiagram(io.Discard, xychart.WithTitle("Sales Revenue")).
		XAxisLabels("Jan", "Feb", "Mar", "Apr", "May", "Jun").
		YAxisRangeWithTitle("Revenue (k$)", 0, 100). //nolint:mnd
		Bar(25, 40, 60, 80, 70, 90).                 //nolint:mnd
		Line(30, 50, 70, 85, 75, 95).                //nolint:mnd
		String()

	err = markdown.NewMarkdown(f).
		H2("XY Chart").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()
	if err != nil {
		panic(err)
	}
}
