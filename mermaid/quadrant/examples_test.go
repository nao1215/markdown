//go:build linux || darwin

package quadrant_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/quadrant"
)

// ExampleChart skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleChart() {
	chart := quadrant.NewChart(
		os.Stdout,
		quadrant.WithTitle("Priority Matrix"),
	).
		XAxis("Low Effort", "High Effort").
		YAxis("Low Impact", "High Impact").
		Quadrant1("Quick Wins").
		Quadrant2("Major Projects").
		Quadrant3("Fill Ins").
		Quadrant4("Time Wasters").
		Point("Task A", 0.8, 0.9).
		Point("Task B", 0.2, 0.7).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Quadrant Chart").
		CodeBlocks(md.SyntaxHighlightMermaid, chart).
		Build()

	// Output:
	// ## Quadrant Chart
	// ```mermaid
	// quadrantChart
	//     title Priority Matrix
	//     x-axis Low Effort --> High Effort
	//     y-axis Low Impact --> High Impact
	//     quadrant-1 Quick Wins
	//     quadrant-2 Major Projects
	//     quadrant-3 Fill Ins
	//     quadrant-4 Time Wasters
	//     Task A: [0.80, 0.90]
	//     Task B: [0.20, 0.70]
	// ```
}

// ExampleChart_withStyling demonstrates the use of point styling and class definitions.
func ExampleChart_withStyling() {
	chart := quadrant.NewChart(
		os.Stdout,
		quadrant.WithTitle("Reach and engagement of campaigns"),
	).
		XAxis("Low Reach", "High Reach").
		YAxis("Low Engagement", "High Engagement").
		Quadrant1("We should expand").
		Quadrant2("Need to promote").
		Quadrant3("Re-evaluate").
		Quadrant4("May be improved").
		PointWithStyle("Campaign A", 0.9, 0.0, "radius: 12").
		PointWithClass("Campaign B", 0.8, 0.1, "class1").
		PointStyled("Campaign C", 0.7, 0.2, quadrant.PointStyle{
			Radius:      25,
			Color:       "#00ff33",
			StrokeColor: "#10f0f0",
		}).
		PointWithClass("Campaign D", 0.5, 0.4, "class2").
		ClassDef("class1", "color: #109060").
		ClassDefStyled("class2", quadrant.ClassStyle{
			Color:       "#908342",
			Radius:      10,
			StrokeColor: "#310085",
			StrokeWidth: "10px",
		}).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Quadrant Chart with Styling").
		CodeBlocks(md.SyntaxHighlightMermaid, chart).
		Build()

	// Output:
	// ## Quadrant Chart with Styling
	// ```mermaid
	// quadrantChart
	//     title Reach and engagement of campaigns
	//     x-axis Low Reach --> High Reach
	//     y-axis Low Engagement --> High Engagement
	//     quadrant-1 We should expand
	//     quadrant-2 Need to promote
	//     quadrant-3 Re-evaluate
	//     quadrant-4 May be improved
	//     Campaign A: [0.90, 0.00] radius: 12
	//     Campaign B:::class1: [0.80, 0.10]
	//     Campaign C: [0.70, 0.20] color: #00ff33, radius: 25, stroke-color: #10f0f0
	//     Campaign D:::class2: [0.50, 0.40]
	//     classDef class1 color: #109060
	//     classDef class2 color: #908342, radius: 10, stroke-color: #310085, stroke-width: 10px
	// ```
}

// ExampleNewChart shows the shape every quadrant chart has: a writer, the axes,
// the quadrant labels, the points, and Build.
func ExampleNewChart() {
	_ = quadrant.NewChart(os.Stdout, quadrant.WithTitle("Reach and Engagement")).
		XAxis("Low Reach", "High Reach").
		YAxis("Low Engagement", "High Engagement").
		Quadrant1("We should expand").
		Quadrant2("Need to promote").
		Quadrant3("Re-evaluate").
		Quadrant4("May be improved").
		Point("Campaign A", 0.3, 0.6).
		Build()

	// Output:
	// quadrantChart
	//     title Reach and Engagement
	//     x-axis Low Reach --> High Reach
	//     y-axis Low Engagement --> High Engagement
	//     quadrant-1 We should expand
	//     quadrant-2 Need to promote
	//     quadrant-3 Re-evaluate
	//     quadrant-4 May be improved
	//     Campaign A: [0.30, 0.60]
}

// ExampleChart_XAxis names the two ends of the horizontal axis. Passing one
// label names the axis itself instead.
func ExampleChart_XAxis() {
	_ = quadrant.NewChart(os.Stdout).
		XAxis("Low Reach", "High Reach").
		Build()

	// Output:
	// quadrantChart
	//     x-axis Low Reach --> High Reach
}

// ExampleChart_YAxis names the two ends of the vertical axis.
func ExampleChart_YAxis() {
	_ = quadrant.NewChart(os.Stdout).
		YAxis("Low Engagement", "High Engagement").
		Build()

	// Output:
	// quadrantChart
	//     y-axis Low Engagement --> High Engagement
}

// ExampleChart_Quadrant1 names the top right quadrant.
func ExampleChart_Quadrant1() {
	_ = quadrant.NewChart(os.Stdout).Quadrant1("We should expand").Build()

	// Output:
	// quadrantChart
	//     quadrant-1 We should expand
}

// ExampleChart_Quadrant2 names the top left quadrant.
func ExampleChart_Quadrant2() {
	_ = quadrant.NewChart(os.Stdout).Quadrant2("Need to promote").Build()

	// Output:
	// quadrantChart
	//     quadrant-2 Need to promote
}

// ExampleChart_Quadrant3 names the bottom left quadrant.
func ExampleChart_Quadrant3() {
	_ = quadrant.NewChart(os.Stdout).Quadrant3("Re-evaluate").Build()

	// Output:
	// quadrantChart
	//     quadrant-3 Re-evaluate
}

// ExampleChart_Quadrant4 names the bottom right quadrant.
func ExampleChart_Quadrant4() {
	_ = quadrant.NewChart(os.Stdout).Quadrant4("May be improved").Build()

	// Output:
	// quadrantChart
	//     quadrant-4 May be improved
}

// ExampleChart_Point places one thing on the chart. Both coordinates run from
// 0.0 to 1.0, and they are written with two digits after the point.
func ExampleChart_Point() {
	_ = quadrant.NewChart(os.Stdout).
		Point("Campaign A", 0.3, 0.6).
		Point("Campaign B", 0.45, 0.23).
		Build()

	// Output:
	// quadrantChart
	//     Campaign A: [0.30, 0.60]
	//     Campaign B: [0.45, 0.23]
}

// ExampleChart_PointWithStyle places a point and styles it with a string
// mermaid reads directly.
func ExampleChart_PointWithStyle() {
	_ = quadrant.NewChart(os.Stdout).
		PointWithStyle("Campaign A", 0.3, 0.6, "radius: 10, color: #ff0000").
		Build()

	// Output:
	// quadrantChart
	//     Campaign A: [0.30, 0.60] radius: 10, color: #ff0000
}

// ExampleChart_PointStyled places a point and styles it with a struct, so the
// caller does not have to build the string mermaid wants.
func ExampleChart_PointStyled() {
	_ = quadrant.NewChart(os.Stdout).
		PointStyled("Campaign A", 0.3, 0.6, quadrant.PointStyle{
			Color:       "#ff0000",
			Radius:      10,
			StrokeColor: "#000000",
			StrokeWidth: "2px",
		}).
		Build()

	// Output:
	// quadrantChart
	//     Campaign A: [0.30, 0.60] color: #ff0000, radius: 10, stroke-color: #000000, stroke-width: 2px
}

// ExampleChart_PointWithClass places a point styled by a class, which is how
// several points share one style.
func ExampleChart_PointWithClass() {
	_ = quadrant.NewChart(os.Stdout).
		ClassDef("winner", "color: #00ff00").
		PointWithClass("Campaign A", 0.3, 0.6, "winner").
		Build()

	// Output:
	// quadrantChart
	//     classDef winner color: #00ff00
	//     Campaign A:::winner: [0.30, 0.60]
}

// ExampleChart_PointWithClassAndStyle places a point styled by a class and then
// adjusted, for the one point that should stand out among its peers.
func ExampleChart_PointWithClassAndStyle() {
	_ = quadrant.NewChart(os.Stdout).
		ClassDef("winner", "color: #00ff00").
		PointWithClassAndStyle("Campaign A", 0.3, 0.6, "winner", "radius: 12").
		Build()

	// Output:
	// quadrantChart
	//     classDef winner color: #00ff00
	//     Campaign A:::winner: [0.30, 0.60] radius: 12
}

// ExampleChart_ClassDef names a style so several points can share it.
func ExampleChart_ClassDef() {
	_ = quadrant.NewChart(os.Stdout).
		ClassDef("winner", "color: #00ff00, radius: 10").
		PointWithClass("Campaign A", 0.3, 0.6, "winner").
		Build()

	// Output:
	// quadrantChart
	//     classDef winner color: #00ff00, radius: 10
	//     Campaign A:::winner: [0.30, 0.60]
}

// ExampleChart_ClassDefStyled names a style described by a struct rather than a
// string.
func ExampleChart_ClassDefStyled() {
	_ = quadrant.NewChart(os.Stdout).
		ClassDefStyled("winner", quadrant.ClassStyle{Color: "#00ff00", Radius: 10}).
		PointWithClass("Campaign A", 0.3, 0.6, "winner").
		Build()

	// Output:
	// quadrantChart
	//     classDef winner color: #00ff00, radius: 10
	//     Campaign A:::winner: [0.30, 0.60]
}

// ExampleChart_String returns the chart without needing a writer, which is how
// it is handed to a markdown code block.
func ExampleChart_String() {
	chart := quadrant.NewChart(io.Discard).Point("Campaign A", 0.3, 0.6).String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, chart).
		Build()

	// Output:
	// ```mermaid
	// quadrantChart
	//     Campaign A: [0.30, 0.60]
	// ```
}

// ExampleChart_Build writes the chart and reports the error the chain recorded.
func ExampleChart_Build() {
	err := quadrant.NewChart(nil).Point("Campaign A", 0.3, 0.6).Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleChart_Error reports the same error Build does, for code that wants to
// look before writing anything.
func ExampleChart_Error() {
	ch := quadrant.NewChart(io.Discard).Point("Campaign A", 0.3, 0.6)
	fmt.Println("error:", ch.Error())

	// Output:
	// error: <nil>
}

// ExampleChart_LF adds a blank line to the chart body.
func ExampleChart_LF() {
	_ = quadrant.NewChart(os.Stdout).
		Point("Campaign A", 0.3, 0.6).
		LF().
		Point("Campaign B", 0.45, 0.23).
		Build()

	// Output:
	// quadrantChart
	//     Campaign A: [0.30, 0.60]
	//
	//     Campaign B: [0.45, 0.23]
}

// ExamplePointStyle shows how one point is styled. Every field is optional, and
// the zero value of one leaves mermaid's default alone.
func ExamplePointStyle() {
	style := quadrant.PointStyle{Color: "#ff0000", Radius: 10}

	_ = quadrant.NewChart(os.Stdout).
		PointStyled("Campaign A", 0.3, 0.6, style).
		Build()

	// Output:
	// quadrantChart
	//     Campaign A: [0.30, 0.60] color: #ff0000, radius: 10
}

// ExamplePointStyle_String renders the style as the string mermaid reads, which
// is what PointStyled writes for you.
func ExamplePointStyle_String() {
	style := quadrant.PointStyle{Color: "#ff0000", Radius: 10, StrokeWidth: "2px"}
	fmt.Println(style.String())

	// Output:
	// color: #ff0000, radius: 10, stroke-width: 2px
}

// ExampleClassStyle shows how a named style is described.
func ExampleClassStyle() {
	style := quadrant.ClassStyle{Color: "#00ff00", Radius: 10}

	_ = quadrant.NewChart(os.Stdout).
		ClassDefStyled("winner", style).
		PointWithClass("Campaign A", 0.3, 0.6, "winner").
		Build()

	// Output:
	// quadrantChart
	//     classDef winner color: #00ff00, radius: 10
	//     Campaign A:::winner: [0.30, 0.60]
}

// ExampleClassStyle_String renders the style as the string mermaid reads, which
// is what ClassDefStyled writes for you.
func ExampleClassStyle_String() {
	style := quadrant.ClassStyle{Color: "#00ff00", Radius: 10}
	fmt.Println(style.String())

	// Output:
	// color: #00ff00, radius: 10
}

// ExampleWithTitle sets the title the chart is drawn with.
func ExampleWithTitle() {
	_ = quadrant.NewChart(os.Stdout, quadrant.WithTitle("Reach and Engagement")).
		Point("Campaign A", 0.3, 0.6).
		Build()

	// Output:
	// quadrantChart
	//     title Reach and Engagement
	//     Campaign A: [0.30, 0.60]
}

// ExampleOption shows what an Option is: a function that changes how the chart
// is written, passed to NewChart.
func ExampleOption() {
	options := []quadrant.Option{quadrant.WithTitle("Reach and Engagement")}

	_ = quadrant.NewChart(os.Stdout, options...).
		Point("Campaign A", 0.3, 0.6).
		Build()

	// Output:
	// quadrantChart
	//     title Reach and Engagement
	//     Campaign A: [0.30, 0.60]
}
