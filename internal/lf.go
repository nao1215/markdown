// Package internal package is used to store the internal implementation of the mermaid package.
package internal

import "runtime"

// LineFeed return line feed for current OS.
func LineFeed() string {
	return lineFeed(runtime.GOOS)
}

// lineFeed returns the line ending of the named operating system.
//
// The operating system is a parameter rather than read here so that both
// answers can be tested on either platform. Reading runtime.GOOS directly left
// the branch for the other platform untested wherever the tests happened to
// run, which is the half of the behavior most likely to be wrong.
func lineFeed(goos string) string {
	if goos == "windows" {
		return "\r\n"
	}
	return "\n"
}
