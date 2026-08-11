package mindmap

// config is the configuration for the mindmap diagram.
type config struct {
	// title is the title of the mindmap diagram.
	title string
}

// newConfig returns a new config with default values.
func newConfig() *config {
	return &config{
		title: noTitle,
	}
}

const (
	// noTitle is a constant for no title.
	noTitle string = ""
)

// Option sets the options for the Diagram struct.
type Option func(*config)

// WithTitle sets the title configuration. The title is emitted as front matter,
// which is the only place mermaid accepts one for a mindmap.
//
// Note that mermaid's mindmap renderer does not draw the title today; it keeps it
// as diagram metadata. Put the wording readers must see in the nodes
// themselves if it matters.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}
