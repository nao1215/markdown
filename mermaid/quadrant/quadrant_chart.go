// Package quadrant is mermaid quadrant chart builder.
package quadrant

import (
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// Chart is a quadrant chart builder.
type Chart struct {
	// body is quadrant chart body.
	body []string
	// config is the configuration for the quadrant chart.
	config *config
	// dest is output destination for quadrant chart body.
	dest io.Writer
	// err manages errors that occur in all parts of the quadrant chart building.
	err error
}

// NewChart returns a new Chart.
func NewChart(w io.Writer, opts ...Option) *Chart {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := []string{"quadrantChart"}
	if c.title != noTitle {
		lines = append(lines, fmt.Sprintf("    title %s", c.title))
	}

	return &Chart{
		body:   lines,
		dest:   w,
		config: c,
	}
}

// String returns the quadrant chart body.
func (ch *Chart) String() string {
	return strings.Join(ch.body, internal.LineFeed())
}

// Error returns the error that occurred during the quadrant chart building.
func (ch *Chart) Error() error {
	return ch.err
}

// Build writes the quadrant chart body to the output destination.
func (ch *Chart) Build() error {
	if _, err := fmt.Fprint(ch.dest, ch.String()); err != nil {
		if ch.err != nil {
			return fmt.Errorf("failed to write: %w: %s", err, ch.err.Error()) //nolint:wrapcheck
		}
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

// XAxis sets the x-axis label.
// If rightLabel is provided, it will be displayed on the right side of the axis.
func (ch *Chart) XAxis(leftLabel string, rightLabel ...string) *Chart {
	if len(rightLabel) > 0 && rightLabel[0] != "" {
		ch.body = append(ch.body, fmt.Sprintf("    x-axis %s --> %s", leftLabel, rightLabel[0]))
	} else {
		ch.body = append(ch.body, fmt.Sprintf("    x-axis %s", leftLabel))
	}
	return ch
}

// YAxis sets the y-axis label.
// If topLabel is provided, it will be displayed on the top of the axis.
func (ch *Chart) YAxis(bottomLabel string, topLabel ...string) *Chart {
	if len(topLabel) > 0 && topLabel[0] != "" {
		ch.body = append(ch.body, fmt.Sprintf("    y-axis %s --> %s", bottomLabel, topLabel[0]))
	} else {
		ch.body = append(ch.body, fmt.Sprintf("    y-axis %s", bottomLabel))
	}
	return ch
}

// Quadrant1 sets the label for quadrant 1 (top-right).
func (ch *Chart) Quadrant1(label string) *Chart {
	ch.body = append(ch.body, fmt.Sprintf("    quadrant-1 %s", label))
	return ch
}

// Quadrant2 sets the label for quadrant 2 (top-left).
func (ch *Chart) Quadrant2(label string) *Chart {
	ch.body = append(ch.body, fmt.Sprintf("    quadrant-2 %s", label))
	return ch
}

// Quadrant3 sets the label for quadrant 3 (bottom-left).
func (ch *Chart) Quadrant3(label string) *Chart {
	ch.body = append(ch.body, fmt.Sprintf("    quadrant-3 %s", label))
	return ch
}

// Quadrant4 sets the label for quadrant 4 (bottom-right).
func (ch *Chart) Quadrant4(label string) *Chart {
	ch.body = append(ch.body, fmt.Sprintf("    quadrant-4 %s", label))
	return ch
}

// Point adds a data point to the quadrant chart.
// x and y should be values between 0 and 1.
func (ch *Chart) Point(name string, x, y float64) *Chart {
	ch.body = append(ch.body, fmt.Sprintf("    %s: [%.2f, %.2f]", name, x, y))
	return ch
}

// PointStyle represents styling options for a point.
type PointStyle struct {
	// Color is the fill color of the point (e.g., "#ff0000").
	Color string
	// Radius is the radius of the point.
	Radius int
	// StrokeColor is the border color of the point.
	StrokeColor string
	// StrokeWidth is the border width of the point (e.g., "5px").
	StrokeWidth string
}

// String returns the style as a mermaid-compatible string.
func (ps PointStyle) String() string {
	var parts []string
	if ps.Color != "" {
		parts = append(parts, fmt.Sprintf("color: %s", ps.Color))
	}
	if ps.Radius > 0 {
		parts = append(parts, fmt.Sprintf("radius: %d", ps.Radius))
	}
	if ps.StrokeColor != "" {
		parts = append(parts, fmt.Sprintf("stroke-color: %s", ps.StrokeColor))
	}
	if ps.StrokeWidth != "" {
		parts = append(parts, fmt.Sprintf("stroke-width: %s", ps.StrokeWidth))
	}
	return strings.Join(parts, ", ")
}

// PointWithStyle adds a data point with custom styling.
// x and y should be values between 0 and 1.
// style can include properties like "radius: 10" or "color: #ff0000".
func (ch *Chart) PointWithStyle(name string, x, y float64, style string) *Chart {
	ch.body = append(ch.body, fmt.Sprintf("    %s: [%.2f, %.2f] %s", name, x, y, style))
	return ch
}

// PointStyled adds a data point with PointStyle struct.
// x and y should be values between 0 and 1.
func (ch *Chart) PointStyled(name string, x, y float64, style PointStyle) *Chart {
	styleStr := style.String()
	if styleStr != "" {
		ch.body = append(ch.body, fmt.Sprintf("    %s: [%.2f, %.2f] %s", name, x, y, styleStr))
	} else {
		ch.body = append(ch.body, fmt.Sprintf("    %s: [%.2f, %.2f]", name, x, y))
	}
	return ch
}

// PointWithClass adds a data point with a class name.
// x and y should be values between 0 and 1.
// className is the name of the class defined by ClassDef.
func (ch *Chart) PointWithClass(name string, x, y float64, className string) *Chart {
	ch.body = append(ch.body, fmt.Sprintf("    %s:::%s: [%.2f, %.2f]", name, className, x, y))
	return ch
}

// PointWithClassAndStyle adds a data point with both a class name and inline style.
// x and y should be values between 0 and 1.
// className is the name of the class defined by ClassDef.
// style can include properties like "radius: 10" or "color: #ff0000".
func (ch *Chart) PointWithClassAndStyle(name string, x, y float64, className, style string) *Chart {
	ch.body = append(ch.body, fmt.Sprintf("    %s:::%s: [%.2f, %.2f] %s", name, className, x, y, style))
	return ch
}

// ClassStyle represents styling options for a class definition.
type ClassStyle struct {
	// Color is the fill color of the point (e.g., "#ff0000").
	Color string
	// Radius is the radius of the point.
	Radius int
	// StrokeColor is the border color of the point.
	StrokeColor string
	// StrokeWidth is the border width of the point (e.g., "10px").
	StrokeWidth string
}

// String returns the class style as a mermaid-compatible string.
func (cs ClassStyle) String() string {
	var parts []string
	if cs.Color != "" {
		parts = append(parts, fmt.Sprintf("color: %s", cs.Color))
	}
	if cs.Radius > 0 {
		parts = append(parts, fmt.Sprintf("radius: %d", cs.Radius))
	}
	if cs.StrokeColor != "" {
		parts = append(parts, fmt.Sprintf("stroke-color: %s", cs.StrokeColor))
	}
	if cs.StrokeWidth != "" {
		parts = append(parts, fmt.Sprintf("stroke-width: %s", cs.StrokeWidth))
	}
	return strings.Join(parts, ", ")
}

// ClassDef defines a class with styling that can be applied to multiple points.
// className is the name of the class.
// style is the styling string (e.g., "color: #ff0000, radius: 10").
func (ch *Chart) ClassDef(className, style string) *Chart {
	ch.body = append(ch.body, fmt.Sprintf("    classDef %s %s", className, style))
	return ch
}

// ClassDefStyled defines a class with ClassStyle struct.
// className is the name of the class.
func (ch *Chart) ClassDefStyled(className string, style ClassStyle) *Chart {
	ch.body = append(ch.body, fmt.Sprintf("    classDef %s %s", className, style.String()))
	return ch
}

// LF adds a line feed to the quadrant chart.
func (ch *Chart) LF() *Chart {
	ch.body = append(ch.body, "")
	return ch
}
