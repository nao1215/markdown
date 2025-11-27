//go:build linux || darwin

package piechart_test

import (
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/piechart"
)

// ExamplePieChart skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExamplePieChart() {
	chart := piechart.NewPieChart(
		os.Stdout,
		piechart.WithTitle("mermaid pie chart builder"),
		piechart.WithShowData(true),
	).
		LabelAndIntValue("A", 10).
		LabelAndFloatValue("B", 20.1).
		LabelAndIntValue("C", 30).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Pie Chart Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, chart).
		Build()

	// Output:
	// ## Pie Chart Diagram
	// ```mermaid
	// %%{init: {"pie": {"textPosition": 0.75}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
	// pie showData
	//     title mermaid pie chart builder
	//     "A" : 10
	//     "B" : 20.100000
	//     "C" : 30
	// ```
}
