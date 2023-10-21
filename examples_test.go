package markdown_test

import (
	"os"
	"runtime"

	"github.com/go-spectest/markdown"
)

func Example() {
	// Skip this test on Windows.
	// The newline codes in the comment section where
	// the expected values are written are represented as '\n',
	// causing failures when testing on Windows.
	if runtime.GOOS == "windows" {
		return
	}

	markdown.NewMarkdown(os.Stdout).
		H1("This is H1").
		PlainText("This is plain text").
		H2f("This is %s with text format", "H2").
		PlainTextf("Package markdown provides functions for text formatting, such as %s and %s, %s styles.",
			markdown.Bold("bold"), markdown.Italic("italic"), markdown.Code("code")).
		H2("Code Block").
		CodeBlocks(markdown.SyntaxHighlightGo,
			`package main
import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`).
		H2("List").
		BulletList("Bullet Item 1", "Bullet Item 2", "Bullet Item 3").
		OrderedList("Ordered Item 1", "Ordered Item 2", "Ordered Item 3").
		H2("CheckBox").
		CheckBox([]markdown.CheckBoxSet{
			{Checked: false, Text: markdown.Code("sample code")},
			{Checked: true, Text: markdown.Link("Go", "https://golang.org")},
			{Checked: false, Text: markdown.Strikethrough("strikethrough")},
		}).
		H2("Blockquote").
		Blockquote("Your time is limited, don't waste it living someone else's life.").
		H3("Horizontal Rule").
		HorizontalRule().
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

	// Output:
	// # This is H1
	// This is plain text
	// ## This is H2 with text format
	// Package markdown provides functions for text formatting, such as **bold** and *italic*, `code` styles.
	// ## Code Block
	// ```go
	// package main
	// import "fmt"
	//
	// func main() {
	// 	fmt.Println("Hello, World!")
	// }
	// ```
	// ## List
	// - Bullet Item 1
	// - Bullet Item 2
	// - Bullet Item 3
	// 1. Ordered Item 1
	// 2. Ordered Item 2
	// 3. Ordered Item 3
	// ## CheckBox
	// - [ ] `sample code`
	// - [x] [Go](https://golang.org)
	// - [ ] ~~strikethrough~~
	// ## Blockquote
	// > Your time is limited, don't waste it living someone else's life.
	// ### Horizontal Rule
	// ---
	// ## Table
	// | NAME  | AGE | COUNTRY |
	// |-------|-----|---------|
	// | David |  23 | USA     |
	// | John  |  30 | UK      |
	// | Bob   |  25 | Canada  |
	//
	// ## Image
	// ![sample_image](./sample.png)
}
