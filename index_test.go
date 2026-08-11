package markdown

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// indexFixtureDir returns the fixture tree GenerateIndex walks in these tests.
//
// It is a directory of its own rather than the whole of testdata, because the
// generated index lists every markdown file it finds: sharing the tree with the
// golden documents of the other tests would make every new golden file change
// the expected index.
//
// The path is joined rather than written as "testdata/index" because
// GenerateIndex strips the target directory from each walked path using the
// platform separator, so a forward slash here produces the wrong links on
// Windows. See the issue linked from the pull request that added this.
func indexFixtureDir() string {
	return filepath.Join("testdata", "index")
}

func TestGenerateIndex(t *testing.T) {
	t.Parallel()

	t.Run("create index", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		if err := GenerateIndex(
			indexFixtureDir(),
			WithTitle("Test Title"),
			WithDescription([]string{"Test Description", "Next Description"}),
			WithWriter(&buf),
		); err != nil {
			t.Fatalf("failed to generate index: %v", err)
		}

		f := filepath.Join(indexFixtureDir(), "expected", "index.md")
		if runtime.GOOS == "windows" {
			f = filepath.Join(indexFixtureDir(), "expected", "index.windows")
		}
		want, err := os.ReadFile(filepath.Clean(f))
		if err != nil {
			t.Fatalf("failed to read expected index: %v", err)
		}

		// Compare the raw bytes. Normalizing line endings away here would hide
		// every change to the document structure, which is exactly what this
		// golden is meant to guard: index.go is the only in-repository user of
		// the fluent API. The platform split above already covers the line
		// ending difference.
		if diff := cmp.Diff(string(want), buf.String()); diff != "" {
			t.Errorf("value is mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestIsMarkdownFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		want bool
	}{
		{path: "README.md", want: true},
		{path: "README.MD", want: true},
		{path: "note.md.bak", want: false},
		{path: "dummy.txt", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()

			if got := isMarkdownFile(tt.path); got != tt.want {
				t.Errorf("isMarkdownFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestGenerateIndexTwice(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	markdownPath := filepath.Join(dir, "sample.md")
	if err := os.WriteFile(markdownPath, []byte("# Sample\n"), 0o600); err != nil {
		t.Fatalf("failed to write markdown file: %v", err)
	}

	if err := GenerateIndex(dir); err != nil {
		t.Fatalf("failed to generate index on first run: %v", err)
	}
	if err := GenerateIndex(dir); err != nil {
		t.Fatalf("failed to generate index on second run: %v", err)
	}

	indexPath := filepath.Join(dir, "index.md")
	//nolint:gosec // indexPath is created from t.TempDir() in this test.
	got, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read generated index: %v", err)
	}
	if strings.Contains(string(got), "(index.md)") {
		t.Fatalf("generated index contains self link: %s", string(got))
	}
}

func TestGenerateIndexClosesFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	markdownPath := filepath.Join(dir, "sample.md")
	if err := os.WriteFile(markdownPath, []byte("# Sample\n"), 0o600); err != nil {
		t.Fatalf("failed to write markdown file: %v", err)
	}

	if err := GenerateIndex(dir); err != nil {
		t.Fatalf("failed to generate index: %v", err)
	}

	indexPath := filepath.Join(dir, "index.md")
	if err := os.Remove(indexPath); err != nil {
		t.Fatalf("failed to remove generated index file: %v", err)
	}
}

// TestGenerateIndexLinksAreRelativeURLs pins the destinations the generated
// index links to.
//
// Two things have to hold for a link to work. It has to be relative to the
// index, which sits in the target directory, and it has to use "/", because a
// markdown link destination is a URL and a backslash in one is an escape
// character rather than a directory separator.
//
// The target directory is given in several shapes because callers write paths
// the way that reads best in Go source, which is usually with forward slashes
// even on Windows.
func TestGenerateIndexLinksAreRelativeURLs(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	docs := filepath.Join(root, "docs")
	if err := os.MkdirAll(filepath.Join(docs, "sub"), 0750); err != nil {
		t.Fatalf("failed to create the fixture tree: %v", err)
	}
	for path, content := range map[string]string{
		filepath.Join(docs, "guide.md"):         "# Guide\n",
		filepath.Join(docs, "sub", "nested.md"): "# Nested\n",
	} {
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatalf("failed to write the fixture %s: %v", path, err)
		}
	}

	tests := map[string]string{
		"the platform separator": docs,
		"a trailing separator":   docs + string(filepath.Separator),
		"forward slashes":        filepath.ToSlash(docs),
	}

	for name, targetDir := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			if err := GenerateIndex(targetDir, WithWriter(&buf)); err != nil {
				t.Fatalf("failed to generate index: %v", err)
			}

			for _, want := range []string{"[Guide](guide.md)", "[Nested](sub/nested.md)"} {
				if !strings.Contains(buf.String(), want) {
					t.Errorf("index = %q, want it to contain %q", buf.String(), want)
				}
			}
			if strings.Contains(buf.String(), `\`) {
				t.Errorf("index = %q, want no backslash: a link destination is a URL", buf.String())
			}
		})
	}
}

// TestLinkToFallsBackToThePathItWasGiven covers the branch where a walked path
// cannot be made relative to the target directory. Dropping the entry would be
// worse than an odd looking link, so the path is used as it is.
func TestLinkToFallsBackToThePathItWasGiven(t *testing.T) {
	t.Parallel()

	// A relative target and an absolute file have no relative path between
	// them, because the working directory is not known to filepath.Rel.
	i := &Index{targetDir: "docs"}

	absolute := filepath.Join(string(filepath.Separator), "elsewhere", "guide.md")
	got := i.linkTo(absolute)
	if want := filepath.ToSlash(absolute); got != want {
		t.Errorf("linkTo(%q) = %q, want %q", absolute, got, want)
	}
}

func TestHasDir(t *testing.T) {
	t.Parallel()

	i := &Index{dir: []*dir{{path: "docs"}, {path: filepath.Join("docs", "api")}}}

	tests := map[string]struct {
		path string
		want bool
	}{
		"a directory the index holds":            {path: "docs", want: true},
		"a nested directory the index holds":     {path: filepath.Join("docs", "api"), want: true},
		"a directory the index does not hold":    {path: "elsewhere", want: false},
		"a prefix of one the index holds":        {path: "doc", want: false},
		"the empty path against a non empty set": {path: "", want: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := i.hasDir(tt.path); got != tt.want {
				t.Errorf("hasDir(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

// TestGenerateIndexWritesIndexFile covers the default destination. Without
// WithWriter the index goes to index.md in the target directory, which is the
// path most callers use and the only one that touches the filesystem.
func TestGenerateIndexWritesIndexFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "guide.md"), []byte("# Guide\n"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v, want nil", err)
	}

	if err := GenerateIndex(dir, WithTitle("Index")); err != nil {
		t.Fatalf("GenerateIndex() = %v, want nil", err)
	}

	got, err := os.ReadFile(filepath.Clean(filepath.Join(dir, "index.md")))
	if err != nil {
		t.Fatalf("ReadFile() = %v, want nil", err)
	}
	for _, want := range []string{"## Index", "[Guide](guide.md)"} {
		if !strings.Contains(string(got), want) {
			t.Errorf("index = %q, want it to contain %q", got, want)
		}
	}

	// Running again has to overwrite rather than append, and the index must not
	// list itself.
	if err := GenerateIndex(dir, WithTitle("Index")); err != nil {
		t.Fatalf("second GenerateIndex() = %v, want nil", err)
	}
	again, err := os.ReadFile(filepath.Clean(filepath.Join(dir, "index.md")))
	if err != nil {
		t.Fatalf("ReadFile() = %v, want nil", err)
	}
	if strings.Count(string(again), "## Index") != 1 {
		t.Errorf("index = %q, want the title exactly once", again)
	}
	if strings.Contains(string(again), "(index.md)") {
		t.Errorf("index = %q, want it not to list itself", again)
	}
}

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
		"indented heading":       {write("indented.md", "   # Title\n"), "Title"},
		// A directory opens without error on some platforms and fails to scan,
		// which is the other way the lookup can come up empty.
		"a directory": {dir, ""},
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
