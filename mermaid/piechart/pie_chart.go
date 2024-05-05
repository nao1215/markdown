// Package piechart is mermaid pie chart builder.
package piechart

import (
	"fmt"
	"io"
	"runtime"
	"strings"
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
		lines = append(lines, fmt.Sprintf("    title %s", c.title))
	}

	return &PieChart{
		body:   lines,
		dest:   w,
		config: c,
	}
}

// String returns the pie chart body.
func (p *PieChart) String() string {
	return strings.Join(p.body, lineFeed())
}

// Build writes the pie chart body to the output destination.
func (p *PieChart) Build() error {
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
	p.body = append(p.body, fmt.Sprintf("    \"%s\" : %d", label, value))
	return p
}

// LabelAndFloatValue adds a label and value to the pie chart.
// The value is formatted with a precision of 6 digits after the decimal point.
func (p *PieChart) LabelAndFloatValue(label string, value float64) *PieChart {
	p.body = append(p.body, fmt.Sprintf("    \"%s\" : %f", label, value))
	return p
}

// lineFeed return line feed for current OS.
func lineFeed() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
