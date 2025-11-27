package gantt

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts []Option
		want string
	}{
		{
			name: "simple chart without options",
			opts: nil,
			want: "gantt",
		},
		{
			name: "chart with title",
			opts: []Option{WithTitle("Project Schedule")},
			want: `gantt
    title Project Schedule`,
		},
		{
			name: "chart with all options",
			opts: []Option{
				WithTitle("Project Schedule"),
				WithDateFormat("YYYY-MM-DD"),
				WithAxisFormat("%Y-%m-%d"),
				WithTickInterval("1week"),
				WithTodayMarker("off"),
				WithExcludes("weekends", "2024-01-01"),
			},
			want: `gantt
    title Project Schedule
    dateFormat YYYY-MM-DD
    axisFormat %Y-%m-%d
    tickInterval 1week
    todayMarker off
    excludes weekends
    excludes 2024-01-01`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewChart(io.Discard, tt.opts...)
			got := strings.ReplaceAll(c.String(), "\r\n", "\n")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestChart_Section(t *testing.T) {
	t.Parallel()

	c := NewChart(io.Discard).
		Section("Phase 1")

	want := `gantt
    section Phase 1`

	got := strings.ReplaceAll(c.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestChart_Task(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func(*Chart) *Chart
		want string
	}{
		{
			name: "simple task",
			fn: func(c *Chart) *Chart {
				return c.Task("Task 1", "2024-01-01", "30d")
			},
			want: `gantt
    Task 1 :2024-01-01, 30d`,
		},
		{
			name: "task with ID",
			fn: func(c *Chart) *Chart {
				return c.TaskWithID("Task 1", "task1", "2024-01-01", "30d")
			},
			want: `gantt
    Task 1 :task1, 2024-01-01, 30d`,
		},
		{
			name: "critical task",
			fn: func(c *Chart) *Chart {
				return c.CriticalTask("Critical Task", "2024-01-01", "7d")
			},
			want: `gantt
    Critical Task :crit, 2024-01-01, 7d`,
		},
		{
			name: "critical task with ID",
			fn: func(c *Chart) *Chart {
				return c.CriticalTaskWithID("Critical Task", "crit1", "2024-01-01", "7d")
			},
			want: `gantt
    Critical Task :crit, crit1, 2024-01-01, 7d`,
		},
		{
			name: "active task",
			fn: func(c *Chart) *Chart {
				return c.ActiveTask("Active Task", "2024-01-01", "5d")
			},
			want: `gantt
    Active Task :active, 2024-01-01, 5d`,
		},
		{
			name: "active task with ID",
			fn: func(c *Chart) *Chart {
				return c.ActiveTaskWithID("Active Task", "active1", "2024-01-01", "5d")
			},
			want: `gantt
    Active Task :active, active1, 2024-01-01, 5d`,
		},
		{
			name: "done task",
			fn: func(c *Chart) *Chart {
				return c.DoneTask("Done Task", "2024-01-01", "3d")
			},
			want: `gantt
    Done Task :done, 2024-01-01, 3d`,
		},
		{
			name: "done task with ID",
			fn: func(c *Chart) *Chart {
				return c.DoneTaskWithID("Done Task", "done1", "2024-01-01", "3d")
			},
			want: `gantt
    Done Task :done, done1, 2024-01-01, 3d`,
		},
		{
			name: "critical active task",
			fn: func(c *Chart) *Chart {
				return c.CriticalActiveTask("Critical Active", "2024-01-01", "2d")
			},
			want: `gantt
    Critical Active :crit, active, 2024-01-01, 2d`,
		},
		{
			name: "critical active task with ID",
			fn: func(c *Chart) *Chart {
				return c.CriticalActiveTaskWithID("Critical Active", "ca1", "2024-01-01", "2d")
			},
			want: `gantt
    Critical Active :crit, active, ca1, 2024-01-01, 2d`,
		},
		{
			name: "critical done task",
			fn: func(c *Chart) *Chart {
				return c.CriticalDoneTask("Critical Done", "2024-01-01", "1d")
			},
			want: `gantt
    Critical Done :crit, done, 2024-01-01, 1d`,
		},
		{
			name: "critical done task with ID",
			fn: func(c *Chart) *Chart {
				return c.CriticalDoneTaskWithID("Critical Done", "cd1", "2024-01-01", "1d")
			},
			want: `gantt
    Critical Done :crit, done, cd1, 2024-01-01, 1d`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewChart(io.Discard)
			tt.fn(c)
			got := strings.ReplaceAll(c.String(), "\r\n", "\n")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestChart_Milestone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func(*Chart) *Chart
		want string
	}{
		{
			name: "milestone",
			fn: func(c *Chart) *Chart {
				return c.Milestone("Release", "2024-01-15")
			},
			want: `gantt
    Release :milestone, 2024-01-15, 0d`,
		},
		{
			name: "milestone with ID",
			fn: func(c *Chart) *Chart {
				return c.MilestoneWithID("Release", "rel1", "2024-01-15")
			},
			want: `gantt
    Release :milestone, rel1, 2024-01-15, 0d`,
		},
		{
			name: "critical milestone",
			fn: func(c *Chart) *Chart {
				return c.CriticalMilestone("Critical Release", "2024-01-15")
			},
			want: `gantt
    Critical Release :crit, milestone, 2024-01-15, 0d`,
		},
		{
			name: "critical milestone with ID",
			fn: func(c *Chart) *Chart {
				return c.CriticalMilestoneWithID("Critical Release", "crel1", "2024-01-15")
			},
			want: `gantt
    Critical Release :crit, milestone, crel1, 2024-01-15, 0d`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewChart(io.Discard)
			tt.fn(c)
			got := strings.ReplaceAll(c.String(), "\r\n", "\n")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestChart_TaskAfter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func(*Chart) *Chart
		want string
	}{
		{
			name: "task after",
			fn: func(c *Chart) *Chart {
				return c.TaskWithID("Task 1", "task1", "2024-01-01", "5d").
					TaskAfter("Task 2", "task1", "3d")
			},
			want: `gantt
    Task 1 :task1, 2024-01-01, 5d
    Task 2 :after task1, 3d`,
		},
		{
			name: "task after with ID",
			fn: func(c *Chart) *Chart {
				return c.TaskWithID("Task 1", "task1", "2024-01-01", "5d").
					TaskAfterWithID("Task 2", "task2", "task1", "3d")
			},
			want: `gantt
    Task 1 :task1, 2024-01-01, 5d
    Task 2 :task2, after task1, 3d`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewChart(io.Discard)
			tt.fn(c)
			got := strings.ReplaceAll(c.String(), "\r\n", "\n")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestChart_Build(t *testing.T) {
	t.Parallel()

	t.Run("success build", func(t *testing.T) {
		t.Parallel()

		buf := new(bytes.Buffer)
		c := NewChart(buf, WithTitle("Test"))
		err := c.Build()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := `gantt
    title Test`
		got := strings.ReplaceAll(buf.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestChart_ComplexExample(t *testing.T) {
	t.Parallel()

	c := NewChart(io.Discard,
		WithTitle("Project Plan"),
		WithDateFormat("YYYY-MM-DD"),
	).
		Section("Planning").
		DoneTaskWithID("Requirements", "req", "2024-01-01", "7d").
		DoneTaskWithID("Design", "design", "2024-01-08", "5d").
		LF().
		Section("Development").
		CriticalActiveTaskWithID("Implementation", "impl", "2024-01-15", "14d").
		TaskAfterWithID("Testing", "test", "impl", "7d").
		LF().
		Section("Deployment").
		TaskAfter("Deploy", "test", "2d").
		CriticalMilestone("Go Live", "2024-02-10")

	want := `gantt
    title Project Plan
    dateFormat YYYY-MM-DD
    section Planning
    Requirements :done, req, 2024-01-01, 7d
    Design :done, design, 2024-01-08, 5d

    section Development
    Implementation :crit, active, impl, 2024-01-15, 14d
    Testing :test, after impl, 7d

    section Deployment
    Deploy :after test, 2d
    Go Live :crit, milestone, 2024-02-10, 0d`

	got := strings.ReplaceAll(c.String(), "\r\n", "\n")

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}
