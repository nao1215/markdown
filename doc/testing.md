# Testing

How markdown is tested, the command that runs each kind of test locally, and the CI job that runs it.

| Technique | Files | Local command | CI |
|---|---|---|---|
| Unit tests and godoc examples | `*_test.go`, `examples_test.go`, `mermaid/*/*_test.go` | `make test` | `unit_test.yml` (Linux, macOS, Windows; Go 1.26 and stable) |
| Golden files | `testdata/golden/*.md`, `internal/golden` | `go test -run Golden .` (`-update` rewrites them) | `unit_test.yml` |
| Parsing the output back | `contract_test.go` | `go test -run 'Parse|Contract' .` | `unit_test.yml` |
| Property-based tests | `contract_test.go` (`TestProperty*`) | `go test -run TestProperty .` | `unit_test.yml` |
| Metamorphic tests | `contract_test.go` (`TestMetamorphic*`) | `go test -run TestMetamorphic .` | `unit_test.yml` |
| Fuzzing | `Fuzz*` in `contract_test.go`, `markdown_test.go`, `testdata/fuzz/` | `go test -run '^$' -fuzz '^FuzzEscapeTableCell$' -fuzztime 1m .` | `fuzz.yml` (seed corpus on every run, real fuzzing on the schedule) |
| Mutation testing | `.github/workflows/mutation.yml` | `gremlins unleash --workers 2 --exclude-files '^doc/' .` | `mutation.yml` (lines a pull request changes; every package weekly) |
| Mermaid render check | `scripts/mermaid-check` | `make render-check` | `doc_render.yml` |

## Parsing the output back

The contract tests read every generated document back with [goldmark](https://github.com/yuin/goldmark) and the GFM extensions, and assert on the structure: how many tables, rows and cells, which heading levels, which list items. goldmark is a test-only dependency. Comparing structure rather than rendered HTML keeps the tests about this library instead of goldmark's version.

## Property-based tests

The properties use [rapid](https://pkg.go.dev/pgregory.net/rapid) and goldmark as the oracle. Each draws 100 inputs per run with a new seed. A failure is shrunk to a small input and prints the seed; `go test -run <Name> . -rapid.seed=<seed>` replays it, and `-rapid.checks=10000` searches longer.

- A table built with `EscapeCells` parses with the rows and columns it was built with, whatever the cells hold (pipes, backslashes, line breaks, backticks, emphasis markers, any Unicode).
- Bullet and ordered lists parse with their items, headings at the level they were written, and a code block holds exactly its content.

## Metamorphic tests

`WithBlockSpacing` only adds blank lines between blocks, so any sequence of builder calls must parse into the same sequence of blocks with and without it. The one intended difference, consecutive `PlainText` lines that form one paragraph without the option, is left out. The relation found that `PlainText(...).HorizontalRule()` used to write a setext heading.

## Mutation testing

[gremlins](https://github.com/go-gremlins/gremlins) v0.6.0 changes one expression at a time and checks that a test fails. A pull request mutates only the lines it changes; the weekly run mutates every package.

| Date | Killed | Lived | Not covered | Efficacy (killed / (killed + lived)) |
|---|---|---|---|---|
| 2026-11-19 | 889 | 79 | 109 | 91.84% |

The sample programs under `doc/` are excluded. The not-covered mutants are mostly in the mermaid builders (79) and `internal/` (22); the lived ones cluster in `markdown.go` (14), `mermaid/venn` (7), `mermaid/xychart` and `internal/golden` (6 each). The CI job does not gate on the score; a drop below the recorded value is a reason to look at the new lived mutants.
