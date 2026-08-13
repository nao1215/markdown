package flowchart

import "fmt"

// Style colors one node outright.
//
// The style is mermaid's own CSS-like syntax, written through unchanged:
// "fill:#f9f,stroke:#333,stroke-width:2px". It is not escaped, because it is
// syntax rather than text; a caller putting a label in it is writing CSS, not a
// label.
func (f *Flowchart) Style(name, style string) *Flowchart {
	if f.err != nil {
		return f
	}

	f.body = append(f.body, fmt.Sprintf("%sstyle %s %s", f.indent(), name, style))
	return f
}

// ClassDef names a style so that several nodes can share it, which is what
// keeps a chart with a colored group from repeating itself.
func (f *Flowchart) ClassDef(className, style string) *Flowchart {
	if f.err != nil {
		return f
	}

	f.body = append(f.body, fmt.Sprintf("%sclassDef %s %s", f.indent(), className, style))
	return f
}

// Class applies a named style to nodes. Several nodes are given as one comma
// separated list, which is what mermaid reads there.
func (f *Flowchart) Class(names, className string) *Flowchart {
	if f.err != nil {
		return f
	}

	f.body = append(f.body, fmt.Sprintf("%sclass %s %s", f.indent(), names, className))
	return f
}

// ClickHref makes a node a link, with the text a browser shows on hover.
//
// The tooltip is escaped the way a label is: it is text, and mermaid reads it
// out of a quoted string the same way. The URL is written through unchanged.
func (f *Flowchart) ClickHref(name, url, tooltip string) *Flowchart {
	if f.err != nil {
		return f
	}

	f.body = append(f.body,
		fmt.Sprintf(`%sclick %s "%s" "%s"`, f.indent(), name, url, escapePlainText(tooltip)))
	return f
}

// ClickCall makes a node call a function in the page when it is clicked.
//
// The callback is written as mermaid wants it, with the parentheses added when
// the caller leaves them off, so that both "showOrder" and "showOrder()" reach
// the drawing as a call.
func (f *Flowchart) ClickCall(name, callback, tooltip string) *Flowchart {
	if f.err != nil {
		return f
	}

	f.body = append(f.body,
		fmt.Sprintf(`%sclick %s call %s "%s"`, f.indent(), name, ensureCall(callback), escapePlainText(tooltip)))
	return f
}

// ensureCall returns the callback with the parentheses mermaid expects.
func ensureCall(callback string) string {
	if len(callback) > 0 && callback[len(callback)-1] == ')' {
		return callback
	}
	return callback + "()"
}
