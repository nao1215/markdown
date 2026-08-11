# Style and Conventions

- Method chaining pattern: all builder methods return `*Diagram`/`*Chart`
- Unexported `config` struct with functional `Option` pattern
- `noTitle` constant for empty/default title sentinel
- Error accumulation in `err` field, surfaced via `Error()` method
- `Build()` writes to `dest io.Writer`; `String()` returns the body joined with platform-specific line feeds
- Tests use `github.com/google/go-cmp/cmp` for diffs
- Example tests are build-tagged `//go:build linux || darwin` due to \r\n differences
- Internal package `internal/lf.go` provides `LineFeed()` returning \r\n on Windows, \n elsewhere
- 4-space indent for mermaid body lines, 8-space for nested block members
- Comments use standard Go doc style
- All exported symbols have doc comments
