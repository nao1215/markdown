// Package markdown is a simple markdown builder.
//
// A document is one method chain: call [NewMarkdown] with a writer, add blocks
// in the order they should appear, and finish with [Markdown.Build]. The output
// follows GitHub Flavored Markdown.
//
// Nested structures, such as a list inside a list item, are out of scope. They
// would turn the chain into a tree.
//
// The builder records errors instead of returning them from every call. Nothing
// panics on bad input, and a rejected call does not stop the document: the
// chain runs to the end, and [Markdown.Error] and [Markdown.Build] both report
// the first error it recorded.
//
// [Markdown.String] returns the document without needing a writer. That is how
// the mermaid subpackages hand a diagram to [Markdown.CodeBlocks].
//
// A builder is not safe for concurrent use. Build one document per goroutine.
package markdown

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/nao1215/markdown/internal"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

// SyntaxHighlight is syntax highlight language.
type SyntaxHighlight string

const (
	// SyntaxHighlightNone is no syntax highlight.
	SyntaxHighlightNone SyntaxHighlight = ""
	// SyntaxHighlightText is syntax highlight for text.
	SyntaxHighlightText SyntaxHighlight = "text"
	// SyntaxHighlightAPIBlueprint is syntax highlight for API Blueprint.
	SyntaxHighlightAPIBlueprint SyntaxHighlight = "markdown"
	// SyntaxHighlightShell is syntax highlight for Shell.
	SyntaxHighlightShell SyntaxHighlight = "shell"
	// SyntaxHighlightGo is syntax highlight for Go.
	SyntaxHighlightGo SyntaxHighlight = "go"
	// SyntaxHighlightJSON is syntax highlight for JSON.
	SyntaxHighlightJSON SyntaxHighlight = "json"
	// SyntaxHighlightYAML is syntax highlight for YAML.
	SyntaxHighlightYAML SyntaxHighlight = "yaml"
	// SyntaxHighlightXML is syntax highlight for XML.
	SyntaxHighlightXML SyntaxHighlight = "xml"
	// SyntaxHighlightHTML is syntax highlight for HTML.
	SyntaxHighlightHTML SyntaxHighlight = "html"
	// SyntaxHighlightCSS is syntax highlight for CSS.
	SyntaxHighlightCSS SyntaxHighlight = "css"
	// SyntaxHighlightJavaScript is syntax highlight for JavaScript.
	SyntaxHighlightJavaScript SyntaxHighlight = "javascript"
	// SyntaxHighlightTypeScript is syntax highlight for TypeScript.
	SyntaxHighlightTypeScript SyntaxHighlight = "typescript"
	// SyntaxHighlightSQL is syntax highlight for SQL.
	SyntaxHighlightSQL SyntaxHighlight = "sql"
	// SyntaxHighlightC is syntax highlight for C.
	SyntaxHighlightC SyntaxHighlight = "c"
	// SyntaxHighlightCSharp is syntax highlight for C#.
	SyntaxHighlightCSharp SyntaxHighlight = "csharp"
	// SyntaxHighlightCPlusPlus is syntax highlight for C++.
	SyntaxHighlightCPlusPlus SyntaxHighlight = "cpp"
	// SyntaxHighlightJava is syntax highlight for Java.
	SyntaxHighlightJava SyntaxHighlight = "java"
	// SyntaxHighlightKotlin is syntax highlight for Kotlin.
	SyntaxHighlightKotlin SyntaxHighlight = "kotlin"
	// SyntaxHighlightPHP is syntax highlight for PHP.
	SyntaxHighlightPHP SyntaxHighlight = "php"
	// SyntaxHighlightPython is syntax highlight for Python.
	SyntaxHighlightPython SyntaxHighlight = "python"
	// SyntaxHighlightRuby is syntax highlight for Ruby.
	SyntaxHighlightRuby SyntaxHighlight = "ruby"
	// SyntaxHighlightSwift is syntax highlight for Swift.
	SyntaxHighlightSwift SyntaxHighlight = "swift"
	// SyntaxHighlightScala is syntax highlight for Scala.
	SyntaxHighlightScala SyntaxHighlight = "scala"
	// SyntaxHighlightRust is syntax highlight for Rust.
	SyntaxHighlightRust SyntaxHighlight = "rust"
	// SyntaxHighlightObjectiveC is syntax highlight for Objective-C.
	SyntaxHighlightObjectiveC SyntaxHighlight = "objectivec"
	// SyntaxHighlightPerl is syntax highlight for Perl.
	SyntaxHighlightPerl SyntaxHighlight = "perl"
	// SyntaxHighlightLua is syntax highlight for Lua.
	SyntaxHighlightLua SyntaxHighlight = "lua"
	// SyntaxHighlightDart is syntax highlight for Dart.
	SyntaxHighlightDart SyntaxHighlight = "dart"
	// SyntaxHighlightClojure is syntax highlight for Clojure.
	SyntaxHighlightClojure SyntaxHighlight = "clojure"
	// SyntaxHighlightGroovy is syntax highlight for Groovy.
	SyntaxHighlightGroovy SyntaxHighlight = "groovy"
	// SyntaxHighlightR is syntax highlight for R.
	SyntaxHighlightR SyntaxHighlight = "r"
	// SyntaxHighlightHaskell is syntax highlight for Haskell.
	SyntaxHighlightHaskell SyntaxHighlight = "haskell"
	// SyntaxHighlightErlang is syntax highlight for Erlang.
	SyntaxHighlightErlang SyntaxHighlight = "erlang"
	// SyntaxHighlightElixir is syntax highlight for Elixir.
	SyntaxHighlightElixir SyntaxHighlight = "elixir"
	// SyntaxHighlightOCaml is syntax highlight for OCaml.
	SyntaxHighlightOCaml SyntaxHighlight = "ocaml"
	// SyntaxHighlightJulia is syntax highlight for Julia.
	SyntaxHighlightJulia SyntaxHighlight = "julia"
	// SyntaxHighlightScheme is syntax highlight for Scheme.
	SyntaxHighlightScheme SyntaxHighlight = "scheme"
	// SyntaxHighlightFSharp is syntax highlight for F#.
	SyntaxHighlightFSharp SyntaxHighlight = "fsharp"
	// SyntaxHighlightCoffeeScript is syntax highlight for CoffeeScript.
	SyntaxHighlightCoffeeScript SyntaxHighlight = "coffeescript"
	// SyntaxHighlightVBNet is syntax highlight for VB.NET.
	SyntaxHighlightVBNet SyntaxHighlight = "vbnet"
	// SyntaxHighlightTeX is syntax highlight for TeX.
	SyntaxHighlightTeX SyntaxHighlight = "tex"
	// SyntaxHighlightDiff is syntax highlight for Diff.
	SyntaxHighlightDiff SyntaxHighlight = "diff"
	// SyntaxHighlightApache is syntax highlight for Apache.
	SyntaxHighlightApache SyntaxHighlight = "apache"
	// SyntaxHighlightDockerfile is syntax highlight for Dockerfile.
	SyntaxHighlightDockerfile SyntaxHighlight = "dockerfile"
	// SyntaxHighlightMermaid is syntax highlight for Mermaid.
	SyntaxHighlightMermaid SyntaxHighlight = "mermaid"
)

// TableOfContentsDepth represents the depth level for table of contents.
type TableOfContentsDepth int

const (
	// TableOfContentsDepthH1 includes only H1 headers in the table of contents.
	TableOfContentsDepthH1 TableOfContentsDepth = 1
	// TableOfContentsDepthH2 includes H1 and H2 headers in the table of contents.
	TableOfContentsDepthH2 TableOfContentsDepth = 2
	// TableOfContentsDepthH3 includes H1, H2, and H3 headers in the table of contents.
	TableOfContentsDepthH3 TableOfContentsDepth = 3
	// TableOfContentsDepthH4 includes H1, H2, H3, and H4 headers in the table of contents.
	TableOfContentsDepthH4 TableOfContentsDepth = 4
	// TableOfContentsDepthH5 includes H1, H2, H3, H4, and H5 headers in the table of contents.
	TableOfContentsDepthH5 TableOfContentsDepth = 5
	// TableOfContentsDepthH6 includes all headers (H1 through H6) in the table of contents.
	TableOfContentsDepthH6 TableOfContentsDepth = 6
)

const (
	// TableOfContentsMarkerBegin is the marker for the beginning of the table of contents.
	TableOfContentsMarkerBegin = "<!-- BEGIN_TOC -->"
	// TableOfContentsMarkerEnd is the marker for the end of the table of contents.
	TableOfContentsMarkerEnd = "<!-- END_TOC -->"
)

// TableOfContentsOptions contains options for generating the table of contents.
type TableOfContentsOptions struct {
	// MinDepth is the minimum header level to include (e.g., 2 for H2 and deeper).
	MinDepth TableOfContentsDepth
	// MaxDepth is the maximum header level to include (e.g., 4 for H4 and shallower).
	MaxDepth TableOfContentsDepth
}

// headerInfo stores information about a header for table of contents generation.
type headerInfo struct {
	level TableOfContentsDepth
	text  string
}

// Markdown is markdown text.
type Markdown struct {
	// body is markdown body.
	body []string
	// dest is output destination for markdown body.
	dest io.Writer
	// err manages errors that occur in all parts of the markdown building.
	err error
	// headers stores header information for table of contents generation.
	headers []headerInfo
	// tocOptions stores the table of contents generation options.
	tocOptions *TableOfContentsOptions
	// tocInserted indicates whether a table of contents placeholder has been generated.
	tocInserted bool
	// blockSpacing separates every block with a blank line.
	blockSpacing bool
}

// Option configures a Markdown at construction time.
type Option func(*Markdown)

// WithBlockSpacing separates every block with a blank line.
//
// The default output only inserts the blank lines markdown cannot do without,
// which keeps documents compact but leaves markdownlint complaining about
// headings, fenced blocks, and tables that touch their neighbors. Tools such
// as mkdocs are stricter than GitHub about this. Turn the option on when the
// document is going to be linted or rendered by something other than GitHub.
func WithBlockSpacing() Option {
	return func(m *Markdown) {
		m.blockSpacing = true
	}
}

// NewMarkdown returns new Markdown.
func NewMarkdown(w io.Writer, opts ...Option) *Markdown {
	m := &Markdown{
		body:    []string{},
		dest:    w,
		headers: []headerInfo{},
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// String returns markdown text.
//
// It returns the document built so far whether or not an error was recorded,
// and it does not need a writer, so it works on a builder constructed with nil.
func (m *Markdown) String() string {
	body := m.body

	if m.tocInserted && m.tocOptions != nil {
		if toc := m.generateTableOfContents(); len(toc) > 0 {
			body = insertTableOfContents(body, toc)
		}
	}

	return joinBlocks(body, m.blockSpacing)
}

// normalizeLineFeeds rewrites every line ending in text to the platform one.
// Text that reaches the builder from elsewhere, such as a table rendered by
// tablewriter, is separated by "\n" regardless of platform.
func normalizeLineFeeds(text string) string {
	return normalizeLineFeedsTo(text, internal.LineFeed())
}

// normalizeLineFeedsTo rewrites every line ending in text to lineFeed.
//
// The target line ending is a parameter rather than read from the platform so
// that both answers can be tested wherever the tests run. Reading the platform
// here left the Windows branch untested on Linux and the other way round.
func normalizeLineFeedsTo(text, lineFeed string) string {
	if lineFeed == "\n" {
		return strings.ReplaceAll(text, "\r\n", "\n")
	}
	return strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\n", lineFeed)
}

// joinBlocks joins the body, adding the blank line that markdown requires
// between certain blocks.
//
// A blockquote or alert swallows whatever follows it by lazy continuation, and
// a list swallows a table or a paragraph that starts on the next line. Both
// produce documents that render wrongly on GitHub while looking fine in the
// source, which is why callers of this package litter their code with manual
// spacer calls. Everything else is joined exactly as before.
func joinBlocks(body []string, always bool) string {
	lf := internal.LineFeed()

	var buf strings.Builder
	for i, block := range body {
		if i > 0 {
			buf.WriteString(lf)
			if needsBlankLine(body[i-1], block, always) {
				buf.WriteString(lf)
			}
		}
		buf.WriteString(block)
	}
	return buf.String()
}

// needsBlankLine reports whether a blank line has to separate two blocks.
func needsBlankLine(prev, next string, always bool) bool {
	// A whitespace-only entry, which is what LF() writes, already separates the
	// blocks; adding another blank line would just pile them up.
	if strings.TrimSpace(prev) == "" || strings.TrimSpace(next) == "" {
		return false
	}
	// Table and Details already end with a line feed, so the join produces the
	// blank line on its own.
	if strings.HasSuffix(prev, internal.LineFeed()) {
		return false
	}
	// An HTML comment renders as nothing and cannot absorb the block before it.
	// The table of contents markers are comments, so this keeps the generated
	// entries tucked against them.
	if strings.HasPrefix(next, "<!--") || strings.HasPrefix(prev, "<!--") {
		return false
	}

	prevList, nextList := listKind(prev), listKind(next)
	sameList := prevList != "" && prevList == nextList

	if always {
		// Consecutive items of one list still belong together; everything else
		// gets the blank line markdownlint expects.
		return !sameList
	}

	switch {
	case isQuoteBlock(prev):
		// Anything on the line after a quote is read as part of it.
		return true
	case prevList != "" && !sameList:
		// A different kind of list starts a new list, so it needs the blank line
		// as much as a paragraph or a table does.
		return true
	default:
		return false
	}
}

// isQuoteBlock reports whether the block is a blockquote or an alert.
func isQuoteBlock(block string) bool {
	return strings.HasPrefix(block, ">")
}

// listKind identifies which kind of list item a block is, or returns an empty
// string when it is not a list item at all.
//
// The kind matters because list items are appended one per entry: consecutive
// items of the same list must stay tight, while a bullet list followed by an
// ordered list is two lists and needs the blank line between them.
func listKind(block string) string {
	trimmed := strings.TrimLeft(block, " ")

	switch {
	case strings.HasPrefix(trimmed, "- [ ] "), strings.HasPrefix(trimmed, "- [x] "):
		return "checkbox"
	case strings.HasPrefix(trimmed, "- "), strings.HasPrefix(trimmed, "* "), strings.HasPrefix(trimmed, "+ "):
		return "bullet"
	}

	digits := 0
	for digits < len(trimmed) && trimmed[digits] >= '0' && trimmed[digits] <= '9' {
		digits++
	}
	if digits > 0 && strings.HasPrefix(trimmed[digits:], ". ") {
		return "ordered"
	}
	return ""
}

// insertTableOfContents places the generated entries between the two marker
// entries in the body.
//
// The markers are separate body entries, so this works on the slice rather than
// on the joined text. Matching the joined text meant matching the exact pair
// "<!-- BEGIN_TOC -->\n<!-- END_TOC -->", which silently stopped matching as
// soon as anything was placed between the markers: the replacement quietly did
// nothing and the document shipped with an empty table of contents.
func insertTableOfContents(body, toc []string) []string {
	begin := -1
	for i, line := range body {
		if line == TableOfContentsMarkerBegin {
			begin = i
			break
		}
	}
	if begin == -1 {
		return body
	}

	end := -1
	for i := begin + 1; i < len(body); i++ {
		if body[i] == TableOfContentsMarkerEnd {
			end = i
			break
		}
	}
	if end == -1 {
		return body
	}

	out := make([]string, 0, len(body)+len(toc))
	out = append(out, body[:begin+1]...)
	out = append(out, toc...)
	out = append(out, body[end:]...)
	return out
}

// Error returns the error the chain recorded, or nil.
//
// It is the same error [Markdown.Build] returns, for callers who would rather
// check before writing than after.
func (m *Markdown) Error() error {
	return m.err
}

// PlainText set plain text
func (m *Markdown) PlainText(text string) *Markdown {
	m.body = append(m.body, text)
	return m
}

// PlainTextf set plain text with format
func (m *Markdown) PlainTextf(format string, args ...interface{}) *Markdown {
	return m.PlainText(fmt.Sprintf(format, args...))
}

// Build writes markdown text to output destination.
//
// It returns the error the chain recorded, or nil. A nil destination and a
// destination that refuses the document are both reported rather than causing a
// panic, and either message carries the earlier error too when there is one.
//
// The document is written with a trailing line ending, so appending a second
// document to the same writer starts it on its own line.
//
// Build may be called more than once; each call writes the document again.
func (m *Markdown) Build() error {
	if m.dest == nil {
		if m.err != nil {
			return fmt.Errorf("failed to write markdown text: destination writer is nil: %s", m.err.Error()) //nolint:wrapcheck
		}
		return errors.New("failed to write markdown text: destination writer is nil")
	}

	// A document written to a file has to end with a newline: markdownlint MD047
	// requires it, and appending a second document to the same writer would
	// otherwise splice its first line onto the last line of this one.
	out := m.String()
	if !strings.HasSuffix(out, internal.LineFeed()) {
		out += internal.LineFeed()
	}

	if _, err := fmt.Fprint(m.dest, out); err != nil {
		if m.err != nil {
			return fmt.Errorf("failed to write markdown text: %w: %s", err, m.err.Error()) //nolint:wrapcheck
		}
		return fmt.Errorf("failed to write markdown text: %w", err)
	}
	return m.err
}

// H1 is markdown header.
// If you set text "Hello", it will be converted to "# Hello".
func (m *Markdown) H1(text string) *Markdown {
	m.headers = append(m.headers, headerInfo{level: TableOfContentsDepthH1, text: text})
	m.body = append(m.body, fmt.Sprintf("# %s", text))
	return m
}

// H1f is markdown header with format.
// If you set format "%s", text "Hello", it will be converted to "# Hello".
func (m *Markdown) H1f(format string, args ...interface{}) *Markdown {
	return m.H1(fmt.Sprintf(format, args...))
}

// H2 is markdown header.
// If you set text "Hello", it will be converted to "## Hello".
func (m *Markdown) H2(text string) *Markdown {
	m.headers = append(m.headers, headerInfo{level: TableOfContentsDepthH2, text: text})
	m.body = append(m.body, fmt.Sprintf("## %s", text))
	return m
}

// H2f is markdown header with format.
// If you set format "%s", text "Hello", it will be converted to "## Hello".
func (m *Markdown) H2f(format string, args ...interface{}) *Markdown {
	return m.H2(fmt.Sprintf(format, args...))
}

// H3 is markdown header.
// If you set text "Hello", it will be converted to "### Hello".
func (m *Markdown) H3(text string) *Markdown {
	m.headers = append(m.headers, headerInfo{level: TableOfContentsDepthH3, text: text})
	m.body = append(m.body, fmt.Sprintf("### %s", text))
	return m
}

// H3f is markdown header with format.
// If you set format "%s", text "Hello", it will be converted to "### Hello".
func (m *Markdown) H3f(format string, args ...interface{}) *Markdown {
	return m.H3(fmt.Sprintf(format, args...))
}

// H4 is markdown header.
// If you set text "Hello", it will be converted to "#### Hello".
func (m *Markdown) H4(text string) *Markdown {
	m.headers = append(m.headers, headerInfo{level: TableOfContentsDepthH4, text: text})
	m.body = append(m.body, fmt.Sprintf("#### %s", text))
	return m
}

// H4f is markdown header with format.
// If you set format "%s", text "Hello", it will be converted to "#### Hello".
func (m *Markdown) H4f(format string, args ...interface{}) *Markdown {
	return m.H4(fmt.Sprintf(format, args...))
}

// H5 is markdown header.
// If you set text "Hello", it will be converted to "##### Hello".
func (m *Markdown) H5(text string) *Markdown {
	m.headers = append(m.headers, headerInfo{level: TableOfContentsDepthH5, text: text})
	m.body = append(m.body, fmt.Sprintf("##### %s", text))
	return m
}

// H5f is markdown header with format.
// If you set format "%s", text "Hello", it will be converted to "##### Hello".
func (m *Markdown) H5f(format string, args ...interface{}) *Markdown {
	return m.H5(fmt.Sprintf(format, args...))
}

// H6 is markdown header.
// If you set text "Hello", it will be converted to "###### Hello".
func (m *Markdown) H6(text string) *Markdown {
	m.headers = append(m.headers, headerInfo{level: TableOfContentsDepthH6, text: text})
	m.body = append(m.body, fmt.Sprintf("###### %s", text))
	return m
}

// H6f is markdown header with format.
// If you set format "%s", text "Hello", it will be converted to "###### Hello".
func (m *Markdown) H6f(format string, args ...interface{}) *Markdown {
	return m.H6(fmt.Sprintf(format, args...))
}

// TableOfContents generates a table of contents placeholder that will be replaced when Build() is called.
// The table of contents will include all headers from H1 to the specified maxDepth.
// Only one table of contents can be generated per document.
//
// Example:
//
//	markdown.NewMarkdown(os.Stdout).
//	   H1("Title").
//	   TableOfContents(markdown.TableOfContentsDepthH3).  // Table of contents will be placed here
//	   H2("Section 1").
//	   H3("Subsection 1.1").
//	   Build()
func (m *Markdown) TableOfContents(maxDepth TableOfContentsDepth) *Markdown {
	return m.TableOfContentsWithRange(TableOfContentsDepthH1, maxDepth)
}

// TableOfContentsWithRange generates a table of contents placeholder with custom depth range.
// The table of contents will include headers from minDepth to maxDepth inclusive.
// Only one table of contents can be generated per document.
//
// Example:
//
//	markdown.NewMarkdown(os.Stdout).
//	   H1("Title").  // This H1 will not appear in table of contents
//	   H2("Table of Contents").
//	   TableOfContentsWithRange(markdown.TableOfContentsDepthH2, markdown.TableOfContentsDepthH4).  // Only include H2-H4 in table of contents
//	   H2("Section 1").
//	   H3("Subsection 1.1").
//	   H4("Detail").
//	   H5("Deep Detail").  // This H5 will not appear in table of contents
//	   Build()
func (m *Markdown) TableOfContentsWithRange(minDepth, maxDepth TableOfContentsDepth) *Markdown {
	if m.tocInserted {
		if m.err == nil {
			m.err = errors.New("table of contents has already been generated")
		}
		return m
	}

	if minDepth < TableOfContentsDepthH1 || minDepth > TableOfContentsDepthH6 {
		if m.err == nil {
			m.err = fmt.Errorf("invalid minDepth: %d (must be between 1 and 6)", minDepth)
		}
		return m
	}

	if maxDepth < TableOfContentsDepthH1 || maxDepth > TableOfContentsDepthH6 {
		if m.err == nil {
			m.err = fmt.Errorf("invalid maxDepth: %d (must be between 1 and 6)", maxDepth)
		}
		return m
	}

	if minDepth > maxDepth {
		if m.err == nil {
			m.err = fmt.Errorf("minDepth (%d) cannot be greater than maxDepth (%d)", minDepth, maxDepth)
		}
		return m
	}

	m.tocOptions = &TableOfContentsOptions{
		MinDepth: minDepth,
		MaxDepth: maxDepth,
	}
	m.tocInserted = true

	// Insert table of contents placeholder markers
	m.body = append(m.body, TableOfContentsMarkerBegin)
	m.body = append(m.body, TableOfContentsMarkerEnd)
	m.body = append(m.body, "")

	return m
}

// generateTableOfContents generates the table of contents based on collected headers and options.
func (m *Markdown) generateTableOfContents() []string {
	if m.tocOptions == nil || len(m.headers) == 0 {
		return []string{}
	}

	tocLines := make([]string, 0, len(m.headers))
	// Indent relative to the shallowest heading that actually appears, not to
	// the requested MinDepth. TableOfContents pins MinDepth at H1, so a document
	// that starts at H2 would otherwise have every entry indented two spaces,
	// producing a list nested under nothing.
	minIndent := int(m.tocOptions.MaxDepth)
	for _, header := range m.headers {
		if header.level < m.tocOptions.MinDepth || header.level > m.tocOptions.MaxDepth {
			continue
		}
		if int(header.level) < minIndent {
			minIndent = int(header.level)
		}
	}
	anchorCounts := make(map[string]int, len(m.headers))

	for _, header := range m.headers {
		// Skip headers outside the specified range
		if header.level < m.tocOptions.MinDepth || header.level > m.tocOptions.MaxDepth {
			continue
		}

		// Calculate relative indentation
		indent := strings.Repeat("  ", int(header.level)-minIndent)

		// Generate anchor following GitHub's convention
		baseAnchor := generateGitHubAnchor(header.text)
		count := anchorCounts[baseAnchor]
		anchor := baseAnchor
		if count > 0 {
			anchor = fmt.Sprintf("%s-%d", baseAnchor, count)
		}
		anchorCounts[baseAnchor] = count + 1

		tocLines = append(tocLines, fmt.Sprintf("%s- [%s](#%s)", indent, header.text, anchor))
	}

	return tocLines
}

func generateGitHubAnchor(text string) string {
	text = strings.ToLower(text)

	var b strings.Builder
	b.Grow(len(text))

	for _, r := range text {
		switch {
		case r == ' ' || r == '-':
			b.WriteRune('-')
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			b.WriteRune(r)
		}
	}

	return b.String()
}

// Details is markdown details.
//
// The body is surrounded by blank lines because an HTML block swallows
// everything up to the next blank line: without them the markdown inside
// <details> renders as literal text, and the block that follows </details>
// disappears into the same HTML block.
func (m *Markdown) Details(summary, text string) *Markdown {
	lf := internal.LineFeed()
	m.body = append(
		m.body,
		fmt.Sprintf("<details>%s<summary>%s</summary>%s%s%s%s%s</details>%s",
			lf, summary, lf, lf, text, lf, lf, lf))
	return m
}

// Detailsf is markdown details with format.
func (m *Markdown) Detailsf(summary, format string, args ...interface{}) *Markdown {
	return m.Details(summary, fmt.Sprintf(format, args...))
}

// BulletList is markdown bullet list.
// If you set text "Hello", it will be converted to "- Hello".
func (m *Markdown) BulletList(text ...string) *Markdown {
	for _, v := range text {
		m.body = append(m.body, fmt.Sprintf("- %s", v))
	}
	return m
}

// OrderedList is markdown number list.
// If you set text "Hello", it will be converted to "1. Hello".
func (m *Markdown) OrderedList(text ...string) *Markdown {
	for i, v := range text {
		m.body = append(m.body, fmt.Sprintf("%d. %s", i+1, v))
	}
	return m
}

// CheckBoxSet is markdown checkbox list.
type CheckBoxSet struct {
	// Checked is whether checked or not.
	Checked bool
	// Text is checkbox text.
	Text string
}

// CheckBox is markdown CheckBox.
func (m *Markdown) CheckBox(set []CheckBoxSet) *Markdown {
	for _, v := range set {
		if v.Checked {
			m.body = append(m.body, fmt.Sprintf("- [x] %s", v.Text))
		} else {
			m.body = append(m.body, fmt.Sprintf("- [ ] %s", v.Text))
		}
	}
	return m
}

// Blockquote is markdown blockquote.
// If you set text "Hello", it will be converted to "> Hello".
func (m *Markdown) Blockquote(text string) *Markdown {
	// Split on "\n" after dropping "\r": splitting on internal.LineFeed() meant
	// a plain Go literal containing "\n" was never split on Windows, and the
	// quote silently covered only its first line.
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("> %s", line)
	}
	// One entry per quote rather than one per line: the whole quote is a single
	// block, and the join has to be able to put a blank line after it without
	// cutting it in half.
	m.body = append(m.body, strings.Join(lines, internal.LineFeed()))
	return m
}

// CodeBlocks is code blocks.
// If you set text "Hello" and lang "go", it will be converted to
// "```go
// Hello
// ```".
//
// The block is fenced with three backticks, so content holding a line that
// starts with three backticks closes it early and everything after that line
// renders as prose. Markdown's answer is a longer fence, which this method does
// not write: a caller whose content can hold a fence has to build the block
// itself with PlainText.
func (m *Markdown) CodeBlocks(lang SyntaxHighlight, text string) *Markdown {
	m.body = append(m.body,
		fmt.Sprintf("```%s%s%s%s```", lang, internal.LineFeed(), text, internal.LineFeed()))
	return m
}

// HorizontalRule is markdown horizontal rule.
// It will be converted to "---".
func (m *Markdown) HorizontalRule() *Markdown {
	m.body = append(m.body, "---")
	return m
}

// TableAlignment represents column alignment in markdown tables.
type TableAlignment int

const (
	// AlignDefault represents no specific alignment (left by default).
	AlignDefault TableAlignment = iota
	// AlignLeft represents left alignment (:------).
	AlignLeft
	// AlignCenter represents center alignment (:-----:).
	AlignCenter
	// AlignRight represents right alignment (------:).
	AlignRight
)

// TableSet is markdown table.
type TableSet struct {
	// Header is table header.
	Header []string
	// Rows is table record.
	Rows [][]string
	// Alignment is column alignment for each column.
	// If nil or shorter than header length, remaining columns use AlignDefault.
	Alignment []TableAlignment
	// EscapeCells runs every header and row cell through EscapeTableCell.
	//
	// It is off by default because cells often hold markup the caller built on
	// purpose, with Link or Bold, and because callers who already escape their
	// own data would end up escaping it twice. Turn it on when the cells carry
	// arbitrary text that may contain a pipe or a newline.
	EscapeCells bool
}

// EscapeTableCell makes text safe to place in a table cell.
//
// A pipe ends the cell, so a pipe in the data silently splits it in two and
// drops the last column of the row; a newline ends the row outright. Neither
// produces an error, and ValidateColumns cannot see either, because the row
// still has the right length before it is serialized.
//
// The function is idempotent: a pipe that the caller already escaped is left
// alone, so passing text through it twice is harmless.
func EscapeTableCell(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			if i+1 < len(s) {
				switch s[i+1] {
				case '|', '\\':
					// Keep an existing escape intact, including the character it
					// escapes, so "\|" does not become "\\|".
					b.WriteByte(s[i])
					i++
					b.WriteByte(s[i])
					continue
				case '\n', '\r':
					// A backslash before a line break is a hard break. The break
					// itself is rewritten below, so drop the backslash rather than
					// letting it carry the line ending through untouched.
					continue
				}
			}
			b.WriteByte(s[i])
		case '|':
			b.WriteString(`\|`)
		case '\r':
			// Swallow the CR of a CRLF pair; a lone CR ends the line too.
			if i+1 < len(s) && s[i+1] == '\n' {
				continue
			}
			b.WriteString("<br>")
		case '\n':
			b.WriteString("<br>")
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// escaped returns the table with every cell escaped, leaving the original
// untouched so the caller's slices are not modified.
func (t TableSet) escaped() TableSet {
	header := make([]string, len(t.Header))
	for i, cell := range t.Header {
		header[i] = EscapeTableCell(cell)
	}

	rows := make([][]string, len(t.Rows))
	for i, row := range t.Rows {
		escapedRow := make([]string, len(row))
		for j, cell := range row {
			escapedRow[j] = EscapeTableCell(cell)
		}
		rows[i] = escapedRow
	}

	return TableSet{
		Header:      header,
		Rows:        rows,
		Alignment:   t.Alignment,
		EscapeCells: t.EscapeCells,
	}
}

// ValidateColumns checks if the number of columns in the header and records match.
func (t *TableSet) ValidateColumns() error {
	headerColumns := len(t.Header)
	for _, record := range t.Rows {
		if len(record) != headerColumns {
			return ErrMismatchColumn
		}
	}
	return nil
}

// Table is markdown table with alignment support.
func (m *Markdown) Table(t TableSet) *Markdown {
	if err := t.ValidateColumns(); err != nil {
		if m.err != nil {
			m.err = fmt.Errorf("failed to validate columns: %w: %s", err, m.err.Error()) //nolint:wrapcheck
		} else {
			m.err = fmt.Errorf("failed to validate columns: %w", err)
		}
		return m
	}

	if len(t.Header) == 0 {
		return m
	}

	if t.EscapeCells {
		t = t.escaped()
	}

	var buf strings.Builder

	// Write header row
	buf.WriteString("|")
	for _, header := range t.Header {
		buf.WriteString(" ")
		buf.WriteString(header)
		buf.WriteString(" |")
	}
	buf.WriteString(internal.LineFeed())

	// Write separator row with alignment
	buf.WriteString("|")
	for i := 0; i < len(t.Header); i++ {
		align := AlignDefault
		if i < len(t.Alignment) {
			align = t.Alignment[i]
		}

		switch align {
		case AlignDefault:
			buf.WriteString("---------|")
		case AlignLeft:
			buf.WriteString(":--------|")
		case AlignCenter:
			buf.WriteString(":-------:|")
		case AlignRight:
			buf.WriteString("--------:|")
		}
	}
	buf.WriteString(internal.LineFeed())

	// Write data rows
	for _, row := range t.Rows {
		buf.WriteString("|")
		for _, cell := range row {
			buf.WriteString(" ")
			buf.WriteString(cell)
			buf.WriteString(" |")
		}
		buf.WriteString(internal.LineFeed())
	}

	m.body = append(m.body, buf.String())
	return m
}

// TableOptions is markdown table options.
type TableOptions struct {
	// AutoWrapText is whether to wrap the text automatically.
	AutoWrapText bool
	// AutoFormatHeaders is whether to format the header automatically.
	AutoFormatHeaders bool
}

// CustomTable is markdown table. This is so not break the original Table function. with Possible breaking changes.
func (m *Markdown) CustomTable(t TableSet, options TableOptions) *Markdown {
	if err := t.ValidateColumns(); err != nil {
		// NOTE: If go version is 1.20, use errors.Join
		if m.err != nil {
			m.err = fmt.Errorf("failed to validate columns: %w: %s", err, m.err.Error()) //nolint:wrapcheck
		} else {
			m.err = fmt.Errorf("failed to validate columns: %w", err)
		}
	}

	if t.EscapeCells {
		t = t.escaped()
	}

	buf := &strings.Builder{}
	table := tablewriter.NewTable(
		buf,
		tablewriter.WithRenderer(
			renderer.NewBlueprint(
				tw.Rendition{
					Symbols: tw.NewSymbolCustom("Markdown").
						WithHeaderLeft("|").
						WithHeaderRight("|").
						WithColumn("|").
						WithMidLeft("|").
						WithMidRight("|").
						WithCenter("|"),
					Borders: tw.Border{
						Left:   tw.On,
						Top:    tw.Off,
						Right:  tw.On,
						Bottom: tw.Off,
					},
				},
			),
		),
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Formatting: tw.CellFormatting{
					AutoFormat: func() tw.State {
						if options.AutoFormatHeaders {
							return tw.Success
						}
						return tw.Fail
					}(),
				},
			},
			Row: tw.CellConfig{
				Formatting: tw.CellFormatting{
					AutoWrap: func() int {
						if options.AutoWrapText {
							return tw.WrapNormal
						}
						return tw.WrapNone
					}(),
					AutoFormat: func() tw.State {
						if options.AutoFormatHeaders {
							return tw.Success
						}
						return tw.Fail
					}(),
				},

				Alignment: tw.CellAlignment{Global: tw.AlignNone},
			},
		}),
	)

	table.Header(t.Header)
	if err := table.Bulk(t.Rows); err != nil {
		m.err = errors.Join(m.err, fmt.Errorf("failed to add rows to table: %w", err))
		return m
	}
	// This is so if the user wants to change the table settings they can
	if err := table.Render(); err != nil {
		m.err = errors.Join(m.err, fmt.Errorf("failed to render table: %w", err))
		return m
	}

	// tablewriter always separates rows with "\n", so on Windows its output
	// would otherwise mix line endings with the rest of the document, and the
	// delimiter row rewrite below would fail to find any rows to work on.
	rendered := normalizeLineFeeds(buf.String())

	m.body = append(m.body, applyAlignmentToDelimiterRow(rendered, t.Alignment))
	return m
}

// applyAlignmentToDelimiterRow rewrites the delimiter row of a rendered table so
// it carries the alignment colons markdown uses.
//
// tablewriter aligns by padding cells, which says nothing to a markdown reader:
// alignment lives in the second row, as :--- / :---: / ---:. Without this,
// TableSet.Alignment is silently dropped by CustomTable while Table honors it.
// Each column keeps its rendered width so the source stays visually aligned.
func applyAlignmentToDelimiterRow(rendered string, alignment []TableAlignment) string {
	if len(alignment) == 0 {
		return rendered
	}

	lineFeed := internal.LineFeed()
	lines := strings.Split(rendered, lineFeed)
	const delimiterRowIndex = 1
	if len(lines) <= delimiterRowIndex {
		return rendered
	}

	// The delimiter row is "|-----|------|"; splitting on "|" yields an empty
	// field at each end.
	// Splitting "|---|---|" on "|" leaves an empty field at each end, so a real
	// delimiter row always yields at least two fields.
	const minDelimiterFields = 2
	columns := strings.Split(lines[delimiterRowIndex], "|")
	if len(columns) < minDelimiterFields {
		return rendered
	}

	for i := 1; i < len(columns)-1; i++ {
		if strings.TrimLeft(columns[i], "-") != "" {
			return rendered // not a delimiter row after all; leave it alone
		}
		align := AlignDefault
		if i-1 < len(alignment) {
			align = alignment[i-1]
		}
		columns[i] = delimiterCell(len(columns[i]), align)
	}

	lines[delimiterRowIndex] = strings.Join(columns, "|")
	return strings.Join(lines, lineFeed)
}

const (
	// oneSidedMarkerWidth is the narrowest cell that fits ":-" or "-:".
	oneSidedMarkerWidth = 2
	// twoSidedMarkerWidth is the narrowest cell that fits ":-:".
	twoSidedMarkerWidth = 3
)

// delimiterCell renders one delimiter cell of the given width for the alignment.
// The width is preserved so the rendered table keeps its columns lined up. A
// column too narrow to hold its marker falls back to a plain rule rather than
// widening the table.
func delimiterCell(width int, align TableAlignment) string {
	switch align {
	case AlignLeft:
		if width < oneSidedMarkerWidth {
			return strings.Repeat("-", width)
		}
		return ":" + strings.Repeat("-", width-1)
	case AlignRight:
		if width < oneSidedMarkerWidth {
			return strings.Repeat("-", width)
		}
		return strings.Repeat("-", width-1) + ":"
	case AlignCenter:
		if width < twoSidedMarkerWidth {
			return strings.Repeat("-", width)
		}
		return ":" + strings.Repeat("-", width-twoSidedMarkerWidth+1) + ":"
	case AlignDefault:
		return strings.Repeat("-", width)
	}
	return strings.Repeat("-", width)
}

// LF is line feed.
//
// It writes a line holding two spaces, which is a hard line break marker. It
// also happens to separate blocks, which is how most callers use it. Use
// BlankLine when a blank line is what you mean.
func (m *Markdown) LF() *Markdown {
	m.body = append(m.body, "  ")
	return m
}

// BlankLine writes an empty line between two blocks.
func (m *Markdown) BlankLine() *Markdown {
	m.body = append(m.body, "")
	return m
}
