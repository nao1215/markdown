[![Go Reference](https://pkg.go.dev/badge/github.com/go-spectest/markdown.svg)](https://pkg.go.dev/github.com/go-spectest/markdown)
[![LinuxUnitTest](https://github.com/go-spectest/markdown/actions/workflows/linux_test.yml/badge.svg)](https://github.com/go-spectest/markdown/actions/workflows/linux_test.yml)
[![MacUnitTest](https://github.com/go-spectest/markdown/actions/workflows/mac_test.yml/badge.svg)](https://github.com/go-spectest/markdown/actions/workflows/mac_test.yml)
[![WindowsUnitTest](https://github.com/go-spectest/markdown/actions/workflows/windows_test.yml/badge.svg)](https://github.com/go-spectest/markdown/actions/workflows/windows_test.yml)
[![reviewdog](https://github.com/go-spectest/markdown/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/go-spectest/markdown/actions/workflows/reviewdog.yml)
[![Gosec](https://github.com/go-spectest/markdown/actions/workflows/gosec.yml/badge.svg)](https://github.com/go-spectest/markdown/actions/workflows/gosec.yml)
![Coverage](https://raw.githubusercontent.com/go-spectest/octocovs-central-repo/main/badges/go-spectest/markdown/coverage.svg)
# What is markdown package
The Package markdown is a simple markdown builder in golang. This library assembles Markdown using method chaining, not uses a template engine like [html/template](https://pkg.go.dev/html/template). 
  
This library was initially developed to display test results in [go-spectest/spectest](https://github.com/go-spectest/spectest). Therefore, it implements the features required by spectest, but there are no plans to add additional functionalities unless requested by someone.
  
Additionally, complex code that increases the complexity of the library, such as generating nested lists, will not be added. I want to keep this library as simple as possible.
  
## Supported OS and go version
- OS: Linux, macOS, Windows
- Go: 1.18 or later
  
## Supported Markdown features
- [x] Heading; H1, H2, H3, H4, H5, H6
- [x] Blockquote 
- [x] Bullet list
- [x] Ordered list
- [x] Checkbox list 
- [x] Code blocks
- [x] Horizontal rule 
- [x] Table
- [x] Text formatting; bold, italic, code, strikethrough, bold italic
- [x] text with link
- [x] text with image
- [x] plain text
- [x] details 

### Features not in Markdown syntax
- Generate badges; RedBadge(), YellowBadge(), GreenBadge().
- Generate an index for a directory full of markdown files; GenerateIndex()
  
## Example
### Basic usage
```go
package main

import (
	"os"

	md "github.com/go-spectest/markdown"
)

func main() {
	md.NewMarkdown(os.Stdout).
		H1("This is H1").
		PlainText("This is plain text").
		H2f("This is %s with text format", "H2").
		PlainTextf("Text formatting, such as %s and %s, %s styles.",
			md.Bold("bold"), md.Italic("italic"), md.Code("code")).
		H2("Code Block").
		CodeBlocks(md.SyntaxHighlightGo,
			`package main
import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`).
		H2("List").
		BulletList("Bullet Item 1", "Bullet Item 2", "Bullet Item 3").
		OrderedList("Ordered Item 1", "Ordered Item 2", "Ordered Item 3").
		H2("CheckBox").
		CheckBox([]md.CheckBoxSet{
			{Checked: false, Text: md.Code("sample code")},
			{Checked: true, Text: md.Link("Go", "https://golang.org")},
			{Checked: false, Text: md.Strikethrough("strikethrough")},
		}).
		H2("Blockquote").
		Blockquote("If you can dream it, you can do it.").
		H3("Horizontal Rule").
		HorizontalRule().
		H2("Table").
		Table(md.TableSet{
			Header: []string{"Name", "Age", "Country"},
			Rows: [][]string{
				{"David", "23", "USA"},
				{"John", "30", "UK"},
				{"Bob", "25", "Canada"},
			},
		}).
		H2("Image").
		PlainTextf(md.Image("sample_image", "./sample.png")).
		Build()
}
```

Output:
````
# This is H1
This is plain text
  
## This is H2 with text format
Text formatting, such as **bold** and *italic*, `code` styles.
  
## Code Block
```go
package main
import "fmt"

func main() {
        fmt.Println("Hello, World!")
}
```
  
## List
- Bullet Item 1
- Bullet Item 2
- Bullet Item 3
1. Ordered Item 1
2. Ordered Item 2
3. Ordered Item 3
  
## CheckBox
- [ ] `sample code`
- [x] [Go](https://golang.org)
- [ ] ~~strikethrough~~
  
## Blockquote
> If you can dream it, you can do it.
  
### Horizontal Rule
---
  
## Table
| NAME  | AGE | COUNTRY |
|-------|-----|---------|
| David |  23 | USA     |
| John  |  30 | UK      |
| Bob   |  25 | Canada  |

## Image
![sample_image](./sample.png)
````

If you want to see how it looks in Markdown, please refer to the following link.
- [sample.md](./doc/generated_example.md)

### Generate status badge
The markdown package can create red, yellow, and green status badges.
[Code example:](./doc/badge/main.go)
```go
	md.NewMarkdown(os.Stdout).
		H1("badge example").
		RedBadge("red_badge").
		YellowBadge("yellow_badge").
		GreenBadge("green_badge").
		Build()
```

[Output:](./doc/badge/generated.md)
````text
# badge example
![Badge](https://img.shields.io/badge/red_badge-red)
![Badge](https://img.shields.io/badge/yellow_badge-yellow)
![Badge](https://img.shields.io/badge/green_badge-green)
````

Your badge will look like this;  
![Badge](https://img.shields.io/badge/red_badge-red)
![Badge](https://img.shields.io/badge/yellow_badge-yellow)
![Badge](https://img.shields.io/badge/green_badge-green)

### Generate Markdown using `"go generate ./..."`
You can generate Markdown using `go generate`. Please define code to generate Markdown first. Then, run `"go generate ./..."` to generate Markdown.

[Code example:](./doc/generate/main.go)
```go
package main

import (
	"os"

	md "github.com/go-spectest/markdown"
)

//go:generate go run main.go

func main() {
	f, err := os.Create("generated.md")
	if err != nil {
		panic(err)
	}

	md.NewMarkdown(f).
		H1("go generate example").
		PlainText("This markdown is generated by `go generate`").
		Build()
}
```

Run below command:
```shell
go generate ./...
```

[Output:](./doc/generate/generated.md)
````text
# go generate example
This markdown is generated by `go generate`
````

## Creating an index for a directory full of markdown files
The markdown package can create an index for Markdown files within the specified directory. This feature was added to generate indexes for Markdown documents produced by [go-spectest/spectest](https://github.com/go-spectest/spectest).
  
For example, consider the following directory structure:

```shell
testdata
├── abc
│   ├── dummy.txt
│   ├── jkl
│   │   └── text.md
│   └── test.md
├── def
│   ├── test.md
│   └── test2.md
├── expected
│   └── index.md
├── ghi
└── test.md
```
  
In the following implementation, it creates an index markdown file containing links to all markdown files located within the testdata directory.

```go
		if err := GenerateIndex(
			"testdata", // target directory that contains markdown files
			WithTitle("Test Title"), // title of index markdown
			WithDescription([]string{"Test Description", "Next Description"}), // description of index markdown
		); err != nil {
			panic(err)
		}
```
  
The index Markdown file is created under "target directory/index.md" by default. If you want to change this path, please use the `WithWriter()` option. The link names in the file will be the first occurrence of H1 or H2 in the target Markdown. If neither H1 nor H2 is present, the link name will be the file name of the destination.  
  
[Output:](./doc/index.md)
```markdown
## Test Title
Test Description
  
Next Description
  
### testdata
- [test.md](test.md)
  
### abc
- [h2 is here](abc/test.md)
  
### jkl
- [text.md](abc/jkl/text.md)
  
### def
- [h2 is first, not h1](def/test.md)
- [h1 is here](def/test2.md)
  
### expected
- [Test Title](expected/index.md)
```
  
## Contribution
First off, thanks for taking the time to contribute! Contributions are not only related to development. For example, GitHub Star motivates me to develop! Please feel free to contribute to this project.

## License
[MIT License](./LICENSE)
