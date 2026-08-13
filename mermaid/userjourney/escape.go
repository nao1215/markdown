package userjourney

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

// A user journey is written entirely in unquoted text: a title is the rest of
// its line, a section is the rest of its line, and a task is
// "name: score: actor, actor". There is nowhere to put a quotation mark, so
// each field's own punctuation has to be written as the entity form mermaid
// decodes. What that punctuation is differs by field, and each set below was
// measured by rendering one character at a time.
const (
	// titleUnsafe is what a title cannot carry. A semicolon ends the statement
	// and loses the whole diagram; a "#" starts a comment, and the title is
	// then quietly cut short at it rather than the diagram failing. A "<" that
	// does not open a "<br/>" passes the parser and is then eaten by the
	// renderer's sanitizer, taking the rest of the title with it. A raw line
	// break splits the statement — the other fields reject one, but a title
	// carries it as the entity a title decodes into a real line break.
	titleUnsafe = "#;<\n"
	// sectionUnsafe is what a section name cannot carry. A colon loses the
	// diagram here, unlike in a title, and a "#" is safe here, unlike in one.
	sectionUnsafe = ";:"
	// taskUnsafe is what a task name cannot carry. A "#" starts a comment and
	// loses the diagram, and a colon ends the name: the diagram still draws
	// then, with the rest of the name read as the score.
	taskUnsafe = "#;:"
	// actorUnsafe is what one actor's name cannot carry. A comma separates one
	// actor from the next, so a comma inside a name quietly splits it in two.
	actorUnsafe = ";:,"
)

// escapeField returns text ready to be written into a field that cannot carry
// the characters in unsafe.
//
// A "#" is escaped whenever it would start an entity, whether or not the field
// lists it: this package writes entities, so a caller's literal "#59;" has to
// come out different from a caller's semicolon. A "#" that starts nothing is
// left alone in the fields that can hold one, which is what keeps output that
// already renders unchanged.
func escapeField(text, unsafe string) string {
	// A CRLF pair and a lone CR are line breaks as much as a lone LF, and the
	// byte loop below knows only the LF; the fields that reject a line break
	// have already done so by the time this runs.
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	if !strings.ContainsAny(text, unsafe+"#") {
		return text
	}

	var b strings.Builder
	b.Grow(len(text))
	for i := 0; i < len(text); i++ {
		switch {
		case text[i] == '<':
			// Only the bare "<" is eaten; "<br/>" is left as the line break
			// the renderer draws.
			if strings.IndexByte(unsafe, '<') >= 0 && !strings.HasPrefix(text[i:], "<br/>") {
				b.WriteString(internal.EntityEscape('<'))
				continue
			}
			b.WriteByte(text[i])
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

// escapeActors returns each actor ready to be written into the comma separated
// list a task ends with.
func escapeActors(actors []string) []string {
	escaped := make([]string, 0, len(actors))
	for _, actor := range actors {
		escaped = append(escaped, escapeField(actor, actorUnsafe))
	}
	return escaped
}
