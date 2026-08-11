package markdown

import (
	"bytes"
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
