<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-5-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->
[![Go Reference](https://pkg.go.dev/badge/github.com/nao1215/markdown.svg)](https://pkg.go.dev/github.com/nao1215/markdown)
[![MultiPlatformUnitTest](https://github.com/nao1215/markdown/actions/workflows/unit_test.yml/badge.svg)](https://github.com/nao1215/markdown/actions/workflows/unit_test.yml)
[![reviewdog](https://github.com/nao1215/markdown/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/markdown/actions/workflows/reviewdog.yml)
[![Gosec](https://github.com/nao1215/markdown/actions/workflows/gosec.yml/badge.svg)](https://github.com/nao1215/markdown/actions/workflows/gosec.yml)
![Coverage](https://raw.githubusercontent.com/nao1215/octocovs-central-repo/main/badges/nao1215/markdown/coverage.svg)

[日本語](./doc/ja/README.md) | [Русский](./doc/ru/README.md) | [中文](./doc/zh-cn/README.md) | [한국어](./doc/ko/README.md) | [Español](./doc/es/README.md) | [Français](./doc/fr/README.md)

# What is markdown package
The Package markdown is a simple markdown builder in golang. The markdown package assembles Markdown using method chaining, not uses a template engine like [html/template](https://pkg.go.dev/html/template). The syntax of Markdown follows **GitHub Markdown**.
  
The markdown package was initially developed to save test results in [nao1215/spectest](https://github.com/nao1215/spectest). Therefore, the markdown package implements the features required by spectest. For example, the markdown package supports **mermaid sequence diagrams (entity relationship diagram, sequence diagram, flowchart, pie chart, state diagram, architecture diagram)**, which was a necessary feature in spectest.
  
Additionally, complex code that increases the complexity of the library, such as generating nested lists, will not be added. I want to keep this library as simple as possible.
  
## Supported OS and go version
- OS: Linux, macOS, Windows
- Go: 1.21 or later
  
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
- [x] Text with link
- [x] Text with image
- [x] Plain text
- [x] Details 
- [x] Alerts; NOTE, TIP, IMPORTANT, CAUTION, WARNING
- [x] mermaid sequence diagram
- [x] mermaid entity relationship diagram
- [x] mermaid flowchart 
- [x] mermaid pie chart
- [x] mermaid state diagram
- [x] mermaid architecture diagram (beta feature) 

### Features not in Markdown syntax
- Generate badges; RedBadge(), YellowBadge(), GreenBadge().
- Generate an index for a directory full of markdown files; GenerateIndex()
  
## Example
### Basic usage
```go
package main

import (
	"os"

	md "github.com/nao1215/markdown"
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

### Generate Markdown using `"go generate ./..."`
You can generate Markdown using `go generate`. Please define code to generate Markdown first. Then, run `"go generate ./..."` to generate Markdown.

[Code example:](./doc/generate/main.go)
```go
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
	defer f.Close()

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

### Alerts syntax
The markdown package can create alerts. Alerts are useful for displaying important information in Markdown. This syntax is supported by GitHub.
[Code example:](./doc/alert/main.go)
```go
	md.NewMarkdown(f).
		H1("Alert example").
		Note("This is note").LF().
		Tip("This is tip").LF().
		Important("This is important").LF().
		Warning("This is warning").LF().
		Caution("This is caution").LF().
		Build()
```
  
[Output:](./doc/alert/generated.md)
````text
# Alert example
> [!NOTE]  
> This is note
  
> [!TIP]  
> This is tip
  
> [!IMPORTANT]  
> This is important
  
> [!WARNING]  
> This is warning
  
> [!CAUTION]  
> This is caution
````

Your alert will look like this;
> [!NOTE]  
> This is note
  
> [!TIP]  
> This is tip
  
> [!IMPORTANT]  
> This is important
  
> [!WARNING]  
> This is warning
  
> [!CAUTION]  
> This is caution

### Status badge syntax
The markdown package can create red, yellow, and green status badges.
[Code example:](./doc/badge/main.go)
```go
	md.NewMarkdown(os.Stdout).
		H1("badge example").
		RedBadge("red_badge").
		YellowBadge("yellow_badge").
		GreenBadge("green_badge").
		BlueBadge("blue_badge").
		Build()
```

[Output:](./doc/badge/generated.md)
````text
# badge example
![Badge](https://img.shields.io/badge/red_badge-red)
![Badge](https://img.shields.io/badge/yellow_badge-yellow)
![Badge](https://img.shields.io/badge/green_badge-green)
![Badge](https://img.shields.io/badge/blue_badge-blue)
````

Your badge will look like this;  
![Badge](https://img.shields.io/badge/red_badge-red)
![Badge](https://img.shields.io/badge/yellow_badge-yellow)
![Badge](https://img.shields.io/badge/green_badge-green)
![Badge](https://img.shields.io/badge/blue_badge-blue)

### Mermaid sequence diagram syntax

```go
package main

import (
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/mermaid/sequence"
)

//go:generate go run main.go

func main() {
	diagram := sequence.NewDiagram(io.Discard).
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

	markdown.NewMarkdown(os.Stdout).
		H2("Sequence Diagram").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()
}
```

Plain text output: [markdown is here](./doc/sequence/generated.md)
````
## Sequence Diagram
```mermaid
sequenceDiagram
    participant Sophia
    participant David
    participant Subaru

    Sophia->>David: Please wake up Subaru
    David-->>Sophia: OK

    loop until Subaru wake up
    David->>Subaru: Wake up!
    Subaru-->>David: zzz
    David->>Subaru: Hey!!!
    break if Subaru wake up
    Subaru-->>David: ......
    end
    end

    David-->>Sophia: wake up, wake up
```
````

Mermaid output:
```mermaid
sequenceDiagram
    participant Sophia
    participant David
    participant Subaru

    Sophia->>David: Please wake up Subaru
    David-->>Sophia: OK

    loop until Subaru wake up
    David->>Subaru: Wake up!
    Subaru-->>David: zzz
    David->>Subaru: Hey!!!
    break if Subaru wake up
    Subaru-->>David: ......
    end
    end

    David-->>Sophia: wake up, wake up
```

### Entity Relationship Diagram syntax

```go
package main

import (
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/er"
)

//go:generate go run main.go

func main() {
	f, err := os.Create("generated.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	teachers := er.NewEntity(
		"teachers",
		[]*er.Attribute{
			{
				Type:         "int",
				Name:         "id",
				IsPrimaryKey: true,
				IsForeignKey: false,
				IsUniqueKey:  true,
				Comment:      "Teacher ID",
			},
			{
				Type:         "string",
				Name:         "name",
				IsPrimaryKey: false,
				IsForeignKey: false,
				IsUniqueKey:  false,
				Comment:      "Teacher Name",
			},
		},
	)
	students := er.NewEntity(
		"students",
		[]*er.Attribute{
			{
				Type:         "int",
				Name:         "id",
				IsPrimaryKey: true,
				IsForeignKey: false,
				IsUniqueKey:  true,
				Comment:      "Student ID",
			},
			{
				Type:         "string",
				Name:         "name",
				IsPrimaryKey: false,
				IsForeignKey: false,
				IsUniqueKey:  false,
				Comment:      "Student Name",
			},
			{
				Type:         "int",
				Name:         "teacher_id",
				IsPrimaryKey: false,
				IsForeignKey: true,
				IsUniqueKey:  true,
				Comment:      "Teacher ID",
			},
		},
	)
	schools := er.NewEntity(
		"schools",
		[]*er.Attribute{
			{
				Type:         "int",
				Name:         "id",
				IsPrimaryKey: true,
				IsForeignKey: false,
				IsUniqueKey:  true,
				Comment:      "School ID",
			},
			{
				Type:         "string",
				Name:         "name",
				IsPrimaryKey: false,
				IsForeignKey: false,
				IsUniqueKey:  false,
				Comment:      "School Name",
			},
			{
				Type:         "int",
				Name:         "teacher_id",
				IsPrimaryKey: false,
				IsForeignKey: true,
				IsUniqueKey:  true,
				Comment:      "Teacher ID",
			},
		},
	)

	erString := er.NewDiagram(f).
		Relationship(
			teachers,
			students,
			er.ExactlyOneRelationship, // "||"
			er.ZeroToMoreRelationship, // "}o"
			er.Identifying,            // "--"
			"Teacher has many students",
		).
		Relationship(
			teachers,
			schools,
			er.OneToMoreRelationship,  // "|}"
			er.ExactlyOneRelationship, // "||"
			er.NonIdentifying,         // ".."
			"School has many teachers",
		).
		String()

	err = markdown.NewMarkdown(f).
		H2("Entity Relationship Diagram").
		CodeBlocks(markdown.SyntaxHighlightMermaid, erString).
		Build()

	if err != nil {
		panic(err)
	}
}
```

Plain text output: [markdown is here](./doc/er/generated.md)
````
## Entity Relationship Diagram
```mermaid
erDiagram
	teachers ||--o{ students : "Teacher has many students"
	teachers }|..|| schools : "School has many teachers"
	schools {
		int id PK,UK "School ID"
		string name  "School Name"
		int teacher_id FK,UK "Teacher ID"
	}
	students {
		int id PK,UK "Student ID"
		string name  "Student Name"
		int teacher_id FK,UK "Teacher ID"
	}
	teachers {
		int id PK,UK "Teacher ID"
		string name  "Teacher Name"
	}

```
````

Mermaid output:
```mermaid
erDiagram
	teachers ||--o{ students : "Teacher has many students"
	teachers }|..|| schools : "School has many teachers"
	schools {
		int id PK,UK "School ID"
		string name  "School Name"
		int teacher_id FK,UK "Teacher ID"
	}
	students {
		int id PK,UK "Student ID"
		string name  "Student Name"
		int teacher_id FK,UK "Teacher ID"
	}
	teachers {
		int id PK,UK "Teacher ID"
		string name  "Teacher Name"
	}
```

### Flowchart syntax

```go
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/flowchart"
)

//go:generate go run main.go

func main() {
	f, err := os.Create("generated.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fc := flowchart.NewFlowchart(
		io.Discard,
		flowchart.WithTitle("mermaid flowchart builder"),
		flowchart.WithOrientalTopToBottom(),
	).
		NodeWithText("A", "Node A").
		StadiumNode("B", "Node B").
		SubroutineNode("C", "Node C").
		DatabaseNode("D", "Database").
		LinkWithArrowHead("A", "B").
		LinkWithArrowHeadAndText("B", "D", "send original data").
		LinkWithArrowHead("B", "C").
		DottedLinkWithText("C", "D", "send filtered data").
		String()

	err = markdown.NewMarkdown(f).
		H2("Flowchart").
		CodeBlocks(markdown.SyntaxHighlightMermaid, fc).
		Build()

	if err != nil {
		panic(err)
	}
}
```

Plain text output: [markdown is here](./doc/flowchart/generated.md)
````
## Flowchart
```mermaid
---
title: mermaid flowchart builder
---
flowchart TB
	A["Node A"]
	B(["Node B"])
	C[["Node C"]]
	D[("Database")]
	A-->B
	B-->|"send original data"|D
	B-->C
	C-. "send filtered data" .-> D
```
````

Mermaid output:
```mermaid
flowchart TB
	A["Node A"]
	B(["Node B"])
	C[["Node C"]]
	D[("Database")]
	A-->B
	B-->|"send original data"|D
	B-->C
	C-. "send filtered data" .-> D
```

### Pie chart syntax

```go
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/piechart"
)

//go:generate go run main.go

func main() {
	f, err := os.Create("generated.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	chart := piechart.NewPieChart(
		io.Discard,
		piechart.WithTitle("mermaid pie chart builder"),
		piechart.WithShowData(true),
	).
		LabelAndIntValue("A", 10).
		LabelAndFloatValue("B", 20.1).
		LabelAndIntValue("C", 30).
		String()

	err = markdown.NewMarkdown(f).
		H2("Pie Chart").
		CodeBlocks(markdown.SyntaxHighlightMermaid, chart).
		Build()

	if err != nil {
		panic(err)
	}
}
```

Plain text output: [markdown is here](./doc/piechart/generated.md)
````
## Pie Chart
```mermaid
%%{init: {"pie": {"textPosition": 0.75}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
pie showData
    title mermaid pie chart builder
    "A" : 10
    "B" : 20.100000
    "C" : 30
```
````

Mermaid output:
```mermaid
%%{init: {"pie": {"textPosition": 0.75}, "themeVariables": {"pieOuterStrokeWidth": "5px"}} }%%
pie showData
    title mermaid pie chart builder
    "A" : 10
    "B" : 20.100000
    "C" : 30
```

### Architecture Diagrams (beta feature)

[The mermaid provides a feature to visualize infrastructure architecture as a beta version](https://mermaid.js.org/syntax/architecture.html), and that feature has been introduced.

```go
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/arch"
)

//go:generate go run main.go

func main() {
	f, err := os.Create("generated.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	diagram := arch.NewArchitecture(io.Discard).
		Service("left_disk", arch.IconDisk, "Disk").
		Service("top_disk", arch.IconDisk, "Disk").
		Service("bottom_disk", arch.IconDisk, "Disk").
		Service("top_gateway", arch.IconInternet, "Gateway").
		Service("bottom_gateway", arch.IconInternet, "Gateway").
		Junction("junctionCenter").
		Junction("junctionRight").
		LF().
		Edges(
			arch.Edge{
				ServiceID: "left_disk",
				Position:  arch.PositionRight,
				Arrow:     arch.ArrowNone,
			},
			arch.Edge{
				ServiceID: "junctionCenter",
				Position:  arch.PositionLeft,
				Arrow:     arch.ArrowNone,
			}).
		Edges(
			arch.Edge{
				ServiceID: "top_disk",
				Position:  arch.PositionBottom,
				Arrow:     arch.ArrowNone,
			},
			arch.Edge{
				ServiceID: "junctionCenter",
				Position:  arch.PositionTop,
				Arrow:     arch.ArrowNone,
			}).
		Edges(
			arch.Edge{
				ServiceID: "bottom_disk",
				Position:  arch.PositionTop,
				Arrow:     arch.ArrowNone,
			},
			arch.Edge{
				ServiceID: "junctionCenter",
				Position:  arch.PositionBottom,
				Arrow:     arch.ArrowNone,
			}).
		Edges(
			arch.Edge{
				ServiceID: "junctionCenter",
				Position:  arch.PositionRight,
				Arrow:     arch.ArrowNone,
			},
			arch.Edge{
				ServiceID: "junctionRight",
				Position:  arch.PositionLeft,
				Arrow:     arch.ArrowNone,
			}).
		Edges(
			arch.Edge{
				ServiceID: "top_gateway",
				Position:  arch.PositionBottom,
				Arrow:     arch.ArrowNone,
			},
			arch.Edge{
				ServiceID: "junctionRight",
				Position:  arch.PositionTop,
				Arrow:     arch.ArrowNone,
			}).
		Edges(
			arch.Edge{
				ServiceID: "bottom_gateway",
				Position:  arch.PositionTop,
				Arrow:     arch.ArrowNone,
			},
			arch.Edge{
				ServiceID: "junctionRight",
				Position:  arch.PositionBottom,
				Arrow:     arch.ArrowNone,
			}).String() //nolint

	err = markdown.NewMarkdown(f).
		H2("Architecture Diagram").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()

	if err != nil {
		panic(err)
	}
```

Plain text output: [markdown is here](./doc/architecture/generated.md)
````
## Architecture Diagram
```mermaid
architecture-beta
    service left_disk(disk)[Disk]
    service top_disk(disk)[Disk]
    service bottom_disk(disk)[Disk]
    service top_gateway(internet)[Gateway]
    service bottom_gateway(internet)[Gateway]
    junction junctionCenter
    junction junctionRight

    left_disk:R -- L:junctionCenter
    top_disk:B -- T:junctionCenter
    bottom_disk:T -- B:junctionCenter
    junctionCenter:R -- L:junctionRight
    top_gateway:B -- T:junctionRight
    bottom_gateway:T -- B:junctionRight
```
````

![Architecture Diagram](./doc/architecture/image.png)

### State Diagram syntax

```go
package main

import (
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/state"
)

//go:generate go run main.go

func main() {
	f, err := os.Create("generated.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	diagram := state.NewDiagram(io.Discard, state.WithTitle("Order State Machine")).
		StartTransition("Pending").
		State("Pending", "Order received").
		State("Processing", "Preparing order").
		State("Shipped", "Order in transit").
		State("Delivered", "Order completed").
		LF().
		TransitionWithNote("Pending", "Processing", "payment confirmed").
		TransitionWithNote("Processing", "Shipped", "items packed").
		TransitionWithNote("Shipped", "Delivered", "customer received").
		LF().
		NoteRight("Pending", "Waiting for payment").
		NoteRight("Processing", "Preparing items").
		LF().
		EndTransition("Delivered").
		String()

	err = markdown.NewMarkdown(f).
		H2("State Diagram").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()

	if err != nil {
		panic(err)
	}
}
```

Plain text output: [markdown is here](./doc/state/generated.md)
````
## State Diagram
```mermaid
---
title: Order State Machine
---
stateDiagram-v2
    [*] --> Pending
    Pending : Order received
    Processing : Preparing order
    Shipped : Order in transit
    Delivered : Order completed

    Pending --> Processing : payment confirmed
    Processing --> Shipped : items packed
    Shipped --> Delivered : customer received

    note right of Pending : Waiting for payment
    note right of Processing : Preparing items

    Delivered --> [*]
```
````

Mermaid output:
```mermaid
---
title: Order State Machine
---
stateDiagram-v2
    [*] --> Pending
    Pending : Order received
    Processing : Preparing order
    Shipped : Order in transit
    Delivered : Order completed

    Pending --> Processing : payment confirmed
    Processing --> Shipped : items packed
    Shipped --> Delivered : customer received

    note right of Pending : Waiting for payment
    note right of Processing : Preparing items

    Delivered --> [*]
```

## Creating an index for a directory full of markdown files
The markdown package can create an index for Markdown files within the specified directory. This feature was added to generate indexes for Markdown documents produced by [nao1215/spectest](https://github.com/nao1215/spectest).
  
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
  
## License
[MIT License](./LICENSE)

## Contribution
First off, thanks for taking the time to contribute! See [CONTRIBUTING.md](./CONTRIBUTING.md) for more information. Contributions are not only related to development. For example, GitHub Star motivates me to develop! Please feel free to contribute to this project.

[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/markdown&type=Date)](https://star-history.com/#nao1215/markdown&Date)

### Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://debimate.jp/"><img src="https://avatars.githubusercontent.com/u/22737008?v=4?s=50" width="50px;" alt="CHIKAMATSU Naohiro"/><br /><sub><b>CHIKAMATSU Naohiro</b></sub></a><br /><a href="https://github.com/nao1215/markdown/commits?author=nao1215" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/varmakarthik12"><img src="https://avatars.githubusercontent.com/u/17958166?v=4?s=50" width="50px;" alt="Karthik Sundari"/><br /><sub><b>Karthik Sundari</b></sub></a><br /><a href="https://github.com/nao1215/markdown/commits?author=varmakarthik12" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Avihuc"><img src="https://avatars.githubusercontent.com/u/32455410?v=4?s=50" width="50px;" alt="Avihuc"/><br /><sub><b>Avihuc</b></sub></a><br /><a href="https://github.com/nao1215/markdown/commits?author=Avihuc" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://www.claranceliberi.me/"><img src="https://avatars.githubusercontent.com/u/60586899?v=4?s=50" width="50px;" alt="Clarance Liberiste Ntwari"/><br /><sub><b>Clarance Liberiste Ntwari</b></sub></a><br /><a href="https://github.com/nao1215/markdown/commits?author=claranceliberi" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/amitaifrey"><img src="https://avatars.githubusercontent.com/u/7527632?v=4?s=50" width="50px;" alt="Amitai Frey"/><br /><sub><b>Amitai Frey</b></sub></a><br /><a href="https://github.com/nao1215/markdown/commits?author=amitaifrey" title="Code">💻</a></td>
    </tr>
  </tbody>
  <tfoot>
    <tr>
      <td align="center" size="13px" colspan="7">
        <img src="https://raw.githubusercontent.com/all-contributors/all-contributors-cli/1b8533af435da9854653492b1327a23a4dbd0a10/assets/logo-small.svg">
          <a href="https://all-contributors.js.org/docs/en/bot/usage">Add your contributions</a>
        </img>
      </td>
    </tr>
  </tfoot>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!
