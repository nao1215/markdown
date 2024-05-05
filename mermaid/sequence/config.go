package sequence

// Config is the configuration for the sequence diagram.
// Ref. https://mermaid.js.org/syntax/sequenceDiagram.html
type Config struct {
	// MirrorActors turns on/off the rendering of actors
	// below the diagram as well as above it.
	// Default is false.
	MirrorActors bool
	// BottomMariginAdjustment Adjusts how far down the graph ended.
	// Wide borders styles with css could generate unwanted clipping which is why this config param exists.
	// Default is 1.
	BottomMariginAdjustment uint
	// ActorFontSize is the font size of the actors.
	// Default is 14.
	ActorFontSize uint
	// ActorFontFamily sets the font family for the actor's description.
	// Default is "Open Sans", sans-serif.
	ActorFontFamily string
	// ActorFontWeight sets the font weight for the actor's description
	// Default is "Open Sans", sans-serif.
	ActorFontWeight string
	// NoteFontSize is the font size of the notes.
	// Default is 14.
	NoteFontSize uint
	// NoteFontFamily sets the font family for the note's description.
	// Default is "trebuchet ms", verdana, arial
	NoteFontFamily string
	// NoteFontWeight sets the font weight for the note's description.
	// Default is "trebuchet ms", verdana, arial
	NoteFontWeight string
	// NoteAlign sets the alignment of the note's description.
	// Default is "center"
	NoteAlign string
	// MessageFontSize is the font size of the messages.
	// Default is 16.
	MessageFontSize uint
	// MessageFontFamily sets the font family for actor<->actor messages
	// Default is "trebuchet ms", verdana, arial
	MessageFontFamily string
	// MessageFontWeight sets the font weight for actor<->actor messages
	// Default is "trebuchet ms", verdana, arial
	MessageFontWeight string
}

// NewConfig returns a new Config with default values.
func NewConfig() *Config {
	return &Config{
		MirrorActors:            false,
		BottomMariginAdjustment: 1,
		ActorFontSize:           14, // nolint:gomnd
		ActorFontFamily:         "Open Sans, sans-serif",
		ActorFontWeight:         "Open Sans, sans-serif",
		NoteFontSize:            14, // nolint:gomnd
		NoteFontFamily:          "trebuchet ms, verdana, arial",
		NoteFontWeight:          "trebuchet ms, verdana, arial",
		NoteAlign:               "center",
		MessageFontSize:         16, // nolint:gomnd
		MessageFontFamily:       "trebuchet ms, verdana, arial",
		MessageFontWeight:       "trebuchet ms, verdana, arial",
	}
}
