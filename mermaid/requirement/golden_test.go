package requirement_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/requirement"
)

// TestGoldenRequirement pins the rendered diagram of every requirement kind,
// every relationship, and every styling method of this package.
//
// Every requirement carries an id and a text because the builder requires both.
func TestGoldenRequirement(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := requirement.NewDiagram(buf, requirement.WithTitle("System Requirements")).
		SetDirection(requirement.DirectionLR).
		Requirement("plain requirement",
			requirement.WithID("1"),
			requirement.WithText("the system shall do the thing"),
			requirement.WithRisk(requirement.RiskLow),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest),
		).
		Requirement("full requirement",
			requirement.WithID("2"),
			requirement.WithText("the system shall do the other thing"),
			requirement.WithRisk(requirement.RiskHigh),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest),
		).
		RequirementOfType(requirement.RequirementTypeRequirement, "typed requirement",
			requirement.WithID("3"),
			requirement.WithText("stated with an explicit type"),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest),
		).
		FunctionalRequirement("functional",
			requirement.WithID("4"),
			requirement.WithText("a functional requirement"),
			requirement.WithRisk(requirement.RiskLow),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest),
		).
		InterfaceRequirement("interface",
			requirement.WithID("5"),
			requirement.WithText("an interface requirement"),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodAnalysis),
		).
		PerformanceRequirement("performance",
			requirement.WithID("6"),
			requirement.WithText("a performance requirement"),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodInspection),
		).
		PhysicalRequirement("physical",
			requirement.WithID("7"),
			requirement.WithText("a physical requirement"),
			requirement.WithRisk(requirement.RiskMedium),
			requirement.WithVerifyMethod(requirement.VerifyMethodDemonstration),
		).
		DesignConstraint("design constraint",
			requirement.WithID("8"),
			requirement.WithText("a design constraint"),
			requirement.WithRisk(requirement.RiskHigh),
			requirement.WithVerifyMethod(requirement.VerifyMethodInspection),
		).
		Requirement("classified requirement",
			requirement.WithID("9"),
			requirement.WithText("a requirement carrying classes"),
			requirement.WithRisk(requirement.RiskLow),
			requirement.WithVerifyMethod(requirement.VerifyMethodAnalysis),
			requirement.WithRequirementClasses("important", "reviewed"),
		).
		LF().
		Element("test suite",
			requirement.WithElementType("simulation"),
			requirement.WithDocRef("./tests"),
		).
		Element("classified element", requirement.WithElementClasses("important")).
		LF().
		Relation("plain requirement", requirement.RelationshipContains, "functional").
		Contains("plain requirement", "interface").
		Copies("functional", "interface").
		Derives("functional", "performance").
		Satisfies("test suite", "functional").
		Verifies("test suite", "interface").
		Refines("performance", "physical").
		Traces("physical", "design constraint").
		LF().
		Style("functional", "fill:#f9f").
		ClassDef("important", "fill:#ffa").
		ClassDefs(
			requirement.Def("reviewed", "stroke:#0a0"),
			requirement.Def("legacy", "stroke-dasharray: 5 5"),
		).
		Class("interface", "important").
		ClassShorthand("performance", "important", "reviewed").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("requirement.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenRequirementSourceRelations pins the fluent relation builder, which
// states the source once and then chains every relationship from it.
func TestGoldenRequirementSourceRelations(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := requirement.NewDiagram(buf).
		Requirement("root",
			requirement.WithID("1"),
			requirement.WithText("the root requirement"),
			requirement.WithRisk(requirement.RiskLow),
			requirement.WithVerifyMethod(requirement.VerifyMethodTest),
		).
		From("root").
		Contains("child").
		Copies("copy").
		Derives("derived").
		Satisfies("satisfied").
		Verifies("verified").
		Refines("refined").
		Traces("traced").
		Relation(requirement.RelationshipContains, "explicit").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("requirement_source_relations.md", buf.String()); err != nil {
		t.Error(err)
	}
}
