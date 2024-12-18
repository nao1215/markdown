// Package arch is mermaid architecture diagram builder.
// The building blocks of an architecture are groups, services, edges, and junctions.
// The arch package incorporates beta features of Mermaid, so the specifications are subject to significant changes.
package arch

import (
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// Architecture is a architecture diagram builder.
type Architecture struct {
	// body is architecture diagram body.
	body []string
	// config is the configuration for the architecture diagram.
	config *config
	// dest is output destination for architecture diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the architecture building.
	err error
}

// NewArchitecture returns a new Architecture.
func NewArchitecture(w io.Writer, opts ...Option) *Architecture {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	return &Architecture{
		body:   []string{"architecture-beta"},
		dest:   w,
		config: c,
	}
}

// String returns the architecture diagram body.
func (a *Architecture) String() string {
	return strings.Join(a.body, internal.LineFeed())
}

// Build writes the architecture diagram body to the output destination.
func (a *Architecture) Build() error {
	if _, err := a.dest.Write([]byte(a.String())); err != nil {
		if a.err != nil {
			return fmt.Errorf("failed to write: %w: %s", err, a.err.Error()) //nolint:wrapcheck
		}
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

// Error returns the error that occurred during the architecture diagram building.
func (a *Architecture) Error() error {
	return a.err
}

// LF add a line feed to the architecture diagram body.
func (a *Architecture) LF() *Architecture {
	a.body = append(a.body, "")
	return a
}
