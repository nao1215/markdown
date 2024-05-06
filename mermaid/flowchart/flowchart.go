// Package flowchart provides a simple way to create flowcharts in mermaid syntax.
package flowchart

import (
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// Flowchart is a flowchart builder.
type Flowchart struct {
	// body is flowchart body.
	body []string
	// dest is output destination for flowchart body.
	dest io.Writer
	// err manages errors that occur in all parts of the flowchart building.
	err error
	// config is the configuration for the flowchart.
	config *config
}

// NewFlowchart returns a new Flowchart.
func NewFlowchart(w io.Writer, opts ...Option) *Flowchart {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := []string{}
	if strings.TrimSpace(c.title) != noTitle {
		lines = append(lines, "---")
		lines = append(lines, fmt.Sprintf("title: %s", c.title))
		lines = append(lines, "---")
	}
	lines = append(lines, fmt.Sprintf("flowchart %s", c.oriental.string()))

	return &Flowchart{
		body:   lines,
		dest:   w,
		config: c,
	}
}

// String returns the flowchart body.
func (f *Flowchart) String() string {
	return strings.Join(f.body, internal.LineFeed())
}

// Build writes the flowchart body to the output destination.
func (f *Flowchart) Build() error {
	if _, err := fmt.Fprint(f.dest, f.String()); err != nil {
		if f.err != nil {
			return fmt.Errorf("failed to write: %w: %s", err, f.err.Error()) //nolint:wrapcheck
		}
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}
