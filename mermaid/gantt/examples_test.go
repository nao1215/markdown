//go:build linux || darwin

package gantt_test

import (
	"fmt"
	"io"
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

// ExampleChart_Task adds a task that has not started.
func ExampleChart_Task() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		Task("Design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Design :2024-01-01, 10d
}

// ExampleChart_ActiveTask adds a task in progress, which mermaid draws hatched.
func ExampleChart_ActiveTask() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		ActiveTask("Build", "2024-01-11", "20d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Build :active, 2024-01-11, 20d
}

// ExampleChart_DoneTask adds a task that is finished, which mermaid draws filled in.
func ExampleChart_DoneTask() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		DoneTask("Design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Design :done, 2024-01-01, 10d
}

// ExampleChart_CriticalTask adds a task on the critical path, which mermaid draws in red.
func ExampleChart_CriticalTask() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		CriticalTask("Sign the contract", "2024-01-01", "5d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Sign the contract :crit, 2024-01-01, 5d
}

// ExampleChart_CriticalActiveTask adds a critical task in progress.
func ExampleChart_CriticalActiveTask() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		CriticalActiveTask("Migrate the data", "2024-01-06", "5d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Migrate the data :crit, active, 2024-01-06, 5d
}

// ExampleChart_CriticalDoneTask adds a critical task that is finished.
func ExampleChart_CriticalDoneTask() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		CriticalDoneTask("Sign the contract", "2024-01-01", "5d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Sign the contract :crit, done, 2024-01-01, 5d
}

// ExampleChart_TaskWithID is Task with an identifier, which is what a later
// task refers to when it says it starts after this one.
func ExampleChart_TaskWithID() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		TaskWithID("Design", "design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Design :design, 2024-01-01, 10d
}

// ExampleChart_ActiveTaskWithID is ActiveTask with an identifier, which is what a later
// task refers to when it says it starts after this one.
func ExampleChart_ActiveTaskWithID() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		ActiveTaskWithID("Build", "build", "2024-01-11", "20d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Build :active, build, 2024-01-11, 20d
}

// ExampleChart_DoneTaskWithID is DoneTask with an identifier, which is what a later
// task refers to when it says it starts after this one.
func ExampleChart_DoneTaskWithID() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		DoneTaskWithID("Design", "design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Design :done, design, 2024-01-01, 10d
}

// ExampleChart_CriticalTaskWithID is CriticalTask with an identifier, which is what a later
// task refers to when it says it starts after this one.
func ExampleChart_CriticalTaskWithID() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		CriticalTaskWithID("Sign the contract", "sign", "2024-01-01", "5d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Sign the contract :crit, sign, 2024-01-01, 5d
}

// ExampleChart_CriticalActiveTaskWithID is CriticalActiveTask with an identifier, which is what a later
// task refers to when it says it starts after this one.
func ExampleChart_CriticalActiveTaskWithID() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		CriticalActiveTaskWithID("Migrate the data", "migrate", "2024-01-06", "5d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Migrate the data :crit, active, migrate, 2024-01-06, 5d
}

// ExampleChart_CriticalDoneTaskWithID is CriticalDoneTask with an identifier, which is what a later
// task refers to when it says it starts after this one.
func ExampleChart_CriticalDoneTaskWithID() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		CriticalDoneTaskWithID("Sign the contract", "sign", "2024-01-01", "5d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Sign the contract :crit, done, sign, 2024-01-01, 5d
}

// ExampleChart_MilestoneWithID is Milestone with an identifier, which is what a later
// task refers to when it says it starts after this one.
func ExampleChart_MilestoneWithID() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		MilestoneWithID("Launch", "launch", "2024-02-01").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Launch :milestone, launch, 2024-02-01, 0d
}

// ExampleChart_CriticalMilestoneWithID is CriticalMilestone with an identifier, which is what a later
// task refers to when it says it starts after this one.
func ExampleChart_CriticalMilestoneWithID() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		CriticalMilestoneWithID("Launch", "launch", "2024-02-01").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Launch :crit, milestone, launch, 2024-02-01, 0d
}

// ExampleNewChart shows the shape every gantt chart has: a writer, a chain of
// calls, and Build.
func ExampleNewChart() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		Task("Design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Design :2024-01-01, 10d
}

// ExampleChart_Section groups the tasks that follow it into a band of the
// chart.
func ExampleChart_Section() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Design").
		Task("Wireframes", "2024-01-01", "5d").
		Section("Build").
		Task("Implement", "2024-01-06", "20d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Design
	//     Wireframes :2024-01-01, 5d
	//     section Build
	//     Implement :2024-01-06, 20d
}

// ExampleChart_Milestone marks a date rather than a span, so it is drawn as a
// diamond with no width.
func ExampleChart_Milestone() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		Milestone("Launch", "2024-02-01").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Launch :milestone, 2024-02-01, 0d
}

// ExampleChart_CriticalMilestone marks a date on the critical path.
func ExampleChart_CriticalMilestone() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		CriticalMilestone("Contract expires", "2024-02-01").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Contract expires :crit, milestone, 2024-02-01, 0d
}

// ExampleChart_TaskAfter starts a task when another finishes rather than on a
// date, so moving the first one moves everything behind it.
func ExampleChart_TaskAfter() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		TaskWithID("Design", "design", "2024-01-01", "10d").
		TaskAfter("Build", "design", "20d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Design :design, 2024-01-01, 10d
	//     Build :after design, 20d
}

// ExampleChart_TaskAfterWithID starts a task after another and gives it an
// identifier, so a third can be hung off it in turn.
func ExampleChart_TaskAfterWithID() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		TaskWithID("Design", "design", "2024-01-01", "10d").
		TaskAfterWithID("Build", "build", "design", "20d").
		TaskAfter("Release", "build", "2d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Design :design, 2024-01-01, 10d
	//     Build :build, after design, 20d
	//     Release :after build, 2d
}

// ExampleChart_String returns the chart without needing a writer, which is how
// it is handed to a markdown code block.
func ExampleChart_String() {
	chart := gantt.NewChart(io.Discard).
		Section("Delivery").
		Task("Design", "2024-01-01", "10d").
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, chart).
		Build()

	// Output:
	// ```mermaid
	// gantt
	//     section Delivery
	//     Design :2024-01-01, 10d
	// ```
}

// ExampleChart_Build writes the chart and reports the error the chain recorded.
func ExampleChart_Build() {
	err := gantt.NewChart(nil).Section("Delivery").Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleChart_Error reports the same error Build does, for code that wants to
// look before writing anything.
func ExampleChart_Error() {
	c := gantt.NewChart(io.Discard).Section("Delivery")
	fmt.Println("error:", c.Error())

	// Output:
	// error: <nil>
}

// ExampleChart_LF adds a blank line to the chart body.
func ExampleChart_LF() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Design").
		Task("Wireframes", "2024-01-01", "5d").
		LF().
		Section("Build").
		Task("Implement", "2024-01-06", "20d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Design
	//     Wireframes :2024-01-01, 5d
	//
	//     section Build
	//     Implement :2024-01-06, 20d
}

// ExampleWithTitle sets the title the chart is drawn with.
func ExampleWithTitle() {
	_ = gantt.NewChart(os.Stdout, gantt.WithTitle("Release Plan")).
		Section("Delivery").
		Task("Design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     title Release Plan
	//     section Delivery
	//     Design :2024-01-01, 10d
}

// ExampleWithDateFormat says how the dates in the chart are written. Without it
// mermaid falls back to its own default, which is rarely the one the dates are
// actually in.
func ExampleWithDateFormat() {
	_ = gantt.NewChart(os.Stdout, gantt.WithDateFormat("YYYY-MM-DD")).
		Section("Delivery").
		Task("Design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Design :2024-01-01, 10d
}

// ExampleWithAxisFormat says how the dates along the bottom are drawn, which is
// separate from how the dates in the chart are written.
func ExampleWithAxisFormat() {
	_ = gantt.NewChart(os.Stdout,
		gantt.WithDateFormat("YYYY-MM-DD"),
		gantt.WithAxisFormat("%d %b"),
	).
		Section("Delivery").
		Task("Design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     axisFormat %d %b
	//     section Delivery
	//     Design :2024-01-01, 10d
}

// ExampleWithTickInterval sets how often a mark appears along the bottom.
func ExampleWithTickInterval() {
	_ = gantt.NewChart(os.Stdout,
		gantt.WithDateFormat("YYYY-MM-DD"),
		gantt.WithTickInterval("1week"),
	).
		Section("Delivery").
		Task("Design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     tickInterval 1week
	//     section Delivery
	//     Design :2024-01-01, 10d
}

// ExampleWithTodayMarker styles the line marking today, or hides it outright
// with "off", which is what a chart of a finished project wants.
func ExampleWithTodayMarker() {
	_ = gantt.NewChart(os.Stdout,
		gantt.WithDateFormat("YYYY-MM-DD"),
		gantt.WithTodayMarker("off"),
	).
		Section("Delivery").
		Task("Design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     todayMarker off
	//     section Delivery
	//     Design :2024-01-01, 10d
}

// ExampleWithExcludes leaves days out of the durations, so a ten day task that
// skips weekends ends two weeks later rather than ten days later.
func ExampleWithExcludes() {
	_ = gantt.NewChart(os.Stdout,
		gantt.WithDateFormat("YYYY-MM-DD"),
		gantt.WithExcludes("weekends"),
	).
		Section("Delivery").
		Task("Design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     dateFormat YYYY-MM-DD
	//     excludes weekends
	//     section Delivery
	//     Design :2024-01-01, 10d
}

// ExampleOption shows what an Option is: a function that changes how the chart
// is written, passed to NewChart.
func ExampleOption() {
	options := []gantt.Option{
		gantt.WithTitle("Release Plan"),
		gantt.WithDateFormat("YYYY-MM-DD"),
	}

	_ = gantt.NewChart(os.Stdout, options...).
		Section("Delivery").
		Task("Design", "2024-01-01", "10d").
		Build()

	// Output:
	// gantt
	//     title Release Plan
	//     dateFormat YYYY-MM-DD
	//     section Delivery
	//     Design :2024-01-01, 10d
}
