// Package userjourney is mermaid user journey diagram builder.
//
// Ref. https://mermaid.js.org/syntax/userJourney.html
package userjourney

import (
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// Score is a task sentiment score in a user journey.
//
// Mermaid supports integer scores from 1 to 5.
type Score int

const (
	// ScoreVeryDissatisfied is the lowest sentiment score.
	ScoreVeryDissatisfied Score = 1
	// ScoreDissatisfied is a low sentiment score.
	ScoreDissatisfied Score = 2
	// ScoreNeutral is a neutral sentiment score.
	ScoreNeutral Score = 3
	// ScoreSatisfied is a high sentiment score.
	ScoreSatisfied Score = 4
	// ScoreVerySatisfied is the highest sentiment score.
	ScoreVerySatisfied Score = 5
)

// Diagram is a user journey diagram builder.
type Diagram struct {
	// body is user journey diagram body.
	body []string
	// config is the configuration for the user journey diagram.
	config *config
	// dest is output destination for user journey diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the user journey diagram building.
	err error
	// currentSection is the latest section set by Section().
	currentSection string
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := []string{"journey"}
	if c.title != noTitle {
		lines = append(lines, fmt.Sprintf("    title %s", c.title))
	}

	return &Diagram{
		body:   lines,
		dest:   w,
		config: c,
	}
}

// String returns the user journey diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the user journey diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the user journey diagram body to the output destination.
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

// Section starts a new section in the user journey.
func (d *Diagram) Section(name string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		d.setError(fmt.Errorf("section name must not be empty"))
		return d
	}

	d.currentSection = trimmed
	d.body = append(d.body, fmt.Sprintf("    section %s", trimmed))
	return d
}

// Task adds a task to the current section.
//
// Section must be called before Task; otherwise an error is recorded.
func (d *Diagram) Task(name string, score Score, actors ...string) *Diagram {
	if d.err != nil {
		return d
	}

	if d.currentSection == "" {
		d.setError(fmt.Errorf("task %q requires a section; call Section first", name))
		return d
	}

	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		d.setError(fmt.Errorf("task name must not be empty"))
		return d
	}

	if !isValidScore(score) {
		d.setError(
			fmt.Errorf(
				"invalid score %d: must be between %d and %d",
				score,
				ScoreVeryDissatisfied,
				ScoreVerySatisfied,
			),
		)
		return d
	}

	taskLine := fmt.Sprintf("        %s: %d", trimmedName, score)
	trimmedActors := normalizeActors(actors...)
	if len(trimmedActors) > 0 {
		taskLine = fmt.Sprintf("%s: %s", taskLine, strings.Join(trimmedActors, ", "))
	}

	d.body = append(d.body, taskLine)
	return d
}

// TaskIn adds a task to the specified section.
//
// If the specified section is different from the current section, TaskIn
// automatically starts the new section.
func (d *Diagram) TaskIn(section, task string, score Score, actors ...string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedSection := strings.TrimSpace(section)
	if trimmedSection == "" {
		d.setError(fmt.Errorf("section name must not be empty"))
		return d
	}

	if d.currentSection != trimmedSection {
		d.Section(trimmedSection)
	}
	return d.Task(task, score, actors...)
}

// LF adds a line feed to the user journey diagram.
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

func isValidScore(score Score) bool {
	return score >= ScoreVeryDissatisfied && score <= ScoreVerySatisfied
}

func normalizeActors(actors ...string) []string {
	normalized := make([]string, 0, len(actors))
	for _, actor := range actors {
		trimmed := strings.TrimSpace(actor)
		// Ignore blank actors so callers can safely pass optional values.
		if trimmed != "" {
			normalized = append(normalized, trimmed)
		}
	}
	return normalized
}
