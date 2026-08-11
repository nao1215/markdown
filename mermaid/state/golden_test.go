package state_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/state"
)

// TestGoldenState pins the rendered diagram of every builder method of this
// package, including the composite state builder.
func TestGoldenState(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := state.NewDiagram(buf, state.WithTitle("Order Lifecycle")).
		SetDirection(state.DirectionLR).
		State("Draft", "The order is being written").
		State("Placed", "The order has been placed").
		State("Shipped", "The order is on its way").
		StartTransition("Draft").
		StartTransitionWithNote("Placed", "restored from a backup").
		Transition("Draft", "Placed").
		TransitionWithNote("Placed", "Shipped", "after payment").
		EndTransition("Shipped").
		EndTransitionWithNote("Placed", "canceled by the customer").
		LF().
		NoteLeft("Draft", "editable").
		NoteRight("Shipped", "immutable").
		NoteLeftMultiLine("Placed", "waiting for payment", "then for the warehouse").
		NoteRightMultiLine("Draft", "no payment yet", "no reservation yet").
		LF().
		Fork("split").
		Join("merge").
		Choice("decide").
		Concurrent().
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("state.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenStateComposite pins the nested block that CompositeState builds,
// which is the only place in this package where a second builder type appears.
func TestGoldenStateComposite(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := state.NewDiagram(buf).
		State("Active", "The order is active").
		CompositeState("Active").
		State("Reserved", "Stock is reserved").
		Transition("Reserved", "Packed").
		TransitionWithNote("Packed", "Handed over", "to the carrier").
		StartTransition("Reserved").
		EndTransition("Handed over").
		End().
		Transition("Active", "Closed").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("state_composite.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenStateDirections pins the header each direction produces. Only one
// direction applies to a diagram, so each needs its own.
func TestGoldenStateDirections(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		golden    string
		direction state.Direction
	}{
		{name: "left to right", golden: "direction_lr.md", direction: state.DirectionLR},
		{name: "right to left", golden: "direction_rl.md", direction: state.DirectionRL},
		{name: "top to bottom", golden: "direction_tb.md", direction: state.DirectionTB},
		{name: "bottom to top", golden: "direction_bt.md", direction: state.DirectionBT},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			err := state.NewDiagram(buf).
				SetDirection(tt.direction).
				State("Draft", "The order is being written").
				Transition("Draft", "Placed").
				Build()
			if err != nil {
				t.Fatalf("Build() = %v, want nil", err)
			}

			if err := golden.Assert(tt.golden, buf.String()); err != nil {
				t.Error(err)
			}
		})
	}
}
