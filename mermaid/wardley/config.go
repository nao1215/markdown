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
// because that is where a Wardley map takes one. Leading and trailing spaces
// are trimmed, the way every builder in this library trims the text it is
// given, and what is left is written as it was: mermaid reads the rest of the
// line and takes a quotation mark, a hash and a semicolon there as themselves.
//
// The exception is "%%", which would open a mermaid comment and drop the rest
// of the title, so a run of them is written as the entity form mermaid decodes.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}
