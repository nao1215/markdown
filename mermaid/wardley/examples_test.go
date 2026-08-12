//go:build linux || darwin

package wardley_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/wardley"
)

// ExampleMap skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleMap() {
	diagram := wardley.NewMap(io.Discard, wardley.WithTitle("Checkout")).
		Anchor("Customer", 0.95, 0.95).
		Component("Checkout", 0.6, 0.8).
		Component("Card network", 0.95, 0.2).
		Link("Customer", "Checkout").
		Link("Checkout", "Card network").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Wardley map").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Wardley map
	// ```mermaid
	// wardley-beta
	//     title Checkout
	//     anchor Customer [0.95, 0.95]
	//     component Checkout [0.6, 0.8]
	//     component Card network [0.95, 0.2]
	//     Customer -> Checkout
	//     Checkout -> Card network
	// ```
}

// ExampleNewMap shows the shape every Wardley map has: a writer, the parts, the
// dependencies between them, and Build.
func ExampleNewMap() {
	_ = wardley.NewMap(os.Stdout).
		Anchor("Customer", 0.95, 0.95).
		Component("Checkout", 0.6, 0.8).
		Link("Customer", "Checkout").
		Build()

	// Output:
	// wardley-beta
	//     anchor Customer [0.95, 0.95]
	//     component Checkout [0.6, 0.8]
	//     Customer -> Checkout
}

// ExampleMap_Component places a part of the system. Evolution runs left to
// right, from something built for the first time to something bought as a
// commodity; visibility runs bottom to top, from the plumbing to what the user
// actually touches.
func ExampleMap_Component() {
	_ = wardley.NewMap(os.Stdout).
		Component("Checkout", 0.6, 0.8).
		Component("Card network", 0.95, 0.2).
		Build()

	// Output:
	// wardley-beta
	//     component Checkout [0.6, 0.8]
	//     component Card network [0.95, 0.2]
}

// ExampleMap_Anchor places the user the map is drawn from the point of view of.
// A map usually has one, near the top, with everything else below it.
func ExampleMap_Anchor() {
	_ = wardley.NewMap(os.Stdout).
		Anchor("Customer", 0.95, 0.95).
		Build()

	// Output:
	// wardley-beta
	//     anchor Customer [0.95, 0.95]
}

// ExampleMap_Link draws a dependency: the first named part needs the second.
func ExampleMap_Link() {
	_ = wardley.NewMap(os.Stdout).
		Component("Checkout", 0.6, 0.8).
		Component("Payment", 0.75, 0.5).
		Link("Checkout", "Payment").
		Build()

	// Output:
	// wardley-beta
	//     component Checkout [0.6, 0.8]
	//     component Payment [0.75, 0.5]
	//     Checkout -> Payment
}

// ExampleMap_Evolve marks a part as moving along the evolution axis, drawn as
// an arrow from where it is to where it is going. It is what turns a map of
// today into an argument about tomorrow.
func ExampleMap_Evolve() {
	_ = wardley.NewMap(os.Stdout).
		Component("Payment", 0.75, 0.5).
		Evolve("Payment", 0.9).
		Build()

	// Output:
	// wardley-beta
	//     component Payment [0.75, 0.5]
	//     evolve Payment 0.9
}

// ExampleMap_String returns the map without needing a writer, which is how it
// is handed to a markdown code block.
func ExampleMap_String() {
	diagram := wardley.NewMap(io.Discard).Component("Checkout", 0.6, 0.8).String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// wardley-beta
	//     component Checkout [0.6, 0.8]
	// ```
}

// ExampleMap_Build writes the map and reports the first error the chain
// recorded.
func ExampleMap_Build() {
	err := wardley.NewMap(nil).Component("Checkout", 0.6, 0.8).Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleMap_Error reports the error the chain recorded. A name mermaid cannot
// read is reported rather than mangled, because a Wardley map writes names
// unquoted and accepts no escape in that position.
func ExampleMap_Error() {
	m := wardley.NewMap(io.Discard).Component("checkout.v2", 0.6, 0.8)
	fmt.Println("error:", m.Error())

	// Output:
	// error: name "checkout.v2" must hold only letters, digits, spaces, underscores, hyphens and parentheses; mermaid reads nothing else there
}

// ExampleMap_LF adds a blank line to the map body.
func ExampleMap_LF() {
	_ = wardley.NewMap(os.Stdout).
		Component("Checkout", 0.6, 0.8).
		LF().
		Component("Payment", 0.75, 0.5).
		Build()

	// Output:
	// wardley-beta
	//     component Checkout [0.6, 0.8]
	//
	//     component Payment [0.75, 0.5]
}

// ExampleWithTitle sets the title the map is drawn with. It is written through
// unchanged: mermaid takes every character probed there, quotation marks and
// hashes included.
func ExampleWithTitle() {
	_ = wardley.NewMap(os.Stdout, wardley.WithTitle(`Checkout "as it stands" #1`)).
		Component("Checkout", 0.6, 0.8).
		Build()

	// Output:
	// wardley-beta
	//     title Checkout "as it stands" #1
	//     component Checkout [0.6, 0.8]
}

// ExampleOption shows what an Option is: a function that changes how the map is
// written, passed to NewMap.
func ExampleOption() {
	options := []wardley.Option{wardley.WithTitle("Checkout")}

	_ = wardley.NewMap(os.Stdout, options...).Component("Checkout", 0.6, 0.8).Build()

	// Output:
	// wardley-beta
	//     title Checkout
	//     component Checkout [0.6, 0.8]
}
