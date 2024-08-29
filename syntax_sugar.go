package markdown

import "fmt"

// Link return text with link format.
// If you set text "Hello" and url "https://example.com",
// it will be converted to "[Hello](https://example.com)".
func Link(text, url string) string {
	return fmt.Sprintf("[%s](%s)", text, url)
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
