package gantt

// config is the configuration for the Gantt chart.
type config struct {
	// title is the title of the Gantt chart.
	title string
	// dateFormat is the date format for the Gantt chart.
	dateFormat string
	// axisFormat is the axis format for the Gantt chart.
	axisFormat string
	// tickInterval is the tick interval for the Gantt chart.
	tickInterval string
	// excludes specifies days to exclude (e.g., "weekends", "2024-01-01").
	excludes []string
	// todayMarker specifies the today marker style (e.g., "off" or CSS style).
	todayMarker string
}

// newConfig returns a new config with default values.
func newConfig() *config {
	return &config{
		title:      noTitle,
		dateFormat: "", // Mermaid default: YYYY-MM-DD
		axisFormat: "",
	}
}

const (
	// noTitle is a constant for no title.
	noTitle string = ""
)

// Option sets the options for the Chart struct.
type Option func(*config)

// WithTitle sets the title configuration.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}

// WithDateFormat sets the date format configuration.
// Common formats: YYYY-MM-DD, DD-MM-YYYY, YYYY-MM-DD HH:mm
func WithDateFormat(format string) Option {
	return func(c *config) {
		c.dateFormat = format
	}
}

// WithAxisFormat sets the axis format configuration.
// Common formats: %Y-%m-%d, %d/%m, %H:%M
func WithAxisFormat(format string) Option {
	return func(c *config) {
		c.axisFormat = format
	}
}

// WithTickInterval sets the tick interval configuration.
// Examples: "1day", "1week", "1month"
func WithTickInterval(interval string) Option {
	return func(c *config) {
		c.tickInterval = interval
	}
}

// WithExcludes sets the excludes configuration.
// Examples: "weekends", "2024-01-01"
func WithExcludes(excludes ...string) Option {
	return func(c *config) {
		c.excludes = excludes
	}
}

// WithTodayMarker sets the today marker configuration.
// Use "off" to disable, or CSS style like "stroke-width:5px"
func WithTodayMarker(marker string) Option {
	return func(c *config) {
		c.todayMarker = marker
	}
}
