package sequence

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDiagramParticipant(t *testing.T) {
	t.Parallel()

	t.Run("should add participant to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.Participant("Alice")

		want := []string{"sequenceDiagram", "    participant Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramCreateDeleteParticipant(t *testing.T) {
	t.Parallel()

	t.Run("should add create and delete participant to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.CreateParticipant("Alice")
		d.DestroyParticipant("Alice")

		want := []string{"sequenceDiagram", "    create participant Alice", "    destroy Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramCreateDeleteActor(t *testing.T) {
	t.Parallel()

	t.Run("should add create and delete actor to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.CreateActor("Alice")
		d.DestroyActor("Alice")

		want := []string{"sequenceDiagram", "    create actor Alice", "    destroy Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramActor(t *testing.T) {
	t.Parallel()

	t.Run("should add actor to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.Actor("Alice")

		want := []string{"sequenceDiagram", "    actor Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAutoNumber(t *testing.T) {
	t.Parallel()

	t.Run("should add autonumber to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AutoNumber()

		want := []string{"sequenceDiagram", "    autonumber"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramBoxStartEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add box to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.Participant("Alice").Participant("Bob")
		d.BoxStart([]string{"Alice", "Bob"})
		d.BoxEnd()

		want := []string{
			"sequenceDiagram",
			"    participant Alice",
			"    participant Bob",
			"    box Alice & Bob",
			"    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}
