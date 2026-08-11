package gitgraph_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/gitgraph"
)

// TestGoldenGitGraph pins the rendered diagram of every builder method and
// every option of this package.
func TestGoldenGitGraph(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := gitgraph.NewDiagram(buf, gitgraph.WithTitle("Release Flow")).
		Commit().
		Commit(gitgraph.WithCommitID("init"), gitgraph.WithCommitTag("v0.1.0")).
		Commit(gitgraph.WithCommitType(gitgraph.CommitTypeNormal)).
		Commit(gitgraph.WithCommitID("reverse"), gitgraph.WithCommitType(gitgraph.CommitTypeReverse)).
		Commit(gitgraph.WithCommitID("highlight"), gitgraph.WithCommitType(gitgraph.CommitTypeHighlight)).
		Branch("develop").
		Branch("feature", gitgraph.WithBranchOrder(2)).
		Checkout("develop").
		Commit(gitgraph.WithCommitID("dev-1")).
		LF().
		Checkout("main").
		Merge("develop", gitgraph.WithCommitTag("v1.0.0")).
		CherryPick("dev-1").
		CherryPick("dev-1", gitgraph.WithCherryPickParent("init")).
		Reset("init").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("gitgraph.md", buf.String()); err != nil {
		t.Error(err)
	}
}
