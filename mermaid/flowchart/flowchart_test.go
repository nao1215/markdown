// Package flowchart provides a simple way to create flowcharts in mermaid syntax.
package flowchart

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
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
title: "mermaid flowchart builder"
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

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewFlowchart(w).NodeWithText("A", "Node A")
	})
}

// TestGoldenFlowchart pins the rendered diagram of every node shape and every
// link style this package can build.
func TestGoldenFlowchart(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewFlowchart(
		buf,
		WithTitle("Every Shape And Link"),
		WithOrientalTopToBottom(),
	).
		Node("plain").
		NodeWithText("text", "With text").
		NodeWithMarkdown("markdown", "**bold** text").
		NodeWithNewLines("newlines", "first\nsecond").
		RoundEdgesNode("round", "Round edges").
		StadiumNode("stadium", "Stadium").
		SubroutineNode("subroutine", "Subroutine").
		CylindricalNode("cylindrical", "Cylindrical").
		DatabaseNode("database", "Database").
		CircleNode("circle", "Circle").
		AsymmetricNode("asymmetric", "Asymmetric").
		RhombusNode("rhombus", "Rhombus").
		HexagonNode("hexagon", "Hexagon").
		ParallelogramNode("parallelogram", "Parallelogram").
		ParallelogramAltNode("parallelogramAlt", "Parallelogram alt").
		TrapezoidNode("trapezoid", "Trapezoid").
		TrapezoidAltNode("trapezoidAlt", "Trapezoid alt").
		DoubleCircleNode("doubleCircle", "Double circle").
		LinkWithArrowHead("plain", "text").
		LinkWithArrowHeadAndText("text", "markdown", "with text").
		OpenLink("markdown", "newlines").
		OpenLinkWithText("newlines", "round", "open with text").
		DottedLink("round", "stadium").
		DottedLinkWithText("stadium", "subroutine", "dotted with text").
		ThickLink("subroutine", "cylindrical").
		ThickLinkWithText("cylindrical", "database", "thick with text").
		InvisibleLink("database", "circle").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("flowchart.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenFlowchartOrientations pins the header each orientation option
// produces. The options are mutually exclusive, so each needs its own diagram.
func TestGoldenFlowchartOrientations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		golden string
		option Option
	}{
		{name: "top to bottom", golden: "orientation_tb.md", option: WithOrientalTopToBottom()},
		{name: "top down", golden: "orientation_td.md", option: WithOrientalTopDown()},
		{name: "bottom to top", golden: "orientation_bt.md", option: WithOrientalBottomToTop()},
		{name: "right to left", golden: "orientation_rl.md", option: WithOrientalRightToLeft()},
		{name: "left to right", golden: "orientation_lr.md", option: WithOrientalLeftToRight()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			err := NewFlowchart(buf, tt.option).
				NodeWithText("A", "Start").
				NodeWithText("B", "End").
				LinkWithArrowHead("A", "B").
				Build()
			if err != nil {
				t.Fatalf("Build() = %v, want nil", err)
			}

			if err := golden.Assert(tt.golden, buf.String()); err != nil {
				t.Error(err)
			}
		})
	}
}

// TestBuildWithNilWriter covers the case where a flowchart is built for String()
// only and Build() is called by mistake. Build() used to dereference the nil
// writer and take the process down; it has to return an error instead.
func TestBuildWithNilWriter(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Build() panicked with a nil writer: %v", r)
		}
	}()

	d := NewFlowchart(nil)

	// String() has always worked without a writer, and callers rely on it.
	_ = d.String()

	err := d.Build()
	if err == nil {
		t.Fatal("Build() with a nil writer must return an error")
	}
	if err.Error() != "output writer must not be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// errWrite is the failure the writer below reports, so the test can assert that
// Build passed it through rather than inventing an error of its own.
var errWrite = errors.New("write failed")

// errWriter fails every write, which is what a full disk or a closed pipe looks
// like to Build.
type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errWrite
}

// TestBuildReportsWriteFailure covers the branch where the destination accepts
// the diagram and then fails. Silently returning nil there would hand the caller
// a document that was never written.
func TestBuildReportsWriteFailure(t *testing.T) {
	t.Parallel()

	err := NewFlowchart(errWriter{}).Build()
	if err == nil {
		t.Fatal("Build must report a failing writer")
	}
	if !errors.Is(err, errWrite) {
		t.Errorf("Build lost the destination error: %v", err)
	}
}

// TestNodeTextEscapesTheQuoteThatEndsIt names the character this escaping
// buys. A double quote in a node label used to reach mermaid unescaped and lose
// the whole diagram: the reader got an error box rather than a picture.
func TestNodeTextEscapesTheQuoteThatEndsIt(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(*Flowchart) *Flowchart
		want  string
	}{
		"a quote in node text": {
			build: func(f *Flowchart) *Flowchart { return f.NodeWithText("A", `say "hi"`) },
			want:  `    A["say #quot;hi#quot;"]`,
		},
		"a quote in a markdown node": {
			build: func(f *Flowchart) *Flowchart { return f.NodeWithMarkdown("A", `say "hi"`) },
			want:  "    A[\"`say #quot;hi#quot;`\"]",
		},
		"a quote in every other node shape": {
			build: func(f *Flowchart) *Flowchart { return f.HexagonNode("A", `"`) },
			want:  `    A{{"#quot;"}}`,
		},
		"a quote in link text": {
			build: func(f *Flowchart) *Flowchart { return f.LinkWithArrowHeadAndText("A", "B", `"`) },
			want:  `    A-->|"#quot;"|B`,
		},
		"a quote in dotted link text": {
			build: func(f *Flowchart) *Flowchart { return f.DottedLinkWithText("A", "B", `"`) },
			want:  `    A-. "#quot;" .-> B`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.build(NewFlowchart(io.Discard)).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("diagram =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}

// TestNodeTextEscapesOnlyTheHashThatStartsAnEntity pins the other half of the
// escaping. mermaid reads "#quot;" and "#123;" as the characters they name, so
// a label holding one has to be escaped or it would draw the same diagram as a
// label holding a quotation mark. A "#" anywhere else is ordinary text, and its
// output is left exactly as it was.
func TestNodeTextEscapesOnlyTheHashThatStartsAnEntity(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		text string
		want string
	}{
		"a plain hash is left alone":            {text: "PR #123 merged", want: `    A["PR #123 merged"]`},
		"a hash at the end is left alone":       {text: "issue #", want: `    A["issue #"]`},
		"a hash before a semicolon is left too": {text: "a#;b", want: `    A["a#;b"]`},
		"a named entity is escaped":             {text: "a#quot;b", want: `    A["a#35;quot;b"]`},
		"a numeric entity is escaped":           {text: "a#123;b", want: `    A["a#35;123;b"]`},
		"a quote and an entity together":        {text: `#39;"`, want: `    A["#35;39;#quot;"]`},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := NewFlowchart(io.Discard).NodeWithText("A", tt.text).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("diagram =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}

// TestErrorReportsTheRecordedError pins the method the v1.0.0 API audit found
// missing. Every other builder in this library lets the recorded error be read
// before anything is written, and this one did not.
func TestErrorReportsTheRecordedError(t *testing.T) {
	t.Parallel()

	f := NewFlowchart(nil).NodeWithText("A", "Start")

	if err := f.Error(); err != nil {
		t.Errorf("Error() = %v before Build, want nil", err)
	}

	// Build reports the same error when it is what stopped the write.
	fromBuild := f.Build()
	if fromBuild == nil {
		t.Fatal("Build() = nil with a nil writer, want an error")
	}
	if f.Error() == nil || f.Error().Error() != fromBuild.Error() {
		t.Errorf("Error() = %v, want the error Build returned, %v", f.Error(), fromBuild)
	}
}
