package mindmap

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
)

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
title: "Product Strategy Mindmap"
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
title: "Product Strategy Mindmap"
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

	d := NewDiagram(errWriter{})
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

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).Root("Product").Child("Market")
	})
}

// TestRecordedErrorContract asserts that a child added before the root surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).Child("a child with no root")
	})
}

// TestGoldenMindmap pins the rendered diagram of every builder method of this
// package, including the explicit depth form of Node.
func TestGoldenMindmap(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(buf, WithTitle("Product Strategy")).
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
	err := NewDiagram(buf).
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

// TestBuildWithNilWriter covers the case where a diagram is built for String()
// only and Build() is called by mistake. Build() used to dereference the nil
// writer and take the process down; it has to return an error instead.
func TestBuildWithNilWriter(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Build() panicked with a nil writer: %v", r)
		}
	}()

	d := NewDiagram(nil)

	// String() has always worked without a writer, and callers rely on it.
	_ = d.String()

	err := d.Build()
	if err == nil {
		t.Fatal("Build() with a nil writer must return an error")
	}
	if err.Error() != "output writer must not be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// errWrite is the failure the writer below reports, so the test can assert that
// Build passed it through rather than inventing an error of its own.
var errWrite = errors.New("write failed")

// errWriter fails every write, which is what a full disk or a closed pipe looks
// like to Build.
type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errWrite
}

// TestBuildReportsWriteFailure covers the branch where the destination accepts
// the diagram and then fails. Silently returning nil there would hand the caller
// a document that was never written.
func TestBuildReportsWriteFailure(t *testing.T) {
	t.Parallel()

	err := NewDiagram(errWriter{}).Build()
	if err == nil {
		t.Fatal("Build must report a failing writer")
	}
	if !errors.Is(err, errWrite) {
		t.Errorf("Build lost the destination error: %v", err)
	}
}
