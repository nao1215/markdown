package internal

import "strings"

// LineBreaksToBr returns s with every line ending written as "<br/>", which is
// the line break mermaid draws.
//
// A raw line ending inside a label never survives: the diagram grammars are
// line oriented, so the text after the break is read as the start of the next
// statement. Measured across every diagram type, that loses the whole diagram
// in twelve of them and silently drops or mangles content in the rest, which is
// worse. "<br/>" is the one spelling of a line break that reaches the drawing,
// so a caller's "\r\n", "\n" and "\r" are each written as one.
func LineBreaksToBr(s string) string {
	if !strings.ContainsAny(s, "\r\n") {
		return s
	}
	s = strings.ReplaceAll(s, "\r\n", "<br/>")
	s = strings.ReplaceAll(s, "\r", "<br/>")
	return strings.ReplaceAll(s, "\n", "<br/>")
}

// EscapeTitle returns a title ready to be written into a `title` statement
// whose diagram type accepts a line break in it: the angles and entity-opening
// hashes are escaped by EscapeTitleAngle, and then every line ending is written
// as "#10;", which was measured to decode to a real line break in a title. The
// order matters: the hash escape runs first so that a caller's literal "#10;"
// stays distinct from a caller's newline.
//
// "<br/>" is not the spelling here because a title draws it as the literal text
// "<br>", measured directly; it is the line the renderer puts between a title
// and a label.
func EscapeTitle(title string) string {
	title = EscapeTitleAngle(title)
	if !strings.ContainsAny(title, "\r\n") {
		return title
	}
	title = strings.ReplaceAll(title, "\r\n", "#10;")
	title = strings.ReplaceAll(title, "\r", "#10;")
	return strings.ReplaceAll(title, "\n", "#10;")
}

// FoldFrontMatterTitleCR returns a title ready for FrontMatterTitle in a
// diagram type whose renderer draws the front matter title.
//
// The YAML quoting keeps every character parseable, but a carriage return
// still loses the title in the drawing where a line feed is drawn as a line
// break, so a CRLF pair or a lone CR is folded into the line feed that works.
//
// A bare "<" is the one character left broken here: the sanitizer eats the
// rest of the title, and unlike in a title statement no escape helps, because
// a front matter title draws every entity form as the literal text it is —
// "#60;", "&lt;" and "&#60;" were each rendered and each came back verbatim.
// There is nothing to escape to, so the limit is documented rather than
// papered over with text the caller never wrote. The diagram types whose
// renderer never draws a title keep their bytes either way: their front matter
// is inert, and inert output that renders is not changed.
func FoldFrontMatterTitleCR(title string) string {
	if !strings.ContainsRune(title, '\r') {
		return title
	}
	title = strings.ReplaceAll(title, "\r\n", "\n")
	return strings.ReplaceAll(title, "\r", "\n")
}

// EscapeBareAngle returns s with every "<" that does not open a "<br/>"
// written as the entity form "#60;", and everything else untouched.
//
// The renderer's sanitizer reads a bare "<" followed by a letter as the start
// of an HTML tag and eats the rest of the text, in quoted labels as much as in
// titles: a block label, a git graph commit id, a requirement name and an xy
// chart label each drew "a" where the caller asked for "a<b then c", measured
// by rendering. "#60;" decodes to the character in every one of them. "<br/>"
// is left alone because it renders today, and run this after
// EscapeEntityOpeners: the "#60;" written here must not be escaped again.
func EscapeBareAngle(s string) string {
	if !strings.ContainsRune(s, '<') {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '<' && !strings.HasPrefix(s[i:], "<br/>") {
			b.WriteString(EntityEscape('<'))
			continue
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// EscapeTitleAngle returns a title ready to be written into the `title`
// statement of a diagram, with every "<" that does not open a "<br/>" written
// as the entity form "#60;".
//
// The renderer's sanitizer reads a bare "<" in a title as the start of an HTML
// tag and eats the rest of the text, so "cost < 10" reaches the drawing as
// "cost " or not at all; sixteen diagram types lose their title this way, and
// "#60;" was measured to decode to a literal "<" in every one of them, front
// matter included. "<br/>" is left alone because it renders today, as the
// literal text "<br>", and output that renders is not changed.
//
// A "#" that would start an entity is escaped for the same reason the label
// escapes do it: a title now written with entities has to keep a caller's
// literal "#60;" distinct from a caller's "<". A "#" anywhere else reaches the
// drawing intact and is left alone.
func EscapeTitleAngle(title string) string {
	if !strings.ContainsAny(title, "<#") {
		return title
	}

	var b strings.Builder
	b.Grow(len(title))
	for i := 0; i < len(title); i++ {
		switch {
		case title[i] == '<' && !strings.HasPrefix(title[i:], "<br/>"):
			b.WriteString(EntityEscape('<'))
		case title[i] == '#' && StartsEntity(title[i+1:]):
			b.WriteString(EntityEscape('#'))
		default:
			b.WriteByte(title[i])
		}
	}
	return b.String()
}
