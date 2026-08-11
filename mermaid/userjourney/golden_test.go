package userjourney_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/userjourney"
)

// TestGoldenUserJourney pins the rendered diagram of every builder method of
// this package, including every score.
func TestGoldenUserJourney(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := userjourney.NewDiagram(buf, userjourney.WithTitle("Sign Up")).
		Section("Discovery").
		Task("Find the site", userjourney.ScoreVerySatisfied, "Visitor").
		Task("Read the docs", userjourney.ScoreSatisfied, "Visitor", "Support").
		Section("Registration").
		Task("Fill the form", userjourney.ScoreNeutral, "Visitor").
		Task("Confirm the mail", userjourney.ScoreDissatisfied, "Visitor").
		Task("Wait for approval", userjourney.ScoreVeryDissatisfied).
		LF().
		Section("Onboarding").
		TaskIn("Onboarding", "First login", userjourney.ScoreSatisfied, "User").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("userjourney.md", buf.String()); err != nil {
		t.Error(err)
	}
}
