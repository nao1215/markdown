package markdown

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unicode"

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
		want := fmt.Sprintf("<details>%s<summary>Hello</summary>%s%sGood World%s%s</details>%s",
			internal.LineFeed(), internal.LineFeed(), internal.LineFeed(),
			internal.LineFeed(), internal.LineFeed(), internal.LineFeed())
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
		// One entry per quote, not per line: the quote is a single block.
		want := []string{
			"> Hello" + internal.LineFeed() + "> Good" + internal.LineFeed() + "> World",
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

	t.Run("Table of contents handles unicode and duplicate headers", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("日本語タイトル").
			H2("同じ").
			H2("同じ").
			TableOfContents(TableOfContentsDepthH2)

		want := "# 日本語タイトル" + internal.LineFeed() +
			"## 同じ" + internal.LineFeed() +
			"## 同じ" + internal.LineFeed() +
			TableOfContentsMarkerBegin + internal.LineFeed() +
			"- [日本語タイトル](#日本語タイトル)" + internal.LineFeed() +
			"  - [同じ](#同じ)" + internal.LineFeed() +
			"  - [同じ](#同じ-1)" + internal.LineFeed() +
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

		want := "<details>" + internal.LineFeed() + "<summary>Summary</summary>" + internal.LineFeed() + internal.LineFeed() +
			"Hidden content" + internal.LineFeed() + internal.LineFeed() + "</details>" + internal.LineFeed()
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Detailsf method", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.Detailsf("Summary", "Hidden %s %d", "content", 42)

		want := "<details>" + internal.LineFeed() + "<summary>Summary</summary>" + internal.LineFeed() + internal.LineFeed() +
			"Hidden content 42" + internal.LineFeed() + internal.LineFeed() + "</details>" + internal.LineFeed()
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

		want := "# Test" + internal.LineFeed()
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

		want := "# Test" + internal.LineFeed()
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

	t.Run("Build with nil writer", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(nil)
		m.H1("Test")

		err := m.Build()
		if err == nil {
			t.Error("Build() should return error when writer is nil")
		}
		if !strings.Contains(err.Error(), "destination writer is nil") {
			t.Errorf("Error should mention nil writer, got: %v", err)
		}
	})
}

func TestMarkdownAlerts(t *testing.T) {
	t.Parallel()

	t.Run("success Notef()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.Notef("%s", "Hello")
		want := []string{"> [!NOTE]  " + internal.LineFeed() + "> Hello"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Warningf()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.Warningf("%s", "Hello")
		want := []string{"> [!WARNING]  " + internal.LineFeed() + "> Hello"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Tipf()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.Tipf("%s", "Hello")
		want := []string{"> [!TIP]  " + internal.LineFeed() + "> Hello"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Importantf()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.Importantf("%s", "Hello")
		want := []string{"> [!IMPORTANT]  " + internal.LineFeed() + "> Hello"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Cautionf()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.Cautionf("%s", "Hello")
		want := []string{"> [!CAUTION]  " + internal.LineFeed() + "> Hello"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestMarkdown_RedBadgef(t *testing.T) {
	t.Parallel()
	t.Run("success RedBadgef()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.RedBadgef("%s", "Hello")
		want := []string{"![Badge](https://img.shields.io/badge/Hello-red)"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success YellowBadgef()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.YellowBadgef("%s", "Hello")
		want := []string{"![Badge](https://img.shields.io/badge/Hello-yellow)"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success GreenBadgef()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.GreenBadgef("%s", "Hello")
		want := []string{"![Badge](https://img.shields.io/badge/Hello-green)"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success BlueBadgef()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.BlueBadgef("%s", "Hello")
		want := []string{"![Badge](https://img.shields.io/badge/Hello-blue)"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

// TestNormalizeLineFeeds covers the conversion applied to text that reaches the
// builder from elsewhere, such as a table rendered by tablewriter.
func TestNormalizeLineFeeds(t *testing.T) {
	t.Parallel()

	lineFeed := internal.LineFeed()

	tests := map[string]struct {
		in   string
		want string
	}{
		"unix input":       {in: "a\nb", want: "a" + lineFeed + "b"},
		"windows input":    {in: "a\r\nb", want: "a" + lineFeed + "b"},
		"mixed input":      {in: "a\r\nb\nc", want: "a" + lineFeed + "b" + lineFeed + "c"},
		"no line feeds":    {in: "abc", want: "abc"},
		"empty":            {in: "", want: ""},
		"trailing newline": {in: "a\n", want: "a" + lineFeed},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := normalizeLineFeeds(tt.in); got != tt.want {
				t.Errorf("normalizeLineFeeds(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestDocumentUsesOneLineEnding builds a document that exercises every
// construct and asserts it contains no line ending other than the platform one.
//
// This is a Windows regression test. It is trivially true on Linux and macOS,
// where the platform line feed is "\n", but on Windows it catches any construct
// that emits a bare "\n" and leaves the document with mixed endings. That is
// exactly how CustomTable shipped alignment that silently did nothing there:
// tablewriter separates rows with "\n" on every platform.
func TestDocumentUsesOneLineEnding(t *testing.T) {
	t.Parallel()

	lineFeed := internal.LineFeed()

	var buf bytes.Buffer
	err := NewMarkdown(&buf).
		H1("Title").
		TableOfContents(TableOfContentsDepthH3).
		PlainText("paragraph").
		H2("Lists").
		BulletList("alpha", "beta").
		OrderedList("one", "two").
		CheckBox([]CheckBoxSet{{Text: "todo"}, {Checked: true, Text: "done"}}).
		H2("Quotes").
		Blockquote("quoted"+lineFeed+"over two lines").
		Note("note"+lineFeed+"over two lines").
		Warning("warning").
		H2("Blocks").
		CodeBlocks(SyntaxHighlightGo, "x := 1").
		Details("summary", "body").
		HorizontalRule().
		H2("Tables").
		Table(TableSet{
			Header:    []string{"name", "value"},
			Rows:      [][]string{{"a", "1"}},
			Alignment: []TableAlignment{AlignLeft, AlignRight},
		}).
		CustomTable(TableSet{
			Header:    []string{"name", "value"},
			Rows:      [][]string{{"a", "1"}},
			Alignment: []TableAlignment{AlignLeft, AlignRight},
		}, TableOptions{}).
		LF().
		Build()
	if err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	if lineFeed == "\n" {
		if strings.Contains(buf.String(), "\r") {
			t.Errorf("document contains a carriage return:\n%q", buf.String())
		}
		return
	}

	// On Windows every "\n" has to be preceded by "\r".
	stripped := strings.ReplaceAll(buf.String(), lineFeed, "")
	if strings.ContainsAny(stripped, "\r\n") {
		t.Errorf("document mixes line endings:\n%q", buf.String())
	}
}

// TestCustomTableAlignmentSurvivesLineFeedConversion pins the specific defect:
// the delimiter row rewrite splits the rendered table on the platform line
// feed, so it found nothing to rewrite while tablewriter was emitting "\n".
func TestCustomTableAlignmentSurvivesLineFeedConversion(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := NewMarkdown(&buf).CustomTable(TableSet{
		Header:    []string{"name", "value"},
		Rows:      [][]string{{"a", "1"}},
		Alignment: []TableAlignment{AlignLeft, AlignRight},
	}, TableOptions{}).Build(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	lines := strings.Split(buf.String(), internal.LineFeed())
	if len(lines) < 2 {
		t.Fatalf("table was not split into rows: %q", buf.String())
	}
	if !strings.HasPrefix(lines[1], "|:") || !strings.HasSuffix(strings.TrimSuffix(lines[1], "|"), ":") {
		t.Errorf("alignment markers missing from the delimiter row: %q", lines[1])
	}
}

// TestNormalizeLineFeedsToEitherPlatform covers both targets on whichever
// platform the tests run on, rather than only the one the platform provides.
func TestNormalizeLineFeedsToEitherPlatform(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in       string
		lineFeed string
		want     string
	}{
		"unix input to unix":       {in: "a\nb", lineFeed: "\n", want: "a\nb"},
		"windows input to unix":    {in: "a\r\nb", lineFeed: "\n", want: "a\nb"},
		"mixed input to unix":      {in: "a\r\nb\nc", lineFeed: "\n", want: "a\nb\nc"},
		"unix input to windows":    {in: "a\nb", lineFeed: "\r\n", want: "a\r\nb"},
		"windows input to windows": {in: "a\r\nb", lineFeed: "\r\n", want: "a\r\nb"},
		"mixed input to windows":   {in: "a\r\nb\nc", lineFeed: "\r\n", want: "a\r\nb\r\nc"},
		"no line feeds":            {in: "abc", lineFeed: "\r\n", want: "abc"},
		"empty":                    {in: "", lineFeed: "\r\n", want: ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := normalizeLineFeedsTo(tt.in, tt.lineFeed); got != tt.want {
				t.Errorf("normalizeLineFeedsTo(%q, %q) = %q, want %q", tt.in, tt.lineFeed, got, tt.want)
			}
		})
	}
}

// lf is the platform line feed, spelled short because these tests are mostly
// exact-output comparisons.
func lf() string { return internal.LineFeed() }

// TestAlertSingleLineIsUnchanged is the compatibility guard for the alert
// rewrite: single-line text, which is what nearly every caller passes, has to
// produce exactly what it produced before.
func TestAlertSingleLineIsUnchanged(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		call func(*Markdown) *Markdown
		want string
	}{
		"note":      {func(m *Markdown) *Markdown { return m.Note("hello") }, "> [!NOTE]  " + lf() + "> hello"},
		"tip":       {func(m *Markdown) *Markdown { return m.Tip("hello") }, "> [!TIP]  " + lf() + "> hello"},
		"important": {func(m *Markdown) *Markdown { return m.Important("hello") }, "> [!IMPORTANT]  " + lf() + "> hello"},
		"warning":   {func(m *Markdown) *Markdown { return m.Warning("hello") }, "> [!WARNING]  " + lf() + "> hello"},
		"caution":   {func(m *Markdown) *Markdown { return m.Caution("hello") }, "> [!CAUTION]  " + lf() + "> hello"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := tt.call(NewMarkdown(nil)).String(); got != tt.want {
				t.Errorf("output changed:\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}

// TestAlertQuotesEveryLine covers the defect: only the first line carried the
// quote marker, so a list or a blank line in the text escaped the callout.
func TestAlertQuotesEveryLine(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		text string
		want []string
	}{
		"prose over two lines": {
			text: "first\nsecond",
			want: []string{"> [!NOTE]  ", "> first", "> second"},
		},
		"list in the body": {
			text: "intro\n- alpha\n- beta",
			want: []string{"> [!NOTE]  ", "> intro", "> - alpha", "> - beta"},
		},
		"blank line stays inside the callout": {
			text: "first\n\nsecond",
			want: []string{"> [!NOTE]  ", "> first", ">", "> second"},
		},
		"windows line endings in the text": {
			text: "first\r\nsecond",
			want: []string{"> [!NOTE]  ", "> first", "> second"},
		},
		"lines the caller already quoted are left alone": {
			// Callers hand-wrote these continuations to work around the bug.
			// Prefixing them again would nest a quote inside the alert.
			text: "first  \n> second  \n> * third",
			want: []string{"> [!NOTE]  ", "> first  ", "> second  ", "> * third"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := strings.Split(NewMarkdown(nil).Note(tt.text).String(), lf())
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("alert body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestBlockquoteSplitsOnPlainNewline covers the Windows case: splitting on the
// platform line feed meant a plain Go literal was never split at all.
func TestBlockquoteSplitsOnPlainNewline(t *testing.T) {
	t.Parallel()

	for name, text := range map[string]string{
		"unix":    "first\nsecond",
		"windows": "first\r\nsecond",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			want := "> first" + lf() + "> second"
			if got := NewMarkdown(nil).Blockquote(text).String(); got != want {
				t.Errorf("blockquote mismatch:\n got: %q\nwant: %q", got, want)
			}
		})
	}
}

// TestDetailsSurroundsBodyWithBlankLines covers the reason nested markdown used
// to render as literal text: an HTML block runs to the next blank line.
func TestDetailsSurroundsBodyWithBlankLines(t *testing.T) {
	t.Parallel()

	// The body is written by the caller and is passed through untouched, so it
	// has to be built with the platform line feed for this comparison.
	body := "- alpha" + lf() + "- beta"
	got := NewMarkdown(nil).Details("summary", body).String()
	want := strings.Join([]string{
		"<details>",
		"<summary>summary</summary>",
		"",
		"- alpha",
		"- beta",
		"",
		"</details>",
		"",
	}, lf())

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("details mismatch (-want +got):\n%s", diff)
	}
}

// TestDetailsDoesNotSwallowTheNextBlock is the other half of the same defect:
// without the trailing blank line the block after </details> stayed inside the
// HTML block and rendered as text.
func TestDetailsDoesNotSwallowTheNextBlock(t *testing.T) {
	t.Parallel()

	got := NewMarkdown(nil).
		Details("summary", "body").
		CodeBlocks(SyntaxHighlightGo, "x := 1").
		String()

	if !strings.Contains(got, "</details>"+lf()+lf()+"```go") {
		t.Errorf("no blank line between </details> and the next block:\n%q", got)
	}
}

// TestBuildTerminatesWithNewline pins MD047 and the append case: writing a
// second document to the same writer used to splice it onto the last line.
func TestBuildTerminatesWithNewline(t *testing.T) {
	t.Parallel()

	t.Run("single document", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		if err := NewMarkdown(&buf).H1("Title").Build(); err != nil {
			t.Fatalf("failed to build: %v", err)
		}
		if want := "# Title" + lf(); buf.String() != want {
			t.Errorf("got %q, want %q", buf.String(), want)
		}
	})

	t.Run("content already ending in a line feed gets no second one", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		if err := NewMarkdown(&buf).Table(TableSet{
			Header: []string{"a"},
			Rows:   [][]string{{"1"}},
		}).Build(); err != nil {
			t.Fatalf("failed to build: %v", err)
		}
		if strings.HasSuffix(buf.String(), lf()+lf()+lf()) {
			t.Errorf("trailing newlines piled up: %q", buf.String())
		}
	})

	t.Run("two documents to one writer stay on separate lines", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		if err := NewMarkdown(&buf).H1("First").Build(); err != nil {
			t.Fatalf("failed to build: %v", err)
		}
		if err := NewMarkdown(&buf).H1("Second").Build(); err != nil {
			t.Fatalf("failed to build: %v", err)
		}
		want := "# First" + lf() + "# Second" + lf()
		if buf.String() != want {
			t.Errorf("got %q, want %q", buf.String(), want)
		}
	})

	t.Run("String is left alone", func(t *testing.T) {
		t.Parallel()

		if got := NewMarkdown(nil).H1("Title").String(); got != "# Title" {
			t.Errorf("String() gained a line feed: %q", got)
		}
	})
}

// TestBlankLineBetweenBlocks covers the separations markdown requires. The
// cases that gain a blank line render wrongly without one; the rest must stay
// byte for byte as they were.
func TestBlankLineBetweenBlocks(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(*Markdown) *Markdown
		want  []string
	}{
		"list then table": {
			build: func(m *Markdown) *Markdown {
				return m.BulletList("alpha").Table(TableSet{Header: []string{"a"}, Rows: [][]string{{"1"}}})
			},
			want: []string{"- alpha", "", "| a |", "|---------|", "| 1 |", ""},
		},
		"list then paragraph": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha").PlainText("after") },
			want:  []string{"- alpha", "", "after"},
		},
		"list then heading": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha").H2("Next") },
			want:  []string{"- alpha", "", "## Next"},
		},
		"consecutive list items stay tight": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha", "beta") },
			want:  []string{"- alpha", "- beta"},
		},
		"ordered list items stay tight": {
			build: func(m *Markdown) *Markdown { return m.OrderedList("alpha", "beta") },
			want:  []string{"1. alpha", "2. beta"},
		},
		"a different list kind is separated": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha").OrderedList("one") },
			want:  []string{"- alpha", "", "1. one"},
		},
		"checkbox items stay tight": {
			build: func(m *Markdown) *Markdown {
				return m.CheckBox([]CheckBoxSet{{Text: "alpha"}, {Checked: true, Text: "beta"}})
			},
			want: []string{"- [ ] alpha", "- [x] beta"},
		},
		"alert then paragraph": {
			build: func(m *Markdown) *Markdown { return m.Note("careful").PlainText("after") },
			want:  []string{"> [!NOTE]  ", "> careful", "", "after"},
		},
		"two alerts do not merge": {
			build: func(m *Markdown) *Markdown { return m.Note("first").Warning("second") },
			want:  []string{"> [!NOTE]  ", "> first", "", "> [!WARNING]  ", "> second"},
		},
		"blockquote then paragraph": {
			build: func(m *Markdown) *Markdown { return m.Blockquote("quoted").PlainText("after") },
			want:  []string{"> quoted", "", "after"},
		},
		"an existing LF spacer is not doubled": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha").LF().PlainText("after") },
			want:  []string{"- alpha", "  ", "after"},
		},
		"an existing empty paragraph is not doubled": {
			build: func(m *Markdown) *Markdown { return m.BulletList("alpha").PlainText("").PlainText("after") },
			want:  []string{"- alpha", "", "after"},
		},
		"headings and paragraphs are untouched": {
			build: func(m *Markdown) *Markdown { return m.H1("Title").PlainText("text").H2("Section") },
			want:  []string{"# Title", "text", "## Section"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := strings.Split(tt.build(NewMarkdown(nil)).String(), lf())
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("block layout mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestTableOfContentsMarkersStayTight makes sure the block separation does not
// push a blank line into the marker pair, and that the entries still land
// between the markers. Matching the marker pair as one string used to be the
// mechanism, and it broke silently the moment anything separated them.
func TestTableOfContentsMarkersStayTight(t *testing.T) {
	t.Parallel()

	got := NewMarkdown(nil).
		H1("Doc").
		TableOfContents(TableOfContentsDepthH2).
		H2("Section").
		String()

	want := strings.Join([]string{
		"# Doc",
		TableOfContentsMarkerBegin,
		"- [Doc](#doc)",
		"  - [Section](#section)",
		TableOfContentsMarkerEnd,
		"",
		"## Section",
	}, lf())

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("table of contents mismatch (-want +got):\n%s", diff)
	}
}

// TestInsertTableOfContentsWithMissingMarkers guards the splice against a body
// that does not hold both markers.
func TestInsertTableOfContentsWithMissingMarkers(t *testing.T) {
	t.Parallel()

	toc := []string{"- [A](#a)"}

	tests := map[string][]string{
		"no markers at all": {"# Title"},
		"only the opening":  {TableOfContentsMarkerBegin, "# Title"},
		"only the closing":  {"# Title", TableOfContentsMarkerEnd},
		"reversed order":    {TableOfContentsMarkerEnd, TableOfContentsMarkerBegin},
	}

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if diff := cmp.Diff(body, insertTableOfContents(body, toc)); diff != "" {
				t.Errorf("body was modified (-want +got):\n%s", diff)
			}
		})
	}
}

// TestWithBlockSpacing covers the opt-in that separates every block, which is
// what markdownlint and mkdocs want. The default is deliberately tighter, so
// both shapes are pinned here.
func TestWithBlockSpacing(t *testing.T) {
	t.Parallel()

	chain := func(m *Markdown) *Markdown {
		return m.H1("Title").
			PlainText("paragraph").
			BulletList("alpha", "beta").
			CodeBlocks(SyntaxHighlightGo, "x := 1")
	}

	compact := chain(NewMarkdown(nil)).String()
	spaced := chain(NewMarkdown(nil, WithBlockSpacing())).String()

	wantCompact := strings.Join([]string{
		"# Title",
		"paragraph",
		"- alpha",
		"- beta",
		"",
		"```go",
		"x := 1",
		"```",
	}, lf())
	if diff := cmp.Diff(wantCompact, compact); diff != "" {
		t.Errorf("default output changed (-want +got):\n%s", diff)
	}

	wantSpaced := strings.Join([]string{
		"# Title",
		"",
		"paragraph",
		"",
		"- alpha",
		"- beta",
		"",
		"```go",
		"x := 1",
		"```",
	}, lf())
	if diff := cmp.Diff(wantSpaced, spaced); diff != "" {
		t.Errorf("spaced output mismatch (-want +got):\n%s", diff)
	}
}

// TestBlankLine covers the explicit spacer, and that it is not doubled by the
// automatic separation.
func TestBlankLine(t *testing.T) {
	t.Parallel()

	got := NewMarkdown(nil).BulletList("alpha").BlankLine().PlainText("after").String()
	want := strings.Join([]string{"- alpha", "", "after"}, lf())

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("blank line mismatch (-want +got):\n%s", diff)
	}
}

func TestLink(t *testing.T) {
	t.Parallel()

	t.Run("success Link()", func(t *testing.T) {
		t.Parallel()

		want := "[Hello](https://example.com)"
		got := Link("Hello", "https://example.com")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestImage(t *testing.T) {
	t.Parallel()

	t.Run("success Image()", func(t *testing.T) {
		t.Parallel()

		want := "![Hello](https://example.com/image.png)"
		got := Image("Hello", "https://example.com/image.png")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestFootnoteReference(t *testing.T) {
	t.Parallel()

	t.Run("success FootnoteReference()", func(t *testing.T) {
		t.Parallel()

		want := "[^1]"
		got := FootnoteReference("1")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestFootnoteDefinition(t *testing.T) {
	t.Parallel()

	t.Run("success FootnoteDefinition()", func(t *testing.T) {
		t.Parallel()

		want := "[^1]: This is footnote"
		got := FootnoteDefinition("1", "This is footnote")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestReferenceLink(t *testing.T) {
	t.Parallel()

	t.Run("success ReferenceLink()", func(t *testing.T) {
		t.Parallel()

		want := "[Go][go-site]"
		got := ReferenceLink("Go", "go-site")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestReferenceLinkDefinition(t *testing.T) {
	t.Parallel()

	t.Run("success ReferenceLinkDefinition() without title", func(t *testing.T) {
		t.Parallel()

		want := "[go-site]: https://golang.org"
		got := ReferenceLinkDefinition("go-site", "https://golang.org")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success ReferenceLinkDefinition() with title", func(t *testing.T) {
		t.Parallel()

		want := "[go-site]: https://golang.org \"The Go Programming Language\""
		got := ReferenceLinkDefinition("go-site", "https://golang.org", "The Go Programming Language")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("ReferenceLinkDefinition() with empty title keeps original format", func(t *testing.T) {
		t.Parallel()

		want := "[go-site]: https://golang.org"
		got := ReferenceLinkDefinition("go-site", "https://golang.org", "")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("ReferenceLinkDefinition() escapes title quotes", func(t *testing.T) {
		t.Parallel()

		want := "[go-site]: https://golang.org \"The \\\"Go\\\" Programming Language\""
		got := ReferenceLinkDefinition("go-site", "https://golang.org", "The \"Go\" Programming Language")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("ReferenceLinkDefinition() escapes title backslashes", func(t *testing.T) {
		t.Parallel()

		want := "[go-site]: https://golang.org \"foo\\\\\""
		got := ReferenceLinkDefinition("go-site", "https://golang.org", "foo\\")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("ReferenceLinkDefinition() escapes path style backslashes in title", func(t *testing.T) {
		t.Parallel()

		want := `[go-site]: https://golang.org "C:\\path\\"`
		got := ReferenceLinkDefinition("go-site", "https://golang.org", `C:\path\`)

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestInlineMath(t *testing.T) {
	t.Parallel()

	t.Run("success InlineMath()", func(t *testing.T) {
		t.Parallel()

		want := "$E=mc^2$"
		got := InlineMath("E=mc^2")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("InlineMath() escapes dollar signs", func(t *testing.T) {
		t.Parallel()

		want := "$price = \\$100$"
		got := InlineMath("price = $100")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestBlockMath(t *testing.T) {
	t.Parallel()

	t.Run("success BlockMath()", func(t *testing.T) {
		t.Parallel()

		want := "$$" + internal.LineFeed() + "x^2 + y^2 = z^2" + internal.LineFeed() + "$$"
		got := BlockMath("x^2 + y^2 = z^2")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("BlockMath() keeps dollar signs as is", func(t *testing.T) {
		t.Parallel()

		want := "$$" + internal.LineFeed() + "cost = $x + $y" + internal.LineFeed() + "$$"
		got := BlockMath("cost = $x + $y")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestStrikethrough(t *testing.T) {
	t.Parallel()

	t.Run("success Strikethrough()", func(t *testing.T) {
		t.Parallel()

		want := "~~Hello~~"
		got := Strikethrough("Hello")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestBold(t *testing.T) {
	t.Parallel()

	t.Run("success Bold()", func(t *testing.T) {
		t.Parallel()

		want := "**Hello**"
		got := Bold("Hello")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestItalic(t *testing.T) {
	t.Parallel()

	t.Run("success Italic()", func(t *testing.T) {
		t.Parallel()

		want := "*Hello*"
		got := Italic("Hello")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestBoldItalic(t *testing.T) {
	t.Parallel()

	t.Run("success BoldItalic()", func(t *testing.T) {
		t.Parallel()

		want := "***Hello***"
		got := BoldItalic("Hello")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestCode(t *testing.T) {
	t.Parallel()

	t.Run("success Code()", func(t *testing.T) {
		t.Parallel()

		want := "`Hello`"
		got := Code("Hello")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestHighlight(t *testing.T) {
	t.Parallel()

	t.Run("success Highlight()", func(t *testing.T) {
		t.Parallel()

		want := "==Hello=="
		got := Highlight("Hello")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

// delimiterRow returns the alignment row of the first table in the rendered
// document, which is where markdown records column alignment.
func delimiterRow(t *testing.T, rendered string) string {
	t.Helper()

	lines := strings.Split(rendered, internal.LineFeed())
	if len(lines) < 2 {
		t.Fatalf("rendered table has %d line(s), want at least 2:\n%s", len(lines), rendered)
	}
	return lines[1]
}

// TestCustomTableHonorsAlignment pins the behavior that CustomTable used to
// drop on the floor: TableSet.Alignment is a documented public field, but it
// never reached the rendered output, so choosing CustomTable silently cost you
// alignment.
func TestCustomTableHonorsAlignment(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		alignment []TableAlignment
		want      string
	}{
		"left, center, right": {
			alignment: []TableAlignment{AlignLeft, AlignCenter, AlignRight},
			want:      "|:----|:------:|-----:|",
		},
		"default keeps a plain rule": {
			alignment: []TableAlignment{AlignDefault, AlignDefault, AlignDefault},
			want:      "|-----|--------|------|",
		},
		"shorter than the header falls back to default": {
			alignment: []TableAlignment{AlignRight},
			want:      "|----:|--------|------|",
		},
		"longer than the header ignores the extra entries": {
			alignment: []TableAlignment{AlignLeft, AlignLeft, AlignLeft, AlignRight, AlignCenter},
			want:      "|:----|:-------|:-----|",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			if err := NewMarkdown(&buf).CustomTable(TableSet{
				Header:    []string{"key", "middle", "last"},
				Rows:      [][]string{{"a", "b", "c"}},
				Alignment: tt.alignment,
			}, TableOptions{}).Build(); err != nil {
				t.Fatalf("failed to build: %v", err)
			}

			if got := delimiterRow(t, buf.String()); got != tt.want {
				t.Errorf("delimiter row mismatch:\n got: %q\nwant: %q\nfull output:\n%s", got, tt.want, buf.String())
			}
		})
	}
}

// TestCustomTableWithoutAlignmentIsUnchanged is the compatibility guard: a
// TableSet with no Alignment must render exactly as it did before alignment
// support was wired in, so existing golden files keep passing.
func TestCustomTableWithoutAlignmentIsUnchanged(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := NewMarkdown(&buf).CustomTable(TableSet{
		Header: []string{"key", "middle", "last"},
		Rows:   [][]string{{"a", "b", "c"}, {"longer", "value", "here"}},
	}, TableOptions{}).Build(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	want := strings.Join([]string{
		"|  key   | middle | last |",
		"|--------|--------|------|",
		"| a      | b      | c    |",
		"| longer | value  | here |",
		"",
	}, internal.LineFeed())

	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("output changed (-want +got):\n%s", diff)
	}
}

// TestTableAndCustomTableAgreeOnAlignment keeps the two table renderers from
// drifting apart again: whatever alignment a caller asks for, both must place
// the colons on the same sides.
func TestTableAndCustomTableAgreeOnAlignment(t *testing.T) {
	t.Parallel()

	set := TableSet{
		Header:    []string{"key", "middle", "last"},
		Rows:      [][]string{{"a", "b", "c"}},
		Alignment: []TableAlignment{AlignRight, AlignCenter, AlignLeft},
	}

	var plain, custom bytes.Buffer
	if err := NewMarkdown(&plain).Table(set).Build(); err != nil {
		t.Fatalf("failed to build Table: %v", err)
	}
	if err := NewMarkdown(&custom).CustomTable(set, TableOptions{}).Build(); err != nil {
		t.Fatalf("failed to build CustomTable: %v", err)
	}

	sides := func(row string) []string {
		cells := strings.Split(strings.Trim(row, "|"), "|")
		out := make([]string, 0, len(cells))
		for _, cell := range cells {
			switch {
			case strings.HasPrefix(cell, ":") && strings.HasSuffix(cell, ":"):
				out = append(out, "center")
			case strings.HasPrefix(cell, ":"):
				out = append(out, "left")
			case strings.HasSuffix(cell, ":"):
				out = append(out, "right")
			default:
				out = append(out, "default")
			}
		}
		return out
	}

	if diff := cmp.Diff(sides(delimiterRow(t, plain.String())), sides(delimiterRow(t, custom.String()))); diff != "" {
		t.Errorf("Table and CustomTable disagree on alignment (-Table +CustomTable):\n%s", diff)
	}
}

// TestDelimiterCell covers the widths where a marker does not fit, which the
// table renderers can reach for a very narrow column.
func TestDelimiterCell(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		width int
		align TableAlignment
		want  string
	}{
		"left needs two characters":    {width: 1, align: AlignLeft, want: "-"},
		"left at the minimum":          {width: 2, align: AlignLeft, want: ":-"},
		"right needs two characters":   {width: 1, align: AlignRight, want: "-"},
		"right at the minimum":         {width: 2, align: AlignRight, want: "-:"},
		"center needs three":           {width: 2, align: AlignCenter, want: "--"},
		"center at the minimum":        {width: 3, align: AlignCenter, want: ":-:"},
		"default is always all dashes": {width: 4, align: AlignDefault, want: "----"},
		"zero width":                   {width: 0, align: AlignCenter, want: ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := delimiterCell(tt.width, tt.align); got != tt.want {
				t.Errorf("delimiterCell(%d, %v) = %q, want %q", tt.width, tt.align, got, tt.want)
			}
		})
	}
}

// TestApplyAlignmentToDelimiterRowLeavesOddInputAlone makes sure the rewrite
// never corrupts something that is not a delimiter row.
func TestApplyAlignmentToDelimiterRowLeavesOddInputAlone(t *testing.T) {
	t.Parallel()

	alignment := []TableAlignment{AlignCenter}

	tests := map[string]string{
		"no alignment requested":    "| a |" + internal.LineFeed() + "|---|",
		"single line":               "| a |",
		"second line is not a rule": "| a |" + internal.LineFeed() + "| b |",
		"empty":                     "",
	}

	for name, rendered := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			use := alignment
			if name == "no alignment requested" {
				use = nil
			}
			if got := applyAlignmentToDelimiterRow(rendered, use); got != rendered {
				t.Errorf("input was rewritten:\n got: %q\nwant: %q", got, rendered)
			}
		})
	}
}

// TestEscapeTableCell covers the characters that end a cell or a row, and the
// idempotence the function promises.
func TestEscapeTableCell(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   string
		want string
	}{
		"plain text is returned unchanged":   {in: "hello", want: "hello"},
		"empty":                              {in: "", want: ""},
		"pipe":                               {in: "a|b", want: `a\|b`},
		"several pipes":                      {in: "a|b|c", want: `a\|b\|c`},
		"pipe at the edges":                  {in: "|a|", want: `\|a\|`},
		"already escaped pipe is left alone": {in: `a\|b`, want: `a\|b`},
		"escaped backslash then pipe":        {in: `a\\|b`, want: `a\\\|b`},
		"newline becomes a break":            {in: "a\nb", want: "a<br>b"},
		"carriage return and newline":        {in: "a\r\nb", want: "a<br>b"},
		"lone carriage return":               {in: "a\rb", want: "a<br>b"},
		"trailing backslash":                 {in: `a\`, want: `a\`},
		"markup the caller built survives":   {in: "[Go](https://go.dev)", want: "[Go](https://go.dev)"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := EscapeTableCell(tt.in)
			if got != tt.want {
				t.Errorf("EscapeTableCell(%q) = %q, want %q", tt.in, got, tt.want)
			}
			if again := EscapeTableCell(got); again != got {
				t.Errorf("not idempotent: %q became %q", got, again)
			}
		})
	}
}

// TestTableEscapeCellsKeepsColumnCount is the point of the option: an unescaped
// pipe splits the cell in two, so GitHub renders the row with the wrong number
// of columns and drops the last value.
func TestTableEscapeCellsKeepsColumnCount(t *testing.T) {
	t.Parallel()

	for name, render := range map[string]func(*Markdown, TableSet) *Markdown{
		"Table":       func(m *Markdown, ts TableSet) *Markdown { return m.Table(ts) },
		"CustomTable": func(m *Markdown, ts TableSet) *Markdown { return m.CustomTable(ts, TableOptions{}) },
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			if err := render(NewMarkdown(&buf), TableSet{
				Header:      []string{"cmd", "desc"},
				Rows:        [][]string{{"a|b", "keeps its value"}},
				EscapeCells: true,
			}).Build(); err != nil {
				t.Fatalf("failed to build: %v", err)
			}

			lines := strings.Split(strings.TrimRight(buf.String(), internal.LineFeed()), internal.LineFeed())
			row := lines[len(lines)-1]
			if cells := unescapedPipeCount(row); cells != 3 {
				t.Errorf("row has %d unescaped pipes, want 3 (one per edge plus one separator): %q", cells, row)
			}
			if !strings.Contains(row, "keeps its value") {
				t.Errorf("the last column was dropped: %q", row)
			}
		})
	}
}

// unescapedPipeCount counts the pipes that actually delimit cells.
func unescapedPipeCount(row string) int {
	count := 0
	for i := 0; i < len(row); i++ {
		if row[i] == '\\' {
			i++
			continue
		}
		if row[i] == '|' {
			count++
		}
	}
	return count
}

// TestTableWithoutEscapeCellsIsUnchanged is the compatibility guard: the option
// is off by default and the rendering must not move.
func TestTableWithoutEscapeCellsIsUnchanged(t *testing.T) {
	t.Parallel()

	set := TableSet{
		Header: []string{"cmd", "desc"},
		Rows:   [][]string{{"a|b", "text"}},
	}

	var buf bytes.Buffer
	if err := NewMarkdown(&buf).Table(set).Build(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	want := strings.Join([]string{
		"| cmd | desc |",
		"|---------|---------|",
		"| a|b | text |",
		"",
	}, internal.LineFeed())

	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("output changed (-want +got):\n%s", diff)
	}
}

// TestEscapeCellsDoesNotMutateTheCallersSlices makes sure enabling the option
// does not write back into the data the caller passed in.
func TestEscapeCellsDoesNotMutateTheCallersSlices(t *testing.T) {
	t.Parallel()

	header := []string{"a|b"}
	rows := [][]string{{"c|d"}}

	var buf bytes.Buffer
	if err := NewMarkdown(&buf).Table(TableSet{
		Header:      header,
		Rows:        rows,
		EscapeCells: true,
	}).Build(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	if header[0] != "a|b" || rows[0][0] != "c|d" {
		t.Errorf("caller data was modified: header=%v rows=%v", header, rows)
	}
}

// tocEntries returns the lines between the table of contents markers.
func tocEntries(t *testing.T, rendered string) []string {
	t.Helper()

	begin := strings.Index(rendered, TableOfContentsMarkerBegin)
	end := strings.Index(rendered, TableOfContentsMarkerEnd)
	if begin == -1 || end == -1 {
		t.Fatalf("markers missing from output:\n%s", rendered)
	}

	body := rendered[begin+len(TableOfContentsMarkerBegin) : end]
	entries := []string{}
	for _, line := range strings.Split(body, internal.LineFeed()) {
		if strings.TrimSpace(line) != "" {
			entries = append(entries, line)
		}
	}
	return entries
}

// TestTableOfContentsIndentsFromTheShallowestHeading covers the case where the
// document has no H1. TableOfContents pins MinDepth at H1, so the indent used
// to be measured from a level that never appears and every entry came out
// indented under nothing.
func TestTableOfContentsIndentsFromTheShallowestHeading(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(*Markdown) *Markdown
		want  []string
	}{
		"document starting at H2": {
			build: func(m *Markdown) *Markdown {
				return m.TableOfContents(TableOfContentsDepthH3).H2("S1").H3("S1.1")
			},
			want: []string{"- [S1](#s1)", "  - [S1.1](#s11)"},
		},
		"document with an H1 is unchanged": {
			build: func(m *Markdown) *Markdown {
				return m.H1("Doc").TableOfContents(TableOfContentsDepthH3).H2("S1").H3("S1.1")
			},
			want: []string{"- [Doc](#doc)", "  - [S1](#s1)", "    - [S1.1](#s11)"},
		},
		"document starting at H3": {
			build: func(m *Markdown) *Markdown {
				return m.TableOfContents(TableOfContentsDepthH4).H3("Deep").H4("Deeper")
			},
			want: []string{"- [Deep](#deep)", "  - [Deeper](#deeper)"},
		},
		"shallowest heading appears late": {
			build: func(m *Markdown) *Markdown {
				return m.TableOfContents(TableOfContentsDepthH3).H3("First").H2("Later")
			},
			want: []string{"  - [First](#first)", "- [Later](#later)"},
		},
		"explicit range is unaffected": {
			build: func(m *Markdown) *Markdown {
				return m.TableOfContentsWithRange(TableOfContentsDepthH2, TableOfContentsDepthH4).
					H2("S1").H3("S1.1")
			},
			want: []string{"- [S1](#s1)", "  - [S1.1](#s11)"},
		},
		"headings outside the range are skipped": {
			build: func(m *Markdown) *Markdown {
				return m.TableOfContentsWithRange(TableOfContentsDepthH2, TableOfContentsDepthH2).
					H1("Title").H2("S1").H3("S1.1")
			},
			want: []string{"- [S1](#s1)"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			if err := tt.build(NewMarkdown(&buf)).Build(); err != nil {
				t.Fatalf("failed to build: %v", err)
			}

			if diff := cmp.Diff(tt.want, tocEntries(t, buf.String())); diff != "" {
				t.Errorf("table of contents mismatch (-want +got):\n%s\nfull output:\n%s", diff, buf.String())
			}
		})
	}
}

// TestTableOfContentsWithNoMatchingHeadings makes sure the indent baseline does
// not misbehave when nothing falls inside the requested range.
func TestTableOfContentsWithNoMatchingHeadings(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := NewMarkdown(&buf).
		TableOfContentsWithRange(TableOfContentsDepthH4, TableOfContentsDepthH6).
		H1("Title").
		H2("Section").
		Build(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	if entries := tocEntries(t, buf.String()); len(entries) != 0 {
		t.Errorf("expected no entries, got %v", entries)
	}
}

// Every defect found in the recent downstream survey was a string-handling
// defect: a pipe in a cell dropped a column, a newline in alert text escaped the
// callout, a heading with punctuation produced an anchor nobody checked. The
// table-driven tests only cover the inputs someone thought of, so these targets
// state the properties instead and let the fuzzer look for the counterexample.

// FuzzEscapeTableCell asserts the two properties the escaper promises: the
// result never contains a character that ends a cell or a row, and escaping
// twice changes nothing.
func FuzzEscapeTableCell(f *testing.F) {
	for _, seed := range []string{
		"", "plain", "a|b", "a|b|c", `a\|b`, `a\\|b`, "a\nb", "a\r\nb", "a\rb",
		`a\`, "[Go](https://go.dev)", "**bold**", "|", "||", "\n", "\r",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, in string) {
		got := EscapeTableCell(in)

		if strings.ContainsAny(got, "\r\n") {
			t.Errorf("EscapeTableCell(%q) = %q, which still ends the row", in, got)
		}
		if unescapedPipeCount(got) != 0 {
			t.Errorf("EscapeTableCell(%q) = %q, which still ends the cell", in, got)
		}
		if again := EscapeTableCell(got); again != got {
			t.Errorf("EscapeTableCell is not idempotent: %q became %q", got, again)
		}
	})
}

// FuzzTableRowKeepsColumnCount asserts the property the escaper exists for: for
// any cell content, an escaped row still has exactly as many cells as the
// header. This is the property the unescaped path violates.
func FuzzTableRowKeepsColumnCount(f *testing.F) {
	for _, seed := range []string{"a", "a|b", "x\ny", "|", `\|`, ""} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, cell string) {
		out := NewMarkdown(nil).Table(TableSet{
			Header:      []string{"one", "two"},
			Rows:        [][]string{{cell, "second"}},
			EscapeCells: true,
		}).String()

		lines := strings.Split(strings.TrimRight(out, internal.LineFeed()), internal.LineFeed())
		row := lines[len(lines)-1]

		// Two edges plus one separator between the cells.
		if got := unescapedPipeCount(row); got != 3 {
			t.Errorf("row %q built from cell %q has %d delimiters, want 3", row, cell, got)
		}
	})
}

// FuzzAlertQuotesEveryLine asserts that no line of an alert can escape the
// callout, whatever the caller puts in the text.
func FuzzAlertQuotesEveryLine(f *testing.F) {
	for _, seed := range []string{
		"", "one line", "first\nsecond", "intro\n- item", "a\n\nb", "a\r\nb",
		"> already quoted", "\n", "   ", "```\ncode\n```",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		out := NewMarkdown(nil).Note(text).String()

		for i, line := range strings.Split(out, internal.LineFeed()) {
			if !strings.HasPrefix(line, ">") {
				t.Errorf("line %d of the alert is outside the callout: %q\nfull output:\n%q", i, line, out)
			}
		}
	})
}

// FuzzGitHubAnchor asserts that a generated anchor holds only the characters
// GitHub keeps, and that generating it twice gives the same answer. A heading
// is arbitrary user text, so this is where punctuation and non-ASCII land.
func FuzzGitHubAnchor(f *testing.F) {
	for _, seed := range []string{
		"", "Simple", "With Space", "punctuation!?", "`code`", "日本語の見出し",
		"a--b", "  leading", "MiXeD CaSe", "1. numbered", "emoji 🎉",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, heading string) {
		anchor := generateGitHubAnchor(heading)

		for _, r := range anchor {
			if r == '-' || unicode.IsLetter(r) || unicode.IsNumber(r) {
				continue
			}
			t.Errorf("anchor %q from heading %q contains %q", anchor, heading, r)
		}
		if again := generateGitHubAnchor(heading); again != anchor {
			t.Errorf("anchor for %q is not stable: %q then %q", heading, anchor, again)
		}
	})
}

// FuzzTableOfContentsAnchorsMatchHeadings asserts that every link the table of
// contents emits points at an anchor derived from a heading in the document.
//
// Uniqueness is deliberately not asserted: GitHub resolves a duplicate anchor to
// the first heading that produces it, so a document containing both "a" twice
// and a literal "a-1" genuinely has a collision. Renaming our anchor would point
// the link at something GitHub does not have.
func FuzzTableOfContentsAnchorsMatchHeadings(f *testing.F) {
	for _, seed := range []string{"", "Same", "Same|Same", "A|B|A", "a-1|a|a", "`code`|日本語"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, joined string) {
		headings := strings.Split(joined, "|")
		if len(headings) > 32 { //nolint:mnd // keep the generated document small
			headings = headings[:32]
		}

		valid := map[string]bool{}
		m := NewMarkdown(nil).TableOfContents(TableOfContentsDepthH2)
		for _, heading := range headings {
			m = m.H2(heading)
			base := generateGitHubAnchor(heading)
			valid[base] = true
			for i := 1; i <= len(headings); i++ {
				valid[fmt.Sprintf("%s-%d", base, i)] = true
			}
		}

		// Scan only between the markers. A heading is arbitrary text and can
		// itself contain "](#", so scanning the whole document would pick up the
		// heading rather than the generated link and report a false failure.
		// The link is the last "](#" on the line for the same reason.
		entries := 0
		inTOC := false
		for _, line := range strings.Split(m.String(), internal.LineFeed()) {
			switch line {
			case TableOfContentsMarkerBegin:
				inTOC = true
				continue
			case TableOfContentsMarkerEnd:
				inTOC = false
				continue
			}
			if !inTOC {
				continue
			}

			open := strings.LastIndex(line, "](#")
			if open == -1 || !strings.HasSuffix(line, ")") {
				continue
			}
			entries++
			anchor := line[open+len("](#") : len(line)-1]
			if !valid[anchor] {
				t.Errorf("anchor %q does not come from any heading in %q", anchor, headings)
			}
		}

		if entries != len(headings) {
			t.Errorf("table of contents has %d entries for %d headings", entries, len(headings))
		}
	})
}
