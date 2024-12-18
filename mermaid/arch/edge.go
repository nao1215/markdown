package arch

import "fmt"

// The Position refers to the location where the Service is placed.
type Position string

const (
	// PositionLeft is the left position.
	PositionLeft Position = "L"
	// PositionRight is the right position.
	PositionRight Position = "R"
	// PositionTop is the top position.
	PositionTop Position = "T"
	// PositionBottom is the bottom position.
	PositionBottom Position = "B"
)

// Arrow can be added to each side of an edge by adding < before the
// direction on the left, and/or > after the direction on the right.
type Arrow string

const (
	// ArrowNone is the default arrow.
	ArrowNone Arrow = ""
	// ArrowRight is the right arrow.
	ArrowRight Arrow = ">"
	// ArrowLeft is the left arrow.
	ArrowLeft Arrow = "<"
)

// Edge represents an edge between two services.
// The edge can be customized with the Position and Arrow.
type Edge struct {
	// ServiceID is edge's service ID.
	// A junction ID can be specified instead of a service ID.
	ServiceID string
	// Position is edge's position. Top, Bottom, Left, Right.
	Position Position
	// Arrow is edge's arrow. None, Left, Right.
	Arrow Arrow
}

// Edges adds a string to connect two services.
// Syntax: {serviceId}:{T|B|L|R} {<}?--{>}? {T|B|L|R}:{serviceId}
func (a *Architecture) Edges(from, to Edge) *Architecture {
	a.body = append(
		a.body,
		fmt.Sprintf("    %s:%s %s--%s %s:%s",
			from.ServiceID, from.Position, from.Arrow,
			to.Arrow, to.Position, to.ServiceID,
		),
	)
	return a
}

// EdgesInAnothorGroup adds a string to connect two services in another group.
// Syntax: {serviceId}{{group}}:{T|B|L|R} {<}?--{>}? {T|B|L|R}:{serviceId}{{group}}
func (a *Architecture) EdgesInAnothorGroup(from, to Edge) *Architecture {
	a.body = append(
		a.body,
		fmt.Sprintf("    %s{group}:%s %s--%s %s:%s{group}",
			from.ServiceID, from.Position, from.Arrow,
			to.Arrow, to.Position, to.ServiceID,
		),
	)
	return a
}
