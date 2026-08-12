// Package piechart is mermaid pie chart builder.
package piechart

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
)

func TestPieChart_Build(t *testing.T) {
	t.Parallel()

	t.Run("Build a pie chart with title", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		p := NewPieChart(
			b,
			WithTitle("mermaid pie chart builder"),
			WithShowData(true),
		)
		p.LabelAndIntValue("A", 10)
		p.LabelAndFloatValue("B", 20.1)
		p.LabelAndIntValue("C", 30)

		if err := p.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := fmt.Sprintf(
			"%s\n%s",
			"%%{init: {\"pie\": {\"textPosition\": 0.75}, \"themeVariables\": {\"pieOuterStrokeWidth\": \"5px\"}} }%%",
			"pie showData\n    title mermaid pie chart builder\n    \"A\" : 10\n    \"B\" : 20.100000\n    \"C\" : 30")
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a pie chart with no title", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		p := NewPieChart(
			b,
			WithShowData(true),
			WithTextPosition(0.5),
		)
		p.LabelAndIntValue("A", 10)
		p.LabelAndFloatValue("B", 20.1)
		p.LabelAndIntValue("C", 30)

		if err := p.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := fmt.Sprintf(
			"%s\n%s",
			"%%{init: {\"pie\": {\"textPosition\": 0.50}, \"themeVariables\": {\"pieOuterStrokeWidth\": \"5px\"}} }%%",
			"pie showData\n    \"A\" : 10\n    \"B\" : 20.100000\n    \"C\" : 30",
		)
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a pie chart with bad text position value(less than 0)", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		p := NewPieChart(
			b,
			WithTitle("mermaid pie chart builder"),
			WithShowData(true),
			WithTextPosition(-0.1),
		)
		p.LabelAndIntValue("A", 10)
		p.LabelAndFloatValue("B", 20.1)
		p.LabelAndIntValue("C", 30)

		if err := p.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := fmt.Sprintf(
			"%s\n%s",
			"%%{init: {\"pie\": {\"textPosition\": 0.75}, \"themeVariables\": {\"pieOuterStrokeWidth\": \"5px\"}} }%%",
			"pie showData\n    title mermaid pie chart builder\n    \"A\" : 10\n    \"B\" : 20.100000\n    \"C\" : 30")
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a pie chart with bad text position value(more than 1)", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		p := NewPieChart(
			b,
			WithTitle("mermaid pie chart builder"),
			WithShowData(true),
			WithTextPosition(1.1),
		)
		p.LabelAndIntValue("A", 10)
		p.LabelAndFloatValue("B", 20.1)
		p.LabelAndIntValue("C", 30)

		if err := p.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := fmt.Sprintf(
			"%s\n%s",
			"%%{init: {\"pie\": {\"textPosition\": 0.75}, \"themeVariables\": {\"pieOuterStrokeWidth\": \"5px\"}} }%%",
			"pie showData\n    title mermaid pie chart builder\n    \"A\" : 10\n    \"B\" : 20.100000\n    \"C\" : 30")
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

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
		return NewPieChart(w).LabelAndIntValue("Go", 120)
	})
}

// TestGoldenPieChart pins the rendered chart of every builder method and every
// option of this package.
func TestGoldenPieChart(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewPieChart(
		buf,
		WithTitle("Language Share"),
		WithShowData(true),
		WithTextPosition(0.75),
	).
		LabelAndIntValue("Go", 120).
		LabelAndFloatValue("Rust", 42.5).
		LabelAndFloatValue("Zig", 0.5).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("piechart.md", buf.String()); err != nil {
		t.Error(err)
	}
}

// TestBuildWithNilWriter covers the case where a pie chart is built for String()
// only and Build() is called by mistake. Build() used to dereference the nil
// writer and take the process down; it has to return an error instead.
func TestBuildWithNilWriter(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Build() panicked with a nil writer: %v", r)
		}
	}()

	d := NewPieChart(nil)

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

	err := NewPieChart(errWriter{}).Build()
	if err == nil {
		t.Fatal("Build must report a failing writer")
	}
	if !errors.Is(err, errWrite) {
		t.Errorf("Build lost the destination error: %v", err)
	}
}

// TestLabelEscapesTheQuoteThatEndsIt names the character this escaping buys. A
// double quote in a slice label used to reach mermaid unescaped and lose the
// whole chart: the reader got an error box rather than a picture.
func TestLabelEscapesTheQuoteThatEndsIt(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		build func(*PieChart) *PieChart
		want  string
	}{
		"a quote in an int label": {
			build: func(p *PieChart) *PieChart {
				return p.LabelAndIntValue(`say "hi"`, 10)
			},
			want: `    "say #quot;hi#quot;" : 10`,
		},
		"a quote in a float label": {
			build: func(p *PieChart) *PieChart {
				return p.LabelAndFloatValue(`"`, 1.5)
			},
			want: `    "#quot;" : 1.500000`,
		},
		"a named entity in a label is escaped": {
			build: func(p *PieChart) *PieChart {
				return p.LabelAndIntValue("a#quot;b", 10)
			},
			want: `    "a#35;quot;b" : 10`,
		},
		"a plain hash in a label is left alone": {
			build: func(p *PieChart) *PieChart {
				return p.LabelAndIntValue("PR #123 merged", 10)
			},
			want: `    "PR #123 merged" : 10`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := tt.build(NewPieChart(io.Discard)).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("chart =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}

// TestTitleEscapesThePercentPairThatCommentsItOut names the other character
// this buys. The title is unquoted, so a "%%" in one opened a mermaid comment
// and the rest of the title was dropped from the drawing without any error.
func TestTitleEscapesThePercentPairThatCommentsItOut(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		title string
		want  string
	}{
		"a percent pair is escaped":          {title: "100%% done", want: "    title 100#37;#37; done"},
		"a longer run is escaped whole":      {title: "a%%%b", want: "    title a#37;#37;#37;b"},
		"a lone percent is left alone":       {title: "100% done", want: "    title 100% done"},
		"percents apart are left alone":      {title: "50% of 10%", want: "    title 50% of 10%"},
		"a named entity is escaped":          {title: "a#quot;b", want: "    title a#35;quot;b"},
		"a plain hash is left alone":         {title: "PR #123 merged", want: "    title PR #123 merged"},
		"a quote needs nothing when bare":    {title: `the "core"`, want: `    title the "core"`},
		"text with neither is left as it is": {title: "Sales", want: "    title Sales"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := NewPieChart(io.Discard, WithTitle(tt.title)).String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("chart =\n%s\nwant it to contain\n%s", got, tt.want)
			}
		})
	}
}

// TestErrorReportsWhatBuildWould pins the method the v1.0.0 API audit found
// missing. Every other builder in this library reports its recorded error from
// Error as well as from Build, and this one did not.
func TestErrorReportsWhatBuildWould(t *testing.T) {
	t.Parallel()

	p := NewPieChart(nil).LabelAndIntValue("Go", 60)

	if err := p.Error(); err != nil {
		t.Errorf("Error() = %v before Build, want nil", err)
	}

	fromBuild := p.Build()
	if fromBuild == nil {
		t.Fatal("Build() = nil with a nil writer, want an error")
	}
	if p.Error() == nil || p.Error().Error() != fromBuild.Error() {
		t.Errorf("Error() = %v, want the error Build returned, %v", p.Error(), fromBuild)
	}
}
