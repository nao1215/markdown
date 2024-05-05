package sequence

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDiagramLoopStartEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add loop to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.LoopStart("description")
		d.LoopEnd()

		want := []string{"sequenceDiagram", "    loop description", "    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAltStartElseEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add alt to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AltStart("description")
		d.AltElse("description")
		d.AltEnd()

		want := []string{"sequenceDiagram", "    alt description", "    else description", "    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramOptStartEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add opt to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.OptStart("description")
		d.OptEnd()

		want := []string{"sequenceDiagram", "    opt description", "    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramParallelStartAndEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add parallel to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.ParallelStart("start")
		d.ParallelAnd("and-description")
		d.ParallelEnd()

		want := []string{"sequenceDiagram", "    par start", "    and and-description", "    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramCriticalStartAndEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add critical to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.CriticalStart("start")
		d.CriticalOption("option-description")
		d.CriticalEnd()

		want := []string{
			"sequenceDiagram",
			"    critical start",
			"    option option-description",
			"    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramBreakStartEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add break to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.BreakStart("description")
		d.BreakEnd()

		want := []string{"sequenceDiagram", "    break description", "    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}
