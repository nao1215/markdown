package class_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/class"
)

// TestGoldenClassDeclarations pins how classes, interfaces, members, and
// annotations are declared.
func TestGoldenClassDeclarations(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := class.NewDiagram(buf, class.WithTitle("Declarations")).
		SetDirection(class.DirectionLR).
		Comment("Every declaration form of the package").
		Class("Empty").
		Class("Account",
			class.WithField(class.VisibilityPublic, "string", "id"),
			class.WithPublicField("string", "owner"),
			class.WithPrivateField("int", "balance"),
			class.WithProtectedField("time.Time", "openedAt"),
			class.WithPackageField("bool", "frozen"),
			class.WithMethod(class.VisibilityPublic, "Transfer", "error", "to Account", "amount int"),
			class.WithPublicMethod("Deposit", "error", "amount int"),
			class.WithPrivateMethod("validate", "error"),
			class.WithProtectedMethod("audit", "void"),
			class.WithPackageMethod("reset", "void"),
		).
		ClassWithLabel("Ledger", "Append only ledger").
		ClassWithAnnotation("Repository", "interface").
		ClassWithMembers("Statement", "+string period", "+Render() string").
		Interface("Persistable").
		Member("Ledger", "+Append(entry Entry) error").
		Annotation("Ledger", "service").
		LF().
		Note("A note that belongs to the diagram").
		NoteFor("Account", "A note that belongs to one class").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("class_declarations.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenClassRelationships pins every relationship form, including the
// cardinality and label variants and the fluent source builder.
func TestGoldenClassRelationships(t *testing.T) {
	t.Parallel()

	relationships := []class.Relationship{
		class.RelationshipInheritance,
		class.RelationshipInheritanceReverse,
		class.RelationshipComposition,
		class.RelationshipCompositionReverse,
		class.RelationshipAggregation,
		class.RelationshipAggregationReverse,
		class.RelationshipAssociation,
		class.RelationshipAssociationReverse,
		class.RelationshipLink,
		class.RelationshipDependency,
		class.RelationshipDependencyReverse,
		class.RelationshipRealization,
		class.RelationshipRealizationReverse,
		class.RelationshipDashedLink,
	}

	buf := &bytes.Buffer{}
	diagram := class.NewDiagram(buf)
	for _, relationship := range relationships {
		diagram = diagram.Relation("Source", relationship, "Target")
	}

	err := diagram.
		LF().
		RelationWithLabel("Source", class.RelationshipAssociation, "Target", "labeled").
		RelationWithCardinality("Source", string(class.CardinalityOne), class.RelationshipAssociation, "Target", string(class.CardinalityMany)).
		RelationWithCardinalityAndLabel("Source", string(class.CardinalityZeroOrOne), class.RelationshipAssociation, "Target", string(class.CardinalityZeroOrMore), "labeled").
		LF().
		Composition("Whole", "Part").
		CompositionWithLabel("Whole", "Part", "contains").
		CompositionWithCardinality("Whole", string(class.CardinalityOne), "Part", string(class.CardinalityOneOrMore)).
		CompositionWithCardinalityAndLabel("Whole", string(class.CardinalityOne), "Part", string(class.CardinalityMany), "contains").
		LF().
		Association("Left", "Right").
		AssociationWithLabel("Left", "Right", "uses").
		AssociationWithCardinality("Left", string(class.CardinalityOne), "Right", string(class.CardinalityOne)).
		AssociationWithCardinalityAndLabel("Left", string(class.CardinalityMany), "Right", string(class.CardinalityMany), "uses").
		LF().
		LollipopInterface("Persistable", "Account").
		LollipopInterfaceReverse("Account", "Persistable").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("class_relationships.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenClassSourceRelations pins the fluent builder that states the source
// class once, together with every relation option it accepts.
func TestGoldenClassSourceRelations(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := class.NewDiagram(buf).
		From("Account").
		Composition("Entry").
		Association("Ledger", class.WithRelationLabel("writes to")).
		Relation(class.RelationshipDependency, "Clock").
		Relation(class.RelationshipAssociation, "OneToOne", class.WithOneToOne()).
		Relation(class.RelationshipAssociation, "OneToMany", class.WithOneToMany()).
		Relation(class.RelationshipAssociation, "ManyToOne", class.WithManyToOne()).
		Relation(class.RelationshipAssociation, "ManyToMany", class.WithManyToMany()).
		Relation(class.RelationshipAssociation, "Explicit",
			class.WithCardinality(class.CardinalityZeroOrOne, class.CardinalityOneOrMore),
			class.WithRelationLabel("explicit cardinality"),
		).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("class_source_relations.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenClassInteractionsAndStyles pins the click handlers and the styling
// methods, which are the parts a reader interacts with rather than reads.
func TestGoldenClassInteractionsAndStyles(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := class.NewDiagram(buf).
		Class("Account").
		Link("Account", "https://example.com/account", "Open the docs").
		Callback("Account", "showDetails", "Show the details").
		ClickCall("Account", "openDetails", "Open the details").
		ClickHref("Account", "https://example.com/account", "Open in a new tab").
		LF().
		Style("Account", "fill:#f9f,stroke:#333").
		ClassDef("highlight", "fill:#ff0").
		CSSClass("Account", "highlight").
		ClassShorthand("Account", "highlight").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("class_interactions.md", buf.String()); err != nil {
		t.Error(err)
	}
}
