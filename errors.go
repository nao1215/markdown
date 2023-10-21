package markdown

import "errors"

var (
	// ErrMismatchColumn is returned when the number of columns in the record doesn't match the header.
	ErrMismatchColumn = errors.New("number of columns in the record doesn't match the header")
)
