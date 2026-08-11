// Package golden compares generated markdown against committed golden files.
//
// The generated output of every builder is part of this library's public
// contract: a document that renders correctly today must keep rendering the
// same way after any refactoring. Unit tests assert the pieces, and the golden
// files here pin the whole document, so a change to a single literal anywhere
// in the library shows up as a failing test instead of as a diff in a user's
// repository.
package golden

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UpdateEnv is the environment variable that switches Assert from comparing to
// rewriting golden files. Run "UPDATE_GOLDEN=1 go test ./..." after an
// intentional output change, then review the resulting diff before committing.
const UpdateEnv = "UPDATE_GOLDEN"

const (
	// dirName is the directory, relative to the package under test, that holds
	// the golden files.
	dirName = "testdata"
	// subDirName is the directory inside dirName that holds the golden files.
	subDirName = "golden"
	// dirPerm is the permission of a created golden directory.
	dirPerm = 0750
	// filePerm is the permission of a written golden file.
	filePerm = 0600
)

// Assert compares got with the contents of the golden file testdata/golden/<name>.
//
// Line endings are normalized to "\n" on both sides before comparing. The
// builders emit the line ending of the current platform, while the golden files
// are committed with "\n", so a byte for byte comparison would fail on Windows
// for output that is in fact correct.
//
// When UpdateEnv is set to a non-empty value the golden file is written with
// got, using "\n" line endings, and Assert returns nil.
func Assert(name, got string) error {
	return assertAt(filepath.Join(dirName, subDirName, name), got)
}

// assertAt is Assert against an explicit path. It exists so the comparison can
// be tested against a temporary directory rather than against the golden files
// of this package.
func assertAt(path, got string) error {
	got = normalize(got)

	if os.Getenv(UpdateEnv) != "" {
		return write(path, got)
	}

	raw, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("read golden file: %w (run the test with %s=1 to create it)", err, UpdateEnv)
	}

	want := normalize(string(raw))
	if want == got {
		return nil
	}
	return fmt.Errorf("output does not match %s:\n%s\n(run the test with %s=1 to update the golden file)",
		path, describe(want, got), UpdateEnv)
}

// write stores content as the golden file at path, creating the directory when
// it does not exist yet.
func write(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("create golden directory: %w", err)
	}
	if err := os.WriteFile(path, []byte(content), filePerm); err != nil {
		return fmt.Errorf("write golden file: %w", err)
	}
	return nil
}

// normalize rewrites every line ending to "\n" so that output produced on
// Windows compares equal to a golden file committed with "\n".
func normalize(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// describe reports the first line at which want and got diverge.
//
// Golden documents are long, so printing both in full buries the one line that
// actually changed. The line number is 1 based to match what an editor shows.
func describe(want, got string) string {
	wantLines := strings.Split(want, "\n")
	gotLines := strings.Split(got, "\n")

	for i := 0; i < len(wantLines) && i < len(gotLines); i++ {
		if wantLines[i] != gotLines[i] {
			return fmt.Sprintf("first difference at line %d:\n  want: %q\n  got:  %q", i+1, wantLines[i], gotLines[i])
		}
	}

	switch {
	case len(gotLines) > len(wantLines):
		return fmt.Sprintf("got %d extra line(s), the first is line %d:\n  got: %q",
			len(gotLines)-len(wantLines), len(wantLines)+1, gotLines[len(wantLines)])
	case len(wantLines) > len(gotLines):
		return fmt.Sprintf("got %d fewer line(s), the first missing one is line %d:\n  want: %q",
			len(wantLines)-len(gotLines), len(gotLines)+1, wantLines[len(gotLines)])
	default:
		return "the documents differ but no differing line was found"
	}
}
