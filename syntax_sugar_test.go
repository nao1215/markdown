package markdown

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal"
)

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
