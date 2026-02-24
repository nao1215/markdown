package kanban

// config is the configuration for the kanban diagram.
type config struct {
	// title is the title of the kanban diagram.
	title string
	// ticketBaseURL is the base URL used for ticket metadata links.
	ticketBaseURL string
}

// newConfig returns a new config with default values.
func newConfig() *config {
	return &config{
		title:         noTitle,
		ticketBaseURL: noTicketBaseURL,
	}
}

const (
	// noTitle is a constant for no title.
	noTitle string = ""
	// noTicketBaseURL is a constant for no ticket base URL.
	noTicketBaseURL string = ""
)

// Option sets the options for the Diagram struct.
type Option func(*config)

// WithTitle sets the title configuration.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}

// WithTicketBaseURL sets the ticket base URL configuration.
func WithTicketBaseURL(baseURL string) Option {
	return func(c *config) {
		c.ticketBaseURL = baseURL
	}
}
