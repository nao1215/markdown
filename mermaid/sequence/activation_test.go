package sequence

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDiagramActivateDeactivate(t *testing.T) {
	t.Parallel()

	t.Run("should add activate and deactivate to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.Activate("Alice")
		d.Deactivate("Alice")

		want := []string{"sequenceDiagram", "    activate Alice", "    deactivate Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramRequestfWithActivation(t *testing.T) {
	t.Parallel()

	t.Run("should add request to the sequence diagram with activation", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.SyncRequestfWithActivation("Alice", "Bob", "Hello %s", "Bob")

		want := []string{"sequenceDiagram", "    Alice->>+Bob: Hello Bob"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramResponsefWithActivation(t *testing.T) {
	t.Parallel()

	t.Run("should add response to the sequence diagram with activation", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.SyncResponsefWithActivation("Alice", "Bob", "Hello %s", "Alice")

		want := []string{"sequenceDiagram", "    Alice-->>-Bob: Hello Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAsyncRequestfWithActivation(t *testing.T) {
	t.Parallel()

	t.Run("should add async request to the sequence diagram with activation", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AsyncRequestfWithActivation("Alice", "Bob", "Hello %s", "Bob")

		want := []string{"sequenceDiagram", "    Alice->>+Bob: Hello Bob"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAsyncResponsefWithActivation(t *testing.T) {
	t.Parallel()

	t.Run("should add async response to the sequence diagram with activation", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AsyncResponsefWithActivation("Alice", "Bob", "Hello %s", "Alice")

		want := []string{"sequenceDiagram", "    Alice-->>-Bob: Hello Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}
