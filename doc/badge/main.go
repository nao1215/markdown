//go:build linux || darwin

// Package main is generating markdown.
package main

import (
	"os"

	md "github.com/nao1215/markdown"
)

//go:generate go run main.go

func main() {
	f, err := os.Create("generated.md")
	if err != nil {
		panic(err)
	}

	if err := md.NewMarkdown(f).
		H1("badge example").
		RedBadge("red_badge").
		YellowBadge("yellow_badge").
		GreenBadge("green_badge").
		Build(); err != nil {
		panic(err)
	}
}
