package state

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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
title: Simple State Diagram
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
