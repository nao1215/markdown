// Package sankey is mermaid sankey diagram builder.
//
// Ref. https://mermaid.js.org/syntax/sankey.html
//
// A sankey diagram is a set of flows, each running from one node to another
// with a quantity attached. The nodes are not declared: they are whatever the
// flows name, so a diagram is built entirely out of Link calls.
//
// The diagram body is CSV, which is why this package quotes rather than
// escapes. A node name holding a comma would otherwise end the field and shift
// the quantity into the target column, producing a diagram that parses and is
// wrong.
//
// Errors are recorded rather than returned from every call: the chain runs to
// the end and the error surfaces from Build. A nil writer and a writer that
// refuses the diagram are both reported rather than causing a panic, and
// String returns the diagram without needing a writer at all.
package sankey

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
	// sankeyLinesCap is the number of lines a small diagram starts with.
	sankeyLinesCap int = 8
	// header is the keyword the diagram starts with.
	//
	// mermaid still spells it "beta". The name is kept out of this package's
	// API so that the API does not have to change when mermaid settles it.
	header = "sankey-beta"
)

// Diagram is a sankey diagram builder.
type Diagram struct {
	// body is the sankey diagram body.
	body []string
	// dest is the output destination for the sankey diagram body.
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

	trimmedTitle := strings.TrimSpace(c.title)
	if containsNewline(trimmedTitle) {
		return &Diagram{
			body: []string{header, ""},
			dest: w,
			err:  errors.New("title must not contain newline characters"),
		}
	}

	lines := make([]string, 0, sankeyLinesCap)
	if trimmedTitle != noTitle {
		lines = append(lines, "---", internal.FrontMatterTitle(trimmedTitle), "---")
	}
	// The blank line after the keyword is what mermaid's own examples use, and
	// it keeps the CSV a block of its own rather than a continuation of the
	// header line.
	lines = append(lines, header, "")

	return &Diagram{
		body: lines,
		dest: w,
	}
}

// String returns the sankey diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the sankey diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the sankey diagram body to the output destination.
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

// LF adds a line feed to the sankey diagram body.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}
	d.body = append(d.body, "")
	return d
}

// Link adds a flow of the given quantity from source to target.
//
// Nodes are never declared: a node exists because a flow names it, and two
// flows naming the same node are two flows through one node.
func (d *Diagram) Link(source, target string, value float64) *Diagram {
	if d.err != nil {
		return d
	}

	sourceField, err := field("source", source)
	if err != nil {
		d.setError(err)
		return d
	}
	targetField, err := field("target", target)
	if err != nil {
		d.setError(err)
		return d
	}
	if err := validateValue(source, target, value); err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, fmt.Sprintf("%s,%s,%s", sourceField, targetField, formatValue(value)))
	return d
}

// setError records the first error and leaves any later one alone, because the
// first is the one that explains the rest.
func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

// field returns value as one CSV field.
//
// A comma in a node name would end the field and shift every column after it,
// so a name that holds one is quoted. Inside quotes a double quote stands for
// itself only when it is doubled, which is what CSV says and what mermaid's
// parser implements.
func field(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s must not be empty", fieldName)
	}
	if containsNewline(trimmed) {
		// CSV allows a newline inside a quoted field, but a sankey row is one
		// flow and mermaid reads the diagram line by line, so a name spanning
		// lines has nowhere to go.
		return "", fmt.Errorf("%s must not contain newline characters", fieldName)
	}

	if !strings.ContainsAny(trimmed, `,"`) {
		return trimmed, nil
	}
	return `"` + strings.ReplaceAll(trimmed, `"`, `""`) + `"`, nil
}

// validateValue rejects the quantities a sankey diagram has no way to draw.
//
// A flow has a width, so a negative one has no meaning, and neither NaN nor an
// infinity survives being written as a decimal number.
func validateValue(source, target string, value float64) error {
	switch {
	case math.IsNaN(value):
		return fmt.Errorf("value of the flow from %q to %q must be a number", source, target)
	case math.IsInf(value, 0):
		return fmt.Errorf("value of the flow from %q to %q must be finite", source, target)
	case value < 0:
		return fmt.Errorf("value of the flow from %q to %q must not be negative", source, target)
	}
	return nil
}

// formatValue renders a quantity the way mermaid's CSV reader expects it.
//
// Plain decimal notation, never the exponent form Go reaches for above 1e21:
// mermaid parses the field as a number and "1e+21" is not one to it.
func formatValue(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

// containsNewline reports whether value holds a line ending of either kind.
func containsNewline(value string) bool {
	return strings.ContainsAny(value, "\n\r")
}
