package xychart

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
		name    string
		opts    []Option
		want    string
		wantErr bool
	}{
		{
			name: "new diagram without options",
			opts: nil,
			want: "xychart",
		},
		{
			name: "new diagram with title",
			opts: []Option{WithTitle("Sales Revenue")},
			want: `xychart
    title "Sales Revenue"`,
		},
		{
			name: "new diagram horizontal",
			opts: []Option{WithHorizontal()},
			want: "xychart horizontal",
		},
		{
			name:    "new diagram with title including newline",
			opts:    []Option{WithTitle("Sales\nRevenue")},
			want:    "xychart",
			wantErr: true,
		},
		{
			name:    "new diagram with invalid orientation",
			opts:    []Option{WithOrientation(Orientation("diagonal"))},
			want:    "xychart",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diagram := NewDiagram(io.Discard, tt.opts...)
			if tt.wantErr && diagram.Error() == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && diagram.Error() != nil {
				t.Fatalf("unexpected error: %v", diagram.Error())
			}

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

	d := NewDiagram(b, WithTitle("Sales Revenue")).
		XAxisLabels("Jan", "Feb", "Mar", "Apr", "May", "Jun").
		YAxisRangeWithTitle("Revenue (k$)", 0, 100).
		Bar(25, 40, 60, 80, 70, 90).
		Line(30, 50, 70, 85, 75, 95)

	if err := d.Build(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `xychart
    title "Sales Revenue"
    x-axis [Jan, Feb, Mar, Apr, May, Jun]
    y-axis "Revenue (k$)" 0 --> 100
    bar [25, 40, 60, 80, 70, 90]
    line [30, 50, 70, 85, 75, 95]`

	got := strings.ReplaceAll(b.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_AxisAndLabels(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard, WithHorizontal()).
		XAxisLabelsWithTitle("Month Name", "Jan", "Feb 2026", `"Mar"`).
		YAxisRange(-10.5, 120.25).
		Line(1, 2.5, -3.75)

	want := `xychart horizontal
    x-axis "Month Name" [Jan, "Feb 2026", Mar]
    y-axis -10.5 --> 120.25
    line [1, 2.5, -3.75]`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_QuoteEscapesSpecialChars(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard, WithTitle(`Revenue "Q1"\FY26`)).
		XAxisLabels(`Jan\2026`, `Feb "2026"`).
		YAxisRangeWithTitle(`Revenue "k$"\path`, 0, 100).
		Bar(1, 2)

	want := `xychart
    title "Revenue &quot;Q1&quot;&#92;FY26"
    x-axis ["Jan&#92;2026", "Feb &quot;2026&quot;"]
    y-axis "Revenue &quot;k$&quot;&#92;path" 0 --> 100
    bar [1, 2]`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func() *Diagram
		want string
	}{
		{
			name: "empty x-axis labels",
			run: func() *Diagram {
				return NewDiagram(io.Discard).XAxisLabels()
			},
			want: "xychart",
		},
		{
			name: "empty x-axis label value",
			run: func() *Diagram {
				return NewDiagram(io.Discard).XAxisLabels("Jan", "")
			},
			want: "xychart",
		},
		{
			name: "newline in x-axis title",
			run: func() *Diagram {
				return NewDiagram(io.Discard).XAxisLabelsWithTitle("Month\nName", "Jan")
			},
			want: "xychart",
		},
		{
			name: "newline in x-axis label",
			run: func() *Diagram {
				return NewDiagram(io.Discard).XAxisLabels("Jan\nuary")
			},
			want: "xychart",
		},
		{
			name: "invalid x-axis range",
			run: func() *Diagram {
				return NewDiagram(io.Discard).XAxisRange(10, 10)
			},
			want: "xychart",
		},
		{
			name: "invalid y-axis range",
			run: func() *Diagram {
				return NewDiagram(io.Discard).YAxisRange(5, 4)
			},
			want: "xychart",
		},
		{
			name: "newline in y-axis title",
			run: func() *Diagram {
				return NewDiagram(io.Discard).YAxisRangeWithTitle("Revenue\n(k$)", 0, 100)
			},
			want: "xychart",
		},
		{
			name: "empty bar values",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Bar()
			},
			want: "xychart",
		},
		{
			name: "empty line values",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Line()
			},
			want: "xychart",
		},
		{
			name: "lf short-circuit after error",
			run: func() *Diagram {
				return NewDiagram(io.Discard).XAxisLabels().LF()
			},
			want: "xychart",
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

func TestDiagram_BuildNilWriter(t *testing.T) {
	t.Parallel()

	d := NewDiagram(nil)
	err := d.Build()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "output writer must not be nil" {
		t.Fatalf("unexpected error: %v", err)
	}

	if d.Error() == nil {
		t.Fatal("expected persisted error, got nil")
	}
	if !errors.Is(d.Error(), err) {
		t.Fatalf("expected Error() to wrap returned error, got %v", d.Error())
	}
}

// TestXAxisRange covers the numeric x-axis, including the rejected range. A
// chart whose axis silently fails to render is worse than one that reports the
// problem, so the error path matters as much as the happy one.
func TestXAxisRange(t *testing.T) {
	t.Parallel()

	t.Run("with a title", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(nil).XAxisRangeWithTitle("Month", 1, 12)
		if err := d.Error(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `    x-axis Month 1 --> 12`
		if !strings.Contains(d.String(), want) {
			t.Errorf("axis missing from output:\n got: %q\nwant to contain: %q", d.String(), want)
		}
	})

	t.Run("min not below max is rejected", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(nil).XAxisRangeWithTitle("Month", 12, 12)
		if d.Error() == nil {
			t.Fatal("an inverted range must be reported")
		}
	})

	t.Run("a later call is a no-op once the diagram holds an error", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(nil).
			XAxisRangeWithTitle("Month", 12, 12).
			XAxisRangeWithTitle("Month", 1, 12)

		if strings.Contains(d.String(), "1 --> 12") {
			t.Errorf("statement was appended after an error:\n%s", d.String())
		}
	})
}

// TestLineFeedStopsOnError pins the same rule for LF: once the diagram is in an
// error state, nothing more is appended.
func TestLineFeedStopsOnError(t *testing.T) {
	t.Parallel()

	clean := NewDiagram(nil).LF()
	if !strings.HasSuffix(clean.String(), "\n") {
		t.Errorf("LF did not append a blank line: %q", clean.String())
	}

	broken := NewDiagram(nil).XAxisRangeWithTitle("Month", 12, 12)
	before := broken.String()
	if after := broken.LF().String(); after != before {
		t.Errorf("LF appended after an error:\n before: %q\n  after: %q", before, after)
	}
}

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).Line(5000, 6000)
	})
}

// TestRecordedErrorContract asserts that an empty set of x-axis labels surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).XAxisLabels()
	})
}

// TestGoldenXYChart pins the rendered chart of every axis form and every series
// kind this package can build.
func TestGoldenXYChart(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(buf, WithTitle("Sales Revenue")).
		XAxisLabelsWithTitle("Month", "jan", "feb", "mar").
		YAxisRangeWithTitle("Revenue (in $)", 4000, 11000).
		Bar(5000, 6000, 7500).
		Line(5000, 6000, 7500).
		LF().
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("xychart.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenXYChartAxisWithoutTitle pins the axis forms that take no title, and
// the numeric x axis range.
func TestGoldenXYChartAxisWithoutTitle(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(buf).
		XAxisRange(1, 12).
		YAxisRange(0, 100).
		Line(10, 20, 30).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("xychart_ranges.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenXYChartTitledRanges pins the axis forms that carry a title and a
// numeric range, which is the one combination the tests above leave out.
func TestGoldenXYChartTitledRanges(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(buf).
		XAxisRangeWithTitle("Month", 1, 12).
		YAxisRangeWithTitle("Revenue (in $)", 0, 11000).
		Bar(5000, 6000).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("xychart_titled_ranges.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenXYChartOrientations pins the header each orientation produces.
func TestGoldenXYChartOrientations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		golden string
		option Option
	}{
		{name: "vertical", golden: "xychart_vertical.md", option: WithOrientation(OrientationVertical)},
		{name: "horizontal", golden: "xychart_horizontal.md", option: WithOrientation(OrientationHorizontal)},
		{name: "horizontal shorthand", golden: "xychart_horizontal_shorthand.md", option: WithHorizontal()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			err := NewDiagram(buf, tt.option).
				XAxisLabels("a", "b").
				Bar(1, 2).
				Build()
			if err != nil {
				t.Fatalf("Build() = %v, want nil", err)
			}

			if err := golden.Assert(tt.golden, buf.String()); err != nil {
				t.Error(err)
			}
		})
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
