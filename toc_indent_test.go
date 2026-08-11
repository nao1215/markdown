package markdown

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal"
)

// tocEntries returns the lines between the table of contents markers.
func tocEntries(t *testing.T, rendered string) []string {
	t.Helper()

	begin := strings.Index(rendered, TableOfContentsMarkerBegin)
	end := strings.Index(rendered, TableOfContentsMarkerEnd)
	if begin == -1 || end == -1 {
		t.Fatalf("markers missing from output:\n%s", rendered)
	}

	body := rendered[begin+len(TableOfContentsMarkerBegin) : end]
	entries := []string{}
	for _, line := range strings.Split(body, internal.LineFeed()) {
		if strings.TrimSpace(line) != "" {
			entries = append(entries, line)
		}
	}
	return entries
}

// TestTableOfContentsIndentsFromTheShallowestHeading covers the case where the
// document has no H1. TableOfContents pins MinDepth at H1, so the indent used
// to be measured from a level that never appears and every entry came out
// indented under nothing.
func TestTableOfContentsIndentsFromTheShallowestHeading(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(*Markdown) *Markdown
		want  []string
	}{
		"document starting at H2": {
			build: func(m *Markdown) *Markdown {
				return m.TableOfContents(TableOfContentsDepthH3).H2("S1").H3("S1.1")
			},
			want: []string{"- [S1](#s1)", "  - [S1.1](#s11)"},
		},
		"document with an H1 is unchanged": {
			build: func(m *Markdown) *Markdown {
				return m.H1("Doc").TableOfContents(TableOfContentsDepthH3).H2("S1").H3("S1.1")
			},
			want: []string{"- [Doc](#doc)", "  - [S1](#s1)", "    - [S1.1](#s11)"},
		},
		"document starting at H3": {
			build: func(m *Markdown) *Markdown {
				return m.TableOfContents(TableOfContentsDepthH4).H3("Deep").H4("Deeper")
			},
			want: []string{"- [Deep](#deep)", "  - [Deeper](#deeper)"},
		},
		"shallowest heading appears late": {
			build: func(m *Markdown) *Markdown {
				return m.TableOfContents(TableOfContentsDepthH3).H3("First").H2("Later")
			},
			want: []string{"  - [First](#first)", "- [Later](#later)"},
		},
		"explicit range is unaffected": {
			build: func(m *Markdown) *Markdown {
				return m.TableOfContentsWithRange(TableOfContentsDepthH2, TableOfContentsDepthH4).
					H2("S1").H3("S1.1")
			},
			want: []string{"- [S1](#s1)", "  - [S1.1](#s11)"},
		},
		"headings outside the range are skipped": {
			build: func(m *Markdown) *Markdown {
				return m.TableOfContentsWithRange(TableOfContentsDepthH2, TableOfContentsDepthH2).
					H1("Title").H2("S1").H3("S1.1")
			},
			want: []string{"- [S1](#s1)"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			if err := tt.build(NewMarkdown(&buf)).Build(); err != nil {
				t.Fatalf("failed to build: %v", err)
			}

			if diff := cmp.Diff(tt.want, tocEntries(t, buf.String())); diff != "" {
				t.Errorf("table of contents mismatch (-want +got):\n%s\nfull output:\n%s", diff, buf.String())
			}
		})
	}
}

// TestTableOfContentsWithNoMatchingHeadings makes sure the indent baseline does
// not misbehave when nothing falls inside the requested range.
func TestTableOfContentsWithNoMatchingHeadings(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := NewMarkdown(&buf).
		TableOfContentsWithRange(TableOfContentsDepthH4, TableOfContentsDepthH6).
		H1("Title").
		H2("Section").
		Build(); err != nil {
		t.Fatalf("failed to build: %v", err)
	}

	if entries := tocEntries(t, buf.String()); len(entries) != 0 {
		t.Errorf("expected no entries, got %v", entries)
	}
}
