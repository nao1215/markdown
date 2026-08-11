// Package internal package is used to store the internal implementation of the mermaid package.
package internal

import (
	"fmt"
	"strings"
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
func quoteYAML(value string) string {
	// The characters a double quoted YAML scalar cannot carry as themselves. The
	// backslash comes first, or it would escape the backslashes that the later
	// replacements introduce.
	escapes := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		"\n", `\n`,
		"\r", `\r`,
		"\t", `\t`,
	)
	return `"` + escapes.Replace(value) + `"`
}
