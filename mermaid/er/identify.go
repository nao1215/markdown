package er

// Identify is a type that represents the relationship between entities
// in an entity relationship diagram. Relationships may be classified as
// either identifying or non-identifying and these are rendered with either
// solid or dashed lines respectively.
type Identify bool

const (
	// Identifying is a constant that represents an identifying relationship.
	// It represents "--" in the entity relationship diagram.
	Identifying Identify = true
	// NonIdentifying is a constant that represents a non-identifying relationship.
	// It represents ".." in the entity relationship diagram.
	NonIdentifying Identify = false
)

// string converts the relationship to a mermaid synatax string.
func (i Identify) string() string {
	if i == Identifying {
		return "--"
	}
	return ".."
}
