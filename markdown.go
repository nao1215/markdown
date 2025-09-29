// Package markdown is markdown builder that includes to convert Markdown to HTML.
package markdown

import (
	"errors"
	"fmt"
	"io"
	"strings"

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
}

// NewMarkdown returns new Markdown.
func NewMarkdown(w io.Writer) *Markdown {
	return &Markdown{
		body:    []string{},
		dest:    w,
		headers: []headerInfo{},
	}
}

// String returns markdown text.
func (m *Markdown) String() string {
	content := strings.Join(m.body, internal.LineFeed())

	// Replace table of contents placeholders with actual table of contents content if present
	if m.tocInserted && m.tocOptions != nil {
		tocContent := m.generateTableOfContents()
		if len(tocContent) > 0 {
			tocText := strings.Join(tocContent, internal.LineFeed())
			placeholder := TableOfContentsMarkerBegin + internal.LineFeed() + TableOfContentsMarkerEnd
			replacement := TableOfContentsMarkerBegin + internal.LineFeed() + tocText + internal.LineFeed() + TableOfContentsMarkerEnd
			content = strings.ReplaceAll(content, placeholder, replacement)
		}
	}

	return content
}

// Error returns error.
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
func (m *Markdown) Build() error {
	if _, err := fmt.Fprint(m.dest, m.String()); err != nil {
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
	minIndent := int(m.tocOptions.MinDepth)

	for _, header := range m.headers {
		// Skip headers outside the specified range
		if header.level < m.tocOptions.MinDepth || header.level > m.tocOptions.MaxDepth {
			continue
		}

		// Calculate relative indentation
		indent := strings.Repeat("  ", int(header.level)-minIndent)

		// Generate anchor following GitHub's convention
		anchor := strings.ToLower(strings.ReplaceAll(header.text, " ", "-"))
		anchor = strings.Map(func(r rune) rune {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
				return r
			}
			return -1
		}, anchor)

		tocLines = append(tocLines, fmt.Sprintf("%s- [%s](#%s)", indent, header.text, anchor))
	}

	return tocLines
}

// Details is markdown details.
func (m *Markdown) Details(summary, text string) *Markdown {
	m.body = append(
		m.body,
		fmt.Sprintf("<details><summary>%s</summary>%s%s%s</details>",
			summary, internal.LineFeed(), text, internal.LineFeed()))
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
	lines := strings.Split(text, internal.LineFeed())
	for _, line := range lines {
		m.body = append(m.body, fmt.Sprintf("> %s", line))
	}
	return m
}

// CodeBlocks is code blocks.
// If you set text "Hello" and lang "go", it will be converted to
// "```go
// Hello
// ```".
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
			m.err = fmt.Errorf("failed to validate columns: %w: %s", err, m.err) //nolint:wrapcheck
		} else {
			m.err = fmt.Errorf("failed to validate columns: %w", err)
		}
		return m
	}

	if len(t.Header) == 0 {
		return m
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
		case AlignLeft:
			buf.WriteString(":--------|")
		case AlignCenter:
			buf.WriteString(":-------:|")
		case AlignRight:
			buf.WriteString("--------:|")
		default: // AlignDefault
			buf.WriteString("---------|")
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
			m.err = fmt.Errorf("failed to validate columns: %s: %s", err, m.err) //nolint:wrapcheck
		} else {
			m.err = fmt.Errorf("failed to validate columns: %s", err)
		}
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

	m.body = append(m.body, buf.String())
	return m
}

// LF is line feed.
func (m *Markdown) LF() *Markdown {
	m.body = append(m.body, "  ")
	return m
}
