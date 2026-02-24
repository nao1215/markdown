//go:build linux || darwin

// Package main is generating mermaid packet diagram.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/packet"
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

	diagram := packet.NewDiagram(io.Discard, packet.WithTitle("UDP Packet")).
		Next(16, "Source Port").                 //nolint:mnd
		Next(16, "Destination Port").            //nolint:mnd
		Field(32, 47, "Length").                 //nolint:mnd
		Field(48, 63, "Checksum").               //nolint:mnd
		Field(64, 95, "Data (variable length)"). //nolint:mnd
		String()

	err = markdown.NewMarkdown(f).
		H2("Packet").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()
	if err != nil {
		panic(err)
	}
}
