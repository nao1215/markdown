package internal

import "testing"

func TestLineBreaksToBr(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   string
		want string
	}{
		"no break is left alone":            {in: "one line", want: "one line"},
		"a newline becomes one break":       {in: "a\nb", want: "a<br/>b"},
		"a carriage return becomes one":     {in: "a\rb", want: "a<br/>b"},
		"a CRLF pair becomes one, not two":  {in: "a\r\nb", want: "a<br/>b"},
		"each break is written separately":  {in: "a\nb\nc", want: "a<br/>b<br/>c"},
		"a break at the edge stays a break": {in: "\nb", want: "<br/>b"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := LineBreaksToBr(tt.in); got != tt.want {
				t.Errorf("LineBreaksToBr(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestEscapeBareAngle(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   string
		want string
	}{
		"no angle is left alone":          {in: "plain", want: "plain"},
		"a bare angle becomes the entity": {in: "a<b then c", want: "a#60;b then c"},
		"a br stays a line break":         {in: "a<br/>b", want: "a<br/>b"},
		"an angle before a br is escaped": {in: "a<<br/>b", want: "a#60;<br/>b"},
		"a closing angle is left alone":   {in: "a>b", want: "a>b"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := EscapeBareAngle(tt.in); got != tt.want {
				t.Errorf("EscapeBareAngle(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestFoldFrontMatterTitleCR(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   string
		want string
	}{
		"a plain title is left alone": {in: "release plan", want: "release plan"},
		// A "<" is left alone deliberately: every escape form draws as its
		// literal text in a front matter title, so there is nothing to write
		// instead. The function's comment records the measurement.
		"an angle is left alone":                 {in: "a<b", want: "a<b"},
		"a CR becomes the LF that draws a break": {in: "a\rb", want: "a\nb"},
		"a CRLF pair becomes one LF":             {in: "a\r\nb", want: "a\nb"},
		"a lone LF is left alone":                {in: "a\nb", want: "a\nb"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := FoldFrontMatterTitleCR(tt.in); got != tt.want {
				t.Errorf("FoldFrontMatterTitleCR(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestEscapeTitle(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   string
		want string
	}{
		"a plain title is left alone":     {in: "release plan", want: "release plan"},
		"a newline becomes the entity":    {in: "a\nb", want: "a#10;b"},
		"a CRLF pair becomes one entity":  {in: "a\r\nb", want: "a#10;b"},
		"an angle and a break together":   {in: "a<b\nc", want: "a#60;b#10;c"},
		"a literal entity stays distinct": {in: "a#10;b\nc", want: "a#35;10;b#10;c"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := EscapeTitle(tt.in); got != tt.want {
				t.Errorf("EscapeTitle(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestEscapeTitleAngle pins where the line falls between the "<" that loses a
// title and the "<br/>" the renderer honors, and keeps the escaping injective
// the same way the label escapes do.
func TestEscapeTitleAngle(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		in   string
		want string
	}{
		"no angle is left alone":            {in: "plain title", want: "plain title"},
		"a bare angle becomes the entity":   {in: "cost < 10", want: "cost #60; 10"},
		"an unfinished tag is escaped":      {in: "a<b", want: "a#60;b"},
		"a br stays a line break":           {in: "a<br/>b", want: "a<br/>b"},
		"an angle before a br is escaped":   {in: "a<<br/>b", want: "a#60;<br/>b"},
		"a closing angle is left alone":     {in: "a>b", want: "a>b"},
		"a literal entity stays literal":    {in: "a#60;b", want: "a#35;60;b"},
		"an ordinary hash is left alone":    {in: "deploy #2", want: "deploy #2"},
		"a br without the slash is escaped": {in: "a<br>b", want: "a#60;br>b"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := EscapeTitleAngle(tt.in); got != tt.want {
				t.Errorf("EscapeTitleAngle(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
