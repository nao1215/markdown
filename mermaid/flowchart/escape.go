package flowchart

import "strings"

// escapeText returns text ready to be written inside a flowchart's quoted
// label.
//
// A double quote ends the label, and mermaid refuses the whole diagram when one
// appears inside it: the reader gets an error box instead of a picture. Neither
// escape that looks obvious works. A backslash is not an escape to mermaid's
// flowchart lexer and fails the same way, and doubling the quote fails too.
// What mermaid does implement is its own entity form, "#quot;", which was found
// by rendering all four.
//
// A "#" that starts one of those entities is escaped for the same reason, and
// only then: mermaid reads "#quot;" and "#123;" as the characters they name, so
// without this a label holding "#quot;" and a label holding a quotation mark
// would produce the same diagram. A "#" anywhere else is ordinary text and is
// left exactly as it was, because that output already renders and is pinned by
// the golden files.
func escapeText(text string) string {
	if !strings.ContainsAny(text, `"#`) {
		return text
	}

	var b strings.Builder
	b.Grow(len(text))
	for i := 0; i < len(text); i++ {
		switch {
		case text[i] == '"':
			b.WriteString("#quot;")
		case text[i] == '#' && startsEntity(text[i+1:]):
			b.WriteString("#35;")
		default:
			b.WriteByte(text[i])
		}
	}
	return b.String()
}

// startsEntity reports whether rest completes an entity that began with a "#",
// that is whether it opens with one or more letters or digits and then a ";".
func startsEntity(rest string) bool {
	for i := 0; i < len(rest); i++ {
		switch {
		case rest[i] == ';':
			return i > 0
		case isEntityByte(rest[i]):
		default:
			return false
		}
	}
	return false
}

// isEntityByte reports whether c may appear in an entity name.
func isEntityByte(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
