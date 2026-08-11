package timeline_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/timeline"
)

// TestGoldenTimeline pins the rendered diagram of every builder method and
// every option of this package.
func TestGoldenTimeline(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := timeline.NewDiagram(buf, timeline.WithTitle("History of Social Media")).
		Period("2002", "LinkedIn").
		Section("Second wave").
		Period("2004", "Facebook", "Google").
		Event("Flickr").
		Period("2005", "YouTube").
		LF().
		Section("Third wave").
		Period("2006", "Twitter").
		Period("09:00 stand up").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("timeline.md", buf.String()); err != nil {
		t.Error(err)
	}
}
