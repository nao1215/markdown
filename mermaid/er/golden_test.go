package er_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/er"
)

// TestGoldenEntityRelationship pins the rendered diagram of every builder
// method of this package, including every relationship cardinality and both
// identification kinds.
func TestGoldenEntityRelationship(t *testing.T) {
	t.Parallel()

	teachers := er.NewEntity("teachers", []*er.Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, IsUniqueKey: true, Comment: "Teacher ID"},
		{Type: "string", Name: "name", Comment: "Teacher name"},
	})
	students := er.NewEntity("students", []*er.Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, IsUniqueKey: true, Comment: "Student ID"},
		{Type: "int", Name: "teacher_id", IsForeignKey: true, Comment: "Teacher ID"},
	})
	schools := er.NewEntity("schools", []*er.Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, Comment: "School ID"},
	})
	clubs := er.NewEntity("clubs", []*er.Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, Comment: "Club ID"},
	})
	// An entity with no attributes at all still has to render its name.
	rooms := er.NewEntity("rooms", nil)

	buf := &bytes.Buffer{}
	err := er.NewDiagram(buf).
		Relationship(teachers, students, er.ExactlyOneRelationship, er.ZeroToMoreRelationship, er.Identifying, "teaches").
		Relationship(schools, teachers, er.OneToMoreRelationship, er.ExactlyOneRelationship, er.NonIdentifying, "employs").
		Relationship(students, clubs, er.ZeroToOneRelationship, er.OneToMoreRelationship, er.Identifying, "joins").
		NoRelationship(rooms).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("entity_relationship.md", buf.String()); err != nil {
		t.Error(err)
	}
}
