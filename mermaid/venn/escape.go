package venn

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

// escapeLabel returns label ready to be written inside the quoted brackets a
// set's label goes in.
//
// A double quote ends the label and mermaid refuses the whole diagram, so it is
// written as the entity form mermaid decodes. A "#" that would start an entity
// is escaped with it, or a label holding "#quot;" and a label holding a
// quotation mark would draw the same picture; a "#" anywhere else reaches the
// drawing and is left alone.
func escapeLabel(label string) string {
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
// The title is the one place a Venn diagram takes unquoted text, so a quotation
// mark needs nothing here. What the unquoted lexer refuses is "#", which starts
// a comment, and ";", which ends the statement; both are written as the entity
// form mermaid decodes back into them. A "<" that does not open a "<br/>" is
// eaten by the renderer's sanitizer instead of the lexer, and gets the same
// treatment. Everything else probed reaches the drawing, including an emoji
// and Japanese.
func escapeTitle(title string) string {
	if !strings.ContainsAny(title, "#;<") {
		return title
	}

	var b strings.Builder
	b.Grow(len(title))
	for i, r := range title {
		switch {
		case r == '#':
			b.WriteString(internal.EntityEscape('#'))
		case r == ';':
			b.WriteString(internal.EntityEscape(';'))
		case r == '<' && !strings.HasPrefix(title[i:], "<br/>"):
			b.WriteString(internal.EntityEscape('<'))
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
