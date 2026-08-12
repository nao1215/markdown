package venn_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/nao1215/markdown/internal"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/venn"
)

// errWrite is the failure the writer below reports, so the test can assert that
// Build passed it through rather than inventing an error of its own.
var errWrite = errors.New("write failed")

// errWriter fails every write, which is what a full disk or a closed pipe looks
// like to Build.
type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errWrite
}

// lines joins the given lines with the line ending the library writes on this
// platform, which is what String returns them joined with.
func lines(want ...string) string {
	return strings.Join(want, internal.LineFeed())
}

func TestDiagram(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(io.Writer) *venn.Diagram
		want  string
	}{
		"an empty diagram is the header alone": {
			build: func(w io.Writer) *venn.Diagram { return venn.NewDiagram(w) },
			want:  "venn-beta",
		},
		"a set is drawn with its name": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w).Set("go").Set("rust")
			},
			want: lines("venn-beta", "    set go", "    set rust"),
		},
		"a set can carry a label instead": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w).SetWithLabel("go", "Written in Go")
			},
			want: lines("venn-beta", `    set go["Written in Go"]`),
		},
		"a title is a statement rather than front matter": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w, venn.WithTitle("What they share"))
			},
			want: lines("venn-beta", "    title What they share"),
		},
		"a title keeps the quotation marks a caller wrote": {
			// mermaid reads the rest of the line, so quoting would draw the
			// quotation marks themselves.
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w, venn.WithTitle(`The "core"`))
			},
			want: lines("venn-beta", `    title The "core"`),
		},
		"a title escapes the punctuation its lexer has taken": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w, venn.WithTitle("Sets #1; v2"))
			},
			want: lines("venn-beta", "    title Sets #35;1#59; v2"),
		},
		"a quotation mark in a label becomes the entity mermaid decodes": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w).SetWithLabel("a", `The "core"`)
			},
			want: lines("venn-beta", `    set a["The #quot;core#quot;"]`),
		},
		"a hash in a label is escaped only where it would start an entity": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w).SetWithLabel("a", "PR #123 and #quot;")
			},
			want: lines("venn-beta", `    set a["PR #123 and #35;quot;"]`),
		},
		"the punctuation a label may carry is left alone": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w).SetWithLabel("a", `x'#;[](){}<br/>🎉日本語:,*-|%%\x`)
			},
			want: lines("venn-beta", `    set a["x'#;[](){}<br/>🎉日本語:,*-|%%\x"]`),
		},
		"a hyphen and an underscore are names": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w).Set("built-in").Set("third_party")
			},
			want: lines("venn-beta", "    set built-in", "    set third_party"),
		},
		"text is trimmed": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w).SetWithLabel("  a  ", "  Alpha  ")
			},
			want: lines("venn-beta", `    set a["Alpha"]`),
		},
		"LF adds a blank line": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w).Set("a").LF().Set("b")
			},
			want: lines("venn-beta", "    set a", "", "    set b"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			d := tt.build(buf)
			if err := d.Build(); err != nil {
				t.Fatalf("Build() = %v, want nil", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("Build() wrote\n%q\nwant\n%q", got, tt.want)
			}
			if got := d.String(); got != tt.want {
				t.Errorf("String() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestDiagramRecordsBadInput(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(io.Writer) *venn.Diagram
		want  string
	}{
		"a title spanning lines": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w, venn.WithTitle("a\nb"))
			},
			want: "title must not contain newline characters",
		},
		"an empty set name": {
			build: func(w io.Writer) *venn.Diagram { return venn.NewDiagram(w).Set("  ") },
			want:  "set name must not be empty",
		},
		"a full stop in a set name": {
			// mermaid loses the whole diagram on one, and there is nothing to
			// escape to, so it is reported rather than mangled.
			build: func(w io.Writer) *venn.Diagram { return venn.NewDiagram(w).Set("a.b") },
			want:  `set name "a.b" must hold only letters, digits, underscores and hyphens`,
		},
		"a space in a set name": {
			build: func(w io.Writer) *venn.Diagram { return venn.NewDiagram(w).Set("a b") },
			want:  `set name "a b" must hold only`,
		},
		"a comma in a set name": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w).SetWithLabel("a,b", "L")
			},
			want: `set name "a,b" must hold only`,
		},
		"an empty label": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w).SetWithLabel("a", "  ")
			},
			want: "set label must not be empty",
		},
		"a label spanning lines": {
			build: func(w io.Writer) *venn.Diagram {
				return venn.NewDiagram(w).SetWithLabel("a", "x\ny")
			},
			want: "set label must not contain newline characters",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			d := tt.build(io.Discard)
			err := d.Build()
			if err == nil {
				t.Fatalf("Build() = nil, want an error mentioning %q", tt.want)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Build() = %v, want it to mention %q", err, tt.want)
			}
			if d.Error() == nil {
				t.Error("Error() = nil after a failed Build()")
			}
		})
	}
}

// TestTheFirstErrorIsKept covers the rule the whole library shares: the chain
// runs to the end after a bad call, and the error that surfaces is the first
// one, because that is the one that explains the rest.
func TestTheFirstErrorIsKept(t *testing.T) {
	t.Parallel()

	err := venn.NewDiagram(io.Discard).Set("a b").Set("c.d").Build()
	if err == nil {
		t.Fatal("Build() = nil, want the first error")
	}
	if want := `set name "a b"`; !strings.Contains(err.Error(), want) {
		t.Errorf("Build() = %v, want the first error %q", err, want)
	}
}

// TestBuildWithNilWriter covers the case where a diagram is built for String()
// only and Build() is called by mistake.
func TestBuildWithNilWriter(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Build() panicked with a nil writer: %v", r)
		}
	}()

	d := venn.NewDiagram(nil).Set("a")
	if d.String() == "" {
		t.Fatal("String() returned nothing for a diagram with a set in it")
	}

	err := d.Build()
	if err == nil {
		t.Fatal("Build() with a nil writer must return an error")
	}
	if err.Error() != "output writer must not be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestBuildReportsWriteFailure covers the branch where the destination accepts
// the diagram and then fails.
func TestBuildReportsWriteFailure(t *testing.T) {
	t.Parallel()

	err := venn.NewDiagram(errWriter{}).Set("a").Build()
	if err == nil {
		t.Fatal("Build must report a failing writer")
	}
	if !errors.Is(err, errWrite) {
		t.Errorf("Build lost the destination error: %v", err)
	}
}

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return venn.NewDiagram(w).Set("a").SetWithLabel("b", "Beta")
	})
}

// TestRecordedErrorContract asserts that a set name mermaid cannot read
// surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return venn.NewDiagram(w).Set("a b")
	})
}

// TestGoldenVenn pins the rendered diagram of every builder method and every
// option of this package.
func TestGoldenVenn(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := venn.NewDiagram(buf, venn.WithTitle(`Languages #1; the "core"`)).
		Set("go").
		SetWithLabel("rust", `Rust ("safe")`).
		LF().
		SetWithLabel("both", "Compiled & typed").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("venn.md", buf.String()); err != nil {
		t.Error(err)
	}
}
