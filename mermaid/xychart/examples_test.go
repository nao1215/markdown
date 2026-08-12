//go:build linux || darwin

package xychart_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/xychart"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := xychart.NewDiagram(
		io.Discard,
		xychart.WithTitle("Sales Revenue"),
	).
		XAxisLabels("Jan", "Feb", "Mar", "Apr", "May", "Jun").
		YAxisRangeWithTitle("Revenue (k$)", 0, 100).
		Bar(25, 40, 60, 80, 70, 90).
		Line(30, 50, 70, 85, 75, 95).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("XY Chart").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## XY Chart
	// ```mermaid
	// xychart
	//     title "Sales Revenue"
	//     x-axis [Jan, Feb, Mar, Apr, May, Jun]
	//     y-axis "Revenue (k$)" 0 --> 100
	//     bar [25, 40, 60, 80, 70, 90]
	//     line [30, 50, 70, 85, 75, 95]
	// ```
}

// ExampleNewDiagram shows the shape every xy chart has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = xychart.NewDiagram(os.Stdout).
		XAxisLabelsWithTitle("Month", "Jan", "Feb").Bar(10, 20).
		Build()

	// Output:
	// xychart
	//     x-axis Month [Jan, Feb]
	//     bar [10, 20]
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := xychart.NewDiagram(io.Discard).
		XAxisLabelsWithTitle("Month", "Jan", "Feb").Bar(10, 20).
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// xychart
	//     x-axis Month [Jan, Feb]
	//     bar [10, 20]
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := xychart.NewDiagram(nil).
		XAxisLabelsWithTitle("Month", "Jan", "Feb").Bar(10, 20).
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := xychart.NewDiagram(io.Discard).
		XAxisRange(10, 0)
	fmt.Println("error:", d.Error())

	// Output:
	// error: x-axis range requires min to be less than max
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = xychart.NewDiagram(os.Stdout).
		XAxisLabelsWithTitle("Month", "Jan", "Feb").Bar(10, 20).
		LF().
		XAxisLabelsWithTitle("Month", "Jan", "Feb").Bar(10, 20).
		Build()

	// Output:
	// xychart
	//     x-axis Month [Jan, Feb]
	//     bar [10, 20]
	//
	//     x-axis Month [Jan, Feb]
	//     bar [10, 20]
}

// ExampleDiagram_full shows a xy chart built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := xychart.NewDiagram(io.Discard).
		XAxisLabelsWithTitle("Month", "Jan", "Feb").Bar(10, 20).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// xychart
	//     x-axis Month [Jan, Feb]
	//     bar [10, 20]
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []xychart.Option{xychart.WithTitle("Overview")}

	_ = xychart.NewDiagram(os.Stdout, options...).
		XAxisLabelsWithTitle("Month", "Jan", "Feb").Bar(10, 20).
		Build()

	// Output:
	// xychart
	//     title "Overview"
	//     x-axis Month [Jan, Feb]
	//     bar [10, 20]
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = xychart.NewDiagram(os.Stdout, xychart.WithTitle("Overview")).
		XAxisLabelsWithTitle("Month", "Jan", "Feb").Bar(10, 20).
		Build()

	// Output:
	// xychart
	//     title "Overview"
	//     x-axis Month [Jan, Feb]
	//     bar [10, 20]
}

// ExampleDiagram_XAxisLabels names the points along the x axis, for data that
// is counted by category rather than measured on a scale.
func ExampleDiagram_XAxisLabels() {
	_ = xychart.NewDiagram(os.Stdout).
		XAxisLabels("Jan", "Feb", "Mar").
		Bar(10, 20, 30).
		Build()

	// Output:
	// xychart
	//     x-axis [Jan, Feb, Mar]
	//     bar [10, 20, 30]
}

// ExampleDiagram_XAxisLabelsWithTitle names the points and the axis itself.
func ExampleDiagram_XAxisLabelsWithTitle() {
	_ = xychart.NewDiagram(os.Stdout).
		XAxisLabelsWithTitle("Month", "Jan", "Feb", "Mar").
		Bar(10, 20, 30).
		Build()

	// Output:
	// xychart
	//     x-axis Month [Jan, Feb, Mar]
	//     bar [10, 20, 30]
}

// ExampleDiagram_XAxisRange measures the x axis on a scale instead, which is
// what a chart of a continuous quantity needs.
func ExampleDiagram_XAxisRange() {
	_ = xychart.NewDiagram(os.Stdout).
		XAxisRange(0, 100).
		Line(10, 20, 30).
		Build()

	// Output:
	// xychart
	//     x-axis 0 --> 100
	//     line [10, 20, 30]
}

// ExampleDiagram_XAxisRangeWithTitle measures the x axis and names it.
func ExampleDiagram_XAxisRangeWithTitle() {
	_ = xychart.NewDiagram(os.Stdout).
		XAxisRangeWithTitle("Load", 0, 100).
		Line(10, 20, 30).
		Build()

	// Output:
	// xychart
	//     x-axis Load 0 --> 100
	//     line [10, 20, 30]
}

// ExampleDiagram_YAxisRange sets the scale the values are drawn against.
// Without it mermaid picks one from the data, so two charts of different data
// cannot be compared by eye.
func ExampleDiagram_YAxisRange() {
	_ = xychart.NewDiagram(os.Stdout).
		XAxisLabels("Jan", "Feb").
		YAxisRange(0, 100).
		Bar(10, 20).
		Build()

	// Output:
	// xychart
	//     x-axis [Jan, Feb]
	//     y-axis 0 --> 100
	//     bar [10, 20]
}

// ExampleDiagram_YAxisRangeWithTitle sets the scale and names it.
func ExampleDiagram_YAxisRangeWithTitle() {
	_ = xychart.NewDiagram(os.Stdout).
		XAxisLabels("Jan", "Feb").
		YAxisRangeWithTitle("Revenue", 0, 100).
		Bar(10, 20).
		Build()

	// Output:
	// xychart
	//     x-axis [Jan, Feb]
	//     y-axis Revenue 0 --> 100
	//     bar [10, 20]
}

// ExampleDiagram_Bar draws the values as columns.
func ExampleDiagram_Bar() {
	_ = xychart.NewDiagram(os.Stdout).
		XAxisLabels("Jan", "Feb", "Mar").
		Bar(10, 20, 30).
		Build()

	// Output:
	// xychart
	//     x-axis [Jan, Feb, Mar]
	//     bar [10, 20, 30]
}

// ExampleDiagram_Line draws the values as a line, and a chart may carry both a
// line and bars.
func ExampleDiagram_Line() {
	_ = xychart.NewDiagram(os.Stdout).
		XAxisLabels("Jan", "Feb", "Mar").
		Bar(10, 20, 30).
		Line(15, 25, 35).
		Build()

	// Output:
	// xychart
	//     x-axis [Jan, Feb, Mar]
	//     bar [10, 20, 30]
	//     line [15, 25, 35]
}

// ExampleWithHorizontal turns the chart on its side, so the categories run down
// the page. It is the readable choice when the labels are long.
func ExampleWithHorizontal() {
	_ = xychart.NewDiagram(os.Stdout, xychart.WithHorizontal()).
		XAxisLabels("Jan", "Feb").
		Bar(10, 20).
		Build()

	// Output:
	// xychart horizontal
	//     x-axis [Jan, Feb]
	//     bar [10, 20]
}

// ExampleWithOrientation says which way the chart runs, for code that decides
// between the two rather than picking one outright.
func ExampleWithOrientation() {
	_ = xychart.NewDiagram(os.Stdout, xychart.WithOrientation(xychart.OrientationHorizontal)).
		XAxisLabels("Jan", "Feb").
		Bar(10, 20).
		Build()

	// Output:
	// xychart horizontal
	//     x-axis [Jan, Feb]
	//     bar [10, 20]
}

// ExampleOrientation shows the two directions a chart can run in.
func ExampleOrientation() {
	for _, orientation := range []xychart.Orientation{
		xychart.OrientationVertical,
		xychart.OrientationHorizontal,
	} {
		_ = xychart.NewDiagram(os.Stdout, xychart.WithOrientation(orientation)).
			XAxisLabels("Jan").
			Bar(10).
			Build()
		fmt.Println()
	}

	// Output:
	// xychart
	//     x-axis [Jan]
	//     bar [10]
	// xychart horizontal
	//     x-axis [Jan]
	//     bar [10]
}
