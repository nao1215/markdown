package wardley

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

// escapeTitle returns title ready to be written after the "title" keyword.
//
// The title is the rest of the line, and almost everything reaches the drawing
// as it was written: a quotation mark, a hash, a semicolon, an emoji and
// Japanese text all survive, which was found by rendering them one at a time.
// What does not is "%%", which opens a mermaid comment, so the rest of the
// title is dropped and the map still draws, saying less than it was asked to.
//
// Each "%" of such a run is written as the entity form mermaid decodes back
// into it. A lone "%" is left alone, so "50% done" comes out as it was given,
// and a "#" is escaped only where it would otherwise start an entity, which
// keeps "PR #12" unchanged too.
func escapeTitle(title string) string {
	if !strings.ContainsAny(title, "%#") {
		return title
	}

	var b strings.Builder
	b.Grow(len(title))
	for i := 0; i < len(title); i++ {
		switch {
		case title[i] == '%' && adjacentPercent(title, i):
			b.WriteString(internal.EntityEscape('%'))
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
