//go:build linux || darwin

package er_test

import (
	"fmt"
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/er"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	teachers := er.NewEntity(
		"teachers",
		[]*er.Attribute{
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
	students := er.NewEntity(
		"students",
		[]*er.Attribute{
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
	schools := er.NewEntity(
		"schools",
		[]*er.Attribute{
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

	erString := er.NewDiagram(os.Stdout).
		Relationship(
			teachers,
			students,
			er.ExactlyOneRelationship, // "||"
			er.ZeroToMoreRelationship, // "}o"
			er.Identifying,            // "--"
			"Teacher has many students",
		).
		Relationship(
			teachers,
			schools,
			er.OneToMoreRelationship,  // "|}"
			er.ExactlyOneRelationship, // "||"
			er.NonIdentifying,         // ".."
			"School has many teachers",
		).
		String()

	_ = markdown.NewMarkdown(os.Stdout).
		H2("Entity Relationship Diagram").
		CodeBlocks(markdown.SyntaxHighlightMermaid, erString).
		Build()

	// Output:
	// ## Entity Relationship Diagram
	// ```mermaid
	// erDiagram
	//     teachers ||--o{ students : "Teacher has many students"
	//     teachers }|..|| schools : "School has many teachers"
	//     schools {
	//         int id PK,UK "School ID"
	//         string name  "School Name"
	//         int teacher_id FK,UK "Teacher ID"
	//     }
	//     students {
	//         int id PK,UK "Student ID"
	//         string name  "Student Name"
	//         int teacher_id FK,UK "Teacher ID"
	//     }
	//     teachers {
	//         int id PK,UK "Teacher ID"
	//         string name  "Teacher Name"
	//     }
	//
	// ```
}

// ExampleNewDiagram shows the shape every entity relationship diagram has: a
// writer, a chain of calls, and Build.
func ExampleNewDiagram() {
	_ = er.NewDiagram(os.Stdout).
		NoRelationship(er.NewEntity("teachers", []*er.Attribute{
			{Type: "int", Name: "id", IsPrimaryKey: true, Comment: "Primary key"},
		})).
		Build()

	// Output:
	// erDiagram
	//     teachers {
	//         int id PK "Primary key"
	//     }
}

// ExampleNewEntity shows how a table is described: a name and its columns. The
// key flags put "PK", "FK" and "UK" beside a column, and the comment is drawn
// under it.
func ExampleNewEntity() {
	students := er.NewEntity("students", []*er.Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, Comment: "Student ID"},
		{Type: "int", Name: "teacher_id", IsForeignKey: true, Comment: "Teacher ID"},
		{Type: "string", Name: "email", IsUniqueKey: true, Comment: "Contact address"},
	})

	_ = er.NewDiagram(os.Stdout).NoRelationship(students).Build()

	// Output:
	// erDiagram
	//     students {
	//         int id PK "Student ID"
	//         int teacher_id FK "Teacher ID"
	//         string email UK "Contact address"
	//     }
}

// ExampleDiagram_NoRelationship adds a table that stands on its own. A diagram
// only draws a table it has been told about, so a table with no relationships
// needs this to appear at all.
func ExampleDiagram_NoRelationship() {
	_ = er.NewDiagram(os.Stdout).
		NoRelationship(er.NewEntity("audit_log", []*er.Attribute{
			{Type: "int", Name: "id", IsPrimaryKey: true, Comment: "Primary key"},
		})).
		Build()

	// Output:
	// erDiagram
	//     audit_log {
	//         int id PK "Primary key"
	//     }
}

// ExampleDiagram_Relationship joins two tables. The two Relationship values are
// the cardinality at each end, read from the left table outwards, and Identify
// says whether the right table can exist without the left one.
func ExampleDiagram_Relationship() {
	teachers := er.NewEntity("teachers", []*er.Attribute{{Type: "int", Name: "id", Comment: "Primary key"}})
	students := er.NewEntity("students", []*er.Attribute{{Type: "int", Name: "teacher_id", Comment: "Teacher ID"}})

	_ = er.NewDiagram(os.Stdout).
		Relationship(
			teachers, students,
			er.ExactlyOneRelationship, er.ZeroToMoreRelationship,
			er.Identifying,
			"teaches",
		).
		Build()

	// Output:
	// erDiagram
	//     teachers ||--o{ students : "teaches"
	//     students {
	//         int teacher_id  "Teacher ID"
	//     }
	//     teachers {
	//         int id  "Primary key"
	//     }
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := er.NewDiagram(io.Discard).
		NoRelationship(er.NewEntity("teachers", []*er.Attribute{{Type: "int", Name: "id", Comment: "Primary key"}})).
		String()

	_ = markdown.NewMarkdown(os.Stdout).
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// erDiagram
	//     teachers {
	//         int id  "Primary key"
	//     }
	//
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded.
func ExampleDiagram_Build() {
	err := er.NewDiagram(nil).
		NoRelationship(er.NewEntity("teachers", []*er.Attribute{{Type: "int", Name: "id", Comment: "Primary key"}})).
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleEntity shows what a table is once described: the value NewEntity
// returns, which the relationship calls take.
func ExampleEntity() {
	teachers := er.NewEntity("teachers", []*er.Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, Comment: "Primary key"},
	})

	_ = er.NewDiagram(os.Stdout).NoRelationship(teachers).Build()

	// Output:
	// erDiagram
	//     teachers {
	//         int id PK "Primary key"
	//     }
}

// ExampleAttribute shows one column of a table. Everything but the type and the
// name is optional.
func ExampleAttribute() {
	id := &er.Attribute{
		Type:         "int",
		Name:         "id",
		IsPrimaryKey: true,
		Comment:      "Primary key",
	}

	_ = er.NewDiagram(os.Stdout).
		NoRelationship(er.NewEntity("teachers", []*er.Attribute{id})).
		Build()

	// Output:
	// erDiagram
	//     teachers {
	//         int id PK "Primary key"
	//     }
}

// ExampleRelationship shows the cardinality at one end of a relationship.
func ExampleRelationship() {
	teachers := er.NewEntity("teachers", []*er.Attribute{{Type: "int", Name: "id", Comment: "Primary key"}})
	students := er.NewEntity("students", []*er.Attribute{{Type: "int", Name: "id", Comment: "Primary key"}})

	_ = er.NewDiagram(os.Stdout).
		Relationship(teachers, students,
			er.ZeroToOneRelationship, er.OneToMoreRelationship,
			er.NonIdentifying, "advises").
		Build()

	// Output:
	// erDiagram
	//     teachers |o..|{ students : "advises"
	//     students {
	//         int id  "Primary key"
	//     }
	//     teachers {
	//         int id  "Primary key"
	//     }
}

// ExampleIdentify says whether the right table can exist without the left one.
// An identifying relationship is drawn with a solid line and a non identifying
// one with a dashed line.
func ExampleIdentify() {
	teachers := er.NewEntity("teachers", []*er.Attribute{{Type: "int", Name: "id", Comment: "Primary key"}})
	students := er.NewEntity("students", []*er.Attribute{{Type: "int", Name: "id", Comment: "Primary key"}})

	_ = er.NewDiagram(os.Stdout).
		Relationship(teachers, students,
			er.ExactlyOneRelationship, er.ZeroToMoreRelationship,
			er.Identifying, "identifying").
		Relationship(teachers, students,
			er.ExactlyOneRelationship, er.ZeroToMoreRelationship,
			er.NonIdentifying, "non identifying").
		Build()

	// Output:
	// erDiagram
	//     teachers ||--o{ students : "identifying"
	//     teachers ||..o{ students : "non identifying"
	//     students {
	//         int id  "Primary key"
	//     }
	//     teachers {
	//         int id  "Primary key"
	//     }
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []er.Option{}

	_ = er.NewDiagram(os.Stdout, options...).
		NoRelationship(er.NewEntity("teachers", []*er.Attribute{{Type: "int", Name: "id", Comment: "Primary key"}})).
		Build()

	// Output:
	// erDiagram
	//     teachers {
	//         int id  "Primary key"
	//     }
}
