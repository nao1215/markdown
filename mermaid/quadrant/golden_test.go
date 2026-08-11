package quadrant_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/quadrant"
)

// TestGoldenQuadrantChart pins the rendered chart of every builder method of
// this package, including every point styling form.
func TestGoldenQuadrantChart(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := quadrant.NewChart(buf, quadrant.WithTitle("Reach And Engagement")).
		XAxis("Low Reach", "High Reach").
		YAxis("Low Engagement", "High Engagement").
		Quadrant1("We should expand").
		Quadrant2("Need to promote").
		Quadrant3("Re-evaluate").
		Quadrant4("May be improved").
		ClassDef("highlight", "color: #ff0000").
		ClassDefStyled("styled", quadrant.ClassStyle{
			Color:       "#00ff00",
			Radius:      10,
			StrokeColor: "#0000ff",
			StrokeWidth: "3px",
		}).
		LF().
		Point("Campaign A", 0.3, 0.6).
		PointWithStyle("Campaign B", 0.45, 0.23, "radius: 10").
		PointStyled("Campaign C", 0.57, 0.69, quadrant.PointStyle{
			Color:       "#ff00ff",
			Radius:      5,
			StrokeColor: "#000000",
			StrokeWidth: "2px",
		}).
		PointWithClass("Campaign D", 0.78, 0.34, "highlight").
		PointWithClassAndStyle("Campaign E", 0.4, 0.34, "styled", "radius: 20").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("quadrant.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenQuadrantChartSingleSidedAxes pins the axis forms that take only one
// label, which is the variadic half of XAxis and YAxis.
func TestGoldenQuadrantChartSingleSidedAxes(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := quadrant.NewChart(buf).
		XAxis("Low Reach").
		YAxis("Low Engagement").
		Point("Campaign A", 0.3, 0.6).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("quadrant_single_sided_axes.md", buf.String()); err != nil {
		t.Error(err)
	}
}
