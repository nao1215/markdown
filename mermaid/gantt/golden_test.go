package gantt_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/gantt"
)

// TestGoldenGanttChart pins the rendered chart of every option and every task
// kind this package can build.
func TestGoldenGanttChart(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := gantt.NewChart(
		buf,
		gantt.WithTitle("Every Task Kind"),
		gantt.WithDateFormat("YYYY-MM-DD"),
		gantt.WithAxisFormat("%m-%d"),
		gantt.WithTickInterval("1week"),
		gantt.WithExcludes("weekends", "2024-01-01"),
		gantt.WithTodayMarker("stroke-width:4px"),
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
