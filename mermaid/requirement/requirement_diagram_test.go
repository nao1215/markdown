package requirement

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type failingWriter struct{}

func (f failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestNewDiagram(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opts    []Option
		want    string
		wantErr bool
	}{
		{
			name: "new diagram without options",
			opts: nil,
			want: "requirementDiagram",
		},
		{
			name: "new diagram with title",
			opts: []Option{WithTitle("Checkout Requirements")},
			want: `---
title: Checkout Requirements
---
requirementDiagram`,
		},
		{
			name:    "new diagram with title including newline",
			opts:    []Option{WithTitle("Checkout\nInjected: malicious")},
			want:    "requirementDiagram",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diagram := NewDiagram(io.Discard, tt.opts...)
			if tt.wantErr && diagram.Error() == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && diagram.Error() != nil {
				t.Fatalf("unexpected error: %v", diagram.Error())
			}

			got := strings.ReplaceAll(diagram.String(), "\r\n", "\n")

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDiagram_Build(t *testing.T) {
	t.Parallel()

	b := new(bytes.Buffer)

	d := NewDiagram(b, WithTitle("Checkout Requirements"))
	d.SetDirection(DirectionTB).
		Requirement(
			"Login",
			WithID("REQ-1"),
			WithText("The system shall support login."),
			WithRisk(RiskHigh),
			WithVerifyMethod(VerifyMethodTest),
			WithRequirementClasses("important"),
		).
		FunctionalRequirement(
			"RememberSession",
			WithID("REQ-2"),
			WithText("The system shall remember the user."),
			WithRisk(RiskMedium),
			WithVerifyMethod(VerifyMethodInspection),
		).
		Element(
			"AuthService",
			WithElementType("system"),
			WithDocRef("docs/auth.md"),
			WithElementClasses("backend"),
		).
		From("AuthService").
		Satisfies("Login").
		From("RememberSession").
		Verifies("Login").
		ClassDefs(
			Def("important", "fill:#f96,stroke:#333,stroke-width:2px"),
			Def("backend", "fill:#9cf,stroke:#333,stroke-width:1px"),
		).
		Class("Login,AuthService", "important").
		Style("AuthService", "fill:#9cf,stroke:#333")

	if err := d.Build(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `---
title: Checkout Requirements
---
requirementDiagram
    direction TB
    requirement Login:::important {
        id: "REQ-1"
        text: "The system shall support login."
        risk: High
        verifymethod: Test
    }
    functionalRequirement RememberSession {
        id: "REQ-2"
        text: "The system shall remember the user."
        risk: Medium
        verifymethod: Inspection
    }
    element AuthService:::backend {
        type: "system"
        docRef: "docs/auth.md"
    }
    AuthService - satisfies -> Login
    RememberSession - verifies -> Login
    classDef important fill:#f96,stroke:#333,stroke-width:2px
    classDef backend fill:#9cf,stroke:#333,stroke-width:1px
    class Login,AuthService important
    style AuthService fill:#9cf,stroke:#333`

	got := strings.ReplaceAll(b.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_FromBuilder(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		From("AuthService").
		Satisfies("Login").
		Verifies("RememberSession")

	if d.Error() != nil {
		t.Fatalf("unexpected error: %v", d.Error())
	}

	want := `requirementDiagram
    AuthService - satisfies -> Login
    AuthService - verifies -> RememberSession`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_QuotedNames(t *testing.T) {
	t.Parallel()

	d := NewDiagram(io.Discard).
		Requirement(
			"User Login",
			WithID("REQ-10"),
			WithText("The user shall be able to sign in."),
			WithRisk(RiskLow),
			WithVerifyMethod(VerifyMethodTest),
		).
		Element("Auth Service").
		Satisfies("Auth Service", "User Login")

	if d.Error() != nil {
		t.Fatalf("unexpected error: %v", d.Error())
	}

	want := `requirementDiagram
    requirement "User Login" {
        id: "REQ-10"
        text: "The user shall be able to sign in."
        risk: Low
        verifymethod: Test
    }
    element "Auth Service" {
    }
    "Auth Service" - satisfies -> "User Login"`

	got := strings.ReplaceAll(d.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestNormalizeQuoted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "paired quotes",
			input: `"hello"`,
			want:  "hello",
		},
		{
			name:  "leading quote only",
			input: `"hello`,
			want:  `"hello`,
		},
		{
			name:  "trailing quote only",
			input: `hello"`,
			want:  `hello"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := normalizeQuoted(tt.input); got != tt.want {
				t.Errorf("normalizeQuoted(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestQuoteEscapesSpecialChars(t *testing.T) {
	t.Parallel()

	got := quote("a\\b\rc\nd\te\"f")
	want := `"a&#92;b&#92;rc&#92;nd&#92;te&quot;f"`
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestDiagram_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func() *Diagram
		want string
	}{
		{
			name: "empty requirement name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Requirement(
						"",
						WithID("REQ-1"),
						WithText("text"),
						WithRisk(RiskLow),
						WithVerifyMethod(VerifyMethodTest),
					)
			},
			want: "requirementDiagram",
		},
		{
			name: "empty requirement id",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Requirement(
						"Login",
						WithID(""),
						WithText("text"),
						WithRisk(RiskLow),
						WithVerifyMethod(VerifyMethodTest),
					)
			},
			want: "requirementDiagram",
		},
		{
			name: "empty requirement text",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Requirement(
						"Login",
						WithID("REQ-1"),
						WithText(""),
						WithRisk(RiskLow),
						WithVerifyMethod(VerifyMethodTest),
					)
			},
			want: "requirementDiagram",
		},
		{
			name: "newline in requirement name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Requirement(
						"Log\nin",
						WithID("REQ-1"),
						WithText("text"),
						WithRisk(RiskLow),
						WithVerifyMethod(VerifyMethodTest),
					)
			},
			want: "requirementDiagram",
		},
		{
			name: "newline in requirement id",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Requirement(
						"Login",
						WithID("REQ\n-1"),
						WithText("text"),
						WithRisk(RiskLow),
						WithVerifyMethod(VerifyMethodTest),
					)
			},
			want: "requirementDiagram",
		},
		{
			name: "newline in requirement text",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Requirement(
						"Login",
						WithID("REQ-1"),
						WithText("line1\nline2"),
						WithRisk(RiskLow),
						WithVerifyMethod(VerifyMethodTest),
					)
			},
			want: "requirementDiagram",
		},
		{
			name: "invalid requirement type",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					RequirementOfType(
						RequirementType("invalid"),
						"Login",
						WithID("REQ-1"),
						WithText("text"),
						WithRisk(RiskLow),
						WithVerifyMethod(VerifyMethodTest),
					)
			},
			want: "requirementDiagram",
		},
		{
			name: "invalid risk",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Requirement(
						"Login",
						WithID("REQ-1"),
						WithText("text"),
						WithRisk(Risk("Critical")),
						WithVerifyMethod(VerifyMethodTest),
					)
			},
			want: "requirementDiagram",
		},
		{
			name: "invalid verify method",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Requirement(
						"Login",
						WithID("REQ-1"),
						WithText("text"),
						WithRisk(RiskLow),
						WithVerifyMethod(VerifyMethod("Proof")),
					)
			},
			want: "requirementDiagram",
		},
		{
			name: "empty element name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Element("")
			},
			want: "requirementDiagram",
		},
		{
			name: "newline in element type",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Element("AuthService", WithElementType("sys\ntem"))
			},
			want: "requirementDiagram",
		},
		{
			name: "newline in element docRef",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Element("AuthService", WithDocRef("docs/\nauth.md"))
			},
			want: "requirementDiagram",
		},
		{
			name: "invalid relationship",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Relation("AuthService", Relationship("connects"), "Login")
			},
			want: "requirementDiagram",
		},
		{
			name: "newline in source name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Contains("Auth\nService", "Login")
			},
			want: "requirementDiagram",
		},
		{
			name: "invalid direction",
			run: func() *Diagram {
				return NewDiagram(io.Discard).SetDirection(Direction("XX"))
			},
			want: "requirementDiagram",
		},
		{
			name: "empty style names",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Style("", "fill:red")
			},
			want: "requirementDiagram",
		},
		{
			name: "empty style value",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Style("Login", "")
			},
			want: "requirementDiagram",
		},
		{
			name: "empty classDef names",
			run: func() *Diagram {
				return NewDiagram(io.Discard).ClassDef("", "fill:red")
			},
			want: "requirementDiagram",
		},
		{
			name: "empty classDefs",
			run: func() *Diagram {
				return NewDiagram(io.Discard).ClassDefs()
			},
			want: "requirementDiagram",
		},
		{
			name: "empty classDef style",
			run: func() *Diagram {
				return NewDiagram(io.Discard).ClassDef("critical", "")
			},
			want: "requirementDiagram",
		},
		{
			name: "empty class names",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Class("", "critical")
			},
			want: "requirementDiagram",
		},
		{
			name: "empty class classNames",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Class("Login", "")
			},
			want: "requirementDiagram",
		},
		{
			name: "class shorthand without class names",
			run: func() *Diagram {
				return NewDiagram(io.Discard).ClassShorthand("Login")
			},
			want: "requirementDiagram",
		},
		{
			name: "empty class name",
			run: func() *Diagram {
				return NewDiagram(io.Discard).
					Requirement(
						"Login",
						WithID("REQ-1"),
						WithText("text"),
						WithRisk(RiskLow),
						WithVerifyMethod(VerifyMethodTest),
						WithRequirementClasses(""),
					)
			},
			want: "requirementDiagram",
		},
		{
			name: "newline in title",
			run: func() *Diagram {
				return NewDiagram(io.Discard, WithTitle("Checkout\nRequirements"))
			},
			want: "requirementDiagram",
		},
		{
			name: "lf short-circuit after error",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Element("").LF()
			},
			want: "requirementDiagram",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := tt.run()
			if d.Error() == nil {
				t.Fatal("expected error, got nil")
			}

			got := strings.ReplaceAll(d.String(), "\r\n", "\n")
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("value is mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDiagram_BuildStoresError(t *testing.T) {
	t.Parallel()

	d := NewDiagram(failingWriter{})
	err := d.Build()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if d.Error() == nil {
		t.Fatal("expected persisted error, got nil")
	}
	if !errors.Is(d.Error(), err) {
		t.Fatalf("expected Error() to wrap returned error, got %v", d.Error())
	}
}
