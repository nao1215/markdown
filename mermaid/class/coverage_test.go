package class_test

import (
	"strings"
	"testing"

	"github.com/nao1215/markdown/internal"
	"github.com/nao1215/markdown/mermaid/class"
)

// build renders a diagram and returns its lines, minus the "classDiagram"
// header, so each test can assert on the statements it produced.
func build(t *testing.T, fn func(*class.Diagram) *class.Diagram) []string {
	t.Helper()

	d := fn(class.NewDiagram(nil))
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
		option class.RelationOption
		want   string
	}{
		"many to one": {class.WithManyToOne(), `    Order "many" *-- "1" LineItem`},
		"one to one":  {class.WithOneToOne(), `    Order "1" *-- "1" LineItem`},
		"many to many": {
			class.WithManyToMany(),
			`    Order "many" *-- "many" LineItem`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := build(t, func(d *class.Diagram) *class.Diagram {
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

	got := build(t, func(d *class.Diagram) *class.Diagram {
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
		build func(*class.Diagram) *class.Diagram
		want  string
	}{
		"composition": {
			build: func(d *class.Diagram) *class.Diagram {
				return d.CompositionWithCardinality("Order", "1", "LineItem", "many")
			},
			want: `    Order "1" *-- "many" LineItem`,
		},
		"composition with label": {
			build: func(d *class.Diagram) *class.Diagram {
				return d.CompositionWithCardinalityAndLabel("Order", "1", "LineItem", "many", "contains")
			},
			want: `    Order "1" *-- "many" LineItem : contains`,
		},
		"association": {
			build: func(d *class.Diagram) *class.Diagram {
				return d.AssociationWithCardinality("Order", "1", "Payment", "0..1")
			},
			want: `    Order "1" --> "0..1" Payment`,
		},
		"association with label": {
			build: func(d *class.Diagram) *class.Diagram {
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

	for name, fn := range map[string]func(*class.Diagram) *class.Diagram{
		"Interface":           func(d *class.Diagram) *class.Diagram { return d.Interface("PaymentGateway") },
		"ClassWithAnnotation": func(d *class.Diagram) *class.Diagram { return d.ClassWithAnnotation("PaymentGateway", "Interface") },
		"Annotation":          func(d *class.Diagram) *class.Diagram { return d.Annotation("PaymentGateway", "<<Interface>>") },
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
