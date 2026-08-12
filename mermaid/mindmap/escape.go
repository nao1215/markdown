package mindmap

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

// shapeDelimiters are the characters that open or close a mindmap node shape.
//
// A mindmap node is unquoted text, and mermaid reads these as the start or the
// end of the shape around it: "root(text)" is a rounded node and "root[text]" a
// square one. One of them inside the text loses the whole diagram, which is
// exactly what a node saying "Deploy (staging)" did.
//
// A closing bracket is not here. mermaid accepts one on its own, since nothing
// opened, and escaping it would change output that already reaches the drawing.
const shapeDelimiters = "[(){}"

// escapeText returns text ready to be written as a mindmap node.
//
// The delimiters above are written as the entity form mermaid decodes, which
// was found by rendering each one. A "#" that would start an entity is escaped
// with them, or a node holding "#40;" and a node holding an opening parenthesis
// would draw the same picture; a "#" anywhere else is ordinary text and comes
// out unchanged.
func escapeText(text string) string {
	if !strings.ContainsAny(text, shapeDelimiters+"#") {
		return text
	}

	var b strings.Builder
	b.Grow(len(text))
	for i := 0; i < len(text); i++ {
		switch {
		case strings.IndexByte(shapeDelimiters, text[i]) >= 0:
			b.WriteString(internal.EntityEscape(rune(text[i])))
		case text[i] == '#' && internal.StartsEntity(text[i+1:]):
			b.WriteString(internal.EntityEscape('#'))
		default:
			b.WriteByte(text[i])
		}
	}
	return b.String()
}
