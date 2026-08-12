//go:build linux || darwin

package state_test

import (
	"fmt"
	"io"
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/state"
)

// ExampleDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleDiagram() {
	diagram := state.NewDiagram(
		os.Stdout,
		state.WithTitle("Simple State Diagram"),
	).
		StartTransition("Still").
		Transition("Still", "Moving").
		TransitionWithNote("Moving", "Crash", "sudden stop").
		EndTransition("Crash").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("State Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## State Diagram
	// ```mermaid
	// ---
	// title: "Simple State Diagram"
	// ---
	// stateDiagram-v2
	//     [*] --> Still
	//     Still --> Moving
	//     Moving --> Crash : sudden stop
	//     Crash --> [*]
	// ```
}

// ExampleNewDiagram shows the shape every state diagram has: a writer, a chain of
// calls, and Build.
func ExampleNewDiagram() {
	_ = state.NewDiagram(os.Stdout).
		State("Draft", "Being written").Transition("Draft", "Review").
		Build()

	// Output:
	// stateDiagram-v2
	//     Draft : Being written
	//     Draft --> Review
}

// ExampleDiagram_String returns the diagram without needing a writer, which is
// how it is handed to a markdown code block.
func ExampleDiagram_String() {
	diagram := state.NewDiagram(io.Discard).
		State("Draft", "Being written").Transition("Draft", "Review").
		String()

	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// stateDiagram-v2
	//     Draft : Being written
	//     Draft --> Review
	// ```
}

// ExampleDiagram_Build writes the diagram and reports the first error the chain
// recorded. Nothing in the chain panics on bad input, so one check at the end
// is enough.
func ExampleDiagram_Build() {
	err := state.NewDiagram(nil).
		State("Draft", "Being written").Transition("Draft", "Review").
		Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleDiagram_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleDiagram_Error() {
	d := state.NewDiagram(io.Discard).
		State("", "")
	fmt.Println("error:", d.Error())

	// Output:
	// error: <nil>
}

// ExampleDiagram_LF adds a blank line to the diagram body.
func ExampleDiagram_LF() {
	_ = state.NewDiagram(os.Stdout).
		State("Draft", "Being written").Transition("Draft", "Review").
		LF().
		State("Draft", "Being written").Transition("Draft", "Review").
		Build()

	// Output:
	// stateDiagram-v2
	//     Draft : Being written
	//     Draft --> Review
	//
	//     Draft : Being written
	//     Draft --> Review
}

// ExampleDiagram_full shows a state diagram built end to end and put into a markdown
// document, which is what this package exists for.
func ExampleDiagram_full() {
	diagram := state.NewDiagram(io.Discard).
		State("Draft", "Being written").Transition("Draft", "Review").
		String()

	_ = md.NewMarkdown(os.Stdout).
		H2("Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Diagram
	// ```mermaid
	// stateDiagram-v2
	//     Draft : Being written
	//     Draft --> Review
	// ```
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewDiagram.
func ExampleOption() {
	options := []state.Option{state.WithTitle("Overview")}

	_ = state.NewDiagram(os.Stdout, options...).
		State("Draft", "Being written").Transition("Draft", "Review").
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// stateDiagram-v2
	//     Draft : Being written
	//     Draft --> Review
}

// ExampleWithTitle sets the title the diagram is drawn with.
func ExampleWithTitle() {
	_ = state.NewDiagram(os.Stdout, state.WithTitle("Overview")).
		State("Draft", "Being written").Transition("Draft", "Review").
		Build()

	// Output:
	// ---
	// title: "Overview"
	// ---
	// stateDiagram-v2
	//     Draft : Being written
	//     Draft --> Review
}

// ExampleDiagram_State declares a state and the description drawn inside it.
func ExampleDiagram_State() {
	_ = state.NewDiagram(os.Stdout).
		State("Draft", "Being written").
		State("Review", "Waiting on a reviewer").
		Build()

	// Output:
	// stateDiagram-v2
	//     Draft : Being written
	//     Review : Waiting on a reviewer
}

// ExampleDiagram_Transition draws an arrow between two states.
func ExampleDiagram_Transition() {
	_ = state.NewDiagram(os.Stdout).
		Transition("Draft", "Review").
		Transition("Review", "Published").
		Build()

	// Output:
	// stateDiagram-v2
	//     Draft --> Review
	//     Review --> Published
}

// ExampleDiagram_TransitionWithNote labels the arrow with what causes it.
func ExampleDiagram_TransitionWithNote() {
	_ = state.NewDiagram(os.Stdout).
		TransitionWithNote("Draft", "Review", "submit").
		TransitionWithNote("Review", "Draft", "changes requested").
		Build()

	// Output:
	// stateDiagram-v2
	//     Draft --> Review : submit
	//     Review --> Draft : changes requested
}

// ExampleDiagram_StartTransition draws the arrow from the entry point, which is
// how a diagram says which state it begins in.
func ExampleDiagram_StartTransition() {
	_ = state.NewDiagram(os.Stdout).
		StartTransition("Draft").
		Transition("Draft", "Review").
		Build()

	// Output:
	// stateDiagram-v2
	//     [*] --> Draft
	//     Draft --> Review
}

// ExampleDiagram_StartTransitionWithNote labels the arrow from the entry point.
func ExampleDiagram_StartTransitionWithNote() {
	_ = state.NewDiagram(os.Stdout).
		StartTransitionWithNote("Draft", "create").
		Build()

	// Output:
	// stateDiagram-v2
	//     [*] --> Draft : create
}

// ExampleDiagram_EndTransition draws the arrow to the exit point.
func ExampleDiagram_EndTransition() {
	_ = state.NewDiagram(os.Stdout).
		Transition("Draft", "Published").
		EndTransition("Published").
		Build()

	// Output:
	// stateDiagram-v2
	//     Draft --> Published
	//     Published --> [*]
}

// ExampleDiagram_EndTransitionWithNote labels the arrow to the exit point.
func ExampleDiagram_EndTransitionWithNote() {
	_ = state.NewDiagram(os.Stdout).
		EndTransitionWithNote("Published", "archive").
		Build()

	// Output:
	// stateDiagram-v2
	//     Published --> [*] : archive
}

// ExampleDiagram_NoteRight puts a note beside a state.
func ExampleDiagram_NoteRight() {
	_ = state.NewDiagram(os.Stdout).
		State("Draft", "Being written").
		NoteRight("Draft", "Only the author can see it").
		Build()

	// Output:
	// stateDiagram-v2
	//     Draft : Being written
	//     note right of Draft : Only the author can see it
}

// ExampleDiagram_NoteLeft puts a note on the other side.
func ExampleDiagram_NoteLeft() {
	_ = state.NewDiagram(os.Stdout).
		State("Draft", "Being written").
		NoteLeft("Draft", "Only the author can see it").
		Build()

	// Output:
	// stateDiagram-v2
	//     Draft : Being written
	//     note left of Draft : Only the author can see it
}

// ExampleDiagram_NoteRightMultiLine puts a note of several lines beside a
// state. Unlike the one line form, mermaid reads each line as text, so a colon
// in one needs nothing done to it.
func ExampleDiagram_NoteRightMultiLine() {
	_ = state.NewDiagram(os.Stdout).
		State("Draft", "Being written").
		NoteRightMultiLine("Draft", "Only the author can see it.", "Visibility: private").
		Build()

	// Output:
	// stateDiagram-v2
	//     Draft : Being written
	//     note right of Draft
	//         Only the author can see it.
	//         Visibility: private
	//     end note
}

// ExampleDiagram_NoteLeftMultiLine puts a note of several lines on the other
// side.
func ExampleDiagram_NoteLeftMultiLine() {
	_ = state.NewDiagram(os.Stdout).
		State("Draft", "Being written").
		NoteLeftMultiLine("Draft", "Only the author can see it.", "Visibility: private").
		Build()

	// Output:
	// stateDiagram-v2
	//     Draft : Being written
	//     note left of Draft
	//         Only the author can see it.
	//         Visibility: private
	//     end note
}

// ExampleDiagram_CompositeState opens a state holding states of its own. The
// calls on the builder it returns belong to the inner diagram, and End closes
// it and hands the outer one back.
func ExampleDiagram_CompositeState() {
	_ = state.NewDiagram(os.Stdout).
		CompositeState("Review").
		State("Reading", "A reviewer is reading it").
		State("Commenting", "A reviewer is writing").
		Transition("Reading", "Commenting").
		End().
		Transition("Draft", "Review").
		Build()

	// Output:
	// stateDiagram-v2
	//     state Review {
	//         Reading : A reviewer is reading it
	//         Commenting : A reviewer is writing
	//         Reading --> Commenting
	//     }
	//     Draft --> Review
}

// ExampleCompositeStateBuilder shows what the composite state builder is: the
// inner diagram, with the same calls as the outer one, until End.
func ExampleCompositeStateBuilder() {
	_ = state.NewDiagram(os.Stdout).
		CompositeState("Review").
		StartTransition("Reading").
		State("Reading", "A reviewer is reading it").
		Transition("Reading", "Commenting").
		TransitionWithNote("Commenting", "Reading", "keep reading").
		EndTransition("Commenting").
		End().
		Build()

	// Output:
	// stateDiagram-v2
	//     state Review {
	//         [*] --> Reading
	//         Reading : A reviewer is reading it
	//         Reading --> Commenting
	//         Commenting --> Reading : keep reading
	//         Commenting --> [*]
	//     }
}

// ExampleCompositeStateBuilder_End closes the composite state and returns the
// outer diagram, so the chain carries on where it left off.
func ExampleCompositeStateBuilder_End() {
	_ = state.NewDiagram(os.Stdout).
		CompositeState("Review").
		State("Reading", "A reviewer is reading it").
		End().
		Transition("Draft", "Review").
		Build()

	// Output:
	// stateDiagram-v2
	//     state Review {
	//         Reading : A reviewer is reading it
	//     }
	//     Draft --> Review
}

// ExampleCompositeStateBuilder_State declares a state inside the composite one.
func ExampleCompositeStateBuilder_State() {
	_ = state.NewDiagram(os.Stdout).
		CompositeState("Review").
		State("Reading", "A reviewer is reading it").
		End().
		Build()

	// Output:
	// stateDiagram-v2
	//     state Review {
	//         Reading : A reviewer is reading it
	//     }
}

// ExampleCompositeStateBuilder_Transition draws an arrow inside the composite
// state.
func ExampleCompositeStateBuilder_Transition() {
	_ = state.NewDiagram(os.Stdout).
		CompositeState("Review").
		Transition("Reading", "Commenting").
		End().
		Build()

	// Output:
	// stateDiagram-v2
	//     state Review {
	//         Reading --> Commenting
	//     }
}

// ExampleCompositeStateBuilder_TransitionWithNote labels an arrow inside the
// composite state.
func ExampleCompositeStateBuilder_TransitionWithNote() {
	_ = state.NewDiagram(os.Stdout).
		CompositeState("Review").
		TransitionWithNote("Reading", "Commenting", "found something").
		End().
		Build()

	// Output:
	// stateDiagram-v2
	//     state Review {
	//         Reading --> Commenting : found something
	//     }
}

// ExampleCompositeStateBuilder_StartTransition says which state the composite
// one begins in.
func ExampleCompositeStateBuilder_StartTransition() {
	_ = state.NewDiagram(os.Stdout).
		CompositeState("Review").
		StartTransition("Reading").
		End().
		Build()

	// Output:
	// stateDiagram-v2
	//     state Review {
	//         [*] --> Reading
	//     }
}

// ExampleCompositeStateBuilder_EndTransition says which state the composite one
// leaves from.
func ExampleCompositeStateBuilder_EndTransition() {
	_ = state.NewDiagram(os.Stdout).
		CompositeState("Review").
		EndTransition("Commenting").
		End().
		Build()

	// Output:
	// stateDiagram-v2
	//     state Review {
	//         Commenting --> [*]
	//     }
}

// ExampleDiagram_Fork splits one path into several that run at the same time.
func ExampleDiagram_Fork() {
	_ = state.NewDiagram(os.Stdout).
		Fork("split").
		Transition("Submitted", "split").
		Transition("split", "Linting").
		Transition("split", "Testing").
		Build()

	// Output:
	// stateDiagram-v2
	//     state split <<fork>>
	//     Submitted --> split
	//     split --> Linting
	//     split --> Testing
}

// ExampleDiagram_Join brings the paths a fork split back together.
func ExampleDiagram_Join() {
	_ = state.NewDiagram(os.Stdout).
		Join("rejoin").
		Transition("Linting", "rejoin").
		Transition("Testing", "rejoin").
		Transition("rejoin", "Merged").
		Build()

	// Output:
	// stateDiagram-v2
	//     state rejoin <<join>>
	//     Linting --> rejoin
	//     Testing --> rejoin
	//     rejoin --> Merged
}

// ExampleDiagram_Choice draws the point a path picks one way or the other.
func ExampleDiagram_Choice() {
	_ = state.NewDiagram(os.Stdout).
		Choice("passed").
		Transition("Testing", "passed").
		TransitionWithNote("passed", "Merged", "yes").
		TransitionWithNote("passed", "Failed", "no").
		Build()

	// Output:
	// stateDiagram-v2
	//     state passed <<choice>>
	//     Testing --> passed
	//     passed --> Merged : yes
	//     passed --> Failed : no
}

// ExampleDiagram_Concurrent separates the regions of a state that run at the
// same time.
func ExampleDiagram_Concurrent() {
	_ = state.NewDiagram(os.Stdout).
		CompositeState("Running").
		State("Linting", "go vet").
		End().
		Concurrent().
		Build()

	// Output:
	// stateDiagram-v2
	//     state Running {
	//         Linting : go vet
	//     }
	//     ---
}

// ExampleDiagram_SetDirection says which way the diagram is laid out.
func ExampleDiagram_SetDirection() {
	_ = state.NewDiagram(os.Stdout).
		SetDirection(state.DirectionLR).
		Transition("Draft", "Review").
		Build()

	// Output:
	// stateDiagram-v2
	//     direction LR
	//     Draft --> Review
}

// ExampleDirection shows the four ways a diagram can be laid out.
func ExampleDirection() {
	for _, direction := range []state.Direction{
		state.DirectionTB, state.DirectionBT, state.DirectionLR, state.DirectionRL,
	} {
		_ = state.NewDiagram(os.Stdout).
			SetDirection(direction).
			Transition("Draft", "Review").
			Build()
		fmt.Println()
	}

	// Output:
	// stateDiagram-v2
	//     direction TB
	//     Draft --> Review
	// stateDiagram-v2
	//     direction BT
	//     Draft --> Review
	// stateDiagram-v2
	//     direction LR
	//     Draft --> Review
	// stateDiagram-v2
	//     direction RL
	//     Draft --> Review
}
