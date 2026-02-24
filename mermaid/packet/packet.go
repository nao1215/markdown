// Package packet is mermaid packet diagram builder.
//
// Ref. https://mermaid.js.org/syntax/packet.html
package packet

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// packetLinesCap is the max initial lines in packet diagram.
	packetLinesCap int = 2
)

// Diagram is a packet diagram builder.
type Diagram struct {
	// body is packet diagram body.
	body []string
	// dest is output destination for packet diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the packet diagram building.
	err error
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := make([]string, 0, packetLinesCap)
	lines = append(lines, "packet")

	trimmedTitle := strings.TrimSpace(c.title)
	if trimmedTitle != noTitle {
		if containsNewline(trimmedTitle) {
			return &Diagram{
				body: []string{"packet"},
				dest: w,
				err:  errors.New("title must not contain newline characters"),
			}
		}
		lines = append(lines, fmt.Sprintf("    title %s", trimmedTitle))
	}

	return &Diagram{
		body: lines,
		dest: w,
	}
}

// String returns the packet diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the packet diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the packet diagram body to the output destination.
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

// LF adds a line feed to the packet diagram.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}

	d.body = append(d.body, "")
	return d
}

// Field adds a packet field by start and end bit positions.
//
// When start and end are the same value, a single-bit field is emitted.
func (d *Diagram) Field(start, end int, label string) *Diagram {
	if d.err != nil {
		return d
	}
	if start < 0 {
		d.setError(errors.New("start bit must be greater than or equal to zero"))
		return d
	}
	if end < 0 {
		d.setError(errors.New("end bit must be greater than or equal to zero"))
		return d
	}
	if start > end {
		d.setError(errors.New("start bit must be less than or equal to end bit"))
		return d
	}

	trimmedLabel, err := validateLabel(label)
	if err != nil {
		d.setError(err)
		return d
	}

	bitRange := strconv.Itoa(start)
	if start != end {
		bitRange = fmt.Sprintf("%d-%d", start, end)
	}

	d.body = append(d.body, fmt.Sprintf("    %s: %s", bitRange, quote(trimmedLabel)))
	return d
}

// Bit adds a single-bit packet field.
func (d *Diagram) Bit(position int, label string) *Diagram {
	return d.Field(position, position, label)
}

// Next adds a packet field using +<bits> syntax.
//
// It is useful when field positions are managed incrementally.
func (d *Diagram) Next(bits int, label string) *Diagram {
	if d.err != nil {
		return d
	}
	if bits <= 0 {
		d.setError(errors.New("bit count must be greater than zero"))
		return d
	}

	trimmedLabel, err := validateLabel(label)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, fmt.Sprintf("    +%d: %s", bits, quote(trimmedLabel)))
	return d
}

func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

func validateLabel(label string) (string, error) {
	trimmed := strings.TrimSpace(label)
	if trimmed == "" {
		return "", errors.New("field label must not be empty")
	}
	if containsNewline(trimmed) {
		return "", errors.New("field label must not contain newline characters")
	}
	return trimmed, nil
}

func containsNewline(value string) bool {
	return strings.ContainsAny(value, "\n\r")
}

func quote(value string) string {
	escaped := strings.ReplaceAll(normalizeQuoted(value), `\`, "&#92;")
	escaped = strings.ReplaceAll(escaped, "\r", "&#92;r")
	escaped = strings.ReplaceAll(escaped, "\n", "&#92;n")
	escaped = strings.ReplaceAll(escaped, "\t", "&#92;t")
	escaped = strings.ReplaceAll(escaped, `"`, "&quot;")
	return `"` + escaped + `"`
}

func normalizeQuoted(value string) string {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) >= 2 && strings.HasPrefix(trimmed, `"`) && strings.HasSuffix(trimmed, `"`) {
		return trimmed[1 : len(trimmed)-1]
	}
	return trimmed
}
