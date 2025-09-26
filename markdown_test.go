// Package markdown is markdown builder that includes to convert Markdown to HTML.
package markdown

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
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

func TestTableOfContents(t *testing.T) {
	t.Parallel()

	t.Run("TOC with depth H2", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Introduction").
			H2("Overview").
			H3("Details").
			H2("Conclusion").
			TableOfContents(TableOfContentsDepthH2)

		want := "# Introduction" + internal.LineFeed() +
			"## Overview" + internal.LineFeed() +
			"### Details" + internal.LineFeed() +
			"## Conclusion" + internal.LineFeed() +
			"- [Introduction](#introduction)" + internal.LineFeed() +
			"  - [Overview](#overview)" + internal.LineFeed() +
			"  - [Conclusion](#conclusion)" + internal.LineFeed() +
			""
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("TOC with depth H3", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Chapter 1").
			H2("Section 1.1").
			H3("Subsection 1.1.1").
			H4("Detail").
			TableOfContents(TableOfContentsDepthH3)

		want := "# Chapter 1" + internal.LineFeed() +
			"## Section 1.1" + internal.LineFeed() +
			"### Subsection 1.1.1" + internal.LineFeed() +
			"#### Detail" + internal.LineFeed() +
			"- [Chapter 1](#chapter-1)" + internal.LineFeed() +
			"  - [Section 1.1](#section-11)" + internal.LineFeed() +
			"    - [Subsection 1.1.1](#subsection-111)" + internal.LineFeed() +
			""
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("TOC with depth H1 only", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("First").
			H2("Should not appear").
			H1("Second").
			TableOfContents(TableOfContentsDepthH1)

		want := "# First" + internal.LineFeed() +
			"## Should not appear" + internal.LineFeed() +
			"# Second" + internal.LineFeed() +
			"- [First](#first)" + internal.LineFeed() +
			"- [Second](#second)" + internal.LineFeed() +
			""
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("TOC with no headers", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.PlainText("Just text").
			TableOfContents(TableOfContentsDepthH2)

		want := "Just text"
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("TOC with special characters in headers", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Hello World!").
			H2("C++ Programming").
			H3("API & SDK").
			TableOfContents(TableOfContentsDepthH3)

		want := "# Hello World!" + internal.LineFeed() +
			"## C++ Programming" + internal.LineFeed() +
			"### API & SDK" + internal.LineFeed() +
			"- [Hello World!](#hello-world)" + internal.LineFeed() +
			"  - [C++ Programming](#c-programming)" + internal.LineFeed() +
			"    - [API & SDK](#api--sdk)" + internal.LineFeed() +
			""
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("TOC with depth H4", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Title").
			H2("Section").
			H3("Subsection").
			H4("Subsubsection").
			H5("Deep").
			TableOfContents(TableOfContentsDepthH4)

		want := "# Title" + internal.LineFeed() +
			"## Section" + internal.LineFeed() +
			"### Subsection" + internal.LineFeed() +
			"#### Subsubsection" + internal.LineFeed() +
			"##### Deep" + internal.LineFeed() +
			"- [Title](#title)" + internal.LineFeed() +
			"  - [Section](#section)" + internal.LineFeed() +
			"    - [Subsection](#subsection)" + internal.LineFeed() +
			"      - [Subsubsection](#subsubsection)" + internal.LineFeed() +
			""
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("TOC with depth H5", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("A").
			H2("B").
			H3("C").
			H4("D").
			H5("E").
			H6("F").
			TableOfContents(TableOfContentsDepthH5)

		want := "# A" + internal.LineFeed() +
			"## B" + internal.LineFeed() +
			"### C" + internal.LineFeed() +
			"#### D" + internal.LineFeed() +
			"##### E" + internal.LineFeed() +
			"###### F" + internal.LineFeed() +
			"- [A](#a)" + internal.LineFeed() +
			"  - [B](#b)" + internal.LineFeed() +
			"    - [C](#c)" + internal.LineFeed() +
			"      - [D](#d)" + internal.LineFeed() +
			"        - [E](#e)" + internal.LineFeed() +
			""
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("TOC with depth H6", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1("Level1").
			H2("Level2").
			H3("Level3").
			H4("Level4").
			H5("Level5").
			H6("Level6").
			TableOfContents(TableOfContentsDepthH6)

		want := "# Level1" + internal.LineFeed() +
			"## Level2" + internal.LineFeed() +
			"### Level3" + internal.LineFeed() +
			"#### Level4" + internal.LineFeed() +
			"##### Level5" + internal.LineFeed() +
			"###### Level6" + internal.LineFeed() +
			"- [Level1](#level1)" + internal.LineFeed() +
			"  - [Level2](#level2)" + internal.LineFeed() +
			"    - [Level3](#level3)" + internal.LineFeed() +
			"      - [Level4](#level4)" + internal.LineFeed() +
			"        - [Level5](#level5)" + internal.LineFeed() +
			"          - [Level6](#level6)" + internal.LineFeed() +
			""
		got := m.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("TOC with H1f, H2f format methods", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(os.Stdout)
		m.H1f("Title %d", 1).
			H2f("Section %s", "A").
			H3f("Sub %s", "B").
			TableOfContents(TableOfContentsDepthH2)

		want := "# Title 1" + internal.LineFeed() +
			"## Section A" + internal.LineFeed() +
			"### Sub B" + internal.LineFeed() +
			"- [Title 1](#title-1)" + internal.LineFeed() +
			"  - [Section A](#section-a)" + internal.LineFeed() +
			""
		got := m.String()

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
