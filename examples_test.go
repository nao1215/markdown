//go:build linux || darwin

package markdown_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	md "github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/piechart"
	"github.com/nao1215/markdown/mermaid/sequence"
)

// ExampleMarkdown skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleMarkdown() {
	_ = md.NewMarkdown(os.Stdout).
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
	//
	// 1. Ordered Item 1
	// 2. Ordered Item 2
	// 3. Ordered Item 3
	//
	// ## CheckBox
	// - [ ] `sample code`
	// - [x] [Go](https://golang.org)
	// - [ ] ~~strikethrough~~
	//
	// ## Blockquote
	// > If you can dream it, you can do it.
	//
	// ### Horizontal Rule
	// ---
	// ## Table
	// | Name | Age | Country |
	// |---------|---------|---------|
	// | David | 23 | USA |
	// | John | 30 | UK |
	// | Bob | 25 | Canada |
	//
	// ## Image
	// ![sample_image](./sample.png)
}

// ExampleNewDiagram skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleNewDiagram() {
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

	_ = md.NewMarkdown(os.Stdout).
		H2("Sequence Diagram").
		CodeBlocks(md.SyntaxHighlightMermaid, diagram).
		Build()

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

// ExampleTableAlignment demonstrates table alignment features.
func ExampleTableAlignment() {
	_ = md.NewMarkdown(os.Stdout).
		H2("Table with Alignments").
		Table(md.TableSet{
			Header: []string{"Left Align", "Center Align", "Right Align"},
			Rows: [][]string{
				{"Content1", "Content2", "Content3"},
				{"Content4", "Content5", "Content6"},
			},
			Alignment: []md.TableAlignment{md.AlignLeft, md.AlignCenter, md.AlignRight},
		}).
		Build()

	// Output:
	// ## Table with Alignments
	// | Left Align | Center Align | Right Align |
	// |:--------|:-------:|--------:|
	// | Content1 | Content2 | Content3 |
	// | Content4 | Content5 | Content6 |
}

// ExampleNewMarkdown shows the shape every document has: a
// writer, a chain of calls, and Build.
func ExampleNewMarkdown() {
	_ = md.NewMarkdown(os.Stdout).
		H1("Release Notes").
		PlainText("Everything that changed in this release.").
		Build()

	// Output:
	// # Release Notes
	// Everything that changed in this release.
}

// ExampleNewMarkdown_writerIsAnyWriter shows that the destination is an
// io.Writer, so a document can be built into memory and used as a string.
func ExampleNewMarkdown_writerIsAnyWriter() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).H2("Summary").Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "## Summary\n"
}

// ExampleWithBlockSpacing shows the option that puts a blank line between
// blocks. Without it the blocks are written one after another.
func ExampleWithBlockSpacing() {
	_ = md.NewMarkdown(os.Stdout, md.WithBlockSpacing()).
		H2("Deploy").
		PlainText("Runs on every merge to main.").
		H2("Rollback").
		PlainText("Re-run the previous release job.").
		Build()

	// Output:
	// ## Deploy
	//
	// Runs on every merge to main.
	//
	// ## Rollback
	//
	// Re-run the previous release job.
}

// ExampleMarkdown_H1 writes a level 1 heading.
func ExampleMarkdown_H1() {
	_ = md.NewMarkdown(os.Stdout).H1("Heading").Build()

	// Output:
	// # Heading
}

// ExampleMarkdown_H1f writes a level 1 heading from a format string.
func ExampleMarkdown_H1f() {
	_ = md.NewMarkdown(os.Stdout).H1f("Heading %d", 1).Build()

	// Output:
	// # Heading 1
}

// ExampleMarkdown_H2 writes a level 2 heading.
func ExampleMarkdown_H2() {
	_ = md.NewMarkdown(os.Stdout).H2("Heading").Build()

	// Output:
	// ## Heading
}

// ExampleMarkdown_H2f writes a level 2 heading from a format string.
func ExampleMarkdown_H2f() {
	_ = md.NewMarkdown(os.Stdout).H2f("Heading %d", 2).Build()

	// Output:
	// ## Heading 2
}

// ExampleMarkdown_H3 writes a level 3 heading.
func ExampleMarkdown_H3() {
	_ = md.NewMarkdown(os.Stdout).H3("Heading").Build()

	// Output:
	// ### Heading
}

// ExampleMarkdown_H3f writes a level 3 heading from a format string.
func ExampleMarkdown_H3f() {
	_ = md.NewMarkdown(os.Stdout).H3f("Heading %d", 3).Build()

	// Output:
	// ### Heading 3
}

// ExampleMarkdown_H4 writes a level 4 heading.
func ExampleMarkdown_H4() {
	_ = md.NewMarkdown(os.Stdout).H4("Heading").Build()

	// Output:
	// #### Heading
}

// ExampleMarkdown_H4f writes a level 4 heading from a format string.
func ExampleMarkdown_H4f() {
	_ = md.NewMarkdown(os.Stdout).H4f("Heading %d", 4).Build()

	// Output:
	// #### Heading 4
}

// ExampleMarkdown_H5 writes a level 5 heading.
func ExampleMarkdown_H5() {
	_ = md.NewMarkdown(os.Stdout).H5("Heading").Build()

	// Output:
	// ##### Heading
}

// ExampleMarkdown_H5f writes a level 5 heading from a format string.
func ExampleMarkdown_H5f() {
	_ = md.NewMarkdown(os.Stdout).H5f("Heading %d", 5).Build()

	// Output:
	// ##### Heading 5
}

// ExampleMarkdown_H6 writes a level 6 heading.
func ExampleMarkdown_H6() {
	_ = md.NewMarkdown(os.Stdout).H6("Heading").Build()

	// Output:
	// ###### Heading
}

// ExampleMarkdown_H6f writes a level 6 heading from a format string.
func ExampleMarkdown_H6f() {
	_ = md.NewMarkdown(os.Stdout).H6f("Heading %d", 6).Build()

	// Output:
	// ###### Heading 6
}

// ExampleMarkdown_PlainText writes a paragraph.
func ExampleMarkdown_PlainText() {
	_ = md.NewMarkdown(os.Stdout).PlainText("A paragraph of text.").Build()

	// Output:
	// A paragraph of text.
}

// ExampleMarkdown_PlainTextf writes a paragraph from a format string.
func ExampleMarkdown_PlainTextf() {
	_ = md.NewMarkdown(os.Stdout).PlainTextf("Built %d documents in %s.", 3, "2s").Build()

	// Output:
	// Built 3 documents in 2s.
}

// ExampleMarkdown_Blockquote writes a quotation. Each line of the text is
// prefixed, so a quotation spanning lines stays one block.
func ExampleMarkdown_Blockquote() {
	_ = md.NewMarkdown(os.Stdout).
		Blockquote("If you can dream it, you can do it.").
		Build()

	// Output:
	// > If you can dream it, you can do it.
}

// ExampleMarkdown_CodeBlocks writes a fenced code block tagged with its
// language.
func ExampleMarkdown_CodeBlocks() {
	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightGo, `fmt.Println("hello")`).
		Build()

	// Output:
	// ```go
	// fmt.Println("hello")
	// ```
}

// ExampleMarkdown_HorizontalRule writes a thematic break.
func ExampleMarkdown_HorizontalRule() {
	_ = md.NewMarkdown(os.Stdout).
		PlainText("Above").
		HorizontalRule().
		PlainText("Below").
		Build()

	// Output:
	// Above
	// ---
	// Below
}

// ExampleMarkdown_BlankLine writes an empty line, for the places a document
// needs one that the block spacing does not give it.
func ExampleMarkdown_BlankLine() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).PlainText("Above").BlankLine().PlainText("Below").Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the empty line carries the two spaces markdown
	// reads as a hard line break, and a godoc Output block cannot hold
	// trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "Above\n\nBelow\n"
}

// ExampleMarkdown_LF is the older name for BlankLine and does the same thing.
func ExampleMarkdown_LF() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).PlainText("Above").LF().PlainText("Below").Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the empty line carries the two spaces markdown
	// reads as a hard line break, and a godoc Output block cannot hold
	// trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "Above\n  \nBelow\n"
}

// ExampleMarkdown_Details writes a collapsible section. Markdown has no
// syntax of its own for one, so this is the single place the library emits
// HTML.
func ExampleMarkdown_Details() {
	_ = md.NewMarkdown(os.Stdout).
		Details("Show the stack trace", "goroutine 1 [running]").
		Build()

	// Output:
	// <details>
	// <summary>Show the stack trace</summary>
	//
	// goroutine 1 [running]
	//
	// </details>
}

// ExampleMarkdown_Detailsf writes a collapsible section from a format string.
func ExampleMarkdown_Detailsf() {
	_ = md.NewMarkdown(os.Stdout).
		Detailsf("Show the log", "exited with status %d", 2).
		Build()

	// Output:
	// <details>
	// <summary>Show the log</summary>
	//
	// exited with status 2
	//
	// </details>
}

// ExampleMarkdown_BulletList writes an unordered list.
func ExampleMarkdown_BulletList() {
	_ = md.NewMarkdown(os.Stdout).
		BulletList("Read the spec", "Write the test", "Write the code").
		Build()

	// Output:
	// - Read the spec
	// - Write the test
	// - Write the code
}

// ExampleMarkdown_OrderedList writes a numbered list. The numbers are
// written out, so the document reads the same as a plain file.
func ExampleMarkdown_OrderedList() {
	_ = md.NewMarkdown(os.Stdout).
		OrderedList("Clone the repository", "Run the tests", "Open a pull request").
		Build()

	// Output:
	// 1. Clone the repository
	// 2. Run the tests
	// 3. Open a pull request
}

// ExampleMarkdown_CheckBox writes a task list. GitHub renders each item as a
// checkbox, and the text of an item may hold any inline markup.
func ExampleMarkdown_CheckBox() {
	_ = md.NewMarkdown(os.Stdout).
		CheckBox([]md.CheckBoxSet{
			{Checked: true, Text: "Write the proposal"},
			{Checked: false, Text: md.Code("go test ./...")},
		}).
		Build()

	// Output:
	// - [x] Write the proposal
	// - [ ] `go test ./...`
}

// ExampleMarkdown_Table writes a table. Every row must have as many cells as
// the header, and a row that does not records an error rather than writing a
// broken table.
func ExampleMarkdown_Table() {
	_ = md.NewMarkdown(os.Stdout).
		Table(md.TableSet{
			Header: []string{"Name", "Age"},
			Rows: [][]string{
				{"David", "23"},
				{"John", "30"},
			},
		}).
		Build()

	// Output:
	// | Name | Age |
	// |---------|---------|
	// | David | 23 |
	// | John | 30 |
}

// ExampleMarkdown_CustomTable writes a table with the alignment row spelled
// out. Without the options every column is left aligned.
func ExampleMarkdown_CustomTable() {
	_ = md.NewMarkdown(os.Stdout).
		CustomTable(md.TableSet{
			Header: []string{"Name", "Age"},
			Rows: [][]string{
				{"David", "23"},
				{"John", "30"},
			},
		}, md.TableOptions{
			AutoFormatHeaders: false,
			AutoWrapText:      false,
		}).
		Build()

	// Output:
	// | Name  | Age |
	// |-------|-----|
	// | David | 23  |
	// | John  | 30  |
}

// ExampleTableSet_ValidateColumns reports a row whose cell count does not
// match the header. Table calls it, so this is only worth calling directly
// when the table is assembled before the document is.
func ExampleTableSet_ValidateColumns() {
	set := md.TableSet{
		Header: []string{"Name", "Age"},
		Rows:   [][]string{{"only one cell"}},
	}
	fmt.Println(set.ValidateColumns())

	// Output:
	// number of columns in the record doesn't match the header
}

// ExampleEscapeTableCell escapes the characters that would end a table cell
// early. A pipe closes the cell it is written in, so text arriving from a
// database or a command's output needs this before it reaches a row.
func ExampleEscapeTableCell() {
	fmt.Println(md.EscapeTableCell("a|b"))
	fmt.Println(md.EscapeTableCell("multi\nline"))

	// Output:
	// a\|b
	// multi<br>line
}

// ExampleMarkdown_TableOfContents writes a table of contents built from the
// headings of the document. It may be called before the headings it lists:
// the list is filled in at Build.
func ExampleMarkdown_TableOfContents() {
	_ = md.NewMarkdown(os.Stdout).
		H1("Guide").
		TableOfContents(md.TableOfContentsDepthH3).
		H2("Install").
		H3("From source").
		H2("Usage").
		Build()

	// Output:
	// # Guide
	// <!-- BEGIN_TOC -->
	// - [Guide](#guide)
	//   - [Install](#install)
	//     - [From source](#from-source)
	//   - [Usage](#usage)
	// <!-- END_TOC -->
	//
	// ## Install
	// ### From source
	// ## Usage
}

// ExampleMarkdown_TableOfContentsWithRange writes a table of contents holding
// only the heading levels between the two given, which is how a document
// leaves its own title out of its contents.
func ExampleMarkdown_TableOfContentsWithRange() {
	_ = md.NewMarkdown(os.Stdout).
		H1("Guide").
		TableOfContentsWithRange(md.TableOfContentsDepthH2, md.TableOfContentsDepthH2).
		H2("Install").
		H3("From source").
		H2("Usage").
		Build()

	// Output:
	// # Guide
	// <!-- BEGIN_TOC -->
	// - [Install](#install)
	// - [Usage](#usage)
	// <!-- END_TOC -->
	//
	// ## Install
	// ### From source
	// ## Usage
}

// ExampleMarkdown_Note writes a GitHub NOTE alert.
func ExampleMarkdown_Note() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).Note("Useful information a reader should know even when skimming.").Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the keyword line ends with the two spaces
	// markdown reads as a hard line break, and a godoc Output block cannot
	// hold trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "> [!NOTE]  \n> Useful information a reader should know even when skimming.\n"
}

// ExampleMarkdown_Notef writes a GitHub NOTE alert from a format string.
func ExampleMarkdown_Notef() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).Notef("Note takes %d minutes.", 5).Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the keyword line ends with the two spaces
	// markdown reads as a hard line break, and a godoc Output block cannot
	// hold trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "> [!NOTE]  \n> Note takes 5 minutes.\n"
}

// ExampleMarkdown_Tip writes a GitHub TIP alert.
func ExampleMarkdown_Tip() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).Tip("Helpful advice for doing things better.").Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the keyword line ends with the two spaces
	// markdown reads as a hard line break, and a godoc Output block cannot
	// hold trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "> [!TIP]  \n> Helpful advice for doing things better.\n"
}

// ExampleMarkdown_Tipf writes a GitHub TIP alert from a format string.
func ExampleMarkdown_Tipf() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).Tipf("Tip takes %d minutes.", 5).Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the keyword line ends with the two spaces
	// markdown reads as a hard line break, and a godoc Output block cannot
	// hold trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "> [!TIP]  \n> Tip takes 5 minutes.\n"
}

// ExampleMarkdown_Important writes a GitHub IMPORTANT alert.
func ExampleMarkdown_Important() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).Important("Key information a reader needs to succeed.").Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the keyword line ends with the two spaces
	// markdown reads as a hard line break, and a godoc Output block cannot
	// hold trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "> [!IMPORTANT]  \n> Key information a reader needs to succeed.\n"
}

// ExampleMarkdown_Importantf writes a GitHub IMPORTANT alert from a format string.
func ExampleMarkdown_Importantf() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).Importantf("Important takes %d minutes.", 5).Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the keyword line ends with the two spaces
	// markdown reads as a hard line break, and a godoc Output block cannot
	// hold trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "> [!IMPORTANT]  \n> Important takes 5 minutes.\n"
}

// ExampleMarkdown_Warning writes a GitHub WARNING alert.
func ExampleMarkdown_Warning() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).Warning("Urgent information needing immediate attention.").Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the keyword line ends with the two spaces
	// markdown reads as a hard line break, and a godoc Output block cannot
	// hold trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "> [!WARNING]  \n> Urgent information needing immediate attention.\n"
}

// ExampleMarkdown_Warningf writes a GitHub WARNING alert from a format string.
func ExampleMarkdown_Warningf() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).Warningf("Warning takes %d minutes.", 5).Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the keyword line ends with the two spaces
	// markdown reads as a hard line break, and a godoc Output block cannot
	// hold trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "> [!WARNING]  \n> Warning takes 5 minutes.\n"
}

// ExampleMarkdown_Caution writes a GitHub CAUTION alert.
func ExampleMarkdown_Caution() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).Caution("Advises about the risks of an action.").Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the keyword line ends with the two spaces
	// markdown reads as a hard line break, and a godoc Output block cannot
	// hold trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "> [!CAUTION]  \n> Advises about the risks of an action.\n"
}

// ExampleMarkdown_Cautionf writes a GitHub CAUTION alert from a format string.
func ExampleMarkdown_Cautionf() {
	buf := &bytes.Buffer{}
	if err := md.NewMarkdown(buf).Cautionf("Caution takes %d minutes.", 5).Build(); err != nil {
		fmt.Println("build:", err)
		return
	}
	// Printed quoted because the keyword line ends with the two spaces
	// markdown reads as a hard line break, and a godoc Output block cannot
	// hold trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "> [!CAUTION]  \n> Caution takes 5 minutes.\n"
}

// ExampleMarkdown_RedBadge writes a red badge, which is an image served by
// shields.io rather than markdown of its own.
func ExampleMarkdown_RedBadge() {
	_ = md.NewMarkdown(os.Stdout).RedBadge("build").Build()

	// Output:
	// ![Badge](https://img.shields.io/badge/build-red)
}

// ExampleMarkdown_RedBadgef writes a red badge from a format string.
func ExampleMarkdown_RedBadgef() {
	_ = md.NewMarkdown(os.Stdout).RedBadgef("coverage %d%%", 96).Build()

	// Output:
	// ![Badge](https://img.shields.io/badge/coverage 96%-red)
}

// ExampleMarkdown_YellowBadge writes a yellow badge, which is an image served by
// shields.io rather than markdown of its own.
func ExampleMarkdown_YellowBadge() {
	_ = md.NewMarkdown(os.Stdout).YellowBadge("build").Build()

	// Output:
	// ![Badge](https://img.shields.io/badge/build-yellow)
}

// ExampleMarkdown_YellowBadgef writes a yellow badge from a format string.
func ExampleMarkdown_YellowBadgef() {
	_ = md.NewMarkdown(os.Stdout).YellowBadgef("coverage %d%%", 96).Build()

	// Output:
	// ![Badge](https://img.shields.io/badge/coverage 96%-yellow)
}

// ExampleMarkdown_GreenBadge writes a green badge, which is an image served by
// shields.io rather than markdown of its own.
func ExampleMarkdown_GreenBadge() {
	_ = md.NewMarkdown(os.Stdout).GreenBadge("build").Build()

	// Output:
	// ![Badge](https://img.shields.io/badge/build-green)
}

// ExampleMarkdown_GreenBadgef writes a green badge from a format string.
func ExampleMarkdown_GreenBadgef() {
	_ = md.NewMarkdown(os.Stdout).GreenBadgef("coverage %d%%", 96).Build()

	// Output:
	// ![Badge](https://img.shields.io/badge/coverage 96%-green)
}

// ExampleMarkdown_BlueBadge writes a blue badge, which is an image served by
// shields.io rather than markdown of its own.
func ExampleMarkdown_BlueBadge() {
	_ = md.NewMarkdown(os.Stdout).BlueBadge("build").Build()

	// Output:
	// ![Badge](https://img.shields.io/badge/build-blue)
}

// ExampleMarkdown_BlueBadgef writes a blue badge from a format string.
func ExampleMarkdown_BlueBadgef() {
	_ = md.NewMarkdown(os.Stdout).BlueBadgef("coverage %d%%", 96).Build()

	// Output:
	// ![Badge](https://img.shields.io/badge/coverage 96%-blue)
}

// ExampleBold returns the inline markup rather than writing it, so it can be
// put inside any text a builder takes.
func ExampleBold() {
	_ = md.NewMarkdown(os.Stdout).
		PlainTextf("This word is %s.", md.Bold("important")).
		Build()

	// Output:
	// This word is **important**.
}

// ExampleItalic returns the inline markup rather than writing it, so it can be
// put inside any text a builder takes.
func ExampleItalic() {
	_ = md.NewMarkdown(os.Stdout).
		PlainTextf("This word is %s.", md.Italic("emphasis")).
		Build()

	// Output:
	// This word is *emphasis*.
}

// ExampleBoldItalic returns the inline markup rather than writing it, so it can be
// put inside any text a builder takes.
func ExampleBoldItalic() {
	_ = md.NewMarkdown(os.Stdout).
		PlainTextf("This word is %s.", md.BoldItalic("both")).
		Build()

	// Output:
	// This word is ***both***.
}

// ExampleCode returns the inline markup rather than writing it, so it can be
// put inside any text a builder takes.
func ExampleCode() {
	_ = md.NewMarkdown(os.Stdout).
		PlainTextf("This word is %s.", md.Code("go test ./...")).
		Build()

	// Output:
	// This word is `go test ./...`.
}

// ExampleStrikethrough returns the inline markup rather than writing it, so it can be
// put inside any text a builder takes.
func ExampleStrikethrough() {
	_ = md.NewMarkdown(os.Stdout).
		PlainTextf("This word is %s.", md.Strikethrough("removed")).
		Build()

	// Output:
	// This word is ~~removed~~.
}

// ExampleHighlight returns the inline markup rather than writing it, so it can be
// put inside any text a builder takes.
func ExampleHighlight() {
	_ = md.NewMarkdown(os.Stdout).
		PlainTextf("This word is %s.", md.Highlight("marked")).
		Build()

	// Output:
	// This word is ==marked==.
}

// ExampleLink returns an inline link.
func ExampleLink() {
	_ = md.NewMarkdown(os.Stdout).
		PlainText(md.Link("The Go Programming Language", "https://go.dev")).
		Build()

	// Output:
	// [The Go Programming Language](https://go.dev)
}

// ExampleImage returns an inline image.
func ExampleImage() {
	_ = md.NewMarkdown(os.Stdout).
		PlainText(md.Image("The Go gopher", "./gopher.png")).
		Build()

	// Output:
	// ![The Go gopher](./gopher.png)
}

// ExampleReferenceLink returns a link that points at a definition written
// elsewhere in the document, which keeps a long URL out of the sentence.
func ExampleReferenceLink() {
	_ = md.NewMarkdown(os.Stdout).
		PlainText(md.ReferenceLink("The Go Programming Language", "go")).
		PlainText(md.ReferenceLinkDefinition("go", "https://go.dev")).
		Build()

	// Output:
	// [The Go Programming Language][go]
	// [go]: https://go.dev
}

// ExampleReferenceLinkDefinition writes the definition a reference link points
// at. The optional third argument is the title a browser shows on hover.
func ExampleReferenceLinkDefinition() {
	fmt.Println(md.ReferenceLinkDefinition("go", "https://go.dev"))
	fmt.Println(md.ReferenceLinkDefinition("go", "https://go.dev", "The Go website"))

	// Output:
	// [go]: https://go.dev
	// [go]: https://go.dev "The Go website"
}

// ExampleFootnoteReference returns the marker that points at a footnote.
func ExampleFootnoteReference() {
	_ = md.NewMarkdown(os.Stdout).
		PlainTextf("Generated by this library%s.", md.FootnoteReference("1")).
		PlainText(md.FootnoteDefinition("1", "github.com/nao1215/markdown")).
		Build()

	// Output:
	// Generated by this library[^1].
	// [^1]: github.com/nao1215/markdown
}

// ExampleFootnoteDefinition writes the text a footnote reference points at.
func ExampleFootnoteDefinition() {
	fmt.Println(md.FootnoteDefinition("1", "github.com/nao1215/markdown"))

	// Output:
	// [^1]: github.com/nao1215/markdown
}

// ExampleInlineMath returns a mathematical expression that sits inside a
// sentence. GitHub renders it with KaTeX.
func ExampleInlineMath() {
	_ = md.NewMarkdown(os.Stdout).
		PlainTextf("The area is %s.", md.InlineMath("\\pi r^2")).
		Build()

	// Output:
	// The area is $\pi r^2$.
}

// ExampleBlockMath returns a mathematical expression that stands on its own.
func ExampleBlockMath() {
	_ = md.NewMarkdown(os.Stdout).
		PlainText(md.BlockMath("\\int_0^1 x^2 dx = \\frac{1}{3}")).
		Build()

	// Output:
	// $$
	// \int_0^1 x^2 dx = \frac{1}{3}
	// $$
}

// ExampleMarkdown_Build writes the document and reports the first error the
// chain recorded. Nothing in the chain panics on bad input, so one check at
// the end is enough.
func ExampleMarkdown_Build() {
	err := md.NewMarkdown(os.Stdout).
		H1("Report").
		Table(md.TableSet{
			Header: []string{"Name", "Age"},
			Rows:   [][]string{{"only one cell"}},
		}).
		Build()
	fmt.Println("error:", err)

	// Output:
	// # Report
	// error: failed to validate columns: number of columns in the record doesn't match the header
}

// ExampleMarkdown_Error reports the same error Build does, for code that wants
// to look before writing anything.
func ExampleMarkdown_Error() {
	m := md.NewMarkdown(io.Discard).
		TableOfContents(md.TableOfContentsDepthH3).
		TableOfContents(md.TableOfContentsDepthH3)
	fmt.Println("error:", m.Error())

	// Output:
	// error: table of contents has already been generated
}

// ExampleMarkdown_String returns the document built so far without needing a
// writer, which is how the mermaid subpackages hand a diagram to CodeBlocks.
func ExampleMarkdown_String() {
	document := md.NewMarkdown(io.Discard).H2("Summary").PlainText("All green.").String()
	fmt.Printf("%q\n", document)

	// Output:
	// "## Summary\nAll green."
}

// ExampleGenerateIndex writes an index of the markdown files under a directory.
// WithWriter sends it somewhere other than the index.md the function would
// otherwise create, which is what makes the output here worth showing.
func ExampleGenerateIndex() {
	parent, err := os.MkdirTemp("", "markdown-index")
	if err != nil {
		fmt.Println("temp dir:", err)
		return
	}
	defer func() { _ = os.RemoveAll(parent) }()

	// The index is headed with the name of the directory it describes, so the
	// directory is named here rather than left as the random one MkdirTemp
	// makes, which would put a different heading in the output on every run.
	dir := filepath.Join(parent, "guide")
	if err := os.Mkdir(dir, 0o750); err != nil {
		fmt.Println("mkdir:", err)
		return
	}
	if err := os.WriteFile(filepath.Join(dir, "install.md"), []byte("# Install\n"), 0o600); err != nil {
		fmt.Println("write:", err)
		return
	}

	if err := md.GenerateIndex(dir, md.WithWriter(os.Stdout)); err != nil {
		fmt.Println("generate:", err)
	}

	// Output:
	// ### guide
	// - [Install](install.md)
}

// ExampleWithTitle sets the heading the generated index opens with.
func ExampleWithTitle() {
	parent, err := os.MkdirTemp("", "markdown-index")
	if err != nil {
		fmt.Println("temp dir:", err)
		return
	}
	defer func() { _ = os.RemoveAll(parent) }()

	// The index is headed with the name of the directory it describes, so the
	// directory is named here rather than left as the random one MkdirTemp
	// makes, which would put a different heading in the output on every run.
	dir := filepath.Join(parent, "guide")
	if err := os.Mkdir(dir, 0o750); err != nil {
		fmt.Println("mkdir:", err)
		return
	}
	if err := os.WriteFile(filepath.Join(dir, "install.md"), []byte("# Install\n"), 0o600); err != nil {
		fmt.Println("write:", err)
		return
	}

	if err := md.GenerateIndex(dir, md.WithTitle("Documentation"), md.WithWriter(os.Stdout)); err != nil {
		fmt.Println("generate:", err)
	}

	// Output:
	// ## Documentation
	// ### guide
	// - [Install](install.md)
}

// ExampleWithDescription sets the paragraphs written under the index heading.
func ExampleWithDescription() {
	parent, err := os.MkdirTemp("", "markdown-index")
	if err != nil {
		fmt.Println("temp dir:", err)
		return
	}
	defer func() { _ = os.RemoveAll(parent) }()

	// The index is headed with the name of the directory it describes, so the
	// directory is named here rather than left as the random one MkdirTemp
	// makes, which would put a different heading in the output on every run.
	dir := filepath.Join(parent, "guide")
	if err := os.Mkdir(dir, 0o750); err != nil {
		fmt.Println("mkdir:", err)
		return
	}
	if err := os.WriteFile(filepath.Join(dir, "install.md"), []byte("# Install\n"), 0o600); err != nil {
		fmt.Println("write:", err)
		return
	}

	buf := &bytes.Buffer{}
	err = md.GenerateIndex(dir,
		md.WithDescription([]string{"Every page in this directory.", "Regenerated on each release."}),
		md.WithWriter(buf),
	)
	if err != nil {
		fmt.Println("generate:", err)
		return
	}
	// Printed quoted because each description line ends with the two spaces
	// markdown reads as a hard line break, and a godoc Output block cannot hold
	// trailing whitespace.
	fmt.Printf("%q\n", buf.String())

	// Output:
	// "Every page in this directory.\n  \nRegenerated on each release.\n  \n### guide\n- [Install](install.md)\n  \n"
}

// ExampleWithWriter sends the index somewhere other than the index.md the
// function creates by default, so it can be inspected or embedded rather than
// written to disk.
func ExampleWithWriter() {
	parent, err := os.MkdirTemp("", "markdown-index")
	if err != nil {
		fmt.Println("temp dir:", err)
		return
	}
	defer func() { _ = os.RemoveAll(parent) }()

	// The index is headed with the name of the directory it describes, so the
	// directory is named here rather than left as the random one MkdirTemp
	// makes, which would put a different heading in the output on every run.
	dir := filepath.Join(parent, "guide")
	if err := os.Mkdir(dir, 0o750); err != nil {
		fmt.Println("mkdir:", err)
		return
	}
	if err := os.WriteFile(filepath.Join(dir, "install.md"), []byte("# Install\n"), 0o600); err != nil {
		fmt.Println("write:", err)
		return
	}

	buf := &bytes.Buffer{}
	if err := md.GenerateIndex(dir, md.WithWriter(buf)); err != nil {
		fmt.Println("generate:", err)
		return
	}
	fmt.Print(buf.String())

	// Output:
	// ### guide
	// - [Install](install.md)
}

// ExampleTableSet shows the shape a table is described with. Every row must
// have as many cells as the header.
func ExampleTableSet() {
	set := md.TableSet{
		Header: []string{"Package", "Coverage"},
		Rows: [][]string{
			{"markdown", "96%"},
			{"internal", "100%"},
		},
	}

	_ = md.NewMarkdown(os.Stdout).Table(set).Build()

	// Output:
	// | Package | Coverage |
	// |---------|---------|
	// | markdown | 96% |
	// | internal | 100% |
}

// ExampleTableOptions shows the options CustomTable takes. They control the
// header casing and the wrapping, not the alignment, which is a field of the
// table itself.
func ExampleTableOptions() {
	_ = md.NewMarkdown(os.Stdout).
		CustomTable(md.TableSet{
			Header: []string{"package", "coverage"},
			Rows:   [][]string{{"markdown", "96%"}},
		}, md.TableOptions{
			AutoFormatHeaders: true,
			AutoWrapText:      false,
		}).
		Build()

	// Output:
	// | PACKAGE  | COVERAGE |
	// |----------|----------|
	// | MARKDOWN | 96 %     |
}

// ExampleTableOfContentsOptions shows the range a table of contents covers.
// Naming both ends is how a document leaves its own title out of its contents.
func ExampleTableOfContentsOptions() {
	options := md.TableOfContentsOptions{
		MinDepth: md.TableOfContentsDepthH2,
		MaxDepth: md.TableOfContentsDepthH2,
	}

	_ = md.NewMarkdown(os.Stdout).
		H1("Guide").
		TableOfContentsWithRange(options.MinDepth, options.MaxDepth).
		H2("Install").
		H3("From source").
		Build()

	// Output:
	// # Guide
	// <!-- BEGIN_TOC -->
	// - [Install](#install)
	// <!-- END_TOC -->
	//
	// ## Install
	// ### From source
}

// ExampleTableOfContentsDepth shows the heading level a table of contents stops
// at.
func ExampleTableOfContentsDepth() {
	_ = md.NewMarkdown(os.Stdout).
		TableOfContents(md.TableOfContentsDepthH2).
		H2("Install").
		H3("From source").
		Build()

	// Output:
	// <!-- BEGIN_TOC -->
	// - [Install](#install)
	// <!-- END_TOC -->
	//
	// ## Install
	// ### From source
}

// ExampleCheckBoxSet shows the shape one item of a task list is described with.
func ExampleCheckBoxSet() {
	_ = md.NewMarkdown(os.Stdout).
		CheckBox([]md.CheckBoxSet{
			{Checked: true, Text: "Write the proposal"},
			{Checked: false, Text: "Get it reviewed"},
		}).
		Build()

	// Output:
	// - [x] Write the proposal
	// - [ ] Get it reviewed
}

// ExampleSyntaxHighlight shows the language a code block is tagged with. The
// constants cover the languages GitHub highlights; any other string works too.
func ExampleSyntaxHighlight() {
	_ = md.NewMarkdown(os.Stdout).
		CodeBlocks(md.SyntaxHighlightGo, `fmt.Println("hello")`).
		CodeBlocks(md.SyntaxHighlightNone, "no highlighting here").
		Build()

	// Output:
	// ```go
	// fmt.Println("hello")
	// ```
	// ```
	// no highlighting here
	// ```
}

// ExampleOption shows what an Option is: a function that changes how a document
// is written, passed to NewMarkdown.
func ExampleOption() {
	options := []md.Option{md.WithBlockSpacing()}

	_ = md.NewMarkdown(os.Stdout, options...).
		H2("Deploy").
		PlainText("Runs on every merge to main.").
		Build()

	// Output:
	// ## Deploy
	//
	// Runs on every merge to main.
}

// ExampleIndexOption shows what an IndexOption is: a function that changes what
// GenerateIndex writes, passed to it after the directory.
func ExampleIndexOption() {
	parent, err := os.MkdirTemp("", "markdown-index")
	if err != nil {
		fmt.Println("temp dir:", err)
		return
	}
	defer func() { _ = os.RemoveAll(parent) }()

	dir := filepath.Join(parent, "guide")
	if err := os.Mkdir(dir, 0o750); err != nil {
		fmt.Println("mkdir:", err)
		return
	}
	if err := os.WriteFile(filepath.Join(dir, "install.md"), []byte("# Install\n"), 0o600); err != nil {
		fmt.Println("write:", err)
		return
	}

	options := []md.IndexOption{md.WithTitle("Documentation"), md.WithWriter(os.Stdout)}
	if err := md.GenerateIndex(dir, options...); err != nil {
		fmt.Println("generate:", err)
	}

	// Output:
	// ## Documentation
	// ### guide
	// - [Install](install.md)
}

// ExampleIndex shows where an Index comes from. The type carries what
// GenerateIndex collected, and nothing exported reaches inside it: the index is
// written rather than inspected.
func ExampleIndex() {
	parent, err := os.MkdirTemp("", "markdown-index")
	if err != nil {
		fmt.Println("temp dir:", err)
		return
	}
	defer func() { _ = os.RemoveAll(parent) }()

	dir := filepath.Join(parent, "guide")
	if err := os.Mkdir(dir, 0o750); err != nil {
		fmt.Println("mkdir:", err)
		return
	}
	if err := os.WriteFile(filepath.Join(dir, "usage.md"), []byte("# Usage\n"), 0o600); err != nil {
		fmt.Println("write:", err)
		return
	}

	if err := md.GenerateIndex(dir, md.WithWriter(os.Stdout)); err != nil {
		fmt.Println("generate:", err)
	}

	// Output:
	// ### guide
	// - [Usage](usage.md)
}

// Example_gitHubActionsJobSummary writes the markdown a GitHub Actions job
// summary is made of: inside a step, open the file named by the
// GITHUB_STEP_SUMMARY environment variable for appending and hand it to
// NewMarkdown in place of os.Stdout; the run's summary page renders the
// result, mermaid diagrams included.
func Example_gitHubActionsJobSummary() {
	coverage := piechart.NewPieChart(
		io.Discard,
		piechart.WithTitle("Coverage"),
		piechart.WithShowData(true),
	).
		LabelAndIntValue("covered", 92).
		LabelAndIntValue("uncovered", 8).
		String()

	err := md.NewMarkdown(os.Stdout, md.WithBlockSpacing()).
		H2("Test Results").
		Table(md.TableSet{
			Header: []string{"Package", "Passed", "Failed"},
			Rows: [][]string{
				{"api", "120", "0"},
				{"core", "89", "2"},
			},
		}).
		CodeBlocks(md.SyntaxHighlightMermaid, coverage).
		Build()
	if err != nil {
		fmt.Println("build:", err)
	}

	// Output:
	// ## Test Results
	//
	// | Package | Passed | Failed |
	// |---------|---------|---------|
	// | api | 120 | 0 |
	// | core | 89 | 2 |
	//
	// ```mermaid
	// %%{init: {"pie": {"textPosition": 0.75}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
	// pie showData
	//     title Coverage
	//     "covered" : 92
	//     "uncovered" : 8
	// ```
}
