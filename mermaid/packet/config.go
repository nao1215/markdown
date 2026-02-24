package packet

// config is the configuration for the packet diagram.
type config struct {
	// title is the title of the packet diagram.
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

// WithTitle sets the title configuration.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}
