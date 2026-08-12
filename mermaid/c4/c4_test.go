package c4_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/nao1215/markdown/internal"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/c4"
)

// errWrite is the failure the writer below reports, so the test can assert that
// Build passed it through rather than inventing an error of its own.
var errWrite = errors.New("write failed")

// errWriter fails every write, which is what a full disk or a closed pipe looks
// like to Build.
type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errWrite
}

// lines joins the given lines with the line ending the library writes on this
// platform, which is what String returns them joined with.
func lines(want ...string) string {
	return strings.Join(want, internal.LineFeed())
}

func TestDiagram(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(io.Writer) *c4.Diagram
		want  string
	}{
		"an empty diagram is the header alone": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w)
			},
			want: "C4Context",
		},
		"a title is a statement rather than front matter": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w, c4.WithTitle("Banking Context"))
			},
			want: lines("C4Context", "    title Banking Context"),
		},
		"a title keeps the quotation marks a caller wrote": {
			// The C4 title is the rest of the line, quotes included, so this
			// package must not add its own: they would be drawn.
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w, c4.WithTitle(`The "Core" Ledger`))
			},
			want: lines("C4Context", `    title The "Core" Ledger`),
		},
		"a title escapes the punctuation its lexer has taken": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w, c4.WithTitle("Ledger #1; v2"))
			},
			want: lines("C4Context", "    title Ledger #35;1#59; v2"),
		},
		"a title is trimmed and a blank one writes no statement": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w, c4.WithTitle("   "))
			},
			want: "C4Context",
		},
		"every element kind writes its own macro": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).
					Person("p", "Customer").
					PersonExt("pe", "Auditor").
					System("s", "Ledger").
					SystemExt("se", "Mail").
					SystemDb("db", "Accounts").
					SystemQueue("q", "Events")
			},
			want: lines(
				"C4Context",
				`    Person(p, "Customer")`,
				`    Person_Ext(pe, "Auditor")`,
				`    System(s, "Ledger")`,
				`    System_Ext(se, "Mail")`,
				`    SystemDb(db, "Accounts")`,
				`    SystemQueue(q, "Events")`,
			),
		},
		"a description is a third argument, and an empty one is left out": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).
					Person("p", "Customer", c4.WithDescription("Holds an account.")).
					System("s", "Ledger", c4.WithDescription("   "))
			},
			want: lines(
				"C4Context",
				`    Person(p, "Customer", "Holds an account.")`,
				`    System(s, "Ledger")`,
			),
		},
		"a relationship carries an optional technology": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).
					Rel("p", "s", "Views balances").
					Rel("p", "s", "Pays", c4.WithTechnology("HTTPS")).
					BiRel("s", "db", "Syncs").
					BiRel("s", "db", "Streams", c4.WithTechnology("gRPC"))
			},
			want: lines(
				"C4Context",
				`    Rel(p, s, "Views balances")`,
				`    Rel(p, s, "Pays", "HTTPS")`,
				`    BiRel(s, db, "Syncs")`,
				`    BiRel(s, db, "Streams", "gRPC")`,
			),
		},
		"boundaries nest and indent what they hold": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).
					EnterpriseBoundary("e", "Big Bank").
					Person("p", "Customer").
					SystemBoundary("b", "Internet Banking").
					System("s", "Ledger").
					BoundaryEnd().
					BoundaryEnd().
					Boundary("r", "Regulator", c4.WithBoundaryType("external")).
					SystemExt("reg", "Reporting").
					BoundaryEnd().
					Rel("p", "s", "Uses")
			},
			want: lines(
				"C4Context",
				`    Enterprise_Boundary(e, "Big Bank") {`,
				`        Person(p, "Customer")`,
				`        System_Boundary(b, "Internet Banking") {`,
				`            System(s, "Ledger")`,
				`        }`,
				`    }`,
				`    Boundary(r, "Regulator", "external") {`,
				`        System_Ext(reg, "Reporting")`,
				`    }`,
				`    Rel(p, s, "Uses")`,
			),
		},
		"a boundary type is trimmed away when it is blank": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).
					Boundary("r", "Regulator", c4.WithBoundaryType("  ")).
					BoundaryEnd()
			},
			want: lines("C4Context", `    Boundary(r, "Regulator") {`, "    }"),
		},
		"a quotation mark in an argument becomes the entity mermaid decodes": {
			// Neither a backslash nor a doubled quote works here: the first
			// makes mermaid refuse the diagram and the second silently splits
			// the argument in two. See the quote doc comment.
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).Person("p", `The "Core" Team`)
			},
			want: lines("C4Context", `    Person(p, "The #quot;Core#quot; Team")`),
		},
		"a hash in an argument is escaped so it cannot start an entity": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).Person("p", "Ticket #quot; and #39;")
			},
			want: lines("C4Context", `    Person(p, "Ticket #35;quot; and #35;39;")`),
		},
		"the punctuation a C4 label may carry is left alone": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).Person("p", `x'#;[](){}<br/>🎉日本語:,*-|%%\x`)
			},
			want: lines("C4Context", `    Person(p, "x'#35;;[](){}<br/>🎉日本語:,*-|%%\x")`),
		},
		"text is trimmed": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).
					Person("  p  ", "  Customer  ", c4.WithDescription("  Holds an account.  ")).
					Rel("  p  ", "  s  ", "  Uses  ", c4.WithTechnology("  HTTPS  "))
			},
			want: lines(
				"C4Context",
				`    Person(p, "Customer", "Holds an account.")`,
				`    Rel(p, s, "Uses", "HTTPS")`,
			),
		},
		"LF adds a blank line": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).Person("p", "Customer").LF().System("s", "Ledger")
			},
			want: lines("C4Context", `    Person(p, "Customer")`, "", `    System(s, "Ledger")`),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			d := tt.build(buf)
			if err := d.Build(); err != nil {
				t.Fatalf("Build() = %v, want nil", err)
			}
			if got := buf.String(); got != tt.want {
				t.Errorf("Build() wrote\n%q\nwant\n%q", got, tt.want)
			}
			if got := d.String(); got != tt.want {
				t.Errorf("String() = \n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestDiagramRecordsBadInput(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(io.Writer) *c4.Diagram
		want  string
	}{
		"a title spanning lines": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w, c4.WithTitle("Ledger\nv2"))
			},
			want: "title must not contain newline characters",
		},
		"an empty element id": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).Person("  ", "Customer")
			},
			want: "Person id must not be empty",
		},
		"an element id spanning lines": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).System("a\nb", "Ledger")
			},
			want: "System id must not contain newline characters",
		},
		"an element id holding macro syntax": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).SystemDb("a,b", "Accounts")
			},
			want: `SystemDb id must not contain ",", which is C4 macro syntax`,
		},
		"an empty element label": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).SystemQueue("q", "  ")
			},
			want: "SystemQueue label must not be empty",
		},
		"an element label spanning lines": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).PersonExt("pe", "Audit\nor")
			},
			want: "Person_Ext label must not contain newline characters",
		},
		"a description spanning lines": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).SystemExt("se", "Mail", c4.WithDescription("a\nb"))
			},
			want: `description of System_Ext "Mail" must not contain newline characters`,
		},
		"an empty boundary id": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).Boundary("", "Regulator")
			},
			want: "boundary id must not be empty",
		},
		"an empty enterprise boundary label": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).EnterpriseBoundary("e", " ")
			},
			want: "boundary label must not be empty",
		},
		"a system boundary id holding macro syntax": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).SystemBoundary("a(b", "Internet Banking")
			},
			want: `boundary id must not contain "(", which is C4 macro syntax`,
		},
		"a boundary type spanning lines": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).Boundary("r", "Regulator", c4.WithBoundaryType("a\nb"))
			},
			want: `type of boundary "Regulator" must not contain newline characters`,
		},
		"closing a boundary that was never opened": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).BoundaryEnd()
			},
			want: "BoundaryEnd was called outside a boundary",
		},
		"leaving a boundary open": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).SystemBoundary("b", "Internet Banking").System("s", "Ledger")
			},
			want: "1 boundary must be closed with BoundaryEnd before Build",
		},
		"an empty relationship source": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).Rel("", "s", "Uses")
			},
			want: "Rel source must not be empty",
		},
		"a relationship destination holding macro syntax": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).BiRel("p", `s"x`, "Uses")
			},
			want: `BiRel destination must not contain "\"", which is C4 macro syntax`,
		},
		"an empty relationship label": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).Rel("p", "s", "")
			},
			want: "Rel label must not be empty",
		},
		"a technology spanning lines": {
			build: func(w io.Writer) *c4.Diagram {
				return c4.NewDiagram(w).Rel("p", "s", "Uses", c4.WithTechnology("a\nb"))
			},
			want: `technology of Rel "Uses" must not contain newline characters`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			d := tt.build(io.Discard)
			err := d.Build()
			if err == nil {
				t.Fatalf("Build() = nil, want an error mentioning %q", tt.want)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Build() = %v, want it to mention %q", err, tt.want)
			}
			if d.Error() == nil {
				t.Error("Error() = nil after a failed Build()")
			}
		})
	}
}

// TestTheFirstErrorIsKept covers the rule the whole library shares: the chain
// runs to the end after a bad call, and the error that surfaces is the first
// one, because that is the one that explains the rest.
func TestTheFirstErrorIsKept(t *testing.T) {
	t.Parallel()

	d := c4.NewDiagram(io.Discard).
		Person("", "Customer").
		System("s,x", "Ledger").
		BoundaryEnd()

	err := d.Build()
	if err == nil {
		t.Fatal("Build() = nil, want the first error")
	}
	if want := "Person id must not be empty"; !strings.Contains(err.Error(), want) {
		t.Errorf("Build() = %v, want the first error %q", err, want)
	}
}

// TestBuildWithNilWriter covers the case where a diagram is built for String()
// only and Build() is called by mistake. Build() must return an error rather
// than dereferencing the nil writer.
func TestBuildWithNilWriter(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Build() panicked with a nil writer: %v", r)
		}
	}()

	d := c4.NewDiagram(nil).Person("p", "Customer")

	// String() has always worked without a writer, and callers rely on it.
	if d.String() == "" {
		t.Fatal("String() returned nothing for a diagram with a person in it")
	}

	err := d.Build()
	if err == nil {
		t.Fatal("Build() with a nil writer must return an error")
	}
	if err.Error() != "output writer must not be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestBuildReportsWriteFailure covers the branch where the destination accepts
// the diagram and then fails. Silently returning nil there would hand the caller
// a document that was never written.
func TestBuildReportsWriteFailure(t *testing.T) {
	t.Parallel()

	err := c4.NewDiagram(errWriter{}).Person("p", "Customer").Build()
	if err == nil {
		t.Fatal("Build must report a failing writer")
	}
	if !errors.Is(err, errWrite) {
		t.Errorf("Build lost the destination error: %v", err)
	}
}

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return c4.NewDiagram(w).Person("p", "Customer").System("s", "Ledger").Rel("p", "s", "Uses")
	})
}

// TestRecordedErrorContract asserts that closing a boundary that was never
// opened surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return c4.NewDiagram(w).BoundaryEnd()
	})
}

// TestGoldenC4 pins the rendered diagram of every builder method and every
// option of this package, including both escapes and three levels of nesting.
func TestGoldenC4(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := c4.NewDiagram(buf, c4.WithTitle(`Ledger #1; the "core"`)).
		EnterpriseBoundary("bank", "Big Bank plc").
		Person("customer", "Personal Banking Customer", c4.WithDescription("Holds a personal account.")).
		SystemBoundary("banking", "Internet Banking").
		System("web", "Internet Banking System", c4.WithDescription(`Shows the "balance" (in full).`)).
		SystemDb("accounts", "Accounts Database").
		SystemQueue("events", "Event Bus", c4.WithDescription("Carries domain events.")).
		BoundaryEnd().
		BoundaryEnd().
		LF().
		PersonExt("auditor", "External Auditor").
		SystemExt("mail", "E-mail System", c4.WithDescription("Microsoft Exchange.")).
		Boundary("regulator", "Regulator", c4.WithBoundaryType("external")).
		SystemExt("reporting", "Regulatory Reporting").
		BoundaryEnd().
		LF().
		Rel("customer", "web", "Views balances", c4.WithTechnology("HTTPS (TLS 1.3)")).
		BiRel("web", "accounts", "Reads from and writes to", c4.WithTechnology("SQL/TCP")).
		Rel("web", "mail", "Sends e-mail using", c4.WithTechnology("SMTP")).
		Rel("auditor", "reporting", `Reads "the report"`).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("c4.md", buf.String()); err != nil {
		t.Error(err)
	}
}
