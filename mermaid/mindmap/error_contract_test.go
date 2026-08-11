package mindmap_test

import (
	"io"
	"testing"

	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/mermaid/mindmap"
)

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return mindmap.NewDiagram(w).Root("Product").Child("Market")
	})
}

// TestRecordedErrorContract asserts that a child added before the root surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return mindmap.NewDiagram(w).Child("a child with no root")
	})
}
