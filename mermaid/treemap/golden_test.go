package treemap_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/treemap"
)

// TestGoldenTreemap pins the rendered diagram of every builder method and every
// option of this package, including the quote doubling and three levels of
// nesting.
func TestGoldenTreemap(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := treemap.NewDiagram(buf, treemap.WithTitle("Budget: 2024")).
		Section(`Ops "core"`).
		Leaf("Salaries: base", 1200).
		Section("Cloud #1").
		Leaf("Compute", 400.5).
		Leaf(`Back\slash`, 0).
		Parent().
		Leaf("Travel; trips", 300).
		Parent().
		LF().
		Section("Marketing").
		Leaf("Ads (paid)", 800).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("treemap.md", buf.String()); err != nil {
		t.Error(err)
	}
}
