package markdown_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/internal/buildertest"
)

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return markdown.NewMarkdown(w).H1("Title")
	})
}

// TestRecordedErrorContract asserts that a table whose rows do not match its
// header surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return markdown.NewMarkdown(w).Table(markdown.TableSet{
			Header: []string{"one", "two"},
			Rows:   [][]string{{"only one cell"}},
		})
	})
}

// TestErrorSurfacesFromBothErrorAndBuild pins that the two ways of asking for
// the error agree. Callers use whichever suits the shape of their code, and a
// builder that reported an error from only one of them would be a trap.
func TestErrorSurfacesFromBothErrorAndBuild(t *testing.T) {
	t.Parallel()

	m := markdown.NewMarkdown(io.Discard).Table(markdown.TableSet{
		Header: []string{"one", "two"},
		Rows:   [][]string{{"only one cell"}},
	})

	fromError := m.Error()
	if fromError == nil {
		t.Fatal("Error() = nil, want the error the chain recorded")
	}
	if !errors.Is(fromError, markdown.ErrMismatchColumn) {
		t.Errorf("Error() = %v, want an error wrapping ErrMismatchColumn", fromError)
	}

	fromBuild := m.Build()
	if fromBuild == nil {
		t.Fatal("Build() = nil, want the error the chain recorded")
	}
	if !errors.Is(fromBuild, markdown.ErrMismatchColumn) {
		t.Errorf("Build() = %v, want an error wrapping ErrMismatchColumn", fromBuild)
	}
}

// TestTheChainContinuesAfterAnError pins that a rejected call does not stop the
// document. Callers write one long chain and check it once at the end, so the
// blocks after a bad table still have to appear.
func TestTheChainContinuesAfterAnError(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := markdown.NewMarkdown(buf).
		H1("Before").
		Table(markdown.TableSet{
			Header: []string{"one", "two"},
			Rows:   [][]string{{"only one cell"}},
		}).
		H2("After").
		Build()
	if err == nil {
		t.Fatal("Build() = nil, want the error the chain recorded")
	}

	for _, want := range []string{"# Before", "## After"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("document = %q, want it to contain %q", buf.String(), want)
		}
	}
}

// TestTheFirstErrorIsKept pins which error a caller sees when a chain records
// more than one. The first is the one that explains the rest, so it is the one
// that survives; the later ones are appended to its message rather than
// replacing it.
func TestTheFirstErrorIsKept(t *testing.T) {
	t.Parallel()

	m := markdown.NewMarkdown(io.Discard).
		TableOfContents(markdown.TableOfContentsDepthH3).
		TableOfContents(markdown.TableOfContentsDepthH3)

	err := m.Error()
	if err == nil {
		t.Fatal("Error() = nil, want the error the chain recorded")
	}
	if want := "table of contents has already been generated"; !strings.Contains(err.Error(), want) {
		t.Errorf("Error() = %v, want it to mention %q", err, want)
	}
}

// TestBuildWithNilWriterKeepsTheEarlierError pins that a nil writer does not
// hide what went wrong before it. Both failures are worth knowing about, so the
// message carries the earlier one too.
func TestBuildWithNilWriterKeepsTheEarlierError(t *testing.T) {
	t.Parallel()

	err := markdown.NewMarkdown(nil).
		Table(markdown.TableSet{
			Header: []string{"one", "two"},
			Rows:   [][]string{{"only one cell"}},
		}).
		Build()
	if err == nil {
		t.Fatal("Build() = nil, want an error")
	}

	for _, want := range []string{"destination writer is nil", markdown.ErrMismatchColumn.Error()} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Build() = %v, want it to mention %q", err, want)
		}
	}
}

// TestBuildWithFailingWriterKeepsTheEarlierError pins the same for a writer
// that refuses the document.
func TestBuildWithFailingWriterKeepsTheEarlierError(t *testing.T) {
	t.Parallel()

	err := markdown.NewMarkdown(buildertest.FailingWriter{}).
		Table(markdown.TableSet{
			Header: []string{"one", "two"},
			Rows:   [][]string{{"only one cell"}},
		}).
		Build()
	if err == nil {
		t.Fatal("Build() = nil, want an error")
	}

	for _, want := range []string{"failed to write markdown text", markdown.ErrMismatchColumn.Error()} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("Build() = %v, want it to mention %q", err, want)
		}
	}
}

// TestBuildEndsTheDocumentWithALineFeed pins the trailing newline. markdownlint
// MD047 requires it, and without it a second document written to the same
// writer would splice its first line onto the last line of this one.
func TestBuildEndsTheDocumentWithALineFeed(t *testing.T) {
	t.Parallel()

	tests := map[string]func(*markdown.Markdown) *markdown.Markdown{
		"a document ending in a heading": func(m *markdown.Markdown) *markdown.Markdown {
			return m.H1("Title")
		},
		"a document ending in a table": func(m *markdown.Markdown) *markdown.Markdown {
			return m.Table(markdown.TableSet{
				Header: []string{"key", "value"},
				Rows:   [][]string{{"a", "b"}},
			})
		},
		"an empty document": func(m *markdown.Markdown) *markdown.Markdown {
			return m
		},
	}

	for name, build := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			if err := build(markdown.NewMarkdown(buf)).Build(); err != nil {
				t.Fatalf("Build() = %v, want nil", err)
			}
			if !strings.HasSuffix(buf.String(), "\n") {
				t.Errorf("document = %q, want it to end with a line feed", buf.String())
			}
		})
	}
}
