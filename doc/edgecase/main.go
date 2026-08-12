//go:build linux || darwin

// Package main generates one document per mermaid subpackage, every label of
// which holds the punctuation that means something to mermaid.
//
// The other generators under doc/ are documentation: they are quoted in the
// README and they have to stay readable. That makes them poor tests, because
// the labels a reader wants to see are exactly the ones that never break a
// renderer. The labels here are the opposite, and the documents they produce
// are rendered in continuous integration alongside the others, so a quoting
// defect fails the build rather than reaching a user's diagram.
//
// One document per diagram type, rather than one document holding all of them,
// so that a failing render names the subpackage that produced it instead of a
// line number that moves whenever a label changes.
package main

import (
	"io"
	"os"
	"strings"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/arch"
	"github.com/nao1215/markdown/mermaid/block"
	"github.com/nao1215/markdown/mermaid/c4"
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
	"github.com/nao1215/markdown/mermaid/radar"
	"github.com/nao1215/markdown/mermaid/requirement"
	"github.com/nao1215/markdown/mermaid/sankey"
	"github.com/nao1215/markdown/mermaid/sequence"
	"github.com/nao1215/markdown/mermaid/state"
	"github.com/nao1215/markdown/mermaid/treemap"
	"github.com/nao1215/markdown/mermaid/userjourney"
	"github.com/nao1215/markdown/mermaid/xychart"
)

// This file is gated by //go:build linux || darwin, so //go:generate is skipped
// on Windows. To regenerate the documents on Windows, run under WSL or via CI.
//go:generate go run main.go

const (
	emoji    = "🎉"
	japanese = "日本語"
)

// supported returns the punctuation the given diagram type can carry today.
//
// Every entry was measured rather than guessed: each character was put through
// its diagram type on its own and rendered, and the ones that survived are
// here. What is missing from an entry is either a quoting gap in this library
// or a limit of the mermaid syntax, and each is recorded with its evidence in
// the tracking issue for mermaid label quoting.
//
// A gap closed by a fix belongs in its entry the moment the fix lands, which is
// what keeps this file honest: the labels only ever get harder.
func supported(diagram string) string {
	return map[string]string{
		// architecture writes service and group labels between square brackets
		// with no quoting at all, and mermaid's architecture-beta grammar takes
		// only word characters there. Every character probed fails, including a
		// plain emoji and Japanese text.
		"architecture": "",
		"block":        `"'#;[](){}<br/>` + emoji + japanese + `:,*-|%%`,
		// c4 escapes a quote and a "#" in a label as the entity form mermaid
		// decodes, and a "#" and a ";" in the unquoted title the same way, so
		// every character below survives both. Only a line break is left out:
		// one macro is one line, and "<br/>" is how a label breaks one.
		"c4": `"'#;[](){}<br/>` + emoji + japanese + `:,*-|%%\`,
		// class: ";" and ":" end a relation label.
		"class": `"'#[](){}` + emoji + japanese + `,*-|%%`,
		// er writes a comment as the entity form mermaid decodes, so the
		// punctuation is all safe.
		"er": `"'#;[](){}<br/>` + emoji + japanese + `:,*-|%%`,
		// flowchart quotes every label and writes a double quote as the entity
		// form mermaid decodes, so the punctuation is all safe. What is left
		// out is a line break: the renderer honours "<br/>" inside a label,
		// which is what a caller wanting one should pass.
		"flowchart": `"'#;[](){}` + emoji + japanese + `:,*-|%%\`,
		// gantt: a colon separates a task name from its data.
		"gantt":    `"'#;[](){}<br/>` + emoji + japanese + `,*-|%%`,
		"gitgraph": `"'#;[](){}` + emoji + japanese + `:,*-|%%`,
		// kanban: quotes end the single quoted metadata, and square brackets
		// and a closing brace end a node.
		"kanban": `#;{<br/>` + emoji + japanese + `:,*-|%%`,
		// mindmap: brackets, parentheses and braces are node shape delimiters.
		"mindmap": `"'#;]<br/>` + emoji + japanese + `:,*-|%%`,
		// packet and piechart put the title in YAML front matter, and mermaid
		// strips a "%%" comment from that before the YAML is read.
		"packet": `"'#;[](){}<br/>` + emoji + japanese + `:,*-|`,
		// piechart quotes a slice label and writes a double quote as the entity
		// form mermaid decodes. Its title is unquoted, so a "%%" there would
		// open a comment and cut the title short; that pair is written as an
		// entity too.
		"piechart": `"'#;[](){}<br/>` + emoji + japanese + `:,*-|%%`,
		// radar quotes every label, so the punctuation is all safe. A line
		// break is the only thing left out, since one label is one line.
		"radar": `"'#;[](){}` + emoji + japanese + `:,*-|%%`,
		// quadrant: axis, quadrant and point labels are all written unquoted.
		"quadrant":    `'#` + emoji + japanese + `,*-`,
		"requirement": `"'#;[](){}` + emoji + japanese + `:,*-|%%`,
		// sankey quotes a node name that needs it, so the punctuation is all
		// safe. What it cannot take is non-ASCII text: mermaid's sankey parser
		// refuses an emoji or Japanese in a node name, the same way its xy
		// chart parser does.
		"sankey": `"'#;[](){}<br/>:,*-|%%`,
		// sequence: ";" and "%%" end a statement, ":" and "," end a participant
		// name, and a parenthesis or a brace in one ends it too once the name
		// holds anything else.
		"sequence": `"'#[]` + emoji + japanese + `*-|`,
		// state: a colon ends the note it is written in. A quote, a hash, a
		// semicolon, a bracket, a brace or a parenthesis each end a state
		// description once it holds anything else, so a description is down to
		// text and the punctuation below.
		"state": emoji + japanese + `,*-|%%`,
		// treemap doubles a quote in a name and everything else, a backslash
		// included, is text inside the quotes, so only a line break is left
		// out: the hierarchy is indentation, and a name spanning lines would
		// read as another node. It is the only type here whose labels carry a
		// backslash, because it is the only one whose quoting was proven to
		// leave one alone.
		"treemap": `"'#;[](){}<br/>` + emoji + japanese + `:,*-|%%\`,
		// userjourney: "#" and ";" end a statement, and ":" separates the
		// fields of a task.
		"userjourney": `"'[](){}<br/>` + emoji + japanese + `,*-|%%`,
		// xychart: an axis label is written unquoted, so non-ASCII breaks it.
		"xychart": `"'#;[](){}` + emoji + `:,*-|%%`,
	}[diagram]
}

// punctuation is the hard part of a diagram's labels.
//
// EDGECASE_ONLY and EDGECASE_LABEL, read here and in main, are how the entries
// above were measured: one diagram and one character at a time.
func punctuation(diagram string) string {
	if v, ok := os.LookupEnv("EDGECASE_LABEL"); ok {
		return v
	}
	return supported(diagram)
}

// label is a label with words around the punctuation, for the diagram types
// that give a label a line of its own.
func label(diagram string) string {
	return "before " + punctuation(diagram) + " after"
}

// shortLabel is a label for the diagram types that lay their labels out on a
// grid, where words around the punctuation make the image unreadable without
// making the parse any harder.
func shortLabel(diagram string) string {
	if p := punctuation(diagram); p != "" {
		return "x" + p + "x"
	}
	// architecture takes no punctuation at all, and no space either.
	return "plainlabel"
}

// title is the diagram title. Some types write it into YAML front matter and
// others into a statement of their own, so it breaks on a different set of
// characters again.
//
// A line break is left out of it because the renderer honours one: the title
// then no longer reads as the text that was asked for, and the check that the
// title reaches the drawing cannot tell that apart from a title that was lost.
func title(diagram string) string {
	return "title " + strings.ReplaceAll(punctuation(diagram), "<br/>", "")
}

func main() {
	only := os.Getenv("EDGECASE_ONLY")
	for _, d := range diagrams() {
		if only != "" && only != d.file {
			continue
		}
		write(d)
	}
}

// write puts one diagram in a document of its own.
func write(d diagram) {
	f, err := os.Create(d.file + ".md")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	err = markdown.NewMarkdown(f, markdown.WithBlockSpacing()).
		H1(d.name+" rendering edge cases").
		PlainText("Every label below holds the punctuation this diagram type can carry. "+
			"Generated by doc/edgecase/main.go and rendered in continuous integration, "+
			"so a quoting defect fails the build rather than reaching a diagram.").
		CodeBlocks(markdown.SyntaxHighlightMermaid, d.build(d.file)).
		Build()
	if err != nil {
		panic(err)
	}
}

// diagram is one mermaid subpackage's document.
type diagram struct {
	// name is the diagram type, as a reader would say it.
	name string
	// file is the basename of the document, and the key into supported.
	file string
	// build returns the diagram, given the key its labels come from.
	build func(diagram string) string
}

// diagrams lists one document per mermaid subpackage.
func diagrams() []diagram {
	return []diagram{
		{name: "Architecture", file: "architecture", build: architecture},
		{name: "Block", file: "block", build: blockDiagram},
		{name: "C4 context", file: "c4", build: c4Context},
		{name: "Class", file: "class", build: classDiagram},
		{name: "Entity relationship", file: "er", build: entityRelationship},
		{name: "Flowchart", file: "flowchart", build: flowchartDiagram},
		{name: "Gantt", file: "gantt", build: ganttChart},
		{name: "Git graph", file: "gitgraph", build: gitGraph},
		{name: "Kanban", file: "kanban", build: kanbanDiagram},
		{name: "Mindmap", file: "mindmap", build: mindmapDiagram},
		{name: "Packet", file: "packet", build: packetDiagram},
		{name: "Pie chart", file: "piechart", build: pieChart},
		{name: "Quadrant", file: "quadrant", build: quadrantChart},
		{name: "Radar", file: "radar", build: radarChart},
		{name: "Requirement", file: "requirement", build: requirementDiagram},
		{name: "Sankey", file: "sankey", build: sankeyDiagram},
		{name: "Sequence", file: "sequence", build: sequenceDiagram},
		{name: "State", file: "state", build: stateDiagram},
		{name: "Treemap", file: "treemap", build: treemapDiagram},
		{name: "User journey", file: "userjourney", build: userJourney},
		{name: "XY chart", file: "xychart", build: xyChart},
	}
}

func architecture(diagram string) string {
	return arch.NewArchitecture(io.Discard).
		Group("cloudgroup", arch.IconCloud, shortLabel(diagram)).
		Service("api", arch.IconServer, shortLabel(diagram)).
		ServiceInGroup("db", arch.IconDatabase, shortLabel(diagram), "cloudgroup").
		Edges(
			arch.Edge{ServiceID: "api", Position: arch.PositionRight, Arrow: arch.ArrowRight},
			arch.Edge{ServiceID: "db", Position: arch.PositionLeft, Arrow: arch.ArrowNone},
		).
		String()
}

func blockDiagram(diagram string) string {
	return block.NewDiagram(io.Discard, block.WithTitle(title(diagram))).
		Columns(2). //nolint:mnd
		Row(
			block.Node("a", block.WithNodeLabel(label(diagram))),
			block.Node("b", block.WithNodeLabel(shortLabel(diagram)), block.WithNodeShape(block.ShapeRhombus)),
		).
		LinkWithLabel("a", shortLabel(diagram), "b").
		String()
}

// c4Context is the C4 context diagram's edge case document.
func c4Context(diagram string) string {
	return c4.NewDiagram(io.Discard, c4.WithTitle(title(diagram))).
		EnterpriseBoundary("enterprise", shortLabel(diagram)).
		Person("person", shortLabel(diagram), c4.WithDescription(label(diagram))).
		SystemBoundary("systems", shortLabel(diagram)+" nested").
		System("system", shortLabel(diagram), c4.WithDescription(label(diagram))).
		BoundaryEnd().
		BoundaryEnd().
		Boundary("outside", shortLabel(diagram), c4.WithBoundaryType(shortLabel(diagram))).
		SystemExt("external", shortLabel(diagram)).
		BoundaryEnd().
		PersonExt("auditor", shortLabel(diagram)).
		SystemDb("db", shortLabel(diagram)).
		SystemQueue("queue", shortLabel(diagram)).
		Rel("person", "system", label(diagram), c4.WithTechnology(shortLabel(diagram))).
		BiRel("system", "db", label(diagram)).
		String()
}

func classDiagram(diagram string) string {
	return class.NewDiagram(io.Discard, class.WithTitle(title(diagram))).
		Class("Account", class.WithPublicField("string", "id")).
		ClassWithLabel("Ledger", label(diagram)).
		Relation("Account", class.RelationshipAssociation, "Ledger").
		RelationWithLabel("Account", class.RelationshipDependency, "Ledger", shortLabel(diagram)).
		NoteFor("Account", label(diagram)).
		String()
}

func entityRelationship(diagram string) string {
	teachers := er.NewEntity("teachers", []*er.Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, Comment: label(diagram)},
	})
	students := er.NewEntity("students", []*er.Attribute{
		{Type: "int", Name: "teacher_id", IsForeignKey: true, Comment: label(diagram)},
	})

	return er.NewDiagram(io.Discard).
		Relationship(teachers, students, er.ExactlyOneRelationship, er.ZeroToMoreRelationship, er.Identifying, label(diagram)).
		String()
}

func flowchartDiagram(diagram string) string {
	return flowchart.NewFlowchart(io.Discard, flowchart.WithTitle(title(diagram))).
		NodeWithText("A", label(diagram)).
		RhombusNode("B", shortLabel(diagram)).
		CircleNode("C", shortLabel(diagram)).
		LinkWithArrowHeadAndText("A", "B", shortLabel(diagram)).
		DottedLinkWithText("B", "C", shortLabel(diagram)).
		String()
}

func ganttChart(diagram string) string {
	return gantt.NewChart(
		io.Discard,
		gantt.WithTitle(title(diagram)),
		gantt.WithDateFormat("YYYY-MM-DD"),
	).
		Section(shortLabel(diagram)).
		Task(shortLabel(diagram), "2024-01-01", "2d").
		MilestoneWithID(shortLabel(diagram), "milestone", "2024-01-03").
		String()
}

func gitGraph(diagram string) string {
	return gitgraph.NewDiagram(io.Discard, gitgraph.WithTitle(title(diagram))).
		Commit(gitgraph.WithCommitID(shortLabel(diagram)), gitgraph.WithCommitTag(shortLabel(diagram))).
		Branch("develop").
		Checkout("develop").
		Commit(gitgraph.WithCommitID("second")).
		String()
}

func kanbanDiagram(diagram string) string {
	return kanban.NewDiagram(io.Discard, kanban.WithTitle(title(diagram))).
		Column(shortLabel(diagram), kanban.WithColumnID("todo")).
		Task(label(diagram), kanban.WithTaskAssigned(shortLabel(diagram)), kanban.WithTaskPriority(kanban.PriorityHigh)).
		String()
}

func mindmapDiagram(diagram string) string {
	return mindmap.NewDiagram(io.Discard, mindmap.WithTitle(title(diagram))).
		Root(shortLabel(diagram)).
		Child(shortLabel(diagram)).
		Sibling(shortLabel(diagram)).
		String()
}

func packetDiagram(diagram string) string {
	return packet.NewDiagram(io.Discard, packet.WithTitle(title(diagram))).
		Field(0, 15, shortLabel(diagram)). //nolint:mnd
		Next(16, shortLabel(diagram)).     //nolint:mnd
		String()
}

func pieChart(diagram string) string {
	return piechart.NewPieChart(io.Discard, piechart.WithTitle(title(diagram)), piechart.WithShowData(true)).
		LabelAndIntValue(shortLabel(diagram), 120).           //nolint:mnd
		LabelAndFloatValue(shortLabel(diagram)+" two", 42.5). //nolint:mnd
		String()
}

func quadrantChart(diagram string) string {
	return quadrant.NewChart(io.Discard, quadrant.WithTitle(title(diagram))).
		XAxis(shortLabel(diagram), shortLabel(diagram)+" high").
		YAxis(shortLabel(diagram), shortLabel(diagram)+" high").
		Quadrant1(shortLabel(diagram)).
		Point(shortLabel(diagram), 0.3, 0.6). //nolint:mnd
		String()
}

func radarChart(diagram string) string {
	return radar.NewDiagram(io.Discard, radar.WithTitle(title(diagram))).
		Axis(shortLabel(diagram), shortLabel(diagram)+" two", label(diagram)).
		Curve(shortLabel(diagram), 85, 90, 80). //nolint:mnd
		Curve(label(diagram), 70, 75, 85).      //nolint:mnd
		Max(100).                               //nolint:mnd
		Min(0).
		String()
}

func requirementDiagram(diagram string) string {
	return requirement.NewDiagram(io.Discard, requirement.WithTitle(title(diagram))).
		Requirement(
			shortLabel(diagram),
			requirement.WithID("1"),
			requirement.WithText(label(diagram)),
			requirement.WithRisk(requirement.RiskLow),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest),
		).
		Element(shortLabel(diagram)+" element", requirement.WithElementType("simulation"), requirement.WithDocRef("./tests")).
		Satisfies(shortLabel(diagram)+" element", shortLabel(diagram)).
		String()
}

func sankeyDiagram(diagram string) string {
	// No title: mermaid parses one for a sankey diagram and does not draw it,
	// which the renderer check reads as a title that went missing.
	return sankey.NewDiagram(io.Discard).
		Link(shortLabel(diagram), shortLabel(diagram)+" target", 100). //nolint:mnd
		Link(shortLabel(diagram)+" target", label(diagram), 42.5).     //nolint:mnd
		String()
}

func sequenceDiagram(diagram string) string {
	return sequence.NewDiagram(io.Discard).
		Participant(shortLabel(diagram)).
		Participant(shortLabel(diagram)+" two").
		SyncRequest(shortLabel(diagram), shortLabel(diagram)+" two", label(diagram)).
		SyncResponse(shortLabel(diagram)+" two", shortLabel(diagram), label(diagram)).
		NoteOver(shortLabel(diagram), label(diagram)).
		String()
}

func stateDiagram(diagram string) string {
	return state.NewDiagram(io.Discard, state.WithTitle(title(diagram))).
		State("Draft", label(diagram)).
		State("Placed", shortLabel(diagram)).
		Transition("Draft", "Placed").
		TransitionWithNote("Placed", "Draft", shortLabel(diagram)).
		NoteRight("Draft", label(diagram)).
		String()
}

func treemapDiagram(diagram string) string {
	return treemap.NewDiagram(io.Discard, treemap.WithTitle(title(diagram))).
		Section(shortLabel(diagram)).
		Leaf(label(diagram), 1200). //nolint:mnd
		Section(shortLabel(diagram)+" nested").
		Leaf(shortLabel(diagram), 400). //nolint:mnd
		Parent().
		Parent().
		Leaf(shortLabel(diagram)+" top", 300). //nolint:mnd
		String()
}

func userJourney(diagram string) string {
	return userjourney.NewDiagram(io.Discard, userjourney.WithTitle(title(diagram))).
		Section(shortLabel(diagram)).
		Task(shortLabel(diagram), userjourney.ScoreSatisfied, shortLabel(diagram)).
		String()
}

func xyChart(diagram string) string {
	return xychart.NewDiagram(io.Discard, xychart.WithTitle(title(diagram))).
		XAxisLabelsWithTitle(shortLabel(diagram), shortLabel(diagram), shortLabel(diagram)+" two").
		YAxisRangeWithTitle(shortLabel(diagram), 0, 100). //nolint:mnd
		Bar(10, 20).                                      //nolint:mnd
		Line(30, 40).                                     //nolint:mnd
		String()
}
