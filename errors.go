package markdown

import "errors"

var (
	// ErrMismatchColumn is returned when the number of columns in the record doesn't match the header.
	ErrMismatchColumn = errors.New("number of columns in the record doesn't match the header")
	// ErrInitMarkdownIndex is returned when the index can't be initialized.
	ErrInitMarkdownIndex = errors.New("markdown index can't be initialized")
	// ErrCreateMarkdownIndex is returned when the index can't be created.
	ErrCreateMarkdownIndex = errors.New("markdown index can't be created")
	// ErrWriteMarkdownIndex is returned when the index can't be written.
	ErrWriteMarkdownIndex = errors.New("markdown index can't be written")
)
