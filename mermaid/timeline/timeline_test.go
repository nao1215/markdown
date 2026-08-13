package timeline_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
	"github.com/nao1215/markdown/mermaid/timeline"
)

func TestDiagram(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(w io.Writer) *timeline.Diagram
		want  []string
	}{
		"a bare timeline": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w)
			},
			want: []string{"timeline"},
		},
		"a title": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w, timeline.WithTitle("History of social media"))
			},
			want: []string{"timeline", "    title History of social media"},
		},
		"a title is trimmed": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w, timeline.WithTitle("   spaced   "))
			},
			want: []string{"timeline", "    title spaced"},
		},
		"a title holding a colon stays as it is": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w, timeline.WithTitle("History: the sequel"))
			},
			// A title statement reads to the end of the line, so a colon in one
			// is not the separator it is on a period line.
			want: []string{"timeline", "    title History: the sequel"},
		},
		"a period with one event": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("2002", "LinkedIn")
			},
			want: []string{"timeline", "    2002 : LinkedIn"},
		},
		"a period with several events": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("2004", "Facebook", "Google")
			},
			want: []string{"timeline", "    2004 : Facebook : Google"},
		},
		"a period with no events at all": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("2002")
			},
			want: []string{"timeline", "    2002"},
		},
		"events added afterwards": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).
					Period("2004", "Facebook").
					Event("Google").
					Event("Flickr")
			},
			want: []string{"timeline", "    2004 : Facebook : Google : Flickr"},
		},
		"a section indents the periods under it": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).
					Period("2002", "LinkedIn").
					Section("Second wave").
					Period("2004", "Facebook").
					Period("2005", "YouTube")
			},
			want: []string{
				"timeline",
				"    2002 : LinkedIn",
				"    section Second wave",
				"        2004 : Facebook",
				"        2005 : YouTube",
			},
		},
		"a colon in a period becomes an entity": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("09:00", "Stand up")
			},
			// The parser would read a bare colon as the separator, so the period
			// would silently become an event of an empty period. The mermaid
			// form "#58;" is drawn as a colon; the HTML form "&#58;" written
			// before v1.0.0 was drawn as "&:", which is why it went.
			want: []string{"timeline", "    09#58;00 : Stand up"},
		},
		"a colon in an event becomes an entity": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("2004", "Launch: phase one").Event("Review: phase two")
			},
			want: []string{"timeline", "    2004 : Launch#58; phase one : Review#58; phase two"},
		},
		"a colon in a section name becomes an entity": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Section("Phase: one")
			},
			want: []string{"timeline", "    section Phase#58; one"},
		},
		"a hash in a period becomes an entity": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("deploy #2", "Ship")
			},
			// A bare "#" opens a comment and the drawing showed "deploy" with
			// no word about the rest, measured by rendering.
			want: []string{"timeline", "    deploy #35;2 : Ship"},
		},
		"a percent run in an event becomes entities": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("2004", "coverage 100%% now")
			},
			// A "%%" run opens a comment and loses the whole line; a lone "%"
			// reaches the drawing and is left alone.
			want: []string{"timeline", "    2004 : coverage 100#37;#37; now"},
		},
		"a lone percent is left alone": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("2004", "50% of traffic")
			},
			want: []string{"timeline", "    2004 : 50% of traffic"},
		},
		"a hash in a section is left alone": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Section("phase #2")
			},
			// A section name carries a bare "#" to the drawing, measured by
			// rendering, so it is not touched.
			want: []string{"timeline", "    section phase #2"},
		},
		"a literal entity in a section stays literal": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Section("a#58;b")
			},
			// The "#" would otherwise decode as this package's own escape, and
			// a caller's literal "#58;" has to stay distinct from a colon.
			want: []string{"timeline", "    section a#35;58;b"},
		},
		"an angle in the title becomes an entity": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w, timeline.WithTitle("cost < 10"))
			},
			// The sanitizer draws a bare "<" as "&lt;"; the entity form is
			// drawn as the character.
			want: []string{"timeline", "    title cost #60; 10"},
		},
		"text is trimmed": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("  2002  ", "  LinkedIn  ")
			},
			want: []string{"timeline", "    2002 : LinkedIn"},
		},
		"a line feed": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("2002", "LinkedIn").LF().Period("2004", "Facebook")
			},
			want: []string{"timeline", "    2002 : LinkedIn", "", "    2004 : Facebook"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			buf := &bytes.Buffer{}
			if err := tt.build(buf).Build(); err != nil {
				t.Fatalf("Build() = %v, want nil", err)
			}

			want := strings.Join(tt.want, "\n")
			got := strings.ReplaceAll(buf.String(), "\r\n", "\n")
			if got != want {
				t.Errorf("diagram =\n%s\nwant\n%s", got, want)
			}
		})
	}
}

func TestDiagramRejectsWhatItCannotWrite(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(w io.Writer) *timeline.Diagram
		want  string
	}{
		"a title holding a newline": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w, timeline.WithTitle("first\nsecond"))
			},
			want: "title must not contain newline characters",
		},
		"a period holding a newline": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("first\nsecond", "event")
			},
			want: "period must not contain newline characters",
		},
		"an event holding a newline": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("2002", "first\nsecond")
			},
			want: "must not contain newline characters",
		},
		"an appended event holding a carriage return": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("2002", "one").Event("first\rsecond")
			},
			want: "event must not contain newline characters",
		},
		"a section holding a newline": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Section("first\nsecond")
			},
			want: "section name must not contain newline characters",
		},
		"an empty period": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("   ")
			},
			want: "period must not be empty",
		},
		"an empty event": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("2002", "  ")
			},
			want: "must not be empty",
		},
		"an empty section name": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Section("")
			},
			want: "section name must not be empty",
		},
		"an event with no period": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Event("orphan")
			},
			want: `event "orphan" requires a period`,
		},
		"an event after a section that has no period yet": {
			build: func(w io.Writer) *timeline.Diagram {
				return timeline.NewDiagram(w).Period("2002", "one").Section("next").Event("orphan")
			},
			want: `event "orphan" requires a period`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			d := tt.build(io.Discard)

			err := d.Build()
			if err == nil {
				t.Fatal("Build() = nil, want an error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Build() = %v, want it to mention %q", err, tt.want)
			}
			if d.Error() == nil {
				t.Error("Error() = nil, want the same error Build returned")
			}
		})
	}
}

// TestTheFirstErrorIsKept pins that a chain records the error that explains the
// rest, and that the calls after it change nothing.
func TestTheFirstErrorIsKept(t *testing.T) {
	t.Parallel()

	d := timeline.NewDiagram(io.Discard).
		Period("").
		Section("").
		Period("2002", "LinkedIn").
		LF().
		Event("Later")

	err := d.Error()
	if err == nil {
		t.Fatal("Error() = nil, want an error")
	}
	if want := "period must not be empty"; !strings.Contains(err.Error(), want) {
		t.Errorf("Error() = %v, want it to mention %q", err, want)
	}
}

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return timeline.NewDiagram(w).Period("2002", "LinkedIn")
	})
}

// TestRecordedErrorContract asserts that an event written before any period
// surfaces from Build.
func TestRecordedErrorContract(t *testing.T) {
	t.Parallel()

	buildertest.RunRecordedErrorContract(t, func(w io.Writer) buildertest.Builder {
		return timeline.NewDiagram(w).Event("an event with no period")
	})
}

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
