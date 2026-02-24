//go:build linux || darwin

package requirement_test

import (
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
	// title: Checkout Requirements
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
