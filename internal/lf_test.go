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
