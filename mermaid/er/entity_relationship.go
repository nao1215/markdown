// Package er is mermaid entity relationship diagram builder.
//
// Errors are recorded rather than returned from every call: the chain runs to
// the end and the error surfaces from Build. A nil writer and a writer that
// refuses the diagram are both reported rather than causing a panic, and
// String returns the diagram without needing a writer at all.
package er

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/nao1215/markdown/internal"
)

// Diagram is a entity relationship diagram builder.
type Diagram struct {
	// body is entity relationship diagram body.
	body []string
	// config is the configuration for the entity relationship diagram.
	config *config
	// dest is output destination for entity relationship diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the entity relationship building.
	err error
	// entities is the list of entities in the diagram.
	entities sync.Map
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	return &Diagram{
		body:     []string{"erDiagram"},
		dest:     w,
		config:   c,
		entities: sync.Map{},
	}
}

// String returns the entity relationship diagram body.
func (d *Diagram) String() string {
	s := strings.Join(d.body, internal.LineFeed())
	s += internal.LineFeed()

	entities := make([]Entity, 0)
	d.entities.Range(func(_, value interface{}) bool {
		e, ok := value.(Entity)
		if !ok {
			return false
		}
		entities = append(entities, e)
		return true
	})

	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Name < entities[j].Name
	})

	for _, e := range entities {
		s += e.string() + internal.LineFeed()
	}
	return s
}

// Build writes the entity relationship body to the output destination.
func (d *Diagram) Build() error {
	if d.dest == nil {
		if d.err == nil {
			d.err = errors.New("output writer must not be nil")
		}
		return d.err
	}

	if _, err := fmt.Fprint(d.dest, d.String()); err != nil {
		if d.err != nil {
			return fmt.Errorf("failed to write: %w: %s", err, d.err.Error()) //nolint:wrapcheck
		}
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}
