package mindmap

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
			want: "mindmap",
		},
		{
			name: "new diagram with title",
			opts: []Option{WithTitle("Product Strategy Mindmap")},
			want: `---
title: Product Strategy Mindmap
---
mindmap`,
		},
		{
			name:    "new diagram with title including newline",
			opts:    []Option{WithTitle("Product\nStrategy")},
			want:    "mindmap",
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

	d := NewDiagram(b, WithTitle("Product Strategy Mindmap"))
	d.Root("Product Strategy").
		Child("Market").
		Child("SMB").
		Sibling("Enterprise").
		Parent().
		Sibling("Execution").
		Child("Q1").
		Sibling("Q2")

	if err := d.Build(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `---
title: Product Strategy Mindmap
---
mindmap
    Product Strategy
        Market
            SMB
            Enterprise
        Execution
            Q1
            Q2`

	got := strings.ReplaceAll(b.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_Node(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		Root("Product Strategy").
		Node(1, "Market").
		Node(2, "SMB").
		Node(1, "Execution")

	want := `mindmap
    Product Strategy
        Market
            SMB
        Execution`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_ParentBacktrack(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		Root("Root").
		Child("A").
		Child("B").
		Parent().
		Parent().
		Child("C")

	want := `mindmap
    Root
        A
            B
        C`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
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
			name: "child before root",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Child("Market")
			},
			want: "mindmap",
		},
		{
			name: "sibling before root",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Sibling("Execution")
			},
			want: "mindmap",
		},
		{
			name: "parent before root",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Parent()
			},
			want: "mindmap",
		},
		{
			name: "second root",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Root("Product Strategy").
					Root("Another Root")
			},
			want: `mindmap
    Product Strategy`,
		},
		{
			name: "node level jump",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Root("Product Strategy").
					Node(3, "Too deep")
			},
			want: `mindmap
    Product Strategy`,
		},
		{
			name: "negative node level",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Node(-1, "Invalid")
			},
			want: "mindmap",
		},
		{
			name: "empty root text",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Root("")
			},
			want: "mindmap",
		},
		{
			name: "newline in root text",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Root("Product\nStrategy")
			},
			want: "mindmap",
		},
		{
			name: "parent at root level",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Root("Product Strategy").
					Parent()
			},
			want: `mindmap
    Product Strategy`,
		},
		{
			name: "lf short-circuit after error",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Root("").LF()
			},
			want: "mindmap",
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
