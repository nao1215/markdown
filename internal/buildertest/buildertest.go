// Package buildertest exercises the error handling that every builder in this
// module shares.
//
// The root package and the seventeen mermaid subpackages each expose a builder
// with the same shape: methods that chain, a String that returns the document
// so far, and a Build that writes it. They also share an error handling
// contract, which callers depend on without it having been written down: no
// call in a chain has to be checked, nothing panics on bad input, and Build is
// the one place an error surfaces.
//
// Keeping the assertions here means each package states which builder it has
// rather than restating the contract, and a change to the contract is one edit
// instead of eighteen.
package buildertest

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

// errFailingWrite is what FailingWriter returns. Tests match on it with
// errors.Is, so a builder that loses the cause instead of wrapping it fails.
var errFailingWrite = errors.New("buildertest: the writer refused the document")

// FailingWriter is an io.Writer whose every write fails. It stands in for a
// closed file or a full disk, which is the failure a builder cannot avoid and
// has to report.
type FailingWriter struct{}

// Write always fails.
func (FailingWriter) Write([]byte) (int, error) {
	return 0, errFailingWrite
}

// Builder is the part of a builder that the shared contract covers. Error is
// left out on purpose: er.Diagram and piechart.PieChart do not have it.
type Builder interface {
	// Build writes the document to the writer the builder was constructed with.
	Build() error
	// String returns the document built so far.
	String() string
}

// New constructs the builder under test, writing to w. The returned builder
// should already have some content, so that a document that fails to be written
// is distinguishable from an empty one.
type New func(w io.Writer) Builder

// RunBuildContract asserts the error handling every builder in this module
// shares. Call it from a package's own test with a constructor for that
// package's builder.
func RunBuildContract(t *testing.T, newBuilder New) {
	t.Helper()

	t.Run("a nil writer is reported rather than panicking", func(t *testing.T) {
		t.Parallel()

		err := newBuilder(nil).Build()
		if err == nil {
			t.Fatal("Build() = nil, want an error: a builder with no destination cannot have written anything")
		}
	})

	t.Run("a failing writer is reported and keeps its cause", func(t *testing.T) {
		t.Parallel()

		err := newBuilder(FailingWriter{}).Build()
		if err == nil {
			t.Fatal("Build() = nil, want an error")
		}
		if !errors.Is(err, errFailingWrite) {
			t.Errorf("Build() = %v, want an error wrapping the writer failure", err)
		}
	})

	t.Run("String returns the document without a writer", func(t *testing.T) {
		t.Parallel()

		// This is how a mermaid diagram reaches CodeBlocks: the diagram is built
		// against no writer at all and handed over as a string.
		if got := newBuilder(nil).String(); got == "" {
			t.Error("String() = \"\", want the document built so far")
		}
	})

	t.Run("Build writes the document String returns", func(t *testing.T) {
		t.Parallel()

		buf := &bytes.Buffer{}
		builder := newBuilder(buf)
		want := builder.String()

		if err := builder.Build(); err != nil {
			t.Fatalf("Build() = %v, want nil", err)
		}
		// The root builder appends the trailing line feed that markdownlint
		// MD047 wants, so the written document starts with, rather than equals,
		// what String returned.
		if !strings.HasPrefix(buf.String(), want) {
			t.Errorf("Build() wrote %q, want it to start with the document String() returned, %q", buf.String(), want)
		}
	})

	t.Run("Build twice writes the document twice", func(t *testing.T) {
		t.Parallel()

		buf := &bytes.Buffer{}
		builder := newBuilder(buf)

		if err := builder.Build(); err != nil {
			t.Fatalf("first Build() = %v, want nil", err)
		}
		once := buf.Len()

		if err := builder.Build(); err != nil {
			t.Fatalf("second Build() = %v, want nil", err)
		}
		if want := once + once; buf.Len() != want {
			t.Errorf("two Build() calls wrote %d bytes, want %d: Build appends, it does not replace", buf.Len(), want)
		}
	})
}

// RunRecordedErrorContract asserts that an error recorded while the chain ran
// surfaces from Build, and that the chain got that far without panicking.
//
// Call it with a constructor that makes the builder reject something, for
// example a name the diagram syntax cannot express. Packages whose builders
// never record an error have nothing to pass here.
func RunRecordedErrorContract(t *testing.T, newRejectingBuilder New) {
	t.Helper()

	t.Run("an error recorded while building surfaces from Build", func(t *testing.T) {
		t.Parallel()

		if err := newRejectingBuilder(io.Discard).Build(); err == nil {
			t.Error("Build() = nil, want the error the chain recorded")
		}
	})

	t.Run("the rest of the chain still runs after an error", func(t *testing.T) {
		t.Parallel()

		// Nothing after a rejected call may panic: callers write one long chain
		// and only check the end of it.
		builder := newRejectingBuilder(io.Discard)
		_ = builder.String()
		_ = builder.Build()
	})
}
