package piechart

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

// escapeLabel returns label ready to be written inside a pie chart's quoted
// slice label.
//
// A double quote ends the label, and mermaid then refuses the whole chart: the
// reader gets an error box instead of a picture. It is written as mermaid's own
// entity form instead. A backslash happens to work here too, unlike in a
// flowchart, but the entity form is what every other builder in this library
// writes and what mermaid documents.
//
// A "#" that would start an entity is escaped for the reason the escape exists
// at all: mermaid reads "#quot;" as a quotation mark, so a label holding one
// has to be spelled differently from a label holding the mark itself. A "#"
// anywhere else is ordinary text and comes out unchanged.
//
// A raw line break ends the quoted label early and the lexer refuses the whole
// chart, so it is written as "<br/>", the line break mermaid draws inside one.
func escapeLabel(label string) string {
	label = internal.LineBreaksToBr(label)
	if !strings.ContainsAny(label, `"#`) {
		return label
	}

	var b strings.Builder
	b.Grow(len(label))
	for i := 0; i < len(label); i++ {
		switch {
		case label[i] == '"':
			b.WriteString(internal.EntityEscape('"'))
		case label[i] == '#' && internal.StartsEntity(label[i+1:]):
			b.WriteString(internal.EntityEscape('#'))
		default:
			b.WriteByte(label[i])
		}
	}
	return b.String()
}

// escapeTitle returns title ready to be written after the "title" keyword.
//
// The title is the one place a pie chart takes unquoted text, so a quotation
// mark needs nothing here. What it cannot take is "%%", which starts a mermaid
// comment: everything from there to the end of the line is dropped, and the
// title is silently cut short rather than the chart failing. Each "%" in such a
// run is written as its entity so the pair never forms.
//
// A lone "%" is left alone, and so is a "#" that starts nothing. Both already
// reach the drawing intact, and their output is pinned by the golden files.
//
// A "<" that does not open a "<br/>" is eaten by the sanitizer, and a raw line
// break splits the statement and loses the chart; both entity forms were
// measured to decode in the drawing, the latter into a real line break.
func escapeTitle(title string) string {
	title = strings.ReplaceAll(title, "\r\n", "\n")
	title = strings.ReplaceAll(title, "\r", "\n")
	if !strings.ContainsAny(title, "%#<\n") {
		return title
	}

	var b strings.Builder
	b.Grow(len(title))
	for i := 0; i < len(title); i++ {
		switch {
		case title[i] == '%' && adjacentPercent(title, i):
			b.WriteString(internal.EntityEscape('%'))
		case title[i] == '<' && !strings.HasPrefix(title[i:], "<br/>"):
			b.WriteString(internal.EntityEscape('<'))
		case title[i] == '\n':
			b.WriteString(internal.EntityEscape('\n'))
		case title[i] == '#' && internal.StartsEntity(title[i+1:]):
			b.WriteString(internal.EntityEscape('#'))
		default:
			b.WriteByte(title[i])
		}
	}
	return b.String()
}

// adjacentPercent reports whether the "%" at index i has another beside it, and
// so is part of the run that would open a comment.
func adjacentPercent(title string, i int) bool {
	return (i > 0 && title[i-1] == '%') || (i+1 < len(title) && title[i+1] == '%')
}
