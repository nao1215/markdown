// Package state is mermaid state diagram builder.
package state

import (
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// Diagram is a state diagram builder.
type Diagram struct {
	// body is state diagram body.
	body []string
	// config is the configuration for the state diagram.
	config *config
	// dest is output destination for state diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the state diagram building.
	err error
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := []string{}
	if c.title != noTitle {
		lines = append(lines, "---")
		lines = append(lines, fmt.Sprintf("title: %s", c.title))
		lines = append(lines, "---")
	}
	lines = append(lines, "stateDiagram-v2")

	return &Diagram{
		body:   lines,
		dest:   w,
		config: c,
	}
}

// String returns the state diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the state diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the state diagram body to the output destination.
func (d *Diagram) Build() error {
	if _, err := fmt.Fprint(d.dest, d.String()); err != nil {
		if d.err != nil {
			return fmt.Errorf("failed to write: %w: %s", err, d.err.Error()) //nolint:wrapcheck
		}
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

// State adds a state to the state diagram.
func (d *Diagram) State(id, description string) *Diagram {
	if description == "" {
		d.body = append(d.body, fmt.Sprintf("    %s", id))
	} else {
		d.body = append(d.body, fmt.Sprintf("    %s : %s", id, description))
	}
	return d
}

// Transition adds a transition between states.
func (d *Diagram) Transition(from, to string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s --> %s", from, to))
	return d
}

// TransitionWithNote adds a transition between states with a note.
func (d *Diagram) TransitionWithNote(from, to, note string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s --> %s : %s", from, to, note))
	return d
}

// StartTransition adds a transition from the start state.
func (d *Diagram) StartTransition(to string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    [*] --> %s", to))
	return d
}

// StartTransitionWithNote adds a transition from the start state with a note.
func (d *Diagram) StartTransitionWithNote(to, note string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    [*] --> %s : %s", to, note))
	return d
}

// EndTransition adds a transition to the end state.
func (d *Diagram) EndTransition(from string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s --> [*]", from))
	return d
}

// EndTransitionWithNote adds a transition to the end state with a note.
func (d *Diagram) EndTransitionWithNote(from, note string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s --> [*] : %s", from, note))
	return d
}

// NoteLeft adds a note on the left side of a state.
func (d *Diagram) NoteLeft(state, note string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    note left of %s : %s", state, note))
	return d
}

// NoteRight adds a note on the right side of a state.
func (d *Diagram) NoteRight(state, note string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    note right of %s : %s", state, note))
	return d
}

// NoteLeftMultiLine adds a multi-line note on the left side of a state.
func (d *Diagram) NoteLeftMultiLine(state string, lines ...string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    note left of %s", state))
	for _, line := range lines {
		d.body = append(d.body, fmt.Sprintf("        %s", line))
	}
	d.body = append(d.body, "    end note")
	return d
}

// NoteRightMultiLine adds a multi-line note on the right side of a state.
func (d *Diagram) NoteRightMultiLine(state string, lines ...string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    note right of %s", state))
	for _, line := range lines {
		d.body = append(d.body, fmt.Sprintf("        %s", line))
	}
	d.body = append(d.body, "    end note")
	return d
}

// CompositeState starts a composite state definition.
func (d *Diagram) CompositeState(id string) *CompositeStateBuilder {
	return &CompositeStateBuilder{
		diagram: d,
		id:      id,
	}
}

// CompositeStateBuilder is a builder for composite states.
type CompositeStateBuilder struct {
	diagram *Diagram
	id      string
	body    []string
}

// State adds a state to the composite state.
func (b *CompositeStateBuilder) State(id, description string) *CompositeStateBuilder {
	if description == "" {
		b.body = append(b.body, fmt.Sprintf("        %s", id))
	} else {
		b.body = append(b.body, fmt.Sprintf("        %s : %s", id, description))
	}
	return b
}

// Transition adds a transition between states in the composite state.
func (b *CompositeStateBuilder) Transition(from, to string) *CompositeStateBuilder {
	b.body = append(b.body, fmt.Sprintf("        %s --> %s", from, to))
	return b
}

// TransitionWithNote adds a transition with a note in the composite state.
func (b *CompositeStateBuilder) TransitionWithNote(from, to, note string) *CompositeStateBuilder {
	b.body = append(b.body, fmt.Sprintf("        %s --> %s : %s", from, to, note))
	return b
}

// StartTransition adds a transition from the start state in the composite state.
func (b *CompositeStateBuilder) StartTransition(to string) *CompositeStateBuilder {
	b.body = append(b.body, fmt.Sprintf("        [*] --> %s", to))
	return b
}

// EndTransition adds a transition to the end state in the composite state.
func (b *CompositeStateBuilder) EndTransition(from string) *CompositeStateBuilder {
	b.body = append(b.body, fmt.Sprintf("        %s --> [*]", from))
	return b
}

// End ends the composite state definition and returns the diagram.
func (b *CompositeStateBuilder) End() *Diagram {
	b.diagram.body = append(b.diagram.body, fmt.Sprintf("    state %s {", b.id))
	b.diagram.body = append(b.diagram.body, b.body...)
	b.diagram.body = append(b.diagram.body, "    }")
	return b.diagram
}

// Fork adds a fork to the state diagram.
func (d *Diagram) Fork(id string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    state %s <<fork>>", id))
	return d
}

// Join adds a join to the state diagram.
func (d *Diagram) Join(id string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    state %s <<join>>", id))
	return d
}

// Choice adds a choice to the state diagram.
func (d *Diagram) Choice(id string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    state %s <<choice>>", id))
	return d
}

// LF adds a line feed to the state diagram.
func (d *Diagram) LF() *Diagram {
	d.body = append(d.body, "")
	return d
}

// Direction sets the direction of the state diagram.
type Direction string

const (
	// DirectionLR is left to right direction.
	DirectionLR Direction = "LR"
	// DirectionRL is right to left direction.
	DirectionRL Direction = "RL"
	// DirectionTB is top to bottom direction.
	DirectionTB Direction = "TB"
	// DirectionBT is bottom to top direction.
	DirectionBT Direction = "BT"
)

// SetDirection sets the direction of the state diagram.
func (d *Diagram) SetDirection(dir Direction) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    direction %s", dir))
	return d
}

// Concurrent adds a concurrent state (indicated by ---).
func (d *Diagram) Concurrent() *Diagram {
	d.body = append(d.body, "    ---")
	return d
}
