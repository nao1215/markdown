package markdown

import (
	"fmt"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// Link return text with link format.
// If you set text "Hello" and url "https://example.com",
// it will be converted to "[Hello](https://example.com)".
func Link(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
}

// FootnoteReference returns text with footnote reference format.
// If you set id "1", it will be converted to "[^1]".
func FootnoteReference(id string) string {
	return fmt.Sprintf("[^%s]", id)
}

// FootnoteDefinition returns text with footnote definition format.
// If you set id "1" and text "Hello", it will be converted to "[^1]: Hello".
func FootnoteDefinition(id, text string) string {
	return fmt.Sprintf("[^%s]: %s", id, text)
}

// ReferenceLink returns text with reference link format.
// If you set text "Go" and id "go-site", it will be converted to "[Go][go-site]".
func ReferenceLink(text, id string) string {
	return fmt.Sprintf("[%s][%s]", text, id)
}

// ReferenceLinkDefinition returns text with reference link definition format.
// If you set id "go-site" and url "https://golang.org",
// it will be converted to "[go-site]: https://golang.org".
// If title is set, it will be converted to
// "[go-site]: https://golang.org \"The Go Programming Language\"".
func ReferenceLinkDefinition(id, url string, title ...string) string {
	if len(title) == 0 || title[0] == "" {
		return fmt.Sprintf("[%s]: %s", id, url)
	}

	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`)
	escapedTitle := r.Replace(title[0])
	return fmt.Sprintf("[%s]: %s \"%s\"", id, url, escapedTitle)
}

// escapeMathExpression escapes '$' in math expression for safe inline usage.
func escapeMathExpression(expression string) string {
	return strings.ReplaceAll(expression, "$", "\\$")
}

// InlineMath returns text with inline mathematical expression format.
// It calls escapeMathExpression, so '$' in expression is escaped as '\$'.
// If you set expression "E=mc^2", it will be converted to "$E=mc^2$".
func InlineMath(expression string) string {
	return fmt.Sprintf("$%s$", escapeMathExpression(expression))
}

// BlockMath returns text with block mathematical expression format.
// BlockMath does not escape expression; it writes the raw expression between '$$' delimiters.
// If you set expression "x^2 + y^2 = z^2", it will be converted to:
//
//	$$
//	x^2 + y^2 = z^2
//	$$
func BlockMath(expression string) string {
	lf := internal.LineFeed()
	return fmt.Sprintf("$$%s%s%s$$", lf, expression, lf)
}

// Image return text with image format.
// If you set text "Hello" and url "https://example.com/image.png",
// it will be converted to "![Hello](https://example.com/image.png)".
func Image(text, url string) string {
	return fmt.Sprintf("![%s](%s)", text, url)
}

// Strikethrough return text with strikethrough format.
// If you set text "Hello", it will be converted to "~~Hello~~".
func Strikethrough(text string) string {
	return fmt.Sprintf("~~%s~~", text)
}

// Bold return text with bold format.
// If you set text "Hello", it will be converted to "**Hello**".
func Bold(text string) string {
	return fmt.Sprintf("**%s**", text)
}

// Italic return text with italic format.
// If you set text "Hello", it will be converted to "*Hello*".
func Italic(text string) string {
	return fmt.Sprintf("*%s*", text)
}

// BoldItalic return text with bold and italic format.
// If you set text "Hello", it will be converted to "***Hello***".
func BoldItalic(text string) string {
	return fmt.Sprintf("***%s***", text)
}

// Code return text with code format.
// If you set text "Hello", it will be converted to "`Hello`".
func Code(text string) string {
	return fmt.Sprintf("`%s`", text)
}

// Highlight return text with highlight format.
// If you set text "Hello", it will be converted to "==Hello==".
func Highlight(text string) string {
	return fmt.Sprintf("==%s==", text)
}
