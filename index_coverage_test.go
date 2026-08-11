package markdown

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
