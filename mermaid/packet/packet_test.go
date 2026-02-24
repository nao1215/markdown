package packet

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
	want := `"a&#92;b&#92;rc&#92;nd&#92;te&quot;f"`
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
