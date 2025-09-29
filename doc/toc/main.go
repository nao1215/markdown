//go:build linux || darwin

// Package main is generating markdown with table of contents.
package main

import (
	"os"

	md "github.com/nao1215/markdown"
)

//go:generate go run main.go

func main() {
	f, err := os.Create("generated.md")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if err := md.NewMarkdown(f).
		H1("Table of Contents Example").
		PlainText("This document demonstrates the table of contents functionality.").
		LF().
		H2("Table of Contents").
		TableOfContentsWithRange(md.TableOfContentsDepthH2, md.TableOfContentsDepthH4).
		LF().
		H2("Introduction").
		PlainText("This is the introduction section. It provides an overview of the document.").
		LF().
		H3("Purpose").
		PlainText("The purpose of this document is to showcase the table of contents feature.").
		LF().
		H3("Scope").
		PlainText("This document covers basic usage of table of contents with different depth levels.").
		LF().
		H2("Features").
		PlainText("The table of contents feature offers several capabilities:").
		LF().
		H3("Basic Table of Contents").
		PlainText("Generate a simple table of contents with all headers up to a specified depth.").
		LF().
		H4("Example Usage").
		CodeBlocks(md.SyntaxHighlightGo, `md.NewMarkdown(os.Stdout).
    H1("Title").
    TableOfContents(md.TableOfContentsDepthH3).
    H2("Section").
    Build()`).
		LF().
		H3("Range-based Table of Contents").
		PlainText("Generate a table of contents that includes only headers within a specific range.").
		LF().
		H4("Advanced Usage").
		CodeBlocks(md.SyntaxHighlightGo, `md.NewMarkdown(os.Stdout).
    H1("Title").  // This H1 will not appear in TOC
    H2("TOC").
    TableOfContentsWithRange(md.TableOfContentsDepthH2, md.TableOfContentsDepthH4).
    H2("Section").
    H3("Subsection").
    H4("Detail").
    H5("Deep Detail").  // This H5 will not appear in TOC
    Build()`).
		LF().
		H4("Benefits").
		BulletList(
			"Control which header levels are included",
			"Exclude document title from TOC",
			"Flexible placement anywhere in document",
			"Automatic anchor generation",
		).
		LF().
		H2("Implementation Details").
		PlainText("The table of contents is implemented using placeholder markers that are replaced during the build process.").
		LF().
		H3("Markers").
		PlainText("The implementation uses HTML comment markers to identify where the TOC should be placed:").
		LF().
		CodeBlocks(md.SyntaxHighlightHTML, `<!-- BEGIN_TOC -->
- [Section 1](#section-1)
  - [Subsection 1.1](#subsection-11)
<!-- END_TOC -->`).
		LF().
		H3("Anchor Generation").
		PlainText("Anchors are automatically generated from header text using GitHub-style conventions:").
		LF().
		BulletList(
			"Convert to lowercase",
			"Replace spaces with hyphens",
			"Remove special characters",
			"Keep only alphanumeric characters and hyphens",
		).
		LF().
		H2("Best Practices").
		PlainText("Follow these best practices when using table of contents:").
		LF().
		OrderedList(
			"Place TOC after the main title and any introductory content",
			"Use range-based TOC to exclude the document title",
			"Choose appropriate depth levels for your document structure",
			"Ensure header text is descriptive and unique",
		).
		LF().
		H2("Conclusion").
		PlainText("The table of contents feature provides a powerful way to improve document navigation and structure.").
		H5("Not include Table Of Contents").
		Build(); err != nil {
		panic(err)
	}
}
