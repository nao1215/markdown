package timeline

// noTitle is a constant for no title.
const noTitle string = ""

// config is the configuration for the timeline diagram.
type config struct {
	// title is the title of the timeline diagram.
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
// The title is emitted as a title statement rather than as front matter,
// because mermaid's timeline renderer draws the former and keeps the latter as
// metadata a reader never sees. A title statement reads to the end of the line,
// so the punctuation that breaks a period line, a colon above all, is safe in
// one.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}
