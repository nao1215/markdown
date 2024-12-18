package er

import (
	"fmt"
	"strings"

	"github.com/nao1215/markdown/internal"
)

// Entity is a entity of entity relationship.
type Entity struct {
	// Name is the name of the entity.
	Name string
	// Attributes is the attributes of the entity.
	Attributes []*Attribute
}

// string returns the string representation of the Entity.
func (e *Entity) string() string {
	attrs := make([]string, 0, len(e.Attributes))
	for _, a := range e.Attributes {
		attrs = append(attrs, a.string())
	}

	return fmt.Sprintf(
		"%s%s {%s%s%s%s}",
		"    ", // indent
		e.Name,
		internal.LineFeed(),
		strings.Join(attrs, internal.LineFeed()),
		internal.LineFeed(),
		"    ", // indent
	)
}

// NewEntity returns a new Entity.
func NewEntity(name string, attrs []*Attribute) Entity {
	return Entity{
		Name:       name,
		Attributes: attrs,
	}
}

// Attribute is a attribute of the entity.
type Attribute struct {
	// Type is the type of the attribute.
	Type string
	// Name is the name of the attribute.
	Name string
	// IsPrimaryKey is the flag that indicates whether the attribute is a primary key.
	IsPrimaryKey bool
	// IsForeignKey is the flag that indicates whether the attribute is a foreign key.
	IsForeignKey bool
	// IsUniqueKey is the flag that indicates whether the attribute is a unique key.
	IsUniqueKey bool
	// Comment is the comment of the attribute.
	Comment string
}

// string returns the string representation of the Attribute.
func (a *Attribute) string() string {
	var keys []string
	if a.IsPrimaryKey {
		keys = append(keys, "PK")
	}
	if a.IsForeignKey {
		keys = append(keys, "FK")
	}
	if a.IsUniqueKey {
		keys = append(keys, "UK")
	}

	s := fmt.Sprintf("        %s %s %s \"%s\"", a.Type, a.Name, strings.Join(keys, ","), a.Comment)
	s = strings.TrimSuffix(s, " ")
	return strings.ReplaceAll(s, "\"\"", "")
}

// Relationship is a relationship of entity relationship.
// leftE: left entity
// rightE: right entity
// leftR: left relationship. You choice from Relationship constants (e.g. ZeroToOneRelationship)
// rightR: right relationship. You choice from Relationship constants (e.g. ZeroToOneRelationship)
// identidy: identify of the relationship. You choice from Identify constants (e.g. Identifying)
func (d *Diagram) Relationship(leftE, rightE Entity, leftR, rightR Relationship, identidy Identify, comment string) *Diagram {
	d.body = append(
		d.body,
		fmt.Sprintf("    %s %s%s%s %s : \"%s\"",
			leftE.Name,
			leftR.string(left),
			identidy.string(),
			rightR.string(right),
			rightE.Name,
			comment,
		),
	)

	d.entities.Store(leftE.Name, leftE)
	d.entities.Store(rightE.Name, rightE)

	return d
}

// NoRelationship adds an entity that has no relationships.
func (d *Diagram) NoRelationship(e Entity) *Diagram {
	d.entities.Store(e.Name, e)
	return d
}
