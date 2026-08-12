package internal

import "strconv"

// Mermaid reads "#name;" and "#123;" inside the text of a diagram as the
// character that name or number stands for. It is the only escape several of
// its grammars have: a construct that takes a quoted label has no way to spell
// a quotation mark otherwise, and a construct that reads the rest of the line
// has no way to spell the punctuation that ends it.
//
// The two helpers here are the fiddly halves of using that escape, shared
// because the diagram types that need it each need both. Which characters a
// given type has to escape, and when, is the type's own business and stays in
// its package: the answers differ, and every one of them was measured by
// rendering rather than guessed.

// EntityEscape returns the mermaid escape for r, for example "#quot;" for a
// double quote.
//
// The numeric form is used for everything except the double quote, which is
// written by name because that is the form mermaid's own documentation shows
// and the form a reader of the generated diagram will recognise.
func EntityEscape(r rune) string {
	if r == '"' {
		return "#quot;"
	}
	return "#" + strconv.Itoa(int(r)) + ";"
}

// StartsEntity reports whether s opens with the rest of a "#name;" or "#123;"
// escape, that is whether a "#" written immediately before it would be read as
// one rather than as an ordinary character.
//
// This is what keeps the escaping injective. A package that writes a quotation
// mark as "#quot;" has to escape a caller's literal "#quot;" as well, or the
// two would produce the same diagram. A "#" anywhere else is ordinary text and
// is left alone, which is what keeps output that already renders unchanged.
func StartsEntity(s string) bool {
	for i := 0; i < len(s); i++ {
		switch {
		case s[i] == ';':
			return i > 0
		case isEntityByte(s[i]):
		default:
			return false
		}
	}
	return false
}

// isEntityByte reports whether c may appear between the "#" and the ";".
func isEntityByte(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
