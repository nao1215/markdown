package class

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
)

func TestNewDiagram(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts []Option
		want string
	}{
		{
			name: "new diagram without options",
			opts: nil,
			want: "classDiagram",
		},
		{
			name: "new diagram with title",
			opts: []Option{WithTitle("Checkout Domain")},
			want: `---
title: "Checkout Domain"
---
classDiagram`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diagram := NewDiagram(io.Discard, tt.opts...)
			got := strings.ReplaceAll(diagram.String(), "\r\n", "\n")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDiagram_Build(t *testing.T) {
	t.Parallel()

	b := new(bytes.Buffer)

	d := NewDiagram(b, WithTitle("Checkout Domain"))
	d.SetDirection(DirectionLR).
		Comment("Domain model").
		ClassWithLabel("Order", "Order Aggregate").
		Class(
			"LineItem",
			WithPublicField("string", "sku"),
			WithPublicField("int", "quantity"),
			WithPublicMethod("Total", "float64"),
		).
		Interface("PaymentGateway").
		Annotation("InventoryService", "<<Service>>").
		Member("Order", "+Create() error")

	d.From("Order").
		Composition("LineItem", WithOneToMany(), WithRelationLabel("contains")).
		Association("PaymentGateway", WithRelationLabel("uses"))

	d.Relation("PaymentGateway", RelationshipRealizationReverse, "StripeGateway").
		RelationWithLabel("LineItem", RelationshipDependency, "InventoryService", "checks stock").
		Note("Simple checkout flow").
		NoteFor("Order", "Aggregate Root").
		ClassDef("important", "fill:#f96,stroke:#333,stroke-width:2px").
		ClassShorthand("Order", "important").
		LF()

	if err := d.Build(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `---
title: "Checkout Domain"
---
classDiagram
    direction LR
    %% Domain model
    class Order["Order Aggregate"]
    class LineItem {
        +string sku
        +int quantity
        +Total() float64
    }
    class PaymentGateway {
        <<Interface>>
    }
    class InventoryService {
        <<Service>>
    }
    Order : +Create() error
    Order "1" *-- "many" LineItem : contains
    Order --> PaymentGateway : uses
    PaymentGateway <|.. StripeGateway
    LineItem ..> InventoryService : checks stock
    note "Simple checkout flow"
    note for Order "Aggregate Root"
    classDef important fill:#f96,stroke:#333,stroke-width:2px
    class Order:::important
`

	got := strings.ReplaceAll(b.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_InteractionsAndStyles(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		Class("CheckoutService").
		Link("CheckoutService", "https://example.com/docs/checkout", "Open docs").
		Callback("CheckoutService", "onCheckoutClick", "Run callback").
		ClickCall("CheckoutService", "onCheckoutClick", "Run callback").
		ClickHref("CheckoutService", "https://example.com/docs/checkout", "Open docs").
		Style("CheckoutService", "fill:#f9f,stroke:#333,stroke-width:2px").
		ClassDef("important", "fill:#f96,stroke:#333,stroke-width:2px").
		CSSClass("\"CheckoutService\"", "important").
		LollipopInterface("PaymentPort", "CheckoutService").
		LollipopInterfaceReverse("CheckoutService", "NotificationPort")

	want := `classDiagram
    class CheckoutService
    link CheckoutService "https://example.com/docs/checkout" "Open docs"
    callback CheckoutService "onCheckoutClick" "Run callback"
    click CheckoutService call onCheckoutClick() "Run callback"
    click CheckoutService href "https://example.com/docs/checkout" "Open docs"
    style CheckoutService fill:#f9f,stroke:#333,stroke-width:2px
    classDef important fill:#f96,stroke:#333,stroke-width:2px
    cssClass "CheckoutService" important;
    PaymentPort ()-- CheckoutService
    CheckoutService --() NotificationPort`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_QuoteEscapesControlChars(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		Class("CheckoutService").
		Note("line1\nline2").
		NoteFor("CheckoutService", "head\rline").
		Link("CheckoutService", "https://example.com/docs\\checkout", "tab\ttooltip")

	want := `classDiagram
    class CheckoutService
    note "line1&#92;nline2"
    note for CheckoutService "head&#92;rline"
    link CheckoutService "https://example.com/docs&#92;checkout" "tab&#92;ttooltip"`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_ClassMemberOptions(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		Class(
			"Order",
			WithPublicField("string", "id"),
			WithPrivateField("int", "version"),
			WithProtectedField("time.Time", "updatedAt"),
			WithPackageField("time.Time", "createdAt"),
			WithPublicMethod("Create", "error", "items []LineItem"),
			WithPrivateMethod("validate", "error"),
			WithProtectedMethod("afterSave", "", "event DomainEvent"),
			WithPackageMethod("snapshot", "OrderSnapshot"),
		)

	want := `classDiagram
    class Order {
        +string id
        -int version
        #time.Time updatedAt
        ~time.Time createdAt
        +Create(items []LineItem) error
        -validate() error
        #afterSave(event DomainEvent)
        ~snapshot() OrderSnapshot
    }`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_RelationshipSugar(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		Composition("Order", "LineItem").
		CompositionWithLabel("Order", "LineItem", "contains").
		Association("Order", "PaymentGateway").
		AssociationWithLabel("Order", "PaymentGateway", "uses")

	want := `classDiagram
    Order *-- LineItem
    Order *-- LineItem : contains
    Order --> PaymentGateway
    Order --> PaymentGateway : uses`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_SourceRelationBuilder(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard)
	d.From("Order").
		Composition("LineItem", WithOneToMany(), WithRelationLabel("contains")).
		Association("PaymentGateway", WithRelationLabel("uses"))

	want := `classDiagram
    Order "1" *-- "many" LineItem : contains
    Order --> PaymentGateway : uses`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_RelationWithCardinality(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		RelationWithCardinality("Order", "1", RelationshipComposition, "LineItem", "many")

	want := `classDiagram
    Order "1" *-- "many" LineItem`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_Error(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard)
	if err := d.Error(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestDiagram_BuildStoresError(t *testing.T) {
	t.Parallel()

	d := NewDiagram(errWriter{})
	err := d.Build()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if d.Error() == nil {
		t.Fatal("expected persisted error, got nil")
	}
	if !errors.Is(d.Error(), err) {
		t.Fatalf("expected Error() to wrap returned error, got %v", d.Error())
	}
}

// build renders a diagram and returns its lines, minus the "classDiagram"
// header, so each test can assert on the statements it produced.
func build(t *testing.T, fn func(*Diagram) *Diagram) []string {
	t.Helper()

	d := fn(NewDiagram(nil))
	lines := strings.Split(d.String(), internal.LineFeed())
	if len(lines) == 0 {
		t.Fatal("diagram produced no output")
	}
	return lines[1:]
}

// TestCardinalityShorthands covers the three shorthands over WithCardinality.
// They exist so a caller does not have to spell out the pair, and a wrong pair
// would be invisible without an assertion on the rendered relationship.
func TestCardinalityShorthands(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		option RelationOption
		want   string
	}{
		"many to one": {WithManyToOne(), `    Order "many" *-- "1" LineItem`},
		"one to one":  {WithOneToOne(), `    Order "1" *-- "1" LineItem`},
		"many to many": {
			WithManyToMany(),
			`    Order "many" *-- "many" LineItem`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := build(t, func(d *Diagram) *Diagram {
				return d.From("Order").Composition("LineItem", tt.option).Diagram
			})
			if got[0] != tt.want {
				t.Errorf("relationship mismatch:\n got: %q\nwant: %q", got[0], tt.want)
			}
		})
	}
}

// TestClassWithMembers covers the untyped member form, which is the escape
// hatch for members the typed helpers cannot express.
func TestClassWithMembers(t *testing.T) {
	t.Parallel()

	got := build(t, func(d *Diagram) *Diagram {
		return d.ClassWithMembers("Order", "+int count", "+Reset()")
	})

	want := []string{
		"    class Order {",
		"        +int count",
		"        +Reset()",
		"    }",
	}
	for i, line := range want {
		if i >= len(got) || got[i] != line {
			t.Fatalf("member block mismatch at line %d:\n got: %#v\nwant: %#v", i, got, want)
		}
	}
}

// TestRelationshipsWithCardinality covers the composition and association
// helpers that take cardinality, with and without a label.
func TestRelationshipsWithCardinality(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(*Diagram) *Diagram
		want  string
	}{
		"composition": {
			build: func(d *Diagram) *Diagram {
				return d.CompositionWithCardinality("Order", "1", "LineItem", "many")
			},
			want: `    Order "1" *-- "many" LineItem`,
		},
		"composition with label": {
			build: func(d *Diagram) *Diagram {
				return d.CompositionWithCardinalityAndLabel("Order", "1", "LineItem", "many", "contains")
			},
			want: `    Order "1" *-- "many" LineItem : contains`,
		},
		"association": {
			build: func(d *Diagram) *Diagram {
				return d.AssociationWithCardinality("Order", "1", "Payment", "0..1")
			},
			want: `    Order "1" --> "0..1" Payment`,
		},
		"association with label": {
			build: func(d *Diagram) *Diagram {
				return d.AssociationWithCardinalityAndLabel("Order", "1", "Payment", "0..1", "settles")
			},
			want: `    Order "1" --> "0..1" Payment : settles`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := build(t, tt.build); got[0] != tt.want {
				t.Errorf("relationship mismatch:\n got: %q\nwant: %q", got[0], tt.want)
			}
		})
	}
}

// TestAnnotationUsesTheClassBody pins the form GitHub accepts. The standalone
// form, "<<Interface>> Name" on its own line, makes GitHub reject the whole
// diagram because it lexes the leading "<" as a relationship token.
func TestAnnotationUsesTheClassBody(t *testing.T) {
	t.Parallel()

	for name, fn := range map[string]func(*Diagram) *Diagram{
		"Interface":           func(d *Diagram) *Diagram { return d.Interface("PaymentGateway") },
		"ClassWithAnnotation": func(d *Diagram) *Diagram { return d.ClassWithAnnotation("PaymentGateway", "Interface") },
		"Annotation":          func(d *Diagram) *Diagram { return d.Annotation("PaymentGateway", "<<Interface>>") },
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := build(t, fn)
			want := []string{
				"    class PaymentGateway {",
				"        <<Interface>>",
				"    }",
			}
			for i, line := range want {
				if i >= len(got) || got[i] != line {
					t.Fatalf("annotation mismatch:\n got: %#v\nwant: %#v", got, want)
				}
			}
			for _, line := range got {
				if strings.HasPrefix(strings.TrimSpace(line), "<<") && strings.Contains(line, ">> ") {
					t.Errorf("standalone annotation form emitted: %q", line)
				}
			}
		})
	}
}

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).Class("Account")
	})
}

// TestGoldenClassDeclarations pins how classes, interfaces, members, and
// annotations are declared.
func TestGoldenClassDeclarations(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(buf, WithTitle("Declarations")).
		SetDirection(DirectionLR).
		Comment("Every declaration form of the package").
		Class("Empty").
		Class("Account",
			WithField(VisibilityPublic, "string", "id"),
			WithPublicField("string", "owner"),
			WithPrivateField("int", "balance"),
			WithProtectedField("time.Time", "openedAt"),
			WithPackageField("bool", "frozen"),
			WithMethod(VisibilityPublic, "Transfer", "error", "to Account", "amount int"),
			WithPublicMethod("Deposit", "error", "amount int"),
			WithPrivateMethod("validate", "error"),
			WithProtectedMethod("audit", "void"),
			WithPackageMethod("reset", "void"),
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

	relationships := []Relationship{
		RelationshipInheritance,
		RelationshipInheritanceReverse,
		RelationshipComposition,
		RelationshipCompositionReverse,
		RelationshipAggregation,
		RelationshipAggregationReverse,
		RelationshipAssociation,
		RelationshipAssociationReverse,
		RelationshipLink,
		RelationshipDependency,
		RelationshipDependencyReverse,
		RelationshipRealization,
		RelationshipRealizationReverse,
		RelationshipDashedLink,
	}

	buf := &bytes.Buffer{}
	diagram := NewDiagram(buf)
	for _, relationship := range relationships {
		diagram = diagram.Relation("Source", relationship, "Target")
	}

	err := diagram.
		LF().
		RelationWithLabel("Source", RelationshipAssociation, "Target", "labeled").
		RelationWithCardinality("Source", string(CardinalityOne), RelationshipAssociation, "Target", string(CardinalityMany)).
		RelationWithCardinalityAndLabel("Source", string(CardinalityZeroOrOne), RelationshipAssociation, "Target", string(CardinalityZeroOrMore), "labeled").
		LF().
		Composition("Whole", "Part").
		CompositionWithLabel("Whole", "Part", "contains").
		CompositionWithCardinality("Whole", string(CardinalityOne), "Part", string(CardinalityOneOrMore)).
		CompositionWithCardinalityAndLabel("Whole", string(CardinalityOne), "Part", string(CardinalityMany), "contains").
		LF().
		Association("Left", "Right").
		AssociationWithLabel("Left", "Right", "uses").
		AssociationWithCardinality("Left", string(CardinalityOne), "Right", string(CardinalityOne)).
		AssociationWithCardinalityAndLabel("Left", string(CardinalityMany), "Right", string(CardinalityMany), "uses").
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
	err := NewDiagram(buf).
		From("Account").
		Composition("Entry").
		Association("Ledger", WithRelationLabel("writes to")).
		Relation(RelationshipDependency, "Clock").
		Relation(RelationshipAssociation, "OneToOne", WithOneToOne()).
		Relation(RelationshipAssociation, "OneToMany", WithOneToMany()).
		Relation(RelationshipAssociation, "ManyToOne", WithManyToOne()).
		Relation(RelationshipAssociation, "ManyToMany", WithManyToMany()).
		Relation(RelationshipAssociation, "Explicit",
			WithCardinality(CardinalityZeroOrOne, CardinalityOneOrMore),
			WithRelationLabel("explicit cardinality"),
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
	err := NewDiagram(buf).
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

// TestBuildWithNilWriter covers the case where a diagram is built for String()
// only and Build() is called by mistake. Build() used to dereference the nil
// writer and take the process down; it has to return an error instead.
func TestBuildWithNilWriter(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Build() panicked with a nil writer: %v", r)
		}
	}()

	d := NewDiagram(nil)

	// String() has always worked without a writer, and callers rely on it.
	_ = d.String()

	err := d.Build()
	if err == nil {
		t.Fatal("Build() with a nil writer must return an error")
	}
	if err.Error() != "output writer must not be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// errWrite is the failure the writer below reports, so the test can assert that
// Build passed it through rather than inventing an error of its own.
var errWrite = errors.New("write failed")

// errWriter fails every write, which is what a full disk or a closed pipe looks
// like to Build.
type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errWrite
}

// TestBuildReportsWriteFailure covers the branch where the destination accepts
// the diagram and then fails. Silently returning nil there would hand the caller
// a document that was never written.
func TestBuildReportsWriteFailure(t *testing.T) {
	t.Parallel()

	err := NewDiagram(errWriter{}).Build()
	if err == nil {
		t.Fatal("Build must report a failing writer")
	}
	if !errors.Is(err, errWrite) {
		t.Errorf("Build lost the destination error: %v", err)
	}
}

// TestRelationLabelEscapesTheCharactersThatEndTheStatement names the two
// characters this escaping buys. The "A --> B : label" and "A : member" lines
// are the only unquoted text a class diagram takes, and a colon or a semicolon
// in one ended the statement: mermaid refused the whole diagram and the reader
// got an error box rather than a picture.
func TestRelationLabelEscapesTheCharactersThatEndTheStatement(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(*Diagram) *Diagram
		want  string
	}{
		"a colon in a relation label": {
			build: func(d *Diagram) *Diagram {
				return d.RelationWithLabel("A", RelationshipAssociation, "B", "owns: many")
			},
			want: "    A --> B : owns#58; many",
		},
		"a semicolon in a relation label": {
			build: func(d *Diagram) *Diagram {
				return d.RelationWithLabel("A", RelationshipAssociation, "B", "a;b")
			},
			want: "    A --> B : a#59;b",
		},
		"a colon in a relation label with cardinality": {
			build: func(d *Diagram) *Diagram {
				return d.RelationWithCardinalityAndLabel("A", "1", RelationshipAssociation, "B", "*", "a:b")
			},
			want: `    A "1" --> "*" B : a#58;b`,
		},
		"a colon in a member": {
			build: func(d *Diagram) *Diagram {
				return d.Member("A", "start: time.Time")
			},
			want: "    A : start#58; time.Time",
		},
		"a named entity in a label is escaped": {
			build: func(d *Diagram) *Diagram {
				return d.RelationWithLabel("A", RelationshipAssociation, "B", "a#58;b")
			},
			want: "    A --> B : a#35;58#59;b",
		},
		"a plain hash in a label is left alone": {
			build: func(d *Diagram) *Diagram {
				return d.RelationWithLabel("A", RelationshipAssociation, "B", "PR #123 merged")
			},
			want: "    A --> B : PR #123 merged",
		},
		"a hash beside a semicolon keeps both": {
			// The "#" starts no entity, so it stays; the ";" becomes one. The
			// "##59;" that comes out reads back as "#" then ";", because
			// mermaid finds the entity at the second "#".
			build: func(d *Diagram) *Diagram {
				return d.RelationWithLabel("A", RelationshipAssociation, "B", "a#;b")
			},
			want: "    A --> B : a##59;b",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.build(NewDiagram(io.Discard)).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("diagram =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}

// TestQuotedTextKeepsItsColonAndSemicolon pins the other half: everywhere else
// in a class diagram is already quoted, and a class body member, a class label
// and a note each take both characters as they are. Escaping there would change
// output that is correct today.
func TestQuotedTextKeepsItsColonAndSemicolon(t *testing.T) {
	t.Parallel()

	got := NewDiagram(io.Discard).
		ClassWithLabel("A", "a:b;c").
		Note("a:b;c").
		NoteFor("A", "a:b;c").
		String()

	for _, want := range []string{`class A["a:b;c"]`, `note "a:b;c"`, `note for A "a:b;c"`} {
		if !strings.Contains(got, want) {
			t.Errorf("diagram =\n%s\nwant it to contain\n%s", got, want)
		}
	}
}
