//go:build linux || darwin

package flowchart_test

import (
	"fmt"
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/flowchart"
)

// ExampleFlowchart skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleFlowchart() {
	fc := flowchart.NewFlowchart(
		os.Stdout,
		flowchart.WithTitle("mermaid flowchart builder"),
		flowchart.WithOrientalTopToBottom(),
	).
		NodeWithText("A", "Node A").
		StadiumNode("B", "Node B").
		SubroutineNode("C", "Node C").
		DatabaseNode("D", "Database").
		LinkWithArrowHead("A", "B").
		LinkWithArrowHeadAndText("B", "D", "send original data").
		LinkWithArrowHead("B", "C").
		DottedLinkWithText("C", "D", "send filtered data").
		String()

	_ = markdown.NewMarkdown(os.Stdout).
		H2("Flowchart").
		CodeBlocks(markdown.SyntaxHighlightMermaid, fc).
		Build()

	// Output:
	// ## Flowchart
	// ```mermaid
	// ---
	// title: "mermaid flowchart builder"
	// ---
	// flowchart TB
	//     A["Node A"]
	//     B(["Node B"])
	//     C[["Node C"]]
	//     D[("Database")]
	//     A-->B
	//     B-->|"send original data"|D
	//     B-->C
	//     C-. "send filtered data" .-> D
	// ```
}

// ExampleFlowchart_NodeWithText draws a rectangle, the default shape.
func ExampleFlowchart_NodeWithText() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A["Text"]
}

// ExampleFlowchart_RoundEdgesNode draws a rectangle with rounded corners.
func ExampleFlowchart_RoundEdgesNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		RoundEdgesNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A("Text")
}

// ExampleFlowchart_StadiumNode draws a stadium, which is a rectangle with round
// ends.
func ExampleFlowchart_StadiumNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		StadiumNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A(["Text"])
}

// ExampleFlowchart_SubroutineNode draws a subroutine, which carries a double
// vertical edge.
func ExampleFlowchart_SubroutineNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		SubroutineNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A[["Text"]]
}

// ExampleFlowchart_CylindricalNode draws a cylinder, which is how a database is
// usually drawn.
func ExampleFlowchart_CylindricalNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		CylindricalNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A[("Text")]
}

// ExampleFlowchart_DatabaseNode draws a cylinder. It is the same as
// CylindricalNode, under a name that says what one is usually for.
func ExampleFlowchart_DatabaseNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		DatabaseNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A[("Text")]
}

// ExampleFlowchart_CircleNode draws a circle.
func ExampleFlowchart_CircleNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		CircleNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A(("Text"))
}

// ExampleFlowchart_DoubleCircleNode draws a double circle, which usually marks
// where a flow ends.
func ExampleFlowchart_DoubleCircleNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		DoubleCircleNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A((("Text")))
}

// ExampleFlowchart_AsymmetricNode draws a flag, flat on one side and pointed on
// the other.
func ExampleFlowchart_AsymmetricNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		AsymmetricNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A>"Text"]
}

// ExampleFlowchart_RhombusNode draws a diamond, which is how a decision is
// usually drawn.
func ExampleFlowchart_RhombusNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		RhombusNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A{"Text"}
}

// ExampleFlowchart_HexagonNode draws a hexagon.
func ExampleFlowchart_HexagonNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		HexagonNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A{{"Text"}}
}

// ExampleFlowchart_ParallelogramNode draws a parallelogram leaning right, which
// is how input or output is usually drawn.
func ExampleFlowchart_ParallelogramNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		ParallelogramNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A[/"Text"/]
}

// ExampleFlowchart_ParallelogramAltNode draws a parallelogram leaning the other way.
func ExampleFlowchart_ParallelogramAltNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		ParallelogramAltNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A[\"Text"\]
}

// ExampleFlowchart_TrapezoidNode draws a trapezoid, narrower at the top.
func ExampleFlowchart_TrapezoidNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		TrapezoidNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A[/"Text"\]
}

// ExampleFlowchart_TrapezoidAltNode draws a trapezoid, narrower at the bottom.
func ExampleFlowchart_TrapezoidAltNode() {
	_ = flowchart.NewFlowchart(os.Stdout).
		TrapezoidAltNode("A", "Text").
		Build()

	// Output:
	// flowchart TB
	//     A[\"Text"/]
}

// ExampleFlowchart_LinkWithArrowHead joins two nodes with an arrow.
func ExampleFlowchart_LinkWithArrowHead() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("A", "Start").
		NodeWithText("B", "End").
		LinkWithArrowHead("A", "B").
		Build()

	// Output:
	// flowchart TB
	//     A["Start"]
	//     B["End"]
	//     A-->B
}

// ExampleFlowchart_LinkWithArrowHeadAndText joins two nodes with an arrow and
// says what it means.
func ExampleFlowchart_LinkWithArrowHeadAndText() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("A", "Start").
		NodeWithText("B", "End").
		LinkWithArrowHeadAndText("A", "B", "then").
		Build()

	// Output:
	// flowchart TB
	//     A["Start"]
	//     B["End"]
	//     A-->|"then"|B
}

// ExampleFlowchart_OpenLink joins two nodes with a plain line, no arrowhead.
func ExampleFlowchart_OpenLink() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("A", "Start").
		NodeWithText("B", "End").
		OpenLink("A", "B").
		Build()

	// Output:
	// flowchart TB
	//     A["Start"]
	//     B["End"]
	//     A --- B
}

// ExampleFlowchart_OpenLinkWithText joins two nodes with a plain line and a label.
func ExampleFlowchart_OpenLinkWithText() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("A", "Start").
		NodeWithText("B", "End").
		OpenLinkWithText("A", "B", "and").
		Build()

	// Output:
	// flowchart TB
	//     A["Start"]
	//     B["End"]
	//     A---|"and"|B
}

// ExampleFlowchart_DottedLink joins two nodes with a dotted line, which usually
// means a weaker relationship than a solid one.
func ExampleFlowchart_DottedLink() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("A", "Start").
		NodeWithText("B", "End").
		DottedLink("A", "B").
		Build()

	// Output:
	// flowchart TB
	//     A["Start"]
	//     B["End"]
	//     A-.->B
}

// ExampleFlowchart_DottedLinkWithText joins two nodes with a dotted line and a label.
func ExampleFlowchart_DottedLinkWithText() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("A", "Start").
		NodeWithText("B", "End").
		DottedLinkWithText("A", "B", "sometimes").
		Build()

	// Output:
	// flowchart TB
	//     A["Start"]
	//     B["End"]
	//     A-. "sometimes" .-> B
}

// ExampleFlowchart_ThickLink joins two nodes with a thick line, for the path
// through a flow that matters most.
func ExampleFlowchart_ThickLink() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("A", "Start").
		NodeWithText("B", "End").
		ThickLink("A", "B").
		Build()

	// Output:
	// flowchart TB
	//     A["Start"]
	//     B["End"]
	//     A ==> B
}

// ExampleFlowchart_ThickLinkWithText joins two nodes with a thick line and a label.
func ExampleFlowchart_ThickLinkWithText() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("A", "Start").
		NodeWithText("B", "End").
		ThickLinkWithText("A", "B", "always").
		Build()

	// Output:
	// flowchart TB
	//     A["Start"]
	//     B["End"]
	//     A == "always" ==> B
}

// ExampleFlowchart_InvisibleLink joins two nodes with a line that is not drawn,
// which is how a diagram is pushed into the layout its author wants without
// saying anything untrue about the flow.
func ExampleFlowchart_InvisibleLink() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("A", "Start").
		NodeWithText("B", "End").
		InvisibleLink("A", "B").
		Build()

	// Output:
	// flowchart TB
	//     A["Start"]
	//     B["End"]
	//     A ~~~ B
}

// ExampleWithOrientalTopDown lays the flowchart out from the top downwards, which is the default.
func ExampleWithOrientalTopDown() {
	_ = flowchart.NewFlowchart(os.Stdout, flowchart.WithOrientalTopDown()).
		NodeWithText("A", "Start").
		Build()

	// Output:
	// flowchart TD
	//     A["Start"]
}

// ExampleWithOrientalTopToBottom lays the flowchart out from the top to the
// bottom. It is the same as WithOrientalTopDown under mermaid's other name for
// it.
func ExampleWithOrientalTopToBottom() {
	_ = flowchart.NewFlowchart(os.Stdout, flowchart.WithOrientalTopToBottom()).
		NodeWithText("A", "Start").
		Build()

	// Output:
	// flowchart TB
	//     A["Start"]
}

// ExampleWithOrientalBottomToTop lays the flowchart out from the bottom upwards.
func ExampleWithOrientalBottomToTop() {
	_ = flowchart.NewFlowchart(os.Stdout, flowchart.WithOrientalBottomToTop()).
		NodeWithText("A", "Start").
		Build()

	// Output:
	// flowchart BT
	//     A["Start"]
}

// ExampleWithOrientalLeftToRight lays the flowchart out from the left to right.
func ExampleWithOrientalLeftToRight() {
	_ = flowchart.NewFlowchart(os.Stdout, flowchart.WithOrientalLeftToRight()).
		NodeWithText("A", "Start").
		Build()

	// Output:
	// flowchart LR
	//     A["Start"]
}

// ExampleWithOrientalRightToLeft lays the flowchart out from the right to left.
func ExampleWithOrientalRightToLeft() {
	_ = flowchart.NewFlowchart(os.Stdout, flowchart.WithOrientalRightToLeft()).
		NodeWithText("A", "Start").
		Build()

	// Output:
	// flowchart RL
	//     A["Start"]
}

// ExampleNewFlowchart shows the shape every flowchart has: a writer, a chain of
// calls, and Build.
func ExampleNewFlowchart() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("A", "Start").
		NodeWithText("B", "End").
		LinkWithArrowHead("A", "B").
		Build()

	// Output:
	// flowchart TB
	//     A["Start"]
	//     B["End"]
	//     A-->B
}

// ExampleFlowchart_Node declares a node by its identifier alone, for one whose
// text is set by a later call or which needs none.
func ExampleFlowchart_Node() {
	_ = flowchart.NewFlowchart(os.Stdout).
		Node("A").
		Build()

	// Output:
	// flowchart TB
	//     A
}

// ExampleFlowchart_NodeWithMarkdown draws a node whose text is markdown, so it
// can carry bold or a link.
func ExampleFlowchart_NodeWithMarkdown() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithMarkdown("A", "**Start** here").
		Build()

	// Output:
	// flowchart TB
	//     A["`**Start** here`"]
}

// ExampleFlowchart_NodeWithNewLines draws a node whose text runs over several
// lines. It is the same form as NodeWithMarkdown, which is what mermaid honors
// a line break inside.
func ExampleFlowchart_NodeWithNewLines() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithNewLines("A", "Start\nhere").
		Build()

	// Output:
	// flowchart TB
	//     A["`Start
	// here`"]
}

// ExampleFlowchart_String returns the flowchart without needing a writer, which
// is how it is handed to a markdown code block.
func ExampleFlowchart_String() {
	diagram := flowchart.NewFlowchart(io.Discard).
		NodeWithText("A", "Start").
		String()

	_ = markdown.NewMarkdown(os.Stdout).
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// flowchart TB
	//     A["Start"]
	// ```
}

// ExampleFlowchart_Build writes the flowchart and reports the error the chain
// recorded.
func ExampleFlowchart_Build() {
	err := flowchart.NewFlowchart(nil).NodeWithText("A", "Start").Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleWithTitle sets the title the flowchart is drawn with.
func ExampleWithTitle() {
	_ = flowchart.NewFlowchart(os.Stdout, flowchart.WithTitle("Checkout")).
		NodeWithText("A", "Start").
		Build()

	// Output:
	// ---
	// title: "Checkout"
	// ---
	// flowchart TB
	//     A["Start"]
}

// ExampleOption shows what an Option is: a function that changes how the
// flowchart is written, passed to NewFlowchart.
func ExampleOption() {
	options := []flowchart.Option{
		flowchart.WithTitle("Checkout"),
		flowchart.WithOrientalLeftToRight(),
	}

	_ = flowchart.NewFlowchart(os.Stdout, options...).
		NodeWithText("A", "Start").
		Build()

	// Output:
	// ---
	// title: "Checkout"
	// ---
	// flowchart LR
	//     A["Start"]
}

// ExampleFlowchart_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleFlowchart_Error() {
	f := flowchart.NewFlowchart(nil).NodeWithText("A", "Start")
	fmt.Println("before Build:", f.Error())
	_ = f.Build()
	fmt.Println("after Build:", f.Error())

	// Output:
	// before Build: <nil>
	// after Build: output writer must not be nil
}

// ExampleFlowchart_Subgraph groups the nodes that follow it into a box.
// SubgraphEnd closes it, and what lies between is indented.
func ExampleFlowchart_Subgraph() {
	_ = flowchart.NewFlowchart(os.Stdout).
		Subgraph("ingest", "Ingest").
		NodeWithText("a", "Fetch").
		NodeWithText("b", "Parse").
		SubgraphEnd().
		NodeWithText("c", "Store").
		Build()

	// Output:
	// flowchart TB
	//     subgraph ingest["Ingest"]
	//         a["Fetch"]
	//         b["Parse"]
	//     end
	//     c["Store"]
}

// ExampleFlowchart_SubgraphEnd closes the subgraph opened last. Leaving one
// open is reported from Build, because mermaid refuses a flowchart whose
// subgraph never ends.
func ExampleFlowchart_SubgraphEnd() {
	_ = flowchart.NewFlowchart(os.Stdout).
		Subgraph("outer", "Outer").
		Subgraph("inner", "Inner").
		NodeWithText("a", "A").
		SubgraphEnd().
		SubgraphEnd().
		Build()

	// Output:
	// flowchart TB
	//     subgraph outer["Outer"]
	//         subgraph inner["Inner"]
	//             a["A"]
	//         end
	//     end
}

// ExampleFlowchart_SubgraphDirection lays one subgraph out across the page
// while the chart around it runs down, which is the readable arrangement for a
// row of steps inside a longer flow.
func ExampleFlowchart_SubgraphDirection() {
	_ = flowchart.NewFlowchart(os.Stdout).
		Subgraph("ingest", "Ingest").
		SubgraphDirection(flowchart.DirectionLR).
		NodeWithText("a", "Fetch").
		NodeWithText("b", "Parse").
		LinkWithArrowHead("a", "b").
		SubgraphEnd().
		Build()

	// Output:
	// flowchart TB
	//     subgraph ingest["Ingest"]
	//         direction LR
	//         a["Fetch"]
	//         b["Parse"]
	//         a-->b
	//     end
}

// ExampleDirection shows the four ways a subgraph can be laid out. The chart's
// own direction is an option on NewFlowchart instead.
func ExampleDirection() {
	for _, direction := range []flowchart.Direction{
		flowchart.DirectionTB,
		flowchart.DirectionBT,
		flowchart.DirectionLR,
		flowchart.DirectionRL,
	} {
		_ = flowchart.NewFlowchart(os.Stdout).
			Subgraph("g", "Group").
			SubgraphDirection(direction).
			SubgraphEnd().
			Build()
		fmt.Println()
	}

	// Output:
	// flowchart TB
	//     subgraph g["Group"]
	//         direction TB
	//     end
	// flowchart TB
	//     subgraph g["Group"]
	//         direction BT
	//     end
	// flowchart TB
	//     subgraph g["Group"]
	//         direction LR
	//     end
	// flowchart TB
	//     subgraph g["Group"]
	//         direction RL
	//     end
}

// ExampleFlowchart_Style colors one node outright. The style is mermaid's own
// CSS-like syntax, written through unchanged.
func ExampleFlowchart_Style() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("a", "Start").
		Style("a", "fill:#f9f,stroke:#333").
		Build()

	// Output:
	// flowchart TB
	//     a["Start"]
	//     style a fill:#f9f,stroke:#333
}

// ExampleFlowchart_ClassDef names a style so that several nodes can share it.
func ExampleFlowchart_ClassDef() {
	_ = flowchart.NewFlowchart(os.Stdout).
		ClassDef("urgent", "fill:#f96,stroke:#333").
		NodeWithText("a", "Start").
		Class("a", "urgent").
		Build()

	// Output:
	// flowchart TB
	//     classDef urgent fill:#f96,stroke:#333
	//     a["Start"]
	//     class a urgent
}

// ExampleFlowchart_Class applies a named style to nodes. Several are given as
// one comma separated list, which is what mermaid reads there.
func ExampleFlowchart_Class() {
	_ = flowchart.NewFlowchart(os.Stdout).
		ClassDef("urgent", "fill:#f96").
		NodeWithText("a", "Start").
		NodeWithText("b", "End").
		Class("a,b", "urgent").
		Build()

	// Output:
	// flowchart TB
	//     classDef urgent fill:#f96
	//     a["Start"]
	//     b["End"]
	//     class a,b urgent
}

// ExampleFlowchart_ClickHref makes a node a link, with the text a browser shows
// on hover. The tooltip is escaped the way a label is.
func ExampleFlowchart_ClickHref() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("a", "Order").
		ClickHref("a", "https://example.com/order", "The Order type").
		Build()

	// Output:
	// flowchart TB
	//     a["Order"]
	//     click a "https://example.com/order" "The Order type"
}

// ExampleFlowchart_ClickCall makes a node call a function in the page when it
// is clicked. The parentheses are added when the caller leaves them off.
func ExampleFlowchart_ClickCall() {
	_ = flowchart.NewFlowchart(os.Stdout).
		NodeWithText("a", "Order").
		ClickCall("a", "showOrder", "Show the order").
		Build()

	// Output:
	// flowchart TB
	//     a["Order"]
	//     click a call showOrder() "Show the order"
}
