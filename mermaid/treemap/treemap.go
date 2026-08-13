// Package treemap is mermaid treemap diagram builder.
//
// Ref. https://mermaid.js.org/syntax/treemap.html
//
// A treemap is a hierarchy of sections holding leaves, and mermaid expresses
// the hierarchy with indentation. This package keeps the flat, chainable API
// the rest of the library has rather than asking for a tree of node objects:
// Section opens a level, Leaf puts a value in the current one, and Parent goes
// back up, the same way mermaid/mindmap walks its own indentation.
//
// Errors are recorded rather than returned from every call: the chain runs to
// the end and the error surfaces from Build. A nil writer and a writer that
// refuses the diagram are both reported rather than causing a panic, and
// String returns the diagram without needing a writer at all.
package treemap

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
	// treemapLinesCap is the number of lines a small diagram starts with.
	treemapLinesCap int = 8
	// header is the keyword the diagram starts with.
	//
	// mermaid still spells it "beta". The name is kept out of this package's
	// API so that the API does not have to change when mermaid settles it.
	header = "treemap-beta"
	// indentUnit is one level of the hierarchy.
	indentUnit = "    "
)

// Diagram is a treemap diagram builder.
type Diagram struct {
	// body is the treemap diagram body.
	body []string
	// dest is the output destination for the treemap diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the diagram building.
	err error
	// depth is how far in the current level is, counted in indentUnits.
	depth int
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
			body: []string{header},
			dest: w,
			err:  errors.New("title must not contain newline characters"),
		}
	}

	lines := make([]string, 0, treemapLinesCap)
	if trimmedTitle != noTitle {
		lines = append(lines, "---", internal.FrontMatterTitle(internal.FoldFrontMatterTitleCR(trimmedTitle)), "---")
	}
	lines = append(lines, header)

	return &Diagram{
		body: lines,
		dest: w,
	}
}

// String returns the treemap diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the treemap diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the treemap diagram body to the output destination.
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

// LF adds a line feed to the treemap diagram body.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}
	d.body = append(d.body, "")
	return d
}

// Section adds a section at the current level and descends into it, so the
// calls that follow belong to it.
//
// A section carries no value of its own: mermaid gives it the sum of what it
// holds.
func (d *Diagram) Section(name string) *Diagram {
	if d.err != nil {
		return d
	}

	text, err := validateText("section name", name)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, d.indent()+quote(text))
	d.depth++
	return d
}

// Leaf adds a leaf with its value at the current level.
//
// A leaf is what a treemap actually draws: the area it gets is its value
// against the whole.
func (d *Diagram) Leaf(name string, value float64) *Diagram {
	if d.err != nil {
		return d
	}

	text, err := validateText("leaf name", name)
	if err != nil {
		d.setError(err)
		return d
	}
	if err := validateValue(name, value); err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, fmt.Sprintf("%s%s: %s", d.indent(), quote(text), formatValue(value)))
	return d
}

// Parent leaves the section opened last, so the calls that follow belong to the
// level above.
//
// Calling it at the top level is an error rather than a silent no-op: a chain
// that has gone up more often than down is not the document its author meant.
func (d *Diagram) Parent() *Diagram {
	if d.err != nil {
		return d
	}
	if d.depth == 0 {
		d.setError(errors.New("Parent was called at the top level; there is nothing to go up to"))
		return d
	}

	d.depth--
	return d
}

// indent returns the prefix of a line at the current level.
func (d *Diagram) indent() string {
	return strings.Repeat(indentUnit, d.depth)
}

// setError records the first error and leaves any later one alone, because the
// first is the one that explains the rest.
func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

// validateText returns value ready to be quoted into the diagram.
func validateText(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s must not be empty", fieldName)
	}
	if containsNewline(trimmed) {
		// The hierarchy is indentation, so a name spanning lines would be read
		// as another node at another level.
		return "", fmt.Errorf("%s must not contain newline characters", fieldName)
	}
	return trimmed, nil
}

// quote returns value as a quoted name.
//
// A double quote ends the name, and it is escaped by doubling rather than with
// a backslash: mermaid's treemap parser reads a backslash as an ordinary
// character and refuses the diagram when one appears before a quote. Doubling
// is what it does implement, which was found by rendering both. Everything
// else, a colon and a backslash included, is text once inside the quotes.
func quote(value string) string {
	return `"` + strings.ReplaceAll(value, `"`, `""`) + `"`
}

// validateValue rejects the numbers a treemap has no way to draw.
//
// An area cannot be negative, and neither NaN nor an infinity survives being
// written as a decimal number.
func validateValue(name string, value float64) error {
	switch {
	case math.IsNaN(value):
		return fmt.Errorf("value of leaf %q must be a number", name)
	case math.IsInf(value, 0):
		return fmt.Errorf("value of leaf %q must be finite", name)
	case value < 0:
		return fmt.Errorf("value of leaf %q must not be negative", name)
	}
	return nil
}

// formatValue renders a value the way mermaid expects it.
//
// Plain decimal notation, never the exponent form Go reaches for above 1e21:
// mermaid parses the token as a number and "1e+21" is not one to it.
func formatValue(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

// containsNewline reports whether value holds a line ending of either kind.
func containsNewline(value string) bool {
	return strings.ContainsAny(value, "\n\r")
}
