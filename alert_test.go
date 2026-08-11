package markdown

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal"
)

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
