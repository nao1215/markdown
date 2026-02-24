// Package mindmap is mermaid mindmap diagram builder.
//
// Ref. https://mermaid.js.org/syntax/mindmap.html
package mindmap

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// mindmapLinesCap is the max initial lines with title frontmatter.
	mindmapLinesCap int = 4
)

// Diagram is a mindmap diagram builder.
type Diagram struct {
	// body is mindmap diagram body.
	body []string
	// dest is output destination for mindmap diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the mindmap diagram building.
	err error
	// currentLevel is the latest node level.
	currentLevel int
	// hasRoot indicates whether the root node is already defined.
	hasRoot bool
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := make([]string, 0, mindmapLinesCap)
	title := strings.TrimSpace(c.title)
	if title != noTitle {
		if containsNewline(title) {
			return &Diagram{
				body: []string{"mindmap"},
				dest: w,
				err:  errors.New("title must not contain newline characters"),
			}
		}
		lines = append(lines, "---")
		lines = append(lines, fmt.Sprintf("title: %s", title))
		lines = append(lines, "---")
	}
	lines = append(lines, "mindmap")

	return &Diagram{
		body:         lines,
		dest:         w,
		currentLevel: -1,
	}
}

// String returns the mindmap diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the mindmap diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the mindmap diagram body to the output destination.
func (d *Diagram) Build() error {
	if d.err != nil {
		return d.err
	}

	if _, err := fmt.Fprint(d.dest, d.String()); err != nil {
		d.err = fmt.Errorf("failed to write: %w", err)
		return d.err
	}
	return nil
}

// Root adds the root node to the mindmap.
func (d *Diagram) Root(text string) *Diagram {
	return d.Node(0, text)
}

// Node adds a node at the specified level.
//
// The root level is 0. A root node must be defined first.
func (d *Diagram) Node(level int, text string) *Diagram {
	if d.err != nil {
		return d
	}

	if level < 0 {
		d.setError(errors.New("level must be greater than or equal to zero"))
		return d
	}

	trimmedText, err := validateName("node text", text)
	if err != nil {
		d.setError(err)
		return d
	}

	if !d.hasRoot {
		if level != 0 {
			d.setError(errors.New("root node must be defined first at level 0"))
			return d
		}
	} else {
		if level == 0 {
			d.setError(errors.New("root node is already defined"))
			return d
		}
		if level > d.currentLevel+1 {
			d.setError(fmt.Errorf("cannot jump from level %d to level %d", d.currentLevel, level))
			return d
		}
	}

	line := fmt.Sprintf("%s%s", strings.Repeat("    ", level+1), trimmedText)
	d.body = append(d.body, line)
	d.currentLevel = level
	if level == 0 {
		d.hasRoot = true
	}
	return d
}

// Child adds a child node under the current node.
func (d *Diagram) Child(text string) *Diagram {
	if d.err != nil {
		return d
	}
	if !d.hasRoot {
		d.setError(errors.New("root node must be defined first"))
		return d
	}

	return d.Node(d.currentLevel+1, text)
}

// Sibling adds a sibling node at the current level.
func (d *Diagram) Sibling(text string) *Diagram {
	if d.err != nil {
		return d
	}
	if !d.hasRoot {
		d.setError(errors.New("root node must be defined first"))
		return d
	}

	return d.Node(d.currentLevel, text)
}

// Parent moves the current level one level up.
func (d *Diagram) Parent() *Diagram {
	if d.err != nil {
		return d
	}
	if !d.hasRoot {
		d.setError(errors.New("root node must be defined first"))
		return d
	}
	if d.currentLevel <= 0 {
		d.setError(errors.New("already at root level"))
		return d
	}

	d.currentLevel--
	return d
}

// LF adds a line feed to the mindmap diagram.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}

	d.body = append(d.body, "")
	return d
}

func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

func validateName(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s must not be empty", fieldName)
	}
	if containsNewline(trimmed) {
		return "", fmt.Errorf("%s must not contain newline characters", fieldName)
	}
	return trimmed, nil
}

func containsNewline(v string) bool {
	return strings.ContainsAny(v, "\n\r")
}
