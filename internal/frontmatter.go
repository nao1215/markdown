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
func quoteYAML(value string) string {
	return strconv.Quote(value)
}
