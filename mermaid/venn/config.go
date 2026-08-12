package venn

// noTitle is a constant for no title.
const noTitle string = ""

// config is the configuration for the Venn diagram.
type config struct {
	// title is the title of the Venn diagram.
	title string
}

// newConfig returns a new config with default values.
func newConfig() *config {
	return &config{
		title: noTitle,
	}
}

// Option sets the options for the Diagram struct.
type Option func(*config)

// WithTitle sets the title configuration.
//
// The title is written as a "title" statement rather than as front matter,
// because that is where a Venn diagram takes one. It is not quoted: mermaid
// reads the rest of the line, so a quoted title would draw its own quotation
// marks.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}
