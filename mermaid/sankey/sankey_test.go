package sankey_test

import (
	"bytes"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/nao1215/markdown/mermaid/sankey"
)

func TestDiagram(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(w io.Writer) *sankey.Diagram
		want  []string
	}{
		"a bare diagram": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w)
			},
			want: []string{"sankey-beta", ""},
		},
		"a title": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w, sankey.WithTitle("Energy flow"))
			},
			want: []string{"---", `title: "Energy flow"`, "---", "sankey-beta", ""},
		},
		"a title holding a colon is quoted by the front matter": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w, sankey.WithTitle("Energy: flow"))
			},
			want: []string{"---", `title: "Energy: flow"`, "---", "sankey-beta", ""},
		},
		"one flow": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Coal", "Electricity", 100)
			},
			want: []string{"sankey-beta", "", "Coal,Electricity,100"},
		},
		"several flows through one node": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).
					Link("Coal", "Electricity", 100).
					Link("Gas", "Electricity", 50).
					Link("Electricity", "Homes", 150)
			},
			want: []string{
				"sankey-beta",
				"",
				"Coal,Electricity,100",
				"Gas,Electricity,50",
				"Electricity,Homes,150",
			},
		},
		"a fractional quantity": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Bio-conversion", "Liquid", 0.597)
			},
			want: []string{"sankey-beta", "", "Bio-conversion,Liquid,0.597"},
		},
		"a zero quantity": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Coal", "Electricity", 0)
			},
			want: []string{"sankey-beta", "", "Coal,Electricity,0"},
		},
		"a quantity large enough that Go would reach for an exponent": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Coal", "Electricity", 1e21)
			},
			// mermaid parses the field as a number, and "1e+21" is not one to it.
			want: []string{"sankey-beta", "", "Coal,Electricity,1000000000000000000000"},
		},
		"a node name holding a comma is quoted": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Losses, and more", "Solid", 280.322)
			},
			// Without the quotes the comma would end the field, and the quantity
			// would land in the target column.
			want: []string{"sankey-beta", "", `"Losses, and more",Solid,280.322`},
		},
		"a node name holding a double quote is quoted and doubled": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link(`He said "hi"`, "Gas", 81.144)
			},
			want: []string{"sankey-beta", "", `"He said ""hi""",Gas,81.144`},
		},
		"a node name holding a single quote needs nothing": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Agricultural 'waste'", "Bio-conversion", 124.729)
			},
			want: []string{"sankey-beta", "", "Agricultural 'waste',Bio-conversion,124.729"},
		},
		"node names are trimmed": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("  Coal  ", "  Electricity  ", 100)
			},
			want: []string{"sankey-beta", "", "Coal,Electricity,100"},
		},
		"a line feed": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Coal", "Electricity", 100).LF().Link("Gas", "Electricity", 50)
			},
			want: []string{"sankey-beta", "", "Coal,Electricity,100", "", "Gas,Electricity,50"},
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
		build func(w io.Writer) *sankey.Diagram
		want  string
	}{
		"a title holding a newline": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w, sankey.WithTitle("first\nsecond"))
			},
			want: "title must not contain newline characters",
		},
		"an empty source": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("  ", "Electricity", 100)
			},
			want: "source must not be empty",
		},
		"an empty target": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Coal", "", 100)
			},
			want: "target must not be empty",
		},
		"a source holding a newline": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("first\nsecond", "Electricity", 100)
			},
			want: "source must not contain newline characters",
		},
		"a target holding a carriage return": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Coal", "first\rsecond", 100)
			},
			want: "target must not contain newline characters",
		},
		"a negative quantity": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Coal", "Electricity", -1)
			},
			want: `must not be negative`,
		},
		"a quantity that is not a number": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Coal", "Electricity", math.NaN())
			},
			want: `must be a number`,
		},
		"an infinite quantity": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Coal", "Electricity", math.Inf(1))
			},
			want: `must be finite`,
		},
		"a negatively infinite quantity": {
			build: func(w io.Writer) *sankey.Diagram {
				return sankey.NewDiagram(w).Link("Coal", "Electricity", math.Inf(-1))
			},
			want: `must be finite`,
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

	d := sankey.NewDiagram(io.Discard).
		Link("", "Electricity", 100).
		Link("Coal", "", 100).
		LF().
		Link("Gas", "Electricity", 50)

	err := d.Error()
	if err == nil {
		t.Fatal("Error() = nil, want an error")
	}
	if want := "source must not be empty"; !strings.Contains(err.Error(), want) {
		t.Errorf("Error() = %v, want it to mention %q", err, want)
	}
}
