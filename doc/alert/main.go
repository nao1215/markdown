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
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if err := md.NewMarkdown(f, md.WithBlockSpacing()).
		H1("Alert example").
		Note("This is note").
		Tip("This is tip").
		Important("This is important").
		Warning("This is warning").
		Caution("This is caution").
		Build(); err != nil {
		panic(err)
	}
}
