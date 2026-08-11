package markdown

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestGenerateIndexErrors covers the three failure paths, each of which wraps a
// different sentinel so the caller can tell what went wrong.
func TestGenerateIndexErrors(t *testing.T) {
	t.Parallel()

	t.Run("an option that fails stops the walk", func(t *testing.T) {
		t.Parallel()

		broken := func(_ *Index) error { return errors.New("bad option") }

		err := GenerateIndex(indexFixtureDir(), broken)
		if !errors.Is(err, ErrInitMarkdownIndex) {
			t.Errorf("got %v, want it to wrap ErrInitMarkdownIndex", err)
		}
	})

	t.Run("a missing directory is reported", func(t *testing.T) {
		t.Parallel()

		err := GenerateIndex(filepath.Join(indexFixtureDir(), "does-not-exist"), WithWriter(os.Stdout))
		if !errors.Is(err, ErrCreateMarkdownIndex) {
			t.Errorf("got %v, want it to wrap ErrCreateMarkdownIndex", err)
		}
	})

	t.Run("a writer that fails is reported", func(t *testing.T) {
		t.Parallel()

		err := GenerateIndex(indexFixtureDir(), WithWriter(&failingWriter{}))
		if !errors.Is(err, ErrWriteMarkdownIndex) {
			t.Errorf("got %v, want it to wrap ErrWriteMarkdownIndex", err)
		}
	})
}

// TestFirstH1orH2 covers the heading lookup that gives each index entry its
// label, including the inputs where there is nothing to find.
func TestFirstH1orH2(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	write := func(name, content string) string {
		t.Helper()
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("failed to write %s: %v", name, err)
		}
		return path
	}

	tests := map[string]struct {
		path string
		want string
	}{
		"h1 wins":                {write("h1.md", "# Title\n## Section\n"), "Title"},
		"h2 when there is no h1": {write("h2.md", "text\n## Section\n"), "Section"},
		"no heading":             {write("none.md", "just text\n"), ""},
		"empty file":             {write("empty.md", ""), ""},
		"missing file":           {filepath.Join(dir, "absent.md"), ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := firstH1orH2(tt.path); got != tt.want {
				t.Errorf("firstH1orH2(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

// TestCustomTableRejectsMismatchedRows covers the validation path: a row with
// the wrong number of cells has to surface as an error rather than a table that
// renders with a column missing.
func TestCustomTableRejectsMismatchedRows(t *testing.T) {
	t.Parallel()

	m := NewMarkdown(nil).CustomTable(TableSet{
		Header: []string{"a", "b"},
		Rows:   [][]string{{"1"}},
	}, TableOptions{})

	if m.Error() == nil {
		t.Fatal("a row with too few cells must be reported")
	}
	if !errors.Is(m.Error(), ErrMismatchColumn) {
		t.Errorf("got %v, want it to wrap ErrMismatchColumn", m.Error())
	}
}

// TestCustomTableOptions covers the two formatting switches, which change the
// rendered table and had no assertion on their effect.
func TestCustomTableOptions(t *testing.T) {
	t.Parallel()

	set := TableSet{
		Header: []string{"name", "description"},
		Rows:   [][]string{{"a", "a description long enough to wrap somewhere"}},
	}

	plain := NewMarkdown(nil).CustomTable(set, TableOptions{}).String()
	formatted := NewMarkdown(nil).CustomTable(set, TableOptions{
		AutoWrapText:      true,
		AutoFormatHeaders: true,
	}).String()

	if plain == formatted {
		t.Errorf("the options changed nothing:\n%s", plain)
	}
}
