// Package block is mermaid block diagram builder.
//
// Ref. https://mermaid.js.org/syntax/block.html
package block

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// blockLinesCap is the max initial lines in block diagram.
	blockLinesCap int = 2
)

// Direction is the direction for block arrows.
type Direction string

const (
	// DirectionRight points to the right.
	DirectionRight Direction = "right"
	// DirectionLeft points to the left.
	DirectionLeft Direction = "left"
	// DirectionUp points upward.
	DirectionUp Direction = "up"
	// DirectionDown points downward.
	DirectionDown Direction = "down"
	// DirectionX points on X axis.
	DirectionX Direction = "x"
	// DirectionY points on Y axis.
	DirectionY Direction = "y"
	// defaultArrowLabel is the default label for arrows without explicit labels.
	defaultArrowLabel string = "&nbsp;"
)

// Shape is the shape of a node token.
type Shape string

const (
	// ShapeRectangle is the default rectangle shape.
	ShapeRectangle Shape = "rectangle"
	// ShapeRound is a round-edged rectangle shape.
	ShapeRound Shape = "round"
	// ShapeStadium is a stadium shape.
	ShapeStadium Shape = "stadium"
	// ShapeSubroutine is a subroutine shape.
	ShapeSubroutine Shape = "subroutine"
	// ShapeCylinder is a cylindrical shape.
	ShapeCylinder Shape = "cylinder"
	// ShapeCircle is a circle shape.
	ShapeCircle Shape = "circle"
	// ShapeAsymmetric is an asymmetric shape.
	ShapeAsymmetric Shape = "asymmetric"
	// ShapeRhombus is a rhombus shape.
	ShapeRhombus Shape = "rhombus"
	// ShapeHexagon is a hexagon shape.
	ShapeHexagon Shape = "hexagon"
	// ShapeParallelogram is a parallelogram shape.
	ShapeParallelogram Shape = "parallelogram"
	// ShapeParallelogramAlt is a mirrored parallelogram shape.
	ShapeParallelogramAlt Shape = "parallelogramAlt"
	// ShapeTrapezoid is a trapezoid shape.
	ShapeTrapezoid Shape = "trapezoid"
	// ShapeTrapezoidAlt is an inverted trapezoid shape.
	ShapeTrapezoidAlt Shape = "trapezoidAlt"
	// ShapeDoubleCircle is a double-circle shape.
	ShapeDoubleCircle Shape = "doubleCircle"
)

// Diagram is a block diagram builder.
type Diagram struct {
	// body is block diagram body.
	body []string
	// dest is output destination for block diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the block diagram building.
	err error
	// indent is current indentation level for nested blocks.
	indent int
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()
	for _, opt := range opts {
		opt(c)
	}

	lines := make([]string, 0, blockLinesCap)
	lines = append(lines, "block")

	trimmedTitle := strings.TrimSpace(c.title)
	if trimmedTitle != noTitle {
		if containsNewline(trimmedTitle) {
			return &Diagram{
				body:   []string{"block"},
				dest:   w,
				err:    errors.New("title must not contain newline characters"),
				indent: 1,
			}
		}
		lines = append(lines, fmt.Sprintf("    title %s", trimmedTitle))
	}

	return &Diagram{
		body:   lines,
		dest:   w,
		indent: 1,
	}
}

// String returns the block diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the block diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the block diagram body to the output destination.
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

// LF adds a line feed to the block diagram.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}

	d.body = append(d.body, "")
	return d
}

// Columns sets the number of columns for subsequent rows.
func (d *Diagram) Columns(count int) *Diagram {
	if d.err != nil {
		return d
	}
	if count <= 0 {
		d.setError(errors.New("column count must be greater than zero"))
		return d
	}

	d.appendLine(fmt.Sprintf("columns %d", count))
	return d
}

// Token is a row token in block diagrams.
type Token struct {
	value string
	err   error
}

func tokenError(err error) Token {
	return Token{err: err}
}

// Literal returns a raw token for Row.
func Literal(value string) Token {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return tokenError(errors.New("token must not be empty"))
	}
	if containsNewline(trimmed) {
		return tokenError(errors.New("token must not contain newline characters"))
	}
	return Token{value: trimmed}
}

// NodeOption sets node token options.
type NodeOption func(*nodeConfig)

type nodeConfig struct {
	label   string
	shape   Shape
	span    int
	hasSpan bool
}

func newNodeConfig() *nodeConfig {
	return &nodeConfig{
		shape: ShapeRectangle,
	}
}

// WithNodeLabel sets node label.
func WithNodeLabel(label string) NodeOption {
	return func(c *nodeConfig) {
		c.label = label
	}
}

// WithNodeShape sets node shape.
func WithNodeShape(shape Shape) NodeOption {
	return func(c *nodeConfig) {
		c.shape = shape
	}
}

// WithNodeSpan sets node span (column width).
func WithNodeSpan(span int) NodeOption {
	return func(c *nodeConfig) {
		c.span = span
		c.hasSpan = true
	}
}

// Node returns a node token.
func Node(id string, opts ...NodeOption) Token {
	trimmedID, err := validateIdentifier("node id", id)
	if err != nil {
		return tokenError(err)
	}

	c := newNodeConfig()
	for _, opt := range opts {
		opt(c)
	}

	if !isValidShape(c.shape) {
		return tokenError(fmt.Errorf("invalid node shape %q", c.shape))
	}
	if c.hasSpan && c.span <= 0 {
		return tokenError(errors.New("node span must be greater than zero"))
	}

	trimmedLabel := strings.TrimSpace(c.label)
	if containsNewline(trimmedLabel) {
		return tokenError(errors.New("node label must not contain newline characters"))
	}

	label := trimmedLabel
	if label == "" && c.shape != ShapeRectangle {
		label = trimmedID
	}

	token := trimmedID
	if label != "" {
		token = formatNodeToken(trimmedID, label, c.shape)
	}
	if c.hasSpan {
		token = fmt.Sprintf("%s:%d", token, c.span)
	}

	return Token{value: token}
}

// Space returns a space token.
func Space(width ...int) Token {
	if len(width) > 1 {
		return tokenError(errors.New("space accepts zero or one width argument"))
	}
	if len(width) == 0 {
		return Token{value: "space"}
	}
	if width[0] <= 0 {
		return tokenError(errors.New("space width must be greater than zero"))
	}
	if width[0] == 1 {
		return Token{value: "space"}
	}
	return Token{value: fmt.Sprintf("space:%d", width[0])}
}

type arrowConfig struct {
	label     string
	secondary *Direction
}

func newArrowConfig() *arrowConfig {
	return &arrowConfig{
		label: defaultArrowLabel,
	}
}

// ArrowOption sets arrow token options.
type ArrowOption func(*arrowConfig)

// WithArrowLabel sets arrow label text.
func WithArrowLabel(label string) ArrowOption {
	return func(c *arrowConfig) {
		c.label = label
	}
}

// WithArrowSecondaryDirection sets optional secondary arrow direction.
func WithArrowSecondaryDirection(direction Direction) ArrowOption {
	return func(c *arrowConfig) {
		d := direction
		c.secondary = &d
	}
}

// Arrow returns a block arrow token.
func Arrow(id string, direction Direction, opts ...ArrowOption) Token {
	trimmedID, err := validateIdentifier("arrow id", id)
	if err != nil {
		return tokenError(err)
	}
	if !isValidDirection(direction) {
		return tokenError(fmt.Errorf("invalid arrow direction %q", direction))
	}

	c := newArrowConfig()
	for _, opt := range opts {
		opt(c)
	}

	args := string(direction)
	if c.secondary != nil {
		if !isValidDirection(*c.secondary) {
			return tokenError(fmt.Errorf("invalid arrow direction %q", *c.secondary))
		}
		args = fmt.Sprintf("%s, %s", direction, *c.secondary)
	}

	trimmedLabel := strings.TrimSpace(c.label)
	if trimmedLabel == "" {
		return tokenError(errors.New("arrow label must not be empty"))
	}
	if containsNewline(trimmedLabel) {
		return tokenError(errors.New("arrow label must not contain newline characters"))
	}

	return Token{
		value: fmt.Sprintf("%s<[%s]>(%s)", trimmedID, quote(trimmedLabel), args),
	}
}

// ArrowRight returns a right-directed arrow token.
func ArrowRight(id string, opts ...ArrowOption) Token {
	return Arrow(id, DirectionRight, opts...)
}

// ArrowLeft returns a left-directed arrow token.
func ArrowLeft(id string, opts ...ArrowOption) Token {
	return Arrow(id, DirectionLeft, opts...)
}

// ArrowUp returns an up-directed arrow token.
func ArrowUp(id string, opts ...ArrowOption) Token {
	return Arrow(id, DirectionUp, opts...)
}

// ArrowDown returns a down-directed arrow token.
func ArrowDown(id string, opts ...ArrowOption) Token {
	return Arrow(id, DirectionDown, opts...)
}

// ArrowX returns an x-directed arrow token.
func ArrowX(id string, opts ...ArrowOption) Token {
	return Arrow(id, DirectionX, opts...)
}

// ArrowY returns a y-directed arrow token.
func ArrowY(id string, opts ...ArrowOption) Token {
	return Arrow(id, DirectionY, opts...)
}

// Row adds a row with one or more tokens.
func (d *Diagram) Row(tokens ...Token) *Diagram {
	if d.err != nil {
		return d
	}
	if len(tokens) == 0 {
		d.setError(errors.New("row must contain at least one token"))
		return d
	}

	parts := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if token.err != nil {
			d.setError(token.err)
			return d
		}
		trimmed := strings.TrimSpace(token.value)
		if trimmed == "" {
			d.setError(errors.New("token must not be empty"))
			return d
		}
		if containsNewline(trimmed) {
			d.setError(errors.New("token must not contain newline characters"))
			return d
		}
		parts = append(parts, trimmed)
	}

	d.appendLine(strings.Join(parts, " "))
	return d
}

type blockConfig struct {
	id      string
	span    int
	hasSpan bool
}

func newBlockConfig() *blockConfig {
	return &blockConfig{}
}

// BlockOption sets composite block options.
type BlockOption func(*blockConfig)

// WithBlockID sets block id for a composite block.
func WithBlockID(id string) BlockOption {
	return func(c *blockConfig) {
		c.id = id
	}
}

// WithBlockSpan sets block span for a composite block.
func WithBlockSpan(span int) BlockOption {
	return func(c *blockConfig) {
		c.span = span
		c.hasSpan = true
	}
}

// Block adds a composite block and invokes build for nested statements.
func (d *Diagram) Block(build func(*Diagram), opts ...BlockOption) *Diagram {
	if d.err != nil {
		return d
	}
	if build == nil {
		d.setError(errors.New("block builder must not be nil"))
		return d
	}

	c := newBlockConfig()
	for _, opt := range opts {
		opt(c)
	}

	line := "block"
	if strings.TrimSpace(c.id) != "" {
		trimmedID, err := validateIdentifier("block id", c.id)
		if err != nil {
			d.setError(err)
			return d
		}
		line += ":" + trimmedID
	}
	if c.hasSpan {
		if c.span <= 0 {
			d.setError(errors.New("block span must be greater than zero"))
			return d
		}
		line += fmt.Sprintf(":%d", c.span)
	}

	d.appendLine(line)
	d.indent++
	build(d)
	d.indent--
	d.appendLine("end")
	return d
}

// Statement adds a raw statement line.
func (d *Diagram) Statement(statement string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmed := strings.TrimSpace(statement)
	if trimmed == "" {
		d.setError(errors.New("statement must not be empty"))
		return d
	}
	if containsNewline(trimmed) {
		d.setError(errors.New("statement must not contain newline characters"))
		return d
	}

	d.appendLine(trimmed)
	return d
}

// Link adds a link between nodes.
func (d *Diagram) Link(from, to string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedFrom, err := validateIdentifier("source node id", from)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedTo, err := validateIdentifier("destination node id", to)
	if err != nil {
		d.setError(err)
		return d
	}

	d.appendLine(fmt.Sprintf("%s --> %s", trimmedFrom, trimmedTo))
	return d
}

// LinkWithLabel adds a labeled link between nodes.
func (d *Diagram) LinkWithLabel(from, label, to string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedFrom, err := validateIdentifier("source node id", from)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedTo, err := validateIdentifier("destination node id", to)
	if err != nil {
		d.setError(err)
		return d
	}

	trimmedLabel := strings.TrimSpace(label)
	if trimmedLabel == "" {
		d.setError(errors.New("link label must not be empty"))
		return d
	}
	if containsNewline(trimmedLabel) {
		d.setError(errors.New("link label must not contain newline characters"))
		return d
	}

	d.appendLine(fmt.Sprintf("%s -- %s --> %s", trimmedFrom, quote(trimmedLabel), trimmedTo))
	return d
}

// Style adds style settings to one or more nodes.
//
// names can be a comma-separated list.
func (d *Diagram) Style(names, style string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedNames, err := validateText("style target names", names)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedStyle, err := validateText("style", style)
	if err != nil {
		d.setError(err)
		return d
	}

	d.appendLine(fmt.Sprintf("style %s %s", trimmedNames, trimmedStyle))
	return d
}

// ClassDef adds a classDef style declaration.
func (d *Diagram) ClassDef(className, style string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedClassName, err := validateText("class name", className)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedStyle, err := validateText("style", style)
	if err != nil {
		d.setError(err)
		return d
	}

	d.appendLine(fmt.Sprintf("classDef %s %s", trimmedClassName, trimmedStyle))
	return d
}

// Class applies a classDef style to one or more nodes.
//
// names can be a comma-separated list.
func (d *Diagram) Class(names, className string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedNames, err := validateText("node names", names)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedClassName, err := validateText("class name", className)
	if err != nil {
		d.setError(err)
		return d
	}

	d.appendLine(fmt.Sprintf("class %s %s", trimmedNames, trimmedClassName))
	return d
}

func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

func (d *Diagram) appendLine(line string) {
	d.body = append(d.body, fmt.Sprintf("%s%s", strings.Repeat("    ", d.indent), line))
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

func containsNewline(value string) bool {
	return strings.ContainsAny(value, "\n\r")
}

func isValidDirection(direction Direction) bool {
	switch direction {
	case DirectionRight, DirectionLeft, DirectionUp, DirectionDown, DirectionX, DirectionY:
		return true
	default:
		return false
	}
}

func isValidShape(shape Shape) bool {
	switch shape {
	case ShapeRectangle,
		ShapeRound,
		ShapeStadium,
		ShapeSubroutine,
		ShapeCylinder,
		ShapeCircle,
		ShapeAsymmetric,
		ShapeRhombus,
		ShapeHexagon,
		ShapeParallelogram,
		ShapeParallelogramAlt,
		ShapeTrapezoid,
		ShapeTrapezoidAlt,
		ShapeDoubleCircle:
		return true
	default:
		return false
	}
}

func formatNodeToken(id, label string, shape Shape) string {
	quoted := quote(label)

	switch shape {
	case ShapeRound:
		return fmt.Sprintf("%s(%s)", id, quoted)
	case ShapeStadium:
		return fmt.Sprintf("%s([%s])", id, quoted)
	case ShapeSubroutine:
		return fmt.Sprintf("%s[[%s]]", id, quoted)
	case ShapeCylinder:
		return fmt.Sprintf("%s[(%s)]", id, quoted)
	case ShapeCircle:
		return fmt.Sprintf("%s((%s))", id, quoted)
	case ShapeAsymmetric:
		return fmt.Sprintf("%s>%s]", id, quoted)
	case ShapeRhombus:
		return fmt.Sprintf("%s{%s}", id, quoted)
	case ShapeHexagon:
		return fmt.Sprintf("%s{{%s}}", id, quoted)
	case ShapeParallelogram:
		return fmt.Sprintf("%s[/%s/]", id, quoted)
	case ShapeParallelogramAlt:
		return fmt.Sprintf("%s[\\%s\\]", id, quoted)
	case ShapeTrapezoid:
		return fmt.Sprintf("%s[/%s\\]", id, quoted)
	case ShapeTrapezoidAlt:
		return fmt.Sprintf("%s[\\%s/]", id, quoted)
	case ShapeDoubleCircle:
		return fmt.Sprintf("%s(((%s)))", id, quoted)
	case ShapeRectangle:
		fallthrough
	default:
		return fmt.Sprintf("%s[%s]", id, quoted)
	}
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
