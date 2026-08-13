package kanban

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

// labelUnsafe is what a card label cannot carry.
//
// A card is written "id[label]", and mermaid ends the label at the first "]".
// A parenthesis or a closing brace ends it too, because the shapes those spell
// elsewhere in mermaid are recognized here. An opening bracket and an opening
// brace are not among them: mermaid takes either inside a label, and escaping
// them would change output that already reaches the drawing. A line break is
// not among them either, because validateCardLabel rejects one before the
// escaping runs.
const labelUnsafe = "])(}"

// metadataUnsafe is what a metadata value cannot carry, beyond the quoting the
// YAML scalar does for itself.
//
// Task metadata is written "@{ ticket: 'value' }", which mermaid reads as YAML
// and then draws. A closing brace ends the block, and a double quote or a
// caret is refused by the kanban lexer before YAML ever sees the line.
const metadataUnsafe = `"}^`

// escapeLabel returns label ready to be written between the brackets of a card.
//
// Each character above is written as the entity form mermaid decodes, found by
// rendering them one at a time. A "#" that would start an entity is escaped
// with them, or a label holding "#93;" and a label holding a closing bracket
// would draw the same card; a "#" anywhere else comes out unchanged.
func escapeLabel(label string) string {
	return escapeEntities(label, labelUnsafe)
}

// escapeMetadata returns value as the single quoted YAML scalar a task's
// metadata takes.
//
// Three escapes meet in one string here and each belongs to a different reader.
// A backslash and the control characters are written the way YAML's escape
// takes them. A single quote is doubled, which is what YAML itself does with
// one inside a single quoted scalar: writing "\'" instead, as this package used
// to, makes the YAML parser refuse the line and lose the diagram. What is left
// is the punctuation the kanban lexer takes before YAML sees it, and that is
// written as a mermaid entity.
func escapeMetadata(value string) string {
	escaped := strings.ReplaceAll(value, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, "\r", `\r`)
	escaped = strings.ReplaceAll(escaped, "\n", `\n`)
	escaped = strings.ReplaceAll(escaped, "\t", `\t`)
	escaped = strings.ReplaceAll(escaped, `'`, `''`)
	return "'" + escapeEntities(escaped, metadataUnsafe) + "'"
}

// escapeEntities writes each character in unsafe, and any "#" that would start
// an entity, as the form mermaid decodes.
func escapeEntities(text, unsafe string) string {
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
