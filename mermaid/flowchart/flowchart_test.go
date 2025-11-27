// Package flowchart provides a simple way to create flowcharts in mermaid syntax.
package flowchart

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFlowchart_Build(t *testing.T) {
	t.Parallel()

	t.Run("Build a flowchart with title", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		f := NewFlowchart(
			b,
			WithTitle("mermaid flowchart builder"),
			WithOrientalTopToBottom(),
		).
			NodeWithText("A", "Node A").
			StadiumNode("B", "Node B").
			SubroutineNode("C", "Node C").
			DatabaseNode("D", "Database").
			LinkWithArrowHead("A", "B").
			LinkWithArrowHeadAndText("B", "D", "send original data").
			LinkWithArrowHead("B", "C").
			DottedLinkWithText("C", "D", "send filtered data")

		if err := f.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `---
title: mermaid flowchart builder
---
flowchart TB
    A["Node A"]
    B(["Node B"])
    C[["Node C"]]
    D[("Database")]
    A-->B
    B-->|"send original data"|D
    B-->C
    C-. "send filtered data" .-> D`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a flowchart, top to bottom", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		f := NewFlowchart(
			b,
			WithOrientalTopToBottom(),
		).NodeWithText("A", "Node A")

		if err := f.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `flowchart TB
    A["Node A"]`
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a flowchart, top down", func(t *testing.T) {
		t.Parallel()
		b := new(bytes.Buffer)

		f := NewFlowchart(
			b,
			WithOrientalTopDown(),
		).NodeWithText("A", "Node A")

		if err := f.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `flowchart TD
    A["Node A"]`
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a flowchart, bottom to top", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		f := NewFlowchart(
			b,
			WithOrientalBottomToTop(),
		).NodeWithText("A", "Node A")

		if err := f.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `flowchart BT
    A["Node A"]`
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a flowchart, right to left", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		f := NewFlowchart(
			b,
			WithOrientalRightToLeft(),
		).NodeWithText("A", "Node A")

		if err := f.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `flowchart RL
    A["Node A"]`
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a flowchart, left to right", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		f := NewFlowchart(
			b,
			WithOrientalLeftToRight(),
		).NodeWithText("A", "Node A")

		if err := f.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `flowchart LR
    A["Node A"]`
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a flowchart with all node and link", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		f := NewFlowchart(
			b,
			WithOrientalTopToBottom(),
		).
			Node("A").
			NodeWithText("B", "Node B").
			NodeWithMarkdown("C", "**Node C**").
			NodeWithNewLines("D", `Node
D`).RoundEdgesNode("E", "Node E").
			StadiumNode("F", "Node F").
			SubroutineNode("G", "Node G").
			CylindricalNode("H", "Node H").
			DatabaseNode("I", "Database").
			CircleNode("J", "Node J").
			AsymmetricNode("K", "Node K").
			RhombusNode("L", "Node L").
			HexagonNode("M", "Node M").
			ParallelogramNode("N", "Node N").
			ParallelogramAltNode("O", "Node O").
			TrapezoidNode("P", "Node P").
			TrapezoidAltNode("Q", "Node Q").
			DoubleCircleNode("R", "Node R").
			LinkWithArrowHead("A", "B").
			LinkWithArrowHeadAndText("B", "C", "send").
			OpenLink("C", "D").
			OpenLinkWithText("D", "E", "send").
			DottedLink("E", "F").
			DottedLinkWithText("F", "G", "send").
			ThickLink("G", "H").
			ThickLinkWithText("H", "I", "send").
			InvisibleLink("I", "J")

		if err := f.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := `flowchart TB
    A
    B["Node B"]
    `
		want += "C[\"`**Node C**`\"]\n"
		want += "    D[\"`Node\n"
		want += "D`\"]"
		want += `
    E("Node E")
    F(["Node F"])
    G[["Node G"]]
    H[("Node H")]
    I[("Database")]
    J(("Node J"))
    K>"Node K"]
    L{"Node L"}
    M{{"Node M"}}
    N[/"Node N"/]
    O[\"Node O"\]
    P[/"Node P"\]
    Q[\"Node Q"/]
    R((("Node R")))
    A-->B
    B-->|"send"|C
    C --- D
    D---|"send"|E
    E-.->F
    F-. "send" .-> G
    G ==> H
    H == "send" ==> I
    I ~~~ J`

		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}
