[![LinuxUnitTest](https://github.com/go-spectest/markdown/actions/workflows/linux_test.yml/badge.svg)](https://github.com/go-spectest/markdown/actions/workflows/linux_test.yml)
[![MacUnitTest](https://github.com/go-spectest/markdown/actions/workflows/mac_test.yml/badge.svg)](https://github.com/go-spectest/markdown/actions/workflows/mac_test.yml)
[![WindowsUnitTest](https://github.com/go-spectest/markdown/actions/workflows/windows_test.yml/badge.svg)](https://github.com/go-spectest/markdown/actions/workflows/windows_test.yml)
[![reviewdog](https://github.com/go-spectest/markdown/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/go-spectest/markdown/actions/workflows/reviewdog.yml)
[![Gosec](https://github.com/go-spectest/markdown/actions/workflows/gosec.yml/badge.svg)](https://github.com/go-spectest/markdown/actions/workflows/gosec.yml)
# What is markdown package
The Package markdown is a simple markdown builder from Go code. This library assembles Markdown using method chaining, not uses a template engine like [html/template](https://pkg.go.dev/html/template). 
  
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
  
## Example
```go
package main

import (
	"os"

	"github.com/go-spectest/markdown"
)

func main() {
	markdown.NewMarkdown(os.Stdout).
		H1("This is H1").
		PlainText("This is plain text").
        LF().
		H2f("This is %s with text format", "H2").
		PlainTextf("Package markdown provides functions for text formatting, such as %s and %s, %s styles.",
			markdown.Bold("bold"), markdown.Italic("italic"), markdown.Code("code")).
        LF().
		H2("Code Block").
		CodeBlocks(markdown.SyntaxHighlightGo,
			`package main
import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`).
        LF().
		H2("List").
		BulletList("Bullet Item 1", "Bullet Item 2", "Bullet Item 3").
		OrderedList("Ordered Item 1", "Ordered Item 2", "Ordered Item 3").
        LF().
		H2("CheckBox").
		CheckBox([]markdown.CheckBoxSet{
			{Checked: false, Text: markdown.Code("sample code")},
			{Checked: true, Text: markdown.Link("Go", "https://golang.org")},
			{Checked: false, Text: markdown.Strikethrough("strikethrough")},
		}).
		H2("Blockquote").
		Blockquote("Your time is limited, don't waste it living someone else's life.").
        LF().
		H3("Horizontal Rule").
		HorizontalRule().
        LF().
		H2("Table").
		Table(markdown.TableSet{
			Header: []string{"Name", "Age", "Country"},
			Rows: [][]string{
				{"David", "23", "USA"},
				{"John", "30", "UK"},
				{"Bob", "25", "Canada"},
			},
		}).
		H2("Image").
		PlainTextf(markdown.Image("sample_image", "./sample.png")).
		Build()
}
```

Output:
````
# This is H1
This is plain text
  
## This is H2 with text format
Package markdown provides functions for text formatting, such as **bold** and *italic*, `code` styles.
  
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
> Your time is limited, don't waste it living someone else's life.
  
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

## License
[MIT License](./LICENSE)