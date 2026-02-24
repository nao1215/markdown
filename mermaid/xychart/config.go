package xychart

// config is the configuration for the xy chart.
type config struct {
	// title is the title of the xy chart.
	title string
	// orientation is the orientation of the xy chart.
	orientation Orientation
}

// newConfig returns a new config with default values.
func newConfig() *config {
	return &config{
		title:       noTitle,
		orientation: OrientationVertical,
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

// WithOrientation sets the orientation configuration.
func WithOrientation(orientation Orientation) Option {
	return func(c *config) {
		c.orientation = orientation
	}
}

// WithHorizontal sets horizontal orientation.
func WithHorizontal() Option {
	return func(c *config) {
		c.orientation = OrientationHorizontal
	}
}
