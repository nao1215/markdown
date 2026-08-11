package kanban_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/kanban"
)

// TestGoldenKanban pins the rendered diagram of every builder method, every
// option, and every priority of this package.
func TestGoldenKanban(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := kanban.NewDiagram(
		buf,
		kanban.WithTitle("Sprint Board"),
		kanban.WithTicketBaseURL("https://example.com/tickets/"),
	).
		Column("Todo").
		Task("Plain task").
		Task("Task with metadata",
			kanban.WithTaskID("task-1"),
			kanban.WithTaskTicket("MB-101"),
			kanban.WithTaskAssigned("Alice"),
			kanban.WithTaskPriority(kanban.PriorityVeryHigh),
		).
		Column("In Progress", kanban.WithColumnID("in-progress")).
		Task("High", kanban.WithTaskPriority(kanban.PriorityHigh)).
		Task("Low", kanban.WithTaskPriority(kanban.PriorityLow)).
		Task("Very low", kanban.WithTaskPriority(kanban.PriorityVeryLow)).
		LF().
		Column("Done").
		TaskIn("Done", "Task added by column name", kanban.WithTaskAssigned("Bob")).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("kanban.md", buf.String()); err != nil {
		t.Error(err)
	}
}
