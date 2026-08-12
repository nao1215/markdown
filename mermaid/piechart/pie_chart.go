// Package piechart is mermaid pie chart builder.
//
// Errors are recorded rather than returned from every call: the chain runs to
// the end and the error surfaces from Build. A nil writer and a writer that
// refuses the diagram are both reported rather than causing a panic, and
// String returns the diagram without needing a writer at all.
package piechart

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// PieChart is a pie chart builder.
type PieChart struct {
	// body is pie chart body.
	body []string
	// dest is output destination for pie chart body.
	dest io.Writer
	// err manages errors that occur in all parts of the pie chart building.
	err error
	// config is the configuration for the pie chart.
	config *config
}

// NewPieChart returns a new PieChart.
func NewPieChart(w io.Writer, opts ...Option) *PieChart {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := []string{}
	lines = append(
		lines,
		fmt.Sprintf(
			"%%%%{init: {\"pie\": {\"textPosition\": %.2f}, \"themeVariables\": {\"pieOuterStrokeWidth\": \"5px\"}} }%%%%",
			c.textPosition,
		))

	baseLine := "pie"
	if c.showData {
		baseLine += " showData"
	}

	if c.title == noTitle {
		lines = append(lines, baseLine)
	} else {
		lines = append(lines, baseLine)
		lines = append(lines, fmt.Sprintf("    title %s", escapeTitle(c.title)))
	}

	return &PieChart{
		body:   lines,
		dest:   w,
		config: c,
	}
}

// String returns the pie chart body.
func (p *PieChart) String() string {
	return strings.Join(p.body, internal.LineFeed())
}

// Error returns the error that occurred during the pie chart building.
//
// It returns the error the chain recorded, for code that wants to look before
// writing anything. Build reports that error too when it stops it writing, but
// the two are not the same call: Build returns nil once it has written the
// document, whatever was recorded on the way.
//
// Every other builder in this library has had this since it was written; this
// one gained it at v1.0.0, when the API audit noticed it missing.
func (p *PieChart) Error() error {
	return p.err
}

// Build writes the pie chart body to the output destination.
func (p *PieChart) Build() error {
	if p.dest == nil {
		if p.err == nil {
			p.err = errors.New("output writer must not be nil")
		}
		return p.err
	}

	if _, err := fmt.Fprint(p.dest, p.String()); err != nil {
		if p.err != nil {
			return fmt.Errorf("failed to write: %w: %s", err, p.err.Error()) //nolint:wrapcheck
		}
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

// LabelAndIntValue adds a label and value to the pie chart.
func (p *PieChart) LabelAndIntValue(label string, value uint64) *PieChart {
	p.body = append(p.body, fmt.Sprintf("    \"%s\" : %d", escapeLabel(label), value))
	return p
}

// LabelAndFloatValue adds a label and value to the pie chart.
// The value is formatted with a precision of 6 digits after the decimal point.
func (p *PieChart) LabelAndFloatValue(label string, value float64) *PieChart {
	p.body = append(p.body, fmt.Sprintf("    \"%s\" : %f", escapeLabel(label), value))
	return p
}
