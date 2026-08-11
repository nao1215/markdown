package markdown_test

import (
	"io"

	"github.com/nao1215/markdown/mermaid/arch"
	"github.com/nao1215/markdown/mermaid/block"
	"github.com/nao1215/markdown/mermaid/class"
	"github.com/nao1215/markdown/mermaid/er"
	"github.com/nao1215/markdown/mermaid/flowchart"
	"github.com/nao1215/markdown/mermaid/gantt"
	"github.com/nao1215/markdown/mermaid/gitgraph"
	"github.com/nao1215/markdown/mermaid/kanban"
	"github.com/nao1215/markdown/mermaid/mindmap"
	"github.com/nao1215/markdown/mermaid/packet"
	"github.com/nao1215/markdown/mermaid/piechart"
	"github.com/nao1215/markdown/mermaid/quadrant"
	"github.com/nao1215/markdown/mermaid/requirement"
	"github.com/nao1215/markdown/mermaid/sequence"
	"github.com/nao1215/markdown/mermaid/state"
	"github.com/nao1215/markdown/mermaid/userjourney"
	"github.com/nao1215/markdown/mermaid/xychart"
)

// mermaidBuilder names one mermaid subpackage and builds a small diagram with
// it.
//
// The list lives in the root package's tests because what it is used for is a
// property of the pair: a diagram is only ever seen by a reader after it has
// been handed to CodeBlocks, and it is the combination that has to hold. Seven
// lines per subpackage here beats a near-identical test file in each of the
// seventeen.
type mermaidBuilder struct {
	// name is the subpackage the diagram comes from.
	name string
	// build returns the diagram body, as CodeBlocks would receive it.
	build func() string
}

// mermaidBuilders holds one builder per mermaid subpackage. A new subpackage
// belongs here as soon as it exists.
func mermaidBuilders() []mermaidBuilder {
	return []mermaidBuilder{
		{
			name: "arch",
			build: func() string {
				return arch.NewArchitecture(io.Discard).
					Service("api", arch.IconServer, "API").
					String()
			},
		},
		{
			name: "block",
			build: func() string {
				return block.NewDiagram(io.Discard).
					Row(block.Node("a", block.WithNodeLabel("A"))).
					String()
			},
		},
		{
			name: "class",
			build: func() string {
				return class.NewDiagram(io.Discard).
					Class("Account", class.WithPublicField("string", "id")).
					String()
			},
		},
		{
			name: "er",
			build: func() string {
				return er.NewDiagram(io.Discard).
					NoRelationship(er.NewEntity("teachers", []*er.Attribute{
						{Type: "int", Name: "id", IsPrimaryKey: true, Comment: "Teacher ID"},
					})).
					String()
			},
		},
		{
			name: "flowchart",
			build: func() string {
				return flowchart.NewFlowchart(io.Discard).
					NodeWithText("A", "Node A").
					String()
			},
		},
		{
			name: "gantt",
			build: func() string {
				return gantt.NewChart(io.Discard).
					Section("Planning").
					Task("Design", "2024-01-01", "2d").
					String()
			},
		},
		{
			name: "gitgraph",
			build: func() string {
				return gitgraph.NewDiagram(io.Discard).
					Commit(gitgraph.WithCommitTag("v1.0.0")).
					String()
			},
		},
		{
			name: "kanban",
			build: func() string {
				return kanban.NewDiagram(io.Discard).
					Column("Todo").
					Task("Define scope").
					String()
			},
		},
		{
			name: "mindmap",
			build: func() string {
				return mindmap.NewDiagram(io.Discard).
					Root("Product").
					Child("Market").
					String()
			},
		},
		{
			name: "packet",
			build: func() string {
				return packet.NewDiagram(io.Discard).
					Field(0, 15, "Source Port").
					String()
			},
		},
		{
			name: "piechart",
			build: func() string {
				return piechart.NewPieChart(io.Discard).
					LabelAndIntValue("Go", 120).
					String()
			},
		},
		{
			name: "quadrant",
			build: func() string {
				return quadrant.NewChart(io.Discard).
					XAxis("Low Reach", "High Reach").
					Point("Campaign A", 0.3, 0.6).
					String()
			},
		},
		{
			name: "requirement",
			build: func() string {
				return requirement.NewDiagram(io.Discard).
					Requirement(
						"a requirement",
						requirement.WithID("1"),
						requirement.WithText("the system shall do the thing"),
						requirement.WithRisk(requirement.RiskLow),
						requirement.WithVerifyMethod(requirement.VerifyMethodTest),
					).
					String()
			},
		},
		{
			name: "sequence",
			build: func() string {
				return sequence.NewDiagram(io.Discard).
					SyncRequest("Client", "Server", "GET /users").
					String()
			},
		},
		{
			name: "state",
			build: func() string {
				return state.NewDiagram(io.Discard).
					State("Draft", "The order is being written").
					Transition("Draft", "Placed").
					String()
			},
		},
		{
			name: "userjourney",
			build: func() string {
				return userjourney.NewDiagram(io.Discard).
					Section("Discovery").
					Task("Find the site", userjourney.ScoreSatisfied, "Visitor").
					String()
			},
		},
		{
			name: "xychart",
			build: func() string {
				return xychart.NewDiagram(io.Discard).
					XAxisLabels("jan", "feb").
					Bar(5000, 6000).
					String()
			},
		},
	}
}
