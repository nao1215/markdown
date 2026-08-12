//go:build linux || darwin

package kanban_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/kanban"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := kanban.NewDiagram(
		io.Discard,
		kanban.WithTitle("Sprint Board"),
		kanban.WithTicketBaseURL("https://example.com/tickets/"),
	).
		Column("Todo").
		Task("Define scope").
		Task(
			"Create login page",
			kanban.WithTaskTicket("MB-101"),
			kanban.WithTaskAssigned("Alice"),
			kanban.WithTaskPriority(kanban.PriorityHigh),
		).
		Column("In Progress").
		Task("Review API", kanban.WithTaskPriority(kanban.PriorityVeryHigh)).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Kanban Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Kanban Diagram
	// ```mermaid
	// ---
	// title: "Sprint Board"
	// config:
	//   kanban:
	//     ticketBaseUrl: 'https://example.com/tickets/'
	// ---
	// kanban
	//     [Todo]
	//         [Define scope]
	//         [Create login page]@{ ticket: 'MB-101', assigned: 'Alice', priority: 'High' }
	//     [In Progress]
	//         [Review API]@{ priority: 'Very High' }
	// ```
}

// ExampleNewDiagram shows the shape every kanban board has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo").Task("Write the spec").
		Build()

	// Output:
	// kanban
	//     [Todo]
	//         [Write the spec]
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := kanban.NewDiagram(io.Discard).
		Column("Todo").Task("Write the spec").
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// kanban
	//     [Todo]
	//         [Write the spec]
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := kanban.NewDiagram(nil).
		Column("Todo").Task("Write the spec").
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := kanban.NewDiagram(io.Discard).
		Task("no column yet")
	fmt.Println("error:", d.Error())

	// Output:
	// error: task "no column yet" requires a column; call Column first
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo").Task("Write the spec").
		LF().
		Column("Todo").Task("Write the spec").
		Build()

	// Output:
	// kanban
	//     [Todo]
	//         [Write the spec]
	//
	//     [Todo]
	//         [Write the spec]
}

// ExampleDiagram_full shows a kanban board built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := kanban.NewDiagram(io.Discard).
		Column("Todo").Task("Write the spec").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// kanban
	//     [Todo]
	//         [Write the spec]
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []kanban.Option{kanban.WithTitle("Overview")}

	_ = kanban.NewDiagram(os.Stdout, options...).
		Column("Todo").Task("Write the spec").
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// kanban
	//     [Todo]
	//         [Write the spec]
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = kanban.NewDiagram(os.Stdout, kanban.WithTitle("Overview")).
		Column("Todo").Task("Write the spec").
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// kanban
	//     [Todo]
	//         [Write the spec]
}

// ExampleDiagram_Column adds a column, and the tasks that follow belong to it.
// A board needs one before any task, and a task added without one is reported
// from Build.
func ExampleDiagram_Column() {
	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo").
		Task("Write the spec").
		Column("Done").
		Task("Read the spec").
		Build()

	// Output:
	// kanban
	//     [Todo]
	//         [Write the spec]
	//     [Done]
	//         [Read the spec]
}

// ExampleDiagram_Task adds a card to the current column.
func ExampleDiagram_Task() {
	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo").
		Task("Write the spec").
		Task("Review the spec").
		Build()

	// Output:
	// kanban
	//     [Todo]
	//         [Write the spec]
	//         [Review the spec]
}

// ExampleDiagram_TaskIn adds a card to a column named outright, which saves
// switching back and forth when the cards of two columns are interleaved in the
// calling code.
func ExampleDiagram_TaskIn() {
	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo").
		Column("Done").
		TaskIn("Todo", "Write the spec").
		TaskIn("Done", "Read the spec").
		Build()

	// Output:
	// kanban
	//     [Todo]
	//     [Done]
	//     [Todo]
	//         [Write the spec]
	//     [Done]
	//         [Read the spec]
}

// ExampleWithColumnID gives a column an identifier of its own, which is what
// TaskIn and a link from outside the board refer to.
func ExampleWithColumnID() {
	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo", kanban.WithColumnID("todo")).
		Task("Write the spec").
		Build()

	// Output:
	// kanban
	//     todo[Todo]
	//         [Write the spec]
}

// ExampleWithTaskID gives a card an identifier of its own.
func ExampleWithTaskID() {
	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo").
		Task("Write the spec", kanban.WithTaskID("kb-1")).
		Build()

	// Output:
	// kanban
	//     [Todo]
	//         kb-1[Write the spec]
}

// ExampleWithTaskTicket puts a ticket reference on a card. With
// WithTicketBaseURL on the board, mermaid draws it as a link.
func ExampleWithTaskTicket() {
	_ = kanban.NewDiagram(os.Stdout, kanban.WithTicketBaseURL("https://example.com/browse/")).
		Column("Todo").
		Task("Write the spec", kanban.WithTaskTicket("KB-1")).
		Build()

	// Output:
	// ---
	// config:
	//   kanban:
	//     ticketBaseUrl: 'https://example.com/browse/'
	// ---
	// kanban
	//     [Todo]
	//         [Write the spec]@{ ticket: 'KB-1' }
}

// ExampleWithTaskAssigned puts a name on a card.
func ExampleWithTaskAssigned() {
	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo").
		Task("Write the spec", kanban.WithTaskAssigned("Alice")).
		Build()

	// Output:
	// kanban
	//     [Todo]
	//         [Write the spec]@{ assigned: 'Alice' }
}

// ExampleWithTaskPriority puts a priority on a card, which mermaid draws as a
// colored stripe down its side.
func ExampleWithTaskPriority() {
	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo").
		Task("Write the spec", kanban.WithTaskPriority(kanban.PriorityVeryHigh)).
		Build()

	// Output:
	// kanban
	//     [Todo]
	//         [Write the spec]@{ priority: 'Very High' }
}

// ExamplePriority shows the four priorities a card can carry.
func ExamplePriority() {
	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo").
		Task("Lowest", kanban.WithTaskPriority(kanban.PriorityVeryLow)).
		Task("Low", kanban.WithTaskPriority(kanban.PriorityLow)).
		Task("High", kanban.WithTaskPriority(kanban.PriorityHigh)).
		Task("Highest", kanban.WithTaskPriority(kanban.PriorityVeryHigh)).
		Build()

	// Output:
	// kanban
	//     [Todo]
	//         [Lowest]@{ priority: 'Very Low' }
	//         [Low]@{ priority: 'Low' }
	//         [High]@{ priority: 'High' }
	//         [Highest]@{ priority: 'Very High' }
}

// ExampleWithTicketBaseURL sets what a card's ticket reference is prefixed with,
// which is what turns "KB-1" into a link a reader can follow.
func ExampleWithTicketBaseURL() {
	_ = kanban.NewDiagram(os.Stdout, kanban.WithTicketBaseURL("https://example.com/browse/")).
		Column("Todo").
		Task("Write the spec", kanban.WithTaskTicket("KB-1")).
		Build()

	// Output:
	// ---
	// config:
	//   kanban:
	//     ticketBaseUrl: 'https://example.com/browse/'
	// ---
	// kanban
	//     [Todo]
	//         [Write the spec]@{ ticket: 'KB-1' }
}

// ExampleColumnOption shows what a ColumnOption is: a function that changes how
// a column is written, passed to Column after its name.
func ExampleColumnOption() {
	options := []kanban.ColumnOption{kanban.WithColumnID("todo")}

	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo", options...).
		Task("Write the spec").
		Build()

	// Output:
	// kanban
	//     todo[Todo]
	//         [Write the spec]
}

// ExampleTaskOption shows what a TaskOption is: a function that changes how a
// card is written, passed to Task after its name.
func ExampleTaskOption() {
	options := []kanban.TaskOption{kanban.WithTaskID("kb-1"), kanban.WithTaskAssigned("Alice")}

	_ = kanban.NewDiagram(os.Stdout).
		Column("Todo").
		Task("Write the spec", options...).
		Build()

	// Output:
	// kanban
	//     [Todo]
	//         kb-1[Write the spec]@{ assigned: 'Alice' }
}
