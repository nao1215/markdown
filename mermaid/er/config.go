package er

// config is the configuration for the entity relationship diagram.
type config struct{}

// newConfig returns a new config.
func newConfig() *config {
	return &config{}
}

// Option sets the options for the PieChart struct.
type Option func(*config)
