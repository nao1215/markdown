package flowchart_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/flowchart"
)

// TestGoldenFlowchart pins the rendered diagram of every node shape and every
// link style this package can build.
func TestGoldenFlowchart(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := flowchart.NewFlowchart(
		buf,
		flowchart.WithTitle("Every Shape And Link"),
		flowchart.WithOrientalTopToBottom(),
	).
		Node("plain").
		NodeWithText("text", "With text").
		NodeWithMarkdown("markdown", "**bold** text").
		NodeWithNewLines("newlines", "first\nsecond").
		RoundEdgesNode("round", "Round edges").
		StadiumNode("stadium", "Stadium").
		SubroutineNode("subroutine", "Subroutine").
		CylindricalNode("cylindrical", "Cylindrical").
		DatabaseNode("database", "Database").
		CircleNode("circle", "Circle").
		AsymmetricNode("asymmetric", "Asymmetric").
		RhombusNode("rhombus", "Rhombus").
		HexagonNode("hexagon", "Hexagon").
		ParallelogramNode("parallelogram", "Parallelogram").
		ParallelogramAltNode("parallelogramAlt", "Parallelogram alt").
		TrapezoidNode("trapezoid", "Trapezoid").
		TrapezoidAltNode("trapezoidAlt", "Trapezoid alt").
		DoubleCircleNode("doubleCircle", "Double circle").
		LinkWithArrowHead("plain", "text").
		LinkWithArrowHeadAndText("text", "markdown", "with text").
		OpenLink("markdown", "newlines").
		OpenLinkWithText("newlines", "round", "open with text").
		DottedLink("round", "stadium").
		DottedLinkWithText("stadium", "subroutine", "dotted with text").
		ThickLink("subroutine", "cylindrical").
		ThickLinkWithText("cylindrical", "database", "thick with text").
		InvisibleLink("database", "circle").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("flowchart.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenFlowchartOrientations pins the header each orientation option
// produces. The options are mutually exclusive, so each needs its own diagram.
func TestGoldenFlowchartOrientations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		golden string
		option flowchart.Option
	}{
		{name: "top to bottom", golden: "orientation_tb.md", option: flowchart.WithOrientalTopToBottom()},
		{name: "top down", golden: "orientation_td.md", option: flowchart.WithOrientalTopDown()},
		{name: "bottom to top", golden: "orientation_bt.md", option: flowchart.WithOrientalBottomToTop()},
		{name: "right to left", golden: "orientation_rl.md", option: flowchart.WithOrientalRightToLeft()},
		{name: "left to right", golden: "orientation_lr.md", option: flowchart.WithOrientalLeftToRight()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			err := flowchart.NewFlowchart(buf, tt.option).
				NodeWithText("A", "Start").
				NodeWithText("B", "End").
				LinkWithArrowHead("A", "B").
				Build()
			if err != nil {
				t.Fatalf("Build() = %v, want nil", err)
			}

			if err := golden.Assert(tt.golden, buf.String()); err != nil {
				t.Error(err)
			}
		})
	}
}
