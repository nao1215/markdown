package class

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

// escapeAfterColon returns text ready to be written after the colon in the
// "A --> B : label" and "A : member" forms.
//
// Those two lines are the only unquoted text a class diagram takes, and mermaid
// reads a colon or a semicolon in them as the end of the statement. Either one
// makes it refuse the whole diagram, so a relationship labeled "owns: many"
// left the reader an error box instead of a picture.
//
// Both are written as the entity form mermaid decodes, which was found by
// rendering. A "#" that would start an entity is escaped with them, or a label
// holding "#58;" and a label holding a colon would draw the same diagram; a "#"
// anywhere else is ordinary text and comes out unchanged.
//
// Everywhere else in this package is already quoted, and a class body member, a
// class label and a note each take both characters as they are. Nothing here
// touches those.
func escapeAfterColon(text string) string {
	if !strings.ContainsAny(text, ":;#") {
		return text
	}

	var b strings.Builder
	b.Grow(len(text))
	for i := 0; i < len(text); i++ {
		switch {
		case text[i] == ':':
			b.WriteString(internal.EntityEscape(':'))
		case text[i] == ';':
			b.WriteString(internal.EntityEscape(';'))
		case text[i] == '#' && internal.StartsEntity(text[i+1:]):
			b.WriteString(internal.EntityEscape('#'))
		default:
			b.WriteByte(text[i])
		}
	}
	return b.String()
}
