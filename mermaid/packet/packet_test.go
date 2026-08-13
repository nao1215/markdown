package packet

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
			want: "packet",
		},
		{
			name: "new diagram with title",
			opts: []Option{WithTitle("UDP Packet")},
			want: `packet
    title UDP Packet`,
		},
		{
			name:    "new diagram with title including newline",
			opts:    []Option{WithTitle("UDP\nPacket")},
			want:    "packet",
			wantErr: true,
		},
		{
			// The sanitizer eats a bare "<" with the rest of the title;
			// "#60;" decodes to the character.
			name: "new diagram with a bare angle in the title",
			opts: []Option{WithTitle("len < 64")},
			want: `packet
    title len #60; 64`,
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

	d := NewDiagram(b, WithTitle("TCP Packet")).
		Field(0, 15, "Source Port").
		Field(16, 31, "Destination Port").
		Field(32, 63, "Sequence Number").
		Bit(106, "URG").
		Next(5, "Flags")

	if err := d.Build(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := `packet
    title TCP Packet
    0-15: "Source Port"
    16-31: "Destination Port"
    32-63: "Sequence Number"
    106: "URG"
    +5: "Flags"`

	got := strings.ReplaceAll(b.String(), "\r\n", "\n")
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("value is mismatch (-want +got):\n%s", diff)
	}
}

func TestQuoteEscapesSpecialChars(t *testing.T) {
	t.Parallel()

	got := quote("a\\b\rc\nd\te\"f")
	want := `"a#92;b<br/>c<br/>d#92;te&quot;f"`
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

func TestDiagram_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func() *Diagram
		want string
	}{
		{
			name: "negative field start",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Field(-1, 10, "Source Port")
			},
			want: "packet",
		},
		{
			name: "negative field end",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Field(0, -1, "Source Port")
			},
			want: "packet",
		},
		{
			name: "start greater than end",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Field(8, 7, "Source Port")
			},
			want: "packet",
		},
		{
			name: "empty field label",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Field(0, 7, "")
			},
			want: "packet",
		},
		{
			name: "newline in field label",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Field(0, 7, "Source\nPort")
			},
			want: "packet",
		},
		{
			name: "invalid next bit count",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Next(0, "Source Port")
			},
			want: "packet",
		},
		{
			name: "empty next label",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Next(8, "")
			},
			want: "packet",
		},
		{
			name: "newline in title",
			run: func() *Diagram {
				return NewDiagram(io.Discard, WithTitle("UDP\nPacket"))
			},
			want: "packet",
		},
		{
			name: "lf short-circuit after error",
			run: func() *Diagram {
				return NewDiagram(io.Discard).Field(-1, 1, "x").LF()
			},
			want: "packet",
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

	d := NewDiagram(errWriter{})
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

func TestDiagram_BuildNilWriter(t *testing.T) {
	t.Parallel()

	d := NewDiagram(nil)
	err := d.Build()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "output writer must not be nil" {
		t.Fatalf("unexpected error: %v", err)
	}

	if d.Error() == nil {
		t.Fatal("expected persisted error, got nil")
	}
	if !errors.Is(d.Error(), err) {
		t.Fatalf("expected Error() to wrap returned error, got %v", d.Error())
	}
}

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).Field(0, 15, "Source Port")
	})
}

// TestRecordedErrorContract asserts that a bit range the syntax cannot express surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).Field(-1, 15, "Source Port")
	})
}

// TestGoldenPacket pins the rendered diagram of every builder method of this
// package.
func TestGoldenPacket(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(buf, WithTitle("TCP Header")).
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
