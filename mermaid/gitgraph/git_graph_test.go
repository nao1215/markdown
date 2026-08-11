package gitgraph

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
			want: "gitGraph",
		},
		{
			name: "new diagram with title",
			opts: []Option{WithTitle("Release Flow")},
			want: `---
title: "Release Flow"
---
gitGraph`,
		},
		{
			name:    "new diagram with title including newline",
			opts:    []Option{WithTitle("Release\nInjected: malicious")},
			want:    "gitGraph",
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

	d := NewDiagram(b, WithTitle("Release Flow"))
	d.Commit(WithCommitID("init"), WithCommitTag("v0.1.0")).
		Branch("develop", WithBranchOrder(2)).
		Checkout("develop").
		Commit(WithCommitType(CommitTypeHighlight)).
		Checkout("main").
		Merge("develop", WithCommitTag("v1.0.0")).
		Branch("cherry-pick").
		Checkout("cherry-pick").
		Commit(WithCommitID("hotfix")).
		Checkout("main").
		CherryPick("hotfix").
		Reset("hotfix")

	if err := d.Build(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `---
title: "Release Flow"
---
gitGraph
    commit id: "init" tag: "v0.1.0"
    branch develop order: 2
    checkout develop
    commit type: HIGHLIGHT
    checkout main
    merge develop tag: "v1.0.0"
    branch "cherry-pick"
    checkout "cherry-pick"
    commit id: "hotfix"
    checkout main
    cherry-pick id: "hotfix"
    reset id: "hotfix"`

	got := strings.ReplaceAll(b.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_CommitWithoutOptions(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).Commit()
	if d.Error() != nil {
		t.Fatalf("unexpected error: %v", d.Error())
	}

	want := `gitGraph
    commit`
	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_BranchOrderZero(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).Branch("develop", WithBranchOrder(0))
	if d.Error() != nil {
		t.Fatalf("unexpected error: %v", d.Error())
	}

	want := `gitGraph
    branch develop order: 0`
	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
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
			name: "empty branch name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Branch("")
			},
			want: "gitGraph",
		},
		{
			name: "invalid branch order",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Branch("develop", WithBranchOrder(-1))
			},
			want: "gitGraph",
		},
		{
			name: "empty checkout branch name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Checkout("")
			},
			want: "gitGraph",
		},
		{
			name: "empty merge branch name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Merge("")
			},
			want: "gitGraph",
		},
		{
			name: "empty cherry-pick id",
			run: func() *Diagram {
				return NewDiagram(io.Discard).CherryPick("")
			},
			want: "gitGraph",
		},
		{
			name: "empty reset id",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Reset("")
			},
			want: "gitGraph",
		},
		{
			name: "invalid commit type",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Commit(WithCommitType(CommitType("INVALID")))
			},
			want: "gitGraph",
		},
		{
			name: "checkout with newline in branch name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Checkout("feature\nbranch")
			},
			want: "gitGraph",
		},
		{
			name: "commit with newline in commit id",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Commit(WithCommitID("abc\n123"))
			},
			want: "gitGraph",
		},
		{
			name: "merge with newline in branch name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Merge("feature\nbranch")
			},
			want: "gitGraph",
		},
		{
			name: "commit with newline in tag",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Commit(WithCommitTag("v1\n.0"))
			},
			want: "gitGraph",
		},
		{
			name: "newline in cherry-pick parent",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					CherryPick("abc123", WithCherryPickParent("def\n456"))
			},
			want: "gitGraph",
		},
		{
			name: "newline in title",
			run: func() *Diagram {
				return NewDiagram(io.Discard, WithTitle("Release\nInjected: malicious"))
			},
			want: "gitGraph",
		},
		{
			name: "lf short-circuit after error",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Branch("").LF()
			},
			want: "gitGraph",
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
		return NewDiagram(w).Commit()
	})
}

// TestRecordedErrorContract asserts that a branch order the syntax cannot express surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).Commit().Branch("develop", WithBranchOrder(-1))
	})
}

// TestGoldenGitGraph pins the rendered diagram of every builder method and
// every option of this package.
func TestGoldenGitGraph(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(buf, WithTitle("Release Flow")).
		Commit().
		Commit(WithCommitID("init"), WithCommitTag("v0.1.0")).
		Commit(WithCommitType(CommitTypeNormal)).
		Commit(WithCommitID("reverse"), WithCommitType(CommitTypeReverse)).
		Commit(WithCommitID("highlight"), WithCommitType(CommitTypeHighlight)).
		Branch("develop").
		Branch("feature", WithBranchOrder(2)).
		Checkout("develop").
		Commit(WithCommitID("dev-1")).
		LF().
		Checkout("main").
		Merge("develop", WithCommitTag("v1.0.0")).
		CherryPick("dev-1").
		CherryPick("dev-1", WithCherryPickParent("init")).
		Reset("init").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("gitgraph.md", buf.String()); err != nil {
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
