package wardley_test

import (
	"bytes"
	"errors"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/nao1215/markdown/internal"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/wardley"
)

// errWrite is the failure the writer below reports, so the test can assert that
// Build passed it through rather than inventing an error of its own.
var errWrite = errors.New("write failed")

// errWriter fails every write, which is what a full disk or a closed pipe looks
// like to Build.
type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errWrite
}

// lines joins the given lines with the line ending the library writes on this
// platform, which is what String returns them joined with.
func lines(want ...string) string {
	return strings.Join(want, internal.LineFeed())
}

func TestMap(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(io.Writer) *wardley.Map
		want  string
	}{
		"an empty map is the header alone": {
			build: func(w io.Writer) *wardley.Map { return wardley.NewMap(w) },
			want:  "wardley-beta",
		},
		"a component carries its two coordinates": {
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w).Component("Checkout", 0.5, 0.8)
			},
			want: lines("wardley-beta", "    component Checkout [0.5, 0.8]"),
		},
		"an anchor is the user the map is drawn for": {
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w).Anchor("Customer", 0.9, 0.95)
			},
			want: lines("wardley-beta", "    anchor Customer [0.9, 0.95]"),
		},
		"a link is a dependency": {
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w).Link("Checkout", "Payment")
			},
			want: lines("wardley-beta", "    Checkout -> Payment"),
		},
		"evolve marks a part as moving": {
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w).Evolve("Checkout", 0.8)
			},
			want: lines("wardley-beta", "    evolve Checkout 0.8"),
		},
		"a title keeps the punctuation mermaid reads there": {
			// A quotation mark, a hash and a semicolon are drawn as themselves,
			// so nothing here is escaped.
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w, wardley.WithTitle(`The "core" #1 map; v2`))
			},
			want: lines("wardley-beta", `    title The "core" #1 map; v2`),
		},
		"a title escapes a hash that would start an entity": {
			// "#1;" is the entity for the character with code 1, so a title
			// holding one has to come out different from a title holding that
			// character.
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w, wardley.WithTitle("issue #1; closed"))
			},
			want: lines("wardley-beta", "    title issue #35;1; closed"),
		},
		"a title escapes the percent run that would comment it out": {
			// "%%" opens a mermaid comment, so the rest of the title is dropped
			// and the map still draws, saying less than it was asked to.
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w, wardley.WithTitle("100%% done"))
			},
			want: lines("wardley-beta", "    title 100#37;#37; done"),
		},
		"a lone percent in a title is left alone": {
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w, wardley.WithTitle("50% done"))
			},
			want: lines("wardley-beta", "    title 50% done"),
		},
		"a name may hold spaces and the punctuation mermaid reads": {
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w).Component("Payment (card) wallet_v2-beta", 0.5, 0.5)
			},
			want: lines("wardley-beta", "    component Payment (card) wallet_v2-beta [0.5, 0.5]"),
		},
		"a coordinate is written in plain decimal": {
			// mermaid parses the token as a number, and "1e-07" is not one to
			// it, which is what Go reaches for at the small end.
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w).Component("Tiny", 0.0000001, 1)
			},
			want: lines("wardley-beta", "    component Tiny [0.0000001, 1]"),
		},
		"text is trimmed": {
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w).Component("  Checkout  ", 0.5, 0.5)
			},
			want: lines("wardley-beta", "    component Checkout [0.5, 0.5]"),
		},
		"LF adds a blank line": {
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w).Component("A", 0.1, 0.9).LF().Component("B", 0.8, 0.2)
			},
			want: lines("wardley-beta", "    component A [0.1, 0.9]", "", "    component B [0.8, 0.2]"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			m := tt.build(buf)
			if err := m.Build(); err != nil {
				t.Fatalf("Build() = %v, want nil", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("Build() wrote\n%q\nwant\n%q", got, tt.want)
			}
			if got := m.String(); got != tt.want {
				t.Errorf("String() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestMapRecordsBadInput(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(io.Writer) *wardley.Map
		want  string
	}{
		"a title spanning lines": {
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w, wardley.WithTitle("a\nb"))
			},
			want: "title must not contain newline characters",
		},
		"an empty name": {
			build: func(w io.Writer) *wardley.Map { return wardley.NewMap(w).Component("  ", 0.5, 0.5) },
			want:  "name must not be empty",
		},
		"a full stop in a name": {
			// mermaid loses the whole map on one, and refuses its own escape
			// form there, so it is reported rather than mangled.
			build: func(w io.Writer) *wardley.Map { return wardley.NewMap(w).Component("a.b", 0.5, 0.5) },
			want:  `name "a.b" must hold only letters, digits, spaces`,
		},
		"a quotation mark in a name": {
			build: func(w io.Writer) *wardley.Map { return wardley.NewMap(w).Anchor(`a"b`, 0.5, 0.5) },
			want:  `must hold only letters, digits, spaces`,
		},
		"an ampersand in a name": {
			// "a&b" renders and "a & b" does not, so the character is refused
			// rather than accepted while nobody puts a space around it.
			build: func(w io.Writer) *wardley.Map { return wardley.NewMap(w).Component("a & b", 0.5, 0.5) },
			want:  `must hold only letters, digits, spaces`,
		},
		"non-ASCII in a name": {
			build: func(w io.Writer) *wardley.Map { return wardley.NewMap(w).Component("日本語", 0.5, 0.5) },
			want:  `must hold only letters, digits, spaces`,
		},
		"a name in a link": {
			build: func(w io.Writer) *wardley.Map { return wardley.NewMap(w).Link("a.b", "c") },
			want:  `name "a.b" must hold only`,
		},
		"a coordinate below zero": {
			build: func(w io.Writer) *wardley.Map { return wardley.NewMap(w).Component("A", -0.1, 0.5) },
			want:  "evolution of A must be between 0.0 and 1.0, not -0.1",
		},
		"a coordinate above one": {
			build: func(w io.Writer) *wardley.Map { return wardley.NewMap(w).Component("A", 0.5, 1.5) },
			want:  "visibility of A must be between 0.0 and 1.0, not 1.5",
		},
		"a coordinate that is not a number": {
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w).Component("A", math.NaN(), 0.5)
			},
			want: "evolution of A must be a number",
		},
		"an infinite coordinate": {
			build: func(w io.Writer) *wardley.Map {
				return wardley.NewMap(w).Component("A", math.Inf(1), 0.5)
			},
			want: "evolution of A must be finite",
		},
		"an evolution out of range": {
			build: func(w io.Writer) *wardley.Map { return wardley.NewMap(w).Evolve("A", 2) },
			want:  "evolution of A must be between 0.0 and 1.0, not 2",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := tt.build(io.Discard)
			err := m.Build()
			if err == nil {
				t.Fatalf("Build() = nil, want an error mentioning %q", tt.want)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Build() = %v, want it to mention %q", err, tt.want)
			}
			if m.Error() == nil {
				t.Error("Error() = nil after a failed Build()")
			}
		})
	}
}

// TestTheFirstErrorIsKept covers the rule the whole library shares: the chain
// runs to the end after a bad call, and the error that surfaces is the first
// one, because that is the one that explains the rest.
func TestTheFirstErrorIsKept(t *testing.T) {
	t.Parallel()

	err := wardley.NewMap(io.Discard).Component("a.b", 0.5, 0.5).Component("c", 2, 0.5).Build()
	if err == nil {
		t.Fatal("Build() = nil, want the first error")
	}
	if want := `name "a.b"`; !strings.Contains(err.Error(), want) {
		t.Errorf("Build() = %v, want the first error %q", err, want)
	}
}

// TestBuildWithNilWriter covers the case where a map is built for String() only
// and Build() is called by mistake.
func TestBuildWithNilWriter(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Build() panicked with a nil writer: %v", r)
		}
	}()

	m := wardley.NewMap(nil).Component("A", 0.5, 0.5)
	if m.String() == "" {
		t.Fatal("String() returned nothing for a map with a component in it")
	}

	err := m.Build()
	if err == nil {
		t.Fatal("Build() with a nil writer must return an error")
	}
	if err.Error() != "output writer must not be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestBuildReportsWriteFailure covers the branch where the destination accepts
// the map and then fails.
func TestBuildReportsWriteFailure(t *testing.T) {
	t.Parallel()

	err := wardley.NewMap(errWriter{}).Component("A", 0.5, 0.5).Build()
	if err == nil {
		t.Fatal("Build must report a failing writer")
	}
	if !errors.Is(err, errWrite) {
		t.Errorf("Build lost the destination error: %v", err)
	}
}

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return wardley.NewMap(w).Component("Checkout", 0.5, 0.8)
	})
}

// TestRecordedErrorContract asserts that a name mermaid cannot read surfaces
// from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return wardley.NewMap(w).Component("a.b", 0.5, 0.5)
	})
}

// TestGoldenWardley pins the rendered map of every builder method and every
// option of this package.
func TestGoldenWardley(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := wardley.NewMap(buf, wardley.WithTitle(`Checkout, "as it stands" #1`)).
		Anchor("Customer", 0.95, 0.95).
		Component("Checkout (web)", 0.6, 0.8).
		Component("Payment (refunds)", 0.75, 0.5).
		LF().
		Component("Card network", 0.95, 0.2).
		Link("Customer", "Checkout (web)").
		Link("Checkout (web)", "Payment (refunds)").
		Link("Payment (refunds)", "Card network").
		Evolve("Payment (refunds)", 0.9).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("wardley.md", buf.String()); err != nil {
		t.Error(err)
	}
}
