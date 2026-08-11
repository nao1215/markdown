package sankey_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/sankey"
)

// TestGoldenSankey pins the rendered diagram of every builder method and every
// option of this package, including both CSV quoting cases.
func TestGoldenSankey(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := sankey.NewDiagram(buf, sankey.WithTitle("Energy flow")).
		Link("Agricultural 'waste'", "Bio-conversion", 124.729).
		Link("Bio-conversion", "Liquid", 0.597).
		Link("Bio-conversion", "Losses, and more", 26.862).
		Link(`He said "hi"`, "Gas", 81.144).
		LF().
		Link("Electricity", "Homes", 150).
		Link("Electricity", "Industry", 0).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("sankey.md", buf.String()); err != nil {
		t.Error(err)
	}
}
