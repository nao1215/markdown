//go:build linux || darwin

package piechart_test

import (
	"fmt"
	"io"
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

// ExampleNewPieChart shows the shape every pie chart has: a writer, a chain of
// calls, and Build.
func ExampleNewPieChart() {
	_ = piechart.NewPieChart(os.Stdout).
		LabelAndIntValue("Go", 60).
		LabelAndIntValue("Rust", 40).
		Build()

	// Output:
	// %%{init: {"pie": {"textPosition": 0.75}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
	// pie
	//     "Go" : 60
	//     "Rust" : 40
}

// ExamplePieChart_LabelAndIntValue adds a slice with a whole number.
func ExamplePieChart_LabelAndIntValue() {
	_ = piechart.NewPieChart(os.Stdout).
		LabelAndIntValue("Go", 60).
		LabelAndIntValue("Rust", 40).
		Build()

	// Output:
	// %%{init: {"pie": {"textPosition": 0.75}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
	// pie
	//     "Go" : 60
	//     "Rust" : 40
}

// ExamplePieChart_LabelAndFloatValue adds a slice with a fractional value. It
// is written with six digits after the point, whatever the value.
func ExamplePieChart_LabelAndFloatValue() {
	_ = piechart.NewPieChart(os.Stdout).
		LabelAndFloatValue("Go", 60.5).
		LabelAndFloatValue("Rust", 39.5).
		Build()

	// Output:
	// %%{init: {"pie": {"textPosition": 0.75}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
	// pie
	//     "Go" : 60.500000
	//     "Rust" : 39.500000
}

// ExamplePieChart_String returns the chart without needing a writer, which is
// how it is handed to a markdown code block.
func ExamplePieChart_String() {
	chart := piechart.NewPieChart(io.Discard).
		LabelAndIntValue("Go", 60).
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, chart).
		Build()

	// Output:
	// ```mermaid
	// %%{init: {"pie": {"textPosition": 0.75}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
	// pie
	//     "Go" : 60
	// ```
}

// ExamplePieChart_Build writes the chart and reports the error the chain
// recorded.
func ExamplePieChart_Build() {
	err := piechart.NewPieChart(nil).LabelAndIntValue("Go", 60).Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleWithTitle sets the title the chart is drawn with.
func ExampleWithTitle() {
	_ = piechart.NewPieChart(os.Stdout, piechart.WithTitle("Languages")).
		LabelAndIntValue("Go", 60).
		Build()

	// Output:
	// %%{init: {"pie": {"textPosition": 0.75}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
	// pie
	//     title Languages
	//     "Go" : 60
}

// ExampleWithShowData draws each slice's value beside its label, which is the
// difference between a chart a reader can read numbers off and one they cannot.
func ExampleWithShowData() {
	_ = piechart.NewPieChart(os.Stdout, piechart.WithShowData(true)).
		LabelAndIntValue("Go", 60).
		LabelAndIntValue("Rust", 40).
		Build()

	// Output:
	// %%{init: {"pie": {"textPosition": 0.75}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
	// pie showData
	//     "Go" : 60
	//     "Rust" : 40
}

// ExampleWithTextPosition moves the labels between the center of the chart and
// its outer edge. It takes 0.0 to 1.0, and anything outside that is ignored in
// favor of the default.
func ExampleWithTextPosition() {
	_ = piechart.NewPieChart(os.Stdout, piechart.WithTextPosition(0.4)).
		LabelAndIntValue("Go", 60).
		Build()

	// Output:
	// %%{init: {"pie": {"textPosition": 0.40}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
	// pie
	//     "Go" : 60
}

// ExampleOption shows what an Option is: a function that changes how the chart
// is written, passed to NewPieChart.
func ExampleOption() {
	options := []piechart.Option{piechart.WithTitle("Languages"), piechart.WithShowData(true)}

	_ = piechart.NewPieChart(os.Stdout, options...).
		LabelAndIntValue("Go", 60).
		Build()

	// Output:
	// %%{init: {"pie": {"textPosition": 0.75}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
	// pie showData
	//     title Languages
	//     "Go" : 60
}

// ExamplePieChart_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExamplePieChart_Error() {
	p := piechart.NewPieChart(nil).LabelAndIntValue("Go", 60)
	fmt.Println("before Build:", p.Error())
	_ = p.Build()
	fmt.Println("after Build:", p.Error())

	// Output:
	// before Build: <nil>
	// after Build: output writer must not be nil
}
