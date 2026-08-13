//go:build linux || darwin

// Package main is generating a GitHub Actions job summary.
//
// Inside a GitHub Actions step, GITHUB_STEP_SUMMARY names a file whose content
// the run's summary page renders as GitHub Flavored Markdown, mermaid diagrams
// included. This program appends to that file when the variable is set, so it
// runs as a step unchanged; outside one it writes the committed sample.
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/piechart"
)

//go:generate go run main.go

func main() {
	path := os.Getenv("GITHUB_STEP_SUMMARY")
	if path == "" {
		path = "generated.md"
	}
	// The variable path is the point of this program: the Actions runner
	// names the file the summary page renders, and this appends to it.
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600) //nolint:gosec
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	coverage := piechart.NewPieChart(
		io.Discard,
		piechart.WithTitle("Coverage"),
		piechart.WithShowData(true),
	).
		LabelAndIntValue("covered", 92).  //nolint:mnd
		LabelAndIntValue("uncovered", 8). //nolint:mnd
		String()

	err = markdown.NewMarkdown(f, markdown.WithBlockSpacing()).
		H2("Test Results").
		Table(markdown.TableSet{
			Header: []string{"Package", "Passed", "Failed"},
			Rows: [][]string{
				{"api", "120", "0"},
				{"core", "89", "2"},
			},
		}).
		Warning("2 tests failed in core; see the failed step for logs.").
		CodeBlocks(markdown.SyntaxHighlightMermaid, coverage).
		Build()

	if err != nil {
		panic(err)
	}
}
