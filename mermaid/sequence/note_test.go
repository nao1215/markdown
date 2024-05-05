package sequence

import (
	"io"
	"reflect"
	"testing"
)

func TestDiagramNoteOver(t *testing.T) {
	t.Parallel()

	t.Run("should add note over to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.NoteOver("Alice", "Hello Alice")

		want := []string{"sequenceDiagram", "    note over Alice: Hello Alice"}
		got := d.body

		if !reflect.DeepEqual(want, got) {
			t.Errorf("value is mismatch want:%v got:%v", want, got)
		}
	})
}

func TestDiagramNoteRightOf(t *testing.T) {
	t.Parallel()

	t.Run("should add note right of to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.NoteRightOf("Alice", "Hello Alice")

		want := []string{"sequenceDiagram", "    note right of Alice: Hello Alice"}
		got := d.body

		if !reflect.DeepEqual(want, got) {
			t.Errorf("value is mismatch want:%v got:%v", want, got)
		}
	})
}

func TestDiagramNoteLeftOf(t *testing.T) {
	t.Parallel()

	t.Run("should add note left of to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.NoteLeftOf("Alice", "Hello Alice")

		want := []string{"sequenceDiagram", "    note left of Alice: Hello Alice"}
		got := d.body

		if !reflect.DeepEqual(want, got) {
			t.Errorf("value is mismatch want:%v got:%v", want, got)
		}
	})
}
