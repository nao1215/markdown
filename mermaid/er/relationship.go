package er

// Relationship is a type that represents the relationship between entities.
type Relationship string

const (
	// ZeroToOneRelationship is a constant that represents a zero to one relationship.
	// e.g. "|o" or "o|"
	ZeroToOneRelationship Relationship = "zero_to_one"
	// ExactlyOneRelationship is a constant that represents an exactly one relationship.
	// e.g. "||"
	ExactlyOneRelationship Relationship = "exactly_one"
	// ZeroToMoreRelationship is a constant that represents a zero to more relationship.
	// e.g. "}o" or "o{}"
	ZeroToMoreRelationship Relationship = "zero_to_more"
	// OneToMoreRelationship is a constant that represents a one to more relationship.
	// e.g. "}|" or "|}"
	OneToMoreRelationship Relationship = "one_to_more"
)

const (
	// left is a constant that represents the left side of the relationship.
	left = true
	// right is a constant that represents the right side of the relationship.
	right = false
)

// string converts the relationship to a mermaid synatax string.
func (r Relationship) string(left bool) string {
	switch r {
	case ZeroToOneRelationship:
		if left {
			return "|o"
		}
		return "o|"
	case ExactlyOneRelationship:
		return "||"
	case ZeroToMoreRelationship:
		if left {
			return "}o"
		}
		return "o{"
	case OneToMoreRelationship:
		if left {
			return "}|"
		}
		return "|{"
	default:
		return ""
	}
}
