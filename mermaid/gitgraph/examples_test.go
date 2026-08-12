//go:build linux || darwin

package gitgraph_test

import (
	"fmt"
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
	// title: "Release Flow"
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

// ExampleNewDiagram shows the shape every git graph has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := gitgraph.NewDiagram(io.Discard).
		Commit(gitgraph.WithCommitID("initial")).
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// gitGraph
	//     commit id: "initial"
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := gitgraph.NewDiagram(nil).
		Commit(gitgraph.WithCommitID("initial")).
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := gitgraph.NewDiagram(io.Discard).
		Checkout("never-branched")
	fmt.Println("error:", d.Error())

	// Output:
	// error: <nil>
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		LF().
		Commit(gitgraph.WithCommitID("initial")).
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
	//
	//     commit id: "initial"
}

// ExampleDiagram_full shows a git graph built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := gitgraph.NewDiagram(io.Discard).
		Commit(gitgraph.WithCommitID("initial")).
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// gitGraph
	//     commit id: "initial"
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []gitgraph.Option{gitgraph.WithTitle("Overview")}

	_ = gitgraph.NewDiagram(os.Stdout, options...).
		Commit(gitgraph.WithCommitID("initial")).
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// gitGraph
	//     commit id: "initial"
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = gitgraph.NewDiagram(os.Stdout, gitgraph.WithTitle("Overview")).
		Commit(gitgraph.WithCommitID("initial")).
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// gitGraph
	//     commit id: "initial"
}

// ExampleDiagram_Commit adds a commit to the current branch.
func ExampleDiagram_Commit() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit().
		Commit(gitgraph.WithCommitID("second")).
		Build()

	// Output:
	// gitGraph
	//     commit
	//     commit id: "second"
}

// ExampleDiagram_Branch starts a branch and switches to it, so the commits that
// follow belong to it.
func ExampleDiagram_Branch() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		Branch("develop").
		Commit(gitgraph.WithCommitID("work")).
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
	//     branch develop
	//     commit id: "work"
}

// ExampleDiagram_Checkout switches back to a branch that already exists.
func ExampleDiagram_Checkout() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		Branch("develop").
		Commit(gitgraph.WithCommitID("work")).
		Checkout("main").
		Commit(gitgraph.WithCommitID("hotfix")).
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
	//     branch develop
	//     commit id: "work"
	//     checkout main
	//     commit id: "hotfix"
}

// ExampleDiagram_Merge brings a branch back into the current one. It takes the
// commit options, because a merge is drawn as a commit and can carry a tag.
func ExampleDiagram_Merge() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		Branch("develop").
		Commit(gitgraph.WithCommitID("work")).
		Checkout("main").
		Merge("develop", gitgraph.WithCommitTag("v1.0.0")).
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
	//     branch develop
	//     commit id: "work"
	//     checkout main
	//     merge develop tag: "v1.0.0"
}

// ExampleDiagram_CherryPick copies one commit onto the current branch. The
// commit it names has to exist already.
func ExampleDiagram_CherryPick() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		Branch("develop").
		Commit(gitgraph.WithCommitID("fix")).
		Checkout("main").
		CherryPick("fix").
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
	//     branch develop
	//     commit id: "fix"
	//     checkout main
	//     cherry-pick id: "fix"
}

// ExampleDiagram_Reset moves the current branch back to a commit.
func ExampleDiagram_Reset() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		Commit(gitgraph.WithCommitID("mistake")).
		Reset("initial").
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
	//     commit id: "mistake"
	//     reset id: "initial"
}

// ExampleWithCommitID names a commit, which is what a cherry pick and a reset
// refer to.
func ExampleWithCommitID() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
}

// ExampleWithCommitTag puts a tag on a commit, drawn beside it.
func ExampleWithCommitTag() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("release"), gitgraph.WithCommitTag("v1.0.0")).
		Build()

	// Output:
	// gitGraph
	//     commit id: "release" tag: "v1.0.0"
}

// ExampleWithCommitType changes the shape a commit is drawn as.
func ExampleWithCommitType() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitType(gitgraph.CommitTypeHighlight)).
		Build()

	// Output:
	// gitGraph
	//     commit type: HIGHLIGHT
}

// ExampleCommitType shows the three shapes a commit can be drawn as.
func ExampleCommitType() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("a"), gitgraph.WithCommitType(gitgraph.CommitTypeNormal)).
		Commit(gitgraph.WithCommitID("b"), gitgraph.WithCommitType(gitgraph.CommitTypeReverse)).
		Commit(gitgraph.WithCommitID("c"), gitgraph.WithCommitType(gitgraph.CommitTypeHighlight)).
		Build()

	// Output:
	// gitGraph
	//     commit id: "a" type: NORMAL
	//     commit id: "b" type: REVERSE
	//     commit id: "c" type: HIGHLIGHT
}

// ExampleWithBranchOrder says where a branch is drawn against the others, which
// is how a diagram keeps the important branch at the top.
func ExampleWithBranchOrder() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		Branch("release", gitgraph.WithBranchOrder(1)).
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
	//     branch release order: 1
}

// ExampleWithCherryPickParent says which parent of a merge commit a cherry pick
// takes, which mermaid needs when the commit being picked is a merge.
func ExampleWithCherryPickParent() {
	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		Branch("develop").
		Commit(gitgraph.WithCommitID("work")).
		Checkout("main").
		Merge("develop", gitgraph.WithCommitID("merge")).
		Branch("release").
		CherryPick("merge", gitgraph.WithCherryPickParent("initial")).
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
	//     branch develop
	//     commit id: "work"
	//     checkout main
	//     merge develop id: "merge"
	//     branch release
	//     cherry-pick id: "merge" parent: "initial"
}

// ExampleCommitOption shows what a CommitOption is: a function that changes how
// a commit is written, passed to Commit or to Merge.
func ExampleCommitOption() {
	options := []gitgraph.CommitOption{
		gitgraph.WithCommitID("release"),
		gitgraph.WithCommitTag("v1.0.0"),
	}

	_ = gitgraph.NewDiagram(os.Stdout).Commit(options...).Build()

	// Output:
	// gitGraph
	//     commit id: "release" tag: "v1.0.0"
}

// ExampleBranchOption shows what a BranchOption is: a function that changes how
// a branch is written, passed to Branch after its name.
func ExampleBranchOption() {
	options := []gitgraph.BranchOption{gitgraph.WithBranchOrder(2)}

	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		Branch("release", options...).
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
	//     branch release order: 2
}

// ExampleCherryPickOption shows what a CherryPickOption is: a function that
// changes how a cherry pick is written, passed to CherryPick after the commit
// it names.
func ExampleCherryPickOption() {
	options := []gitgraph.CherryPickOption{}

	_ = gitgraph.NewDiagram(os.Stdout).
		Commit(gitgraph.WithCommitID("initial")).
		Branch("develop").
		Commit(gitgraph.WithCommitID("fix")).
		Checkout("main").
		CherryPick("fix", options...).
		Build()

	// Output:
	// gitGraph
	//     commit id: "initial"
	//     branch develop
	//     commit id: "fix"
	//     checkout main
	//     cherry-pick id: "fix"
}
