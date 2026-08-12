// Package internal package is used to store the internal implementation of the mermaid package.
package internal

import (
	"fmt"
	"strconv"
)

// FrontMatterTitle returns the `title:` line of a mermaid front matter block.
//
// The value is always a double quoted YAML scalar, because mermaid runs the
// front matter through a YAML parser before it draws anything and a bare scalar
// is not safe to build from arbitrary text. "Checkout: API" and "*ref" make that
// parser throw, which loses the whole diagram; "Checkout # API" is truncated at
// the comment, and "~", "# Checkout" and "&anchor" resolve to something that is
// not the title at all. Quoting removes every one of those readings.
func FrontMatterTitle(title string) string {
	return fmt.Sprintf("title: %s", quoteYAML(title))
}

// quoteYAML returns value as a double quoted YAML scalar.
//
// strconv.Quote does the escaping because Go's double quoted form and YAML's
// agree on every escape it emits: \\, \", the \a \b \f \n \r \t \v shorthands,
// and the \xNN, \uNNNN and \UNNNNNNNN forms for everything else. Hand rolling
// the replacements misses the control characters that have no shorthand, and a
// literal control character inside a quoted scalar is a parse error, which loses
// the whole diagram rather than mangling one line of it.
//
// The two forms agree for every title that is valid UTF-8, which is every title
// a caller has. They part company below that: strconv.Quote writes a byte that
// no valid UTF-8 sequence covers as \xNN, meaning that byte, and YAML reads
// \xNN as the code point U+00NN. A title holding the single byte 0xC9 is written
// as "\xc9" and read back as "É".
//
// That is left alone deliberately. YAML is defined over Unicode, so a byte
// string that is not UTF-8 has no faithful representation in it, and every way
// out mangles the title one way or another: replacing the byte with U+FFFD
// loses it just as thoroughly and changes bytes this library has already
// shipped. Reinterpretation at least keeps the diagram, keeps the front matter
// parseable, and is documented here rather than being a surprise.
func quoteYAML(value string) string {
	return strconv.Quote(value)
}
