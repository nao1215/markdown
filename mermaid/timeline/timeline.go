// Package timeline is mermaid timeline diagram builder.
//
// Ref. https://mermaid.js.org/syntax/timeline.html
//
// A timeline is a list of time periods, each holding one or more events, and
// optionally grouped into sections. A period is written once and its events
// follow it, so the builder mirrors that: Period starts one, and Event adds to
// whichever period is current.
//
// Errors are recorded rather than returned from every call: the chain runs to
// the end and the error surfaces from Build. A nil writer and a writer that
// refuses the diagram are both reported rather than causing a panic, and
// String returns the diagram without needing a writer at all.
package timeline

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// timelineLinesCap is the number of lines a small timeline starts with.
	timelineLinesCap int = 8
	// statementIndent is the indent of a statement directly under "timeline".
	statementIndent = "    "
	// sectionIndent is the indent of a period inside a section.
	sectionIndent = "        "
)

// Diagram is a timeline diagram builder.
type Diagram struct {
	// body is the timeline diagram body.
	body []string
	// dest is the output destination for the timeline diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the diagram building.
	err error
	// inSection reports whether a section has been opened, which decides how
	// far a period is indented.
	inSection bool
	// hasPeriod reports whether a period has been written, which is what an
	// event needs in order to have somewhere to go.
	hasPeriod bool
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()
	for _, opt := range opts {
		opt(c)
	}

	trimmedTitle := strings.TrimSpace(c.title)
	if containsNewline(trimmedTitle) {
		return &Diagram{
			body: []string{"timeline"},
			dest: w,
			err:  errors.New("title must not contain newline characters"),
		}
	}

	lines := make([]string, 0, timelineLinesCap)
	lines = append(lines, "timeline")
	if trimmedTitle != noTitle {
		// A title statement reads to the end of the line, so almost every
		// character goes in as it is; the newline that would end it early is
		// refused above, and a "<" is escaped because the renderer's sanitizer
		// otherwise draws it as "&lt;".
		lines = append(lines, statementIndent+"title "+internal.EscapeTitleAngle(trimmedTitle))
	}

	return &Diagram{
		body: lines,
		dest: w,
	}
}

// String returns the timeline diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the timeline diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the timeline diagram body to the output destination.
func (d *Diagram) Build() error {
	if d.err != nil {
		return d.err
	}
	if d.dest == nil {
		d.err = errors.New("output writer must not be nil")
		return d.err
	}

	if _, err := fmt.Fprint(d.dest, d.String()); err != nil {
		d.err = fmt.Errorf("failed to write: %w", err)
		return d.err
	}
	return nil
}

// LF adds a line feed to the timeline diagram body.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}
	d.body = append(d.body, "")
	return d
}

// Section starts a section, which groups the periods that follow it.
//
// A timeline needs no sections at all; periods written before the first one
// stand on their own.
func (d *Diagram) Section(name string) *Diagram {
	if d.err != nil {
		return d
	}

	text, err := validateText("section name", name)
	if err != nil {
		d.setError(err)
		return d
	}

	d.inSection = true
	d.hasPeriod = false
	d.body = append(d.body, statementIndent+"section "+escapeSection(text))
	return d
}

// Period adds a time period and the events that belong to it.
//
// The period is what a reader sees on the axis, so it is usually a year or a
// date, but mermaid does not require that: any text will do. Events may also be
// added afterwards with Event.
func (d *Diagram) Period(period string, events ...string) *Diagram {
	if d.err != nil {
		return d
	}

	text, err := validateText("period", period)
	if err != nil {
		d.setError(err)
		return d
	}

	var line strings.Builder
	line.WriteString(d.periodIndent())
	line.WriteString(escapeText(text))
	for i, event := range events {
		eventText, err := validateText(fmt.Sprintf("event %d of period %q", i+1, period), event)
		if err != nil {
			d.setError(err)
			return d
		}
		line.WriteString(" : ")
		line.WriteString(escapeText(eventText))
	}

	d.hasPeriod = true
	d.body = append(d.body, line.String())
	return d
}

// Event adds an event to the period Period opened last.
//
// Period must be called before Event; otherwise an error is recorded, because
// an event with no period has nowhere to go on the axis.
func (d *Diagram) Event(event string) *Diagram {
	if d.err != nil {
		return d
	}
	if !d.hasPeriod {
		d.setError(fmt.Errorf("event %q requires a period; call Period first", event))
		return d
	}

	text, err := validateText("event", event)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body[len(d.body)-1] += " : " + escapeText(text)
	return d
}

// periodIndent returns the indent a period is written at, which is one level
// deeper inside a section.
func (d *Diagram) periodIndent() string {
	if d.inSection {
		return sectionIndent
	}
	return statementIndent
}

// setError records the first error and leaves any later one alone, because the
// first is the one that explains the rest.
func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

// validateText returns value trimmed and checked, ready for the escaping its
// field needs.
//
// A newline ends a statement, so it cannot be carried through as it is. It is
// rejected rather than encoded, because a timeline entry is a single line by
// construction and silently joining the lines would say something the caller
// did not. The punctuation a field cannot carry is escaped by the caller,
// because a section and a period lose a different set: see escape.go.
func validateText(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s must not be empty", fieldName)
	}
	if containsNewline(trimmed) {
		return "", fmt.Errorf("%s must not contain newline characters", fieldName)
	}
	return trimmed, nil
}

// containsNewline reports whether value holds a line ending of either kind.
func containsNewline(value string) bool {
	return strings.ContainsAny(value, "\n\r")
}
