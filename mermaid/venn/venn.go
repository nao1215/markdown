// Package venn is mermaid Venn diagram builder.
//
// Ref. https://mermaid.js.org/syntax/venn.html
//
// A Venn diagram is a handful of sets and what they have in common. This
// package keeps the flat, chainable API the rest of the library has: a set is
// one call, and where sets overlap is mermaid's business rather than something
// the caller lays out.
//
// mermaid still spells the keyword "venn-beta", so its syntax may change. The
// name is kept out of this package's API so that the API does not have to
// change when mermaid settles it.
//
// Errors are recorded rather than returned from every call: the chain runs to
// the end and the error surfaces from Build. A nil writer and a writer that
// refuses the diagram are both reported rather than causing a panic, and
// String returns the diagram without needing a writer at all.
package venn

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// vennLinesCap is the number of lines a small diagram starts with.
	vennLinesCap int = 8
	// header is the keyword the diagram starts with.
	header = "venn-beta"
	// indentUnit is the indent every statement carries.
	indentUnit = "    "
)

// Diagram is a Venn diagram builder.
type Diagram struct {
	// body is the Venn diagram body.
	body []string
	// dest is the output destination for the Venn diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the diagram building.
	err error
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()
	for _, opt := range opts {
		opt(c)
	}

	lines := make([]string, 0, vennLinesCap)
	lines = append(lines, header)

	d := &Diagram{
		body: lines,
		dest: w,
	}

	trimmedTitle := strings.TrimSpace(c.title)
	if trimmedTitle == noTitle {
		return d
	}
	if containsNewline(trimmedTitle) {
		// mermaid reads the rest of the line as the title, so a second line
		// would be read as a statement of its own.
		d.err = errors.New("title must not contain newline characters")
		return d
	}
	d.body = append(d.body, indentUnit+"title "+escapeTitle(trimmedTitle))
	return d
}

// String returns the Venn diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the Venn diagram building.
//
// It returns the error the chain recorded, for code that wants to look before
// writing anything. Build reports that error too when it stops it writing.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the Venn diagram body to the output destination.
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

// LF adds a line feed to the Venn diagram body.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}
	d.body = append(d.body, "")
	return d
}

// Set adds a set drawn with its identifier inside.
//
// Where two sets overlap is decided by mermaid from the sets it is given; there
// is nothing to declare about an intersection, which is why this package has no
// call for one.
func (d *Diagram) Set(id string) *Diagram {
	if d.err != nil {
		return d
	}

	validID, err := validateIdentifier(id)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, fmt.Sprintf("%sset %s", indentUnit, validID))
	return d
}

// SetWithLabel adds a set drawn with the given label rather than its
// identifier, for a set whose name is not what a reader should see.
func (d *Diagram) SetWithLabel(id, label string) *Diagram {
	if d.err != nil {
		return d
	}

	validID, err := validateIdentifier(id)
	if err != nil {
		d.setError(err)
		return d
	}
	text, err := validateText("set label", label)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, fmt.Sprintf(`%sset %s["%s"]`, indentUnit, validID, escapeLabel(text)))
	return d
}

// setError records the first error and leaves any later one alone, because the
// first is the one that explains the rest.
func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

// validateText returns value ready to be quoted into a label.
func validateText(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s must not be empty", fieldName)
	}
	if containsNewline(trimmed) {
		// One set is one line, so a label spanning lines would be read as a
		// statement of its own. "<br/>" is how a label breaks a line.
		return "", fmt.Errorf("%s must not contain newline characters", fieldName)
	}
	return trimmed, nil
}

// validateIdentifier returns value ready to be written as a set's bare name.
//
// A set name is written unquoted, and mermaid's lexer takes only word
// characters and a hyphen there: a full stop, a space and a comma each lose the
// whole diagram, which was found by rendering them one at a time. There is
// nothing to escape to, so a name outside that set is reported rather than
// mangled into one that draws something else.
func validateIdentifier(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", errors.New("set name must not be empty")
	}
	for _, r := range trimmed {
		if !isNameRune(r) {
			return "", fmt.Errorf(
				"set name %q must hold only letters, digits, underscores and hyphens; mermaid reads nothing else there",
				trimmed,
			)
		}
	}
	return trimmed, nil
}

// isNameRune reports whether r may appear in a set name.
func isNameRune(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
		r == '_' || r == '-'
}

// containsNewline reports whether value holds a line ending of either kind.
func containsNewline(value string) bool {
	return strings.ContainsAny(value, "\n\r")
}
