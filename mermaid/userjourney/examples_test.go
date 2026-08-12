//go:build linux || darwin

package userjourney_test

import (
	"fmt"
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

// ExampleNewDiagram shows the shape every user journey has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = userjourney.NewDiagram(os.Stdout).
		Section("Checkout").Task("Add to basket", userjourney.ScoreSatisfied, "Customer").
		Build()

	// Output:
	// journey
	//     section Checkout
	//         Add to basket: 4: Customer
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := userjourney.NewDiagram(io.Discard).
		Section("Checkout").Task("Add to basket", userjourney.ScoreSatisfied, "Customer").
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// journey
	//     section Checkout
	//         Add to basket: 4: Customer
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := userjourney.NewDiagram(nil).
		Section("Checkout").Task("Add to basket", userjourney.ScoreSatisfied, "Customer").
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := userjourney.NewDiagram(io.Discard).
		Task("no section yet", userjourney.ScoreSatisfied)
	fmt.Println("error:", d.Error())

	// Output:
	// error: task "no section yet" requires a section; call Section first
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = userjourney.NewDiagram(os.Stdout).
		Section("Checkout").Task("Add to basket", userjourney.ScoreSatisfied, "Customer").
		LF().
		Section("Checkout").Task("Add to basket", userjourney.ScoreSatisfied, "Customer").
		Build()

	// Output:
	// journey
	//     section Checkout
	//         Add to basket: 4: Customer
	//
	//     section Checkout
	//         Add to basket: 4: Customer
}

// ExampleDiagram_full shows a user journey built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := userjourney.NewDiagram(io.Discard).
		Section("Checkout").Task("Add to basket", userjourney.ScoreSatisfied, "Customer").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// journey
	//     section Checkout
	//         Add to basket: 4: Customer
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []userjourney.Option{userjourney.WithTitle("Overview")}

	_ = userjourney.NewDiagram(os.Stdout, options...).
		Section("Checkout").Task("Add to basket", userjourney.ScoreSatisfied, "Customer").
		Build()

	// Output:
	// journey
	//     title Overview
	//     section Checkout
	//         Add to basket: 4: Customer
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = userjourney.NewDiagram(os.Stdout, userjourney.WithTitle("Overview")).
		Section("Checkout").Task("Add to basket", userjourney.ScoreSatisfied, "Customer").
		Build()

	// Output:
	// journey
	//     title Overview
	//     section Checkout
	//         Add to basket: 4: Customer
}

// ExampleDiagram_Section starts a stage of the journey. A journey needs one
// before any task, and a task added without one is reported from Build.
func ExampleDiagram_Section() {
	_ = userjourney.NewDiagram(os.Stdout).
		Section("Browse").
		Task("Search the catalogue", userjourney.ScoreSatisfied, "Customer").
		Section("Checkout").
		Task("Pay", userjourney.ScoreNeutral, "Customer", "Payment Service").
		Build()

	// Output:
	// journey
	//     section Browse
	//         Search the catalogue: 4: Customer
	//     section Checkout
	//         Pay: 3: Customer, Payment Service
}

// ExampleDiagram_Task adds a step to the current section, with how the actors
// felt about it and who was involved. The actors are optional.
func ExampleDiagram_Task() {
	_ = userjourney.NewDiagram(os.Stdout).
		Section("Checkout").
		Task("Add to basket", userjourney.ScoreVerySatisfied, "Customer").
		Task("Enter card details", userjourney.ScoreDissatisfied, "Customer", "Payment Service").
		Task("Confirm", userjourney.ScoreSatisfied).
		Build()

	// Output:
	// journey
	//     section Checkout
	//         Add to basket: 5: Customer
	//         Enter card details: 2: Customer, Payment Service
	//         Confirm: 4
}

// ExampleDiagram_TaskIn adds a step to a section named outright, which saves
// switching back and forth when the steps of two sections are interleaved in
// the calling code.
func ExampleDiagram_TaskIn() {
	_ = userjourney.NewDiagram(os.Stdout).
		Section("Browse").
		Section("Checkout").
		TaskIn("Browse", "Search the catalogue", userjourney.ScoreSatisfied, "Customer").
		TaskIn("Checkout", "Pay", userjourney.ScoreNeutral, "Customer").
		Build()

	// Output:
	// journey
	//     section Browse
	//     section Checkout
	//     section Browse
	//         Search the catalogue: 4: Customer
	//     section Checkout
	//         Pay: 3: Customer
}

// ExampleScore shows the five sentiments a task can carry. mermaid draws a
// happier face for a higher one, and the numbers it writes are 1 to 5.
func ExampleScore() {
	_ = userjourney.NewDiagram(os.Stdout).
		Section("Checkout").
		Task("Very dissatisfied", userjourney.ScoreVeryDissatisfied, "Customer").
		Task("Dissatisfied", userjourney.ScoreDissatisfied, "Customer").
		Task("Neutral", userjourney.ScoreNeutral, "Customer").
		Task("Satisfied", userjourney.ScoreSatisfied, "Customer").
		Task("Very satisfied", userjourney.ScoreVerySatisfied, "Customer").
		Build()

	// Output:
	// journey
	//     section Checkout
	//         Very dissatisfied: 1: Customer
	//         Dissatisfied: 2: Customer
	//         Neutral: 3: Customer
	//         Satisfied: 4: Customer
	//         Very satisfied: 5: Customer
}
