// Package class is mermaid class diagram builder.
//
// Ref. https://mermaid.js.org/syntax/classDiagram.html
package class

import (
	"fmt"
	"io"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// Diagram is a class diagram builder.
type Diagram struct {
	// body is class diagram body.
	body []string
	// config is the configuration for the class diagram.
	config *config
	// dest is output destination for class diagram body.
	dest io.Writer
	// err manages errors that occur in all parts of the class diagram building.
	err error
}

// NewDiagram returns a new Diagram.
func NewDiagram(w io.Writer, opts ...Option) *Diagram {
	c := newConfig()

	for _, opt := range opts {
		opt(c)
	}

	lines := []string{}
	if c.title != noTitle {
		lines = append(lines, "---")
		lines = append(lines, fmt.Sprintf("title: %s", c.title))
		lines = append(lines, "---")
	}
	lines = append(lines, "classDiagram")

	return &Diagram{
		body:   lines,
		dest:   w,
		config: c,
	}
}

// String returns the class diagram body.
func (d *Diagram) String() string {
	return strings.Join(d.body, internal.LineFeed())
}

// Error returns the error that occurred during the class diagram building.
func (d *Diagram) Error() error {
	return d.err
}

// Build writes the class diagram body to the output destination.
func (d *Diagram) Build() error {
	if _, err := fmt.Fprint(d.dest, d.String()); err != nil {
		d.err = fmt.Errorf("failed to write: %w", err)
		return d.err
	}
	return nil
}

// LF adds a line feed to the class diagram.
func (d *Diagram) LF() *Diagram {
	d.body = append(d.body, "")
	return d
}

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

// SetDirection sets the direction of the class diagram.
func (d *Diagram) SetDirection(dir Direction) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    direction %s", dir))
	return d
}

// Visibility is the member visibility marker in class diagrams.
type Visibility string

const (
	// VisibilityPublic is public member visibility.
	VisibilityPublic Visibility = "+"
	// VisibilityPrivate is private member visibility.
	VisibilityPrivate Visibility = "-"
	// VisibilityProtected is protected member visibility.
	VisibilityProtected Visibility = "#"
	// VisibilityPackage is package/internal member visibility.
	VisibilityPackage Visibility = "~"
)

// ClassMemberOption sets class members for class block syntax.
type ClassMemberOption func(*classMemberConfig)

type classMemberConfig struct {
	members []string
}

// Cardinality is the cardinality text used in class relationships.
type Cardinality string

const (
	// CardinalityOne means exactly one.
	CardinalityOne Cardinality = "1"
	// CardinalityMany means many.
	CardinalityMany Cardinality = "many"
	// CardinalityZeroOrOne means zero or one.
	CardinalityZeroOrOne Cardinality = "0..1"
	// CardinalityZeroOrMore means zero or more.
	CardinalityZeroOrMore Cardinality = "0..*"
	// CardinalityOneOrMore means one or more.
	CardinalityOneOrMore Cardinality = "1..*"
)

// RelationOption sets options for relationship sugar methods.
type RelationOption func(*relationConfig)

type relationConfig struct {
	fromCardinality string
	toCardinality   string
	label           string
}

// WithCardinality sets cardinality for relationship sugar methods.
func WithCardinality(from, to Cardinality) RelationOption {
	return func(c *relationConfig) {
		c.fromCardinality = string(from)
		c.toCardinality = string(to)
	}
}

// WithOneToMany sets one-to-many cardinality.
func WithOneToMany() RelationOption {
	return WithCardinality(CardinalityOne, CardinalityMany)
}

// WithManyToOne sets many-to-one cardinality.
func WithManyToOne() RelationOption {
	return WithCardinality(CardinalityMany, CardinalityOne)
}

// WithOneToOne sets one-to-one cardinality.
func WithOneToOne() RelationOption {
	return WithCardinality(CardinalityOne, CardinalityOne)
}

// WithManyToMany sets many-to-many cardinality.
func WithManyToMany() RelationOption {
	return WithCardinality(CardinalityMany, CardinalityMany)
}

// WithRelationLabel sets label for relationship sugar methods.
func WithRelationLabel(label string) RelationOption {
	return func(c *relationConfig) {
		c.label = label
	}
}

// WithField adds a field member option.
func WithField(visibility Visibility, fieldType, fieldName string) ClassMemberOption {
	return func(c *classMemberConfig) {
		c.members = append(c.members, formatFieldMember(visibility, fieldType, fieldName))
	}
}

// WithPublicField adds a public field member option.
func WithPublicField(fieldType, fieldName string) ClassMemberOption {
	return WithField(VisibilityPublic, fieldType, fieldName)
}

// WithPrivateField adds a private field member option.
func WithPrivateField(fieldType, fieldName string) ClassMemberOption {
	return WithField(VisibilityPrivate, fieldType, fieldName)
}

// WithProtectedField adds a protected field member option.
func WithProtectedField(fieldType, fieldName string) ClassMemberOption {
	return WithField(VisibilityProtected, fieldType, fieldName)
}

// WithPackageField adds a package/internal field member option.
func WithPackageField(fieldType, fieldName string) ClassMemberOption {
	return WithField(VisibilityPackage, fieldType, fieldName)
}

// WithMethod adds a method member option.
func WithMethod(
	visibility Visibility,
	methodName,
	returnType string,
	params ...string,
) ClassMemberOption {
	return func(c *classMemberConfig) {
		c.members = append(c.members, formatMethodMember(visibility, methodName, returnType, params...))
	}
}

// WithPublicMethod adds a public method member option.
func WithPublicMethod(methodName, returnType string, params ...string) ClassMemberOption {
	return WithMethod(VisibilityPublic, methodName, returnType, params...)
}

// WithPrivateMethod adds a private method member option.
func WithPrivateMethod(methodName, returnType string, params ...string) ClassMemberOption {
	return WithMethod(VisibilityPrivate, methodName, returnType, params...)
}

// WithProtectedMethod adds a protected method member option.
func WithProtectedMethod(methodName, returnType string, params ...string) ClassMemberOption {
	return WithMethod(VisibilityProtected, methodName, returnType, params...)
}

// WithPackageMethod adds a package/internal method member option.
func WithPackageMethod(methodName, returnType string, params ...string) ClassMemberOption {
	return WithMethod(VisibilityPackage, methodName, returnType, params...)
}

// Comment adds a Mermaid comment line.
func (d *Diagram) Comment(comment string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %%%% %s", comment))
	return d
}

// Class adds a class declaration.
// If member options are provided, Class generates class block syntax.
func (d *Diagram) Class(name string, opts ...ClassMemberOption) *Diagram {
	if len(opts) == 0 {
		d.body = append(d.body, fmt.Sprintf("    class %s", name))
		return d
	}

	c := &classMemberConfig{}
	for _, opt := range opts {
		opt(c)
	}
	d.appendClassMembers(name, c.members)
	return d
}

func (d *Diagram) appendClassMembers(name string, members []string) {
	if len(members) == 0 {
		d.body = append(d.body, fmt.Sprintf("    class %s", name))
		return
	}

	d.body = append(d.body, fmt.Sprintf("    class %s {", name))
	for _, member := range members {
		d.body = append(d.body, fmt.Sprintf("        %s", member))
	}
	d.body = append(d.body, "    }")
}

// ClassWithLabel adds a class declaration with a display label.
func (d *Diagram) ClassWithLabel(name, label string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    class %s[%s]", name, quote(label)))
	return d
}

// ClassWithAnnotation adds a class declaration with an annotation.
// The annotation can be specified with or without enclosing << >>.
func (d *Diagram) ClassWithAnnotation(name, annotation string) *Diagram {
	// Use separate-line annotation form for wider renderer compatibility.
	d.body = append(d.body, fmt.Sprintf("    class %s", name))
	d.body = append(d.body, fmt.Sprintf("    <<%s>> %s", normalizeAnnotation(annotation), name))
	return d
}

// ClassWithMembers adds a class block with member definitions.
// Members can include properties and methods, e.g. "+int count" and "+Reset()".
// For a typed API with visibility sugar, use Class(name, WithPublicField(...), ...).
func (d *Diagram) ClassWithMembers(name string, members ...string) *Diagram {
	d.appendClassMembers(name, members)
	return d
}

// Interface adds an interface annotation for a class.
func (d *Diagram) Interface(name string) *Diagram {
	return d.ClassWithAnnotation(name, "Interface")
}

// From starts relationship chaining from a source class.
func (d *Diagram) From(from string) *SourceRelationBuilder {
	return &SourceRelationBuilder{
		Diagram: d,
		from:    from,
	}
}

// SourceRelationBuilder is a relationship builder that keeps source class context.
type SourceRelationBuilder struct {
	*Diagram
	from string
}

// Member adds a member using Class : Member syntax.
func (d *Diagram) Member(className, member string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s : %s", className, member))
	return d
}

// Annotation adds a separate annotation declaration for a class.
// The annotation can be specified with or without enclosing << >>.
func (d *Diagram) Annotation(className, annotation string) *Diagram {
	d.body = append(
		d.body,
		fmt.Sprintf("    <<%s>> %s", normalizeAnnotation(annotation), className),
	)
	return d
}

// Relation adds a relationship between classes.
func (d *Diagram) Relation(from string, relationship Relationship, to string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s %s %s", from, relationship, to))
	return d
}

func (d *Diagram) relationWithOptions(
	from string,
	relationship Relationship,
	to string,
	opts ...RelationOption,
) *Diagram {
	cfg := &relationConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.fromCardinality != "" && cfg.toCardinality != "" {
		if cfg.label != "" {
			return d.RelationWithCardinalityAndLabel(
				from,
				cfg.fromCardinality,
				relationship,
				to,
				cfg.toCardinality,
				cfg.label,
			)
		}
		return d.RelationWithCardinality(
			from,
			cfg.fromCardinality,
			relationship,
			to,
			cfg.toCardinality,
		)
	}

	if cfg.label != "" {
		return d.RelationWithLabel(from, relationship, to, cfg.label)
	}

	return d.Relation(from, relationship, to)
}

// Relation adds a relationship from source class to target class.
func (b *SourceRelationBuilder) Relation(
	relationship Relationship,
	to string,
	opts ...RelationOption,
) *SourceRelationBuilder {
	b.relationWithOptions(b.from, relationship, to, opts...)
	return b
}

// Composition adds a composition relationship from source class to target class.
func (b *SourceRelationBuilder) Composition(to string, opts ...RelationOption) *SourceRelationBuilder {
	return b.Relation(RelationshipComposition, to, opts...)
}

// Association adds an association relationship from source class to target class.
func (b *SourceRelationBuilder) Association(to string, opts ...RelationOption) *SourceRelationBuilder {
	return b.Relation(RelationshipAssociation, to, opts...)
}

// Composition adds a composition relationship between classes.
func (d *Diagram) Composition(from, to string) *Diagram {
	return d.Relation(from, RelationshipComposition, to)
}

// CompositionWithLabel adds a composition relationship with a label.
func (d *Diagram) CompositionWithLabel(from, to, label string) *Diagram {
	return d.RelationWithLabel(from, RelationshipComposition, to, label)
}

// CompositionWithCardinality adds a composition relationship with cardinality values.
func (d *Diagram) CompositionWithCardinality(
	from,
	fromCardinality,
	to,
	toCardinality string,
) *Diagram {
	return d.RelationWithCardinality(from, fromCardinality, RelationshipComposition, to, toCardinality)
}

// CompositionWithCardinalityAndLabel adds a composition relationship with cardinality values and a label.
func (d *Diagram) CompositionWithCardinalityAndLabel(
	from,
	fromCardinality,
	to,
	toCardinality,
	label string,
) *Diagram {
	return d.RelationWithCardinalityAndLabel(
		from,
		fromCardinality,
		RelationshipComposition,
		to,
		toCardinality,
		label,
	)
}

// Association adds an association relationship between classes.
func (d *Diagram) Association(from, to string) *Diagram {
	return d.Relation(from, RelationshipAssociation, to)
}

// AssociationWithLabel adds an association relationship with a label.
func (d *Diagram) AssociationWithLabel(from, to, label string) *Diagram {
	return d.RelationWithLabel(from, RelationshipAssociation, to, label)
}

// AssociationWithCardinality adds an association relationship with cardinality values.
func (d *Diagram) AssociationWithCardinality(
	from,
	fromCardinality,
	to,
	toCardinality string,
) *Diagram {
	return d.RelationWithCardinality(from, fromCardinality, RelationshipAssociation, to, toCardinality)
}

// AssociationWithCardinalityAndLabel adds an association relationship with cardinality values and a label.
func (d *Diagram) AssociationWithCardinalityAndLabel(
	from,
	fromCardinality,
	to,
	toCardinality,
	label string,
) *Diagram {
	return d.RelationWithCardinalityAndLabel(
		from,
		fromCardinality,
		RelationshipAssociation,
		to,
		toCardinality,
		label,
	)
}

// RelationWithLabel adds a relationship with a label.
func (d *Diagram) RelationWithLabel(from string, relationship Relationship, to, label string) *Diagram {
	d.body = append(
		d.body,
		fmt.Sprintf("    %s %s %s : %s", from, relationship, to, label),
	)
	return d
}

// RelationWithCardinality adds a relationship with cardinality values.
func (d *Diagram) RelationWithCardinality(
	from,
	fromCardinality string,
	relationship Relationship,
	to,
	toCardinality string,
) *Diagram {
	d.body = append(
		d.body,
		fmt.Sprintf("    %s %s %s %s %s", from, quote(fromCardinality), relationship, quote(toCardinality), to),
	)
	return d
}

// RelationWithCardinalityAndLabel adds a relationship with cardinality values and a label.
func (d *Diagram) RelationWithCardinalityAndLabel(
	from,
	fromCardinality string,
	relationship Relationship,
	to,
	toCardinality,
	label string,
) *Diagram {
	d.body = append(
		d.body,
		fmt.Sprintf("    %s %s %s %s %s : %s",
			from,
			quote(fromCardinality),
			relationship,
			quote(toCardinality),
			to,
			label,
		),
	)
	return d
}

// LollipopInterface adds a lollipop interface relationship.
// It renders as: interfaceName ()-- className
func (d *Diagram) LollipopInterface(interfaceName, className string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s ()-- %s", interfaceName, className))
	return d
}

// LollipopInterfaceReverse adds a reverse lollipop interface relationship.
// It renders as: className --() interfaceName
func (d *Diagram) LollipopInterfaceReverse(className, interfaceName string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    %s --() %s", className, interfaceName))
	return d
}

// Note adds a note for the class diagram.
func (d *Diagram) Note(note string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    note %s", quote(note)))
	return d
}

// NoteFor adds a note associated with a class.
func (d *Diagram) NoteFor(className, note string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    note for %s %s", className, quote(note)))
	return d
}

// Link adds an interaction link for a class.
func (d *Diagram) Link(className, url, tooltip string) *Diagram {
	d.body = append(
		d.body,
		fmt.Sprintf("    link %s %s %s", className, quote(url), quote(tooltip)),
	)
	return d
}

// Callback adds a callback interaction for a class.
func (d *Diagram) Callback(className, callbackName, tooltip string) *Diagram {
	d.body = append(
		d.body,
		fmt.Sprintf("    callback %s %s %s", className, quote(callbackName), quote(tooltip)),
	)
	return d
}

// ClickCall adds a click action using callback call syntax.
func (d *Diagram) ClickCall(className, callbackName, tooltip string) *Diagram {
	d.body = append(
		d.body,
		fmt.Sprintf(
			"    click %s call %s %s",
			className,
			ensureFunctionCall(callbackName),
			quote(tooltip),
		),
	)
	return d
}

// ClickHref adds a click action using href syntax.
func (d *Diagram) ClickHref(className, url, tooltip string) *Diagram {
	d.body = append(
		d.body,
		fmt.Sprintf("    click %s href %s %s", className, quote(url), quote(tooltip)),
	)
	return d
}

// Style adds style settings for a class.
func (d *Diagram) Style(className, style string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    style %s %s", className, style))
	return d
}

// ClassDef adds a classDef style declaration.
func (d *Diagram) ClassDef(classNames, style string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    classDef %s %s", classNames, style))
	return d
}

// CSSClass applies a classDef style to one or more classes.
// classNames can be a comma-separated list.
// This method applies surrounding quotes required by Mermaid.
func (d *Diagram) CSSClass(classNames, classDefName string) *Diagram {
	d.body = append(
		d.body,
		fmt.Sprintf("    cssClass %s %s;", quote(normalizeQuoted(classNames)), classDefName),
	)
	return d
}

// ClassShorthand applies classDef style using shorthand syntax.
func (d *Diagram) ClassShorthand(className, classDefName string) *Diagram {
	d.body = append(d.body, fmt.Sprintf("    class %s:::%s", className, classDefName))
	return d
}

func quote(s string) string {
	escaped := strings.ReplaceAll(s, `\`, "&#92;")
	escaped = strings.ReplaceAll(escaped, "\r", "&#92;r")
	escaped = strings.ReplaceAll(escaped, "\n", "&#92;n")
	escaped = strings.ReplaceAll(escaped, "\t", "&#92;t")
	escaped = strings.ReplaceAll(escaped, `"`, "&quot;")
	return `"` + escaped + `"`
}

func normalizeAnnotation(annotation string) string {
	a := strings.TrimSpace(annotation)
	a = strings.TrimPrefix(a, "<<")
	a = strings.TrimSuffix(a, ">>")
	return a
}

func ensureFunctionCall(callbackName string) string {
	trimmed := strings.TrimSpace(callbackName)
	if strings.HasSuffix(trimmed, ")") {
		return trimmed
	}
	return fmt.Sprintf("%s()", trimmed)
}

func normalizeQuoted(v string) string {
	trimmed := strings.TrimSpace(v)
	trimmed = strings.TrimPrefix(trimmed, "\"")
	trimmed = strings.TrimSuffix(trimmed, "\"")
	return trimmed
}

func visibilityMarker(v Visibility) string {
	return string(v)
}

func formatFieldMember(visibility Visibility, fieldType, fieldName string) string {
	return fmt.Sprintf("%s%s %s", visibilityMarker(visibility), fieldType, fieldName)
}

func formatMethodMember(
	visibility Visibility,
	methodName,
	returnType string,
	params ...string,
) string {
	paramString := strings.Join(params, ", ")
	method := fmt.Sprintf("%s%s(%s)", visibilityMarker(visibility), methodName, paramString)
	if returnType != "" {
		method += fmt.Sprintf(" %s", returnType)
	}
	return method
}
