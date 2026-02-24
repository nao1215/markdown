//go:build linux || darwin

package kanban_test

import (
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
	// title: Sprint Board
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
