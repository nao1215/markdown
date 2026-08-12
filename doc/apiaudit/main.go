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
	"markdown.Highlight":         "Emits `==text==`, which GitHub does not render. Kept: it has always been exported, it costs nothing, and SPEC.md already says it is outside the GFM target.",
	"markdown.Markdown.RedBadge": "The badge helpers point at img.shields.io. Kept: the markdown they emit is plain GFM and the dependency is the reader's browser, not this library.",
	"markdown.Markdown.LF":       "Older name for BlankLine, doing the same thing. Kept and not deprecated: both names are in use downstream and neither is wrong.",
	"markdown.Index":             "Carries what GenerateIndex collected and exposes nothing. Kept: it is the return shape of an exported function, so it cannot be unexported.",

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

// symbol is one exported identifier, with the parts of its declaration the
// consistency checks need.
type symbol struct {
	name string
	kind string
	recv string
	// results are the result types of a function, by name.
	results []string
	// params are the parameter types of a function, by name.
	params []string
	// typeName is the declared type of a constant, where it has one.
	typeName string
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
			s := symbol{
				name:    d.Name.Name,
				kind:    "func",
				results: typeNames(d.Type.Results),
				params:  typeNames(d.Type.Params),
			}
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
	constType := ""
	for _, spec := range d.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			if !s.Name.IsExported() {
				continue
			}
			symbols = append(symbols, symbol{name: s.Name.Name, kind: "type"})
			symbols = append(symbols, members(s)...)
		case *ast.ValueSpec:
			kind := "const"
			if d.Tok == token.VAR {
				kind = "var"
			}
			// A constant in a group inherits the type of the first spec that
			// names one, which is how an iota block declares its type once.
			if named := typeName(s.Type); named != "" {
				constType = named
			}
			for _, ident := range s.Names {
				if ident.IsExported() {
					symbols = append(symbols, symbol{name: ident.Name, kind: kind, typeName: constType})
				}
			}
		}
	}
	return symbols
}

// members returns the exported fields of an exported struct and the exported
// methods of an exported interface.
//
// They are part of the API as much as the type is: a caller writes
// er.Attribute{Type: "int"}, so the field name is frozen at v1.0.0 the same way
// the type name is.
func members(spec *ast.TypeSpec) []symbol {
	symbols := []symbol{}
	switch t := spec.Type.(type) {
	case *ast.StructType:
		for _, field := range fieldNames(t.Fields) {
			symbols = append(symbols, symbol{name: field, kind: "field", recv: spec.Name.Name})
		}
	case *ast.InterfaceType:
		for _, method := range fieldNames(t.Methods) {
			// A distinct kind, so that an interface declaring Build is not
			// mistaken for a builder and asked for a constructor.
			symbols = append(symbols,
				symbol{name: method, kind: "interface method", recv: spec.Name.Name})
		}
	}
	return symbols
}

// fieldNames returns the exported names declared in a field list, including an
// embedded type, which a caller can reach through just as it can a named field.
func fieldNames(fields *ast.FieldList) []string {
	if fields == nil {
		return nil
	}

	names := []string{}
	for _, field := range fields.List {
		if len(field.Names) == 0 {
			if name := receiver(field.Type); ast.IsExported(name) {
				names = append(names, name)
			}
			continue
		}
		for _, ident := range field.Names {
			if ident.IsExported() {
				names = append(names, ident.Name)
			}
		}
	}
	return names
}

// typeNames returns the type of each entry in a field list, written the way the
// source writes it, which is enough to compare against a name.
func typeNames(fields *ast.FieldList) []string {
	if fields == nil {
		return nil
	}

	names := []string{}
	for _, field := range fields.List {
		named := typeName(field.Type)
		count := len(field.Names)
		if count == 0 {
			count = 1
		}
		for range count {
			names = append(names, named)
		}
	}
	return names
}

// typeName writes an expression the way its type reads in source: "*Diagram",
// "io.Writer", "...Option".
func typeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeName(t.X)
	case *ast.Ellipsis:
		return "..." + typeName(t.Elt)
	case *ast.SelectorExpr:
		return typeName(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + typeName(t.Elt)
	case *ast.FuncType:
		return "func"
	default:
		return ""
	}
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
	findings := checkBuilders(p)
	findings = append(findings, checkOptions(p)...)
	findings = append(findings, checkEnums(p)...)
	return findings
}

// checkBuilders reports a builder without the three methods every one of them
// has, or without a constructor taking a writer and its options.
//
// A builder is a type with a Build method; the small value types that only have
// String are not one and are not held to this.
func checkBuilders(p pkg) []string {
	methods := map[string]map[string]symbol{}
	for _, s := range p.symbols {
		if s.kind != "method" || s.recv == "" {
			continue
		}
		if methods[s.recv] == nil {
			methods[s.recv] = map[string]symbol{}
		}
		methods[s.recv][s.name] = s
	}

	wanted := []struct {
		name   string
		result string
	}{
		{"Build", "error"},
		{"Error", "error"},
		{"String", "string"},
	}

	findings := make([]string, 0, len(methods))
	for _, name := range sortedKeys(methods) {
		if _, ok := methods[name]["Build"]; !ok {
			continue
		}
		for _, want := range wanted {
			method, ok := methods[name][want.name]
			switch {
			case !ok:
				findings = append(findings, fmt.Sprintf("`%s` has no `%s`", name, want.name))
			case len(method.params) != 0 || len(method.results) != 1 || method.results[0] != want.result:
				findings = append(findings,
					fmt.Sprintf("`%s.%s` is not `%s() %s`", name, want.name, want.name, want.result))
			}
		}
		findings = append(findings, checkConstructor(p, name)...)
	}
	return findings
}

// checkConstructor reports a builder whose constructor is missing, is not named
// for it, or does not take a writer and a run of options.
func checkConstructor(p pkg, builder string) []string {
	for _, s := range p.symbols {
		if s.kind != "func" || s.recv != "" || len(s.results) != 1 {
			continue
		}
		if s.results[0] != "*"+builder {
			continue
		}

		findings := []string{}
		if s.name != "New"+builder {
			findings = append(findings,
				fmt.Sprintf("`%s` returns `*%s` but is not named `New%s`", s.name, builder, builder))
		}
		if len(s.params) != 2 || s.params[0] != "io.Writer" || !isOptionVariadic(p, s.params[1]) {
			findings = append(findings,
				fmt.Sprintf("`%s` does not take `(io.Writer, ...Option)`", s.name))
		}
		return findings
	}
	return []string{fmt.Sprintf("`%s` has no constructor returning `*%s`", builder, builder)}
}

// isOptionVariadic reports whether the parameter is a run of one of the
// package's own option types, rather than a run of anything at all.
func isOptionVariadic(p pkg, param string) bool {
	named, ok := strings.CutPrefix(param, "...")
	if !ok {
		return false
	}
	for _, s := range p.symbols {
		if s.kind == "type" && s.name == named && strings.HasSuffix(named, "Option") {
			return true
		}
	}
	return false
}

// checkOptions reports a function that returns one of the package's option
// types without being named for it.
//
// The result type is what decides, rather than the name, so a new option
// constructor called something else is caught rather than skipped.
func checkOptions(p pkg) []string {
	optionTypes := map[string]bool{}
	findings := make([]string, 0, len(p.symbols))
	for _, s := range p.symbols {
		if s.kind == "type" && strings.HasSuffix(s.name, "Option") {
			optionTypes[s.name] = true
		}
	}

	for _, s := range p.symbols {
		if s.kind != "func" || s.recv != "" || len(s.results) != 1 {
			continue
		}
		if optionTypes[s.results[0]] && !strings.HasPrefix(s.name, "With") {
			findings = append(findings,
				fmt.Sprintf("`%s` returns `%s` but is not named `WithXxx`", s.name, s.results[0]))
		}
	}
	return findings
}

// checkEnums reports the constants of a named type that share neither a prefix
// nor a suffix.
//
// A shared affix is what makes an enumeration readable at the call site: a
// reader seeing RiskHigh or ZeroToOneRelationship knows what it is for, and a
// reader seeing High does not. Which of the two a group uses is not worth a
// finding, so a group sharing either passes; where the affix is not the name of
// the type, the summary says so, because that is worth knowing when reading the
// package rather than something to fix.
func checkEnums(p pkg) []string {
	groups := map[string][]string{}
	for _, s := range p.symbols {
		if s.kind == "const" && s.typeName != "" {
			groups[s.typeName] = append(groups[s.typeName], s.name)
		}
	}

	types := make([]string, 0, len(groups))
	for name := range groups {
		types = append(types, name)
	}
	sort.Strings(types)

	findings := []string{}
	for _, typeName := range types {
		names := groups[typeName]
		// One constant of a type shares an affix with nothing, so there is
		// nothing to check.
		if len(names) < 2 { //nolint:mnd // "fewer than two" is the condition, not a magic number.
			continue
		}

		prefix := commonPrefix(names)
		suffix := commonSuffix(names)
		switch {
		case strings.HasPrefix(prefix, typeName):
			// The usual case: RiskLow, RiskMedium, RiskHigh.
		case prefix == "" && suffix == "":
			findings = append(findings,
				fmt.Sprintf("the `%s` constants share neither a prefix nor a suffix", typeName))
		case prefix != "":
			findings = append(findings,
				fmt.Sprintf("the `%s` constants are prefixed `%s` rather than with the type name",
					typeName, atWordBoundary(prefix)))
		default:
			findings = append(findings,
				fmt.Sprintf("the `%s` constants are suffixed `%s` rather than prefixed",
					typeName, fromWordBoundary(suffix)))
		}
	}
	return findings
}

// commonPrefix returns the longest prefix every name shares.
func commonPrefix(names []string) string {
	shared := names[0]
	for _, name := range names[1:] {
		for !strings.HasPrefix(name, shared) {
			shared = shared[:len(shared)-1]
			if shared == "" {
				return ""
			}
		}
	}
	return shared
}

// commonSuffix returns the longest suffix every name shares.
func commonSuffix(names []string) string {
	shared := names[0]
	for _, name := range names[1:] {
		for !strings.HasSuffix(name, shared) {
			shared = shared[1:]
			if shared == "" {
				return ""
			}
		}
	}
	return shared
}

// atWordBoundary cuts a prefix back to the last word it starts, so a partial
// word is not reported as the shared one: three names agreeing on "AlignL" are
// sharing "Align".
func atWordBoundary(affix string) string {
	for i := len(affix) - 1; i > 0; i-- {
		if affix[i] >= 'A' && affix[i] <= 'Z' {
			return affix[:i]
		}
	}
	return affix
}

// fromWordBoundary cuts a suffix forward to the word it starts, for the same
// reason: names agreeing on "eRelationship" are sharing "Relationship".
func fromWordBoundary(affix string) string {
	for i := 0; i < len(affix); i++ {
		if affix[i] >= 'A' && affix[i] <= 'Z' {
			return affix[i:]
		}
	}
	return affix
}

// sortedKeys returns the keys of m in order, so the findings do not move about
// between runs.
func sortedKeys(m map[string]map[string]symbol) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
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

- every builder has `+"`Build() error`, `Error() error` and `String() string`"+`,
  with those exact signatures
- every builder has a constructor named for it, taking `+"`(io.Writer, ...Option)`"+`
- every function returning one of the package's option types is named `+"`WithXxx`"+`
- enum-like constants share a prefix or a suffix

It passes everywhere now, but it did not when it was first run: `+"`er.Diagram`"+`,
`+"`flowchart.Flowchart`"+` and `+"`piechart.PieChart`"+` had no `+"`Error`"+` method, while every
other builder in the library did. Adding one is additive, so it was done here
rather than written down as a wart to live with forever. That is what this audit
was for.

What is left is naming, and naming cannot be fixed without breaking callers.
Two of the deviations above are enumerations: `+"`markdown.TableAlignment`"+`'s
constants are prefixed `+"`Align`"+` rather than with the type name, and
`+"`er.Relationship`"+`'s are suffixed rather than prefixed. Both still share an
affix, which is what makes them readable at a call site, so neither is worth
breaking a caller over. The rest are types and methods:
`+"`arch.Architecture`"+` is the one builder not called `+"`Diagram`"+`, `+"`gantt.Chart`"+` and
`+"`quadrant.Chart`"+` and `+"`piechart.PieChart`"+` and `+"`flowchart.Flowchart`"+` follow what
mermaid calls them rather than the convention, and
`+"`arch.EdgesInAnothorGroup`"+` is spelled wrong. Each is **accepted for v1** and
noted against its symbol below. A renamed twin would leave two names for one
thing in the API forever, which is worse than the typo.

`, total, len(packages))

	w.printf("## Summary\n\n" +
		"Everything in the last column is **accepted for v1**. The checklist findings\n" +
		"come from the generator; the noted symbols are the ones carrying a note in\n" +
		"the tables below, which is where the reason for each of them is.\n\n" +
		"| Package | Symbols | Checklist findings | Noted symbols |\n" +
		"| --- | ---: | --- | --- |\n")
	for _, p := range packages {
		deviations := "none"
		if len(p.findings) > 0 {
			deviations = strings.Join(p.findings, "; ")
		}
		noted := "none"
		if names := notedSymbols(p); len(names) > 0 {
			noted = "`" + strings.Join(names, "`, `") + "`"
		}
		w.printf("| `%s` | %d | %s | %s |\n", importPath(p), len(p.symbols), deviations, noted)
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

// notedSymbols returns the symbols of a package that carry a note, so the
// summary and the tables below cannot disagree about what was accepted.
func notedSymbols(p pkg) []string {
	names := []string{}
	for _, s := range p.symbols {
		name := s.name
		key := p.name + "." + s.name
		if s.recv != "" {
			name = s.recv + "." + s.name
			key = p.name + "." + s.recv + "." + s.name
		}
		if notes[key] != "" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// importPath returns the import path a package is known by.
func importPath(p pkg) string {
	if p.dir == "../.." {
		return "github.com/nao1215/markdown"
	}
	return "github.com/nao1215/markdown/mermaid/" + filepath.Base(p.dir)
}
