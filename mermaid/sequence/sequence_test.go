// Package sequence is mermaid sequence diagram builder.
package sequence

import (
	"fmt"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestString(t *testing.T) {
	t.Parallel()

	t.Run("should return the sequence diagram body", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.Participant("Alice")

		want := fmt.Sprintf("sequenceDiagram%s    participant Alice", lineFeed())
		got := d.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("expected and actual are different: %s", diff)
		}
	})
}

func TestDiagramRequestf(t *testing.T) {
	t.Parallel()

	t.Run("should add request to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.SyncRequestf("Alice", "Bob", "Hello %s", "Bob")

		want := []string{"sequenceDiagram", "    Alice->>Bob: Hello Bob"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramResponsef(t *testing.T) {
	t.Parallel()

	t.Run("should add response to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.SyncResponsef("Alice", "Bob", "Hello %s", "Alice")

		want := []string{"sequenceDiagram", "    Alice-->>Bob: Hello Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramRequestErrorf(t *testing.T) {
	t.Parallel()

	t.Run("should add request error to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.RequestErrorf("Alice", "Bob", "Hello %s", "Bob")

		want := []string{"sequenceDiagram", "    Alice-xBob: Hello Bob"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramResponseErrorf(t *testing.T) {
	t.Parallel()

	t.Run("should add response error to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.ResponseErrorf("Alice", "Bob", "Hello %s", "Alice")

		want := []string{"sequenceDiagram", "    Alice--xBob: Hello Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAsyncRequestf(t *testing.T) {
	t.Parallel()

	t.Run("should add async request to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AsyncRequestf("Alice", "Bob", "Hello %s", "Bob")

		want := []string{"sequenceDiagram", "    Alice->)Bob: Hello Bob"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAsyncResponsef(t *testing.T) {
	t.Parallel()

	t.Run("should add async response to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AsyncResponsef("Alice", "Bob", "Hello %s", "Alice")

		want := []string{"sequenceDiagram", "    Alice--)Bob: Hello Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramError(t *testing.T) {
	t.Parallel()

	t.Run("should return the error", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.err = fmt.Errorf("error")

		if d.Error().Error() != "error" {
			t.Error("value is mismatch, want error")
		}
	})
}
