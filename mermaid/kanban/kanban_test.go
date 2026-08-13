package kanban

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
			want: "kanban",
		},
		{
			name: "new diagram with title",
			opts: []Option{WithTitle("Sprint Board")},
			want: `---
title: "Sprint Board"
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
title: "Sprint Board"
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
title: "Sprint Board"
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
        [Fix parser]@{ ticket: 'KB-\\123', assigned: 'O''Reilly' }`

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

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).Column("Todo").Task("Define scope")
	})
}

// TestRecordedErrorContract asserts that a task added before any column surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).Task("a task with no column")
	})
}

// TestGoldenKanban pins the rendered diagram of every builder method, every
// option, and every priority of this package.
func TestGoldenKanban(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(
		buf,
		WithTitle("Sprint Board"),
		WithTicketBaseURL("https://example.com/tickets/"),
	).
		Column("Todo").
		Task("Plain task").
		Task("Task with metadata",
			WithTaskID("task-1"),
			WithTaskTicket("MB-101"),
			WithTaskAssigned("Alice"),
			WithTaskPriority(PriorityVeryHigh),
		).
		Column("In Progress", WithColumnID("in-progress")).
		Task("High", WithTaskPriority(PriorityHigh)).
		Task("Low", WithTaskPriority(PriorityLow)).
		Task("Very low", WithTaskPriority(PriorityVeryLow)).
		LF().
		Column("Done").
		TaskIn("Done", "Task added by column name", WithTaskAssigned("Bob")).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("kanban.md", buf.String()); err != nil {
		t.Error(err)
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

// TestCardLabelEscapesWhatEndsIt names the characters this escaping buys in a
// card. A card is written "id[label]" and mermaid ends the label at a closing
// bracket, a parenthesis or a closing brace, so a task called "Fix parser (v2)"
// lost the whole diagram. A bracket used to be rejected outright instead, which
// is not this library's job: a caller passes data and the builder encodes it.
func TestCardLabelEscapesWhatEndsIt(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(io.Writer) *Diagram
		want  string
	}{
		"parentheses in a column": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Column("Todo (soon)") },
			want:  "    [Todo #40;soon#41;]",
		},
		"a closing bracket in a task": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).Column("Todo").Task("Fix parser]")
			},
			want: "        [Fix parser#93;]",
		},
		"a closing brace in a task": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).Column("Todo").Task("a}b")
			},
			want: "        [a#125;b]",
		},
		"an opening bracket and brace are left alone": {
			// mermaid takes either inside a label, so escaping them would
			// change output that already reaches the drawing.
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).Column("Todo").Task("a[b{c")
			},
			want: "        [a[b{c]",
		},
		"a quote needs nothing in a label": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).Column("Todo").Task(`say "hi" to O'Reilly`)
			},
			want: `        [say "hi" to O'Reilly]`,
		},
		"a named entity is escaped": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Column("a#93;b") },
			want:  "    [a#35;93;b]",
		},
		"a plain hash is left alone": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Column("PR #123") },
			want:  "    [PR #123]",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.build(io.Discard).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("diagram =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}

// TestMetadataEscapingMeetsThreeReaders names what a metadata value has to get
// past. mermaid reads "@{ ticket: 'value' }" as YAML and then draws the result,
// so a single quote is doubled the way YAML wants it, while the punctuation the
// kanban lexer takes before YAML ever sees the line is written as a mermaid
// entity. Writing "\'" for the quote, as this package used to, makes the YAML
// parser refuse the line and lose the diagram.
func TestMetadataEscapingMeetsThreeReaders(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		value string
		want  string
	}{
		"a single quote is doubled":     {value: "O'Reilly", want: "ticket: 'O''Reilly'"},
		"a double quote is an entity":   {value: `say "hi"`, want: `ticket: 'say #quot;hi#quot;'`},
		"a closing brace is an entity":  {value: "a}b", want: "ticket: 'a#125;b'"},
		"a backslash keeps its escape":  {value: `KB-\123`, want: `ticket: 'KB-\\123'`},
		"an opening brace is left":      {value: "a{b", want: "ticket: 'a{b'"},
		"a plain hash is left alone":    {value: "PR #123", want: "ticket: 'PR #123'"},
		"a named entity gets escaped":   {value: "a#125;b", want: "ticket: 'a#35;125;b'"},
		"ordinary text is left as text": {value: "KB-1", want: "ticket: 'KB-1'"},
		"a caret is an entity": {
			// The kanban lexer refuses a bare caret before YAML sees the line,
			// measured by rendering: "assigned: 'x^x'" lost the whole board.
			value: "a^b", want: "ticket: 'a#94;b'",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := NewDiagram(io.Discard).Column("Todo").Task("t", WithTaskTicket(tt.value)).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("diagram =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}
