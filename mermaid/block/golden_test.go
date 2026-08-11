package block_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/block"
)

// TestGoldenBlock pins the rendered diagram of every builder method, every
// token constructor, and every token option of this package.
func TestGoldenBlock(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := block.NewDiagram(buf, block.WithTitle("Every Token")).
		Columns(3).
		Row(
			block.Node("plain"),
			block.Node("labelled", block.WithNodeLabel("With a label")),
			block.Node("wide", block.WithNodeSpan(2)),
		).
		Row(
			block.Space(),
			block.Space(2),
			block.Literal("raw:3"),
		).
		Row(
			block.Arrow("explicit", block.DirectionRight),
			block.Arrow("labelled", block.DirectionRight, block.WithArrowLabel("to the right")),
			block.Arrow("secondary", block.DirectionRight, block.WithArrowSecondaryDirection(block.DirectionDown)),
		).
		Row(
			block.ArrowRight("right"),
			block.ArrowLeft("left"),
			block.ArrowUp("up"),
		).
		Row(
			block.ArrowDown("down"),
			block.ArrowX("x"),
			block.ArrowY("y"),
		).
		LF().
		Block(func(d *block.Diagram) {
			d.Columns(2).
				Row(block.Node("inner1"), block.Node("inner2")).
				Statement("style inner1 fill:#f9f")
		}, block.WithBlockID("grouped"), block.WithBlockSpan(2)).
		Statement("style plain fill:#eee").
		Link("plain", "labelled").
		LinkWithLabel("labelled", "connects to", "wide").
		Style("plain", "fill:#fff").
		ClassDef("highlight", "fill:#ff0").
		Class("plain", "highlight").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("block.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenBlockShapes pins one node per shape, because the shape decides the
// brackets that surround the label and there is one pair per shape.
func TestGoldenBlockShapes(t *testing.T) {
	t.Parallel()

	shapes := []block.Shape{
		block.ShapeRectangle,
		block.ShapeRound,
		block.ShapeStadium,
		block.ShapeSubroutine,
		block.ShapeCylinder,
		block.ShapeCircle,
		block.ShapeAsymmetric,
		block.ShapeRhombus,
		block.ShapeHexagon,
		block.ShapeParallelogram,
		block.ShapeParallelogramAlt,
		block.ShapeTrapezoid,
		block.ShapeTrapezoidAlt,
		block.ShapeDoubleCircle,
	}

	buf := &bytes.Buffer{}
	diagram := block.NewDiagram(buf)
	for _, shape := range shapes {
		diagram = diagram.Row(block.Node(
			string(shape),
			block.WithNodeLabel(string(shape)),
			block.WithNodeShape(shape),
		))
	}
	if err := diagram.Build(); err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("block_shapes.md", buf.String()); err != nil {
		t.Error(err)
	}
}
