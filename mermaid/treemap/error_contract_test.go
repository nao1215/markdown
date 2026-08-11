package treemap_test

import (
	"io"
	"testing"

	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/mermaid/treemap"
)

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return treemap.NewDiagram(w).Section("Ops").Leaf("Salaries", 1200)
	})
}

// TestRecordedErrorContract asserts that going up from the top level surfaces
// from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return treemap.NewDiagram(w).Parent()
	})
}
