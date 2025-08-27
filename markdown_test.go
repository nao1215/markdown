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
