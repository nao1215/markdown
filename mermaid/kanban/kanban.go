// Package kanban is mermaid kanban diagram builder.
//
// Ref. https://mermaid.js.org/syntax/kanban.html
package kanban

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// kanbanLinesCap is the max initial lines in kanban.
	kanbanLinesCap int = 8
	// taskMetadataCap is the max metadata entries on a task line.
	taskMetadataCap int = 3
)

// Priority is a task priority in kanban metadata.
type Priority string

const (
	// PriorityVeryLow is the lowest task priority.
	PriorityVeryLow Priority = "Very Low"
	// PriorityLow is a low task priority.
	PriorityLow Priority = "Low"
	// PriorityHigh is a high task priority.
	PriorityHigh Priority = "High"
	// PriorityVeryHigh is the highest task priority.
	PriorityVeryHigh Priority = "Very High"
)

// Diagram is a kanban diagram builder.
type Diagram struct {
	// body is kanban diagram body.
	body []string
	// dest is output destination for kanban diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the kanban diagram building.
	err error
	// currentColumn is the latest column set by Column().
	currentColumn string
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()
	for _, opt := range opts {
		opt(c)
	}

	trimmedTitle := strings.TrimSpace(c.title)
	if containsNewline(trimmedTitle) {
		return &Diagram{
			body: []string{"kanban"},
			dest: w,
			err:  errors.New("title must not contain newline characters"),
		}
	}

	trimmedTicketBaseURL := strings.TrimSpace(c.ticketBaseURL)
	if containsNewline(trimmedTicketBaseURL) {
		return &Diagram{
			body: []string{"kanban"},
			dest: w,
			err:  errors.New("ticket base URL must not contain newline characters"),
		}
	}

	lines := make([]string, 0, kanbanLinesCap)
	if trimmedTitle != noTitle || trimmedTicketBaseURL != noTicketBaseURL {
		lines = append(lines, "---")
		if trimmedTitle != noTitle {
			lines = append(lines, fmt.Sprintf("title: %s", trimmedTitle))
		}
		if trimmedTicketBaseURL != noTicketBaseURL {
			lines = append(lines, "config:")
			lines = append(lines, "  kanban:")
			lines = append(lines, fmt.Sprintf("    ticketBaseUrl: %s", quoteYAML(trimmedTicketBaseURL)))
		}
		lines = append(lines, "---")
	}
	lines = append(lines, "kanban")

	return &Diagram{
		body: lines,
		dest: w,
	}
}

// String returns the kanban diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the kanban diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the kanban diagram body to the output destination.
func (d *Diagram) Build() error {
	if d.err != nil {
		return d.err
	}
	if d.dest == nil {
		d.err = errors.New("output writer must not be nil")
		return d.err
	}

	if _, err := fmt.Fprint(d.dest, d.String()); err != nil {
		d.err = fmt.Errorf("failed to write: %w", err)
		return d.err
	}
	return nil
}

// ColumnOption sets column options.
type ColumnOption func(*columnConfig)

type columnConfig struct {
	id string
}

func newColumnConfig() *columnConfig {
	return &columnConfig{}
}

// WithColumnID sets the column identifier.
func WithColumnID(id string) ColumnOption {
	return func(c *columnConfig) {
		c.id = id
	}
}

// Column adds a kanban column.
func (d *Diagram) Column(name string, opts ...ColumnOption) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedName, err := validateCardLabel("column name", name)
	if err != nil {
		d.setError(err)
		return d
	}

	c := newColumnConfig()
	for _, opt := range opts {
		opt(c)
	}

	trimmedID, err := validateOptionalIdentifier("column id", c.id)
	if err != nil {
		d.setError(err)
		return d
	}

	d.currentColumn = trimmedName
	d.body = append(d.body, "    "+formatCard(trimmedID, trimmedName))
	return d
}

// TaskOption sets task options.
type TaskOption func(*taskConfig)

type taskConfig struct {
	id       string
	ticket   string
	assigned string
	priority Priority
}

func newTaskConfig() *taskConfig {
	return &taskConfig{}
}

// WithTaskID sets the task identifier.
func WithTaskID(id string) TaskOption {
	return func(c *taskConfig) {
		c.id = id
	}
}

// WithTaskTicket sets the task ticket metadata.
func WithTaskTicket(ticket string) TaskOption {
	return func(c *taskConfig) {
		c.ticket = ticket
	}
}

// WithTaskAssigned sets the task assignee metadata.
func WithTaskAssigned(assigned string) TaskOption {
	return func(c *taskConfig) {
		c.assigned = assigned
	}
}

// WithTaskPriority sets the task priority metadata.
func WithTaskPriority(priority Priority) TaskOption {
	return func(c *taskConfig) {
		c.priority = priority
	}
}

// Task adds a task to the current column.
//
// Column must be called before Task; otherwise an error is recorded.
func (d *Diagram) Task(name string, opts ...TaskOption) *Diagram {
	if d.err != nil {
		return d
	}
	if d.currentColumn == "" {
		d.setError(fmt.Errorf("task %q requires a column; call Column first", name))
		return d
	}

	trimmedName, err := validateCardLabel("task name", name)
	if err != nil {
		d.setError(err)
		return d
	}

	c := newTaskConfig()
	for _, opt := range opts {
		opt(c)
	}

	trimmedID, err := validateOptionalIdentifier("task id", c.id)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedTicket, err := validateOptionalText("task ticket", c.ticket)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedAssigned, err := validateOptionalText("task assignee", c.assigned)
	if err != nil {
		d.setError(err)
		return d
	}

	normalizedPriority := ""
	if strings.TrimSpace(string(c.priority)) != "" {
		var ok bool
		normalizedPriority, ok = normalizePriority(c.priority)
		if !ok {
			d.setError(fmt.Errorf("invalid priority %q", c.priority))
			return d
		}
	}

	line := formatCard(trimmedID, trimmedName)
	metadata := make([]string, 0, taskMetadataCap)
	if trimmedTicket != "" {
		metadata = append(metadata, fmt.Sprintf("ticket: %s", quoteMetadata(trimmedTicket)))
	}
	if trimmedAssigned != "" {
		metadata = append(metadata, fmt.Sprintf("assigned: %s", quoteMetadata(trimmedAssigned)))
	}
	if normalizedPriority != "" {
		metadata = append(metadata, fmt.Sprintf("priority: %s", quoteMetadata(normalizedPriority)))
	}
	if len(metadata) > 0 {
		line += fmt.Sprintf("@{ %s }", strings.Join(metadata, ", "))
	}

	d.body = append(d.body, "        "+line)
	return d
}

// TaskIn adds a task to the specified column.
//
// If the specified column is different from the current column, TaskIn
// automatically starts the new column.
func (d *Diagram) TaskIn(column, task string, opts ...TaskOption) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedColumn, err := validateCardLabel("column name", column)
	if err != nil {
		d.setError(err)
		return d
	}
	if d.currentColumn != trimmedColumn {
		d.Column(trimmedColumn)
	}

	return d.Task(task, opts...)
}

// LF adds a line feed to the kanban diagram.
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

func normalizePriority(priority Priority) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(string(priority))) {
	case strings.ToLower(string(PriorityVeryLow)):
		return string(PriorityVeryLow), true
	case strings.ToLower(string(PriorityLow)):
		return string(PriorityLow), true
	case strings.ToLower(string(PriorityHigh)):
		return string(PriorityHigh), true
	case strings.ToLower(string(PriorityVeryHigh)):
		return string(PriorityVeryHigh), true
	default:
		return "", false
	}
}

func validateCardLabel(fieldName, value string) (string, error) {
	trimmed, err := validateText(fieldName, value)
	if err != nil {
		return "", err
	}

	normalized := normalizeBracketed(trimmed)
	if normalized == "" {
		return "", fmt.Errorf("%s must not be empty", fieldName)
	}
	if strings.ContainsAny(normalized, "[]") {
		return "", fmt.Errorf("%s must not contain '[' or ']'", fieldName)
	}
	return normalized, nil
}

func validateOptionalIdentifier(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}
	return validateIdentifier(fieldName, trimmed)
}

func validateIdentifier(fieldName, value string) (string, error) {
	trimmed, err := validateText(fieldName, value)
	if err != nil {
		return "", err
	}
	if strings.ContainsAny(trimmed, " \t") {
		return "", fmt.Errorf("%s must not contain whitespace", fieldName)
	}
	return trimmed, nil
}

func validateOptionalText(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}
	return validateText(fieldName, trimmed)
}

func validateText(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s must not be empty", fieldName)
	}
	if containsNewline(trimmed) {
		return "", fmt.Errorf("%s must not contain newline characters", fieldName)
	}
	return trimmed, nil
}

func formatCard(id, label string) string {
	if id == "" {
		return fmt.Sprintf("[%s]", label)
	}
	return fmt.Sprintf("%s[%s]", id, label)
}

func normalizeBracketed(value string) string {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) >= 2 && strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
		return strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	}
	return trimmed
}

func containsNewline(value string) bool {
	return strings.ContainsAny(value, "\n\r")
}

func quoteMetadata(value string) string {
	escaped := strings.ReplaceAll(value, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `'`, `\'`)
	escaped = strings.ReplaceAll(escaped, "\r", `\r`)
	escaped = strings.ReplaceAll(escaped, "\n", `\n`)
	escaped = strings.ReplaceAll(escaped, "\t", `\t`)
	return "'" + escaped + "'"
}

func quoteYAML(value string) string {
	escaped := strings.ReplaceAll(value, `'`, `''`)
	return "'" + escaped + "'"
}
