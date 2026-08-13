package gantt

import (
	"strings"

	"github.com/nao1215/markdown/internal"
)

// escapeTaskName returns name ready to be written before the colon that starts
// a task's data.
//
// A gantt task is written "name :data", and mermaid splits the line at the
// first colon. A colon inside the name therefore ends it early and everything
// after becomes task data: the chart still draws, with a task named "Deploy"
// where the caller asked for "Deploy: staging". That is the quiet kind of
// failure, worse than a diagram that is lost, because nothing reports it.
//
// mermaid reads its own entity form in a task name, so a colon is written
// "#58;" and reaches the drawing as a colon. A "#" that would start an entity
// is escaped with it, or a name holding "#58;" and a name holding a colon would
// draw the same task. A "#" anywhere else is ordinary text and comes out
// unchanged, which is what keeps the golden files as they are.
//
// A raw line break ends the statement early too, and the rest of the name is
// read as a line of its own; it is written as "<br/>", the line break this
// chart draws in a task name.
//
// The section name and the title need none of this: mermaid reads each as the
// rest of its line, and a colon in one already reaches the drawing intact.
func escapeTaskName(name string) string {
	name = internal.LineBreaksToBr(name)
	if !strings.ContainsAny(name, ":#") {
		return name
	}

	var b strings.Builder
	b.Grow(len(name))
	for i := 0; i < len(name); i++ {
		switch {
		case name[i] == ':':
			b.WriteString(internal.EntityEscape(':'))
		case name[i] == '#' && internal.StartsEntity(name[i+1:]):
			b.WriteString(internal.EntityEscape('#'))
		default:
			b.WriteByte(name[i])
		}
	}
	return b.String()
}
