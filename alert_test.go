package markdown

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMarkdownAlerts(t *testing.T) {
	t.Parallel()

	t.Run("success Notef()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.Notef("%s", "Hello")
		want := []string{"> [!NOTE]  \n> Hello"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Warningf()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.Warningf("%s", "Hello")
		want := []string{"> [!WARNING]  \n> Hello"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Tipf()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.Tipf("%s", "Hello")
		want := []string{"> [!TIP]  \n> Hello"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Importantf()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.Importantf("%s", "Hello")
		want := []string{"> [!IMPORTANT]  \n> Hello"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("success Cautionf()", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdown(io.Discard)
		m.Cautionf("%s", "Hello")
		want := []string{"> [!CAUTION]  \n> Hello"}
		got := m.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}
