//go:build linux || darwin

package requirement_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/requirement"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := requirement.NewDiagram(
		io.Discard,
		requirement.WithTitle("Checkout Requirements"),
	).
		Requirement(
			"Login",
			requirement.WithID("REQ-1"),
			requirement.WithText("The system shall support login."),
			requirement.WithRisk(requirement.RiskHigh),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest),
		).
		Element("AuthService", requirement.WithElementType("system")).
		From("AuthService").
		Satisfies("Login").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Requirement Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Requirement Diagram
	// ```mermaid
	// ---
	// title: "Checkout Requirements"
	// ---
	// requirementDiagram
	//     requirement Login {
	//         id: "REQ-1"
	//         text: "The system shall support login."
	//         risk: High
	//         verifymethod: Test
	//     }
	//     element AuthService {
	//         type: "system"
	//     }
	//     AuthService - satisfies -> Login
	// ```
}

// ExampleDiagram_Requirement adds a plain requirement.
func ExampleDiagram_Requirement() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("The system shall log in", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Build()

	// Output:
	// requirementDiagram
	//     requirement "The system shall log in" {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleDiagram_FunctionalRequirement adds a requirement about what the system does.
func ExampleDiagram_FunctionalRequirement() {
	_ = requirement.NewDiagram(os.Stdout).
		FunctionalRequirement("The system shall log in", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Build()

	// Output:
	// requirementDiagram
	//     functionalRequirement "The system shall log in" {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleDiagram_InterfaceRequirement adds a requirement about how it is talked to.
func ExampleDiagram_InterfaceRequirement() {
	_ = requirement.NewDiagram(os.Stdout).
		InterfaceRequirement("The system shall log in", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Build()

	// Output:
	// requirementDiagram
	//     interfaceRequirement "The system shall log in" {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleDiagram_PerformanceRequirement adds a requirement about how fast or how much.
func ExampleDiagram_PerformanceRequirement() {
	_ = requirement.NewDiagram(os.Stdout).
		PerformanceRequirement("The system shall log in", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Build()

	// Output:
	// requirementDiagram
	//     performanceRequirement "The system shall log in" {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleDiagram_PhysicalRequirement adds a requirement about the hardware it runs on.
func ExampleDiagram_PhysicalRequirement() {
	_ = requirement.NewDiagram(os.Stdout).
		PhysicalRequirement("The system shall log in", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Build()

	// Output:
	// requirementDiagram
	//     physicalRequirement "The system shall log in" {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleDiagram_DesignConstraint adds a limit on how the system may be built.
func ExampleDiagram_DesignConstraint() {
	_ = requirement.NewDiagram(os.Stdout).
		DesignConstraint("The system shall log in", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Build()

	// Output:
	// requirementDiagram
	//     designConstraint "The system shall log in" {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleDiagram_Contains says that the first named thing holds the second.
func ExampleDiagram_Contains() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		Contains("login", "auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - contains -> "auth service"
}

// ExampleDiagram_Copies says that the first named thing copies the second.
func ExampleDiagram_Copies() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		Copies("login", "auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - copies -> "auth service"
}

// ExampleDiagram_Derives says that the first named thing is derived from the second.
func ExampleDiagram_Derives() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		Derives("login", "auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - derives -> "auth service"
}

// ExampleDiagram_Satisfies says that the first named thing satisfies the second.
func ExampleDiagram_Satisfies() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		Satisfies("login", "auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - satisfies -> "auth service"
}

// ExampleDiagram_Verifies says that the first named thing verifies the second.
func ExampleDiagram_Verifies() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		Verifies("login", "auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - verifies -> "auth service"
}

// ExampleDiagram_Refines says that the first named thing refines the second.
func ExampleDiagram_Refines() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		Refines("login", "auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - refines -> "auth service"
}

// ExampleDiagram_Traces says that the first named thing traces to the second.
func ExampleDiagram_Traces() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		Traces("login", "auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - traces -> "auth service"
}

// ExampleSourceRelationBuilder_Contains says that the source named by From holds
// the thing given here.
func ExampleSourceRelationBuilder_Contains() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		From("login").
		Contains("auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - contains -> "auth service"
}

// ExampleSourceRelationBuilder_Copies says that the source named by From copies
// the thing given here.
func ExampleSourceRelationBuilder_Copies() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		From("login").
		Copies("auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - copies -> "auth service"
}

// ExampleSourceRelationBuilder_Derives says that the source named by From is derived from
// the thing given here.
func ExampleSourceRelationBuilder_Derives() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		From("login").
		Derives("auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - derives -> "auth service"
}

// ExampleSourceRelationBuilder_Satisfies says that the source named by From satisfies
// the thing given here.
func ExampleSourceRelationBuilder_Satisfies() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		From("login").
		Satisfies("auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - satisfies -> "auth service"
}

// ExampleSourceRelationBuilder_Verifies says that the source named by From verifies
// the thing given here.
func ExampleSourceRelationBuilder_Verifies() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		From("login").
		Verifies("auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - verifies -> "auth service"
}

// ExampleSourceRelationBuilder_Refines says that the source named by From refines
// the thing given here.
func ExampleSourceRelationBuilder_Refines() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		From("login").
		Refines("auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - refines -> "auth service"
}

// ExampleSourceRelationBuilder_Traces says that the source named by From traces to
// the thing given here.
func ExampleSourceRelationBuilder_Traces() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest), requirement.WithText("The system shall log a user in.")).
		Element("auth service").
		From("login").
		Traces("auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//     }
	//     login - traces -> "auth service"
}

// ExampleNewDiagram shows the shape every requirement diagram has: a writer,
// the requirements, what they relate to, and Build.
//
// A requirement needs all four of its fields. mermaid draws the block from
// them, and leaving one out is recorded as an error rather than drawn short.
func ExampleNewDiagram() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest)).
		Element("auth service", requirement.WithElementType("service")).
		Satisfies("login", "auth service").
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	//     element "auth service" {
	//         type: "service"
	//     }
	//     login - satisfies -> "auth service"
}

// ExampleDiagram_RequirementOfType adds a requirement whose kind is named
// outright, for code that decides between the kinds rather than picking one.
func ExampleDiagram_RequirementOfType() {
	_ = requirement.NewDiagram(os.Stdout).
		RequirementOfType(requirement.RequirementTypePerformance, "latency", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest)).
		Build()

	// Output:
	// requirementDiagram
	//     performanceRequirement latency {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleDiagram_Element adds something that is not a requirement: the part of
// the system a requirement is about.
func ExampleDiagram_Element() {
	_ = requirement.NewDiagram(os.Stdout).
		Element("auth service", requirement.WithElementType("service")).
		Build()

	// Output:
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
}

// ExampleDiagram_Relation joins two things with a relationship named outright.
func ExampleDiagram_Relation() {
	_ = requirement.NewDiagram(os.Stdout).
		Element("auth service", requirement.WithElementType("service")).
		Relation("login", requirement.RelationshipSatisfies, "auth service").
		Build()

	// Output:
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
	//     login - satisfies -> "auth service"
}

// ExampleDiagram_From names a source once and hangs several relationships off
// it, which keeps a requirement with many links readable.
func ExampleDiagram_From() {
	_ = requirement.NewDiagram(os.Stdout).
		Element("auth service", requirement.WithElementType("service")).
		Element("session store", requirement.WithElementType("service")).
		From("login").
		Satisfies("auth service").
		Traces("session store").
		Build()

	// Output:
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
	//     element "session store" {
	//         type: "service"
	//     }
	//     login - satisfies -> "auth service"
	//     login - traces -> "session store"
}

// ExampleSourceRelationBuilder shows what From returns: the same source, with
// one relationship per call.
func ExampleSourceRelationBuilder() {
	_ = requirement.NewDiagram(os.Stdout).
		Element("auth service", requirement.WithElementType("service")).
		From("login").
		Satisfies("auth service").
		Build()

	// Output:
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
	//     login - satisfies -> "auth service"
}

// ExampleSourceRelationBuilder_Relation joins the source to something with a
// relationship named outright.
func ExampleSourceRelationBuilder_Relation() {
	_ = requirement.NewDiagram(os.Stdout).
		Element("auth service", requirement.WithElementType("service")).
		From("login").
		Relation(requirement.RelationshipSatisfies, "auth service").
		Build()

	// Output:
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
	//     login - satisfies -> "auth service"
}

// ExampleDiagram_SetDirection says which way the diagram is laid out.
func ExampleDiagram_SetDirection() {
	_ = requirement.NewDiagram(os.Stdout).
		SetDirection(requirement.DirectionLR).
		Element("auth service", requirement.WithElementType("service")).
		Build()

	// Output:
	// requirementDiagram
	//     direction LR
	//     element "auth service" {
	//         type: "service"
	//     }
}

// ExampleDiagram_Style colors one thing outright.
func ExampleDiagram_Style() {
	_ = requirement.NewDiagram(os.Stdout).
		Element("auth service", requirement.WithElementType("service")).
		Style("auth service", "fill:#f9f").
		Build()

	// Output:
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
	//     style auth service fill:#f9f
}

// ExampleDiagram_ClassDef names a style so several things can share it.
func ExampleDiagram_ClassDef() {
	_ = requirement.NewDiagram(os.Stdout).
		ClassDef("urgent", "fill:#f96").
		Element("auth service", requirement.WithElementType("service")).
		Class("auth service", "urgent").
		Build()

	// Output:
	// requirementDiagram
	//     classDef urgent fill:#f96
	//     element "auth service" {
	//         type: "service"
	//     }
	//     class auth service urgent
}

// ExampleDiagram_ClassDefs names several styles at once.
func ExampleDiagram_ClassDefs() {
	_ = requirement.NewDiagram(os.Stdout).
		ClassDefs(
			requirement.Def("urgent", "fill:#f96"),
			requirement.Def("done", "fill:#9f6"),
		).
		Build()

	// Output:
	// requirementDiagram
	//     classDef urgent fill:#f96
	//     classDef done fill:#9f6
}

// ExampleDef describes one named style, for passing to ClassDefs.
func ExampleDef() {
	_ = requirement.NewDiagram(os.Stdout).
		ClassDefs(requirement.Def("urgent", "fill:#f96")).
		Build()

	// Output:
	// requirementDiagram
	//     classDef urgent fill:#f96
}

// ExampleDiagram_Class applies named styles to things.
func ExampleDiagram_Class() {
	_ = requirement.NewDiagram(os.Stdout).
		ClassDef("urgent", "fill:#f96").
		Element("auth service", requirement.WithElementType("service")).
		Class("auth service", "urgent").
		Build()

	// Output:
	// requirementDiagram
	//     classDef urgent fill:#f96
	//     element "auth service" {
	//         type: "service"
	//     }
	//     class auth service urgent
}

// ExampleDiagram_ClassShorthand applies styles with mermaid's ":::" form, which
// says the same thing as Class in one token.
func ExampleDiagram_ClassShorthand() {
	_ = requirement.NewDiagram(os.Stdout).
		ClassDef("urgent", "fill:#f96").
		Element("auth service", requirement.WithElementType("service")).
		ClassShorthand("auth service", "urgent").
		Build()

	// Output:
	// requirementDiagram
	//     classDef urgent fill:#f96
	//     element "auth service" {
	//         type: "service"
	//     }
	//     "auth service":::urgent
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := requirement.NewDiagram(io.Discard).
		Element("auth service", requirement.WithElementType("service")).
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the error the chain
// recorded.
func ExampleDiagram_Build() {
	err := requirement.NewDiagram(nil).
		Element("auth service", requirement.WithElementType("service")).
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does. A requirement missing
// one of its four fields is exactly the kind of thing it reports.
func ExampleDiagram_Error() {
	d := requirement.NewDiagram(io.Discard).Requirement("login", requirement.WithID("1"))
	fmt.Println("error:", d.Error())

	// Output:
	// error: requirement text must not be empty
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = requirement.NewDiagram(os.Stdout).
		Element("auth service", requirement.WithElementType("service")).
		LF().
		Element("session store", requirement.WithElementType("service")).
		Build()

	// Output:
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
	//
	//     element "session store" {
	//         type: "service"
	//     }
}

// ExampleWithID gives a requirement its number, which is how a specification
// refers to it outside the diagram.
func ExampleWithID() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest)).
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleWithText spells the requirement out, where the name is only a handle.
func ExampleWithText() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest)).
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleWithRisk says how much is at stake if the requirement is not met.
func ExampleWithRisk() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest)).
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleRisk shows the three levels a requirement can carry.
func ExampleRisk() {
	for _, risk := range []requirement.Risk{
		requirement.RiskLow, requirement.RiskMedium, requirement.RiskHigh,
	} {
		_ = requirement.NewDiagram(os.Stdout).
			Requirement("login",
				requirement.WithID("1"),
				requirement.WithText("The system shall log a user in."),
				requirement.WithRisk(risk),
				requirement.WithVerifyMethod(requirement.VerifyMethodTest),
			).
			Build()
		fmt.Println()
	}

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Low
	//         verifymethod: Test
	//     }
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: High
	//         verifymethod: Test
	//     }
}

// ExampleWithVerifyMethod says how the requirement will be shown to be met.
func ExampleWithVerifyMethod() {
	_ = requirement.NewDiagram(os.Stdout).
		Requirement("login", requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest)).
		Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleVerifyMethod shows the four ways a requirement can be shown to be met.
func ExampleVerifyMethod() {
	for _, method := range []requirement.VerifyMethod{
		requirement.VerifyMethodAnalysis,
		requirement.VerifyMethodInspection,
		requirement.VerifyMethodTest,
		requirement.VerifyMethodDemonstration,
	} {
		_ = requirement.NewDiagram(os.Stdout).
			Requirement("login",
				requirement.WithID("1"),
				requirement.WithText("The system shall log a user in."),
				requirement.WithRisk(requirement.RiskMedium),
				requirement.WithVerifyMethod(method),
			).
			Build()
		fmt.Println()
	}

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Analysis
	//     }
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Inspection
	//     }
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Demonstration
	//     }
}

// ExampleWithRequirementClasses styles a requirement as it is declared, which
// saves a separate Class call.
func ExampleWithRequirementClasses() {
	_ = requirement.NewDiagram(os.Stdout).
		ClassDef("urgent", "fill:#f96").
		Requirement("login",
			requirement.WithID("1"),
			requirement.WithText("The system shall log a user in."),
			requirement.WithRisk(requirement.RiskHigh),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest),
			requirement.WithRequirementClasses("urgent"),
		).
		Build()

	// Output:
	// requirementDiagram
	//     classDef urgent fill:#f96
	//     requirement login:::urgent {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: High
	//         verifymethod: Test
	//     }
}

// ExampleWithElementType says what kind of thing an element is.
func ExampleWithElementType() {
	_ = requirement.NewDiagram(os.Stdout).
		Element("auth service", requirement.WithElementType("service")).
		Build()

	// Output:
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
}

// ExampleWithDocRef points an element at where it is written down.
func ExampleWithDocRef() {
	_ = requirement.NewDiagram(os.Stdout).
		Element("auth service",
			requirement.WithElementType("service"),
			requirement.WithDocRef("docs/auth.md"),
		).
		Build()

	// Output:
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//         docRef: "docs/auth.md"
	//     }
}

// ExampleWithElementClasses styles an element as it is declared.
func ExampleWithElementClasses() {
	_ = requirement.NewDiagram(os.Stdout).
		ClassDef("external", "fill:#eee").
		Element("auth service",
			requirement.WithElementType("service"),
			requirement.WithElementClasses("external"),
		).
		Build()

	// Output:
	// requirementDiagram
	//     classDef external fill:#eee
	//     element "auth service":::external {
	//         type: "service"
	//     }
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = requirement.NewDiagram(os.Stdout, requirement.WithTitle("Login")).
		Element("auth service", requirement.WithElementType("service")).
		Build()

	// Output:
	// ---
	// title: "Login"
	// ---
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
}

// ExampleRequirementType shows the six kinds a requirement can be, for the form
// that names one outright.
func ExampleRequirementType() {
	for _, kind := range []requirement.RequirementType{
		requirement.RequirementTypeRequirement,
		requirement.RequirementTypeFunctional,
		requirement.RequirementTypeInterface,
		requirement.RequirementTypePerformance,
		requirement.RequirementTypePhysical,
		requirement.RequirementTypeDesignConstraint,
	} {
		_ = requirement.NewDiagram(os.Stdout).
			RequirementOfType(kind, "login",
				requirement.WithID("1"),
				requirement.WithText("The system shall log a user in."),
				requirement.WithRisk(requirement.RiskMedium),
				requirement.WithVerifyMethod(requirement.VerifyMethodTest),
			).
			Build()
		fmt.Println()
	}

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	// requirementDiagram
	//     functionalRequirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	// requirementDiagram
	//     interfaceRequirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	// requirementDiagram
	//     performanceRequirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	// requirementDiagram
	//     physicalRequirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
	// requirementDiagram
	//     designConstraint login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: Medium
	//         verifymethod: Test
	//     }
}

// ExampleRelationship shows the seven ways two things can be related.
func ExampleRelationship() {
	_ = requirement.NewDiagram(os.Stdout).
		Element("auth service", requirement.WithElementType("service")).
		Relation("login", requirement.RelationshipSatisfies, "auth service").
		Relation("login", requirement.RelationshipTraces, "auth service").
		Build()

	// Output:
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
	//     login - satisfies -> "auth service"
	//     login - traces -> "auth service"
}

// ExampleDirection shows the four ways a diagram can be laid out.
func ExampleDirection() {
	for _, direction := range []requirement.Direction{
		requirement.DirectionTB,
		requirement.DirectionBT,
		requirement.DirectionLR,
		requirement.DirectionRL,
	} {
		_ = requirement.NewDiagram(os.Stdout).
			SetDirection(direction).
			Element("auth service", requirement.WithElementType("service")).
			Build()
		fmt.Println()
	}

	// Output:
	// requirementDiagram
	//     direction TB
	//     element "auth service" {
	//         type: "service"
	//     }
	// requirementDiagram
	//     direction BT
	//     element "auth service" {
	//         type: "service"
	//     }
	// requirementDiagram
	//     direction LR
	//     element "auth service" {
	//         type: "service"
	//     }
	// requirementDiagram
	//     direction RL
	//     element "auth service" {
	//         type: "service"
	//     }
}

// ExampleClassDefSpec is one named style, as Def returns it.
func ExampleClassDefSpec() {
	specs := []requirement.ClassDefSpec{
		requirement.Def("urgent", "fill:#f96"),
		requirement.Def("done", "fill:#9f6"),
	}

	_ = requirement.NewDiagram(os.Stdout).ClassDefs(specs...).Build()

	// Output:
	// requirementDiagram
	//     classDef urgent fill:#f96
	//     classDef done fill:#9f6
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []requirement.Option{requirement.WithTitle("Login")}

	_ = requirement.NewDiagram(os.Stdout, options...).
		Element("auth service", requirement.WithElementType("service")).
		Build()

	// Output:
	// ---
	// title: "Login"
	// ---
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
}

// ExampleRequirementOption shows what a RequirementOption is: a function that
// changes how a requirement is written, passed after its name.
func ExampleRequirementOption() {
	options := []requirement.RequirementOption{
		requirement.WithID("1"),
		requirement.WithText("The system shall log a user in."),
		requirement.WithRisk(requirement.RiskHigh),
		requirement.WithVerifyMethod(requirement.VerifyMethodTest),
	}

	_ = requirement.NewDiagram(os.Stdout).Requirement("login", options...).Build()

	// Output:
	// requirementDiagram
	//     requirement login {
	//         id: "1"
	//         text: "The system shall log a user in."
	//         risk: High
	//         verifymethod: Test
	//     }
}

// ExampleElementOption shows what an ElementOption is: a function that changes
// how an element is written, passed after its name.
func ExampleElementOption() {
	options := []requirement.ElementOption{requirement.WithElementType("service")}

	_ = requirement.NewDiagram(os.Stdout).Element("auth service", options...).Build()

	// Output:
	// requirementDiagram
	//     element "auth service" {
	//         type: "service"
	//     }
}
