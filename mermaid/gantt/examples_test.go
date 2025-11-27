//go:build linux || darwin

package gantt_test

import (
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/gantt"
)

// ExampleChart skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleChart() {
	chart := gantt.NewChart(
		os.Stdout,
		gantt.WithTitle("Project Schedule"),
		gantt.WithDateFormat("YYYY-MM-DD"),
	).
		Section("Planning").
		DoneTaskWithID("Requirements", "req", "2024-01-01", "5d").
		DoneTaskWithID("Design", "design", "2024-01-08", "3d").
		Section("Development").
		CriticalActiveTaskWithID("Coding", "code", "2024-01-12", "10d").
		TaskAfterWithID("Review", "review", "code", "2d").
		Section("Release").
		MilestoneWithID("Launch", "launch", "2024-01-26").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Gantt Chart").
		CodeBlocks(md.SyntaxHighlightMermaid, chart).
		Build()

	// Output:
	// ## Gantt Chart
	// ```mermaid
	// gantt
	//     title Project Schedule
	//     dateFormat YYYY-MM-DD
	//     section Planning
	//     Requirements :done, req, 2024-01-01, 5d
	//     Design :done, design, 2024-01-08, 3d
	//     section Development
	//     Coding :crit, active, code, 2024-01-12, 10d
	//     Review :review, after code, 2d
	//     section Release
	//     Launch :milestone, launch, 2024-01-26, 0d
	// ```
}
