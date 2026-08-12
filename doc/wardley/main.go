//go:build linux || darwin

// Package main is generating mermaid Wardley map.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/wardley"
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

	diagram := wardley.NewMap(io.Discard, wardley.WithTitle("Checkout, as it stands")).
		Anchor("Customer", 0.95, 0.95).          //nolint:mnd
		Component("Checkout (web)", 0.6, 0.8).   //nolint:mnd
		Component("Payment service", 0.75, 0.5). //nolint:mnd
		Component("Card network", 0.95, 0.2).    //nolint:mnd
		Link("Customer", "Checkout (web)").
		Link("Checkout (web)", "Payment service").
		Link("Payment service", "Card network").
		Evolve("Payment service", 0.9). //nolint:mnd
		String()

	err = markdown.NewMarkdown(f, markdown.WithBlockSpacing()).
		H2("Wardley map").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()
	if err != nil {
		panic(err)
	}
}
