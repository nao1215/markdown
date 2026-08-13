package quadrant

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

const (
	// labelUnsafe is what an axis label and a quadrant label cannot carry.
	//
	// A quadrant chart writes every one of them unquoted, and mermaid's grammar
	// takes almost no punctuation there: each character below loses the whole
	// chart. An angle bracket is among them because "<br/>", which reads as a
	// line break in most other diagram types, is not accepted here at all.
	// A percent run is here too: it opens a comment, and on an axis or a
	// quadrant line that loses the chart while on a point line it drops the
	// point without a word. The tilde, the at sign and the caret joined the set
	// when the probe grew past the first measurement: each one alone is a
	// lexical error that loses the chart, measured the same way as the rest.
	// A line break is in the set as well, written as "#10;", which this chart
	// was measured to decode back into a real line break.
	labelUnsafe = "\";[](){}:|<>%~@^\n"
)

// escapeLabel returns an axis or quadrant label ready to be written.
func escapeLabel(label string) string {
	return escape(label, labelUnsafe)
}

// escapePointName returns a point's name ready to be written. A point name
// loses the same characters an axis label does, measured the same way.
func escapePointName(name string) string {
	return escape(name, labelUnsafe)
}

// escape returns text ready to be written into a construct that cannot carry
// the characters in unsafe.
//
// Each is written as the entity form mermaid decodes, found by rendering them
// one at a time. A "#" that would start an entity is escaped with them, because
// this package now writes entities and a caller's literal "#59;" has to come
// out different from a caller's semicolon; a "#" that starts nothing reaches
// every label here intact and is left alone.
//
// The title is handled separately: it is the one construct in a quadrant chart
// that mermaid reads as the rest of its line, and it takes every character
// probed except the "<" and the line break that internal.EscapeTitle covers.
func escape(text, unsafe string) string {
	// A CRLF pair and a lone CR are line breaks as much as a lone LF, and the
	// byte loop below knows only the LF, so they are folded into one first.
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	if !strings.ContainsAny(text, unsafe+"#") {
		return text
	}

	var b strings.Builder
	b.Grow(len(text))
	for i := 0; i < len(text); i++ {
		switch {
		case text[i] == '%':
			// Only a run opens a comment, so "50% of traffic" is untouched.
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
