package markdown

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
}
