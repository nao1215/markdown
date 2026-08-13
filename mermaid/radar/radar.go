// Package radar is mermaid radar chart builder.
//
// Ref. https://mermaid.js.org/syntax/radar.html
//
// A radar chart is a set of axes and a set of curves, where each curve gives
// one value per axis. Axes are declared once, in order, and every curve lists
// its values in that same order.
//
// Errors are recorded rather than returned from every call: the chain runs to
// the end and the error surfaces from Build. A nil writer and a writer that
// refuses the diagram are both reported rather than causing a panic, and
// String returns the diagram without needing a writer at all.
package radar

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
	// radarLinesCap is the number of lines a small chart starts with.
	radarLinesCap int = 8
	// header is the keyword the chart starts with.
	//
	// mermaid still spells it "beta". The name is kept out of this package's
	// API so that the API does not have to change when mermaid settles it.
	header = "radar-beta"
	// statementIndent is the indent of a statement under the header.
	statementIndent = "  "
)

// Diagram is a radar chart builder.
type Diagram struct {
	// body is the radar chart body.
	body []string
	// dest is the output destination for the radar chart body.
	dest io.Writer
	// err manages errors that occur in all parts of the chart building.
	err error
	// axes counts the axes declared so far, which is what an identifier is
	// numbered from.
	axes int
	// curves counts the curves declared so far, for the same reason.
	curves int
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

	lines := make([]string, 0, radarLinesCap)
	if trimmedTitle != noTitle {
		lines = append(lines, "---", internal.FrontMatterTitle(internal.FoldFrontMatterTitleCR(trimmedTitle)), "---")
	}
	lines = append(lines, header)

	return &Diagram{
		body: lines,
		dest: w,
	}
}

// String returns the radar chart body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the radar chart building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the radar chart body to the output destination.
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

// LF adds a line feed to the radar chart body.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}
	d.body = append(d.body, "")
	return d
}

// Axis adds one or more axes, in the order a curve's values are read in.
//
// mermaid wants an identifier in front of every axis label. Nothing in a radar
// chart refers to one, so this package numbers them and the caller passes only
// the labels.
func (d *Diagram) Axis(labels ...string) *Diagram {
	if d.err != nil {
		return d
	}
	if len(labels) == 0 {
		d.setError(errors.New("axis requires at least one label"))
		return d
	}

	entries := make([]string, 0, len(labels))
	for i, label := range labels {
		text, err := validateText(fmt.Sprintf("axis label %d", i+1), label)
		if err != nil {
			d.setError(err)
			return d
		}
		d.axes++
		entries = append(entries, fmt.Sprintf("a%d[%s]", d.axes, quote(text)))
	}

	d.body = append(d.body, statementIndent+"axis "+strings.Join(entries, ", "))
	return d
}

// Curve adds a curve, giving one value per axis in the order the axes were
// declared.
//
// mermaid draws whatever it is given, so a curve with fewer values than there
// are axes is not refused here: it simply does not reach the axes it has no
// value for.
func (d *Diagram) Curve(label string, values ...float64) *Diagram {
	if d.err != nil {
		return d
	}

	text, err := validateText("curve label", label)
	if err != nil {
		d.setError(err)
		return d
	}
	if len(values) == 0 {
		d.setError(fmt.Errorf("curve %q requires at least one value", label))
		return d
	}

	formatted := make([]string, 0, len(values))
	for i, value := range values {
		if err := validateValue(fmt.Sprintf("value %d of curve %q", i+1, label), value); err != nil {
			d.setError(err)
			return d
		}
		formatted = append(formatted, formatValue(value))
	}

	d.curves++
	d.body = append(d.body, fmt.Sprintf("%scurve c%d[%s]{%s}",
		statementIndent, d.curves, quote(text), strings.Join(formatted, ", ")))
	return d
}

// Max sets the outer edge of the chart, the value an axis is full at.
func (d *Diagram) Max(value float64) *Diagram {
	return d.scale("max", value)
}

// Min sets the centre of the chart, the value an axis is empty at.
func (d *Diagram) Min(value float64) *Diagram {
	return d.scale("min", value)
}

// scale writes one of the two scale statements.
func (d *Diagram) scale(keyword string, value float64) *Diagram {
	if d.err != nil {
		return d
	}
	if err := validateValue(keyword, value); err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, fmt.Sprintf("%s%s %s", statementIndent, keyword, formatValue(value)))
	return d
}

// setError records the first error and leaves any later one alone, because the
// first is the one that explains the rest.
func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

// validateText returns value ready to be quoted into the chart.
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

// quote returns value as a quoted label.
//
// A double quote ends the label, and a backslash before the closing one would
// swallow it, so both are escaped. Everything else, a colon and a hash
// included, is text once it is inside the quotes.
func quote(value string) string {
	escaped := strings.ReplaceAll(value, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	return `"` + escaped + `"`
}

// validateValue rejects the numbers a chart has no way to draw.
func validateValue(fieldName string, value float64) error {
	switch {
	case math.IsNaN(value):
		return fmt.Errorf("%s must be a number", fieldName)
	case math.IsInf(value, 0):
		return fmt.Errorf("%s must be finite", fieldName)
	}
	return nil
}

// formatValue renders a number the way mermaid expects it.
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
