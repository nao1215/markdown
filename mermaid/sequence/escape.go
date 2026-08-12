package sequence

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// textUnsafe is what a message, a note or a block description cannot carry.
	//
	// A semicolon ends the statement and loses the whole diagram. A "#" opens a
	// comment, and the text is then quietly cut short at it: "deploy #2 of 3"
	// reached the drawing as "deploy" and nothing said so.
	textUnsafe = "#;"
	// participantUnsafe is what a participant's name cannot carry.
	//
	// A semicolon ends the statement here too. A colon separates a name from
	// the message that follows it, a comma separates one name from the next,
	// and an angle bracket opens the arrow syntax. A hyphen is the start of
	// every arrow, so "web-server" could be declared but never sent a message.
	// A "%%" opens a comment on the message line and drops the message that
	// followed it without a word. A parenthesis is safe on its own and a pair
	// is not, so both are written out: an unmatched one is not a name anybody
	// writes, and the matched pair is what a name like "Deploy (prod)" has. A
	// "#" is safe in a name, unlike in the text above.
	participantUnsafe = ";:,<>-%()"
	// noteParticipantUnsafe is what the participant a note is placed over
	// cannot carry.
	//
	// It is participantUnsafe without the comma, because a note may be placed
	// over two participants at once and the comma between them is the syntax
	// that says so. A caller passing "Alice,Bob" there means both, and this
	// package has no way to tell that apart from one name holding a comma.
	noteParticipantUnsafe = ";:<>-%()"
)

// escape returns text ready to be written into a construct that cannot carry
// the characters in unsafe.
//
// A sequence diagram takes no quoted text at all, so each character is written
// as the entity form mermaid decodes, found by rendering them one at a time. A
// "#" that would start an entity is escaped whether or not the construct lists
// it, because this package writes entities and a caller's literal "#59;" has to
// come out different from a caller's semicolon. A "#" that starts nothing is
// left alone wherever it reaches the drawing, which is what keeps output that
// already renders unchanged.
func escape(text, unsafe string) string {
	if !strings.ContainsAny(text, unsafe+"#") {
		return text
	}

	var b strings.Builder
	b.Grow(len(text))
	for i := 0; i < len(text); i++ {
		switch {
		case text[i] == '%':
			// Only a run opens a comment. A lone "%" reaches every drawing
			// here, so "50% of traffic" comes out as it always has.
			if strings.IndexByte(unsafe, '%') >= 0 && adjacentPercent(text, i) {
				b.WriteString(internal.EntityEscape('%'))
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

// adjacentPercent reports whether the "%" at index i has another beside it, and
// so is part of the run that would open a comment.
func adjacentPercent(text string, i int) bool {
	return (i > 0 && text[i-1] == '%') || (i+1 < len(text) && text[i+1] == '%')
}

// escapeText returns message, note or description text ready to be written.
func escapeText(text string) string {
	return escape(text, textUnsafe)
}

// escapeParticipant returns a participant's name ready to be written.
func escapeParticipant(name string) string {
	return escape(name, participantUnsafe)
}

// escapeNoteParticipant returns the participant a note is placed over, ready to
// be written.
func escapeNoteParticipant(name string) string {
	return escape(name, noteParticipantUnsafe)
}

// escapeParticipants returns each name ready to be written.
func escapeParticipants(names []string) []string {
	escaped := make([]string, 0, len(names))
	for _, name := range names {
		escaped = append(escaped, escapeParticipant(name))
	}
	return escaped
}
