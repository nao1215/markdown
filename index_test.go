package markdown

import (
	"bytes"
	"os"
	"path/filepath"
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

		want, err := os.ReadFile(filepath.Join("testdata", "expected", "index.md"))
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
