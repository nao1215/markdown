package radar

// noTitle is a constant for no title.
const noTitle string = ""

// config is the configuration for the radar chart.
type config struct {
	// title is the title of the radar chart.
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
// which is where mermaid takes one for a radar chart, and the renderer draws it.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}
