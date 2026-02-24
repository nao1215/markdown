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
