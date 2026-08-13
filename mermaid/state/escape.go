package state

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// statementUnsafe is what a state description or a transition label cannot
	// carry.
	//
	// A semicolon ends the statement it is written in. What follows is then
	// parsed as a statement of its own, so whether the diagram survives depends
	// on what happens to come next: "a;b" draws, because "b" reads as a state
	// id, and "a;b{c" loses the whole diagram because "{" does not read as
	// anything. A description is prose and cannot be relied on to end in a word.
	statementUnsafe = ";"
	// noteUnsafe is what a one line note cannot carry. "note right of s1 : text"
	// is the one construct here that reads a colon as syntax, and a second one
	// ends the note.
	noteUnsafe = ";:"
)

// escapeStatement returns text ready to be written into a construct that cannot
// carry the characters in unsafe.
//
// Each is written as the entity form mermaid decodes, found by rendering them
// one at a time. A "#" that would start an entity is escaped with them, or text
// holding "#59;" and text holding a semicolon would draw the same diagram; a
// "#" anywhere else is ordinary text and comes out unchanged.
//
// A raw line break splits the statement the same way a semicolon does, and the
// second line drew as a stray state; it is written as "<br/>", the line break
// mermaid draws in these constructs.
//
// The lines of a multi line note need none of this: mermaid reads each as text
// until "end note" and takes every character probed, colon and semicolon
// included.
func escapeStatement(text, unsafe string) string {
	text = internal.LineBreaksToBr(text)
	if !strings.ContainsAny(text, unsafe+"#") {
		return text
	}

	var b strings.Builder
	b.Grow(len(text))
	for i := 0; i < len(text); i++ {
		switch {
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
