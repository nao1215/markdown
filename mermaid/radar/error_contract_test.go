package radar_test

import (
	"io"
	"testing"

	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/mermaid/radar"
)

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return radar.NewDiagram(w).Axis("Math", "Science").Curve("Alice", 85, 90)
	})
}

// TestRecordedErrorContract asserts that an axis with no labels surfaces from
// Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return radar.NewDiagram(w).Axis()
	})
}
