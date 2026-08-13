package timeline

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// textUnsafe is what a period or an event cannot carry.
	//
	// A colon separates a period from its events, so a literal one splits the
	// text in two. A "#" opens a comment and the rest of the line is quietly
	// dropped: "deploy #2" reached the drawing as "deploy". A "%%" run opens a
	// comment of its own and loses the whole line. Each was measured by
	// rendering, and the entity forms were measured to decode back to the
	// character they stand for.
	textUnsafe = ":#%"
	// sectionUnsafe is what a section name cannot carry.
	//
	// A literal colon in a section name is a parse error that loses the whole
	// diagram. The "#" and "%%" a period cannot carry both reach the drawing
	// from a section name and are left alone.
	sectionUnsafe = ":"
)

// escape returns text ready to be written into a construct that cannot carry
// the characters in unsafe.
//
// A timeline takes no quoted text at all, so each character is written as the
// entity form mermaid decodes, found by rendering them one at a time. A "#"
// that would start an entity is escaped whether or not the construct lists it,
// because this package writes entities and a caller's literal "#58;" has to
// come out different from a caller's colon. A "#" that starts nothing is left
// alone wherever it reaches the drawing.
func escape(text, unsafe string) string {
	if !strings.ContainsAny(text, unsafe+"#") {
		return text
	}

	var b strings.Builder
	b.Grow(len(text))
	for i := 0; i < len(text); i++ {
		switch {
		case text[i] == '%':
			// Only a run opens a comment, so "50% of traffic" is untouched.
			if strings.IndexByte(unsafe, '%') >= 0 && adjacentPercent(text, i) {
				b.WriteString(internal.EntityEscape('%'))
				continue
			}
			b.WriteByte(text[i])
		case text[i] == '#' && strings.IndexByte(unsafe, '#') >= 0:
			b.WriteString(internal.EntityEscape('#'))
		case strings.IndexByte(unsafe, text[i]) >= 0:
			b.WriteString(internal.EntityEscape(rune(text[i])))
		case text[i] == '#' && internal.StartsEntity(text[i+1:]):
			b.WriteString(internal.EntityEscape('#'))
		default:
			b.WriteByte(text[i])
		}
	}
	return b.String()
}

// adjacentPercent reports whether the "%" at index i has another beside it, and
// so is part of the run that would open a comment.
func adjacentPercent(text string, i int) bool {
	return (i > 0 && text[i-1] == '%') || (i+1 < len(text) && text[i+1] == '%')
}

// escapeText returns a period or an event ready to be written.
func escapeText(text string) string {
	return escape(text, textUnsafe)
}

// escapeSection returns a section name ready to be written.
func escapeSection(name string) string {
	return escape(name, sectionUnsafe)
}
