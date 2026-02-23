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

func TestGenerateIndex(t *testing.T) {
	t.Parallel()

	t.Run("create index", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		if err := GenerateIndex(
			"testdata",
			WithTitle("Test Title"),
			WithDescription([]string{"Test Description", "Next Description"}),
			WithWriter(&buf),
		); err != nil {
			t.Fatalf("failed to generate index: %v", err)
		}

		f := filepath.Join("testdata", "expected", "index.md")
		if runtime.GOOS == "windows" {
			f = filepath.Join("testdata", "expected", "index.windows")
		}
		want, err := os.ReadFile(filepath.Clean(f))
		if err != nil {
			t.Fatalf("failed to read expected index: %v", err)
		}

		expect := strings.ReplaceAll(string(want), "\r\n", "\n")
		expect = strings.ReplaceAll(expect, "\n", "")
		got := strings.ReplaceAll(buf.String(), "\r\n", "\n")
		got = strings.ReplaceAll(got, "\n", "")

		if diff := cmp.Diff(expect, got); diff != "" {
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
		tt := tt
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()

			if got := isMarkdownFile(tt.path); got != tt.want {
				t.Errorf("isMarkdownFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
