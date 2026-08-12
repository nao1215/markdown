// Package arch is mermaid architecture diagram builder.
// The building blocks of an architecture are groups, services, edges, and junctions.
// The arch package incorporates beta features of Mermaid, so the specifications are subject to significant changes.
//
// # Labels take letters, digits, underscores and spaces, and nothing else
//
// A group title and a service title are written between square brackets, and
// mermaid's architecture-beta grammar accepts only [A-Za-z0-9_ ] there. Every
// other character makes it refuse the whole diagram, so the reader gets an error
// box instead of a picture. That includes a quotation mark, an apostrophe, a
// hyphen, a full stop, a comma, a colon, an emoji and any non-ASCII text: a
// service titled "Order-Service" or "注文サービス" does not draw.
//
// Unlike every other diagram type in this library, there is nothing to escape
// to. The "#name;" entity form mermaid decodes elsewhere is refused by this
// lexer before any decoding happens, so "#quot;" fails exactly as a quotation
// mark does. This package therefore passes a title through as it was given
// rather than mangling it into something that renders but says something else,
// and the limit is recorded in SPEC.md. It is a limit of mermaid's beta
// grammar, and it will lift when that grammar does.
//
// Errors are recorded rather than returned from every call: the chain runs to
// the end and the error surfaces from Build. A nil writer and a writer that
// refuses the diagram are both reported rather than causing a panic, and
// String returns the diagram without needing a writer at all.
package arch

import (
	"errors"
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
	if a.dest == nil {
		if a.err == nil {
			a.err = errors.New("output writer must not be nil")
		}
		return a.err
	}

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
