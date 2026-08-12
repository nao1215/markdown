//go:build linux || darwin

// Package main is generating mermaid C4 context diagram.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/c4"
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

	diagram := c4.NewDiagram(io.Discard, c4.WithTitle("System Context: Internet Banking")).
		EnterpriseBoundary("bank", "Big Bank plc").
		Person("customer", "Personal Banking Customer", c4.WithDescription("A customer of the bank.")).
		SystemBoundary("banking", "Internet Banking").
		System("web", "Internet Banking System", c4.WithDescription("Shows account information.")).
		SystemDb("accounts", "Accounts Database").
		BoundaryEnd().
		BoundaryEnd().
		SystemExt("mail", "E-mail System", c4.WithDescription("The internal Microsoft Exchange system.")).
		Rel("customer", "web", "Views balances", c4.WithTechnology("HTTPS")).
		BiRel("web", "accounts", "Reads from and writes to", c4.WithTechnology("SQL/TCP")).
		Rel("web", "mail", "Sends e-mail using", c4.WithTechnology("SMTP")).
		String()

	err = markdown.NewMarkdown(f, markdown.WithBlockSpacing()).
		H2("C4 Context").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()
	if err != nil {
		panic(err)
	}
}
