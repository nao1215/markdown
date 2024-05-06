package flowchart

const (
	// noTitle is a constant for no title.
	noTitle string = ""
)

// config is a flowchart configuration.
type config struct {
	// title is the title of the flowchart.
	title string
	// oriental is the oriental of the flowchart.
	// Default is top to bottom.
	oriental oriental
}

// newConfig returns a new Config with default values.
func newConfig() *config {
	return &config{
		oriental: tb,
	}
}

// Option sets the options for the Flowchart struct.
type Option func(*config)

// WithTitle sets the title configuration.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}

// WithOrientalTopToBottom sets the oriental configuration to top to bottom.
func WithOrientalTopToBottom() Option {
	return func(c *config) {
		c.oriental = tb
	}
}

// WithOrientalTopDown sets the oriental configuration to top down.
// Same as top to bottom.
func WithOrientalTopDown() Option {
	return func(c *config) {
		c.oriental = td
	}
}

// WithOrientalBottomToTop sets the oriental configuration to bottom to top.
func WithOrientalBottomToTop() Option {
	return func(c *config) {
		c.oriental = bt
	}
}

// WithOrientalRightToLeft sets the oriental configuration to right to left.
func WithOrientalRightToLeft() Option {
	return func(c *config) {
		c.oriental = rl
	}
}

// WithOrientalLeftToRight sets the oriental configuration to left to right.
func WithOrientalLeftToRight() Option {
	return func(c *config) {
		c.oriental = lr
	}
}
