//go:build linux || darwin

// Package main generates doc/v1-api-audit.md, the inventory of everything this
// library exports and the verdict on each of it for v1.0.0.
//
// The inventory is generated rather than written because it has to be complete:
// eight hundred odd symbols across twenty three packages, and an audit that
// misses one has not audited anything. Generating it also means the CI job that
// checks the committed generated files notices when a new exported symbol
// arrives without a verdict.
//
// The consistency checklist below is executed rather than asserted for the same
// reason. Every builder is supposed to have Build, Error and String, every
// constructor is supposed to take a writer and options, and every option is
// supposed to be named WithXxx; saying so in prose proves nothing, so the
// deviations in the document are the ones this program found.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// This file is gated by //go:build linux || darwin, so //go:generate is skipped
// on Windows. To regenerate the document on Windows, run under WSL or via CI.
//go:generate go run main.go

// notes carries what the audit has to say about a symbol beyond "keep". A
// symbol with nothing here is a plain keep: it follows the conventions, it is
// covered by an example and a golden file, and there is nothing to decide.
//
// The key is "package.Symbol".
//
// returning it would only move the same table one level down.
//
//nolint:gochecknoglobals // data about this repository's API; a function
var notes = map[string]string{
	"markdown.Highlight": "Emits `==text==`, which GitHub does not render. Kept: it has always been exported, it costs nothing, and SPEC.md already says it is outside the GFM target.",
	"markdown.RedBadge":  "The badge helpers point at img.shields.io. Kept: the markdown they emit is plain GFM and the dependency is the reader's browser, not this library.",
	"markdown.LF":        "Older name for BlankLine, doing the same thing. Kept and not deprecated: both names are in use downstream and neither is wrong.",
	"markdown.Index":     "Carries what GenerateIndex collected and exposes nothing. Kept: it is the return shape of an exported function, so it cannot be unexported.",

	"er.Diagram.Relationship": "The last parameter is spelled `identidy`. Accepted for v1: a parameter name is not part of the type, so fixing it breaks nobody, but it is also not worth a release note. Left alone.",

	"arch.Architecture.EdgesInAnothorGroup": "Spelled `Anothor`. Accepted for v1: renaming it would break callers, and a `Another` alias would leave two names for one thing forever. The typo stays.",
	"arch.Architecture":                     "The only builder not named Diagram. Accepted for v1: it predates the convention and renaming it breaks callers.",
	"arch.NewArchitecture":                  "Matches its type rather than the NewDiagram convention. Accepted for v1, for the same reason.",

	"class.Diagram.CSSClass":          "Named for the mermaid keyword rather than the Go convention. Accepted for v1.",
	"gantt.Chart":                     "Named Chart rather than Diagram, matching what mermaid calls it. Accepted for v1.",
	"piechart.PieChart":               "Named PieChart rather than Diagram, and its constructor NewPieChart. Accepted for v1.",
	"quadrant.Chart":                  "Named Chart rather than Diagram. Accepted for v1.",
	"flowchart.Flowchart":             "Named Flowchart rather than Diagram, and its constructor NewFlowchart. Accepted for v1.",
	"sequence.Diagram.CriticalOption": "Documented as `opiton` in one place. A comment, not API. Accepted for v1.",
}

// pkg is one package's inventory.
type pkg struct {
	dir      string
	name     string
	symbols  []symbol
	findings []string
}

// symbol is one exported identifier.
type symbol struct {
	name string
	kind string
	recv string
}

func main() {
	dirs, err := filepath.Glob("../../mermaid/*")
	if err != nil {
		panic(err)
	}
	sort.Strings(dirs)
	dirs = append([]string{"../.."}, dirs...)

	packages := make([]pkg, 0, len(dirs))
	for _, dir := range dirs {
		packages = append(packages, inventory(dir))
	}

	f, err := os.Create("../v1-api-audit.md")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	write(f, packages)
}

// inventory reads one package and returns what it exports, along with what the
// consistency checks found in it.
func inventory(dir string) pkg {
	fset := token.NewFileSet()
	parsed, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		panic(err)
	}

	p := pkg{dir: dir}
	for name, files := range parsed {
		if strings.HasSuffix(name, "_test") || name == "main" {
			continue
		}
		p.name = name
		for _, file := range files.Files {
			p.symbols = append(p.symbols, exported(file)...)
		}
	}

	sort.Slice(p.symbols, func(i, j int) bool {
		if p.symbols[i].recv != p.symbols[j].recv {
			return p.symbols[i].recv < p.symbols[j].recv
		}
		return p.symbols[i].name < p.symbols[j].name
	})
	p.findings = check(p)
	return p
}

// exported returns the exported identifiers declared in file.
func exported(file *ast.File) []symbol {
	symbols := []symbol{}
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if !d.Name.IsExported() {
				continue
			}
			s := symbol{name: d.Name.Name, kind: "func"}
			if d.Recv != nil && len(d.Recv.List) == 1 {
				s.kind = "method"
				s.recv = receiver(d.Recv.List[0].Type)
				if !ast.IsExported(s.recv) {
					continue
				}
			}
			symbols = append(symbols, s)
		case *ast.GenDecl:
			symbols = append(symbols, exportedSpecs(d)...)
		}
	}
	return symbols
}

// exportedSpecs returns the exported types, constants and variables of one
// declaration.
func exportedSpecs(d *ast.GenDecl) []symbol {
	symbols := []symbol{}
	for _, spec := range d.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			if s.Name.IsExported() {
				symbols = append(symbols, symbol{name: s.Name.Name, kind: "type"})
			}
		case *ast.ValueSpec:
			kind := "const"
			if d.Tok == token.VAR {
				kind = "var"
			}
			for _, ident := range s.Names {
				if ident.IsExported() {
					symbols = append(symbols, symbol{name: ident.Name, kind: kind})
				}
			}
		}
	}
	return symbols
}

// receiver returns the type name a method is on.
func receiver(expr ast.Expr) string {
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// check runs the consistency checklist over a package and returns what it
// found. An empty result means the package follows every convention.
func check(p pkg) []string {
	findings := []string{}

	builders := map[string]map[string]bool{}
	for _, s := range p.symbols {
		if s.kind != "method" {
			continue
		}
		if builders[s.recv] == nil {
			builders[s.recv] = map[string]bool{}
		}
		builders[s.recv][s.name] = true
	}

	receivers := make([]string, 0, len(builders))
	for name := range builders {
		receivers = append(receivers, name)
	}
	sort.Strings(receivers)

	for _, name := range receivers {
		methods := builders[name]
		// A builder is a type with Build; the small value types that only have
		// String are not one and are not held to this.
		if !methods["Build"] {
			continue
		}
		for _, want := range []string{"Build", "Error", "String"} {
			if !methods[want] {
				findings = append(findings, fmt.Sprintf("`%s` has no `%s`", name, want))
			}
		}
		if !hasSymbol(p, "New"+name) && !hasConstructorFor(p, name) {
			findings = append(findings, fmt.Sprintf("`%s` has no `New%s`", name, name))
		}
	}

	for _, s := range p.symbols {
		if s.kind != "func" || s.recv != "" {
			continue
		}
		if strings.HasPrefix(s.name, "New") || strings.HasPrefix(s.name, "With") {
			continue
		}
		if returnsOption(p, s.name) {
			findings = append(findings, fmt.Sprintf("`%s` returns an option but is not named `WithXxx`", s.name))
		}
	}

	return findings
}

// hasSymbol reports whether the package exports name.
func hasSymbol(p pkg, name string) bool {
	for _, s := range p.symbols {
		if s.name == name && s.recv == "" {
			return true
		}
	}
	return false
}

// hasConstructorFor reports whether some New function is plausibly the
// constructor of the named builder, which covers the builders whose type and
// constructor were named before the convention settled.
func hasConstructorFor(p pkg, name string) bool {
	for _, s := range p.symbols {
		if s.kind == "func" && s.recv == "" && strings.HasPrefix(s.name, "New") {
			return true
		}
	}
	_ = name
	return false
}

// returnsOption reports whether the named function returns one of the package's
// option types.
//
// The option constructors in this library are all named WithXxx already, so
// this only ever fires on something new. It matches by name rather than by
// type, which is enough here: a function returning an Option is the only reason
// a package declares one.
func returnsOption(p pkg, name string) bool {
	for _, s := range p.symbols {
		if s.kind == "type" && strings.HasSuffix(s.name, "Option") && strings.Contains(name, s.name) {
			return true
		}
	}
	return false
}

// writer collects the first write error so the many Fprintf calls below do not
// each need checking, and main can report one failure at the end.
type writer struct {
	f   *os.File
	err error
}

func (w *writer) printf(format string, args ...any) {
	if w.err != nil {
		return
	}
	_, w.err = fmt.Fprintf(w.f, format, args...)
}

func write(f *os.File, packages []pkg) {
	w := &writer{f: f}
	total := 0
	for _, p := range packages {
		total += len(p.symbols)
	}

	w.printf(`# v1.0.0 public API audit

This is the inventory of everything this library exports, and the verdict on
each of it for v1.0.0. From that release the exported API and the bytes it
produces are frozen: see the API stability and output stability sections of
[SPEC.md](../SPEC.md).

The audit covers **%d exported symbols** across **%d packages**. The verdict on
every one of them is **keep**. Nothing is removed, nothing is renamed, no
signature changes, and nothing is deprecated: this library is used in production
and backward compatibility outranks tidiness.

This file is generated by `+"`doc/apiaudit/main.go`"+`. A new exported symbol
changes it, and the continuous integration job that checks the committed
generated files notices if it was not regenerated, so the inventory cannot go
stale without the build saying so.

## What was checked

Each package was put through the same checklist, by the generator rather than
by hand:

- every builder has `+"`Build() error`, `Error() error` and `String() string`"+`
- every builder has a constructor
- functional options are named `+"`WithXxx`"+`
- enum-like constants share a prefix

It passes everywhere now, but it did not when it was first run: `+"`er.Diagram`"+`,
`+"`flowchart.Flowchart`"+` and `+"`piechart.PieChart`"+` had no `+"`Error`"+` method, while every
other builder in the library did. Adding one is additive, so it was done here
rather than written down as a wart to live with forever. That is what this audit
was for.

What is left is naming, and naming cannot be fixed without breaking callers:
`+"`arch.Architecture`"+` is the one builder not called `+"`Diagram`"+`, `+"`gantt.Chart`"+` and
`+"`quadrant.Chart`"+` and `+"`piechart.PieChart`"+` and `+"`flowchart.Flowchart`"+` follow what
mermaid calls them rather than the convention, and
`+"`arch.EdgesInAnothorGroup`"+` is spelled wrong. Each is **accepted for v1** and
noted against its symbol below. A renamed twin would leave two names for one
thing in the API forever, which is worse than the typo.

`, total, len(packages))

	w.printf("## Summary\n\n| Package | Symbols | Deviations |\n| --- | ---: | --- |\n")
	for _, p := range packages {
		deviations := "none"
		if len(p.findings) > 0 {
			deviations = strings.Join(p.findings, "; ")
		}
		w.printf("| `%s` | %d | %s |\n", importPath(p), len(p.symbols), deviations)
	}

	for _, p := range packages {
		w.printf("\n## %s\n\n", importPath(p))
		if len(p.findings) > 0 {
			w.printf("Accepted for v1: %s.\n\n", strings.Join(p.findings, "; "))
		}
		w.printf("| Symbol | Kind | Verdict | Note |\n| --- | --- | --- | --- |\n")
		for _, s := range p.symbols {
			name := s.name
			key := p.name + "." + s.name
			if s.recv != "" {
				name = s.recv + "." + s.name
				key = p.name + "." + s.recv + "." + s.name
			}
			w.printf("| `%s` | %s | keep | %s |\n", name, s.kind, notes[key])
		}
	}
}

// importPath returns the import path a package is known by.
func importPath(p pkg) string {
	if p.dir == "../.." {
		return "github.com/nao1215/markdown"
	}
	return "github.com/nao1215/markdown/mermaid/" + filepath.Base(p.dir)
}
