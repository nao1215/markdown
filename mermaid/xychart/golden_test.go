package xychart_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/xychart"
)

// TestGoldenXYChart pins the rendered chart of every axis form and every series
// kind this package can build.
func TestGoldenXYChart(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := xychart.NewDiagram(buf, xychart.WithTitle("Sales Revenue")).
		XAxisLabelsWithTitle("Month", "jan", "feb", "mar").
		YAxisRangeWithTitle("Revenue (in $)", 4000, 11000).
		Bar(5000, 6000, 7500).
		Line(5000, 6000, 7500).
		LF().
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("xychart.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenXYChartAxisWithoutTitle pins the axis forms that take no title, and
// the numeric x axis range.
func TestGoldenXYChartAxisWithoutTitle(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := xychart.NewDiagram(buf).
		XAxisRange(1, 12).
		YAxisRange(0, 100).
		Line(10, 20, 30).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("xychart_ranges.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenXYChartTitledRanges pins the axis forms that carry a title and a
// numeric range, which is the one combination the tests above leave out.
func TestGoldenXYChartTitledRanges(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := xychart.NewDiagram(buf).
		XAxisRangeWithTitle("Month", 1, 12).
		YAxisRangeWithTitle("Revenue (in $)", 0, 11000).
		Bar(5000, 6000).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("xychart_titled_ranges.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenXYChartOrientations pins the header each orientation produces.
func TestGoldenXYChartOrientations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		golden string
		option xychart.Option
	}{
		{name: "vertical", golden: "xychart_vertical.md", option: xychart.WithOrientation(xychart.OrientationVertical)},
		{name: "horizontal", golden: "xychart_horizontal.md", option: xychart.WithOrientation(xychart.OrientationHorizontal)},
		{name: "horizontal shorthand", golden: "xychart_horizontal_shorthand.md", option: xychart.WithHorizontal()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			err := xychart.NewDiagram(buf, tt.option).
				XAxisLabels("a", "b").
				Bar(1, 2).
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
