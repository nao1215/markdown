package state

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
)

func TestDiagram_Build(t *testing.T) {
	t.Parallel()

	t.Run("Build a simple state diagram", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		d := NewDiagram(b)
		d.StartTransition("Still").
			Transition("Still", "Moving").
			Transition("Moving", "Crash").
			EndTransition("Crash")

		if err := d.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `stateDiagram-v2
    [*] --> Still
    Still --> Moving
    Moving --> Crash
    Crash --> [*]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a state diagram with title", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		d := NewDiagram(b, WithTitle("Simple State Diagram"))
		d.StartTransition("Still").
			Transition("Still", "Moving").
			EndTransition("Moving")

		if err := d.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `---
title: "Simple State Diagram"
---
stateDiagram-v2
    [*] --> Still
    Still --> Moving
    Moving --> [*]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a state diagram with transition notes", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		d := NewDiagram(b)
		d.StartTransition("Still").
			TransitionWithNote("Still", "Moving", "A transition").
			EndTransition("Moving")

		if err := d.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `stateDiagram-v2
    [*] --> Still
    Still --> Moving : A transition
    Moving --> [*]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a state diagram with state descriptions", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		d := NewDiagram(b)
		d.State("s1", "This is state 1").
			State("s2", "This is state 2").
			StartTransition("s1").
			Transition("s1", "s2").
			EndTransition("s2")

		if err := d.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `stateDiagram-v2
    s1 : This is state 1
    s2 : This is state 2
    [*] --> s1
    s1 --> s2
    s2 --> [*]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a state diagram with notes", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		d := NewDiagram(b)
		d.StartTransition("First").
			NoteLeft("First", "This is a note").
			Transition("First", "Second").
			NoteRight("Second", "Another note").
			EndTransition("Second")

		if err := d.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `stateDiagram-v2
    [*] --> First
    note left of First : This is a note
    First --> Second
    note right of Second : Another note
    Second --> [*]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a state diagram with composite state", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		d := NewDiagram(b)
		d.StartTransition("First").
			CompositeState("First").
			StartTransition("fir").
			Transition("fir", "sec").
			EndTransition("sec").
			End().
			EndTransition("First")

		if err := d.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `stateDiagram-v2
    [*] --> First
    state First {
        [*] --> fir
        fir --> sec
        sec --> [*]
    }
    First --> [*]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a state diagram with fork and join", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		d := NewDiagram(b)
		d.StartTransition("Start").
			Fork("fork_state").
			Transition("Start", "fork_state").
			Transition("fork_state", "State2").
			Transition("fork_state", "State3").
			Join("join_state").
			Transition("State2", "join_state").
			Transition("State3", "join_state").
			EndTransition("join_state")

		if err := d.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `stateDiagram-v2
    [*] --> Start
    state fork_state <<fork>>
    Start --> fork_state
    fork_state --> State2
    fork_state --> State3
    state join_state <<join>>
    State2 --> join_state
    State3 --> join_state
    join_state --> [*]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a state diagram with choice", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		d := NewDiagram(b)
		d.StartTransition("First").
			Transition("First", "Second").
			Choice("second_choice").
			Transition("Second", "second_choice").
			TransitionWithNote("second_choice", "Third", "if n > 0").
			TransitionWithNote("second_choice", "Fourth", "if n <= 0").
			EndTransition("Third").
			EndTransition("Fourth")

		if err := d.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `stateDiagram-v2
    [*] --> First
    First --> Second
    state second_choice <<choice>>
    Second --> second_choice
    second_choice --> Third : if n > 0
    second_choice --> Fourth : if n <= 0
    Third --> [*]
    Fourth --> [*]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a state diagram with direction", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		d := NewDiagram(b)
		d.SetDirection(DirectionLR).
			StartTransition("A").
			Transition("A", "B").
			EndTransition("B")

		if err := d.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `stateDiagram-v2
    direction LR
    [*] --> A
    A --> B
    B --> [*]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagram_String(t *testing.T) {
	t.Parallel()

	t.Run("String returns the state diagram body", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		d := NewDiagram(b)
		d.StartTransition("A").Transition("A", "B")

		want := `stateDiagram-v2
    [*] --> A
    A --> B`

		got := strings.ReplaceAll(d.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagram_Error(t *testing.T) {
	t.Parallel()

	t.Run("Error returns nil when no error", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)
		d := NewDiagram(b)

		if err := d.Error(); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

func TestDiagram_NoteMultiLine(t *testing.T) {
	t.Parallel()

	t.Run("Build a state diagram with multi-line notes", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		d := NewDiagram(b)
		d.StartTransition("First").
			NoteLeftMultiLine("First", "This is a", "multi-line note").
			Transition("First", "Second").
			NoteRightMultiLine("Second", "Another", "multi-line", "note").
			EndTransition("Second")

		if err := d.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `stateDiagram-v2
    [*] --> First
    note left of First
        This is a
        multi-line note
    end note
    First --> Second
    note right of Second
        Another
        multi-line
        note
    end note
    Second --> [*]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

// statements renders the diagram and returns its lines without the header.
func statements(t *testing.T, fn func(*Diagram) *Diagram) []string {
	t.Helper()

	lines := strings.Split(fn(NewDiagram(nil)).String(), internal.LineFeed())
	if len(lines) == 0 {
		t.Fatal("diagram produced no output")
	}
	return lines[1:]
}

// TestStateAndTransitions covers the statement builders that had no test: a
// wrong separator or a missing label here would render a valid-looking but
// wrong diagram, which is the failure mode nobody notices.
func TestStateAndTransitions(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(*Diagram) *Diagram
		want  string
	}{
		"state without a description": {
			build: func(d *Diagram) *Diagram { return d.State("Idle", "") },
			want:  "    Idle",
		},
		"state with a description": {
			build: func(d *Diagram) *Diagram { return d.State("Idle", "waiting for work") },
			want:  "    Idle : waiting for work",
		},
		"transition with a note": {
			build: func(d *Diagram) *Diagram { return d.TransitionWithNote("Idle", "Busy", "job arrives") },
			want:  "    Idle --> Busy : job arrives",
		},
		"start transition with a note": {
			build: func(d *Diagram) *Diagram { return d.StartTransitionWithNote("Idle", "boot") },
			want:  "    [*] --> Idle : boot",
		},
		"end transition with a note": {
			build: func(d *Diagram) *Diagram { return d.EndTransitionWithNote("Done", "shutdown") },
			want:  "    Done --> [*] : shutdown",
		},
		"concurrent separator": {
			build: func(d *Diagram) *Diagram { return d.Concurrent() },
			want:  "    ---",
		},
		"line feed": {
			build: func(d *Diagram) *Diagram { return d.LF() },
			want:  "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := statements(t, tt.build)
			if len(got) == 0 || got[0] != tt.want {
				t.Errorf("statement mismatch:\n got: %#v\nwant: %q", got, tt.want)
			}
		})
	}
}

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).State("Draft", "The order is being written")
	})
}

// TestGoldenState pins the rendered diagram of every builder method of this
// package, including the composite state builder.
func TestGoldenState(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(buf, WithTitle("Order Lifecycle")).
		SetDirection(DirectionLR).
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
	err := NewDiagram(buf).
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
		direction Direction
	}{
		{name: "left to right", golden: "direction_lr.md", direction: DirectionLR},
		{name: "right to left", golden: "direction_rl.md", direction: DirectionRL},
		{name: "top to bottom", golden: "direction_tb.md", direction: DirectionTB},
		{name: "bottom to top", golden: "direction_bt.md", direction: DirectionBT},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			err := NewDiagram(buf).
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
