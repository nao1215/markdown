# Table of Contents Example
This document demonstrates the table of contents functionality.
  
## Table of Contents
<!-- BEGIN_TOC -->
- [Table of Contents](#table-of-contents)
- [Introduction](#introduction)
  - [Purpose](#purpose)
  - [Scope](#scope)
- [Features](#features)
  - [Basic Table of Contents](#basic-table-of-contents)
    - [Example Usage](#example-usage)
  - [Range-based Table of Contents](#range-based-table-of-contents)
    - [Advanced Usage](#advanced-usage)
    - [Benefits](#benefits)
- [Implementation Details](#implementation-details)
  - [Markers](#markers)
  - [Anchor Generation](#anchor-generation)
- [Best Practices](#best-practices)
- [Conclusion](#conclusion)
<!-- END_TOC -->

  
## Introduction
This is the introduction section. It provides an overview of the document.
  
### Purpose
The purpose of this document is to showcase the table of contents feature.
  
### Scope
This document covers basic usage of table of contents with different depth levels.
  
## Features
The table of contents feature offers several capabilities:
  
### Basic Table of Contents
Generate a simple table of contents with all headers up to a specified depth.
  
#### Example Usage
```go
md.NewMarkdown(os.Stdout).
    H1("Title").
    TableOfContents(md.TableOfContentsDepthH3).
    H2("Section").
    Build()
```
  
### Range-based Table of Contents
Generate a table of contents that includes only headers within a specific range.
  
#### Advanced Usage
```go
md.NewMarkdown(os.Stdout).
    H1("Title").  // This H1 will not appear in TOC
    H2("TOC").
    TableOfContentsWithRange(md.TableOfContentsDepthH2, md.TableOfContentsDepthH4).
    H2("Section").
    H3("Subsection").
    H4("Detail").
    H5("Deep Detail").  // This H5 will not appear in TOC
    Build()
```
  
#### Benefits
- Control which header levels are included
- Exclude document title from TOC
- Flexible placement anywhere in document
- Automatic anchor generation
  
## Implementation Details
The table of contents is implemented using placeholder markers that are replaced during the build process.
  
### Markers
The implementation uses HTML comment markers to identify where the TOC should be placed:
  
```html
<!-- BEGIN_TOC -->
- [Section 1](#section-1)
  - [Subsection 1.1](#subsection-11)
<!-- END_TOC -->
```
  
### Anchor Generation
Anchors are automatically generated from header text using GitHub-style conventions:
  
- Convert to lowercase
- Replace spaces with hyphens
- Remove special characters
- Keep only alphanumeric characters and hyphens
  
## Best Practices
Follow these best practices when using table of contents:
  
1. Place TOC after the main title and any introductory content
2. Use range-based TOC to exclude the document title
3. Choose appropriate depth levels for your document structure
4. Ensure header text is descriptive and unique
  
## Conclusion
The table of contents feature provides a powerful way to improve document navigation and structure.
##### Not include Table Of Contents