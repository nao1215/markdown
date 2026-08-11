package markdown

import (
	"fmt"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// alert renders a GitHub alert, quoting every line of the body.
//
// Only the first line used to carry the "> " marker. Prose survived that by
// lazy continuation, but a list, a blank line, or a fenced block in the text
// escaped the callout and rendered as a sibling of it. Because alert text is
// usually a variable rather than a literal, the newline that caused it was
// never visible at the call site.
func alert(kind, text string) string {
	lf := internal.LineFeed()

	// Split on "\n" after dropping "\r" so a plain Go literal containing "\n"
	// is handled on Windows too, where internal.LineFeed() is "\r\n".
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	for i, line := range lines {
		switch {
		case strings.TrimSpace(line) == "":
			lines[i] = ">"
		case strings.HasPrefix(line, ">"):
			// Already quoted by the caller. Callers wrote these continuations by
			// hand to work around this very bug; prefixing them again would turn
			// them into a quote nested inside the alert.
		default:
			lines[i] = "> " + line
		}
	}

	return fmt.Sprintf("> [!%s]  %s%s", kind, lf, strings.Join(lines, lf))
}

// Note set text with note format.
func (m *Markdown) Note(text string) *Markdown {
	m.body = append(m.body, alert("NOTE", text))
	return m
}

// Notef set text with note format. It is similar to fmt.Sprintf.
func (m *Markdown) Notef(format string, args ...interface{}) *Markdown {
	return m.Note(fmt.Sprintf(format, args...))
}

// Tip set text with tip format.
func (m *Markdown) Tip(text string) *Markdown {
	m.body = append(m.body, alert("TIP", text))
	return m
}

// Tipf set text with tip format. It is similar to fmt.Sprintf.
func (m *Markdown) Tipf(format string, args ...interface{}) *Markdown {
	return m.Tip(fmt.Sprintf(format, args...))
}

// Important set text with important format.
func (m *Markdown) Important(text string) *Markdown {
	m.body = append(m.body, alert("IMPORTANT", text))
	return m
}

// Importantf set text with important format. It is similar to fmt.Sprintf.
func (m *Markdown) Importantf(format string, args ...interface{}) *Markdown {
	return m.Important(fmt.Sprintf(format, args...))
}

// Warning set text with warning format.
func (m *Markdown) Warning(text string) *Markdown {
	m.body = append(m.body, alert("WARNING", text))
	return m
}

// Warningf set text with warning format. It is similar to fmt.Sprintf.
func (m *Markdown) Warningf(format string, args ...interface{}) *Markdown {
	return m.Warning(fmt.Sprintf(format, args...))
}

// Caution set text with caution format.
func (m *Markdown) Caution(text string) *Markdown {
	m.body = append(m.body, alert("CAUTION", text))
	return m
}

// Cautionf set text with caution format. It is similar to fmt.Sprintf.
func (m *Markdown) Cautionf(format string, args ...interface{}) *Markdown {
	return m.Caution(fmt.Sprintf(format, args...))
}
