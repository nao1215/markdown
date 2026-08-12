//go:build linux || darwin

package c4_test

import (
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/c4"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
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

	_ = md.NewMarkdown(os.Stdout).
		H2("C4 Context").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## C4 Context
	// ```mermaid
	// C4Context
	//     title System Context: Internet Banking
	//     Enterprise_Boundary(bank, "Big Bank plc") {
	//         Person(customer, "Personal Banking Customer", "A customer of the bank.")
	//         System_Boundary(banking, "Internet Banking") {
	//             System(web, "Internet Banking System", "Shows account information.")
	//             SystemDb(accounts, "Accounts Database")
	//         }
	//     }
	//     System_Ext(mail, "E-mail System", "The internal Microsoft Exchange system.")
	//     Rel(customer, web, "Views balances", "HTTPS")
	//     BiRel(web, accounts, "Reads from and writes to", "SQL/TCP")
	//     Rel(web, mail, "Sends e-mail using", "SMTP")
	// ```
}

// ExampleDiagram_boundary shows the pair of calls a boundary is: Boundary opens
// one and BoundaryEnd closes it, and what lies between belongs to it. They nest,
// and leaving one open is reported from Build rather than written out.
func ExampleDiagram_boundary() {
	diagram := c4.NewDiagram(io.Discard).
		Boundary("region", "eu-west-1", c4.WithBoundaryType("region")).
		System("api", "Public API").
		BoundaryEnd().
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// C4Context
	//     Boundary(region, "eu-west-1", "region") {
	//         System(api, "Public API")
	//     }
	// ```
}

// ExampleDiagram_quoting shows what this package does with the punctuation C4
// has taken for itself. A quotation mark in a label becomes "#quot;" and a "#"
// becomes "#35;", because a backslash makes mermaid refuse the diagram and a
// doubled quote silently splits the argument in two. A title is the one place
// the text is not quoted, so there a "#" and a ";" are the escaped pair.
func ExampleDiagram_quoting() {
	diagram := c4.NewDiagram(io.Discard, c4.WithTitle(`Ledger #1; the "core"`)).
		System("ledger", `The "Core" Ledger`, c4.WithDescription("Tracks #1 of everything.")).
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// C4Context
	//     title Ledger #35;1#59; the "core"
	//     System(ledger, "The #quot;Core#quot; Ledger", "Tracks #35;1 of everything.")
	// ```
}
