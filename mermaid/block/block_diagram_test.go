package block

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type failingWriter struct{}

func (f failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestNewDiagram(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    []Option
		want    string
		wantErr bool
	}{
		{
			name: "new diagram without options",
			opts: nil,
			want: "block",
		},
		{
			name: "new diagram with title",
			opts: []Option{WithTitle("Checkout Architecture")},
			want: `block
    title Checkout Architecture`,
		},
		{
			name:    "new diagram with title including newline",
			opts:    []Option{WithTitle("Checkout\nArchitecture")},
			want:    "block",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diagram := NewDiagram(io.Discard, tt.opts...)
			if tt.wantErr && diagram.Error() == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && diagram.Error() != nil {
				t.Fatalf("unexpected error: %v", diagram.Error())
			}

			got := strings.ReplaceAll(diagram.String(), "\r\n", "\n")
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDiagram_Build(t *testing.T) {
	t.Parallel()

	b := new(bytes.Buffer)

	d := NewDiagram(b, WithTitle("Checkout Architecture")).
		Columns(3).
		Row(
			Node("Frontend"),
			ArrowRight("toBackend", WithArrowLabel("calls")),
			Node("Backend"),
		).
		Row(
			Space(2),
			ArrowDown("toDB"),
		).
		Row(
			Node("Database", WithNodeLabel("Primary DB"), WithNodeShape(ShapeCylinder)),
			Space(),
			Node("Cache", WithNodeLabel("Cache"), WithNodeShape(ShapeRound)),
		).
		Link("Backend", "Database").
		LinkWithLabel("Backend", "reads from", "Cache").
		Style("Backend", "fill:#9cf,stroke:#333").
		ClassDef("service", "fill:#9cf,stroke:#333").
		Class("Frontend,Backend,Database", "service").
		Block(func(blk *Diagram) {
			blk.Columns(1).Row(Node("PaymentService"))
		}, WithBlockID("payments"), WithBlockSpan(2))

	if err := d.Build(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `block
    title Checkout Architecture
    columns 3
    Frontend toBackend<["calls"]>(right) Backend
    space:2 toDB<["&nbsp;"]>(down)
    Database[("Primary DB")] space Cache("Cache")
    Backend --> Database
    Backend -- "reads from" --> Cache
    style Backend fill:#9cf,stroke:#333
    classDef service fill:#9cf,stroke:#333
    class Frontend,Backend,Database service
    block:payments:2
        columns 1
        PaymentService
    end`

	got := strings.ReplaceAll(b.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestNormalizeQuoted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "paired quotes",
			input: `"hello"`,
			want:  "hello",
		},
		{
			name:  "leading quote only",
			input: `"hello`,
			want:  `"hello`,
		},
		{
			name:  "trailing quote only",
			input: `hello"`,
			want:  `hello"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := normalizeQuoted(tt.input); got != tt.want {
				t.Errorf("normalizeQuoted(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestQuoteEscapesSpecialChars(t *testing.T) {
	t.Parallel()

	got := quote("a\\b\rc\nd\te\"f")
	want := `"a&#92;b&#92;rc&#92;nd&#92;te&quot;f"`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func() *Diagram
		want string
	}{
		{
			name: "invalid columns",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Columns(0)
			},
			want: "block",
		},
		{
			name: "empty row",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Row()
			},
			want: "block",
		},
		{
			name: "invalid literal token",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Row(Literal(""))
			},
			want: "block",
		},
		{
			name: "invalid node id with whitespace",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Row(Node("bad id"))
			},
			want: "block",
		},
		{
			name: "invalid arrow direction",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Row(
					Arrow("arrow", Direction("side"), WithArrowLabel("label")),
				)
			},
			want: "block",
		},
		{
			name: "invalid arrow secondary direction",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Row(
					Arrow(
						"arrow",
						DirectionRight,
						WithArrowLabel("label"),
						WithArrowSecondaryDirection(Direction("north")),
					),
				)
			},
			want: "block",
		},
		{
			name: "invalid block span",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Block(func(*Diagram) {}, WithBlockSpan(0))
			},
			want: `block`,
		},
		{
			name: "block callback error does not append end",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Block(func(blk *Diagram) {
						blk.Columns(0)
					})
			},
			want: "block\n    block",
		},
		{
			name: "nil block builder",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Block(nil)
			},
			want: "block",
		},
		{
			name: "empty link label",
			run: func() *Diagram {
				return NewDiagram(io.Discard).LinkWithLabel("A", "", "B")
			},
			want: "block",
		},
		{
			name: "empty style names",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Style("", "fill:#ccc")
			},
			want: "block",
		},
		{
			name: "empty class style",
			run: func() *Diagram {
				return NewDiagram(io.Discard).ClassDef("service", "")
			},
			want: "block",
		},
		{
			name: "newline in title",
			run: func() *Diagram {
				return NewDiagram(io.Discard, WithTitle("Checkout\nArchitecture"))
			},
			want: "block",
		},
		{
			name: "lf short-circuit after error",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Columns(0).LF()
			},
			want: "block",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := tt.run()
			if d.Error() == nil {
				t.Fatal("expected error, got nil")
			}

			got := strings.ReplaceAll(d.String(), "\r\n", "\n")
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDiagram_BuildStoresError(t *testing.T) {
	t.Parallel()

	d := NewDiagram(failingWriter{})
	err := d.Build()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if d.Error() == nil {
		t.Fatal("expected persisted error, got nil")
	}
	if !errors.Is(d.Error(), err) {
		t.Fatalf("expected Error() to wrap returned error, got %v", d.Error())
	}
}

func TestDiagram_BuildNilWriter(t *testing.T) {
	t.Parallel()

	d := NewDiagram(nil)
	err := d.Build()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "output writer must not be nil" {
		t.Fatalf("unexpected error: %v", err)
	}

	if d.Error() == nil {
		t.Fatal("expected persisted error, got nil")
	}
	if !errors.Is(d.Error(), err) {
		t.Fatalf("expected Error() to wrap returned error, got %v", d.Error())
	}
}

func TestNodeShapeFormatting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		shape Shape
		want  string
	}{
		{name: "rectangle", shape: ShapeRectangle, want: `svc["Service"]`},
		{name: "round", shape: ShapeRound, want: `svc("Service")`},
		{name: "stadium", shape: ShapeStadium, want: `svc(["Service"])`},
		{name: "subroutine", shape: ShapeSubroutine, want: `svc[["Service"]]`},
		{name: "cylinder", shape: ShapeCylinder, want: `svc[("Service")]`},
		{name: "circle", shape: ShapeCircle, want: `svc(("Service"))`},
		{name: "asymmetric", shape: ShapeAsymmetric, want: `svc>"Service"]`},
		{name: "rhombus", shape: ShapeRhombus, want: `svc{"Service"}`},
		{name: "hexagon", shape: ShapeHexagon, want: `svc{{"Service"}}`},
		{name: "parallelogram", shape: ShapeParallelogram, want: `svc[/"Service"/]`},
		{name: "parallelogram alt", shape: ShapeParallelogramAlt, want: `svc[\"Service"\]`},
		{name: "trapezoid", shape: ShapeTrapezoid, want: `svc[/"Service"\]`},
		{name: "trapezoid alt", shape: ShapeTrapezoidAlt, want: `svc[\"Service"/]`},
		{name: "double circle", shape: ShapeDoubleCircle, want: `svc((("Service")))`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			token := Node("svc", WithNodeLabel("Service"), WithNodeShape(tt.shape))
			if token.err != nil {
				t.Fatalf("unexpected error: %v", token.err)
			}
			if diff := cmp.Diff(tt.want, token.value); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}

	token := Node("svc", WithNodeShape(ShapeCircle))
	if token.err != nil {
		t.Fatalf("unexpected error: %v", token.err)
	}
	if diff := cmp.Diff(`svc(("svc"))`, token.value); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}

	if got := formatNodeToken("svc", "Service", Shape("unsupported")); got != `svc["Service"]` {
		t.Fatalf("unexpected fallback shape output: %s", got)
	}
}

func TestTokenValidation(t *testing.T) {
	t.Parallel()

	valid := []Token{
		Literal("  Frontend  "),
		Space(),
		Space(1),
		Space(2),
		Node("svc", WithNodeSpan(2)),
		ArrowLeft("toLeft"),
		ArrowUp("toUp"),
		ArrowX("toX"),
		ArrowY("toY"),
	}

	for _, token := range valid {
		if token.err != nil {
			t.Fatalf("unexpected error: %v", token.err)
		}
	}

	invalid := []Token{
		Literal(""),
		Literal("line1\nline2"),
		Space(0),
		Space(1, 2),
		Node("bad id"),
		Node("svc", WithNodeSpan(0)),
		Node("svc", WithNodeLabel("line1\nline2")),
		Node("svc", WithNodeShape(Shape("unsupported"))),
		Arrow("to", Direction("invalid")),
		Arrow("to", DirectionRight, WithArrowLabel("")),
		Arrow("to", DirectionRight, WithArrowLabel("line1\nline2")),
		Arrow("to", DirectionRight, WithArrowSecondaryDirection(Direction("invalid"))),
	}

	for _, token := range invalid {
		if token.err == nil {
			t.Fatal("expected token error, got nil")
		}
	}
}

func TestDiagram_RowAndStatement(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		LF().
		Row(Literal("Frontend"), Space(1), Literal("Backend")).
		Row(Space(2)).
		Row(Node("svc", WithNodeSpan(2))).
		Statement("raw statement")

	if d.Error() != nil {
		t.Fatalf("unexpected error: %v", d.Error())
	}

	want := `block

    Frontend space Backend
    space:2
    svc:2
    raw statement`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}

	errCases := []struct {
		name string
		run  func() *Diagram
	}{
		{
			name: "row empty trimmed token",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Row(Token{value: "   "})
			},
		},
		{
			name: "row token with newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Row(Token{value: "a\nb"})
			},
		},
		{
			name: "statement empty",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Statement(" ")
			},
		},
		{
			name: "statement with newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Statement("a\nb")
			},
		},
	}

	for _, tt := range errCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diagram := tt.run()
			if diagram.Error() == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestDiagram_LinkStyleAndClass(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		Link("Frontend", "Backend").
		LinkWithLabel("Frontend", "calls", "Backend").
		Style("Frontend,Backend", "fill:#9cf").
		ClassDef("service", "fill:#9cf").
		Class("Frontend,Backend", "service")

	if d.Error() != nil {
		t.Fatalf("unexpected error: %v", d.Error())
	}

	want := `block
    Frontend --> Backend
    Frontend -- "calls" --> Backend
    style Frontend,Backend fill:#9cf
    classDef service fill:#9cf
    class Frontend,Backend service`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}

	errCases := []struct {
		name string
		run  func() *Diagram
	}{
		{
			name: "link invalid source",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Link("bad id", "B")
			},
		},
		{
			name: "link invalid destination",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Link("A", "bad id")
			},
		},
		{
			name: "link label newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).LinkWithLabel("A", "line1\nline2", "B")
			},
		},
		{
			name: "style with newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Style("A", "line1\nline2")
			},
		},
		{
			name: "classdef with newline class name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).ClassDef("ser\nvice", "fill:#9cf")
			},
		},
		{
			name: "class with newline class name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Class("A", "ser\nvice")
			},
		},
	}

	for _, tt := range errCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diagram := tt.run()
			if diagram.Error() == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestDiagram_BuildWithExistingError(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).Columns(0)
	if d.Error() == nil {
		t.Fatal("expected setup error, got nil")
	}

	err := d.Build()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(d.Error(), err) {
		t.Fatalf("expected Error() to wrap returned error, got %v", d.Error())
	}
}
