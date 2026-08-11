// Package internal package is used to store the internal implementation of the mermaid package.
package internal

import "testing"

func TestFrontMatterTitle(t *testing.T) {
	t.Parallel()

	// Every case here is a title that a bare YAML scalar reads as something
	// other than the title. The colon and the alias make mermaid's front matter
	// parser throw, which loses the whole diagram; the rest are quietly
	// truncated, dropped, or turned into another YAML type.
	tests := []struct {
		name  string
		title string
		want  string
	}{
		{
			name:  "plain title",
			title: "Checkout Architecture",
			want:  `title: "Checkout Architecture"`,
		},
		{
			name:  "colon would end the mapping key",
			title: "Checkout: API",
			want:  `title: "Checkout: API"`,
		},
		{
			name:  "hash would start a comment",
			title: "Checkout # API",
			want:  `title: "Checkout # API"`,
		},
		{
			name:  "leading hash would comment out the value",
			title: "# Checkout",
			want:  `title: "# Checkout"`,
		},
		{
			name:  "asterisk would be an alias",
			title: "*ref",
			want:  `title: "*ref"`,
		},
		{
			name:  "ampersand would be an anchor",
			title: "&anchor",
			want:  `title: "&anchor"`,
		},
		{
			name:  "tilde would be null",
			title: "~",
			want:  `title: "~"`,
		},
		{
			name:  "brackets would be a sequence",
			title: "[Checkout]",
			want:  `title: "[Checkout]"`,
		},
		{
			name:  "true would be a boolean",
			title: "true",
			want:  `title: "true"`,
		},
		{
			name:  "digits would be a number",
			title: "123",
			want:  `title: "123"`,
		},
		{
			name:  "double quote is escaped",
			title: `Checkout "API"`,
			want:  `title: "Checkout \"API\""`,
		},
		{
			name:  "backslash is escaped",
			title: `Checkout\API`,
			want:  `title: "Checkout\\API"`,
		},
		{
			name:  "backslash before a quote stays two separate escapes",
			title: `Checkout\"API`,
			want:  `title: "Checkout\\\"API"`,
		},
		{
			name:  "control characters are escaped",
			title: "Checkout\tAPI\r\nv2",
			want:  `title: "Checkout\tAPI\r\nv2"`,
		},
		{
			name:  "empty title",
			title: "",
			want:  `title: ""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := FrontMatterTitle(tt.title); got != tt.want {
				t.Errorf("FrontMatterTitle(%q) = %q, want %q", tt.title, got, tt.want)
			}
		})
	}
}
