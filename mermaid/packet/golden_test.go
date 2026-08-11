package packet_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/packet"
)

// TestGoldenPacket pins the rendered diagram of every builder method of this
// package.
func TestGoldenPacket(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := packet.NewDiagram(buf, packet.WithTitle("TCP Header")).
		Field(0, 15, "Source Port").
		Field(16, 31, "Destination Port").
		Next(32, "Sequence Number").
		Bit(64, "URG").
		Bit(65, "ACK").
		LF().
		Next(16, "Window").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("packet.md", buf.String()); err != nil {
		t.Error(err)
	}
}
