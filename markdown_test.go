package markdown

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal"
)

func TestPlainText(t *testing.T) {
	t.Parallel()

	t.Run("success PlainText()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.PlainText("Hello")
		want := []string{"Hello"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestMarkdownHeader(t *testing.T) {
	t.Parallel()

	t.Run("success H1f()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1f("%s", "Hello")
		want := "# Hello"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success H2f()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H2f("%s", "Hello")
		want := "## Hello"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success H3f()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H3f("%s", "Hello")
		want := "### Hello"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success H4f()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H4f("%s", "Hello")
		want := "#### Hello"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success H5f()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H5f("%s", "Hello")
		want := "##### Hello"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success H6f()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H6f("%s", "Hello")
		want := "###### Hello"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestMarkdownDetailsf(t *testing.T) {
	t.Parallel()

	t.Run("success Detailsf()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.Detailsf("Hello", "Good %s", "World")
		want := fmt.Sprintf("<details><summary>Hello</summary>%sGood World%s</details>", internal.LineFeed(), internal.LineFeed())
		got := m.body[0]

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestMarkdownBulletList(t *testing.T) {
	t.Parallel()

	t.Run("success BulletList()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.BulletList("Hello", "World")
		want := []string{"- Hello", "- World"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestMarkdownNumberList(t *testing.T) {
	t.Parallel()

	t.Run("success NumberList()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.OrderedList("Hello", "World")
		want := []string{"1. Hello", "2. World"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestMarkdownCheckBoxf(t *testing.T) {
	t.Run("success CheckBoxf(); check [x]", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		set := []CheckBoxSet{
			{Text: "Hello", Checked: true},
			{Text: "World", Checked: false},
		}
		m.CheckBox(set)
		want := []string{
			"- [x] Hello",
			"- [ ] World",
		}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestMarkdownBlockquote(t *testing.T) {
	t.Parallel()

	t.Run("success Blockquote()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.Blockquote(fmt.Sprintf("%s%s%s%s%s", "Hello", internal.LineFeed(), "Good", internal.LineFeed(), "World"))
		want := []string{
			"> Hello",
			"> Good",
			"> World",
		}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestMarkdownCodeBlocks(t *testing.T) {
	t.Run("success CodeBlock()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.CodeBlocks(SyntaxHighlightGo, "Hello")
		want := []string{fmt.Sprintf("```go%sHello%s```", internal.LineFeed(), internal.LineFeed())}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestHorizontalRule(t *testing.T) {
	t.Parallel()

	t.Run("success HorizontalRule()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.HorizontalRule()
		want := []string{"---"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestTableSetValidateColumns(t *testing.T) {
	t.Parallel()
	t.Run("success TableSet.ValidateColumns()", func(t *testing.T) {
		t.Parallel()

		set := TableSet{
			Header: []string{"Name", "Age"},
			Rows:   [][]string{{"David", "23"}},
		}

		err := set.ValidateColumns()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})

	t.Run("failed TableSet.ValidateColumns(); invalid header", func(t *testing.T) {
		t.Parallel()

		set := TableSet{
			Header: []string{"Name", "Age"},
			Rows:   [][]string{{"David"}},
		}

		err := set.ValidateColumns()
		if err == nil {
			t.Error("expected error, but not occurred")
		}
	})
}

func TestMarkdownTable(t *testing.T) {
	t.Parallel()

	t.Run("success Table() without alignment", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		set := TableSet{
			Header: []string{"Name", "Age"},
			Rows:   [][]string{{"David", "23"}},
		}
		m.Table(set)
		want := []string{
			fmt.Sprintf("| Name | Age |%s|---------|---------|%s| David | 23 |%s",
				internal.LineFeed(), internal.LineFeed(), internal.LineFeed()),
		}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Table() with left alignment", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		set := TableSet{
			Header:    []string{"Left Align", "Normal"},
			Rows:      [][]string{{"Content1", "Content2"}},
			Alignment: []TableAlignment{AlignLeft},
		}
		m.Table(set)
		want := []string{
			fmt.Sprintf("| Left Align | Normal |%s|:--------|---------|%s| Content1 | Content2 |%s",
				internal.LineFeed(), internal.LineFeed(), internal.LineFeed()),
		}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Table() with center alignment", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		set := TableSet{
			Header:    []string{"Center Align", "Normal"},
			Rows:      [][]string{{"Content1", "Content2"}},
			Alignment: []TableAlignment{AlignCenter},
		}
		m.Table(set)
		want := []string{
			fmt.Sprintf("| Center Align | Normal |%s|:-------:|---------|%s| Content1 | Content2 |%s",
				internal.LineFeed(), internal.LineFeed(), internal.LineFeed()),
		}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Table() with right alignment", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		set := TableSet{
			Header:    []string{"Right Align", "Normal"},
			Rows:      [][]string{{"Content1", "Content2"}},
			Alignment: []TableAlignment{AlignRight},
		}
		m.Table(set)
		want := []string{
			fmt.Sprintf("| Right Align | Normal |%s|--------:|---------|%s| Content1 | Content2 |%s",
				internal.LineFeed(), internal.LineFeed(), internal.LineFeed()),
		}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Table() with mixed alignments", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		set := TableSet{
			Header:    []string{"Left Align", "Center Align", "Right Align"},
			Rows:      [][]string{{"Content1", "Content2", "Content3"}, {"Content4", "Content5", "Content6"}},
			Alignment: []TableAlignment{AlignLeft, AlignCenter, AlignRight},
		}
		m.Table(set)
		want := []string{
			fmt.Sprintf("| Left Align | Center Align | Right Align |%s|:--------|:-------:|--------:|%s| Content1 | Content2 | Content3 |%s| Content4 | Content5 | Content6 |%s",
				internal.LineFeed(), internal.LineFeed(), internal.LineFeed(), internal.LineFeed()),
		}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Table() with partial alignment specification", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		set := TableSet{
			Header:    []string{"Left", "Default", "Center"},
			Rows:      [][]string{{"A", "B", "C"}},
			Alignment: []TableAlignment{AlignLeft, AlignCenter}, // Only specify first 2 columns
		}
		m.Table(set)
		want := []string{
			fmt.Sprintf("| Left | Default | Center |%s|:--------|:-------:|---------|%s| A | B | C |%s",
				internal.LineFeed(), internal.LineFeed(), internal.LineFeed()),
		}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("empty table headers", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		set := TableSet{
			Header: []string{},
			Rows:   [][]string{},
		}
		m.Table(set)
		want := []string{}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestMarkdownBuildError(t *testing.T) {
	t.Parallel()

	t.Run("Error() return nil", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		if err := m.H1("sample").Build(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})

	t.Run("Error() return error", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.Table(TableSet{
			Header: []string{"Name", "Age"},
			Rows:   [][]string{{"David"}},
		})
		if err := m.Build(); err == nil {
			t.Error("expected error, but not occurred")
		}
	})
}

func TestMarkdownLF(t *testing.T) {
	t.Parallel()
	t.Run("success Markdown.LF()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.LF()
		want := []string{"  "}
		got := m.body

		if !reflect.DeepEqual(want, got) {
			t.Errorf("value is mismatch want: %v, got: %v", want, got)
		}
	})
}

func TestMarkdownError(t *testing.T) {
	t.Parallel()

	t.Run("Error() return nil", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		if err := m.H1("sample").Error(); err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})

	t.Run("Error() return error", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.Table(TableSet{
			Header: []string{"Name", "Age"},
			Rows:   [][]string{{"David"}},
		})
		if err := m.Error(); err == nil {
			t.Error("expected error, but not occurred")
		}
	})

	t.Run("Error() return error Custom Table", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.CustomTable(TableSet{
			Header: []string{"Name", "Age"},
			Rows:   [][]string{{"David"}},
		}, TableOptions{
			AutoWrapText: false,
		})
		if err := m.Error(); err == nil {
			t.Error("expected error, but not occurred")
		}
	})
}

func TestMarkdownCustomTable(t *testing.T) {
	t.Parallel()
	t.Run("success Table()", func(t *testing.T) {
		t.Parallel()

		if runtime.GOOS == "windows" {
			t.Skip("Skip test on Windows due to line feed mismatch")
		}

		m := NewMarkdown(os.Stdout)
		set := TableSet{
			Header: []string{"Name", "Age"},
			Rows:   [][]string{{"David", "23"}},
		}
		m.CustomTable(set, TableOptions{
			AutoWrapText:      false,
			AutoFormatHeaders: false,
		})
		want := []string{
			fmt.Sprintf("| Name  | Age |%s|-------|-----|%s| David | 23  |%s",
				internal.LineFeed(), internal.LineFeed(), internal.LineFeed()),
		}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestTableOfContents(t *testing.T) {
	t.Parallel()

	t.Run("TableOfContents places table of contents at correct position", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Title").
			TableOfContents(TableOfContentsDepthH3).
			H2("Section 1").
			H3("Subsection 1.1").
			H2("Section 2")

		want := "# Title" + internal.LineFeed() +
			TableOfContentsMarkerBegin + internal.LineFeed() +
			"- [Title](#title)" + internal.LineFeed() +
			"  - [Section 1](#section-1)" + internal.LineFeed() +
			"    - [Subsection 1.1](#subsection-11)" + internal.LineFeed() +
			"  - [Section 2](#section-2)" + internal.LineFeed() +
			TableOfContentsMarkerEnd + internal.LineFeed() +
			"" + internal.LineFeed() +
			"## Section 1" + internal.LineFeed() +
			"### Subsection 1.1" + internal.LineFeed() +
			"## Section 2"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("TableOfContents prevents duplicate generation", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Title").
			TableOfContents(TableOfContentsDepthH2).
			H2("Section").
			TableOfContents(TableOfContentsDepthH3)

		if m.Error() == nil {
			t.Error("expected error when generating table of contents twice")
		}
	})
}

func TestTableOfContentsWithRange(t *testing.T) {
	t.Parallel()

	t.Run("Table of contents with custom range excludes H1", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Document Title").
			H2("Table of Contents").
			TableOfContentsWithRange(TableOfContentsDepthH2, TableOfContentsDepthH4).
			H2("Introduction").
			H3("Overview").
			H4("Details").
			H5("Deep Details")

		want := "# Document Title" + internal.LineFeed() +
			"## Table of Contents" + internal.LineFeed() +
			TableOfContentsMarkerBegin + internal.LineFeed() +
			"- [Table of Contents](#table-of-contents)" + internal.LineFeed() +
			"- [Introduction](#introduction)" + internal.LineFeed() +
			"  - [Overview](#overview)" + internal.LineFeed() +
			"    - [Details](#details)" + internal.LineFeed() +
			TableOfContentsMarkerEnd + internal.LineFeed() +
			"" + internal.LineFeed() +
			"## Introduction" + internal.LineFeed() +
			"### Overview" + internal.LineFeed() +
			"#### Details" + internal.LineFeed() +
			"##### Deep Details"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Table of contents with range H3 to H5", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Title").
			H2("Section").
			H3("Subsection").
			TableOfContentsWithRange(TableOfContentsDepthH3, TableOfContentsDepthH5).
			H4("Detail").
			H5("Deep Detail").
			H6("Very Deep Detail")

		want := "# Title" + internal.LineFeed() +
			"## Section" + internal.LineFeed() +
			"### Subsection" + internal.LineFeed() +
			TableOfContentsMarkerBegin + internal.LineFeed() +
			"- [Subsection](#subsection)" + internal.LineFeed() +
			"  - [Detail](#detail)" + internal.LineFeed() +
			"    - [Deep Detail](#deep-detail)" + internal.LineFeed() +
			TableOfContentsMarkerEnd + internal.LineFeed() +
			"" + internal.LineFeed() +
			"#### Detail" + internal.LineFeed() +
			"##### Deep Detail" + internal.LineFeed() +
			"###### Very Deep Detail"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Invalid depth ranges return errors", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name                string
			minDepth, maxDepth  TableOfContentsDepth
			expectedErrorSubstr string
		}{
			{"minDepth too low", 0, 3, "invalid minDepth: 0"},
			{"minDepth too high", 7, 6, "invalid minDepth: 7"},
			{"maxDepth too low", 1, 0, "invalid maxDepth: 0"},
			{"maxDepth too high", 1, 7, "invalid maxDepth: 7"},
			{"minDepth > maxDepth", 4, 2, "minDepth (4) cannot be greater than maxDepth (2)"},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				m := NewMarkdown(os.Stdout)
				m.H1("Test").TableOfContentsWithRange(tt.minDepth, tt.maxDepth)

				if m.Error() == nil {
					t.Errorf("expected error for %s", tt.name)
				} else if !strings.Contains(m.Error().Error(), tt.expectedErrorSubstr) {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectedErrorSubstr, m.Error().Error())
				}
			})
		}
	})

	t.Run("Empty table of contents when no headers in range", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Title").
			TableOfContentsWithRange(TableOfContentsDepthH3, TableOfContentsDepthH5).
			PlainText("Some content")

		want := "# Title" + internal.LineFeed() +
			TableOfContentsMarkerBegin + internal.LineFeed() +
			TableOfContentsMarkerEnd + internal.LineFeed() +
			"" + internal.LineFeed() +
			"Some content"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestTableOfContentsWithSpecialCharacters(t *testing.T) {
	t.Parallel()

	t.Run("Table of contents handles special characters in headers", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("API & SDK").
			H2("C++ Programming").
			H3("Configuration: Advanced Settings").
			TableOfContents(TableOfContentsDepthH3)

		want := "# API & SDK" + internal.LineFeed() +
			"## C++ Programming" + internal.LineFeed() +
			"### Configuration: Advanced Settings" + internal.LineFeed() +
			TableOfContentsMarkerBegin + internal.LineFeed() +
			"- [API & SDK](#api--sdk)" + internal.LineFeed() +
			"  - [C++ Programming](#c-programming)" + internal.LineFeed() +
			"    - [Configuration: Advanced Settings](#configuration-advanced-settings)" + internal.LineFeed() +
			TableOfContentsMarkerEnd + internal.LineFeed() +
			""
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestTableOfContentsMethodCompatibility(t *testing.T) {
	t.Parallel()

	t.Run("TableOfContents method works correctly", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Title").
			H2("Section").
			H3("Subsection").
			TableOfContents(TableOfContentsDepthH2)

		want := "# Title" + internal.LineFeed() +
			"## Section" + internal.LineFeed() +
			"### Subsection" + internal.LineFeed() +
			TableOfContentsMarkerBegin + internal.LineFeed() +
			"- [Title](#title)" + internal.LineFeed() +
			"  - [Section](#section)" + internal.LineFeed() +
			TableOfContentsMarkerEnd + internal.LineFeed() +
			""
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Table of contents usage example from your description", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("MyTitle").
			H2("Table of contents").
			TableOfContentsWithRange(TableOfContentsDepthH2, TableOfContentsDepthH5).
			H2("Section 1").
			H2("Section 2")

		want := "# MyTitle" + internal.LineFeed() +
			"## Table of contents" + internal.LineFeed() +
			TableOfContentsMarkerBegin + internal.LineFeed() +
			"- [Table of contents](#table-of-contents)" + internal.LineFeed() +
			"- [Section 1](#section-1)" + internal.LineFeed() +
			"- [Section 2](#section-2)" + internal.LineFeed() +
			TableOfContentsMarkerEnd + internal.LineFeed() +
			"" + internal.LineFeed() +
			"## Section 1" + internal.LineFeed() +
			"## Section 2"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

// Tests for header format methods that were uncovered
func TestHeaderFormatMethods(t *testing.T) {
	t.Parallel()

	t.Run("H1f method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1f("Title %s", "Test")

		want := "# Title Test"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("H3f method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H3f("Section %d.%d", 1, 1)

		want := "### Section 1.1"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("H4f method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H4f("Subsection %s", "A")

		want := "#### Subsection A"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("H5f method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H5f("Detail %s", "X")

		want := "##### Detail X"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("H6f method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H6f("Deep %s", "Content")

		want := "###### Deep Content"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

// Tests for Details methods
func TestDetailsMethod(t *testing.T) {
	t.Parallel()

	t.Run("Details method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.Details("Summary", "Hidden content")

		want := "<details><summary>Summary</summary>" + internal.LineFeed() + "Hidden content" + internal.LineFeed() + "</details>"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Detailsf method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.Detailsf("Summary", "Hidden %s %d", "content", 42)

		want := "<details><summary>Summary</summary>" + internal.LineFeed() + "Hidden content 42" + internal.LineFeed() + "</details>"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

// Tests for PlainTextf method
func TestPlainTextfMethod(t *testing.T) {
	t.Parallel()

	t.Run("PlainTextf method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.PlainTextf("Hello %s, you have %d messages", "John", 5)

		want := "Hello John, you have 5 messages"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

// Tests for list methods
func TestListMethods(t *testing.T) {
	t.Parallel()

	t.Run("BulletList with multiple items", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.BulletList("Item 1", "Item 2", "Item 3")

		want := "- Item 1" + internal.LineFeed() + "- Item 2" + internal.LineFeed() + "- Item 3"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("OrderedList with multiple items", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.OrderedList("First", "Second", "Third")

		want := "1. First" + internal.LineFeed() + "2. Second" + internal.LineFeed() + "3. Third"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("CheckBox with mixed states", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		checkboxes := []CheckBoxSet{
			{Checked: true, Text: "Completed task"},
			{Checked: false, Text: "Incomplete task"},
		}
		m.CheckBox(checkboxes)

		want := "- [x] Completed task" + internal.LineFeed() + "- [ ] Incomplete task"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

// Tests for other methods
func TestOtherMethods(t *testing.T) {
	t.Parallel()

	t.Run("Blockquote method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.Blockquote("This is a quote")

		want := "> This is a quote"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("CodeBlocks method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.CodeBlocks(SyntaxHighlightGo, "func main() {}")

		want := "```go" + internal.LineFeed() + "func main() {}" + internal.LineFeed() + "```"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("CodeBlocks with no syntax highlighting", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.CodeBlocks(SyntaxHighlightNone, "plain text")

		want := "```" + internal.LineFeed() + "plain text" + internal.LineFeed() + "```"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

// Tests for Build method and error handling
func TestBuildMethodAndErrors(t *testing.T) {
	t.Parallel()

	t.Run("Build method success", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		m := NewMarkdown(&buf)
		m.H1("Test")

		err := m.Build()
		if err != nil {
			t.Errorf("Build() returned error: %v", err)
		}

		want := "# Test"
		got := buf.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Error method returns nil when no error", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Test")

		if err := m.Error(); err != nil {
			t.Errorf("Error() should return nil when no error, got: %v", err)
		}
	})
}

// Tests for table methods
func TestTableMethods(t *testing.T) {
	t.Parallel()

	t.Run("Table method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		tableSet := TableSet{
			Header: []string{"Name", "Age"},
			Rows: [][]string{
				{"Alice", "30"},
				{"Bob", "25"},
			},
		}
		m.Table(tableSet)

		got := m.String()

		// Check that it contains the expected table elements
		if !strings.Contains(got, "Name") || !strings.Contains(got, "Age") {
			t.Errorf("Table should contain headers")
		}
		if !strings.Contains(got, "Alice") || !strings.Contains(got, "Bob") {
			t.Errorf("Table should contain row data")
		}
	})

	t.Run("CustomTable method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		tableSet := TableSet{
			Header: []string{"Column 1", "Column 2"},
			Rows: [][]string{
				{"Data 1", "Data 2"},
			},
		}
		options := TableOptions{}
		m.CustomTable(tableSet, options)

		got := m.String()

		// Check that it contains the expected table elements
		if !strings.Contains(got, "Column 1") || !strings.Contains(got, "Column 2") {
			t.Errorf("CustomTable should contain headers")
		}
		if !strings.Contains(got, "Data 1") || !strings.Contains(got, "Data 2") {
			t.Errorf("CustomTable should contain row data")
		}
	})

	t.Run("ValidateColumns with valid data", func(t *testing.T) {
		t.Parallel()

		tableSet := TableSet{
			Header: []string{"A", "B"},
			Rows: [][]string{
				{"1", "2"},
				{"3", "4"},
			},
		}

		err := tableSet.ValidateColumns()
		if err != nil {
			t.Errorf("ValidateColumns should not return error for valid data: %v", err)
		}
	})

	t.Run("ValidateColumns with invalid data", func(t *testing.T) {
		t.Parallel()

		tableSet := TableSet{
			Header: []string{"A", "B"},
			Rows: [][]string{
				{"1", "2", "3"}, // Too many columns
			},
		}

		err := tableSet.ValidateColumns()
		if err == nil {
			t.Error("ValidateColumns should return error for invalid data")
		}
	})
}

// Tests for edge cases and error conditions to improve coverage
func TestEdgeCasesAndErrors(t *testing.T) {
	t.Parallel()

	t.Run("Build with existing error", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		m := NewMarkdown(&buf)
		m.TableOfContents(TableOfContentsDepthH2)
		m.TableOfContents(TableOfContentsDepthH3) // This should cause an error

		err := m.Build()
		if err == nil {
			t.Error("Build() should return error when there's an existing error")
		}
		if !strings.Contains(err.Error(), "table of contents has already been generated") {
			t.Errorf("Error should mention table of contents duplication, got: %v", err)
		}
	})

	t.Run("Error method returns error when present", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.TableOfContents(TableOfContentsDepthH2)
		m.TableOfContents(TableOfContentsDepthH3) // This should cause an error

		err := m.Error()
		if err == nil {
			t.Error("Error() should return error when error is present")
		}
	})

	t.Run("NewMarkdown with different writer", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		m := NewMarkdown(&buf)
		m.H1("Test")

		if m.dest != &buf {
			t.Error("NewMarkdown should set the writer correctly")
		}

		err := m.Build()
		if err != nil {
			t.Errorf("Build should not return error: %v", err)
		}

		want := "# Test"
		got := buf.String()
		if got != want {
			t.Errorf("Expected %q, got %q", want, got)
		}
	})

	t.Run("String method with table of contents replacement", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Title")
		m.TableOfContents(TableOfContentsDepthH1)
		m.H2("Section")

		result := m.String()

		if !strings.Contains(result, "<!-- BEGIN_TOC -->") {
			t.Error("String should contain TOC begin marker")
		}
		if !strings.Contains(result, "<!-- END_TOC -->") {
			t.Error("String should contain TOC end marker")
		}
		if !strings.Contains(result, "- [Title](#title)") {
			t.Error("String should contain TOC content")
		}
	})

	t.Run("BulletList with empty slice", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.BulletList()

		want := ""
		got := m.String()
		if got != want {
			t.Errorf("BulletList with no items should produce empty string, got: %q", got)
		}
	})

	t.Run("OrderedList with single item", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.OrderedList("Only item")

		want := "1. Only item"
		got := m.String()
		if got != want {
			t.Errorf("OrderedList with single item: want %q, got %q", want, got)
		}
	})

	t.Run("CheckBox with empty slice", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.CheckBox([]CheckBoxSet{})

		want := ""
		got := m.String()
		if got != want {
			t.Errorf("CheckBox with empty slice should produce empty string, got: %q", got)
		}
	})

	t.Run("Blockquote with multiline text", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		multilineText := "Line 1" + internal.LineFeed() + "Line 2"
		m.Blockquote(multilineText)

		want := "> Line 1" + internal.LineFeed() + "> Line 2"
		got := m.String()
		if got != want {
			t.Errorf("Blockquote with multiline: want %q, got %q", want, got)
		}
	})

	t.Run("CodeBlocks with different syntax highlighting", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.CodeBlocks(SyntaxHighlightPython, "print('hello')")

		want := "```python" + internal.LineFeed() + "print('hello')" + internal.LineFeed() + "```"
		got := m.String()
		if got != want {
			t.Errorf("CodeBlocks with Python: want %q, got %q", want, got)
		}
	})
}

// Test table validation edge cases
func TestTableValidationEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("ValidateColumns with empty header", func(t *testing.T) {
		t.Parallel()

		tableSet := TableSet{
			Header: []string{},
			Rows:   [][]string{{"data"}},
		}

		err := tableSet.ValidateColumns()
		if err == nil {
			t.Error("ValidateColumns should return error for empty header with data rows")
		}
	})

	t.Run("ValidateColumns with empty rows", func(t *testing.T) {
		t.Parallel()

		tableSet := TableSet{
			Header: []string{"Col1", "Col2"},
			Rows:   [][]string{},
		}

		err := tableSet.ValidateColumns()
		if err != nil {
			t.Errorf("ValidateColumns should not return error for empty rows: %v", err)
		}
	})

	t.Run("ValidateColumns with row having fewer columns", func(t *testing.T) {
		t.Parallel()

		tableSet := TableSet{
			Header: []string{"Col1", "Col2", "Col3"},
			Rows: [][]string{
				{"data1", "data2"}, // Missing one column
			},
		}

		err := tableSet.ValidateColumns()
		if err == nil {
			t.Error("ValidateColumns should return error for row with fewer columns")
		}
	})
}

// Test table generation with various configurations
func TestTableGeneration(t *testing.T) {
	t.Parallel()

	t.Run("Table with alignment", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		tableSet := TableSet{
			Header: []string{"Left", "Center", "Right"},
			Rows: [][]string{
				{"L1", "C1", "R1"},
				{"L2", "C2", "R2"},
			},
		}
		m.Table(tableSet)

		got := m.String()

		// Should contain table structure
		if !strings.Contains(got, "Left") || !strings.Contains(got, "Center") || !strings.Contains(got, "Right") {
			t.Error("Table should contain all headers")
		}
		if !strings.Contains(got, "L1") || !strings.Contains(got, "C1") || !strings.Contains(got, "R1") {
			t.Error("Table should contain first row data")
		}
	})

	t.Run("CustomTable with options", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		tableSet := TableSet{
			Header: []string{"Name", "Value"},
			Rows: [][]string{
				{"Test", "123"},
			},
		}
		options := TableOptions{
			AutoWrapText: true,
		}
		m.CustomTable(tableSet, options)

		got := m.String()

		if !strings.Contains(got, "Name") || !strings.Contains(got, "Value") {
			t.Error("CustomTable should contain headers")
		}
		if !strings.Contains(got, "Test") || !strings.Contains(got, "123") {
			t.Error("CustomTable should contain row data")
		}
	})
}

// Test error handling with mock writer that fails
type failingWriter struct{}

func (fw *failingWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write failed")
}

func TestBuildWithWriteError(t *testing.T) {
	t.Parallel()

	t.Run("Build with write error", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(&failingWriter{})
		m.H1("Test")

		err := m.Build()
		if err == nil {
			t.Error("Build() should return error when write fails")
		}
		if !strings.Contains(err.Error(), "failed to write markdown text") {
			t.Errorf("Error should mention write failure, got: %v", err)
		}
	})

	t.Run("Build with write error and existing markdown error", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(&failingWriter{})
		m.TableOfContents(TableOfContentsDepthH2)
		m.TableOfContents(TableOfContentsDepthH3) // This creates an existing error
		m.H1("Test")

		err := m.Build()
		if err == nil {
			t.Error("Build() should return error when write fails")
		}
		// Should contain both the write error and the existing error
		errMsg := err.Error()
		if !strings.Contains(errMsg, "failed to write markdown text") {
			t.Errorf("Error should mention write failure, got: %v", err)
		}
		if !strings.Contains(errMsg, "table of contents has already been generated") {
			t.Errorf("Error should mention existing error, got: %v", err)
		}
	})
}
