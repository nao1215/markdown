package radar_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/radar"
)

// TestGoldenRadar pins the rendered chart of every builder method and every
// option of this package, including both escaping cases.
func TestGoldenRadar(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := radar.NewDiagram(buf, radar.WithTitle("Grades: term one")).
		Axis("Math: advanced", `Quote "lit"`, `Back\slash`).
		Axis("History #1", "Art (fine)").
		LF().
		Curve("Alice; team", 85, 90, 80, 70, 75).
		Curve("Bob", 70.5, 75, 85, 80, 90).
		Max(100).
		Min(0).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("radar.md", buf.String()); err != nil {
		t.Error(err)
	}
}
