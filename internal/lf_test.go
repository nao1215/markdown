// Package internal package is used to store the internal implementation of the mermaid package.
package internal

import (
	"runtime"
	"testing"
)

func TestLineFeed(t *testing.T) {
	t.Parallel()

	t.Run("should return line feed for current OS", func(t *testing.T) {
		t.Parallel()

		got := LineFeed()

		switch runtime.GOOS {
		case "windows":
			if got != "\r\n" {
				t.Errorf("expected \\r\\n, but got %s", got)
			}
		default:
			if got != "\n" {
				t.Errorf("expected \\n, but got %s", got)
			}
		}
	})
}

// TestLineFeedPerOperatingSystem covers both answers on whichever platform the
// tests run on. The branch that is not the current platform's is the half most
// likely to be wrong, because nothing exercises it by accident.
func TestLineFeedPerOperatingSystem(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"windows": "\r\n",
		"linux":   "\n",
		"darwin":  "\n",
		"plan9":   "\n",
	}

	for goos, want := range tests {
		t.Run(goos, func(t *testing.T) {
			t.Parallel()

			if got := lineFeed(goos); got != want {
				t.Errorf("lineFeed(%q) = %q, want %q", goos, got, want)
			}
		})
	}
}
