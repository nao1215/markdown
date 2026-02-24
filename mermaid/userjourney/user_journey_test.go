package userjourney

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
			want: "journey",
		},
		{
			name: "new diagram with title",
			opts: []Option{WithTitle("Checkout Journey")},
			want: `journey
    title Checkout Journey`,
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

	d := NewDiagram(b, WithTitle("Checkout Journey"))
	d.Section("Discover").
		Task("Browse products", ScoreVerySatisfied, "Customer").
		Task("Add item to cart", ScoreSatisfied, "Customer").
		LF().
		Section("Checkout").
		Task("Enter shipping details", ScoreNeutral, "Customer").
		Task("Complete payment", ScoreSatisfied, "Customer", "Payment Service")

	if err := d.Build(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `journey
    title Checkout Journey
    section Discover
        Browse products: 5: Customer
        Add item to cart: 4: Customer

    section Checkout
        Enter shipping details: 3: Customer
        Complete payment: 4: Customer, Payment Service`

	got := strings.ReplaceAll(b.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_TaskIn(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		TaskIn("Discover", "Browse products", ScoreVerySatisfied, "Customer").
		TaskIn("Discover", "Add item to cart", ScoreSatisfied, "Customer").
		TaskIn("Checkout", "Complete payment", ScoreNeutral, "Customer", "Payment Service")

	want := `journey
    section Discover
        Browse products: 5: Customer
        Add item to cart: 4: Customer
    section Checkout
        Complete payment: 3: Customer, Payment Service`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_Error(t *testing.T) {
	t.Parallel()

	t.Run("task before section", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard).
			Task("Browse products", ScoreVerySatisfied, "Customer")

		if d.Error() == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid score lower bound", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard).
			Section("Discover").
			Task("Browse products", Score(0), "Customer")

		if d.Error() == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid score upper bound", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard).
			Section("Discover").
			Task("Browse products", Score(6), "Customer")

		if d.Error() == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid score negative", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard).
			Section("Discover").
			Task("Browse products", Score(-1), "Customer")

		if d.Error() == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty section name", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard).Section("")

		if d.Error() == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("empty task name", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard).
			Section("Discover").
			Task("", ScoreNeutral, "Customer")

		if d.Error() == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("taskin with empty section", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard).
			TaskIn("", "Browse products", ScoreNeutral, "Customer")

		if d.Error() == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestDiagram_ActorNormalization(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		Section("Discover").
		Task("Browse products", ScoreVerySatisfied, " ", "Alice", "", " Bob ").
		Task("Add item to cart", ScoreSatisfied)

	want := `journey
    section Discover
        Browse products: 5: Alice, Bob
        Add item to cart: 4`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_LFAfterErrorShortCircuits(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		Section("").
		LF()

	if d.Error() == nil {
		t.Fatal("expected error, got nil")
	}

	want := "journey"
	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_NewlineValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func() *Diagram
		want string
	}{
		{
			name: "section name with newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Section("Discover\nInjected")
			},
			want: "journey",
		},
		{
			name: "task name with newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Section("Discover").
					Task("Browse\nproducts", ScoreNeutral, "Customer")
			},
			want: `journey
    section Discover`,
		},
		{
			name: "taskin section with newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					TaskIn("Discover\nInjected", "Browse products", ScoreNeutral, "Customer")
			},
			want: "journey",
		},
		{
			name: "actor name with newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Section("Discover").
					Task("Browse products", ScoreNeutral, "Customer\nOps")
			},
			want: `journey
    section Discover`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := tt.run()
			if d.Error() == nil {
				t.Fatal("expected error, got nil")
			}

			got := strings.ReplaceAll(d.String(), "\r\n", "\n")
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
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
