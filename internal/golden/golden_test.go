package golden

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssertAt(t *testing.T) {
	t.Parallel()

	t.Run("identical content passes", func(t *testing.T) {
		t.Parallel()

		path := writeTemp(t, "# Title\n\nbody\n")
		if err := assertAt(path, "# Title\n\nbody\n"); err != nil {
			t.Errorf("assertAt() = %v, want nil", err)
		}
	})

	t.Run("CRLF output matches an LF golden file", func(t *testing.T) {
		t.Parallel()

		path := writeTemp(t, "# Title\n\nbody\n")
		if err := assertAt(path, "# Title\r\n\r\nbody\r\n"); err != nil {
			t.Errorf("assertAt() = %v, want nil", err)
		}
	})

	t.Run("a changed line is reported with its line number", func(t *testing.T) {
		t.Parallel()

		path := writeTemp(t, "# Title\nfirst\nsecond\n")
		err := assertAt(path, "# Title\nchanged\nsecond\n")
		if err == nil {
			t.Fatal("assertAt() = nil, want an error")
		}
		for _, want := range []string{"line 2", `"first"`, `"changed"`} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("assertAt() error = %q, want it to contain %q", err, want)
			}
		}
	})

	t.Run("an extra line is reported", func(t *testing.T) {
		t.Parallel()

		path := writeTemp(t, "a\nb")
		err := assertAt(path, "a\nb\nc")
		if err == nil {
			t.Fatal("assertAt() = nil, want an error")
		}
		if !strings.Contains(err.Error(), "extra line") {
			t.Errorf("assertAt() error = %q, want it to mention an extra line", err)
		}
	})

	t.Run("a missing line is reported", func(t *testing.T) {
		t.Parallel()

		path := writeTemp(t, "a\nb\nc")
		err := assertAt(path, "a\nb")
		if err == nil {
			t.Fatal("assertAt() = nil, want an error")
		}
		if !strings.Contains(err.Error(), "fewer line") {
			t.Errorf("assertAt() error = %q, want it to mention a missing line", err)
		}
	})

	t.Run("a missing golden file explains how to create it", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "absent.md")
		err := assertAt(path, "anything")
		if err == nil {
			t.Fatal("assertAt() = nil, want an error")
		}
		if !strings.Contains(err.Error(), UpdateEnv) {
			t.Errorf("assertAt() error = %q, want it to mention %s", err, UpdateEnv)
		}
	})
}

func TestAssertAtUpdatesGoldenFile(t *testing.T) {
	// The update mode is driven by an environment variable, so this test cannot
	// run in parallel with the comparison tests above.
	t.Setenv(UpdateEnv, "1")

	path := filepath.Join(t.TempDir(), "nested", "created.md")
	if err := assertAt(path, "# Created\r\n"); err != nil {
		t.Fatalf("assertAt() = %v, want nil", err)
	}

	got, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("ReadFile() = %v, want nil", err)
	}
	if want := "# Created\n"; string(got) != want {
		t.Errorf("golden file = %q, want %q (line endings must be normalized)", got, want)
	}
}

func TestAssertUsesTheGoldenDirectory(t *testing.T) {
	t.Parallel()

	// The fixture below lives at testdata/golden/fixture.md, which is the layout
	// every package using this helper relies on.
	if err := Assert("fixture.md", "# Fixture\n"); err != nil {
		t.Errorf("Assert() = %v, want nil", err)
	}
}

func TestDescribeWithoutADifferingLine(t *testing.T) {
	t.Parallel()

	if got, want := describe("same", "same"), "no differing line"; !strings.Contains(got, want) {
		t.Errorf("describe() = %q, want it to contain %q", got, want)
	}
}

func TestWriteReportsAnUnusableDirectory(t *testing.T) {
	t.Parallel()

	// A regular file cannot hold children, so MkdirAll fails on this path.
	file := filepath.Join(t.TempDir(), "regular")
	if err := os.WriteFile(file, []byte("x"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v, want nil", err)
	}

	if err := write(filepath.Join(file, "child", "golden.md"), "content"); err == nil {
		t.Error("write() = nil, want an error")
	}
}

// writeTemp writes content to a fresh temporary file and returns its path.
func writeTemp(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "golden.md")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("WriteFile() = %v, want nil", err)
	}
	return path
}
