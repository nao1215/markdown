// Package c4 is mermaid C4 context diagram builder.
//
// Ref. https://mermaid.js.org/syntax/c4.html
//
// A C4 context diagram is the top level of the C4 model: the people and the
// software systems around the system being described, and the relationships
// between them. This package builds the C4Context form of it, and keeps the
// flat, chainable API the rest of the library has: an element is one call,
// a relationship is one call, and a boundary is a pair of calls that the
// elements between them belong to, the same way mermaid/treemap walks its own
// nesting.
//
// mermaid marks its C4 support experimental and says the syntax may change.
// The scope here is deliberately narrow because of that: C4Context with its
// elements, boundaries and relationships. C4Container, C4Component, C4Dynamic,
// C4Deployment and the UpdateElementStyle family are left out, and each can
// arrive later as a new method without changing anything that exists.
//
// Errors are recorded rather than returned from every call: the chain runs to
// the end and the error surfaces from Build. A nil writer and a writer that
// refuses the diagram are both reported rather than causing a panic, and
// String returns the diagram without needing a writer at all.
package c4

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// c4LinesCap is the number of lines a small diagram starts with.
	c4LinesCap int = 8
	// header is the keyword the diagram starts with.
	header = "C4Context"
	// indentUnit is one level of boundary nesting.
	indentUnit = "    "
)

// Diagram is a C4 context diagram builder.
type Diagram struct {
	// body is the C4 context diagram body.
	body []string
	// dest is the output destination for the C4 context diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the diagram building.
	err error
	// depth is how many boundaries are open, counted in indentUnits.
	depth int
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()
	for _, opt := range opts {
		opt(c)
	}

	lines := make([]string, 0, c4LinesCap)
	lines = append(lines, header)

	d := &Diagram{
		body: lines,
		dest: w,
	}

	trimmedTitle := strings.TrimSpace(c.title)
	if trimmedTitle == noTitle {
		return d
	}
	if containsNewline(trimmedTitle) {
		// mermaid reads the rest of the line as the title, so a second line
		// would be read as a statement of its own.
		d.err = errors.New("title must not contain newline characters")
		return d
	}
	d.body = append(d.body, indentUnit+"title "+escapeStatement(trimmedTitle))
	return d
}

// String returns the C4 context diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the C4 context diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the C4 context diagram body to the output destination.
//
// A boundary left open is reported here rather than written out: mermaid
// refuses a diagram whose brace never closes, which loses the whole drawing.
func (d *Diagram) Build() error {
	if d.err != nil {
		return d.err
	}
	if d.depth != 0 {
		d.err = fmt.Errorf("%d boundary must be closed with BoundaryEnd before Build", d.depth)
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

// LF adds a line feed to the C4 context diagram body.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}
	d.body = append(d.body, "")
	return d
}

// Person adds a person inside the enterprise being described.
func (d *Diagram) Person(id, label string, opts ...ElementOption) *Diagram {
	return d.element("Person", id, label, opts)
}

// PersonExt adds a person outside the enterprise being described. mermaid
// spells the macro Person_Ext and draws it in a colder color than Person.
func (d *Diagram) PersonExt(id, label string, opts ...ElementOption) *Diagram {
	return d.element("Person_Ext", id, label, opts)
}

// System adds a software system inside the enterprise being described.
func (d *Diagram) System(id, label string, opts ...ElementOption) *Diagram {
	return d.element("System", id, label, opts)
}

// SystemExt adds a software system outside the enterprise being described.
// mermaid spells the macro System_Ext.
func (d *Diagram) SystemExt(id, label string, opts ...ElementOption) *Diagram {
	return d.element("System_Ext", id, label, opts)
}

// SystemDb adds a software system drawn as a database.
func (d *Diagram) SystemDb(id, label string, opts ...ElementOption) *Diagram {
	return d.element("SystemDb", id, label, opts)
}

// SystemQueue adds a software system drawn as a queue.
func (d *Diagram) SystemQueue(id, label string, opts ...ElementOption) *Diagram {
	return d.element("SystemQueue", id, label, opts)
}

// Boundary opens a boundary. Everything added until the matching BoundaryEnd is
// drawn inside it.
//
// Boundaries nest, and the pair is explicit rather than a callback taking a
// nested builder: the chain stays flat, which is how every other builder in
// this library reads, and how mermaid/treemap walks its own nesting with
// Section and Parent.
func (d *Diagram) Boundary(id, label string, opts ...BoundaryOption) *Diagram {
	if d.err != nil {
		return d
	}

	c := newBoundaryConfig()
	for _, opt := range opts {
		opt(c)
	}

	args, err := d.boundaryArgs(id, label, strings.TrimSpace(c.boundaryType))
	if err != nil {
		d.setError(err)
		return d
	}
	return d.openBoundary("Boundary", args)
}

// EnterpriseBoundary opens a boundary drawn as the enterprise. mermaid spells
// the macro Enterprise_Boundary.
func (d *Diagram) EnterpriseBoundary(id, label string) *Diagram {
	if d.err != nil {
		return d
	}

	args, err := d.boundaryArgs(id, label, "")
	if err != nil {
		d.setError(err)
		return d
	}
	return d.openBoundary("Enterprise_Boundary", args)
}

// SystemBoundary opens a boundary drawn as a software system. mermaid spells
// the macro System_Boundary.
func (d *Diagram) SystemBoundary(id, label string) *Diagram {
	if d.err != nil {
		return d
	}

	args, err := d.boundaryArgs(id, label, "")
	if err != nil {
		d.setError(err)
		return d
	}
	return d.openBoundary("System_Boundary", args)
}

// BoundaryEnd closes the boundary opened last.
//
// Calling it outside a boundary is an error rather than a silent no-op: a chain
// that has closed more boundaries than it opened is not the diagram its author
// meant.
func (d *Diagram) BoundaryEnd() *Diagram {
	if d.err != nil {
		return d
	}
	if d.depth == 0 {
		d.setError(errors.New("BoundaryEnd was called outside a boundary; there is nothing to close"))
		return d
	}

	d.depth--
	d.body = append(d.body, d.indent()+"}")
	return d
}

// Rel adds a one way relationship from one element to another.
func (d *Diagram) Rel(from, to, label string, opts ...RelationOption) *Diagram {
	return d.relation("Rel", from, to, label, opts)
}

// BiRel adds a relationship drawn with an arrowhead at each end.
func (d *Diagram) BiRel(from, to, label string, opts ...RelationOption) *Diagram {
	return d.relation("BiRel", from, to, label, opts)
}

// element writes one element macro.
func (d *Diagram) element(macro, id, label string, opts []ElementOption) *Diagram {
	if d.err != nil {
		return d
	}

	c := newElementConfig()
	for _, opt := range opts {
		opt(c)
	}

	validID, err := validateIdentifier(macro+" id", id)
	if err != nil {
		d.setError(err)
		return d
	}
	text, err := validateText(macro+" label", label)
	if err != nil {
		d.setError(err)
		return d
	}

	args := []string{validID, quote(text)}
	description := strings.TrimSpace(c.description)
	if description != "" {
		if containsNewline(description) {
			d.setError(fmt.Errorf("description of %s %q must not contain newline characters", macro, text))
			return d
		}
		args = append(args, quote(description))
	}

	d.body = append(d.body, fmt.Sprintf("%s%s(%s)", d.indent(), macro, strings.Join(args, ", ")))
	return d
}

// relation writes one relationship macro.
func (d *Diagram) relation(macro, from, to, label string, opts []RelationOption) *Diagram {
	if d.err != nil {
		return d
	}

	c := newRelationConfig()
	for _, opt := range opts {
		opt(c)
	}

	validFrom, err := validateIdentifier(macro+" source", from)
	if err != nil {
		d.setError(err)
		return d
	}
	validTo, err := validateIdentifier(macro+" destination", to)
	if err != nil {
		d.setError(err)
		return d
	}
	text, err := validateText(macro+" label", label)
	if err != nil {
		d.setError(err)
		return d
	}

	args := []string{validFrom, validTo, quote(text)}
	technology := strings.TrimSpace(c.technology)
	if technology != "" {
		if containsNewline(technology) {
			d.setError(fmt.Errorf("technology of %s %q must not contain newline characters", macro, text))
			return d
		}
		args = append(args, quote(technology))
	}

	d.body = append(d.body, fmt.Sprintf("%s%s(%s)", d.indent(), macro, strings.Join(args, ", ")))
	return d
}

// boundaryArgs validates a boundary's parts and returns them ready to be
// written between the macro's parentheses.
func (d *Diagram) boundaryArgs(id, label, boundaryType string) ([]string, error) {
	validID, err := validateIdentifier("boundary id", id)
	if err != nil {
		return nil, err
	}
	text, err := validateText("boundary label", label)
	if err != nil {
		return nil, err
	}

	args := []string{validID, quote(text)}
	if boundaryType != "" {
		if containsNewline(boundaryType) {
			return nil, fmt.Errorf("type of boundary %q must not contain newline characters", text)
		}
		args = append(args, quote(boundaryType))
	}
	return args, nil
}

// openBoundary writes the opening line of a boundary and descends into it.
func (d *Diagram) openBoundary(macro string, args []string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("%s%s(%s) {", d.indent(), macro, strings.Join(args, ", ")))
	d.depth++
	return d
}

// indent returns the prefix of a line at the current nesting.
func (d *Diagram) indent() string {
	return indentUnit + strings.Repeat(indentUnit, d.depth)
}

// setError records the first error and leaves any later one alone, because the
// first is the one that explains the rest.
func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

// validateText returns value ready to be quoted into a macro argument.
func validateText(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s must not be empty", fieldName)
	}
	if containsNewline(trimmed) {
		// One macro is one line, so a label spanning lines would be read as a
		// statement of its own. "<br/>" is how a C4 label breaks a line.
		return "", fmt.Errorf("%s must not contain newline characters", fieldName)
	}
	return trimmed, nil
}

// validateIdentifier returns value ready to be written as a macro's bare
// identifier.
//
// An identifier names an element for the relationships that point at it and is
// never drawn, so nothing is escaped here: what a caller passes is what the
// other macros have to spell. A comma ends the argument it is written in, and
// the remaining characters below are the macro's own punctuation. mermaid
// happens to accept some of them today, because its lexer stops at the first
// comma, but an identifier that only works by accident is worth rejecting while
// this package is new.
func validateIdentifier(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s must not be empty", fieldName)
	}
	if containsNewline(trimmed) {
		return "", fmt.Errorf("%s must not contain newline characters", fieldName)
	}
	if i := strings.IndexAny(trimmed, `,()"{}`); i >= 0 {
		return "", fmt.Errorf("%s must not contain %q, which is C4 macro syntax", fieldName, trimmed[i:i+1])
	}
	return trimmed, nil
}

// quote returns value as a quoted macro argument.
//
// A double quote ends the argument, and neither of the two escapes that look
// obvious works: a backslash makes mermaid refuse the diagram, and doubling the
// quote silently splits the argument in two, so Rel(a, b, "x""y") draws "y" as
// the technology rather than as part of the label. What mermaid does implement
// is its own entity form, "#quot;", which was found by rendering all three.
//
// A "#" is escaped first and for the same reason: mermaid reads "#word;" and
// "#123;" as entities, so a label holding "#quot;" or "#39;" would otherwise
// come back as a character the caller never wrote. Both replacements happen in
// one pass, because a second pass over the first one's output would escape the
// "#" it just wrote.
func quote(value string) string {
	var b strings.Builder
	b.Grow(len(value) + 2) //nolint:mnd // the two quotes around the value.
	b.WriteByte('"')
	for _, r := range value {
		switch r {
		case '#':
			b.WriteString("#35;")
		case '"':
			b.WriteString("#quot;")
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

// escapeStatement returns value ready to be written as the rest of a statement
// line.
//
// The title is the one place a C4 diagram takes unquoted text, so the argument
// escaping above does not apply: a "#quot;" would be read as an entity but a
// plain quotation mark is drawn as itself. What the unquoted lexer refuses is
// "#", which starts a comment, and ";", which ends the statement; both are
// written as the entity form mermaid decodes back into them. A "<" that does
// not open a "<br/>" passes the lexer and is then eaten by the renderer's
// sanitizer, so it is written the same way.
func escapeStatement(value string) string {
	var b strings.Builder
	b.Grow(len(value))
	for i, r := range value {
		switch {
		case r == '#':
			b.WriteString("#35;")
		case r == ';':
			b.WriteString("#59;")
		case r == '<' && !strings.HasPrefix(value[i:], "<br/>"):
			b.WriteString("#60;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// containsNewline reports whether value holds a line ending of either kind.
func containsNewline(value string) bool {
	return strings.ContainsAny(value, "\n\r")
}
