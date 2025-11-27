package quadrant

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestChart_Build(t *testing.T) {
	t.Parallel()

	t.Run("Build a simple quadrant chart", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		ch := NewChart(b)
		ch.XAxis("Low", "High").
			YAxis("Low", "High").
			Quadrant1("Promote").
			Quadrant2("Plan").
			Quadrant3("Eliminate").
			Quadrant4("Delegate").
			Point("Task A", 0.8, 0.9).
			Point("Task B", 0.3, 0.7).
			Point("Task C", 0.2, 0.3).
			Point("Task D", 0.7, 0.2)

		if err := ch.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `quadrantChart
    x-axis Low --> High
    y-axis Low --> High
    quadrant-1 Promote
    quadrant-2 Plan
    quadrant-3 Eliminate
    quadrant-4 Delegate
    Task A: [0.80, 0.90]
    Task B: [0.30, 0.70]
    Task C: [0.20, 0.30]
    Task D: [0.70, 0.20]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a quadrant chart with title", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		ch := NewChart(b, WithTitle("Priority Matrix"))
		ch.XAxis("Effort", "").
			YAxis("Impact", "").
			Point("Feature A", 0.5, 0.8)

		if err := ch.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `quadrantChart
    title Priority Matrix
    x-axis Effort
    y-axis Impact
    Feature A: [0.50, 0.80]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a quadrant chart with styled points", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		ch := NewChart(b)
		ch.XAxis("Low Priority", "High Priority").
			YAxis("Low Risk", "High Risk").
			Quadrant1("Critical").
			Quadrant2("Caution").
			Quadrant3("Safe").
			Quadrant4("Monitor").
			PointWithStyle("Project A", 0.9, 0.8, "radius: 12").
			PointWithStyle("Project B", 0.2, 0.9, "color: #ff0000")

		if err := ch.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `quadrantChart
    x-axis Low Priority --> High Priority
    y-axis Low Risk --> High Risk
    quadrant-1 Critical
    quadrant-2 Caution
    quadrant-3 Safe
    quadrant-4 Monitor
    Project A: [0.90, 0.80] radius: 12
    Project B: [0.20, 0.90] color: #ff0000`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a quadrant chart with line feeds", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		ch := NewChart(b, WithTitle("Eisenhower Matrix"))
		ch.XAxis("Not Urgent", "Urgent").
			YAxis("Not Important", "Important").
			LF().
			Quadrant1("Do First").
			Quadrant2("Schedule").
			Quadrant3("Eliminate").
			Quadrant4("Delegate").
			LF().
			Point("Task 1", 0.9, 0.9).
			Point("Task 2", 0.1, 0.9)

		if err := ch.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `quadrantChart
    title Eisenhower Matrix
    x-axis Not Urgent --> Urgent
    y-axis Not Important --> Important

    quadrant-1 Do First
    quadrant-2 Schedule
    quadrant-3 Eliminate
    quadrant-4 Delegate

    Task 1: [0.90, 0.90]
    Task 2: [0.10, 0.90]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a quadrant chart with class definitions", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		ch := NewChart(b, WithTitle("Reach and engagement of campaigns"))
		ch.XAxis("Low Reach", "High Reach").
			YAxis("Low Engagement", "High Engagement").
			Quadrant1("We should expand").
			Quadrant2("Need to promote").
			Quadrant3("Re-evaluate").
			Quadrant4("May be improved").
			PointWithStyle("Campaign A", 0.9, 0.0, "radius: 12").
			PointWithClassAndStyle("Campaign B", 0.8, 0.1, "class1", "color: #ff3300, radius: 10").
			PointWithStyle("Campaign C", 0.7, 0.2, "radius: 25, color: #00ff33, stroke-color: #10f0f0").
			PointWithStyle("Campaign D", 0.6, 0.3, "radius: 15, stroke-color: #00ff0f, stroke-width: 5px, color: #ff33f0").
			PointWithClass("Campaign E", 0.5, 0.4, "class2").
			PointWithClassAndStyle("Campaign F", 0.4, 0.5, "class3", "color: #0000ff").
			ClassDef("class1", "color: #109060").
			ClassDef("class2", "color: #908342, radius: 10, stroke-color: #310085, stroke-width: 10px").
			ClassDef("class3", "color: #f00fff, radius: 10")

		if err := ch.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `quadrantChart
    title Reach and engagement of campaigns
    x-axis Low Reach --> High Reach
    y-axis Low Engagement --> High Engagement
    quadrant-1 We should expand
    quadrant-2 Need to promote
    quadrant-3 Re-evaluate
    quadrant-4 May be improved
    Campaign A: [0.90, 0.00] radius: 12
    Campaign B:::class1: [0.80, 0.10] color: #ff3300, radius: 10
    Campaign C: [0.70, 0.20] radius: 25, color: #00ff33, stroke-color: #10f0f0
    Campaign D: [0.60, 0.30] radius: 15, stroke-color: #00ff0f, stroke-width: 5px, color: #ff33f0
    Campaign E:::class2: [0.50, 0.40]
    Campaign F:::class3: [0.40, 0.50] color: #0000ff
    classDef class1 color: #109060
    classDef class2 color: #908342, radius: 10, stroke-color: #310085, stroke-width: 10px
    classDef class3 color: #f00fff, radius: 10`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a quadrant chart with PointStyled", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		ch := NewChart(b)
		ch.XAxis("X", "").
			YAxis("Y", "").
			PointStyled("Point A", 0.9, 0.8, PointStyle{
				Color:       "#ff0000",
				Radius:      12,
				StrokeColor: "#00ff00",
				StrokeWidth: "5px",
			}).
			PointStyled("Point B", 0.5, 0.5, PointStyle{
				Color:  "#0000ff",
				Radius: 8,
			}).
			PointStyled("Point C", 0.2, 0.2, PointStyle{})

		if err := ch.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `quadrantChart
    x-axis X
    y-axis Y
    Point A: [0.90, 0.80] color: #ff0000, radius: 12, stroke-color: #00ff00, stroke-width: 5px
    Point B: [0.50, 0.50] color: #0000ff, radius: 8
    Point C: [0.20, 0.20]`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a quadrant chart with ClassDefStyled", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		ch := NewChart(b)
		ch.XAxis("X", "").
			YAxis("Y", "").
			PointWithClass("Point A", 0.9, 0.8, "myClass").
			ClassDefStyled("myClass", ClassStyle{
				Color:       "#ff0000",
				Radius:      15,
				StrokeColor: "#00ff00",
				StrokeWidth: "3px",
			})

		if err := ch.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `quadrantChart
    x-axis X
    y-axis Y
    Point A:::myClass: [0.90, 0.80]
    classDef myClass color: #ff0000, radius: 15, stroke-color: #00ff00, stroke-width: 3px`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestChart_String(t *testing.T) {
	t.Parallel()

	t.Run("String returns the quadrant chart body", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		ch := NewChart(b)
		ch.XAxis("X", "").YAxis("Y", "").Point("A", 0.5, 0.5)

		want := `quadrantChart
    x-axis X
    y-axis Y
    A: [0.50, 0.50]`

		got := strings.ReplaceAll(ch.String(), "\r\n", "\n")
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestChart_Error(t *testing.T) {
	t.Parallel()

	t.Run("Error returns nil when no error", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)
		ch := NewChart(b)

		if err := ch.Error(); err != nil {
			t.Errorf("expected nil, got %v", err)
		}
	})
}

func TestPointStyle_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		style PointStyle
		want  string
	}{
		{
			name:  "empty style",
			style: PointStyle{},
			want:  "",
		},
		{
			name: "color only",
			style: PointStyle{
				Color: "#ff0000",
			},
			want: "color: #ff0000",
		},
		{
			name: "radius only",
			style: PointStyle{
				Radius: 10,
			},
			want: "radius: 10",
		},
		{
			name: "all properties",
			style: PointStyle{
				Color:       "#ff0000",
				Radius:      12,
				StrokeColor: "#00ff00",
				StrokeWidth: "5px",
			},
			want: "color: #ff0000, radius: 12, stroke-color: #00ff00, stroke-width: 5px",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.style.String()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):%s", diff)
			}
		})
	}
}

func TestClassStyle_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		style ClassStyle
		want  string
	}{
		{
			name:  "empty style",
			style: ClassStyle{},
			want:  "",
		},
		{
			name: "color only",
			style: ClassStyle{
				Color: "#ff0000",
			},
			want: "color: #ff0000",
		},
		{
			name: "all properties",
			style: ClassStyle{
				Color:       "#ff0000",
				Radius:      12,
				StrokeColor: "#00ff00",
				StrokeWidth: "5px",
			},
			want: "color: #ff0000, radius: 12, stroke-color: #00ff00, stroke-width: 5px",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.style.String()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):%s", diff)
			}
		})
	}
}
