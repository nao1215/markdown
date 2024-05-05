//go:build linux || darwin

package markdown_test

import (
	"os"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/sequence"
)

// ExamleMarkdown skip this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleMarkdown() {
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

	// Output:
	// # This is H1
	// This is plain text
	// ## This is H2 with text format
	// Text formatting, such as **bold** and *italic*, `code` styles.
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
	// > If you can dream it, you can do it.
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

// ExampleSequence skip this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleSequence() {
	diagram := sequence.NewDiagram(os.Stdout).
		Participant("Sophia").
		Participant("David").
		Participant("Subaru").
		LF().
		SyncRequest("Sophia", "David", "Please wake up Subaru").
		SyncResponse("David", "Sophia", "OK").
		LF().
		LoopStart("until Subaru wake up").
		SyncRequest("David", "Subaru", "Wake up!").
		SyncResponse("Subaru", "David", "zzz").
		SyncRequest("David", "Subaru", "Hey!!!").
		BreakStart("if Subaru wake up").
		SyncResponse("Subaru", "David", "......").
		BreakEnd().
		LoopEnd().
		LF().
		SyncResponse("David", "Sophia", "wake up, wake up").
		String()

	md.NewMarkdown(os.Stdout).
		H2("Sequence Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build() //nolint

	// Output:
	// ## Sequence Diagram
	// ```mermaid
	// sequenceDiagram
	//     participant Sophia
	//     participant David
	//     participant Subaru
	//
	//     Sophia->>David: Please wake up Subaru
	//     David-->>Sophia: OK
	//
	//     loop until Subaru wake up
	//     David->>Subaru: Wake up!
	//     Subaru-->>David: zzz
	//     David->>Subaru: Hey!!!
	//     break if Subaru wake up
	//     Subaru-->>David: ......
	//     end
	//     end
	//
	//     David-->>Sophia: wake up, wake up
	// ```
}
