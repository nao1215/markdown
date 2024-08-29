package markdown

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
