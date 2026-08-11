package arch_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/arch"
)

// TestGoldenArchitecture pins the rendered diagram of every builder method of
// this package.
func TestGoldenArchitecture(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := arch.NewArchitecture(buf).
		Group("api", arch.IconCloud, "API").
		GroupInParentGroup("storage", arch.IconDatabase, "Storage", "api").
		Service("gateway", arch.IconInternet, "Gateway").
		ServiceInGroup("db", arch.IconDatabase, "Database", "storage").
		ServiceInGroup("disk", arch.IconDisk, "Disk", "storage").
		ServiceInGroup("worker", arch.IconServer, "Worker", "api").
		Junction("hub").
		JunctionsInParent("inner", "storage").
		LF().
		Edges(
			arch.Edge{ServiceID: "gateway", Position: arch.PositionRight, Arrow: arch.ArrowRight},
			arch.Edge{ServiceID: "hub", Position: arch.PositionLeft, Arrow: arch.ArrowNone},
		).
		Edges(
			arch.Edge{ServiceID: "hub", Position: arch.PositionBottom, Arrow: arch.ArrowNone},
			arch.Edge{ServiceID: "worker", Position: arch.PositionTop, Arrow: arch.ArrowLeft},
		).
		EdgesInAnothorGroup(
			arch.Edge{ServiceID: "worker", Position: arch.PositionRight, Arrow: arch.ArrowNone},
			arch.Edge{ServiceID: "db", Position: arch.PositionLeft, Arrow: arch.ArrowRight},
		).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("architecture.md", buf.String()); err != nil {
		t.Error(err)
	}
}
