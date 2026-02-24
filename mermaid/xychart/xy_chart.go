// Package xychart is mermaid XY chart builder.
//
// Ref. https://mermaid.js.org/syntax/xyChart.html
package xychart

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/nao1215/markdown/internal"
)

const (
	// xychartLinesCap is the max initial lines in xychart.
	xychartLinesCap int = 4
)

// Orientation is the chart orientation.
type Orientation string

const (
	// OrientationVertical is default top-to-bottom orientation.
	OrientationVertical Orientation = "vertical"
	// OrientationHorizontal is left-to-right orientation.
	OrientationHorizontal Orientation = "horizontal"
)

// Diagram is an XY chart builder.
type Diagram struct {
	// body is xychart body.
	body []string
	// dest is output destination for xychart body.
	dest io.Writer
	// err manages errors that occur in all parts of the xychart building.
	err error
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	if !isValidOrientation(c.orientation) {
		return &Diagram{
			body: []string{"xychart"},
			dest: w,
			err:  fmt.Errorf("invalid orientation %q", c.orientation),
		}
	}

	lines := make([]string, 0, xychartLinesCap)
	base := "xychart"
	if c.orientation == OrientationHorizontal {
		base += " horizontal"
	}
	lines = append(lines, base)

	trimmedTitle := strings.TrimSpace(c.title)
	if trimmedTitle != noTitle {
		if containsNewline(trimmedTitle) {
			return &Diagram{
				body: []string{base},
				dest: w,
				err:  errors.New("title must not contain newline characters"),
			}
		}
		lines = append(lines, fmt.Sprintf("    title %s", quote(trimmedTitle)))
	}

	return &Diagram{
		body: lines,
		dest: w,
	}
}

// String returns the xychart body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the xychart building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the xychart body to the output destination.
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

// LF adds a line feed to the xychart.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}

	d.body = append(d.body, "")
	return d
}

// XAxisLabels sets categorical labels for x-axis.
func (d *Diagram) XAxisLabels(labels ...string) *Diagram {
	return d.xAxisLabels("", labels...)
}

// XAxisLabelsWithTitle sets x-axis title and categorical labels.
func (d *Diagram) XAxisLabelsWithTitle(title string, labels ...string) *Diagram {
	return d.xAxisLabels(title, labels...)
}

func (d *Diagram) xAxisLabels(title string, labels ...string) *Diagram {
	if d.err != nil {
		return d
	}

	if len(labels) == 0 {
		d.setError(errors.New("x-axis labels must not be empty"))
		return d
	}

	labelList, err := formatLabelList(labels...)
	if err != nil {
		d.setError(err)
		return d
	}

	trimmedTitle, err := validateOptionalText("x-axis title", title)
	if err != nil {
		d.setError(err)
		return d
	}

	line := "    x-axis"
	if trimmedTitle != "" {
		line += " " + formatTextToken(trimmedTitle)
	}
	line += " [" + labelList + "]"

	d.body = append(d.body, line)
	return d
}

// XAxisRange sets numeric range for x-axis.
func (d *Diagram) XAxisRange(min, max float64) *Diagram {
	return d.xAxisRange("", min, max)
}

// XAxisRangeWithTitle sets x-axis title and numeric range.
func (d *Diagram) XAxisRangeWithTitle(title string, min, max float64) *Diagram {
	return d.xAxisRange(title, min, max)
}

func (d *Diagram) xAxisRange(title string, min, max float64) *Diagram {
	if d.err != nil {
		return d
	}

	if min >= max {
		d.setError(errors.New("x-axis range requires min to be less than max"))
		return d
	}

	trimmedTitle, err := validateOptionalText("x-axis title", title)
	if err != nil {
		d.setError(err)
		return d
	}

	line := "    x-axis"
	if trimmedTitle != "" {
		line += " " + formatTextToken(trimmedTitle)
	}
	line += " " + formatNumber(min) + " --> " + formatNumber(max)

	d.body = append(d.body, line)
	return d
}

// YAxisRange sets numeric range for y-axis.
func (d *Diagram) YAxisRange(min, max float64) *Diagram {
	return d.yAxisRange("", min, max)
}

// YAxisRangeWithTitle sets y-axis title and numeric range.
func (d *Diagram) YAxisRangeWithTitle(title string, min, max float64) *Diagram {
	return d.yAxisRange(title, min, max)
}

func (d *Diagram) yAxisRange(title string, min, max float64) *Diagram {
	if d.err != nil {
		return d
	}

	if min >= max {
		d.setError(errors.New("y-axis range requires min to be less than max"))
		return d
	}

	trimmedTitle, err := validateOptionalText("y-axis title", title)
	if err != nil {
		d.setError(err)
		return d
	}

	line := "    y-axis"
	if trimmedTitle != "" {
		line += " " + formatTextToken(trimmedTitle)
	}
	line += " " + formatNumber(min) + " --> " + formatNumber(max)

	d.body = append(d.body, line)
	return d
}

// Bar adds a bar series.
func (d *Diagram) Bar(values ...float64) *Diagram {
	return d.series("bar", values...)
}

// Line adds a line series.
func (d *Diagram) Line(values ...float64) *Diagram {
	return d.series("line", values...)
}

func (d *Diagram) series(seriesType string, values ...float64) *Diagram {
	if d.err != nil {
		return d
	}

	if len(values) == 0 {
		d.setError(fmt.Errorf("%s values must not be empty", seriesType))
		return d
	}

	formatted := make([]string, 0, len(values))
	for _, value := range values {
		formatted = append(formatted, formatNumber(value))
	}

	d.body = append(d.body, fmt.Sprintf("    %s [%s]", seriesType, strings.Join(formatted, ", ")))
	return d
}

func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

func isValidOrientation(orientation Orientation) bool {
	switch orientation {
	case OrientationVertical, OrientationHorizontal:
		return true
	default:
		return false
	}
}

func formatLabelList(labels ...string) (string, error) {
	formatted := make([]string, 0, len(labels))
	for _, label := range labels {
		trimmed, err := validateText("x-axis label", label)
		if err != nil {
			return "", err
		}
		formatted = append(formatted, formatTextToken(trimmed))
	}
	return strings.Join(formatted, ", "), nil
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

func validateOptionalText(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", nil
	}
	if containsNewline(trimmed) {
		return "", fmt.Errorf("%s must not contain newline characters", fieldName)
	}
	return trimmed, nil
}

func containsNewline(value string) bool {
	return strings.ContainsAny(value, "\n\r")
}

func formatTextToken(value string) string {
	trimmed := normalizeQuoted(strings.TrimSpace(value))
	if shouldQuote(trimmed) || isKeyword(trimmed) {
		return quote(trimmed)
	}
	return trimmed
}

func shouldQuote(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if !isIdentifierChar(r) {
			return true
		}
	}
	return false
}

func isIdentifierChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.'
}

func isKeyword(value string) bool {
	switch strings.ToLower(value) {
	case "xychart", "horizontal", "title", "x-axis", "y-axis", "bar", "line":
		return true
	default:
		return false
	}
}

func formatNumber(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func quote(value string) string {
	escaped := normalizeQuoted(value)
	escaped = strings.ReplaceAll(escaped, `\`, "&#92;")
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
