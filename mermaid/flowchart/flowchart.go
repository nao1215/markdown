// Package flowchart provides a simple way to create flowcharts in mermaid syntax.
//
// Errors are recorded rather than returned from every call: the chain runs to
// the end and the error surfaces from Build. A nil writer and a writer that
// refuses the diagram are both reported rather than causing a panic, and
// String returns the diagram without needing a writer at all.
package flowchart

import (
	"errors"
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
	// depth is how many subgraphs are open, counted in indentUnits.
	depth int
}

// indentUnit is one level of subgraph nesting.
const indentUnit = "    "

// NewFlowchart returns a new Flowchart.
func NewFlowchart(w io.Writer, opts ...Option) *Flowchart {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := []string{}
	if strings.TrimSpace(c.title) != noTitle {
		lines = append(lines, "---")
		lines = append(lines, internal.FrontMatterTitle(internal.FoldFrontMatterTitleCR(c.title)))
		lines = append(lines, "---")
	}
	lines = append(lines, fmt.Sprintf("flowchart %s", c.oriental.string()))

	return &Flowchart{
		body:   lines,
		dest:   w,
		config: c,
	}
}

// indent returns the prefix of a line at the current nesting.
//
// A flowchart with no subgraph in it indents by one unit, which is what every
// line of every flowchart written before subgraphs existed used, so the output
// of a chain that opens none is unchanged.
func (f *Flowchart) indent() string {
	return indentUnit + strings.Repeat(indentUnit, f.depth)
}

// setError records the first error and leaves any later one alone, because the
// first is the one that explains the rest.
func (f *Flowchart) setError(err error) {
	if f.err == nil {
		f.err = err
	}
}

// String returns the flowchart body.
func (f *Flowchart) String() string {
	return strings.Join(f.body, internal.LineFeed())
}

// Error returns the error that occurred during the flowchart building.
//
// It returns the error the chain recorded, for code that wants to look before
// writing anything. Build reports that error too when it stops it writing, but
// the two are not the same call: Build returns nil once it has written the
// document, whatever was recorded on the way.
//
// Every other builder in this library has had this since it was written; this
// one gained it at v1.0.0, when the API audit noticed it missing.
func (f *Flowchart) Error() error {
	return f.err
}

// Build writes the flowchart body to the output destination.
//
// A subgraph left open is reported here rather than written out: mermaid
// refuses a flowchart whose subgraph never ends, which loses the whole drawing.
// So is any error the chain recorded, which is what every other builder in this
// library does.
func (f *Flowchart) Build() error {
	if f.depth != 0 {
		f.setError(fmt.Errorf("%d subgraph must be closed with SubgraphEnd before Build", f.depth))
	}
	// A recorded error stops the write, the way it does in every other builder
	// here. Nothing recorded one before subgraphs existed, so no chain written
	// against an earlier release reaches this.
	if f.err != nil {
		return f.err
	}
	if f.dest == nil {
		if f.err == nil {
			f.err = errors.New("output writer must not be nil")
		}
		return f.err
	}

	if _, err := fmt.Fprint(f.dest, f.String()); err != nil {
		if f.err != nil {
			return fmt.Errorf("failed to write: %w: %s", err, f.err.Error()) //nolint:wrapcheck
		}
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}
