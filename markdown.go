// Package markdown is markdown builder that includes to convert Markdown to HTML.
package markdown

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/olekukonko/tablewriter"
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

// Markdown is markdown text.
type Markdown struct {
	// body is markdown body.
	body []string
	// dest is output destination for markdown body.
	dest io.Writer
	// err manages errors that occur in all parts of the markdown building.
	err error
}

// NewMarkdown returns new Markdown.
func NewMarkdown(w io.Writer) *Markdown {
	return &Markdown{
		body: []string{},
		dest: w,
	}
}

// String returns markdown text.
func (m *Markdown) String() string {
	return strings.Join(m.body, lineFeed())
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
	m.body = append(m.body, fmt.Sprintf("###### %s", text))
	return m
}

// H6f is markdown header with format.
// If you set format "%s", text "Hello", it will be converted to "###### Hello".
func (m *Markdown) H6f(format string, args ...interface{}) *Markdown {
	return m.H6(fmt.Sprintf(format, args...))
}

// Details is markdown details.
func (m *Markdown) Details(summary, text string) *Markdown {
	m.body = append(
		m.body,
		fmt.Sprintf("<details><summary>%s</summary>%s%s%s</details>", summary, lineFeed(), text, lineFeed()))
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
	lines := strings.Split(text, lineFeed())
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
		fmt.Sprintf("```%s%s%s%s```", lang, lineFeed(), text, lineFeed()))
	return m
}

// HorizontalRule is markdown horizontal rule.
// It will be converted to "---".
func (m *Markdown) HorizontalRule() *Markdown {
	m.body = append(m.body, "---")
	return m
}

// TableSet is markdown table.
type TableSet struct {
	// Header is table header.
	Header []string
	// Rows is table record.
	Rows [][]string
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

// Table is markdown table.
func (m *Markdown) Table(t TableSet) *Markdown {
	if err := t.ValidateColumns(); err != nil {
		// NOTE: If go version is 1.20, use errors.Join
		if m.err != nil {
			m.err = fmt.Errorf("failed to validate columns: %w: %s", err, m.err) //nolint:wrapcheck
		} else {
			m.err = fmt.Errorf("failed to validate columns: %w", err)
		}
	}

	buf := &strings.Builder{}
	table := tablewriter.NewWriter(buf)
	table.SetNewLine(lineFeed())
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetHeader(t.Header)
	for _, v := range t.Rows {
		table.Append(v)
	}
	table.Render()

	m.body = append(m.body, buf.String())
	return m
}

type TableOptions struct {
	// AutoWrapText is whether to wrap the text automatically.
	AutoWrapText bool
}

// CustomTable is markdown table. This is so not break the original Table function. with Possible breaking changes.
func (m *Markdown) CustomTable(t TableSet, options TableOptions) *Markdown {
	if err := t.ValidateColumns(); err != nil {
		// NOTE: If go version is 1.20, use errors.Join
		if m.err != nil {
			m.err = fmt.Errorf("failed to validate columns: %w: %s", err, m.err) //nolint:wrapcheck
		} else {
			m.err = fmt.Errorf("failed to validate columns: %w", err)
		}
	}

	buf := &strings.Builder{}
	table := tablewriter.NewWriter(buf)
	table.SetNewLine(lineFeed())
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(options.AutoWrapText)

	table.SetHeader(t.Header)
	for _, v := range t.Rows {
		table.Append(v)
	}
	// This is so if the user wants to change the table settings they can
	table.Render()

	m.body = append(m.body, buf.String())
	return m
}

// LF is line feed.
func (m *Markdown) LF() *Markdown {
	m.body = append(m.body, "  ")
	return m
}

// lineFeed return line feed for current OS.
func lineFeed() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
