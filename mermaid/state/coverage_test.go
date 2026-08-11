package state_test

import (
	"strings"
	"testing"

	"github.com/nao1215/markdown/mermaid/state"
)

// statements renders the diagram and returns its lines without the header.
func statements(t *testing.T, fn func(*state.Diagram) *state.Diagram) []string {
	t.Helper()

	lines := strings.Split(fn(state.NewDiagram(nil)).String(), "\n")
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
		build func(*state.Diagram) *state.Diagram
		want  string
	}{
		"state without a description": {
			build: func(d *state.Diagram) *state.Diagram { return d.State("Idle", "") },
			want:  "    Idle",
		},
		"state with a description": {
			build: func(d *state.Diagram) *state.Diagram { return d.State("Idle", "waiting for work") },
			want:  "    Idle : waiting for work",
		},
		"transition with a note": {
			build: func(d *state.Diagram) *state.Diagram { return d.TransitionWithNote("Idle", "Busy", "job arrives") },
			want:  "    Idle --> Busy : job arrives",
		},
		"start transition with a note": {
			build: func(d *state.Diagram) *state.Diagram { return d.StartTransitionWithNote("Idle", "boot") },
			want:  "    [*] --> Idle : boot",
		},
		"end transition with a note": {
			build: func(d *state.Diagram) *state.Diagram { return d.EndTransitionWithNote("Done", "shutdown") },
			want:  "    Done --> [*] : shutdown",
		},
		"concurrent separator": {
			build: func(d *state.Diagram) *state.Diagram { return d.Concurrent() },
			want:  "    ---",
		},
		"line feed": {
			build: func(d *state.Diagram) *state.Diagram { return d.LF() },
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
