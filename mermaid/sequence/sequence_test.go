// Package sequence is mermaid sequence diagram builder.
package sequence

import (
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal"
)

func TestString(t *testing.T) {
	t.Parallel()

	t.Run("should return the sequence diagram body", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.Participant("Alice")

		want := fmt.Sprintf("sequenceDiagram%s    participant Alice", internal.LineFeed())
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

func TestNewDiagram(t *testing.T) {
	t.Parallel()

	t.Run("with all options", func(t *testing.T) {
		t.Parallel()

		got := NewDiagram(
			io.Discard,
			WithMirrorActors(true),
			WithBottomMariginAdjustment(2),
			WithActorFontSize(12),
			WithActorFontFamily("Arial"),
			WithActorFontWeight("bold"),
			WithNoteFontFamily("Arial"),
			WithNoteFontSize(12),
			WithNoteFontWeight("bold"),
			WithNoteAlign("left"),
			WithMessageFontFamily("Arial"),
			WithMessageFontSize(12),
			WithMessageFontWeight("bold"),
		)

		want := &Diagram{
			body: []string{"sequenceDiagram"},
			dest: io.Discard,
			config: &config{
				mirrorActors:            true,
				bottomMariginAdjustment: 2,
				actorFontSize:           12,
				actorFontFamily:         "Arial",
				actorFontWeight:         "bold",
				noteFontSize:            12,
				noteFontFamily:          "Arial",
				noteFontWeight:          "bold",
				noteAlign:               "left",
				messageFontSize:         12,
				messageFontFamily:       "Arial",
				messageFontWeight:       "bold",
			},
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("value is mismatch, want %v, got %v", want, got)
		}
	})
}
