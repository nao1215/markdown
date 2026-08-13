// Package er is mermaid entity relationship diagram builder.
package er

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
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

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).NoRelationship(NewEntity("teachers", nil))
	})
}

// TestGoldenEntityRelationship pins the rendered diagram of every builder
// method of this package, including every relationship cardinality and both
// identification kinds.
func TestGoldenEntityRelationship(t *testing.T) {
	t.Parallel()

	teachers := NewEntity("teachers", []*Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, IsUniqueKey: true, Comment: "Teacher ID"},
		{Type: "string", Name: "name", Comment: "Teacher name"},
	})
	students := NewEntity("students", []*Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, IsUniqueKey: true, Comment: "Student ID"},
		{Type: "int", Name: "teacher_id", IsForeignKey: true, Comment: "Teacher ID"},
	})
	schools := NewEntity("schools", []*Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, Comment: "School ID"},
	})
	clubs := NewEntity("clubs", []*Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, Comment: "Club ID"},
	})
	// An entity with no attributes at all still has to render its name.
	rooms := NewEntity("rooms", nil)

	buf := &bytes.Buffer{}
	err := NewDiagram(buf).
		Relationship(teachers, students, ExactlyOneRelationship, ZeroToMoreRelationship, Identifying, "teaches").
		Relationship(schools, teachers, OneToMoreRelationship, ExactlyOneRelationship, NonIdentifying, "employs").
		Relationship(students, clubs, ZeroToOneRelationship, OneToMoreRelationship, Identifying, "joins").
		NoRelationship(rooms).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("entity_relationship.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestBuildWithNilWriter covers the case where a diagram is built for String()
// only and Build() is called by mistake. Build() used to dereference the nil
// writer and take the process down; it has to return an error instead.
func TestBuildWithNilWriter(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Build() panicked with a nil writer: %v", r)
		}
	}()

	d := NewDiagram(nil)

	// String() has always worked without a writer, and callers rely on it.
	_ = d.String()

	err := d.Build()
	if err == nil {
		t.Fatal("Build() with a nil writer must return an error")
	}
	if err.Error() != "output writer must not be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRelationship_string(t *testing.T) {
	t.Parallel()

	type args struct {
		lr bool
	}
	tests := []struct {
		name string
		r    Relationship
		args args
		want string
	}{
		{
			name: "ZeroToOneRelationship, left",
			r:    ZeroToOneRelationship,
			args: args{lr: left},
			want: "|o",
		},
		{
			name: "ZeroToOneRelationship, right",
			r:    ZeroToOneRelationship,
			args: args{lr: right},
			want: "o|",
		},
		{
			name: "ExactlyOneRelationship, left",
			r:    ExactlyOneRelationship,
			args: args{lr: left},
			want: "||",
		},
		{
			name: "ExactlyOneRelationship, right",
			r:    ExactlyOneRelationship,
			args: args{lr: right},
			want: "||",
		},
		{
			name: "ZeroToMoreRelationship, left",
			r:    ZeroToMoreRelationship,
			args: args{lr: left},
			want: "}o",
		},
		{
			name: "ZeroToMoreRelationship, right",
			r:    ZeroToMoreRelationship,
			args: args{lr: right},
			want: "o{",
		},
		{
			name: "OneToMoreRelationship, left",
			r:    OneToMoreRelationship,
			args: args{lr: left},
			want: "}|",
		},
		{
			name: "OneToMoreRelationship, right",
			r:    OneToMoreRelationship,
			args: args{lr: right},
			want: "|{",
		},
		{
			name: "default",
			r:    "default",
			args: args{lr: left},
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.r.string(tt.args.lr); got != tt.want {
				t.Errorf("Relationship.string() = %v, want %v", got, tt.want)
			}
		})
	}
}

// errWrite is the failure the writer below reports, so the test can assert that
// Build passed it through rather than inventing an error of its own.
var errWrite = errors.New("write failed")

// errWriter fails every write, which is what a full disk or a closed pipe looks
// like to Build.
type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errWrite
}

// TestBuildReportsWriteFailure covers the branch where the destination accepts
// the diagram and then fails. Silently returning nil there would hand the caller
// a document that was never written.
func TestBuildReportsWriteFailure(t *testing.T) {
	t.Parallel()

	err := NewDiagram(errWriter{}).Build()
	if err == nil {
		t.Fatal("Build must report a failing writer")
	}
	if !errors.Is(err, errWrite) {
		t.Errorf("Build lost the destination error: %v", err)
	}
}

// TestCommentEscapesTheQuoteThatEndsIt names the character this escaping buys.
// A double quote in an attribute comment or a relationship comment used to
// reach mermaid unescaped and lose the whole diagram: the reader got an error
// box rather than a picture.
func TestCommentEscapesTheQuoteThatEndsIt(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(io.Writer) *Diagram
		want  string
	}{
		"a quote in an attribute comment": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).NoRelationship(NewEntity("teachers", []*Attribute{
					{Type: "int", Name: "id", IsPrimaryKey: true, Comment: `the "primary" key`},
				}))
			},
			want: `        int id PK "the #quot;primary#quot; key"`,
		},
		"a quote in a relationship comment": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).Relationship(
					NewEntity("teachers", []*Attribute{{Type: "int", Name: "id"}}),
					NewEntity("students", []*Attribute{{Type: "int", Name: "id"}}),
					ExactlyOneRelationship, ZeroToMoreRelationship, Identifying, `"teaches"`,
				)
			},
			want: `"#quot;teaches#quot;"`,
		},
		"a named entity in a comment is escaped": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).NoRelationship(NewEntity("teachers", []*Attribute{
					{Type: "int", Name: "id", Comment: "a#quot;b"},
				}))
			},
			want: `"a#35;quot;b"`,
		},
		"a plain hash in a comment is left alone": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).NoRelationship(NewEntity("teachers", []*Attribute{
					{Type: "int", Name: "id", Comment: "PR #123 merged"},
				}))
			},
			want: `"PR #123 merged"`,
		},
		"a line break in a comment becomes the break mermaid draws": {
			// Raw it was swallowed and "first\nsecond" drew "firstsecond".
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).NoRelationship(NewEntity("teachers", []*Attribute{
					{Type: "int", Name: "id", Comment: "first\nsecond"},
				}))
			},
			want: `"first<br/>second"`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.build(io.Discard).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("diagram =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}

// TestErrorReportsTheRecordedError pins the method the v1.0.0 API audit found
// missing. Every other builder in this library lets the recorded error be read
// before anything is written, and this one did not.
func TestErrorReportsTheRecordedError(t *testing.T) {
	t.Parallel()

	d := NewDiagram(nil).NoRelationship(NewEntity("teachers", []*Attribute{
		{Type: "int", Name: "id", Comment: "Primary key"},
	}))

	if err := d.Error(); err != nil {
		t.Errorf("Error() = %v before Build, want nil", err)
	}

	// Build reports the same error when it is what stopped the write.
	fromBuild := d.Build()
	if fromBuild == nil {
		t.Fatal("Build() = nil with a nil writer, want an error")
	}
	if d.Error() == nil || d.Error().Error() != fromBuild.Error() {
		t.Errorf("Error() = %v, want the error Build returned, %v", d.Error(), fromBuild)
	}
}

// TestZeroValueDiagramDoesNotPanic pins that er.Diagram{} takes calls, which it
// did when the entity set was a sync.Map and would not have when it became a
// plain map without this.
//
// The type is exported, so a caller can write the zero value, and nothing in
// this library panics on how it is called: a bad call records an error and the
// chain runs on. A nil map assignment would have broken that.
func TestZeroValueDiagramDoesNotPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("the zero value panicked: %v", r)
		}
	}()

	var d Diagram
	d.NoRelationship(NewEntity("teachers", []*Attribute{
		{Type: "int", Name: "id", IsPrimaryKey: true, Comment: "Primary key"},
	}))

	if want := "teachers {"; !strings.Contains(d.String(), want) {
		t.Errorf("diagram = %q, want it to contain %q", d.String(), want)
	}
}
