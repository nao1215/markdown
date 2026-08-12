//go:build linux || darwin

package c4_test

import (
	"fmt"
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

// ExampleNewDiagram shows the shape every C4 context diagram has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = c4.NewDiagram(os.Stdout).
		Person("customer", "Customer").System("ledger", "Ledger").
		Build()

	// Output:
	// C4Context
	//     Person(customer, "Customer")
	//     System(ledger, "Ledger")
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := c4.NewDiagram(io.Discard).
		Person("customer", "Customer").System("ledger", "Ledger").
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// C4Context
	//     Person(customer, "Customer")
	//     System(ledger, "Ledger")
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := c4.NewDiagram(nil).
		Person("customer", "Customer").System("ledger", "Ledger").
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := c4.NewDiagram(io.Discard).
		BoundaryEnd()
	fmt.Println("error:", d.Error())

	// Output:
	// error: BoundaryEnd was called outside a boundary; there is nothing to close
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = c4.NewDiagram(os.Stdout).
		Person("customer", "Customer").System("ledger", "Ledger").
		LF().
		Person("customer", "Customer").System("ledger", "Ledger").
		Build()

	// Output:
	// C4Context
	//     Person(customer, "Customer")
	//     System(ledger, "Ledger")
	//
	//     Person(customer, "Customer")
	//     System(ledger, "Ledger")
}

// ExampleDiagram_full shows a C4 context diagram built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := c4.NewDiagram(io.Discard).
		Person("customer", "Customer").System("ledger", "Ledger").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// C4Context
	//     Person(customer, "Customer")
	//     System(ledger, "Ledger")
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []c4.Option{c4.WithTitle("Overview")}

	_ = c4.NewDiagram(os.Stdout, options...).
		Person("customer", "Customer").System("ledger", "Ledger").
		Build()

	// Output:
	// C4Context
	//     title Overview
	//     Person(customer, "Customer")
	//     System(ledger, "Ledger")
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = c4.NewDiagram(os.Stdout, c4.WithTitle("Overview")).
		Person("customer", "Customer").System("ledger", "Ledger").
		Build()

	// Output:
	// C4Context
	//     title Overview
	//     Person(customer, "Customer")
	//     System(ledger, "Ledger")
}

// ExampleDiagram_Person adds someone inside the enterprise being described.
func ExampleDiagram_Person() {
	_ = c4.NewDiagram(os.Stdout).
		Person("customer", "Personal Banking Customer").
		Build()

	// Output:
	// C4Context
	//     Person(customer, "Personal Banking Customer")
}

// ExampleDiagram_PersonExt adds someone outside it, which mermaid draws in a
// colder color.
func ExampleDiagram_PersonExt() {
	_ = c4.NewDiagram(os.Stdout).
		PersonExt("auditor", "External Auditor").
		Build()

	// Output:
	// C4Context
	//     Person_Ext(auditor, "External Auditor")
}

// ExampleDiagram_System adds a software system inside the enterprise.
func ExampleDiagram_System() {
	_ = c4.NewDiagram(os.Stdout).
		System("banking", "Internet Banking System").
		Build()

	// Output:
	// C4Context
	//     System(banking, "Internet Banking System")
}

// ExampleDiagram_SystemExt adds a software system outside it.
func ExampleDiagram_SystemExt() {
	_ = c4.NewDiagram(os.Stdout).
		SystemExt("mail", "E-mail System").
		Build()

	// Output:
	// C4Context
	//     System_Ext(mail, "E-mail System")
}

// ExampleDiagram_SystemDb adds a software system drawn as a database.
func ExampleDiagram_SystemDb() {
	_ = c4.NewDiagram(os.Stdout).
		SystemDb("accounts", "Accounts Database").
		Build()

	// Output:
	// C4Context
	//     SystemDb(accounts, "Accounts Database")
}

// ExampleDiagram_SystemQueue adds a software system drawn as a queue.
func ExampleDiagram_SystemQueue() {
	_ = c4.NewDiagram(os.Stdout).
		SystemQueue("events", "Event Bus").
		Build()

	// Output:
	// C4Context
	//     SystemQueue(events, "Event Bus")
}

// ExampleDiagram_Boundary opens a box the elements after it are drawn inside.
// BoundaryEnd closes it, and leaving one open is reported from Build because
// mermaid refuses a diagram whose brace never closes.
func ExampleDiagram_Boundary() {
	_ = c4.NewDiagram(os.Stdout).
		Boundary("region", "eu-west-1", c4.WithBoundaryType("region")).
		System("api", "Public API").
		BoundaryEnd().
		Build()

	// Output:
	// C4Context
	//     Boundary(region, "eu-west-1", "region") {
	//         System(api, "Public API")
	//     }
}

// ExampleDiagram_EnterpriseBoundary opens a boundary drawn as the enterprise.
func ExampleDiagram_EnterpriseBoundary() {
	_ = c4.NewDiagram(os.Stdout).
		EnterpriseBoundary("bank", "Big Bank plc").
		Person("customer", "Customer").
		BoundaryEnd().
		Build()

	// Output:
	// C4Context
	//     Enterprise_Boundary(bank, "Big Bank plc") {
	//         Person(customer, "Customer")
	//     }
}

// ExampleDiagram_SystemBoundary opens a boundary drawn as a software system,
// and boundaries nest.
func ExampleDiagram_SystemBoundary() {
	_ = c4.NewDiagram(os.Stdout).
		EnterpriseBoundary("bank", "Big Bank plc").
		SystemBoundary("banking", "Internet Banking").
		System("web", "Web Application").
		BoundaryEnd().
		BoundaryEnd().
		Build()

	// Output:
	// C4Context
	//     Enterprise_Boundary(bank, "Big Bank plc") {
	//         System_Boundary(banking, "Internet Banking") {
	//             System(web, "Web Application")
	//         }
	//     }
}

// ExampleDiagram_BoundaryEnd closes the boundary opened last. Calling it
// outside one is an error rather than a silent no-op.
func ExampleDiagram_BoundaryEnd() {
	_ = c4.NewDiagram(os.Stdout).
		Boundary("region", "eu-west-1").
		System("api", "Public API").
		BoundaryEnd().
		System("outside", "Somewhere Else").
		Build()

	// Output:
	// C4Context
	//     Boundary(region, "eu-west-1") {
	//         System(api, "Public API")
	//     }
	//     System(outside, "Somewhere Else")
}

// ExampleDiagram_Rel draws a one way relationship.
func ExampleDiagram_Rel() {
	_ = c4.NewDiagram(os.Stdout).
		Person("customer", "Customer").
		System("ledger", "Ledger").
		Rel("customer", "ledger", "Views balances", c4.WithTechnology("HTTPS")).
		Build()

	// Output:
	// C4Context
	//     Person(customer, "Customer")
	//     System(ledger, "Ledger")
	//     Rel(customer, ledger, "Views balances", "HTTPS")
}

// ExampleDiagram_BiRel draws a relationship with an arrowhead at each end.
func ExampleDiagram_BiRel() {
	_ = c4.NewDiagram(os.Stdout).
		System("web", "Web Application").
		SystemDb("accounts", "Accounts Database").
		BiRel("web", "accounts", "Reads from and writes to", c4.WithTechnology("SQL/TCP")).
		Build()

	// Output:
	// C4Context
	//     System(web, "Web Application")
	//     SystemDb(accounts, "Accounts Database")
	//     BiRel(web, accounts, "Reads from and writes to", "SQL/TCP")
}

// ExampleWithDescription puts a sentence under an element's label.
func ExampleWithDescription() {
	_ = c4.NewDiagram(os.Stdout).
		Person("customer", "Customer", c4.WithDescription("A customer of the bank.")).
		Build()

	// Output:
	// C4Context
	//     Person(customer, "Customer", "A customer of the bank.")
}

// ExampleWithTechnology says how two things talk to each other, drawn in
// brackets on the arrow.
func ExampleWithTechnology() {
	_ = c4.NewDiagram(os.Stdout).
		Person("customer", "Customer").
		System("ledger", "Ledger").
		Rel("customer", "ledger", "Views balances", c4.WithTechnology("HTTPS (TLS 1.3)")).
		Build()

	// Output:
	// C4Context
	//     Person(customer, "Customer")
	//     System(ledger, "Ledger")
	//     Rel(customer, ledger, "Views balances", "HTTPS (TLS 1.3)")
}

// ExampleWithBoundaryType puts a tag in brackets beside a boundary's label.
// EnterpriseBoundary and SystemBoundary carry their own, so this applies to
// Boundary alone.
func ExampleWithBoundaryType() {
	_ = c4.NewDiagram(os.Stdout).
		Boundary("regulator", "Regulator", c4.WithBoundaryType("external")).
		SystemExt("reporting", "Regulatory Reporting").
		BoundaryEnd().
		Build()

	// Output:
	// C4Context
	//     Boundary(regulator, "Regulator", "external") {
	//         System_Ext(reporting, "Regulatory Reporting")
	//     }
}

// ExampleElementOption shows what an ElementOption is: a function that changes
// how an element is written, passed after its label.
func ExampleElementOption() {
	options := []c4.ElementOption{c4.WithDescription("A customer of the bank.")}

	_ = c4.NewDiagram(os.Stdout).Person("customer", "Customer", options...).Build()

	// Output:
	// C4Context
	//     Person(customer, "Customer", "A customer of the bank.")
}

// ExampleRelationOption shows what a RelationOption is: a function that changes
// how a relationship is written, passed after its label.
func ExampleRelationOption() {
	options := []c4.RelationOption{c4.WithTechnology("HTTPS")}

	_ = c4.NewDiagram(os.Stdout).
		Person("customer", "Customer").
		System("ledger", "Ledger").
		Rel("customer", "ledger", "Views balances", options...).
		Build()

	// Output:
	// C4Context
	//     Person(customer, "Customer")
	//     System(ledger, "Ledger")
	//     Rel(customer, ledger, "Views balances", "HTTPS")
}

// ExampleBoundaryOption shows what a BoundaryOption is: a function that changes
// how a boundary is written, passed after its label.
func ExampleBoundaryOption() {
	options := []c4.BoundaryOption{c4.WithBoundaryType("region")}

	_ = c4.NewDiagram(os.Stdout).
		Boundary("region", "eu-west-1", options...).
		System("api", "Public API").
		BoundaryEnd().
		Build()

	// Output:
	// C4Context
	//     Boundary(region, "eu-west-1", "region") {
	//         System(api, "Public API")
	//     }
}
