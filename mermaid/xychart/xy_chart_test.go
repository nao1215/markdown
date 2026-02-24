package xychart

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
