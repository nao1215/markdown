//go:build linux || darwin

package sequence_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/sequence"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := sequence.NewDiagram(os.Stdout).
		Participant("Sophia").
		Participant("David").
		Participant("Subaru").
		LF().
		SyncRequest("Sophia", "David", "Please wake up Subaru").
		SyncResponse("David", "Sophia", "OK").
		LF().
		LoopStart("until Subaru wake up").
		SyncRequest("David", "Subaru", "Wake up!").
		SyncResponse("Subaru", "David", "zzz").
		SyncRequest("David", "Subaru", "Hey!!!").
		BreakStart("if Subaru wake up").
		SyncResponse("Subaru", "David", "......").
		BreakEnd().
		LoopEnd().
		LF().
		SyncResponse("David", "Sophia", "wake up, wake up").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Sequence Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Sequence Diagram
	// ```mermaid
	// sequenceDiagram
	//     participant Sophia
	//     participant David
	//     participant Subaru
	//
	//     Sophia->>David: Please wake up Subaru
	//     David-->>Sophia: OK
	//
	//     loop until Subaru wake up
	//     David->>Subaru: Wake up!
	//     Subaru-->>David: zzz
	//     David->>Subaru: Hey!!!
	//     break if Subaru wake up
	//     Subaru-->>David: ......
	//     end
	//     end
	//
	//     David-->>Sophia: wake up, wake up
	// ```
}

// ExampleDiagram_SyncRequest draws a solid arrow with a filled head, the ordinary call.
func ExampleDiagram_SyncRequest() {
	_ = sequence.NewDiagram(os.Stdout).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleDiagram_SyncRequestf draws the same from a format string.
func ExampleDiagram_SyncRequestf() {
	_ = sequence.NewDiagram(os.Stdout).
		SyncRequestf("Alice", "Bob", "Retry %d of %d", 2, 3).
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: Retry 2 of 3
}

// ExampleDiagram_SyncResponse draws a dashed arrow with a filled head, the ordinary reply.
func ExampleDiagram_SyncResponse() {
	_ = sequence.NewDiagram(os.Stdout).
		SyncResponse("Bob", "Alice", "Fine, thanks").
		Build()

	// Output:
	// sequenceDiagram
	//     Bob-->>Alice: Fine, thanks
}

// ExampleDiagram_SyncResponsef draws the same from a format string.
func ExampleDiagram_SyncResponsef() {
	_ = sequence.NewDiagram(os.Stdout).
		SyncResponsef("Bob", "Alice", "%d items", 7).
		Build()

	// Output:
	// sequenceDiagram
	//     Bob-->>Alice: 7 items
}

// ExampleDiagram_AsyncRequest draws a call that does not wait for a reply.
func ExampleDiagram_AsyncRequest() {
	_ = sequence.NewDiagram(os.Stdout).
		AsyncRequest("Alice", "Bob", "Start the job").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->)Bob: Start the job
}

// ExampleDiagram_AsyncRequestf draws the same from a format string.
func ExampleDiagram_AsyncRequestf() {
	_ = sequence.NewDiagram(os.Stdout).
		AsyncRequestf("Alice", "Bob", "Start job %d", 7).
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->)Bob: Start job 7
}

// ExampleDiagram_AsyncResponse draws a reply to a call that was not waiting.
func ExampleDiagram_AsyncResponse() {
	_ = sequence.NewDiagram(os.Stdout).
		AsyncResponse("Bob", "Alice", "Job finished").
		Build()

	// Output:
	// sequenceDiagram
	//     Bob--)Alice: Job finished
}

// ExampleDiagram_AsyncResponsef draws the same from a format string.
func ExampleDiagram_AsyncResponsef() {
	_ = sequence.NewDiagram(os.Stdout).
		AsyncResponsef("Bob", "Alice", "Job %d finished", 7).
		Build()

	// Output:
	// sequenceDiagram
	//     Bob--)Alice: Job 7 finished
}

// ExampleDiagram_RequestError draws a call that failed, drawn with a cross at the end.
func ExampleDiagram_RequestError() {
	_ = sequence.NewDiagram(os.Stdout).
		RequestError("Alice", "Bob", "Timed out").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice-xBob: Timed out
}

// ExampleDiagram_RequestErrorf draws the same from a format string.
func ExampleDiagram_RequestErrorf() {
	_ = sequence.NewDiagram(os.Stdout).
		RequestErrorf("Alice", "Bob", "Timed out after %ds", 30).
		Build()

	// Output:
	// sequenceDiagram
	//     Alice-xBob: Timed out after 30s
}

// ExampleDiagram_ResponseError draws a reply that failed.
func ExampleDiagram_ResponseError() {
	_ = sequence.NewDiagram(os.Stdout).
		ResponseError("Bob", "Alice", "Refused").
		Build()

	// Output:
	// sequenceDiagram
	//     Bob--xAlice: Refused
}

// ExampleDiagram_ResponseErrorf draws the same from a format string.
func ExampleDiagram_ResponseErrorf() {
	_ = sequence.NewDiagram(os.Stdout).
		ResponseErrorf("Bob", "Alice", "Refused with %d", 503).
		Build()

	// Output:
	// sequenceDiagram
	//     Bob--xAlice: Refused with 503
}

// ExampleDiagram_SyncRequestWithActivation draws a call that also turns the receiver's activation bar on.
func ExampleDiagram_SyncRequestWithActivation() {
	_ = sequence.NewDiagram(os.Stdout).
		SyncRequestWithActivation("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>+Bob: How are you?
}

// ExampleDiagram_SyncRequestfWithActivation draws the same from a format string.
func ExampleDiagram_SyncRequestfWithActivation() {
	_ = sequence.NewDiagram(os.Stdout).
		SyncRequestfWithActivation("Alice", "Bob", "Retry %d", 2).
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>+Bob: Retry 2
}

// ExampleDiagram_SyncResponseWithActivation draws a reply that also turns the bar off.
func ExampleDiagram_SyncResponseWithActivation() {
	_ = sequence.NewDiagram(os.Stdout).
		SyncResponseWithActivation("Bob", "Alice", "Fine, thanks").
		Build()

	// Output:
	// sequenceDiagram
	//     Bob-->>-Alice: Fine, thanks
}

// ExampleDiagram_SyncResponsefWithActivation draws the same from a format string.
func ExampleDiagram_SyncResponsefWithActivation() {
	_ = sequence.NewDiagram(os.Stdout).
		SyncResponsefWithActivation("Bob", "Alice", "%d items", 7).
		Build()

	// Output:
	// sequenceDiagram
	//     Bob-->>-Alice: 7 items
}

// ExampleDiagram_AsyncRequestWithActivation draws an asynchronous call that turns the bar on.
func ExampleDiagram_AsyncRequestWithActivation() {
	_ = sequence.NewDiagram(os.Stdout).
		AsyncRequestWithActivation("Alice", "Bob", "Start the job").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>+Bob: Start the job
}

// ExampleDiagram_AsyncRequestfWithActivation draws the same from a format string.
func ExampleDiagram_AsyncRequestfWithActivation() {
	_ = sequence.NewDiagram(os.Stdout).
		AsyncRequestfWithActivation("Alice", "Bob", "Start job %d", 7).
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>+Bob: Start job 7
}

// ExampleDiagram_AsyncResponseWithActivation draws an asynchronous reply that turns the bar off.
func ExampleDiagram_AsyncResponseWithActivation() {
	_ = sequence.NewDiagram(os.Stdout).
		AsyncResponseWithActivation("Bob", "Alice", "Job finished").
		Build()

	// Output:
	// sequenceDiagram
	//     Bob-->>-Alice: Job finished
}

// ExampleDiagram_AsyncResponsefWithActivation draws the same from a format string.
func ExampleDiagram_AsyncResponsefWithActivation() {
	_ = sequence.NewDiagram(os.Stdout).
		AsyncResponsefWithActivation("Bob", "Alice", "Job %d finished", 7).
		Build()

	// Output:
	// sequenceDiagram
	//     Bob-->>-Alice: Job 7 finished
}

// ExampleNewDiagram shows the shape every sequence diagram has: a writer, the
// participants, the messages between them, and Build.
func ExampleNewDiagram() {
	_ = sequence.NewDiagram(os.Stdout).
		Participant("Alice").
		Participant("Bob").
		SyncRequest("Alice", "Bob", "How are you?").
		SyncResponse("Bob", "Alice", "Fine, thanks").
		Build()

	// Output:
	// sequenceDiagram
	//     participant Alice
	//     participant Bob
	//     Alice->>Bob: How are you?
	//     Bob-->>Alice: Fine, thanks
}

// ExampleDiagram_Participant declares someone taking part, drawn as a box.
// Declaring one is what fixes the order they appear in; a participant a message
// names is drawn anyway, at the point it is first mentioned.
func ExampleDiagram_Participant() {
	_ = sequence.NewDiagram(os.Stdout).
		Participant("Bob").
		Participant("Alice").
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     participant Bob
	//     participant Alice
	//     Alice->>Bob: How are you?
}

// ExampleDiagram_Actor declares someone taking part, drawn as a stick figure
// rather than a box.
func ExampleDiagram_Actor() {
	_ = sequence.NewDiagram(os.Stdout).
		Actor("Alice").
		Participant("Bob").
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     actor Alice
	//     participant Bob
	//     Alice->>Bob: How are you?
}

// ExampleDiagram_CreateParticipant brings a participant into being partway
// through, which is how a diagram shows something being started rather than
// having been there all along.
func ExampleDiagram_CreateParticipant() {
	_ = sequence.NewDiagram(os.Stdout).
		Participant("Alice").
		CreateParticipant("Worker").
		SyncRequest("Alice", "Worker", "Start").
		Build()

	// Output:
	// sequenceDiagram
	//     participant Alice
	//     create participant Worker
	//     Alice->>Worker: Start
}

// ExampleDiagram_DestroyParticipant ends a participant, drawn with a cross.
func ExampleDiagram_DestroyParticipant() {
	_ = sequence.NewDiagram(os.Stdout).
		Participant("Alice").
		Participant("Worker").
		SyncRequest("Alice", "Worker", "Stop").
		DestroyParticipant("Worker").
		Build()

	// Output:
	// sequenceDiagram
	//     participant Alice
	//     participant Worker
	//     Alice->>Worker: Stop
	//     destroy Worker
}

// ExampleDiagram_CreateActor brings an actor into being partway through.
func ExampleDiagram_CreateActor() {
	_ = sequence.NewDiagram(os.Stdout).
		Participant("System").
		CreateActor("Reviewer").
		SyncRequest("System", "Reviewer", "Please review").
		Build()

	// Output:
	// sequenceDiagram
	//     participant System
	//     create actor Reviewer
	//     System->>Reviewer: Please review
}

// ExampleDiagram_DestroyActor ends an actor.
func ExampleDiagram_DestroyActor() {
	_ = sequence.NewDiagram(os.Stdout).
		Actor("Reviewer").
		DestroyActor("Reviewer").
		Build()

	// Output:
	// sequenceDiagram
	//     actor Reviewer
	//     destroy Reviewer
}

// ExampleDiagram_Activate turns a participant's activation bar on, which shows
// it doing something rather than waiting.
func ExampleDiagram_Activate() {
	_ = sequence.NewDiagram(os.Stdout).
		SyncRequest("Alice", "Bob", "How are you?").
		Activate("Bob").
		SyncResponse("Bob", "Alice", "Fine, thanks").
		Deactivate("Bob").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
	//     activate Bob
	//     Bob-->>Alice: Fine, thanks
	//     deactivate Bob
}

// ExampleDiagram_Deactivate turns the bar off again.
func ExampleDiagram_Deactivate() {
	_ = sequence.NewDiagram(os.Stdout).
		Activate("Bob").
		Deactivate("Bob").
		Build()

	// Output:
	// sequenceDiagram
	//     activate Bob
	//     deactivate Bob
}

// ExampleDiagram_NoteOver puts a note across one or more participants. Two
// named with a comma between them is a note spanning both.
func ExampleDiagram_NoteOver() {
	_ = sequence.NewDiagram(os.Stdout).
		NoteOver("Alice", "Thinking it over").
		NoteOver("Alice,Bob", "Both are waiting").
		Build()

	// Output:
	// sequenceDiagram
	//     note over Alice: Thinking it over
	//     note over Alice,Bob: Both are waiting
}

// ExampleDiagram_NoteRightOf puts a note to the right of a participant.
func ExampleDiagram_NoteRightOf() {
	_ = sequence.NewDiagram(os.Stdout).
		NoteRightOf("Bob", "Bob is busy").
		Build()

	// Output:
	// sequenceDiagram
	//     note right of Bob: Bob is busy
}

// ExampleDiagram_NoteLeftOf puts a note to the left of a participant.
func ExampleDiagram_NoteLeftOf() {
	_ = sequence.NewDiagram(os.Stdout).
		NoteLeftOf("Alice", "Alice is waiting").
		Build()

	// Output:
	// sequenceDiagram
	//     note left of Alice: Alice is waiting
}

// ExampleDiagram_LoopStart opens a block that repeats, and LoopEnd closes it.
func ExampleDiagram_LoopStart() {
	_ = sequence.NewDiagram(os.Stdout).
		LoopStart("every minute").
		SyncRequest("Alice", "Bob", "Still there?").
		LoopEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     loop every minute
	//     Alice->>Bob: Still there?
	//     end
}

// ExampleDiagram_LoopEnd closes the block LoopStart opened.
func ExampleDiagram_LoopEnd() {
	_ = sequence.NewDiagram(os.Stdout).
		LoopStart("every minute").
		SyncRequest("Alice", "Bob", "Still there?").
		LoopEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     loop every minute
	//     Alice->>Bob: Still there?
	//     end
}

// ExampleDiagram_AltStart opens a choice between paths. AltElse begins each
// alternative and AltEnd closes the whole thing.
func ExampleDiagram_AltStart() {
	_ = sequence.NewDiagram(os.Stdout).
		AltStart("is logged in").
		SyncResponse("Bob", "Alice", "Here you are").
		AltElse("is not").
		SyncResponse("Bob", "Alice", "Please log in").
		AltEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     alt is logged in
	//     Bob-->>Alice: Here you are
	//     else is not
	//     Bob-->>Alice: Please log in
	//     end
}

// ExampleDiagram_AltElse begins an alternative path.
func ExampleDiagram_AltElse() {
	_ = sequence.NewDiagram(os.Stdout).
		AltStart("is logged in").
		AltElse("is not").
		AltEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     alt is logged in
	//     else is not
	//     end
}

// ExampleDiagram_AltEnd closes the choice AltStart opened.
func ExampleDiagram_AltEnd() {
	_ = sequence.NewDiagram(os.Stdout).
		AltStart("is logged in").
		AltEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     alt is logged in
	//     end
}

// ExampleDiagram_OptStart opens a block that may or may not happen, which is a
// choice with only one path.
func ExampleDiagram_OptStart() {
	_ = sequence.NewDiagram(os.Stdout).
		OptStart("if the cache is cold").
		SyncRequest("Bob", "Database", "Fetch").
		OptEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     opt if the cache is cold
	//     Bob->>Database: Fetch
	//     end
}

// ExampleDiagram_OptEnd closes the block OptStart opened.
func ExampleDiagram_OptEnd() {
	_ = sequence.NewDiagram(os.Stdout).
		OptStart("if the cache is cold").
		OptEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     opt if the cache is cold
	//     end
}

// ExampleDiagram_ParallelStart opens paths that run at the same time.
// ParallelAnd begins each one after the first and ParallelEnd closes them.
func ExampleDiagram_ParallelStart() {
	_ = sequence.NewDiagram(os.Stdout).
		ParallelStart("lint").
		SyncRequest("CI", "Linter", "Run").
		ParallelAnd("test").
		SyncRequest("CI", "Tester", "Run").
		ParallelEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     par lint
	//     CI->>Linter: Run
	//     and test
	//     CI->>Tester: Run
	//     end
}

// ExampleDiagram_ParallelAnd begins another path running at the same time.
func ExampleDiagram_ParallelAnd() {
	_ = sequence.NewDiagram(os.Stdout).
		ParallelStart("lint").
		ParallelAnd("test").
		ParallelEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     par lint
	//     and test
	//     end
}

// ExampleDiagram_ParallelEnd closes the paths ParallelStart opened.
func ExampleDiagram_ParallelEnd() {
	_ = sequence.NewDiagram(os.Stdout).
		ParallelStart("lint").
		ParallelEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     par lint
	//     end
}

// ExampleDiagram_CriticalStart opens a block that has to happen, with
// CriticalOption for each thing that can go wrong instead.
func ExampleDiagram_CriticalStart() {
	_ = sequence.NewDiagram(os.Stdout).
		CriticalStart("connect to the database").
		SyncRequest("Service", "Database", "Connect").
		CriticalOption("the network is down").
		SyncRequest("Service", "Operator", "Page").
		CriticalEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     critical connect to the database
	//     Service->>Database: Connect
	//     option the network is down
	//     Service->>Operator: Page
	//     end
}

// ExampleDiagram_CriticalOption begins what happens instead when the critical
// block cannot go ahead.
func ExampleDiagram_CriticalOption() {
	_ = sequence.NewDiagram(os.Stdout).
		CriticalStart("connect to the database").
		CriticalOption("the network is down").
		CriticalEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     critical connect to the database
	//     option the network is down
	//     end
}

// ExampleDiagram_CriticalEnd closes the block CriticalStart opened.
func ExampleDiagram_CriticalEnd() {
	_ = sequence.NewDiagram(os.Stdout).
		CriticalStart("connect to the database").
		CriticalEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     critical connect to the database
	//     end
}

// ExampleDiagram_BreakStart opens a block that stops the flow when it happens.
func ExampleDiagram_BreakStart() {
	_ = sequence.NewDiagram(os.Stdout).
		BreakStart("the request is invalid").
		SyncResponse("Bob", "Alice", "400 Bad Request").
		BreakEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     break the request is invalid
	//     Bob-->>Alice: 400 Bad Request
	//     end
}

// ExampleDiagram_BreakEnd closes the block BreakStart opened.
func ExampleDiagram_BreakEnd() {
	_ = sequence.NewDiagram(os.Stdout).
		BreakStart("the request is invalid").
		BreakEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     break the request is invalid
	//     end
}

// ExampleDiagram_BoxStart draws a box around several participants, which is how
// a diagram says which of them belong to one system.
func ExampleDiagram_BoxStart() {
	_ = sequence.NewDiagram(os.Stdout).
		BoxStart([]string{"Alice", "Bob"}).
		Participant("Alice").
		Participant("Bob").
		BoxEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     box Alice & Bob
	//     participant Alice
	//     participant Bob
	//     end
}

// ExampleDiagram_BoxEnd closes the box BoxStart opened.
func ExampleDiagram_BoxEnd() {
	_ = sequence.NewDiagram(os.Stdout).
		BoxStart([]string{"Alice"}).
		Participant("Alice").
		BoxEnd().
		Build()

	// Output:
	// sequenceDiagram
	//     box Alice
	//     participant Alice
	//     end
}

// ExampleDiagram_AutoNumber numbers the messages, so a discussion can refer to
// one of them by number.
func ExampleDiagram_AutoNumber() {
	_ = sequence.NewDiagram(os.Stdout).
		AutoNumber().
		SyncRequest("Alice", "Bob", "How are you?").
		SyncResponse("Bob", "Alice", "Fine, thanks").
		Build()

	// Output:
	// sequenceDiagram
	//     autonumber
	//     Alice->>Bob: How are you?
	//     Bob-->>Alice: Fine, thanks
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := sequence.NewDiagram(io.Discard).
		SyncRequest("Alice", "Bob", "How are you?").
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// sequenceDiagram
	//     Alice->>Bob: How are you?
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the error the chain
// recorded.
func ExampleDiagram_Build() {
	err := sequence.NewDiagram(nil).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := sequence.NewDiagram(io.Discard).SyncRequest("Alice", "Bob", "How are you?")
	fmt.Println("error:", d.Error())

	// Output:
	// error: <nil>
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = sequence.NewDiagram(os.Stdout).
		SyncRequest("Alice", "Bob", "How are you?").
		LF().
		SyncResponse("Bob", "Alice", "Fine, thanks").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
	//
	//     Bob-->>Alice: Fine, thanks
}

// ExampleWithMirrorActors draws the actors along the bottom as well as the top,
// which a long diagram wants so a reader does not have to scroll back.
func ExampleWithMirrorActors() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithMirrorActors(true)).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleDiagram_AutoNumber_second shows that numbering is a call rather than
// an option, so it can be turned on partway through a chain.
func ExampleDiagram_AutoNumber_second() {
	_ = sequence.NewDiagram(os.Stdout).
		AutoNumber().
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     autonumber
	//     Alice->>Bob: How are you?
}

// ExampleNotePosition shows where a note can be placed. The constants are used
// by the note methods rather than passed to them, and WithNoteAlign takes the
// same wording as a string.
func ExampleNotePosition() {
	fmt.Println(sequence.NotePositionOver)
	fmt.Println(sequence.NotePositionLeft)
	fmt.Println(sequence.NotePositionRight)

	// Output:
	// over
	// left of
	// right of
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []sequence.Option{
		sequence.WithMirrorActors(true),
		sequence.WithActorFontSize(18),
	}

	_ = sequence.NewDiagram(os.Stdout, options...).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleWithActorFontSize sets the size of the actor names.
func ExampleWithActorFontSize() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithActorFontSize(18)).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleWithActorFontFamily sets the typeface of the actor names.
func ExampleWithActorFontFamily() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithActorFontFamily("Helvetica")).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleWithActorFontWeight sets the weight of the actor names.
func ExampleWithActorFontWeight() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithActorFontWeight("bold")).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleWithMessageFontSize sets the size of the message text.
func ExampleWithMessageFontSize() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithMessageFontSize(14)).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleWithMessageFontFamily sets the typeface of the message text.
func ExampleWithMessageFontFamily() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithMessageFontFamily("Helvetica")).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleWithMessageFontWeight sets the weight of the message text.
func ExampleWithMessageFontWeight() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithMessageFontWeight("bold")).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleWithNoteFontSize sets the size of the note text.
func ExampleWithNoteFontSize() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithNoteFontSize(12)).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleWithNoteFontFamily sets the typeface of the note text.
func ExampleWithNoteFontFamily() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithNoteFontFamily("Helvetica")).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleWithNoteFontWeight sets the weight of the note text.
func ExampleWithNoteFontWeight() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithNoteFontWeight("bold")).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleWithNoteAlign sets which way the note text is aligned.
func ExampleWithNoteAlign() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithNoteAlign("center")).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}

// ExampleWithBottomMariginAdjustment adds space below the diagram.
func ExampleWithBottomMariginAdjustment() {
	_ = sequence.NewDiagram(os.Stdout, sequence.WithBottomMariginAdjustment(10)).
		SyncRequest("Alice", "Bob", "How are you?").
		Build()

	// Output:
	// sequenceDiagram
	//     Alice->>Bob: How are you?
}
