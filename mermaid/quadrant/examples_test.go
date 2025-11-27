//go:build linux || darwin

package quadrant_test

import (
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
