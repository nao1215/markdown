package treemap_test

import (
	"bytes"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/treemap"
)

func TestDiagram(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(w io.Writer) *treemap.Diagram
		want  []string
	}{
		"a bare diagram": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w)
			},
			want: []string{"treemap-beta"},
		},
		"a title": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w, treemap.WithTitle("Budget"))
			},
			want: []string{"---", `title: "Budget"`, "---", "treemap-beta"},
		},
		"a leaf at the top level": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf("Salaries", 1200)
			},
			want: []string{"treemap-beta", `"Salaries": 1200`},
		},
		"a section holding leaves": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).
					Section("Ops").
					Leaf("Salaries", 1200).
					Leaf("Travel", 300)
			},
			want: []string{
				"treemap-beta",
				`"Ops"`,
				`    "Salaries": 1200`,
				`    "Travel": 300`,
			},
		},
		"nesting and coming back up": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).
					Section("Ops").
					Leaf("Salaries", 1200).
					Section("Cloud").
					Leaf("Compute", 400).
					Parent().
					Leaf("Travel", 300).
					Parent().
					Section("Marketing").
					Leaf("Ads", 800)
			},
			want: []string{
				"treemap-beta",
				`"Ops"`,
				`    "Salaries": 1200`,
				`    "Cloud"`,
				`        "Compute": 400`,
				`    "Travel": 300`,
				`"Marketing"`,
				`    "Ads": 800`,
			},
		},
		"a fractional value": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf("Salaries", 1200.75)
			},
			want: []string{"treemap-beta", `"Salaries": 1200.75`},
		},
		"a zero value": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf("Salaries", 0)
			},
			want: []string{"treemap-beta", `"Salaries": 0`},
		},
		"a value large enough that Go would reach for an exponent": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf("Salaries", 1e21)
			},
			want: []string{"treemap-beta", `"Salaries": 1000000000000000000000`},
		},
		"a double quote is doubled": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Section(`Ops "core"`)
			},
			// A backslash escape is what mermaid's treemap parser refuses;
			// doubling is what it implements.
			want: []string{"treemap-beta", `"Ops ""core"""`},
		},
		"a backslash is left alone": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf(`Back\slash`, 50)
			},
			want: []string{"treemap-beta", `"Back\slash": 50`},
		},
		"punctuation that is text inside quotes": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf("Salaries: base #1; gross (net)", 1200)
			},
			want: []string{"treemap-beta", `"Salaries: base #1; gross (net)": 1200`},
		},
		"names are trimmed": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Section("  Ops  ").Leaf("  Salaries  ", 1200)
			},
			want: []string{"treemap-beta", `"Ops"`, `    "Salaries": 1200`},
		},
		"a line feed": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf("Salaries", 1200).LF().Leaf("Travel", 300)
			},
			want: []string{"treemap-beta", `"Salaries": 1200`, "", `"Travel": 300`},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			if err := tt.build(buf).Build(); err != nil {
				t.Fatalf("Build() = %v, want nil", err)
			}

			want := strings.Join(tt.want, "\n")
			got := strings.ReplaceAll(buf.String(), "\r\n", "\n")
			if got != want {
				t.Errorf("diagram =\n%s\nwant\n%s", got, want)
			}
		})
	}
}

func TestDiagramRejectsWhatItCannotWrite(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(w io.Writer) *treemap.Diagram
		want  string
	}{
		"a title holding a newline": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w, treemap.WithTitle("first\nsecond"))
			},
			want: "title must not contain newline characters",
		},
		"an empty section name": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Section("  ")
			},
			want: "section name must not be empty",
		},
		"a section name holding a newline": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Section("first\nsecond")
			},
			want: "section name must not contain newline characters",
		},
		"an empty leaf name": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf("", 1200)
			},
			want: "leaf name must not be empty",
		},
		"a leaf name holding a carriage return": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf("first\rsecond", 1200)
			},
			want: "leaf name must not contain newline characters",
		},
		"a negative value": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf("Salaries", -1)
			},
			want: `value of leaf "Salaries" must not be negative`,
		},
		"a value that is not a number": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf("Salaries", math.NaN())
			},
			want: `value of leaf "Salaries" must be a number`,
		},
		"an infinite value": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Leaf("Salaries", math.Inf(1))
			},
			want: `value of leaf "Salaries" must be finite`,
		},
		"going up from the top level": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Parent()
			},
			want: "Parent was called at the top level",
		},
		"going up more often than down": {
			build: func(w io.Writer) *treemap.Diagram {
				return treemap.NewDiagram(w).Section("Ops").Parent().Parent()
			},
			want: "Parent was called at the top level",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			d := tt.build(io.Discard)

			err := d.Build()
			if err == nil {
				t.Fatal("Build() = nil, want an error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Build() = %v, want it to mention %q", err, tt.want)
			}
			if d.Error() == nil {
				t.Error("Error() = nil, want the same error Build returned")
			}
		})
	}
}

// TestTheFirstErrorIsKept pins that a chain records the error that explains the
// rest, and that the calls after it change nothing.
func TestTheFirstErrorIsKept(t *testing.T) {
	t.Parallel()

	d := treemap.NewDiagram(io.Discard).
		Section("").
		Leaf("", 1).
		LF().
		Parent().
		Leaf("Salaries", 1200)

	err := d.Error()
	if err == nil {
		t.Fatal("Error() = nil, want an error")
	}
	if want := "section name must not be empty"; !strings.Contains(err.Error(), want) {
		t.Errorf("Error() = %v, want it to mention %q", err, want)
	}
}

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return treemap.NewDiagram(w).Section("Ops").Leaf("Salaries", 1200)
	})
}

// TestRecordedErrorContract asserts that going up from the top level surfaces
// from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return treemap.NewDiagram(w).Parent()
	})
}

// TestGoldenTreemap pins the rendered diagram of every builder method and every
// option of this package, including the quote doubling and three levels of
// nesting.
func TestGoldenTreemap(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := treemap.NewDiagram(buf, treemap.WithTitle("Budget: 2024")).
		Section(`Ops "core"`).
		Leaf("Salaries: base", 1200).
		Section("Cloud #1").
		Leaf("Compute", 400.5).
		Leaf(`Back\slash`, 0).
		Parent().
		Leaf("Travel; trips", 300).
		Parent().
		LF().
		Section("Marketing").
		Leaf("Ads (paid)", 800).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("treemap.md", buf.String()); err != nil {
		t.Error(err)
	}
}
