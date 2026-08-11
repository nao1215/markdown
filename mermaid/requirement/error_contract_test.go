package requirement_test

import (
	"io"
	"testing"

	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/mermaid/requirement"
)

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return requirement.NewDiagram(w).Requirement(
			"a requirement",
			requirement.WithID("1"),
			requirement.WithText("the system shall do the thing"),
			requirement.WithRisk(requirement.RiskLow),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest),
		)
	})
}

// TestRecordedErrorContract asserts that a requirement missing the fields the syntax needs surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return requirement.NewDiagram(w).Requirement("a requirement with no id")
	})
}
