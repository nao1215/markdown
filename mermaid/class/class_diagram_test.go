package class

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type failingWriter struct{}

func (f failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

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
title: Checkout Domain
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
title: Checkout Domain
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
    class PaymentGateway
    <<Interface>> PaymentGateway
    <<Service>> InventoryService
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

	d := NewDiagram(failingWriter{})
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
