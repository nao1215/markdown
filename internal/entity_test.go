package internal

import "testing"

func TestEntityEscape(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		r    rune
		want string
	}{
		"a double quote is written by name": {r: '"', want: "#quot;"},
		"a hash is written by number":       {r: '#', want: "#35;"},
		"a percent is written by number":    {r: '%', want: "#37;"},
		"a semicolon is written by number":  {r: ';', want: "#59;"},
		"a colon is written by number":      {r: ':', want: "#58;"},
		"an emoji keeps its code point":     {r: '🎉', want: "#127881;"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := EntityEscape(tt.r); got != tt.want {
				t.Errorf("EntityEscape(%q) = %q, want %q", tt.r, got, tt.want)
			}
		})
	}
}

// TestStartsEntity pins where the line falls between a "#" that has to be
// escaped and one that is ordinary text. Getting this wrong in either direction
// is a defect: too eager and output that already renders changes, too shy and a
// caller's literal "#quot;" draws a quotation mark.
func TestStartsEntity(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		rest string
		want bool
	}{
		"a name and a semicolon":            {rest: "quot;", want: true},
		"digits and a semicolon":            {rest: "123;", want: true},
		"a name followed by more text":      {rest: "quot;b", want: true},
		"one letter is enough":              {rest: "a;", want: true},
		"nothing at all":                    {rest: "", want: false},
		"a semicolon on its own":            {rest: ";", want: false},
		"a name that never closes":          {rest: "quot", want: false},
		"a space before the semicolon":      {rest: "a b;", want: false},
		"punctuation before the semicolon":  {rest: "a-b;", want: false},
		"another hash before the semicolon": {rest: "a#b;", want: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := StartsEntity(tt.rest); got != tt.want {
				t.Errorf("StartsEntity(%q) = %v, want %v", tt.rest, got, tt.want)
			}
		})
	}
}
