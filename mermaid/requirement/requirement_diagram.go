// Package requirement is mermaid requirement diagram builder.
//
// Ref. https://mermaid.js.org/syntax/requirementDiagram.html
package requirement

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/nao1215/markdown/internal"
)

const (
	// requirementLinesCap is the max initial lines with title frontmatter.
	requirementLinesCap int = 4
)

// Direction is the diagram direction.
type Direction string

const (
	// DirectionTB is top to bottom direction.
	DirectionTB Direction = "TB"
	// DirectionBT is bottom to top direction.
	DirectionBT Direction = "BT"
	// DirectionLR is left to right direction.
	DirectionLR Direction = "LR"
	// DirectionRL is right to left direction.
	DirectionRL Direction = "RL"
)

// RequirementType is the requirement type keyword.
type RequirementType string

const (
	// RequirementTypeRequirement represents a generic requirement.
	RequirementTypeRequirement RequirementType = "requirement"
	// RequirementTypeFunctional represents a functional requirement.
	RequirementTypeFunctional RequirementType = "functionalRequirement"
	// RequirementTypeInterface represents an interface requirement.
	RequirementTypeInterface RequirementType = "interfaceRequirement"
	// RequirementTypePerformance represents a performance requirement.
	RequirementTypePerformance RequirementType = "performanceRequirement"
	// RequirementTypePhysical represents a physical requirement.
	RequirementTypePhysical RequirementType = "physicalRequirement"
	// RequirementTypeDesignConstraint represents a design constraint.
	RequirementTypeDesignConstraint RequirementType = "designConstraint"
)

// Risk is the risk level used in requirement blocks.
type Risk string

const (
	// RiskLow is low risk.
	RiskLow Risk = "Low"
	// RiskMedium is medium risk.
	RiskMedium Risk = "Medium"
	// RiskHigh is high risk.
	RiskHigh Risk = "High"
)

// VerifyMethod is the verify method used in requirement blocks.
type VerifyMethod string

const (
	// VerifyMethodAnalysis means analysis.
	VerifyMethodAnalysis VerifyMethod = "Analysis"
	// VerifyMethodInspection means inspection.
	VerifyMethodInspection VerifyMethod = "Inspection"
	// VerifyMethodTest means test.
	VerifyMethodTest VerifyMethod = "Test"
	// VerifyMethodDemonstration means demonstration.
	VerifyMethodDemonstration VerifyMethod = "Demonstration"
)

// Relationship is the relationship type in requirement diagrams.
type Relationship string

const (
	// RelationshipContains represents "contains".
	RelationshipContains Relationship = "contains"
	// RelationshipCopies represents "copies".
	RelationshipCopies Relationship = "copies"
	// RelationshipDerives represents "derives".
	RelationshipDerives Relationship = "derives"
	// RelationshipSatisfies represents "satisfies".
	RelationshipSatisfies Relationship = "satisfies"
	// RelationshipVerifies represents "verifies".
	RelationshipVerifies Relationship = "verifies"
	// RelationshipRefines represents "refines".
	RelationshipRefines Relationship = "refines"
	// RelationshipTraces represents "traces".
	RelationshipTraces Relationship = "traces"
)

// Diagram is a requirement diagram builder.
type Diagram struct {
	// body is requirement diagram body.
	body []string
	// dest is output destination for requirement diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the requirement diagram building.
	err error
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := make([]string, 0, requirementLinesCap)
	title := strings.TrimSpace(c.title)
	if title != noTitle {
		if containsNewline(title) {
			return &Diagram{
				body: []string{"requirementDiagram"},
				dest: w,
				err:  errors.New("title must not contain newline characters"),
			}
		}
		lines = append(lines, "---")
		lines = append(lines, fmt.Sprintf("title: %s", title))
		lines = append(lines, "---")
	}
	lines = append(lines, "requirementDiagram")

	return &Diagram{
		body: lines,
		dest: w,
	}
}

// String returns the requirement diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the requirement diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the requirement diagram body to the output destination.
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

// LF adds a line feed to the requirement diagram.
func (d *Diagram) LF() *Diagram {
	if d.err != nil {
		return d
	}

	d.body = append(d.body, "")
	return d
}

// SetDirection sets the direction of the requirement diagram.
func (d *Diagram) SetDirection(dir Direction) *Diagram {
	if d.err != nil {
		return d
	}
	if !isValidDirection(dir) {
		d.setError(fmt.Errorf("invalid direction %q", dir))
		return d
	}

	d.body = append(d.body, fmt.Sprintf("    direction %s", dir))
	return d
}

type requirementConfig struct {
	id           string
	text         string
	risk         Risk
	verifyMethod VerifyMethod
	classNames   []string
}

func newRequirementConfig() *requirementConfig {
	return &requirementConfig{}
}

// RequirementOption sets requirement options.
type RequirementOption func(*requirementConfig)

// WithID sets id property for a requirement.
func WithID(id string) RequirementOption {
	return func(c *requirementConfig) {
		c.id = id
	}
}

// WithText sets text property for a requirement.
func WithText(text string) RequirementOption {
	return func(c *requirementConfig) {
		c.text = text
	}
}

// WithRisk sets risk property for a requirement.
func WithRisk(risk Risk) RequirementOption {
	return func(c *requirementConfig) {
		c.risk = risk
	}
}

// WithVerifyMethod sets verifymethod property for a requirement.
func WithVerifyMethod(method VerifyMethod) RequirementOption {
	return func(c *requirementConfig) {
		c.verifyMethod = method
	}
}

// WithRequirementClasses applies Mermaid class shorthand to a requirement.
func WithRequirementClasses(classNames ...string) RequirementOption {
	return func(c *requirementConfig) {
		c.classNames = append(c.classNames, classNames...)
	}
}

// Requirement adds a requirement block with type "requirement".
//
// WithID, WithText, WithRisk, and WithVerifyMethod are required;
// omitting any of them records an error.
func (d *Diagram) Requirement(name string, opts ...RequirementOption) *Diagram {
	return d.RequirementOfType(RequirementTypeRequirement, name, opts...)
}

// RequirementOfType adds a requirement block with an explicit type.
func (d *Diagram) RequirementOfType(
	reqType RequirementType,
	name string,
	opts ...RequirementOption,
) *Diagram {
	if d.err != nil {
		return d
	}

	normalizedType, ok := normalizeRequirementType(reqType)
	if !ok {
		d.setError(fmt.Errorf("invalid requirement type %q", reqType))
		return d
	}

	trimmedName, err := validateName("requirement name", name)
	if err != nil {
		d.setError(err)
		return d
	}

	c := newRequirementConfig()
	for _, opt := range opts {
		opt(c)
	}

	trimmedID, err := validateName("requirement id", c.id)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedText, err := validateName("requirement text", c.text)
	if err != nil {
		d.setError(err)
		return d
	}

	normalizedRisk, ok := normalizeRisk(c.risk)
	if !ok {
		d.setError(fmt.Errorf("invalid risk %q", c.risk))
		return d
	}
	normalizedVerifyMethod, ok := normalizeVerifyMethod(c.verifyMethod)
	if !ok {
		d.setError(fmt.Errorf("invalid verify method %q", c.verifyMethod))
		return d
	}

	classShorthand, err := formatClassShorthand(c.classNames...)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(
		d.body,
		fmt.Sprintf("    %s %s%s {", normalizedType, formatName(trimmedName), classShorthand),
	)
	d.body = append(d.body, fmt.Sprintf("        id: %s", quote(normalizeQuoted(trimmedID))))
	d.body = append(d.body, fmt.Sprintf("        text: %s", quote(normalizeQuoted(trimmedText))))
	d.body = append(d.body, fmt.Sprintf("        risk: %s", normalizedRisk))
	d.body = append(d.body, fmt.Sprintf("        verifymethod: %s", normalizedVerifyMethod))
	d.body = append(d.body, "    }")
	return d
}

// FunctionalRequirement adds a requirement block with type "functionalRequirement".
func (d *Diagram) FunctionalRequirement(name string, opts ...RequirementOption) *Diagram {
	return d.RequirementOfType(RequirementTypeFunctional, name, opts...)
}

// InterfaceRequirement adds a requirement block with type "interfaceRequirement".
func (d *Diagram) InterfaceRequirement(name string, opts ...RequirementOption) *Diagram {
	return d.RequirementOfType(RequirementTypeInterface, name, opts...)
}

// PerformanceRequirement adds a requirement block with type "performanceRequirement".
func (d *Diagram) PerformanceRequirement(name string, opts ...RequirementOption) *Diagram {
	return d.RequirementOfType(RequirementTypePerformance, name, opts...)
}

// PhysicalRequirement adds a requirement block with type "physicalRequirement".
func (d *Diagram) PhysicalRequirement(name string, opts ...RequirementOption) *Diagram {
	return d.RequirementOfType(RequirementTypePhysical, name, opts...)
}

// DesignConstraint adds a requirement block with type "designConstraint".
func (d *Diagram) DesignConstraint(name string, opts ...RequirementOption) *Diagram {
	return d.RequirementOfType(RequirementTypeDesignConstraint, name, opts...)
}

type elementConfig struct {
	eType      string
	docRef     string
	classNames []string
}

func newElementConfig() *elementConfig {
	return &elementConfig{}
}

// ElementOption sets element options.
type ElementOption func(*elementConfig)

// WithElementType sets type property for an element.
func WithElementType(elementType string) ElementOption {
	return func(c *elementConfig) {
		c.eType = elementType
	}
}

// WithDocRef sets docRef property for an element.
func WithDocRef(docRef string) ElementOption {
	return func(c *elementConfig) {
		c.docRef = docRef
	}
}

// WithElementClasses applies Mermaid class shorthand to an element.
func WithElementClasses(classNames ...string) ElementOption {
	return func(c *elementConfig) {
		c.classNames = append(c.classNames, classNames...)
	}
}

// Element adds an element block.
func (d *Diagram) Element(name string, opts ...ElementOption) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedName, err := validateName("element name", name)
	if err != nil {
		d.setError(err)
		return d
	}

	c := newElementConfig()
	for _, opt := range opts {
		opt(c)
	}

	trimmedType, err := validateOptionalValue("element type", c.eType)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedDocRef, err := validateOptionalValue("element docRef", c.docRef)
	if err != nil {
		d.setError(err)
		return d
	}

	classShorthand, err := formatClassShorthand(c.classNames...)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(
		d.body,
		fmt.Sprintf("    element %s%s {", formatName(trimmedName), classShorthand),
	)
	if trimmedType != "" {
		d.body = append(d.body, fmt.Sprintf("        type: %s", quote(normalizeQuoted(trimmedType))))
	}
	if trimmedDocRef != "" {
		d.body = append(d.body, fmt.Sprintf("        docRef: %s", quote(normalizeQuoted(trimmedDocRef))))
	}
	d.body = append(d.body, "    }")
	return d
}

// Relation adds a relationship between requirement diagram nodes.
func (d *Diagram) Relation(from string, relationship Relationship, to string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedFrom, err := validateName("source name", from)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedTo, err := validateName("destination name", to)
	if err != nil {
		d.setError(err)
		return d
	}

	normalizedRelationship, ok := normalizeRelationship(relationship)
	if !ok {
		d.setError(fmt.Errorf("invalid relationship %q", relationship))
		return d
	}

	d.body = append(
		d.body,
		fmt.Sprintf(
			"    %s - %s -> %s",
			formatName(trimmedFrom),
			normalizedRelationship,
			formatName(trimmedTo),
		),
	)
	return d
}

// From returns a builder that fixes relation source to the specified node.
func (d *Diagram) From(from string) *SourceRelationBuilder {
	trimmedFrom := strings.TrimSpace(from)
	if d.err == nil {
		var err error
		trimmedFrom, err = validateName("source name", from)
		if err != nil {
			d.setError(err)
		}
	}

	return &SourceRelationBuilder{
		Diagram: d,
		from:    trimmedFrom,
	}
}

// SourceRelationBuilder builds relationships from a fixed source node.
type SourceRelationBuilder struct {
	*Diagram
	from string
}

// Relation adds a relationship from the fixed source node to a target node.
func (b *SourceRelationBuilder) Relation(
	relationship Relationship,
	to string,
) *SourceRelationBuilder {
	b.Diagram.Relation(b.from, relationship, to)
	return b
}

// Contains adds a "contains" relationship from fixed source node.
func (b *SourceRelationBuilder) Contains(to string) *SourceRelationBuilder {
	return b.Relation(RelationshipContains, to)
}

// Copies adds a "copies" relationship from fixed source node.
func (b *SourceRelationBuilder) Copies(to string) *SourceRelationBuilder {
	return b.Relation(RelationshipCopies, to)
}

// Derives adds a "derives" relationship from fixed source node.
func (b *SourceRelationBuilder) Derives(to string) *SourceRelationBuilder {
	return b.Relation(RelationshipDerives, to)
}

// Satisfies adds a "satisfies" relationship from fixed source node.
func (b *SourceRelationBuilder) Satisfies(to string) *SourceRelationBuilder {
	return b.Relation(RelationshipSatisfies, to)
}

// Verifies adds a "verifies" relationship from fixed source node.
func (b *SourceRelationBuilder) Verifies(to string) *SourceRelationBuilder {
	return b.Relation(RelationshipVerifies, to)
}

// Refines adds a "refines" relationship from fixed source node.
func (b *SourceRelationBuilder) Refines(to string) *SourceRelationBuilder {
	return b.Relation(RelationshipRefines, to)
}

// Traces adds a "traces" relationship from fixed source node.
func (b *SourceRelationBuilder) Traces(to string) *SourceRelationBuilder {
	return b.Relation(RelationshipTraces, to)
}

// Contains adds a "contains" relationship.
func (d *Diagram) Contains(from, to string) *Diagram {
	return d.Relation(from, RelationshipContains, to)
}

// Copies adds a "copies" relationship.
func (d *Diagram) Copies(from, to string) *Diagram {
	return d.Relation(from, RelationshipCopies, to)
}

// Derives adds a "derives" relationship.
func (d *Diagram) Derives(from, to string) *Diagram {
	return d.Relation(from, RelationshipDerives, to)
}

// Satisfies adds a "satisfies" relationship.
func (d *Diagram) Satisfies(from, to string) *Diagram {
	return d.Relation(from, RelationshipSatisfies, to)
}

// Verifies adds a "verifies" relationship.
func (d *Diagram) Verifies(from, to string) *Diagram {
	return d.Relation(from, RelationshipVerifies, to)
}

// Refines adds a "refines" relationship.
func (d *Diagram) Refines(from, to string) *Diagram {
	return d.Relation(from, RelationshipRefines, to)
}

// Traces adds a "traces" relationship.
func (d *Diagram) Traces(from, to string) *Diagram {
	return d.Relation(from, RelationshipTraces, to)
}

// Style adds style settings to one or more nodes.
//
// names can be a comma-separated list.
func (d *Diagram) Style(names, style string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedNames, err := validateName("style target names", names)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedStyle, err := validateName("style", style)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, fmt.Sprintf("    style %s %s", trimmedNames, trimmedStyle))
	return d
}

// ClassDef adds a classDef style declaration.
func (d *Diagram) ClassDef(classNames, style string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedClassNames, err := validateName("class names", classNames)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedStyle, err := validateName("style", style)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, fmt.Sprintf("    classDef %s %s", trimmedClassNames, trimmedStyle))
	return d
}

// ClassDefSpec is a style declaration input for ClassDefs.
type ClassDefSpec struct {
	ClassNames string
	Style      string
}

// Def returns a ClassDefSpec for use with ClassDefs.
func Def(classNames, style string) ClassDefSpec {
	return ClassDefSpec{
		ClassNames: classNames,
		Style:      style,
	}
}

// ClassDefs adds multiple classDef style declarations at once.
func (d *Diagram) ClassDefs(defs ...ClassDefSpec) *Diagram {
	if d.err != nil {
		return d
	}
	if len(defs) == 0 {
		d.setError(errors.New("at least one classDef is required"))
		return d
	}

	for _, def := range defs {
		d.ClassDef(def.ClassNames, def.Style)
	}
	return d
}

// Class applies a classDef style to one or more nodes.
//
// names can be a comma-separated list.
func (d *Diagram) Class(names, classNames string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedNames, err := validateName("node names", names)
	if err != nil {
		d.setError(err)
		return d
	}
	trimmedClassNames, err := validateName("class names", classNames)
	if err != nil {
		d.setError(err)
		return d
	}

	d.body = append(d.body, fmt.Sprintf("    class %s %s", trimmedNames, trimmedClassNames))
	return d
}

// ClassShorthand applies classDef style using shorthand syntax.
func (d *Diagram) ClassShorthand(name string, classNames ...string) *Diagram {
	if d.err != nil {
		return d
	}

	trimmedName, err := validateName("node name", name)
	if err != nil {
		d.setError(err)
		return d
	}
	if len(classNames) == 0 {
		d.setError(errors.New("at least one class name is required"))
		return d
	}

	classShorthand, err := formatClassShorthand(classNames...)
	if err != nil {
		d.setError(err)
		return d
	}
	d.body = append(d.body, fmt.Sprintf("    %s%s", formatName(trimmedName), classShorthand))
	return d
}

func (d *Diagram) setError(err error) {
	if d.err == nil {
		d.err = err
	}
}

func isValidDirection(dir Direction) bool {
	switch dir {
	case DirectionTB, DirectionBT, DirectionLR, DirectionRL:
		return true
	default:
		return false
	}
}

func normalizeRequirementType(reqType RequirementType) (string, bool) {
	switch strings.ToLower(string(reqType)) {
	case strings.ToLower(string(RequirementTypeRequirement)):
		return string(RequirementTypeRequirement), true
	case strings.ToLower(string(RequirementTypeFunctional)):
		return string(RequirementTypeFunctional), true
	case strings.ToLower(string(RequirementTypeInterface)):
		return string(RequirementTypeInterface), true
	case strings.ToLower(string(RequirementTypePerformance)):
		return string(RequirementTypePerformance), true
	case strings.ToLower(string(RequirementTypePhysical)):
		return string(RequirementTypePhysical), true
	case strings.ToLower(string(RequirementTypeDesignConstraint)):
		return string(RequirementTypeDesignConstraint), true
	default:
		return "", false
	}
}

func normalizeRisk(risk Risk) (string, bool) {
	switch strings.ToLower(string(risk)) {
	case "low":
		return "Low", true
	case "medium":
		return "Medium", true
	case "high":
		return "High", true
	default:
		return "", false
	}
}

func normalizeVerifyMethod(method VerifyMethod) (string, bool) {
	switch strings.ToLower(string(method)) {
	case "analysis":
		return "Analysis", true
	case "inspection":
		return "Inspection", true
	case "test":
		return "Test", true
	case "demonstration":
		return "Demonstration", true
	default:
		return "", false
	}
}

func normalizeRelationship(relationship Relationship) (string, bool) {
	switch strings.ToLower(string(relationship)) {
	case strings.ToLower(string(RelationshipContains)):
		return string(RelationshipContains), true
	case strings.ToLower(string(RelationshipCopies)):
		return string(RelationshipCopies), true
	case strings.ToLower(string(RelationshipDerives)):
		return string(RelationshipDerives), true
	case strings.ToLower(string(RelationshipSatisfies)):
		return string(RelationshipSatisfies), true
	case strings.ToLower(string(RelationshipVerifies)):
		return string(RelationshipVerifies), true
	case strings.ToLower(string(RelationshipRefines)):
		return string(RelationshipRefines), true
	case strings.ToLower(string(RelationshipTraces)):
		return string(RelationshipTraces), true
	default:
		return "", false
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

func validateOptionalValue(fieldName, value string) (string, error) {
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

func formatClassShorthand(classNames ...string) (string, error) {
	if len(classNames) == 0 {
		return "", nil
	}

	var b strings.Builder
	for _, className := range classNames {
		trimmed, err := validateName("class name", className)
		if err != nil {
			return "", err
		}
		b.WriteString(":::")
		b.WriteString(trimmed)
	}
	return b.String(), nil
}

func formatName(name string) string {
	trimmed := normalizeQuoted(strings.TrimSpace(name))
	if isMermaidKeyword(trimmed) || shouldQuote(trimmed) {
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
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func isMermaidKeyword(value string) bool {
	switch strings.ToLower(value) {
	case "requirementdiagram",
		"direction",
		strings.ToLower(string(RequirementTypeRequirement)),
		strings.ToLower(string(RequirementTypeFunctional)),
		strings.ToLower(string(RequirementTypeInterface)),
		strings.ToLower(string(RequirementTypePerformance)),
		strings.ToLower(string(RequirementTypePhysical)),
		strings.ToLower(string(RequirementTypeDesignConstraint)),
		"element",
		"classdef",
		"class",
		"style",
		"id",
		"text",
		"risk",
		"verifymethod",
		"type",
		"docref",
		strings.ToLower(string(RelationshipContains)),
		strings.ToLower(string(RelationshipCopies)),
		strings.ToLower(string(RelationshipDerives)),
		strings.ToLower(string(RelationshipSatisfies)),
		strings.ToLower(string(RelationshipVerifies)),
		strings.ToLower(string(RelationshipRefines)),
		strings.ToLower(string(RelationshipTraces)):
		return true
	default:
		return false
	}
}

func quote(value string) string {
	escaped := strings.ReplaceAll(value, `\`, "&#92;")
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
