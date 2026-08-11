package radar_test

import (
	"bytes"
	"io"
	"math"
	"strings"
	"testing"

	"github.com/nao1215/markdown/mermaid/radar"
)

func TestDiagram(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(w io.Writer) *radar.Diagram
		want  []string
	}{
		"a bare chart": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w)
			},
			want: []string{"radar-beta"},
		},
		"a title": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w, radar.WithTitle("Grades"))
			},
			want: []string{"---", `title: "Grades"`, "---", "radar-beta"},
		},
		"one axis": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Axis("Math")
			},
			want: []string{"radar-beta", `  axis a1["Math"]`},
		},
		"several axes in one call": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Axis("Math", "Science", "English")
			},
			want: []string{"radar-beta", `  axis a1["Math"], a2["Science"], a3["English"]`},
		},
		"axes numbered across calls": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Axis("Math", "Science").Axis("English")
			},
			// The identifiers carry on counting, because a chart with two axes
			// named a1 is a chart with one axis.
			want: []string{
				"radar-beta",
				`  axis a1["Math"], a2["Science"]`,
				`  axis a3["English"]`,
			},
		},
		"a curve": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Axis("Math", "Science").Curve("Alice", 85, 90)
			},
			want: []string{
				"radar-beta",
				`  axis a1["Math"], a2["Science"]`,
				`  curve c1["Alice"]{85, 90}`,
			},
		},
		"curves numbered in order": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Curve("Alice", 85).Curve("Bob", 70)
			},
			want: []string{
				"radar-beta",
				`  curve c1["Alice"]{85}`,
				`  curve c2["Bob"]{70}`,
			},
		},
		"fractional and negative values": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Curve("Alice", 85.5, -10, 0)
			},
			want: []string{"radar-beta", `  curve c1["Alice"]{85.5, -10, 0}`},
		},
		"a value large enough that Go would reach for an exponent": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Curve("Alice", 1e21)
			},
			want: []string{"radar-beta", `  curve c1["Alice"]{1000000000000000000000}`},
		},
		"the scale": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Max(100).Min(0)
			},
			want: []string{"radar-beta", "  max 100", "  min 0"},
		},
		"a label holding a double quote is escaped": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Axis(`Quote "lit"`)
			},
			want: []string{"radar-beta", `  axis a1["Quote \"lit\""]`},
		},
		"a label holding a backslash is escaped": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Axis(`Back\slash`)
			},
			// A trailing backslash would otherwise swallow the closing quote.
			want: []string{"radar-beta", `  axis a1["Back\\slash"]`},
		},
		"a label holding punctuation that is text inside quotes": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Curve("Alice: team #1; lead (senior)", 85)
			},
			want: []string{"radar-beta", `  curve c1["Alice: team #1; lead (senior)"]{85}`},
		},
		"labels are trimmed": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Axis("  Math  ").Curve("  Alice  ", 85)
			},
			want: []string{"radar-beta", `  axis a1["Math"]`, `  curve c1["Alice"]{85}`},
		},
		"a line feed": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Axis("Math").LF().Curve("Alice", 85)
			},
			want: []string{"radar-beta", `  axis a1["Math"]`, "", `  curve c1["Alice"]{85}`},
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
				t.Errorf("chart =\n%s\nwant\n%s", got, want)
			}
		})
	}
}

func TestDiagramRejectsWhatItCannotWrite(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(w io.Writer) *radar.Diagram
		want  string
	}{
		"a title holding a newline": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w, radar.WithTitle("first\nsecond"))
			},
			want: "title must not contain newline characters",
		},
		"an axis with no labels": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Axis()
			},
			want: "axis requires at least one label",
		},
		"an empty axis label": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Axis("Math", "  ")
			},
			want: "axis label 2 must not be empty",
		},
		"an axis label holding a newline": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Axis("first\nsecond")
			},
			want: "axis label 1 must not contain newline characters",
		},
		"an empty curve label": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Curve("", 85)
			},
			want: "curve label must not be empty",
		},
		"a curve label holding a carriage return": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Curve("first\rsecond", 85)
			},
			want: "curve label must not contain newline characters",
		},
		"a curve with no values": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Curve("Alice")
			},
			want: `curve "Alice" requires at least one value`,
		},
		"a curve value that is not a number": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Curve("Alice", 85, math.NaN())
			},
			want: `value 2 of curve "Alice" must be a number`,
		},
		"an infinite curve value": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Curve("Alice", math.Inf(1))
			},
			want: `value 1 of curve "Alice" must be finite`,
		},
		"a max that is not a number": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Max(math.NaN())
			},
			want: "max must be a number",
		},
		"an infinite min": {
			build: func(w io.Writer) *radar.Diagram {
				return radar.NewDiagram(w).Min(math.Inf(-1))
			},
			want: "min must be finite",
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

	d := radar.NewDiagram(io.Discard).
		Axis("").
		Curve("").
		LF().
		Max(math.NaN()).
		Min(math.Inf(1)).
		Curve("Alice", 85)

	err := d.Error()
	if err == nil {
		t.Fatal("Error() = nil, want an error")
	}
	if want := "axis label 1 must not be empty"; !strings.Contains(err.Error(), want) {
		t.Errorf("Error() = %v, want it to mention %q", err, want)
	}
}
