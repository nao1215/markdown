package wardley

// noTitle is a constant for no title.
const noTitle string = ""

// config is the configuration for the Wardley map.
type config struct {
	// title is the title of the Wardley map.
	title string
}

// newConfig returns a new config with default values.
func newConfig() *config {
	return &config{
		title: noTitle,
	}
}

// Option sets the options for the Map struct.
type Option func(*config)

// WithTitle sets the title configuration.
//
// The title is written as a "title" statement rather than as front matter,
// because that is where a Wardley map takes one, and it is written through
// unchanged: mermaid reads the rest of the line and takes every character
// probed there, quotation marks and hashes included.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}
