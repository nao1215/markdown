package mindmap_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/mindmap"
)

// TestGoldenMindmap pins the rendered diagram of every builder method of this
// package, including the explicit depth form of Node.
func TestGoldenMindmap(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := mindmap.NewDiagram(buf, mindmap.WithTitle("Product Strategy")).
		Root("Product Strategy").
		Child("Market").
		Child("SMB").
		Sibling("Enterprise").
		Parent().
		Sibling("Execution").
		Child("Q1").
		Sibling("Q2").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("mindmap.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenMindmapExplicitDepth pins the explicit depth form of the builder,
// where the caller states the level of every node instead of walking the tree
// with Child, Sibling, and Parent.
func TestGoldenMindmapExplicitDepth(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := mindmap.NewDiagram(buf).
		Root("Root").
		Node(1, "Level one").
		Node(2, "Level two").
		LF().
		Node(2, "Another level two").
		Node(1, "Back to level one").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("mindmap_explicit_depth.md", buf.String()); err != nil {
		t.Error(err)
	}
}
