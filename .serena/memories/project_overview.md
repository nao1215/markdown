# Project Overview

- **Name**: github.com/nao1215/markdown
- **Purpose**: Go library for building Markdown programmatically using method chaining (no templates). GitHub-flavored Markdown.
- **Language**: Go 1.23+
- **License**: MIT

## Key Packages
- Root: Markdown builder (`markdown.go`, `syntax_sugar.go`, `alert.go`, `badge.go`, `index.go`)
- `mermaid/sequence`, `mermaid/er`, `mermaid/flowchart`, `mermaid/piechart`, `mermaid/quadrant`, `mermaid/state`, `mermaid/class`, `mermaid/gantt`, `mermaid/arch`

## Mermaid Package Pattern
Each mermaid sub-package follows the same structure:
- `config.go` — config struct, Option type, WithTitle, newConfig
- Main file — Diagram/Chart struct with `body []string`, `config *config`, `dest io.Writer`, `err error`; NewXxx constructor, String(), Error(), Build() methods; domain-specific builder methods returning `*Diagram`/`*Chart` for chaining
- `_test.go` — unit tests with `go-cmp`
- `examples_test.go` — example tests (build-tagged `linux || darwin`)
- `doc/<name>/main.go` + `generated.md` — go:generate example

## Commands
- `make test` — run tests with coverage
- `make lint` — run golangci-lint
- `go generate ./...` — regenerate doc examples
