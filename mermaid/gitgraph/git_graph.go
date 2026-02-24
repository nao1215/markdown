// Package gitgraph is mermaid git graph diagram builder.
//
// Ref. https://mermaid.js.org/syntax/gitgraph.html
package gitgraph

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// CommitType is the commit style in git graph.
type CommitType string

const (
	// CommitTypeNormal is default commit style.
	CommitTypeNormal CommitType = "NORMAL"
	// CommitTypeReverse is reverse commit style.
	CommitTypeReverse CommitType = "REVERSE"
	// CommitTypeHighlight is highlight commit style.
	CommitTypeHighlight CommitType = "HIGHLIGHT"
	// gitGraphLinesCap is the max initial lines with title frontmatter.
	gitGraphLinesCap int = 4
)

// Diagram is a git graph diagram builder.
type Diagram struct {
	// body is git graph diagram body.
	body []string
	// dest is output destination for git graph diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the git graph diagram building.
	err error
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := make([]string, 0, gitGraphLinesCap)
	if c.title != noTitle {
		if containsNewline(c.title) {
			return &Diagram{
				body: []string{"gitGraph"},
				dest: w,
				err:  errors.New("title must not contain newline characters"),
			}
		}
		lines = append(lines, "---")
		lines = append(lines, fmt.Sprintf("title: %s", c.title))
		lines = append(lines, "---")
	}
	lines = append(lines, "gitGraph")

	return &Diagram{
		body: lines,
		dest: w,
	}
}

// String returns the git graph diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the git graph diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the git graph diagram body to the output destination.
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

type commitConfig struct {
	id    string
	tag   string
	cType CommitType
}

func newCommitConfig() *commitConfig {
	return &commitConfig{}
}

// CommitOption sets commit options.
type CommitOption func(*commitConfig)

// WithCommitID sets the commit id.
func WithCommitID(id string) CommitOption {
	return func(c *commitConfig) {
		c.id = id
	}
}

// WithCommitTag sets the commit tag.
func WithCommitTag(tag string) CommitOption {
	return func(c *commitConfig) {
		c.tag = tag
	}
}

// WithCommitType sets the commit type.
func WithCommitType(cType CommitType) CommitOption {
	return func(c *commitConfig) {
		c.cType = cType
	}
}

// Commit adds a commit command to the git graph.
func (d *Diagram) Commit(opts ...CommitOption) *Diagram {
	if d.err != nil {
		return d
	}

	c := newCommitConfig()
	for _, opt := range opts {
		opt(c)
	}

	line, err := formatCommitLikeLine("commit", c)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, line)
	return d
}

type branchConfig struct {
	order    int
	hasOrder bool
}

func newBranchConfig() *branchConfig {
	return &branchConfig{}
}

// BranchOption sets branch options.
type BranchOption func(*branchConfig)

// WithBranchOrder sets the branch order.
// Order must be zero or greater.
func WithBranchOrder(order int) BranchOption {
	return func(c *branchConfig) {
		c.order = order
		c.hasOrder = true
	}
}

// Branch adds a branch command to the git graph.
func (d *Diagram) Branch(name string, opts ...BranchOption) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedName, err := validateRefName("branch name", name)
	if err != nil {
		d.setError(err)
		return d
	}

	c := newBranchConfig()
	for _, opt := range opts {
		opt(c)
	}
	if c.hasOrder && c.order < 0 {
		d.setError(errors.New("branch order must be greater than or equal to zero"))
		return d
	}

	line := fmt.Sprintf("    branch %s", formatRefName(trimmedName))
	if c.hasOrder {
		line = fmt.Sprintf("%s order: %d", line, c.order)
	}

	d.body = append(d.body, line)
	return d
}

// Checkout adds a checkout command to the git graph.
func (d *Diagram) Checkout(name string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedName, err := validateRefName("branch name", name)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, fmt.Sprintf("    checkout %s", formatRefName(trimmedName)))
	return d
}

// Merge adds a merge command to the git graph.
func (d *Diagram) Merge(branch string, opts ...CommitOption) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedBranch, err := validateRefName("branch name", branch)
	if err != nil {
		d.setError(err)
		return d
	}

	c := newCommitConfig()
	for _, opt := range opts {
		opt(c)
	}

	line, err := formatCommitLikeLine(
		fmt.Sprintf("merge %s", formatRefName(trimmedBranch)),
		c,
	)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, line)
	return d
}

type cherryPickConfig struct {
	parentID string
}

func newCherryPickConfig() *cherryPickConfig {
	return &cherryPickConfig{}
}

// CherryPickOption sets cherry-pick options.
type CherryPickOption func(*cherryPickConfig)

// WithCherryPickParent sets the cherry-pick parent commit id.
func WithCherryPickParent(parentID string) CherryPickOption {
	return func(c *cherryPickConfig) {
		c.parentID = parentID
	}
}

// CherryPick adds a cherry-pick command to the git graph.
func (d *Diagram) CherryPick(id string, opts ...CherryPickOption) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedID, err := validateName("commit id", id)
	if err != nil {
		d.setError(err)
		return d
	}

	c := newCherryPickConfig()
	for _, opt := range opts {
		opt(c)
	}

	line := fmt.Sprintf("    cherry-pick id: %s", quote(trimmedID))
	if strings.TrimSpace(c.parentID) != "" {
		trimmedParentID, err := validateName("parent commit id", c.parentID)
		if err != nil {
			d.setError(err)
			return d
		}
		line = fmt.Sprintf("%s parent: %s", line, quote(trimmedParentID))
	}

	d.body = append(d.body, line)
	return d
}

// Reset adds a reset command to the git graph.
func (d *Diagram) Reset(id string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedID, err := validateName("commit id", id)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, fmt.Sprintf("    reset id: %s", quote(trimmedID)))
	return d
}

// LF adds a line feed to the git graph.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}

	d.body = append(d.body, "")
	return d
}

func formatCommitLikeLine(command string, c *commitConfig) (string, error) {
	line := fmt.Sprintf("    %s", command)

	if strings.TrimSpace(c.id) != "" {
		trimmedID, err := validateName("commit id", c.id)
		if err != nil {
			return "", err
		}
		line = fmt.Sprintf("%s id: %s", line, quote(trimmedID))
	}

	if strings.TrimSpace(c.tag) != "" {
		trimmedTag, err := validateName("commit tag", c.tag)
		if err != nil {
			return "", err
		}
		line = fmt.Sprintf("%s tag: %s", line, quote(trimmedTag))
	}

	if c.cType != "" {
		if !isValidCommitType(c.cType) {
			return "", fmt.Errorf("invalid commit type %q", c.cType)
		}
		line = fmt.Sprintf("%s type: %s", line, c.cType)
	}

	return line, nil
}

func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

func isValidCommitType(cType CommitType) bool {
	return cType == CommitTypeNormal || cType == CommitTypeReverse || cType == CommitTypeHighlight
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

func validateRefName(fieldName, value string) (string, error) {
	trimmed, err := validateName(fieldName, value)
	if err != nil {
		return "", err
	}
	return normalizeQuoted(trimmed), nil
}

func formatRefName(name string) string {
	if containsWhitespace(name) || isKeyword(name) {
		return quote(name)
	}
	return name
}

func quote(v string) string {
	s := normalizeQuoted(v)
	escaped := strings.ReplaceAll(s, `\`, "&#92;")
	escaped = strings.ReplaceAll(escaped, `"`, "&quot;")
	return `"` + escaped + `"`
}

func normalizeQuoted(v string) string {
	trimmed := strings.TrimSpace(v)
	if len(trimmed) >= 2 && strings.HasPrefix(trimmed, `"`) && strings.HasSuffix(trimmed, `"`) {
		return trimmed[1 : len(trimmed)-1]
	}
	return trimmed
}

func containsWhitespace(v string) bool {
	return strings.ContainsAny(v, " \t")
}

func containsNewline(v string) bool {
	return strings.ContainsAny(v, "\n\r")
}

func isKeyword(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "commit", "branch", "checkout", "merge", "cherry-pick", "reset":
		return true
	default:
		return false
	}
}
