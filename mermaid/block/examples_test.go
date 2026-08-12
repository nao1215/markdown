//go:build linux || darwin

package block_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/block"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := block.NewDiagram(
		io.Discard,
		block.WithTitle("Checkout Architecture"),
	).
		Columns(3).
		Row(
			block.Node("Frontend"),
			block.ArrowRight("toBackend", block.WithArrowLabel("calls")),
			block.Node("Backend"),
		).
		Row(
			block.Space(2),
			block.ArrowDown("toDB"),
		).
		Row(
			block.Node("Database", block.WithNodeLabel("Primary DB"), block.WithNodeShape(block.ShapeCylinder)),
			block.Space(),
			block.Node("Cache", block.WithNodeLabel("Cache"), block.WithNodeShape(block.ShapeRound)),
		).
		Link("Backend", "Database").
		LinkWithLabel("Backend", "reads from", "Cache").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Block Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Block Diagram
	// ```mermaid
	// ---
	// title: "Checkout Architecture"
	// ---
	// block
	//     columns 3
	//     Frontend toBackend<["calls"]>(right) Backend
	//     space:2 toDB<["&nbsp;"]>(down)
	//     Database[("Primary DB")] space Cache("Cache")
	//     Backend --> Database
	//     Backend -- "reads from" --> Cache
	// ```
}

// ExampleNewDiagram shows the shape every block diagram has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a", block.WithNodeLabel("A"))).
		Build()

	// Output:
	// block
	//     a["A"]
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := block.NewDiagram(io.Discard).
		Row(block.Node("a", block.WithNodeLabel("A"))).
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// block
	//     a["A"]
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := block.NewDiagram(nil).
		Row(block.Node("a", block.WithNodeLabel("A"))).
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := block.NewDiagram(io.Discard).
		Columns(0)
	fmt.Println("error:", d.Error())

	// Output:
	// error: column count must be greater than zero
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a", block.WithNodeLabel("A"))).
		LF().
		Row(block.Node("a", block.WithNodeLabel("A"))).
		Build()

	// Output:
	// block
	//     a["A"]
	//
	//     a["A"]
}

// ExampleDiagram_full shows a block diagram built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := block.NewDiagram(io.Discard).
		Row(block.Node("a", block.WithNodeLabel("A"))).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// block
	//     a["A"]
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []block.Option{block.WithTitle("Overview")}

	_ = block.NewDiagram(os.Stdout, options...).
		Row(block.Node("a", block.WithNodeLabel("A"))).
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// block
	//     a["A"]
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = block.NewDiagram(os.Stdout, block.WithTitle("Overview")).
		Row(block.Node("a", block.WithNodeLabel("A"))).
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// block
	//     a["A"]
}

// ExampleDiagram_Row puts a line of blocks in the diagram. Everything on a row
// is drawn side by side.
func ExampleDiagram_Row() {
	_ = block.NewDiagram(os.Stdout).
		Row(
			block.Node("a", block.WithNodeLabel("Ingest")),
			block.Node("b", block.WithNodeLabel("Process")),
		).
		Build()

	// Output:
	// block
	//     a["Ingest"] b["Process"]
}

// ExampleDiagram_Columns fixes how many columns the diagram is laid out in, so
// a long row wraps where the caller wants it to rather than where mermaid does.
func ExampleDiagram_Columns() {
	_ = block.NewDiagram(os.Stdout).
		Columns(2).
		Row(
			block.Node("a", block.WithNodeLabel("A")),
			block.Node("b", block.WithNodeLabel("B")),
			block.Node("c", block.WithNodeLabel("C")),
		).
		Build()

	// Output:
	// block
	//     columns 2
	//     a["A"] b["B"] c["C"]
}

// ExampleDiagram_Block opens a group of blocks drawn inside a box. The calls in
// the function belong to it.
func ExampleDiagram_Block() {
	_ = block.NewDiagram(os.Stdout).
		Block(func(d *block.Diagram) {
			d.Row(block.Node("a", block.WithNodeLabel("Inside")))
		}, block.WithBlockID("group")).
		Build()

	// Output:
	// block
	//     block:group
	//         a["Inside"]
	//     end
}

// ExampleDiagram_Link joins two blocks by their identifiers.
func ExampleDiagram_Link() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.Node("b")).
		Link("a", "b").
		Build()

	// Output:
	// block
	//     a b
	//     a --> b
}

// ExampleDiagram_LinkWithLabel joins two blocks and says what the line means.
func ExampleDiagram_LinkWithLabel() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.Node("b")).
		LinkWithLabel("a", "publishes", "b").
		Build()

	// Output:
	// block
	//     a b
	//     a -- "publishes" --> b
}

// ExampleDiagram_Statement writes a line of mermaid through unchanged, for the
// parts of the block syntax this package has no method for.
func ExampleDiagram_Statement() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a")).
		Statement("style a fill:#f9f").
		Build()

	// Output:
	// block
	//     a
	//     style a fill:#f9f
}

// ExampleDiagram_Style colors one block outright.
func ExampleDiagram_Style() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a")).
		Style("a", "fill:#f9f,stroke:#333").
		Build()

	// Output:
	// block
	//     a
	//     style a fill:#f9f,stroke:#333
}

// ExampleDiagram_ClassDef names a style so several blocks can share it.
func ExampleDiagram_ClassDef() {
	_ = block.NewDiagram(os.Stdout).
		ClassDef("warning", "fill:#f96").
		Row(block.Node("a"), block.Node("b")).
		Class("a,b", "warning").
		Build()

	// Output:
	// block
	//     classDef warning fill:#f96
	//     a b
	//     class a,b warning
}

// ExampleDiagram_Class applies a named style to blocks.
func ExampleDiagram_Class() {
	_ = block.NewDiagram(os.Stdout).
		ClassDef("warning", "fill:#f96").
		Row(block.Node("a")).
		Class("a", "warning").
		Build()

	// Output:
	// block
	//     classDef warning fill:#f96
	//     a
	//     class a warning
}

// ExampleNode is one block. Without a label it is drawn with its identifier
// inside, which is often all a diagram needs.
func ExampleNode() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.Node("b", block.WithNodeLabel("With a label"))).
		Build()

	// Output:
	// block
	//     a b["With a label"]
}

// ExampleLiteral writes a token through unchanged, for a shape this package has
// no option for.
func ExampleLiteral() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Literal(`a(("Circle"))`)).
		Build()

	// Output:
	// block
	//     a(("Circle"))
}

// ExampleSpace leaves a gap in a row, which is how a diagram lines up two rows
// that hold different numbers of blocks.
func ExampleSpace() {
	_ = block.NewDiagram(os.Stdout).
		Columns(3).
		Row(block.Node("a"), block.Space(), block.Node("b")).
		Row(block.Node("c"), block.Space(2), block.Node("d")).
		Build()

	// Output:
	// block
	//     columns 3
	//     a space b
	//     c space:2 d
}

// ExampleArrow draws an arrow token pointing the way it is told.
func ExampleArrow() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.Arrow("ar", block.DirectionRight), block.Node("b")).
		Build()

	// Output:
	// block
	//     a ar<["&nbsp;"]>(right) b
}

// ExampleArrowRight draws an arrow pointing right, which is the same as Arrow
// with DirectionRight and shorter to read in a chain.
func ExampleArrowRight() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.ArrowRight("ar"), block.Node("b")).
		Build()

	// Output:
	// block
	//     a ar<["&nbsp;"]>(right) b
}

// ExampleArrowLeft draws an arrow pointing left.
func ExampleArrowLeft() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.ArrowLeft("al"), block.Node("b")).
		Build()

	// Output:
	// block
	//     a al<["&nbsp;"]>(left) b
}

// ExampleArrowUp draws an arrow pointing up.
func ExampleArrowUp() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.ArrowUp("au"), block.Node("b")).
		Build()

	// Output:
	// block
	//     a au<["&nbsp;"]>(up) b
}

// ExampleArrowDown draws an arrow pointing down.
func ExampleArrowDown() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.ArrowDown("ad"), block.Node("b")).
		Build()

	// Output:
	// block
	//     a ad<["&nbsp;"]>(down) b
}

// ExampleArrowX draws an arrow pointing both left and right.
func ExampleArrowX() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.ArrowX("ax"), block.Node("b")).
		Build()

	// Output:
	// block
	//     a ax<["&nbsp;"]>(x) b
}

// ExampleArrowY draws an arrow pointing both up and down.
func ExampleArrowY() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.ArrowY("ay"), block.Node("b")).
		Build()

	// Output:
	// block
	//     a ay<["&nbsp;"]>(y) b
}

// ExampleWithNodeLabel puts text inside a block instead of its identifier.
func ExampleWithNodeLabel() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a", block.WithNodeLabel("Ingest"))).
		Build()

	// Output:
	// block
	//     a["Ingest"]
}

// ExampleWithNodeShape draws a block as something other than a rectangle.
func ExampleWithNodeShape() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a", block.WithNodeLabel("Decide"), block.WithNodeShape(block.ShapeRhombus))).
		Build()

	// Output:
	// block
	//     a{"Decide"}
}

// ExampleWithNodeSpan makes one block as wide as several columns.
func ExampleWithNodeSpan() {
	_ = block.NewDiagram(os.Stdout).
		Columns(3).
		Row(block.Node("a", block.WithNodeLabel("Wide"), block.WithNodeSpan(2)), block.Node("b")).
		Build()

	// Output:
	// block
	//     columns 3
	//     a["Wide"]:2 b
}

// ExampleWithBlockID names a group of blocks.
func ExampleWithBlockID() {
	_ = block.NewDiagram(os.Stdout).
		Block(func(d *block.Diagram) {
			d.Row(block.Node("a"))
		}, block.WithBlockID("group")).
		Build()

	// Output:
	// block
	//     block:group
	//         a
	//     end
}

// ExampleWithBlockSpan makes a group as wide as several columns.
func ExampleWithBlockSpan() {
	_ = block.NewDiagram(os.Stdout).
		Columns(3).
		Block(func(d *block.Diagram) {
			d.Row(block.Node("a"))
		}, block.WithBlockID("group"), block.WithBlockSpan(2)).
		Build()

	// Output:
	// block
	//     columns 3
	//     block:group:2
	//         a
	//     end
}

// ExampleWithArrowLabel says what an arrow means.
func ExampleWithArrowLabel() {
	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.ArrowRight("ar", block.WithArrowLabel("publishes")), block.Node("b")).
		Build()

	// Output:
	// block
	//     a ar<["publishes"]>(right) b
}

// ExampleWithArrowSecondaryDirection points an arrow two ways at once.
func ExampleWithArrowSecondaryDirection() {
	_ = block.NewDiagram(os.Stdout).
		Row(
			block.Node("a"),
			block.ArrowRight("ar", block.WithArrowSecondaryDirection(block.DirectionDown)),
			block.Node("b"),
		).
		Build()

	// Output:
	// block
	//     a ar<["&nbsp;"]>(right, down) b
}

// ExampleToken is what a row is made of: a node, an arrow, a gap or a literal.
func ExampleToken() {
	tokens := []block.Token{
		block.Node("a", block.WithNodeLabel("Ingest")),
		block.ArrowRight("ar"),
		block.Node("b", block.WithNodeLabel("Process")),
	}

	_ = block.NewDiagram(os.Stdout).Row(tokens...).Build()

	// Output:
	// block
	//     a["Ingest"] ar<["&nbsp;"]>(right) b["Process"]
}

// ExampleShape shows the shapes a block can be drawn as.
func ExampleShape() {
	_ = block.NewDiagram(os.Stdout).
		Row(
			block.Node("a", block.WithNodeShape(block.ShapeRound)),
			block.Node("b", block.WithNodeShape(block.ShapeStadium)),
			block.Node("c", block.WithNodeShape(block.ShapeRhombus)),
		).
		Build()

	// Output:
	// block
	//     a("a") b(["b"]) c{"c"}
}

// ExampleDirection shows the ways an arrow can point.
func ExampleDirection() {
	_ = block.NewDiagram(os.Stdout).
		Row(
			block.Arrow("a", block.DirectionRight),
			block.Arrow("b", block.DirectionLeft),
			block.Arrow("c", block.DirectionUp),
			block.Arrow("d", block.DirectionDown),
		).
		Build()

	// Output:
	// block
	//     a<["&nbsp;"]>(right) b<["&nbsp;"]>(left) c<["&nbsp;"]>(up) d<["&nbsp;"]>(down)
}

// ExampleNodeOption shows what a NodeOption is: a function that changes how a
// block is written, passed to Node after its identifier.
func ExampleNodeOption() {
	options := []block.NodeOption{block.WithNodeLabel("Ingest"), block.WithNodeSpan(2)}

	_ = block.NewDiagram(os.Stdout).Columns(3).Row(block.Node("a", options...)).Build()

	// Output:
	// block
	//     columns 3
	//     a["Ingest"]:2
}

// ExampleBlockOption shows what a BlockOption is: a function that changes how a
// group is written, passed to Block after the function that fills it.
func ExampleBlockOption() {
	options := []block.BlockOption{block.WithBlockID("group")}

	_ = block.NewDiagram(os.Stdout).
		Block(func(d *block.Diagram) { d.Row(block.Node("a")) }, options...).
		Build()

	// Output:
	// block
	//     block:group
	//         a
	//     end
}

// ExampleArrowOption shows what an ArrowOption is: a function that changes how
// an arrow is written, passed to one of the arrow constructors.
func ExampleArrowOption() {
	options := []block.ArrowOption{block.WithArrowLabel("publishes")}

	_ = block.NewDiagram(os.Stdout).
		Row(block.Node("a"), block.ArrowRight("ar", options...), block.Node("b")).
		Build()

	// Output:
	// block
	//     a ar<["publishes"]>(right) b
}
