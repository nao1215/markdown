package er

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

// escapeComment returns comment ready to be written inside the quoted string an
// entity relationship diagram puts a comment in.
//
// A double quote ends that string, and mermaid then refuses the whole diagram:
// the reader gets an error box instead of a picture. A comment is prose, which
// is exactly where a quotation mark comes from, and escaping it is the
// builder's job rather than the caller's.
//
// A backslash is not an escape to mermaid's entity relationship lexer and fails
// the same way the bare quote does. What mermaid implements is its own entity
// form, "#quot;", which was found by rendering both.
//
// A "#" that would start an entity is escaped for the same reason, and only
// then: mermaid reads "#quot;" and "#123;" as the characters they name, so
// without this a comment holding "#quot;" and a comment holding a quotation
// mark would produce the same diagram. A "#" anywhere else is ordinary text and
// comes out unchanged, which is what keeps the golden files as they are.
func escapeComment(comment string) string {
	if !strings.ContainsAny(comment, `"#`) {
		return comment
	}

	var b strings.Builder
	b.Grow(len(comment))
	for i := 0; i < len(comment); i++ {
		switch {
		case comment[i] == '"':
			b.WriteString(internal.EntityEscape('"'))
		case comment[i] == '#' && internal.StartsEntity(comment[i+1:]):
			b.WriteString(internal.EntityEscape('#'))
		default:
			b.WriteByte(comment[i])
		}
	}
	return b.String()
}
