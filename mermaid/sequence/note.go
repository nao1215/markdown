package sequence

import "fmt"

// NotePosition is a note position.
type NotePosition string

const (
	// NotePositionOver is a note position.
	NotePositionOver NotePosition = "over"
	// NotePositionRight is a note position.
	NotePositionRight NotePosition = "right of"
	// NotePositionLeft is a note position.
	NotePositionLeft NotePosition = "left of"
)

// NoteOver add a note to the sequence diagram.
func (d *Diagram) NoteOver(participant, message string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    note over %s: %s", participant, message))
	return d
}

// NoteRightOf add a note to the sequence diagram.
func (d *Diagram) NoteRightOf(participant, message string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    note right of %s: %s", participant, message))
	return d
}

// NoteLeftOf add a note to the sequence diagram.
func (d *Diagram) NoteLeftOf(participant, message string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    note left of %s: %s", participant, message))
	return d
}
