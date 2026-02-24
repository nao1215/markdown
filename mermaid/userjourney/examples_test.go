//go:build linux || darwin

package userjourney_test

import (
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/userjourney"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := userjourney.NewDiagram(
		io.Discard,
		userjourney.WithTitle("Checkout Journey"),
	).
		Section("Discover").
		Task("Browse products", userjourney.ScoreVerySatisfied, "Customer").
		Task("Add item to cart", userjourney.ScoreSatisfied, "Customer").
		Section("Checkout").
		Task("Enter shipping details", userjourney.ScoreNeutral, "Customer").
		Task("Complete payment", userjourney.ScoreSatisfied, "Customer", "Payment Service").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("User Journey Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## User Journey Diagram
	// ```mermaid
	// journey
	//     title Checkout Journey
	//     section Discover
	//         Browse products: 5: Customer
	//         Add item to cart: 4: Customer
	//     section Checkout
	//         Enter shipping details: 3: Customer
	//         Complete payment: 4: Customer, Payment Service
	// ```
}
