//go:build linux || darwin

package class_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/class"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := class.NewDiagram(
		io.Discard,
		class.WithTitle("Checkout Domain"),
	).
		SetDirection(class.DirectionLR).
		Class(
			"Order",
			class.WithPublicField("string", "id"),
			class.WithPublicMethod("Create", "error", "items []LineItem"),
			class.WithPublicMethod("Pay", "error"),
		).
		Class(
			"LineItem",
			class.WithPublicField("string", "sku"),
			class.WithPublicField("int", "quantity"),
			class.WithPublicMethod("Subtotal", "int"),
		).
		Interface("PaymentGateway")

	diagram.From("Order").
		Composition("LineItem", class.WithOneToMany(), class.WithRelationLabel("contains")).
		Association("PaymentGateway", class.WithRelationLabel("uses"))

	diagramString := diagram.
		NoteFor("Order", "Aggregate Root").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Class Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagramString).
		Build()

	// Output:
	// ## Class Diagram
	// ```mermaid
	// ---
	// title: "Checkout Domain"
	// ---
	// classDiagram
	//     direction LR
	//     class Order {
	//         +string id
	//         +Create(items []LineItem) error
	//         +Pay() error
	//     }
	//     class LineItem {
	//         +string sku
	//         +int quantity
	//         +Subtotal() int
	//     }
	//     class PaymentGateway {
	//         <<Interface>>
	//     }
	//     Order "1" *-- "many" LineItem : contains
	//     Order --> PaymentGateway : uses
	//     note for Order "Aggregate Root"
	// ```
}

// ExampleDiagram_Composition draws a composition, the diamond that says the whole owns the part.
func ExampleDiagram_Composition() {
	_ = class.NewDiagram(os.Stdout).
		Composition("Order", "LineItem").
		Build()

	// Output:
	// classDiagram
	//     Order *-- LineItem
}

// ExampleDiagram_CompositionWithLabel draws it and says what it means.
func ExampleDiagram_CompositionWithLabel() {
	_ = class.NewDiagram(os.Stdout).
		CompositionWithLabel("Order", "LineItem", "holds").
		Build()

	// Output:
	// classDiagram
	//     Order *-- LineItem : holds
}

// ExampleDiagram_CompositionWithCardinality draws it and says how many of each.
func ExampleDiagram_CompositionWithCardinality() {
	_ = class.NewDiagram(os.Stdout).
		CompositionWithCardinality("Order", "1", "LineItem", "1..*").
		Build()

	// Output:
	// classDiagram
	//     Order "1" *-- "1..*" LineItem
}

// ExampleDiagram_CompositionWithCardinalityAndLabel draws it with both.
func ExampleDiagram_CompositionWithCardinalityAndLabel() {
	_ = class.NewDiagram(os.Stdout).
		CompositionWithCardinalityAndLabel("Order", "1", "LineItem", "1..*", "holds").
		Build()

	// Output:
	// classDiagram
	//     Order "1" *-- "1..*" LineItem : holds
}

// ExampleDiagram_Association draws an association, the plain arrow between two classes.
func ExampleDiagram_Association() {
	_ = class.NewDiagram(os.Stdout).
		Association("Order", "LineItem").
		Build()

	// Output:
	// classDiagram
	//     Order --> LineItem
}

// ExampleDiagram_AssociationWithLabel draws it and says what it means.
func ExampleDiagram_AssociationWithLabel() {
	_ = class.NewDiagram(os.Stdout).
		AssociationWithLabel("Order", "LineItem", "holds").
		Build()

	// Output:
	// classDiagram
	//     Order --> LineItem : holds
}

// ExampleDiagram_AssociationWithCardinality draws it and says how many of each.
func ExampleDiagram_AssociationWithCardinality() {
	_ = class.NewDiagram(os.Stdout).
		AssociationWithCardinality("Order", "1", "LineItem", "1..*").
		Build()

	// Output:
	// classDiagram
	//     Order "1" --> "1..*" LineItem
}

// ExampleDiagram_AssociationWithCardinalityAndLabel draws it with both.
func ExampleDiagram_AssociationWithCardinalityAndLabel() {
	_ = class.NewDiagram(os.Stdout).
		AssociationWithCardinalityAndLabel("Order", "1", "LineItem", "1..*", "holds").
		Build()

	// Output:
	// classDiagram
	//     Order "1" --> "1..*" LineItem : holds
}

// ExampleWithPublicField declares a field anything can read.
func ExampleWithPublicField() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order", class.WithPublicField("string", "id")).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         +string id
	//     }
}

// ExampleWithPrivateField declares a field only the class can read.
func ExampleWithPrivateField() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order", class.WithPrivateField("int", "total")).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         -int total
	//     }
}

// ExampleWithProtectedField declares a field the class and its children can read.
func ExampleWithProtectedField() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order", class.WithProtectedField("string", "status")).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         #string status
	//     }
}

// ExampleWithPackageField declares a field the package can read.
func ExampleWithPackageField() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order", class.WithPackageField("time.Time", "createdAt")).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         ~time.Time createdAt
	//     }
}

// ExampleWithPublicMethod declares a method anything can call.
func ExampleWithPublicMethod() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order", class.WithPublicMethod("Pay", "error")).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         +Pay() error
	//     }
}

// ExampleWithPrivateMethod declares a method only the class can call.
func ExampleWithPrivateMethod() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order", class.WithPrivateMethod("recalculate", "")).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         -recalculate()
	//     }
}

// ExampleWithProtectedMethod declares a method the class and its children can call.
func ExampleWithProtectedMethod() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order", class.WithProtectedMethod("validate", "error")).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         #validate() error
	//     }
}

// ExampleWithPackageMethod declares a method the package can call.
func ExampleWithPackageMethod() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order", class.WithPackageMethod("audit", "error")).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         ~audit() error
	//     }
}

// ExampleNewDiagram shows the shape every class diagram has: a writer, the
// classes, what relates them, and Build.
func ExampleNewDiagram() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order", class.WithPublicField("string", "id"), class.WithPublicMethod("Pay", "error")).
		Class("LineItem", class.WithPublicField("string", "sku")).
		CompositionWithCardinality("Order", "1", "LineItem", "1..*").
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         +string id
	//         +Pay() error
	//     }
	//     class LineItem {
	//         +string sku
	//     }
	//     Order "1" *-- "1..*" LineItem
}

// ExampleDiagram_Class declares a class with its members. Without any it is
// drawn as an empty box.
func ExampleDiagram_Class() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order",
			class.WithPublicField("string", "id"),
			class.WithPrivateField("int", "total"),
			class.WithPublicMethod("Pay", "error"),
		).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         +string id
	//         -int total
	//         +Pay() error
	//     }
}

// ExampleDiagram_ClassWithMembers declares a class from members already built,
// for a diagram assembled from data rather than written out.
func ExampleDiagram_ClassWithMembers() {
	_ = class.NewDiagram(os.Stdout).
		ClassWithMembers("Order", "+string id", "+Pay() error").
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         +string id
	//         +Pay() error
	//     }
}

// ExampleDiagram_ClassWithLabel gives a class a label to be drawn under,
// where its identifier is not what a reader should see.
func ExampleDiagram_ClassWithLabel() {
	_ = class.NewDiagram(os.Stdout).
		ClassWithLabel("Order", "Customer Order").
		Build()

	// Output:
	// classDiagram
	//     class Order["Customer Order"]
}

// ExampleDiagram_ClassWithAnnotation declares a class and what kind of thing it
// is, drawn above the name in guillemets.
func ExampleDiagram_ClassWithAnnotation() {
	_ = class.NewDiagram(os.Stdout).
		ClassWithAnnotation("Payable", "Interface").
		Build()

	// Output:
	// classDiagram
	//     class Payable {
	//         <<Interface>>
	//     }
}

// ExampleDiagram_Annotation says what kind of thing a class is, separately from
// declaring it.
func ExampleDiagram_Annotation() {
	_ = class.NewDiagram(os.Stdout).
		Class("Payable").
		Annotation("Payable", "Interface").
		Build()

	// Output:
	// classDiagram
	//     class Payable
	//     class Payable {
	//         <<Interface>>
	//     }
}

// ExampleDiagram_Interface annotates a class as an interface, which is the same
// as Annotation with "Interface" and shorter to read in a chain.
func ExampleDiagram_Interface() {
	_ = class.NewDiagram(os.Stdout).
		Class("Payable", class.WithPublicMethod("Pay", "error")).
		Interface("Payable").
		Build()

	// Output:
	// classDiagram
	//     class Payable {
	//         +Pay() error
	//     }
	//     class Payable {
	//         <<Interface>>
	//     }
}

// ExampleDiagram_Member adds one member to a class already declared, using
// mermaid's "Class : member" form.
func ExampleDiagram_Member() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order").
		Member("Order", "+string id").
		Build()

	// Output:
	// classDiagram
	//     class Order
	//     Order : +string id
}

// ExampleDiagram_Relation joins two classes with a relationship named outright,
// for code that decides between them rather than picking one.
func ExampleDiagram_Relation() {
	_ = class.NewDiagram(os.Stdout).
		Relation("Order", class.RelationshipInheritance, "Document").
		Build()

	// Output:
	// classDiagram
	//     Order <|-- Document
}

// ExampleDiagram_RelationWithLabel joins two classes and says what the line
// means.
func ExampleDiagram_RelationWithLabel() {
	_ = class.NewDiagram(os.Stdout).
		RelationWithLabel("Order", class.RelationshipAssociation, "Customer", "belongs to").
		Build()

	// Output:
	// classDiagram
	//     Order --> Customer : belongs to
}

// ExampleDiagram_RelationWithCardinality joins two classes and says how many of
// each.
func ExampleDiagram_RelationWithCardinality() {
	_ = class.NewDiagram(os.Stdout).
		RelationWithCardinality("Order", "many",
			class.RelationshipAssociation, "Customer", "1").
		Build()

	// Output:
	// classDiagram
	//     Order "many" --> "1" Customer
}

// ExampleDiagram_RelationWithCardinalityAndLabel joins two classes with both.
func ExampleDiagram_RelationWithCardinalityAndLabel() {
	_ = class.NewDiagram(os.Stdout).
		RelationWithCardinalityAndLabel("Order", "many",
			class.RelationshipAssociation, "Customer", "1", "belongs to").
		Build()

	// Output:
	// classDiagram
	//     Order "many" --> "1" Customer : belongs to
}

// ExampleDiagram_From names a source once and hangs several relationships off
// it, which keeps a class with many links readable.
func ExampleDiagram_From() {
	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Composition("LineItem", class.WithOneToMany()).
		Association("Customer", class.WithRelationLabel("belongs to")).
		Build()

	// Output:
	// classDiagram
	//     Order "1" *-- "many" LineItem
	//     Order --> Customer : belongs to
}

// ExampleSourceRelationBuilder shows what From returns: the same source, with
// one relationship per call.
func ExampleSourceRelationBuilder() {
	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Composition("LineItem").
		Build()

	// Output:
	// classDiagram
	//     Order *-- LineItem
}

// ExampleSourceRelationBuilder_Composition draws a composition from the source.
func ExampleSourceRelationBuilder_Composition() {
	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Composition("LineItem", class.WithOneToMany()).
		Build()

	// Output:
	// classDiagram
	//     Order "1" *-- "many" LineItem
}

// ExampleSourceRelationBuilder_Association draws an association from the
// source.
func ExampleSourceRelationBuilder_Association() {
	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Association("Customer", class.WithRelationLabel("belongs to")).
		Build()

	// Output:
	// classDiagram
	//     Order --> Customer : belongs to
}

// ExampleSourceRelationBuilder_Relation draws a relationship named outright
// from the source.
func ExampleSourceRelationBuilder_Relation() {
	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Relation(class.RelationshipInheritance, "Document").
		Build()

	// Output:
	// classDiagram
	//     Order <|-- Document
}

// ExampleDiagram_LollipopInterface draws the circle-on-a-stick that says a
// class provides an interface.
func ExampleDiagram_LollipopInterface() {
	_ = class.NewDiagram(os.Stdout).
		LollipopInterface("Payable", "Order").
		Build()

	// Output:
	// classDiagram
	//     Payable ()-- Order
}

// ExampleDiagram_LollipopInterfaceReverse draws it the other way round, from
// the class to the interface.
func ExampleDiagram_LollipopInterfaceReverse() {
	_ = class.NewDiagram(os.Stdout).
		LollipopInterfaceReverse("Order", "Payable").
		Build()

	// Output:
	// classDiagram
	//     Order --() Payable
}

// ExampleDiagram_Note puts a note in the diagram, attached to nothing.
func ExampleDiagram_Note() {
	_ = class.NewDiagram(os.Stdout).
		Note("Generated from the domain model").
		Build()

	// Output:
	// classDiagram
	//     note "Generated from the domain model"
}

// ExampleDiagram_NoteFor puts a note beside one class.
func ExampleDiagram_NoteFor() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order").
		NoteFor("Order", "Immutable once paid").
		Build()

	// Output:
	// classDiagram
	//     class Order
	//     note for Order "Immutable once paid"
}

// ExampleDiagram_Comment writes a mermaid comment, which is in the source and
// not in the drawing.
func ExampleDiagram_Comment() {
	_ = class.NewDiagram(os.Stdout).
		Comment("generated by the schema exporter").
		Class("Order").
		Build()

	// Output:
	// classDiagram
	//     %% generated by the schema exporter
	//     class Order
}

// ExampleDiagram_Link makes a class a link, with a tooltip.
func ExampleDiagram_Link() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order").
		Link("Order", "https://example.com/order", "The Order type").
		Build()

	// Output:
	// classDiagram
	//     class Order
	//     link Order "https://example.com/order" "The Order type"
}

// ExampleDiagram_Callback makes a class call a function in the page when it is
// clicked.
func ExampleDiagram_Callback() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order").
		Callback("Order", "showOrder", "Show the order").
		Build()

	// Output:
	// classDiagram
	//     class Order
	//     callback Order "showOrder" "Show the order"
}

// ExampleDiagram_ClickCall is the click form of Callback.
func ExampleDiagram_ClickCall() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order").
		ClickCall("Order", "showOrder", "Show the order").
		Build()

	// Output:
	// classDiagram
	//     class Order
	//     click Order call showOrder() "Show the order"
}

// ExampleDiagram_ClickHref is the click form of Link.
func ExampleDiagram_ClickHref() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order").
		ClickHref("Order", "https://example.com/order", "The Order type").
		Build()

	// Output:
	// classDiagram
	//     class Order
	//     click Order href "https://example.com/order" "The Order type"
}

// ExampleDiagram_Style colors one class outright.
func ExampleDiagram_Style() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order").
		Style("Order", "fill:#f9f").
		Build()

	// Output:
	// classDiagram
	//     class Order
	//     style Order fill:#f9f
}

// ExampleDiagram_ClassDef names a style so several classes can share it.
func ExampleDiagram_ClassDef() {
	_ = class.NewDiagram(os.Stdout).
		ClassDef("aggregate", "fill:#f96").
		Class("Order").
		CSSClass("Order", "aggregate").
		Build()

	// Output:
	// classDiagram
	//     classDef aggregate fill:#f96
	//     class Order
	//     cssClass "Order" aggregate;
}

// ExampleDiagram_CSSClass applies a named style to classes.
func ExampleDiagram_CSSClass() {
	_ = class.NewDiagram(os.Stdout).
		ClassDef("aggregate", "fill:#f96").
		Class("Order").
		CSSClass("Order", "aggregate").
		Build()

	// Output:
	// classDiagram
	//     classDef aggregate fill:#f96
	//     class Order
	//     cssClass "Order" aggregate;
}

// ExampleDiagram_ClassShorthand applies a style with mermaid's ":::" form,
// which says the same thing as CSSClass in one token.
func ExampleDiagram_ClassShorthand() {
	_ = class.NewDiagram(os.Stdout).
		ClassDef("aggregate", "fill:#f96").
		ClassShorthand("Order", "aggregate").
		Build()

	// Output:
	// classDiagram
	//     classDef aggregate fill:#f96
	//     class Order:::aggregate
}

// ExampleDiagram_SetDirection says which way the diagram is laid out.
func ExampleDiagram_SetDirection() {
	_ = class.NewDiagram(os.Stdout).
		SetDirection(class.DirectionLR).
		Class("Order").
		Build()

	// Output:
	// classDiagram
	//     direction LR
	//     class Order
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := class.NewDiagram(io.Discard).Class("Order").String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// classDiagram
	//     class Order
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the error the chain
// recorded.
func ExampleDiagram_Build() {
	err := class.NewDiagram(nil).Class("Order").Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := class.NewDiagram(io.Discard).Class("Order")
	fmt.Println("error:", d.Error())

	// Output:
	// error: <nil>
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order").
		LF().
		Class("LineItem").
		Build()

	// Output:
	// classDiagram
	//     class Order
	//
	//     class LineItem
}

// ExampleWithField declares a member with its visibility named outright, for
// code that decides between them rather than picking one.
func ExampleWithField() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order", class.WithField(class.VisibilityPublic, "string", "id")).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         +string id
	//     }
}

// ExampleWithMethod declares a method with its visibility named outright, and
// with the parameters it takes.
func ExampleWithMethod() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order",
			class.WithMethod(class.VisibilityPublic, "Pay", "error", "amount int"),
		).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         +Pay(amount int) error
	//     }
}

// ExampleWithRelationLabel says what a relationship means, for the From form.
func ExampleWithRelationLabel() {
	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Association("Customer", class.WithRelationLabel("belongs to")).
		Build()

	// Output:
	// classDiagram
	//     Order --> Customer : belongs to
}

// ExampleWithCardinality says how many of each end a relationship has, for the
// From form.
func ExampleWithCardinality() {
	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Composition("LineItem",
			class.WithCardinality(class.CardinalityOne, class.CardinalityOneOrMore)).
		Build()

	// Output:
	// classDiagram
	//     Order "1" *-- "1..*" LineItem
}

// ExampleWithOneToOne is WithCardinality for the commonest pairing, and reads
// as what it means.
func ExampleWithOneToOne() {
	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Association("Invoice", class.WithOneToOne()).
		Build()

	// Output:
	// classDiagram
	//     Order "1" --> "1" Invoice
}

// ExampleWithOneToMany says one of the source and many of the target.
func ExampleWithOneToMany() {
	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Composition("LineItem", class.WithOneToMany()).
		Build()

	// Output:
	// classDiagram
	//     Order "1" *-- "many" LineItem
}

// ExampleWithManyToOne says many of the source and one of the target.
func ExampleWithManyToOne() {
	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Association("Customer", class.WithManyToOne()).
		Build()

	// Output:
	// classDiagram
	//     Order "many" --> "1" Customer
}

// ExampleWithManyToMany says many of both.
func ExampleWithManyToMany() {
	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Association("Tag", class.WithManyToMany()).
		Build()

	// Output:
	// classDiagram
	//     Order "many" --> "many" Tag
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = class.NewDiagram(os.Stdout, class.WithTitle("Checkout Domain")).
		Class("Order").
		Build()

	// Output:
	// ---
	// title: "Checkout Domain"
	// ---
	// classDiagram
	//     class Order
}

// ExampleCardinality shows how many of one end a relationship has.
func ExampleCardinality() {
	_ = class.NewDiagram(os.Stdout).
		RelationWithCardinality("Order", string(class.CardinalityOne),
			class.RelationshipComposition, "LineItem", string(class.CardinalityOneOrMore)).
		RelationWithCardinality("Order", string(class.CardinalityZeroOrOne),
			class.RelationshipAssociation, "Coupon", string(class.CardinalityZeroOrMore)).
		Build()

	// Output:
	// classDiagram
	//     Order "1" *-- "1..*" LineItem
	//     Order "0..1" --> "0..*" Coupon
}

// ExampleVisibility shows the four visibilities a member can have.
func ExampleVisibility() {
	_ = class.NewDiagram(os.Stdout).
		Class("Order",
			class.WithField(class.VisibilityPublic, "string", "id"),
			class.WithField(class.VisibilityPrivate, "int", "total"),
			class.WithField(class.VisibilityProtected, "string", "status"),
			class.WithField(class.VisibilityPackage, "string", "region"),
		).
		Build()

	// Output:
	// classDiagram
	//     class Order {
	//         +string id
	//         -int total
	//         #string status
	//         ~string region
	//     }
}

// ExampleRelationship shows the lines two classes can be joined with. Each has
// a reverse, which points the same relationship the other way.
func ExampleRelationship() {
	_ = class.NewDiagram(os.Stdout).
		Relation("Order", class.RelationshipInheritance, "Document").
		Relation("Order", class.RelationshipComposition, "LineItem").
		Relation("Order", class.RelationshipAggregation, "Coupon").
		Relation("Order", class.RelationshipDependency, "Clock").
		Build()

	// Output:
	// classDiagram
	//     Order <|-- Document
	//     Order *-- LineItem
	//     Order o-- Coupon
	//     Order ..> Clock
}

// ExampleDirection shows the four ways a diagram can be laid out.
func ExampleDirection() {
	for _, direction := range []class.Direction{
		class.DirectionTB, class.DirectionBT, class.DirectionLR, class.DirectionRL,
	} {
		_ = class.NewDiagram(os.Stdout).SetDirection(direction).Class("Order").Build()
		fmt.Println()
	}

	// Output:
	// classDiagram
	//     direction TB
	//     class Order
	// classDiagram
	//     direction BT
	//     class Order
	// classDiagram
	//     direction LR
	//     class Order
	// classDiagram
	//     direction RL
	//     class Order
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []class.Option{class.WithTitle("Checkout Domain")}

	_ = class.NewDiagram(os.Stdout, options...).Class("Order").Build()

	// Output:
	// ---
	// title: "Checkout Domain"
	// ---
	// classDiagram
	//     class Order
}

// ExampleClassMemberOption shows what a ClassMemberOption is: a function that
// declares one member, passed to Class after its name.
func ExampleClassMemberOption() {
	options := []class.ClassMemberOption{
		class.WithPublicField("string", "id"),
		class.WithPublicMethod("Pay", "error"),
	}

	_ = class.NewDiagram(os.Stdout).Class("Order", options...).Build()

	// Output:
	// classDiagram
	//     class Order {
	//         +string id
	//         +Pay() error
	//     }
}

// ExampleRelationOption shows what a RelationOption is: a function that changes
// how a relationship from From is written.
func ExampleRelationOption() {
	options := []class.RelationOption{
		class.WithOneToMany(),
		class.WithRelationLabel("holds"),
	}

	_ = class.NewDiagram(os.Stdout).
		From("Order").
		Composition("LineItem", options...).
		Build()

	// Output:
	// classDiagram
	//     Order "1" *-- "many" LineItem : holds
}
