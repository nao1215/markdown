package piechart

// config is a pie chart configuration.
type config struct {
	// The axial position of the pie slice labels,
	// from 0.0 at the center to 1.0 at the outside edge of the circle.
	textPosition float64
	// showData is a flag to show the data in the pie chart.
	showData bool
	// title is the title of the pie chart.
	title string
}

// newConfig returns a new Config with default values.
func newConfig() *config {
	return &config{
		textPosition: defaultTextPosition,
	}
}

const (
	// defaultTextPosition is the default axial position of the pie slice labels.
	defaultTextPosition float64 = 0.75
	// minTextPosition is the minimum axial position of the pie slice labels.
	minTextPosition float64 = 0.0
	// maxTextPosition is the maximum axial position of the pie slice labels.
	maxTextPosition float64 = 1.0
	// noTitle is a constant for no title.
	noTitle string = ""
)

// Option sets the options for the PieChart struct.
type Option func(*config)

// WithTextPosition sets the axial position of the pie slice labels.
// The axial position of the pie slice labels, from 0.0 at the center
// to 1.0 at the outside edge of the circle.
// If pos is less than 0.0 or greater than 1.0, it will be set to the default value (0.75)
func WithTextPosition(pos float64) Option {
	return func(c *config) {
		if pos < minTextPosition || pos > maxTextPosition {
			pos = defaultTextPosition
		}
		c.textPosition = pos
	}
}

// WithShowData sets the showData configuration.
func WithShowData(showData bool) Option {
	return func(c *config) {
		c.showData = showData
	}
}

// WithTitle sets the title configuration.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}
