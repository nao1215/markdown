package piechart_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/piechart"
)

// TestGoldenPieChart pins the rendered chart of every builder method and every
// option of this package.
func TestGoldenPieChart(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := piechart.NewPieChart(
		buf,
		piechart.WithTitle("Language Share"),
		piechart.WithShowData(true),
		piechart.WithTextPosition(0.75),
	).
		LabelAndIntValue("Go", 120).
		LabelAndFloatValue("Rust", 42.5).
		LabelAndFloatValue("Zig", 0.5).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("piechart.md", buf.String()); err != nil {
		t.Error(err)
	}
}
