package markdown_test

import (
	"io"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/flowchart"
	"github.com/yuin/goldmark/ast"
	"gopkg.in/yaml.v3"
)

// Every defect this package has shipped lived on a boundary between the text a
// caller passes in and the syntax that text is dropped into: a pipe ends a table
// cell, a newline ends a row, a colon ends a YAML key, a backtick run ends a
// fence. The targets in fuzz_test.go cover the table and alert boundaries. The
// ones here cover the rest, with a real parser as the oracle rather than a
// string comparison, because the failure mode is always a document that looks
// fine and means something else.

// boundarySeeds are the inputs that have historically broken one boundary or
// another, plus the shapes of text that reach a builder from a database or a
// command's output.
func boundarySeeds() []string {
	return []string{
		"",
		"plain text",
		"a: b",
		"# comment",
		"*ref",
		"&anchor",
		"~",
		"true",
		"123",
		"  leading and trailing  ",
		"---",
		"...",
		"line\nbreak",
		"line\r\nbreak",
		"carriage\rreturn",
		"null\x00byte",
		"bell\x07",
		"tab\tseparated",
		`quote"inside`,
		`back\slash`,
		"emoji 🎉 here",
		"日本語のタイトル",
		"مرحبا بالعالم",
		"```",
		"````",
		"a|b",
		"[link](https://go.dev)",
		strings.Repeat("long ", 2048),
	}
}

// FuzzFrontMatterTitle asserts that a diagram title survives the YAML front
// matter it is written into.
//
// mermaid runs the front matter through a YAML parser before it draws anything,
// so a title is not a string in a template but a value in a document. A bare
// scalar built from arbitrary text can end the value early, turn into a comment,
// resolve to a boolean, or fail to parse at all, and any of those loses the
// whole diagram rather than one line of it.
func FuzzFrontMatterTitle(f *testing.F) {
	for _, seed := range boundarySeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, title string) {
		if strings.TrimSpace(title) == "" {
			// A title that is blank once trimmed writes no front matter at all,
			// which is a different property and is covered by the unit tests.
			return
		}
		if !utf8.ValidString(title) {
			// YAML is defined over Unicode, so a title that is not valid UTF-8
			// has no faithful representation in it. See the note on this target
			// in the pull request that added it.
			return
		}

		diagram := flowchart.NewFlowchart(io.Discard, flowchart.WithTitle(title)).
			NodeWithText("A", "Node A").
			String()

		block, ok := frontMatter(diagram)
		if !ok {
			t.Fatalf("the diagram built with title %q has no front matter block:\n%s", title, diagram)
		}

		var parsed struct {
			Title string `yaml:"title"`
		}
		if err := yaml.Unmarshal([]byte(block), &parsed); err != nil {
			t.Fatalf("the front matter of the diagram built with title %q is not valid YAML: %v\n%s", title, err, block)
		}
		if parsed.Title != title {
			t.Errorf("the front matter title parsed as %q, want %q", parsed.Title, title)
		}
	})
}

// frontMatter returns the body between the first two "---" lines of a diagram.
func frontMatter(diagram string) (string, bool) {
	lines := strings.Split(strings.ReplaceAll(diagram, "\r\n", "\n"), "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return "", false
	}
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			return strings.Join(lines[1:i], "\n"), true
		}
	}
	return "", false
}

// FuzzCodeBlockContent asserts that content placed in a fenced block comes back
// out of a parser unchanged.
//
// The property only holds for content that does not open a fence of its own:
// three backticks at the start of a line end the block, and the rest of the
// document spills onto the page as prose. That case is a documented limitation
// of CodeBlocks rather than something this target can assert, so it is skipped
// here and pinned separately by TestACodeBlockEndsAtAFenceInItsContent.
func FuzzCodeBlockContent(f *testing.F) {
	for _, seed := range boundarySeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, content string) {
		if holdsFence(content) {
			return
		}

		document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
			return m.CodeBlocks(markdown.SyntaxHighlightGo, content)
		})

		root, source := parse(t, document)
		blocks := nodesOfKind(root, ast.KindFencedCodeBlock)
		if len(blocks) != 1 {
			t.Fatalf("content %q produced %d fenced code blocks, want 1:\n%s", content, len(blocks), document)
		}

		block, ok := blocks[0].(*ast.FencedCodeBlock)
		if !ok {
			t.Fatalf("the parsed node is %T, want *ast.FencedCodeBlock", blocks[0])
		}

		var got strings.Builder
		for i := 0; i < block.Lines().Len(); i++ {
			line := block.Lines().At(i)
			got.Write(line.Value(source))
		}

		// The builder writes a line ending after the content, and a parser
		// reports the block's lines with theirs, so compare without the one at
		// the end. Content whose own last line is empty is indistinguishable
		// from that, which is why both sides are trimmed rather than one.
		if want, have := normalizeContent(content), normalizeContent(got.String()); want != have {
			t.Errorf("content %q came back as %q", want, have)
		}
	})
}

// holdsFence reports whether content has a line that opens or closes a fence.
func holdsFence(content string) bool {
	for _, line := range strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n") {
		if strings.HasPrefix(strings.TrimLeft(line, " "), "```") {
			return true
		}
	}
	return false
}

// normalizeContent puts code block content in the one shape a round trip can be
// compared in: every line ending resolved to "\n", and the trailing one dropped.
//
// A parser reads "\r", "\n", and "\r\n" as the same thing, so a round trip
// through a document cannot preserve which one the caller wrote. Nor can a
// document tell content ending in line breaks apart from content that does not,
// since the closing fence has to start on its own line either way, which is why
// every trailing one is dropped rather than a single one.
func normalizeContent(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.TrimRight(s, "\n")
}

// TestACodeBlockEndsAtAFenceInItsContent pins the limitation FuzzCodeBlockContent
// skips. CodeBlocks writes a three backtick fence, so content holding a line of
// three backticks closes the block early and everything after it renders as
// prose. Callers whose content can hold a fence have to say so themselves; this
// test exists so the limitation is a stated one rather than a surprise.
func TestACodeBlockEndsAtAFenceInItsContent(t *testing.T) {
	t.Parallel()

	document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
		return m.CodeBlocks(markdown.SyntaxHighlightMermaid, "before\n```\nafter")
	})

	root, _ := parse(t, document)
	if got := len(nodesOfKind(root, ast.KindFencedCodeBlock)); got < 2 {
		t.Errorf("parsed %d fenced code blocks, want at least 2: the limitation this test documents no longer holds", got)
	}
}

// FuzzHeadingText asserts that a heading built from arbitrary text stays one
// heading, for the text a heading can hold.
//
// A heading is a single line by construction, so text carrying a line break is
// out of scope here: markdown has no multi-line heading to put it in. The rest
// of the punctuation that means something at the start of a line, "#" above all,
// has to be harmless.
func FuzzHeadingText(f *testing.F) {
	for _, seed := range boundarySeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		if strings.ContainsAny(text, "\r\n") {
			return
		}

		document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
			return m.H2(text)
		})

		root, _ := parse(t, document)
		if got := len(nodesOfKind(root, ast.KindHeading)); got != 1 {
			t.Errorf("H2(%q) produced %d headings, want 1:\n%s", text, got, document)
		}
	})
}

// FuzzLinkAndImageText asserts that a link or an image stays inline.
//
// The destination is not asserted for arbitrary input, because a URL holding a
// space or an unbalanced parenthesis has no inline form and the caller has to
// encode it. What has to hold for any input is that the helper cannot introduce
// block structure: a link built from a database value must not turn into a
// heading, a list, or a fenced block halfway through a paragraph.
func FuzzLinkAndImageText(f *testing.F) {
	for _, seed := range boundarySeeds() {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		if strings.ContainsAny(text, "\r\n") {
			// A line break in the text ends the paragraph, which is a property
			// of markdown rather than of these helpers.
			return
		}

		for name, inline := range map[string]string{
			"link":  markdown.Link(text, "https://example.com"),
			"image": markdown.Image(text, "https://example.com/image.png"),
			"url":   markdown.Link("the text", text),
		} {
			document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
				return m.PlainTextf("before %s after", inline)
			})

			root, _ := parse(t, document)
			for _, kind := range []ast.NodeKind{
				ast.KindHeading,
				ast.KindList,
				ast.KindFencedCodeBlock,
				ast.KindBlockquote,
				ast.KindThematicBreak,
			} {
				if got := len(nodesOfKind(root, kind)); got != 0 {
					t.Errorf("%s built from %q produced %d %s nodes, want 0:\n%s", name, text, got, kind, document)
				}
			}
		}
	})
}

// FuzzSafeLinkDestinationRoundTrips asserts the destination itself, for the
// destinations that have an inline form at all.
func FuzzSafeLinkDestinationRoundTrips(f *testing.F) {
	for _, seed := range []string{
		"https://go.dev",
		"https://example.com/page?a=1&b=2",
		"./relative/path.md",
		"#anchor",
		"mailto:someone@example.com",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, url string) {
		if url == "" || !isInlineDestination(url) {
			return
		}

		document := build(t, func(m *markdown.Markdown) *markdown.Markdown {
			return m.PlainText(markdown.Link("the text", url))
		})

		root, _ := parse(t, document)
		links := nodesOfKind(root, ast.KindLink)
		if len(links) != 1 {
			t.Fatalf("url %q produced %d links, want 1:\n%s", url, len(links), document)
		}

		link, ok := links[0].(*ast.Link)
		if !ok {
			t.Fatalf("the parsed node is %T, want *ast.Link", links[0])
		}
		if got := string(link.Destination); got != url {
			t.Errorf("url %q came back as %q", url, got)
		}
	})
}

// isInlineDestination reports whether url can be written between parentheses
// without being encoded first. Whitespace ends the destination, parentheses and
// angle brackets have to balance, and a control character is not a URL at all.
func isInlineDestination(url string) bool {
	if strings.ContainsAny(url, "()<>\\") {
		return false
	}
	for _, r := range url {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return false
		}
	}
	return true
}
