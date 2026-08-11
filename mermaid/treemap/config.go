package treemap

// noTitle is a constant for no title.
const noTitle string = ""

// config is the configuration for the treemap diagram.
type config struct {
	// title is the title of the treemap diagram.
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

// WithTitle sets the title configuration. The title is emitted as front matter,
// which is where mermaid takes one for a treemap, and the renderer draws it.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}
