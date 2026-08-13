// Package sequence is mermaid sequence diagram builder.
package sequence

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
)

func TestString(t *testing.T) {
	t.Parallel()

	t.Run("should return the sequence diagram body", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.Participant("Alice")

		want := fmt.Sprintf("sequenceDiagram%s    participant Alice", internal.LineFeed())
		got := d.String()

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("expected and actual are different: %s", diff)
		}
	})
}

func TestDiagramRequestf(t *testing.T) {
	t.Parallel()

	t.Run("should add request to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.SyncRequestf("Alice", "Bob", "Hello %s", "Bob")

		want := []string{"sequenceDiagram", "    Alice->>Bob: Hello Bob"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramResponsef(t *testing.T) {
	t.Parallel()

	t.Run("should add response to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.SyncResponsef("Alice", "Bob", "Hello %s", "Alice")

		want := []string{"sequenceDiagram", "    Alice-->>Bob: Hello Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramRequestErrorf(t *testing.T) {
	t.Parallel()

	t.Run("should add request error to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.RequestErrorf("Alice", "Bob", "Hello %s", "Bob")

		want := []string{"sequenceDiagram", "    Alice-xBob: Hello Bob"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramResponseErrorf(t *testing.T) {
	t.Parallel()

	t.Run("should add response error to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.ResponseErrorf("Alice", "Bob", "Hello %s", "Alice")

		want := []string{"sequenceDiagram", "    Alice--xBob: Hello Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAsyncRequestf(t *testing.T) {
	t.Parallel()

	t.Run("should add async request to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AsyncRequestf("Alice", "Bob", "Hello %s", "Bob")

		want := []string{"sequenceDiagram", "    Alice->)Bob: Hello Bob"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAsyncResponsef(t *testing.T) {
	t.Parallel()

	t.Run("should add async response to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AsyncResponsef("Alice", "Bob", "Hello %s", "Alice")

		want := []string{"sequenceDiagram", "    Alice--)Bob: Hello Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramError(t *testing.T) {
	t.Parallel()

	t.Run("should return the error", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.err = fmt.Errorf("error")

		if d.Error().Error() != "error" {
			t.Error("value is mismatch, want error")
		}
	})
}

func TestNewDiagram(t *testing.T) {
	t.Parallel()

	t.Run("with all options", func(t *testing.T) {
		t.Parallel()

		got := NewDiagram(
			io.Discard,
			WithMirrorActors(true),
			WithBottomMariginAdjustment(2),
			WithActorFontSize(12),
			WithActorFontFamily("Arial"),
			WithActorFontWeight("bold"),
			WithNoteFontFamily("Arial"),
			WithNoteFontSize(12),
			WithNoteFontWeight("bold"),
			WithNoteAlign("left"),
			WithMessageFontFamily("Arial"),
			WithMessageFontSize(12),
			WithMessageFontWeight("bold"),
		)

		want := &Diagram{
			body: []string{"sequenceDiagram"},
			dest: io.Discard,
			config: &config{
				mirrorActors:            true,
				bottomMariginAdjustment: 2,
				actorFontSize:           12,
				actorFontFamily:         "Arial",
				actorFontWeight:         "bold",
				noteFontSize:            12,
				noteFontFamily:          "Arial",
				noteFontWeight:          "bold",
				noteAlign:               "left",
				messageFontSize:         12,
				messageFontFamily:       "Arial",
				messageFontWeight:       "bold",
			},
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("value is mismatch, want %v, got %v", want, got)
		}
	})
}

func TestDiagramActivateDeactivate(t *testing.T) {
	t.Parallel()

	t.Run("should add activate and deactivate to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.Activate("Alice")
		d.Deactivate("Alice")

		want := []string{"sequenceDiagram", "    activate Alice", "    deactivate Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramRequestfWithActivation(t *testing.T) {
	t.Parallel()

	t.Run("should add request to the sequence diagram with activation", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.SyncRequestfWithActivation("Alice", "Bob", "Hello %s", "Bob")

		want := []string{"sequenceDiagram", "    Alice->>+Bob: Hello Bob"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramResponsefWithActivation(t *testing.T) {
	t.Parallel()

	t.Run("should add response to the sequence diagram with activation", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.SyncResponsefWithActivation("Alice", "Bob", "Hello %s", "Alice")

		want := []string{"sequenceDiagram", "    Alice-->>-Bob: Hello Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAsyncRequestfWithActivation(t *testing.T) {
	t.Parallel()

	t.Run("should add async request to the sequence diagram with activation", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AsyncRequestfWithActivation("Alice", "Bob", "Hello %s", "Bob")

		want := []string{"sequenceDiagram", "    Alice->>+Bob: Hello Bob"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAsyncResponsefWithActivation(t *testing.T) {
	t.Parallel()

	t.Run("should add async response to the sequence diagram with activation", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AsyncResponsefWithActivation("Alice", "Bob", "Hello %s", "Alice")

		want := []string{"sequenceDiagram", "    Alice-->>-Bob: Hello Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramParticipant(t *testing.T) {
	t.Parallel()

	t.Run("should add participant to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.Participant("Alice")

		want := []string{"sequenceDiagram", "    participant Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramCreateDeleteParticipant(t *testing.T) {
	t.Parallel()

	t.Run("should add create and delete participant to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.CreateParticipant("Alice")
		d.DestroyParticipant("Alice")

		want := []string{"sequenceDiagram", "    create participant Alice", "    destroy Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramCreateDeleteActor(t *testing.T) {
	t.Parallel()

	t.Run("should add create and delete actor to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.CreateActor("Alice")
		d.DestroyActor("Alice")

		want := []string{"sequenceDiagram", "    create actor Alice", "    destroy Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramActor(t *testing.T) {
	t.Parallel()

	t.Run("should add actor to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.Actor("Alice")

		want := []string{"sequenceDiagram", "    actor Alice"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAutoNumber(t *testing.T) {
	t.Parallel()

	t.Run("should add autonumber to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AutoNumber()

		want := []string{"sequenceDiagram", "    autonumber"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramBoxStartEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add box to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.Participant("Alice").Participant("Bob")
		d.BoxStart([]string{"Alice", "Bob"})
		d.BoxEnd()

		want := []string{
			"sequenceDiagram",
			"    participant Alice",
			"    participant Bob",
			"    box Alice & Bob",
			"    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewDiagram(w).SyncRequest("Client", "Server", "GET /users")
	})
}

// TestGoldenSequence pins the rendered diagram of every message kind, every
// block construct, and every participant lifecycle method of this package.
func TestGoldenSequence(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(buf).
		AutoNumber().
		BoxStart([]string{"Client", "Server"}).
		Participant("Client").
		Participant("Server").
		BoxEnd().
		Actor("Operator").
		LF().
		SyncRequest("Client", "Server", "GET /users").
		SyncRequestf("Client", "Server", "GET /users/%d", 1).
		SyncResponse("Server", "Client", "200 OK").
		SyncResponsef("Server", "Client", "%d OK", 200).
		AsyncRequest("Client", "Server", "publish event").
		AsyncRequestf("Client", "Server", "publish %s", "event").
		AsyncResponse("Server", "Client", "ack").
		AsyncResponsef("Server", "Client", "%s", "ack").
		RequestError("Client", "Server", "malformed body").
		RequestErrorf("Client", "Server", "malformed %s", "body").
		ResponseError("Server", "Client", "500 Internal Server Error").
		ResponseErrorf("Server", "Client", "%d Internal Server Error", 500).
		LF().
		Activate("Server").
		SyncRequestWithActivation("Client", "Server", "activate on request").
		SyncRequestfWithActivation("Client", "Server", "activate on %s", "request").
		SyncResponseWithActivation("Server", "Client", "deactivate on response").
		SyncResponsefWithActivation("Server", "Client", "deactivate on %s", "response").
		AsyncRequestWithActivation("Client", "Server", "async activate").
		AsyncRequestfWithActivation("Client", "Server", "async %s", "activate").
		AsyncResponseWithActivation("Server", "Client", "async deactivate").
		AsyncResponsefWithActivation("Server", "Client", "async %s", "deactivate").
		Deactivate("Server").
		LF().
		NoteOver("Server", "a note over the participant").
		NoteRightOf("Server", "a note to the right").
		NoteLeftOf("Client", "a note to the left").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("sequence.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenSequenceBlocks pins the block constructs, which have to nest and
// close in the right order to render at all.
func TestGoldenSequenceBlocks(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(buf).
		CreateParticipant("Worker").
		CreateActor("Auditor").
		LoopStart("every minute").
		SyncRequest("Client", "Worker", "poll").
		LoopEnd().
		AltStart("the queue is empty").
		SyncResponse("Worker", "Client", "nothing to do").
		AltElse("the queue has work").
		SyncResponse("Worker", "Client", "one job").
		AltEnd().
		OptStart("the caller asked for details").
		SyncResponse("Worker", "Client", "job details").
		OptEnd().
		ParallelStart("fan out").
		SyncRequest("Worker", "Auditor", "record start").
		ParallelAnd("and").
		SyncRequest("Worker", "Auditor", "record end").
		ParallelEnd().
		CriticalStart("acquire the lock").
		SyncRequest("Worker", "Auditor", "lock").
		CriticalOption("the lock is taken").
		SyncRequest("Worker", "Auditor", "wait").
		CriticalEnd().
		BreakStart("the job failed").
		SyncResponse("Worker", "Client", "error").
		BreakEnd().
		DestroyActor("Auditor").
		DestroyParticipant("Worker").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("sequence_blocks.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestGoldenSequenceOptions pins the configuration block that the construction
// time options produce.
func TestGoldenSequenceOptions(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewDiagram(
		buf,
		WithMirrorActors(true),
		WithBottomMariginAdjustment(2),
		WithActorFontSize(16),
		WithActorFontFamily("Helvetica"),
		WithActorFontWeight("bold"),
		WithNoteFontSize(12),
		WithNoteFontFamily("Courier"),
		WithNoteFontWeight("normal"),
		WithNoteAlign("left"),
		WithMessageFontSize(14),
		WithMessageFontFamily("Arial"),
		WithMessageFontWeight("lighter"),
	).
		SyncRequest("Client", "Server", "GET /users").
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("sequence_options.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestBuildWithNilWriter covers the case where a diagram is built for String()
// only and Build() is called by mistake. Build() used to dereference the nil
// writer and take the process down; it has to return an error instead.
func TestBuildWithNilWriter(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Build() panicked with a nil writer: %v", r)
		}
	}()

	d := NewDiagram(nil)

	// String() has always worked without a writer, and callers rely on it.
	_ = d.String()

	err := d.Build()
	if err == nil {
		t.Fatal("Build() with a nil writer must return an error")
	}
	if err.Error() != "output writer must not be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDiagramNoteOver(t *testing.T) {
	t.Parallel()

	t.Run("should add note over to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.NoteOver("Alice", "Hello Alice")

		want := []string{"sequenceDiagram", "    note over Alice: Hello Alice"}
		got := d.body

		if !reflect.DeepEqual(want, got) {
			t.Errorf("value is mismatch want:%v got:%v", want, got)
		}
	})
}

func TestDiagramNoteRightOf(t *testing.T) {
	t.Parallel()

	t.Run("should add note right of to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.NoteRightOf("Alice", "Hello Alice")

		want := []string{"sequenceDiagram", "    note right of Alice: Hello Alice"}
		got := d.body

		if !reflect.DeepEqual(want, got) {
			t.Errorf("value is mismatch want:%v got:%v", want, got)
		}
	})
}

func TestDiagramNoteLeftOf(t *testing.T) {
	t.Parallel()

	t.Run("should add note left of to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.NoteLeftOf("Alice", "Hello Alice")

		want := []string{"sequenceDiagram", "    note left of Alice: Hello Alice"}
		got := d.body

		if !reflect.DeepEqual(want, got) {
			t.Errorf("value is mismatch want:%v got:%v", want, got)
		}
	})
}

func TestDiagramLoopStartEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add loop to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.LoopStart("description")
		d.LoopEnd()

		want := []string{"sequenceDiagram", "    loop description", "    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramAltStartElseEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add alt to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.AltStart("description")
		d.AltElse("description")
		d.AltEnd()

		want := []string{"sequenceDiagram", "    alt description", "    else description", "    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramOptStartEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add opt to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.OptStart("description")
		d.OptEnd()

		want := []string{"sequenceDiagram", "    opt description", "    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramParallelStartAndEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add parallel to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.ParallelStart("start")
		d.ParallelAnd("and-description")
		d.ParallelEnd()

		want := []string{"sequenceDiagram", "    par start", "    and and-description", "    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramCriticalStartAndEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add critical to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.CriticalStart("start")
		d.CriticalOption("option-description")
		d.CriticalEnd()

		want := []string{
			"sequenceDiagram",
			"    critical start",
			"    option option-description",
			"    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestDiagramBreakStartEnd(t *testing.T) {
	t.Parallel()

	t.Run("should add break to the sequence diagram", func(t *testing.T) {
		t.Parallel()

		d := NewDiagram(io.Discard)
		d.BreakStart("description")
		d.BreakEnd()

		want := []string{"sequenceDiagram", "    break description", "    end"}
		got := d.body

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
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

// TestMessageTextEscapesWhatEndsIt names the characters this escaping buys in
// the text of a diagram. A sequence diagram takes no quoted text at all, so a
// message, a note and a block description each go out bare.
func TestMessageTextEscapesWhatEndsIt(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(io.Writer) *Diagram
		want  string
	}{
		"a semicolon in a message loses the diagram": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).SyncRequest("A", "B", "a;b") },
			want:  "    A->>B: a#59;b",
		},
		"a hash in a message cuts it short": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).SyncRequest("A", "B", "deploy #2 of 3")
			},
			want: "    A->>B: deploy #35;2 of 3",
		},
		"a hash in a note": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).NoteOver("A", "run #1") },
			want:  "    note over A: run #35;1",
		},
		"a semicolon in a loop description": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).LoopStart("retry; twice") },
			want:  "    loop retry#59; twice",
		},
		"a colon in a message is left alone": {
			// Only the first colon on the line is syntax, and this package
			// writes that one itself.
			build: func(w io.Writer) *Diagram { return NewDiagram(w).SyncRequest("A", "B", "a:b") },
			want:  "    A->>B: a:b",
		},
		"a percent pair in a message is left alone": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).SyncRequest("A", "B", "100%% done") },
			want:  "    A->>B: 100%% done",
		},
		"a line break in a message becomes the entity": {
			// A raw line break splits the statement and loses the diagram;
			// "#10;" was measured to decode back into a real line break.
			build: func(w io.Writer) *Diagram { return NewDiagram(w).SyncRequest("A", "B", "a\nb") },
			want:  "    A->>B: a#10;b",
		},
		"a CRLF pair in a note is one line break": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).NoteOver("A", "a\r\nb") },
			want:  "    note over A: a#10;b",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.build(io.Discard).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("diagram =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}

// TestParticipantNameEscapesWhatEndsIt names the characters a participant's
// name loses, which are a different set from the text above. The name is
// escaped identically where it is declared and where a message refers to it, so
// the two still match.
func TestParticipantNameEscapesWhatEndsIt(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(io.Writer) *Diagram
		want  string
	}{
		"a hyphen would be read as an arrow": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Participant("web-server") },
			want:  "    participant web#45;server",
		},
		"the declaration and the reference agree": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).Participant("web-server").SyncRequest("web-server", "db", "read")
			},
			want: "    web#45;server->>db: read",
		},
		"a colon and a comma in an actor": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Actor("Ops: EU, US") },
			want:  "    actor Ops#58; EU#44; US",
		},
		"parentheses in a name": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Participant("Deploy (prod)") },
			want:  "    participant Deploy #40;prod#41;",
		},
		"an angle bracket in a name": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Participant("a<b") },
			want:  "    participant a#60;b",
		},
		"a percent pair in a name": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Participant("a%%b") },
			want:  "    participant a#37;#37;b",
		},
		"a lone percent in a name is left alone": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Participant("50% traffic") },
			want:  "    participant 50% traffic",
		},
		"a hash in a name is left alone": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Participant("PR #12") },
			want:  "    participant PR #12",
		},
		"a comma between two participants a note spans is a separator": {
			// A note may be placed over two participants at once, and this
			// package cannot tell that apart from one name holding a comma.
			build: func(w io.Writer) *Diagram { return NewDiagram(w).NoteOver("Alice,Bob", "m") },
			want:  "    note over Alice,Bob: m",
		},
		"a plus in a name could be declared but never messaged": {
			build: func(w io.Writer) *Diagram {
				return NewDiagram(w).Participant("a+b").SyncRequest("a+b", "c", "read")
			},
			want: "    a#43;b->>c: read",
		},
		"an at sign in a name cannot even be declared": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Participant("ops@eu") },
			want:  "    participant ops#64;eu",
		},
		"a line break in a name becomes the entity": {
			build: func(w io.Writer) *Diagram { return NewDiagram(w).Participant("a\nb") },
			want:  "    participant a#10;b",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.build(io.Discard).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("diagram =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}
