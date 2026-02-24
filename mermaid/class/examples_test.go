//go:build linux || darwin

package class_test

import (
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/class"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := class.NewDiagram(
		os.Stdout,
		class.WithTitle("Checkout Domain"),
	).
		SetDirection(class.DirectionLR).
		Class(
			"Order",
			class.WithPublicField("string", "id"),
			class.WithPublicMethod("Create", "error", "items []LineItem"),
			class.WithPublicMethod("Pay", "error"),
		).
		Class(
			"LineItem",
			class.WithPublicField("string", "sku"),
			class.WithPublicField("int", "quantity"),
			class.WithPublicMethod("Subtotal", "int"),
		).
		Interface("PaymentGateway")

	diagram.From("Order").
		Composition("LineItem", class.WithOneToMany(), class.WithRelationLabel("contains")).
		Association("PaymentGateway", class.WithRelationLabel("uses"))

	diagramString := diagram.
		NoteFor("Order", "Aggregate Root").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Class Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagramString).
		Build()

	// Output:
	// ## Class Diagram
	// ```mermaid
	// ---
	// title: Checkout Domain
	// ---
	// classDiagram
	//     direction LR
	//     class Order {
	//         +string id
	//         +Create(items []LineItem) error
	//         +Pay() error
	//     }
	//     class LineItem {
	//         +string sku
	//         +int quantity
	//         +Subtotal() int
	//     }
	//     class PaymentGateway
	//     <<Interface>> PaymentGateway
	//     Order "1" *-- "many" LineItem : contains
	//     Order --> PaymentGateway : uses
	//     note for Order "Aggregate Root"
	// ```
}
