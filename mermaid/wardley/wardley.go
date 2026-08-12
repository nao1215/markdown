// Package wardley is mermaid Wardley map builder.
//
// Ref. https://mermaid.js.org/syntax/wardley.html
//
// A Wardley map places the parts of a system on two axes: how visible each part
// is to the user, and how evolved it is, from something built for the first
// time to something bought as a commodity. This package keeps the flat,
// chainable API the rest of the library has: a part is one call, a dependency
// between two is one call.
//
// mermaid still spells the keyword "wardley-beta", so its syntax may change.
// The name is kept out of this package's API so that the API does not have to
// change when mermaid settles it.
//
// Errors are recorded rather than returned from every call: the chain runs to
// the end and the error surfaces from Build. A nil writer and a writer that
// refuses the map are both reported rather than causing a panic, and String
// returns the map without needing a writer at all.
package wardley

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// wardleyLinesCap is the number of lines a small map starts with.
	wardleyLinesCap int = 8
	// header is the keyword the map starts with.
	header = "wardley-beta"
	// indentUnit is the indent every statement carries.
	indentUnit = "    "
)

// Map is a Wardley map builder.
type Map struct {
	// body is the Wardley map body.
	body []string
	// dest is the output destination for the Wardley map body.
	dest io.Writer
	// err manages errors that occur in all parts of the map building.
	err error
}

// NewMap returns a new Map.
func NewMap(w io.Writer, opts ...Option) *Map {
	c := newConfig()
	for _, opt := range opts {
		opt(c)
	}

	lines := make([]string, 0, wardleyLinesCap)
	lines = append(lines, header)

	m := &Map{
		body: lines,
		dest: w,
	}

	trimmedTitle := strings.TrimSpace(c.title)
	if trimmedTitle == noTitle {
		return m
	}
	if containsNewline(trimmedTitle) {
		// mermaid reads the rest of the line as the title, so a second line
		// would be read as a statement of its own.
		m.err = errors.New("title must not contain newline characters")
		return m
	}
	m.body = append(m.body, indentUnit+"title "+escapeTitle(trimmedTitle))
	return m
}

// String returns the Wardley map body.
func (m *Map) String() string {
	return strings.Join(m.body, internal.LineFeed())
}

// Error returns the error that occurred during the Wardley map building.
//
// It returns the error the chain recorded, for code that wants to look before
// writing anything. Build reports that error too when it stops it writing.
func (m *Map) Error() error {
	return m.err
}

// Build writes the Wardley map body to the output destination.
func (m *Map) Build() error {
	if m.err != nil {
		return m.err
	}
	if m.dest == nil {
		m.err = errors.New("output writer must not be nil")
		return m.err
	}

	if _, err := fmt.Fprint(m.dest, m.String()); err != nil {
		m.err = fmt.Errorf("failed to write: %w", err)
		return m.err
	}
	return nil
}

// LF adds a line feed to the Wardley map body.
func (m *Map) LF() *Map {
	if m.err != nil {
		return m
	}
	m.body = append(m.body, "")
	return m
}

// Component adds a part of the system at the given position.
//
// The two coordinates are how evolved the part is and how visible it is to the
// user, each from 0.0 to 1.0. Evolution runs left to right, from something
// built for the first time to something bought as a commodity; visibility runs
// bottom to top, from the plumbing to what the user actually touches.
func (m *Map) Component(name string, evolution, visibility float64) *Map {
	return m.place("component", name, evolution, visibility)
}

// Anchor adds the user the map is drawn from the point of view of. A map
// usually has one, at the top, and everything else hangs below it.
func (m *Map) Anchor(name string, evolution, visibility float64) *Map {
	return m.place("anchor", name, evolution, visibility)
}

// Link draws a dependency: the first named part needs the second.
func (m *Map) Link(from, to string) *Map {
	if m.err != nil {
		return m
	}

	validFrom, err := validateName(from)
	if err != nil {
		m.setError(err)
		return m
	}
	validTo, err := validateName(to)
	if err != nil {
		m.setError(err)
		return m
	}

	m.body = append(m.body, fmt.Sprintf("%s%s -> %s", indentUnit, validFrom, validTo))
	return m
}

// Evolve marks a part as moving along the evolution axis, drawn as an arrow
// from where it is to where it is going.
func (m *Map) Evolve(name string, evolution float64) *Map {
	if m.err != nil {
		return m
	}

	validName, err := validateName(name)
	if err != nil {
		m.setError(err)
		return m
	}
	if err := validatePosition("evolution of "+validName, evolution); err != nil {
		m.setError(err)
		return m
	}

	m.body = append(m.body, fmt.Sprintf("%sevolve %s %s", indentUnit, validName, position(evolution)))
	return m
}

// place writes one component or anchor.
func (m *Map) place(keyword, name string, evolution, visibility float64) *Map {
	if m.err != nil {
		return m
	}

	validName, err := validateName(name)
	if err != nil {
		m.setError(err)
		return m
	}
	if err := validatePosition("evolution of "+validName, evolution); err != nil {
		m.setError(err)
		return m
	}
	if err := validatePosition("visibility of "+validName, visibility); err != nil {
		m.setError(err)
		return m
	}

	m.body = append(m.body, fmt.Sprintf("%s%s %s [%s, %s]",
		indentUnit, keyword, validName, position(evolution), position(visibility)))
	return m
}

// setError records the first error and leaves any later one alone, because the
// first is the one that explains the rest.
func (m *Map) setError(err error) {
	if m.err == nil {
		m.err = err
	}
}

// validateName returns value ready to be written as a part's bare name.
//
// A name is written unquoted, and mermaid's lexer takes letters, digits,
// spaces, underscores, hyphens and parentheses there and nothing else: a full
// stop, a slash, a quotation mark, an emoji and Japanese text each lose the
// whole map, which was found by rendering them one at a time. It refuses its
// own "#name;" escape in that position too, so there is nothing to encode to
// and a name outside the set is reported rather than mangled into one that
// draws something else.
//
// An ampersand is left out although "a&b" renders: "a & b" does not, so it
// works only while nobody puts a space around it. A character that depends on
// its neighbors is worse to hand a caller than one that is simply refused.
func validateName(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", errors.New("name must not be empty")
	}
	for _, r := range trimmed {
		if !isNameRune(r) {
			return "", fmt.Errorf(
				"name %q must hold only letters, digits, spaces, underscores, hyphens and parentheses; mermaid reads nothing else there",
				trimmed,
			)
		}
	}
	return trimmed, nil
}

// isNameRune reports whether r may appear in a part's name.
func isNameRune(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
		r == ' ' || r == '_' || r == '-' || r == '(' || r == ')'
}

// validatePosition rejects the coordinates a map has no way to draw.
func validatePosition(what string, value float64) error {
	switch {
	case math.IsNaN(value):
		return fmt.Errorf("%s must be a number", what)
	case math.IsInf(value, 0):
		return fmt.Errorf("%s must be finite", what)
	case value < 0 || value > 1:
		return fmt.Errorf("%s must be between 0.0 and 1.0, not %s", what, position(value))
	}
	return nil
}

// position renders a coordinate the way mermaid expects it.
//
// Plain decimal notation, never the exponent form Go reaches for at the small
// end: mermaid parses the token as a number and "1e-07" is not one to it.
func position(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

// containsNewline reports whether value holds a line ending of either kind.
func containsNewline(value string) bool {
	return strings.ContainsAny(value, "\n\r")
}
