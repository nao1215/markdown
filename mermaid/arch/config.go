package arch

// config is the configuration for the Architecture Diagrams.
type config struct{}

// newConfig returns a new config.
func newConfig() *config {
	return &config{}
}

// Option sets the options for the Architecture struct.
type Option func(*config)
