//go:build linux || darwin

// Package main is generating mermaid Gantt chart.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/gantt"
)

//go:generate go run main.go

func main() {
	f, err := os.Create("generated.md")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	chart := gantt.NewChart(
		io.Discard,
		gantt.WithTitle("Software Development Schedule"),
		gantt.WithDateFormat("YYYY-MM-DD"),
	).
		Section("Planning").
		DoneTaskWithID("Requirements Analysis", "req", "2024-01-01", "7d").
		DoneTaskWithID("System Design", "design", "2024-01-08", "5d").
		LF().
		Section("Development").
		CriticalActiveTaskWithID("Backend Development", "backend", "2024-01-15", "14d"). //nolint:mnd
		ActiveTaskWithID("Frontend Development", "frontend", "2024-01-15", "14d").       //nolint:mnd
		TaskAfterWithID("Integration", "integrate", "backend", "5d").                    //nolint:mnd
		LF().
		Section("Testing").
		TaskAfterWithID("Unit Testing", "unit", "integrate", "3d").      //nolint:mnd
		TaskAfterWithID("Integration Testing", "inttest", "unit", "4d"). //nolint:mnd
		TaskAfterWithID("UAT", "uat", "inttest", "5d").                  //nolint:mnd
		LF().
		Section("Deployment").
		TaskAfter("Staging Deploy", "uat", "2d").              //nolint:mnd
		CriticalMilestone("Production Release", "2024-03-01"). //nolint:mnd
		String()

	err = markdown.NewMarkdown(f).
		H2("Gantt Chart").
		CodeBlocks(markdown.SyntaxHighlightMermaid, chart).
		Build()

	if err != nil {
		panic(err)
	}
}
