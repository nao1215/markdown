package gitgraph_test

import (
	"io"
	"testing"

	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/mermaid/gitgraph"
)

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return gitgraph.NewDiagram(w).Commit()
	})
}

// TestRecordedErrorContract asserts that a branch order the syntax cannot express surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return gitgraph.NewDiagram(w).Commit().Branch("develop", gitgraph.WithBranchOrder(-1))
	})
}
