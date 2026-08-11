package sequence_test

import (
	"bytes"
	"testing"

	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/sequence"
)

// TestGoldenSequence pins the rendered diagram of every message kind, every
// block construct, and every participant lifecycle method of this package.
func TestGoldenSequence(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := sequence.NewDiagram(buf).
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
	err := sequence.NewDiagram(buf).
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
	err := sequence.NewDiagram(
		buf,
		sequence.WithMirrorActors(true),
		sequence.WithBottomMariginAdjustment(2),
		sequence.WithActorFontSize(16),
		sequence.WithActorFontFamily("Helvetica"),
		sequence.WithActorFontWeight("bold"),
		sequence.WithNoteFontSize(12),
		sequence.WithNoteFontFamily("Courier"),
		sequence.WithNoteFontWeight("normal"),
		sequence.WithNoteAlign("left"),
		sequence.WithMessageFontSize(14),
		sequence.WithMessageFontFamily("Arial"),
		sequence.WithMessageFontWeight("lighter"),
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
