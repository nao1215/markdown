package userjourney

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).Section("Discovery").Task("Find the site", ScoreSatisfied)
	})
}

// TestRecordedErrorContract asserts that an empty section name surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).Section("")
	})
}

// TestGoldenUserJourney pins the rendered diagram of every builder method of
// this package, including every score.
func TestGoldenUserJourney(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(buf, WithTitle("Sign Up")).
		Section("Discovery").
		Task("Find the site", ScoreVerySatisfied, "Visitor").
		Task("Read the docs", ScoreSatisfied, "Visitor", "Support").
		Section("Registration").
		Task("Fill the form", ScoreNeutral, "Visitor").
		Task("Confirm the mail", ScoreDissatisfied, "Visitor").
		Task("Wait for approval", ScoreVeryDissatisfied).
		LF().
		Section("Onboarding").
		TaskIn("Onboarding", "First login", ScoreSatisfied, "User").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("userjourney.md", buf.String()); err != nil {
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

// TestFieldsEscapeTheirOwnPunctuation names the characters this escaping buys.
// A user journey is written entirely in unquoted text, so each of its four
// fields loses a different set of characters, and each set below was measured
// by rendering one character at a time.
func TestFieldsEscapeTheirOwnPunctuation(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(io.Writer) *Diagram
		want  string
	}{
		"a semicolon in a title loses the diagram": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w, WithTitle("a;b")) },
			want:  "    title a#59;b",
		},
		"a hash in a title cuts it short": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w, WithTitle("PR #12")) },
			want:  "    title PR #35;12",
		},
		"a semicolon and a colon in a section": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Section("a;b:c") },
			want:  "    section a#59;b#58;c",
		},
		"a hash in a section is left alone": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Section("PR #12") },
			want:  "    section PR #12",
		},
		"a hash, a semicolon and a colon in a task name": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).Section("S").Task("a#b;c:d", ScoreSatisfied)
			},
			want: "        a#35;b#59;c#58;d: 4",
		},
		"a comma in an actor splits it in two": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).Section("S").Task("t", ScoreSatisfied, "Ops, EU")
			},
			want: "        t: 4: Ops#44; EU",
		},
		"actors are still separated by a comma": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).Section("S").Task("t", ScoreSatisfied, "Ops", "EU")
			},
			want: "        t: 4: Ops, EU",
		},
		"a comma in a task name is left alone": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).Section("S").Task("a,b", ScoreSatisfied)
			},
			want: "        a,b: 4",
		},
		"a named entity is escaped wherever entities are written": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Section("a#59;b") },
			want:  "    section a#35;59#59;b",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.build(io.Discard).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("diagram =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}
