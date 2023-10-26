package markdown

import "fmt"

// RedBadge return text with red badge format.
func (m *Markdown) RedBadge(text string) *Markdown {
	m.body = append(m.body, fmt.Sprintf("![Badge](https://img.shields.io/badge/%s-red)", text))
	return m
}

// RedBadgef return text with red badge format.
func (m *Markdown) RedBadgef(format string, args ...interface{}) *Markdown {
	return m.RedBadge(fmt.Sprintf(format, args...))
}

// YellowBadge return text with yellow badge format.
func (m *Markdown) YellowBadge(text string) *Markdown {
	m.body = append(m.body, fmt.Sprintf("![Badge](https://img.shields.io/badge/%s-yellow)", text))
	return m
}

// YellowBadgef return text with yellow badge format.
func (m *Markdown) YellowBadgef(format string, args ...interface{}) *Markdown {
	return m.YellowBadge(fmt.Sprintf(format, args...))
}

// GreenBadge return text with green badge format.
func (m *Markdown) GreenBadge(text string) *Markdown {
	m.body = append(m.body, fmt.Sprintf("![Badge](https://img.shields.io/badge/%s-green)", text))
	return m
}

// GreenBadgef return text with green badge format.
func (m *Markdown) GreenBadgef(format string, args ...interface{}) *Markdown {
	return m.GreenBadge(fmt.Sprintf(format, args...))
}
