//go:build linux || darwin

package gitgraph_test

import (
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/gitgraph"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := gitgraph.NewDiagram(
		io.Discard,
		gitgraph.WithTitle("Release Flow"),
	).
		Commit(gitgraph.WithCommitID("init"), gitgraph.WithCommitTag("v0.1.0")).
		Branch("develop").
		Checkout("develop").
		Commit(gitgraph.WithCommitType(gitgraph.CommitTypeHighlight)).
		Checkout("main").
		Merge("develop", gitgraph.WithCommitTag("v1.0.0")).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Git Graph").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Git Graph
	// ```mermaid
	// ---
	// title: Release Flow
	// ---
	// gitGraph
	//     commit id: "init" tag: "v0.1.0"
	//     branch develop
	//     checkout develop
	//     commit type: HIGHLIGHT
	//     checkout main
	//     merge develop tag: "v1.0.0"
	// ```
}
