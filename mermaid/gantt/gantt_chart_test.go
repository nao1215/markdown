package gantt

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

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewChart(w).Section("Planning").Task("Design", "2024-01-01", "2d")
	})
}

// TestChartError covers the accessor callers use to find out why a chart came
// out wrong, which nothing exercised.
func TestChartError(t *testing.T) {
	t.Parallel()

	t.Run("a well formed chart reports no error", func(t *testing.T) {
		t.Parallel()

		c := NewChart(nil).Section("build").Task("compile", "2024-01-01", "1d")
		if err := c.Error(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("a nil writer is reported through Error after Build", func(t *testing.T) {
		t.Parallel()

		c := NewChart(nil)
		if err := c.Build(); err == nil {
			t.Fatal("Build with a nil writer must fail")
		}
		if c.Error() == nil {
			t.Error("Error must report the failure Build returned")
		}
		if !strings.Contains(c.Error().Error(), "nil") {
			t.Errorf("unexpected error: %v", c.Error())
		}
	})
}

// TestGoldenGanttChart pins the rendered chart of every option and every task
// kind this package can build.
func TestGoldenGanttChart(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewChart(
		buf,
		WithTitle("Every Task Kind"),
		WithDateFormat("YYYY-MM-DD"),
		WithAxisFormat("%m-%d"),
		WithTickInterval("1week"),
		WithExcludes("weekends", "2024-01-01"),
		WithTodayMarker("stroke-width:4px"),
	).
		Section("Plain tasks").
		Task("Task", "2024-01-01", "2d").
		TaskWithID("Task with id", "task-id", "2024-01-03", "2d").
		TaskAfter("Task after", "task-id", "1d").
		TaskAfterWithID("Task after with id", "after-id", "task-id", "1d").
		Section("Marked tasks").
		CriticalTask("Critical", "2024-01-06", "1d").
		CriticalTaskWithID("Critical with id", "crit-id", "2024-01-07", "1d").
		ActiveTask("Active", "2024-01-08", "1d").
		ActiveTaskWithID("Active with id", "active-id", "2024-01-09", "1d").
		DoneTask("Done", "2024-01-10", "1d").
		DoneTaskWithID("Done with id", "done-id", "2024-01-11", "1d").
		CriticalActiveTask("Critical active", "2024-01-12", "1d").
		CriticalActiveTaskWithID("Critical active with id", "crit-active-id", "2024-01-13", "1d").
		CriticalDoneTask("Critical done", "2024-01-14", "1d").
		CriticalDoneTaskWithID("Critical done with id", "crit-done-id", "2024-01-15", "1d").
		LF().
		Section("Milestones").
		Milestone("Milestone", "2024-01-16").
		MilestoneWithID("Milestone with id", "milestone-id", "2024-01-17").
		CriticalMilestone("Critical milestone", "2024-01-18").
		CriticalMilestoneWithID("Critical milestone with id", "crit-milestone-id", "2024-01-19").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("gantt.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestBuildWithNilWriter covers the case where a chart is built for String()
// only and Build() is called by mistake. Build() used to dereference the nil
// writer and take the process down; it has to return an error instead.
func TestBuildWithNilWriter(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Build() panicked with a nil writer: %v", r)
		}
	}()

	d := NewChart(nil)

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

	err := NewChart(errWriter{}).Build()
	if err == nil {
		t.Fatal("Build must report a failing writer")
	}
	if !errors.Is(err, errWrite) {
		t.Errorf("Build lost the destination error: %v", err)
	}
}

// TestTaskNameEscapesTheColonThatEndsIt names the character this escaping buys.
// A gantt task is written "name :data" and mermaid splits the line at the first
// colon, so a colon in the name used to end it early: the chart still drew, with
// a task called "Deploy" where the caller asked for "Deploy: staging". Nothing
// reported it, which is what makes it worth fixing.
func TestTaskNameEscapesTheColonThatEndsIt(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(*Chart) *Chart
		want  string
	}{
		"a colon in a plain task": {
			build: func(g *Chart) *Chart {
				return g.Task("Deploy: staging", "2024-01-01", "2d")
			},
			want: "    Deploy#58; staging :2024-01-01, 2d",
		},
		"a colon in a critical task": {
			build: func(g *Chart) *Chart {
				return g.CriticalTask("a:b", "2024-01-01", "2d")
			},
			want: "    a#58;b :crit, 2024-01-01, 2d",
		},
		"a colon in a milestone": {
			build: func(g *Chart) *Chart {
				return g.Milestone("a:b", "2024-01-01")
			},
			want: "    a#58;b :milestone, 2024-01-01, 0d",
		},
		"a colon in a task that follows another": {
			build: func(g *Chart) *Chart {
				return g.TaskAfter("a:b", "t1", "2d")
			},
			want: "    a#58;b :after t1, 2d",
		},
		"a named entity in a task name is escaped": {
			build: func(g *Chart) *Chart {
				return g.Task("a#quot;b", "2024-01-01", "2d")
			},
			want: "    a#35;quot;b :2024-01-01, 2d",
		},
		"a plain hash in a task name is left alone": {
			build: func(g *Chart) *Chart {
				return g.Task("PR #123 merged", "2024-01-01", "2d")
			},
			want: "    PR #123 merged :2024-01-01, 2d",
		},
		"a line break in a task name becomes the break mermaid draws": {
			// Raw it ended the statement and handed the second line to the
			// parser as a line of its own.
			build: func(g *Chart) *Chart {
				return g.Task("first\nsecond", "2024-01-01", "2d")
			},
			want: "    first<br/>second :2024-01-01, 2d",
		},
		"a line break in a section name becomes the break mermaid draws": {
			build: func(g *Chart) *Chart {
				return g.Section("first\r\nsecond")
			},
			want: "    section first<br/>second",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.build(NewChart(io.Discard)).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("chart =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}

// TestSectionAndTitleKeepTheirColon pins the other half: mermaid reads each of
// those as the rest of its line, so a colon in one already reaches the drawing
// and escaping it would change output that is correct today.
func TestSectionAndTitleKeepTheirColon(t *testing.T) {
	t.Parallel()

	got := NewChart(io.Discard, WithTitle("Q1: plan")).Section("Ops: core").String()

	for _, want := range []string{"    title Q1: plan", "    section Ops: core"} {
		if !strings.Contains(got, want) {
			t.Errorf("chart =\n%s\nwant it to contain\n%s", got, want)
		}
	}
}

// TestTitleEscapesAngleAndLineBreak covers the two characters a gantt title
// cannot carry as they are: a bare "<" is eaten by the sanitizer, and a line
// break splits the statement. Both entity forms decode in the drawing.
func TestTitleEscapesAngleAndLineBreak(t *testing.T) {
	t.Parallel()

	got := NewChart(io.Discard, WithTitle("cost < 10\nper unit")).String()

	if want := "    title cost #60; 10#10;per unit"; !strings.Contains(got, want) {
		t.Errorf("chart =\n%s\nwant it to contain\n%s", got, want)
	}
}
