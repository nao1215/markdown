// Package internal package is used to store the internal implementation of the mermaid package.
package internal

import "runtime"

// LineFeed return line feed for current OS.
func LineFeed() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
