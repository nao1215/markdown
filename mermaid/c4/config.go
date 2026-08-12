package c4

// noTitle is a constant for no title.
const noTitle string = ""

// config is the configuration for the C4 context diagram.
type config struct {
	// title is the title of the C4 context diagram.
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

// WithTitle sets the title configuration.
//
// The title is written as a "title" statement rather than as front matter,
// because a C4 diagram draws the statement and ignores the front matter. It is
// also the one place in this package where the text is not quoted: mermaid
// takes the rest of the line, quotes and all, so a quoted title would draw its
// own quotation marks.
func WithTitle(title string) Option {
	return func(c *config) {
		c.title = title
	}
}

// elementConfig is the optional part of an element.
type elementConfig struct {
	// description is the sentence drawn under an element's label.
	description string
}

func newElementConfig() *elementConfig {
	return &elementConfig{}
}

// ElementOption sets the options an element takes.
type ElementOption func(*elementConfig)

// WithDescription sets the description drawn under an element's label.
//
// An empty description is the same as none: mermaid draws nothing for it.
func WithDescription(description string) ElementOption {
	return func(c *elementConfig) {
		c.description = description
	}
}

// boundaryConfig is the optional part of a boundary.
type boundaryConfig struct {
	// boundaryType is the tag mermaid draws in brackets beside the label.
	boundaryType string
}

func newBoundaryConfig() *boundaryConfig {
	return &boundaryConfig{}
}

// BoundaryOption sets the options a boundary takes.
type BoundaryOption func(*boundaryConfig)

// WithBoundaryType sets the tag mermaid draws in brackets beside a boundary
// label. EnterpriseBoundary and SystemBoundary carry their own, so this applies
// to Boundary alone.
func WithBoundaryType(boundaryType string) BoundaryOption {
	return func(c *boundaryConfig) {
		c.boundaryType = boundaryType
	}
}

// relationConfig is the optional part of a relationship.
type relationConfig struct {
	// technology is the wording drawn in brackets on the arrow.
	technology string
}

func newRelationConfig() *relationConfig {
	return &relationConfig{}
}

// RelationOption sets the options a relationship takes.
type RelationOption func(*relationConfig)

// WithTechnology sets the technology drawn in brackets on the arrow, which is
// where a C4 diagram says how two things talk to each other.
//
// An empty technology is the same as none: mermaid draws no brackets for it.
func WithTechnology(technology string) RelationOption {
	return func(c *relationConfig) {
		c.technology = technology
	}
}
