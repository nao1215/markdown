package kanban

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
			want: "kanban",
		},
		{
			name: "new diagram with title",
			opts: []Option{WithTitle("Sprint Board")},
			want: `---
title: Sprint Board
---
kanban`,
		},
		{
			name: "new diagram with ticket base URL",
			opts: []Option{WithTicketBaseURL("https://example.com/tickets/")},
			want: `---
config:
  kanban:
    ticketBaseUrl: 'https://example.com/tickets/'
---
kanban`,
		},
		{
			name: "new diagram with title and ticket base URL",
			opts: []Option{
				WithTitle("Sprint Board"),
				WithTicketBaseURL("https://example.com/tickets/"),
			},
			want: `---
title: Sprint Board
config:
  kanban:
    ticketBaseUrl: 'https://example.com/tickets/'
---
kanban`,
		},
		{
			name: "new diagram with single quote in ticket base URL",
			opts: []Option{
				WithTicketBaseURL("https://example.com/o'hare"),
			},
			want: `---
config:
  kanban:
    ticketBaseUrl: 'https://example.com/o''hare'
---
kanban`,
		},
		{
			name:    "new diagram with title including newline",
			opts:    []Option{WithTitle("Sprint\nBoard")},
			want:    "kanban",
			wantErr: true,
		},
		{
			name:    "new diagram with ticket base URL including newline",
			opts:    []Option{WithTicketBaseURL("https://example.com/\ntickets/")},
			want:    "kanban",
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

	d := NewDiagram(
		b,
		WithTitle("Sprint Board"),
		WithTicketBaseURL("https://example.com/tickets/"),
	).
		Column("Todo", WithColumnID("todo")).
		Task("Define scope").
		Task(
			"Create login page",
			WithTaskID("k1"),
			WithTaskTicket("MB-101"),
			WithTaskAssigned("Alice"),
			WithTaskPriority(PriorityHigh),
		).
		LF().
		Column("In Progress", WithColumnID("doing")).
		Task("Review API", WithTaskPriority(PriorityVeryHigh))

	if err := d.Build(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `---
title: Sprint Board
config:
  kanban:
    ticketBaseUrl: 'https://example.com/tickets/'
---
kanban
    todo[Todo]
        [Define scope]
        k1[Create login page]@{ ticket: 'MB-101', assigned: 'Alice', priority: 'High' }

    doing[In Progress]
        [Review API]@{ priority: 'Very High' }`

	got := strings.ReplaceAll(b.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_TaskIn(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		TaskIn("Todo", "Design UI").
		TaskIn("Todo", "Implement UI", WithTaskPriority(PriorityLow)).
		TaskIn("Done", "Ship release")

	want := `kanban
    [Todo]
        [Design UI]
        [Implement UI]@{ priority: 'Low' }
    [Done]
        [Ship release]`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_BracketedNamesAndMetadataEscaping(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		Column("[Todo]").
		Task(
			"[Fix parser]",
			WithTaskAssigned("O'Reilly"),
			WithTaskTicket(`KB-\123`),
		)

	want := `kanban
    [Todo]
        [Fix parser]@{ ticket: 'KB-\\123', assigned: 'O\'Reilly' }`

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
			name: "task before column",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Task("Define scope")
			},
			want: "kanban",
		},
		{
			name: "empty column name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Column("")
			},
			want: "kanban",
		},
		{
			name: "column name with newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Column("To\ndo")
			},
			want: "kanban",
		},
		{
			name: "column name with bracket char",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Column("To[do]")
			},
			want: "kanban",
		},
		{
			name: "column id with whitespace",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Column("Todo", WithColumnID("to do"))
			},
			want: "kanban",
		},
		{
			name: "empty task name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Column("Todo").Task("")
			},
			want: `kanban
    [Todo]`,
		},
		{
			name: "task name with newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Column("Todo").Task("Define\nscope")
			},
			want: `kanban
    [Todo]`,
		},
		{
			name: "task name with bracket char",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Column("Todo").Task("Fix [parser]")
			},
			want: `kanban
    [Todo]`,
		},
		{
			name: "task id with whitespace",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Column("Todo").Task("Define", WithTaskID("k 1"))
			},
			want: `kanban
    [Todo]`,
		},
		{
			name: "task ticket with newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Column("Todo").
					Task("Define", WithTaskTicket("MB-\n1"))
			},
			want: `kanban
    [Todo]`,
		},
		{
			name: "task assignee with newline",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Column("Todo").
					Task("Define", WithTaskAssigned("Ali\nce"))
			},
			want: `kanban
    [Todo]`,
		},
		{
			name: "invalid task priority",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Column("Todo").
					Task("Define", WithTaskPriority(Priority("urgent")))
			},
			want: `kanban
    [Todo]`,
		},
		{
			name: "taskin with empty column",
			run: func() *Diagram {
				return NewDiagram(io.Discard).TaskIn("", "Define")
			},
			want: "kanban",
		},
		{
			name: "lf short-circuit after error",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Column("").LF()
			},
			want: "kanban",
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
