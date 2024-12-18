// Package er is mermaid entity relationship diagram builder.
package er

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDiagram_Build(t *testing.T) {
	t.Parallel()

	t.Run("should write the entity relationship diagram body to the output destination", func(t *testing.T) {
		t.Parallel()

		teachers := NewEntity(
			"teachers",
			[]*Attribute{
				{
					Type:         "int",
					Name:         "id",
					IsPrimaryKey: true,
					IsForeignKey: false,
					IsUniqueKey:  true,
					Comment:      "Teacher ID",
				},
				{
					Type:         "string",
					Name:         "name",
					IsPrimaryKey: false,
					IsForeignKey: false,
					IsUniqueKey:  false,
					Comment:      "Teacher Name",
				},
			},
		)
		students := NewEntity(
			"students",
			[]*Attribute{
				{
					Type:         "int",
					Name:         "id",
					IsPrimaryKey: true,
					IsForeignKey: false,
					IsUniqueKey:  true,
					Comment:      "Student ID",
				},
				{
					Type:         "string",
					Name:         "name",
					IsPrimaryKey: false,
					IsForeignKey: false,
					IsUniqueKey:  false,
					Comment:      "Student Name",
				},
				{
					Type:         "int",
					Name:         "teacher_id",
					IsPrimaryKey: false,
					IsForeignKey: true,
					IsUniqueKey:  true,
					Comment:      "Teacher ID",
				},
			},
		)
		schools := NewEntity(
			"schools",
			[]*Attribute{
				{
					Type:         "int",
					Name:         "id",
					IsPrimaryKey: true,
					IsForeignKey: false,
					IsUniqueKey:  true,
					Comment:      "School ID",
				},
				{
					Type:         "string",
					Name:         "name",
					IsPrimaryKey: false,
					IsForeignKey: false,
					IsUniqueKey:  false,
					Comment:      "School Name",
				},
				{
					Type:         "int",
					Name:         "teacher_id",
					IsPrimaryKey: false,
					IsForeignKey: true,
					IsUniqueKey:  true,
					Comment:      "Teacher ID",
				},
			},
		)
		personalComputers := NewEntity(
			"personal_computers",
			[]*Attribute{
				{
					Type:         "int",
					Name:         "id",
					IsPrimaryKey: true,
					IsForeignKey: false,
					IsUniqueKey:  true,
					Comment:      "Personal Computer ID",
				},
			},
		)

		b := new(bytes.Buffer)
		d := NewDiagram(b).
			Relationship(
				teachers,
				students,
				ExactlyOneRelationship,
				ZeroToMoreRelationship,
				Identifying,
				"Teacher has many students",
			).
			Relationship(
				teachers,
				schools,
				OneToMoreRelationship,
				ExactlyOneRelationship,
				NonIdentifying,
				"School has many teachers",
			).
			NoRelationship(personalComputers)

		if err := d.Build(); err != nil {
			t.Fatalf("error should be nil: %v", err)
		}

		want := `erDiagram
    teachers ||--o{ students : "Teacher has many students"
    teachers }|..|| schools : "School has many teachers"
    personal_computers {
        int id PK,UK "Personal Computer ID"
    }
    schools {
        int id PK,UK "School ID"
        string name  "School Name"
        int teacher_id FK,UK "Teacher ID"
    }
    students {
        int id PK,UK "Student ID"
        string name  "Student Name"
        int teacher_id FK,UK "Teacher ID"
    }
    teachers {
        int id PK,UK "Teacher ID"
        string name  "Teacher Name"
    }
`
		want = strings.ReplaceAll(want, "\r\n", "\n")
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}
