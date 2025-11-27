// Package gantt is a mermaid Gantt chart builder.
package gantt

import (
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// Chart is a Gantt chart builder.
type Chart struct {
	// body is Gantt chart body.
	body []string
	// config is the configuration for the Gantt chart.
	config *config
	// dest is output destination for Gantt chart body.
	dest io.Writer
	// err manages errors that occur in all parts of the Gantt chart building.
	err error
}

// NewChart returns a new Chart.
func NewChart(w io.Writer, opts ...Option) *Chart {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := []string{"gantt"}
	if c.title != noTitle {
		lines = append(lines, fmt.Sprintf("    title %s", c.title))
	}
	if c.dateFormat != "" {
		lines = append(lines, fmt.Sprintf("    dateFormat %s", c.dateFormat))
	}
	if c.axisFormat != "" {
		lines = append(lines, fmt.Sprintf("    axisFormat %s", c.axisFormat))
	}
	if c.tickInterval != "" {
		lines = append(lines, fmt.Sprintf("    tickInterval %s", c.tickInterval))
	}
	if c.todayMarker != "" {
		lines = append(lines, fmt.Sprintf("    todayMarker %s", c.todayMarker))
	}
	for _, exclude := range c.excludes {
		lines = append(lines, fmt.Sprintf("    excludes %s", exclude))
	}

	return &Chart{
		body:   lines,
		dest:   w,
		config: c,
	}
}

// String returns the Gantt chart body.
func (c *Chart) String() string {
	return strings.Join(c.body, internal.LineFeed())
}

// Error returns the error that occurred during the Gantt chart building.
func (c *Chart) Error() error {
	return c.err
}

// Build writes the Gantt chart body to the output destination.
func (c *Chart) Build() error {
	if _, err := fmt.Fprint(c.dest, c.String()); err != nil {
		if c.err != nil {
			return fmt.Errorf("failed to write: %w: %s", err, c.err.Error()) //nolint:wrapcheck
		}
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

// Section adds a section to the Gantt chart.
func (c *Chart) Section(name string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    section %s", name))
	return c
}

// Task adds a task to the Gantt chart.
// startDate can be a date string (e.g., "2024-01-01") or "after taskID".
// duration can be a duration string (e.g., "30d", "1w") or an end date.
func (c *Chart) Task(name, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :%s, %s", name, startDate, duration))
	return c
}

// TaskWithID adds a task with an ID to the Gantt chart.
func (c *Chart) TaskWithID(name, id, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :%s, %s, %s", name, id, startDate, duration))
	return c
}

// CriticalTask adds a critical task to the Gantt chart.
func (c *Chart) CriticalTask(name, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :crit, %s, %s", name, startDate, duration))
	return c
}

// CriticalTaskWithID adds a critical task with an ID to the Gantt chart.
func (c *Chart) CriticalTaskWithID(name, id, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :crit, %s, %s, %s", name, id, startDate, duration))
	return c
}

// ActiveTask adds an active task to the Gantt chart.
func (c *Chart) ActiveTask(name, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :active, %s, %s", name, startDate, duration))
	return c
}

// ActiveTaskWithID adds an active task with an ID to the Gantt chart.
func (c *Chart) ActiveTaskWithID(name, id, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :active, %s, %s, %s", name, id, startDate, duration))
	return c
}

// DoneTask adds a done task to the Gantt chart.
func (c *Chart) DoneTask(name, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :done, %s, %s", name, startDate, duration))
	return c
}

// DoneTaskWithID adds a done task with an ID to the Gantt chart.
func (c *Chart) DoneTaskWithID(name, id, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :done, %s, %s, %s", name, id, startDate, duration))
	return c
}

// CriticalActiveTask adds a critical and active task to the Gantt chart.
func (c *Chart) CriticalActiveTask(name, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :crit, active, %s, %s", name, startDate, duration))
	return c
}

// CriticalActiveTaskWithID adds a critical and active task with an ID to the Gantt chart.
func (c *Chart) CriticalActiveTaskWithID(name, id, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :crit, active, %s, %s, %s", name, id, startDate, duration))
	return c
}

// CriticalDoneTask adds a critical and done task to the Gantt chart.
func (c *Chart) CriticalDoneTask(name, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :crit, done, %s, %s", name, startDate, duration))
	return c
}

// CriticalDoneTaskWithID adds a critical and done task with an ID to the Gantt chart.
func (c *Chart) CriticalDoneTaskWithID(name, id, startDate, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :crit, done, %s, %s, %s", name, id, startDate, duration))
	return c
}

// Milestone adds a milestone to the Gantt chart.
// A milestone is a task with 0 duration.
func (c *Chart) Milestone(name, date string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :milestone, %s, 0d", name, date))
	return c
}

// MilestoneWithID adds a milestone with an ID to the Gantt chart.
func (c *Chart) MilestoneWithID(name, id, date string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :milestone, %s, %s, 0d", name, id, date))
	return c
}

// CriticalMilestone adds a critical milestone to the Gantt chart.
func (c *Chart) CriticalMilestone(name, date string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :crit, milestone, %s, 0d", name, date))
	return c
}

// CriticalMilestoneWithID adds a critical milestone with an ID to the Gantt chart.
func (c *Chart) CriticalMilestoneWithID(name, id, date string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :crit, milestone, %s, %s, 0d", name, id, date))
	return c
}

// TaskAfter adds a task that starts after another task.
func (c *Chart) TaskAfter(name, afterTaskID, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :after %s, %s", name, afterTaskID, duration))
	return c
}

// TaskAfterWithID adds a task with ID that starts after another task.
func (c *Chart) TaskAfterWithID(name, id, afterTaskID, duration string) *Chart {
	c.body = append(c.body, fmt.Sprintf("    %s :%s, after %s, %s", name, id, afterTaskID, duration))
	return c
}

// LF adds a line feed to the Gantt chart.
func (c *Chart) LF() *Chart {
	c.body = append(c.body, "")
	return c
}
