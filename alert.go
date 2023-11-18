package markdown

import "fmt"

// Note set text with note format.
func (m *Markdown) Note(text string) *Markdown {
	m.body = append(m.body, fmt.Sprintf("> [!NOTE]  \n> %s", text))
	return m
}

// Notef set text with note format. It is similar to fmt.Sprintf.
func (m *Markdown) Notef(format string, args ...interface{}) *Markdown {
	return m.Note(fmt.Sprintf(format, args...))
}

// Tip set text with tip format.
func (m *Markdown) Tip(text string) *Markdown {
	m.body = append(m.body, fmt.Sprintf("> [!TIP]  \n> %s", text))
	return m
}

// Tipf set text with tip format. It is similar to fmt.Sprintf.
func (m *Markdown) Tipf(format string, args ...interface{}) *Markdown {
	return m.Tip(fmt.Sprintf(format, args...))
}

// Important set text with important format.
func (m *Markdown) Important(text string) *Markdown {
	m.body = append(m.body, fmt.Sprintf("> [!IMPORTANT]  \n> %s", text))
	return m
}

// Importantf set text with important format. It is similar to fmt.Sprintf.
func (m *Markdown) Importantf(format string, args ...interface{}) *Markdown {
	return m.Important(fmt.Sprintf(format, args...))
}

// Warning set text with warning format.
func (m *Markdown) Warning(text string) *Markdown {
	m.body = append(m.body, fmt.Sprintf("> [!WARNING]  \n> %s", text))
	return m
}

// Warningf set text with warning format. It is similar to fmt.Sprintf.
func (m *Markdown) Warningf(format string, args ...interface{}) *Markdown {
	return m.Warning(fmt.Sprintf(format, args...))
}

// Caution set text with caution format.
func (m *Markdown) Caution(text string) *Markdown {
	m.body = append(m.body, fmt.Sprintf("> [!CAUTION]  \n> %s", text))
	return m
}

// Cautionf set text with caution format. It is similar to fmt.Sprintf.
func (m *Markdown) Cautionf(format string, args ...interface{}) *Markdown {
	return m.Caution(fmt.Sprintf(format, args...))
}
