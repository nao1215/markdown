//go:build linux || darwin

package arch_test

import (
	"fmt"
	"io"
	"os"

	"github.com/nao1215/markdown"
	"github.com/nao1215/markdown/mermaid/arch"
)

// ExampleArchitecture skips this test on Windows.
// The newline codes in the comment section where
// the expected values are written are represented as '\n',
// causing failures when testing on Windows.
func ExampleArchitecture() {
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

	_ = markdown.NewMarkdown(os.Stdout).
		H2("Architecture Diagram").
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ## Architecture Diagram
	// ```mermaid
	// architecture-beta
	//     service left_disk(disk)[Disk]
	//     service top_disk(disk)[Disk]
	//     service bottom_disk(disk)[Disk]
	//     service top_gateway(internet)[Gateway]
	//     service bottom_gateway(internet)[Gateway]
	//     junction junctionCenter
	//     junction junctionRight
	//
	//     left_disk:R -- L:junctionCenter
	//     top_disk:B -- T:junctionCenter
	//     bottom_disk:T -- B:junctionCenter
	//     junctionCenter:R -- L:junctionRight
	//     top_gateway:B -- T:junctionRight
	//     bottom_gateway:T -- B:junctionRight
	// ```
}

// ExampleNewArchitecture shows the shape every architecture diagram has: a
// writer, a chain of calls, and Build.
//
// A title here takes only letters, digits, underscores and spaces. mermaid's
// architecture-beta grammar accepts nothing else and refuses the whole diagram
// when it finds it, and there is no escape to reach for; the package
// documentation says so at more length.
func ExampleNewArchitecture() {
	_ = arch.NewArchitecture(os.Stdout).
		Service("api", arch.IconServer, "API").
		Build()

	// Output:
	// architecture-beta
	//     service api(server)[API]
}

// ExampleArchitecture_Service adds one piece of the system, with the icon it is
// drawn as.
func ExampleArchitecture_Service() {
	_ = arch.NewArchitecture(os.Stdout).
		Service("api", arch.IconServer, "API").
		Service("db", arch.IconDatabase, "Database").
		Build()

	// Output:
	// architecture-beta
	//     service api(server)[API]
	//     service db(database)[Database]
}

// ExampleArchitecture_ServiceInGroup puts a service inside a group, which is
// how a diagram says which things sit together.
func ExampleArchitecture_ServiceInGroup() {
	_ = arch.NewArchitecture(os.Stdout).
		Group("cloud", arch.IconCloud, "Cloud").
		ServiceInGroup("api", arch.IconServer, "API", "cloud").
		Build()

	// Output:
	// architecture-beta
	//     group cloud(cloud)[Cloud]
	//     service api(server)[API] in cloud
}

// ExampleArchitecture_Group adds a box the services can sit in.
func ExampleArchitecture_Group() {
	_ = arch.NewArchitecture(os.Stdout).
		Group("cloud", arch.IconCloud, "Cloud").
		Group("onprem", arch.IconServer, "On premise").
		Build()

	// Output:
	// architecture-beta
	//     group cloud(cloud)[Cloud]
	//     group onprem(server)[On premise]
}

// ExampleArchitecture_GroupInParentGroup nests one group inside another, for a
// region holding a cluster holding services.
func ExampleArchitecture_GroupInParentGroup() {
	_ = arch.NewArchitecture(os.Stdout).
		Group("cloud", arch.IconCloud, "Cloud").
		GroupInParentGroup("cluster", arch.IconServer, "Cluster", "cloud").
		ServiceInGroup("api", arch.IconServer, "API", "cluster").
		Build()

	// Output:
	// architecture-beta
	//     group cloud(cloud)[Cloud]
	//     group cluster(server)[Cluster] in cloud
	//     service api(server)[API] in cluster
}

// ExampleArchitecture_Edges joins two services. Each end says which side of its
// service the line leaves from and whether it carries an arrowhead.
func ExampleArchitecture_Edges() {
	_ = arch.NewArchitecture(os.Stdout).
		Service("api", arch.IconServer, "API").
		Service("db", arch.IconDatabase, "Database").
		Edges(
			arch.Edge{ServiceID: "api", Position: arch.PositionRight, Arrow: arch.ArrowNone},
			arch.Edge{ServiceID: "db", Position: arch.PositionLeft, Arrow: arch.ArrowRight},
		).
		Build()

	// Output:
	// architecture-beta
	//     service api(server)[API]
	//     service db(database)[Database]
	//     api:R --> L:db
}

// ExampleArchitecture_EdgesInAnothorGroup joins two services that are in
// different groups, which mermaid draws differently from an edge inside one.
func ExampleArchitecture_EdgesInAnothorGroup() {
	_ = arch.NewArchitecture(os.Stdout).
		Group("cloud", arch.IconCloud, "Cloud").
		Group("onprem", arch.IconServer, "On premise").
		ServiceInGroup("api", arch.IconServer, "API", "cloud").
		ServiceInGroup("db", arch.IconDatabase, "Database", "onprem").
		EdgesInAnothorGroup(
			arch.Edge{ServiceID: "api", Position: arch.PositionRight, Arrow: arch.ArrowNone},
			arch.Edge{ServiceID: "db", Position: arch.PositionLeft, Arrow: arch.ArrowRight},
		).
		Build()

	// Output:
	// architecture-beta
	//     group cloud(cloud)[Cloud]
	//     group onprem(server)[On premise]
	//     service api(server)[API] in cloud
	//     service db(database)[Database] in onprem
	//     api{group}:R --> L:db{group}
}

// ExampleArchitecture_Junction adds a point several edges can meet at, so a
// diagram can route lines around each other rather than across.
func ExampleArchitecture_Junction() {
	_ = arch.NewArchitecture(os.Stdout).
		Service("api", arch.IconServer, "API").
		Junction("j1").
		Build()

	// Output:
	// architecture-beta
	//     service api(server)[API]
	//     junction j1
}

// ExampleArchitecture_JunctionsInParent puts a junction inside a group.
func ExampleArchitecture_JunctionsInParent() {
	_ = arch.NewArchitecture(os.Stdout).
		Group("cloud", arch.IconCloud, "Cloud").
		JunctionsInParent("j1", "cloud").
		Build()

	// Output:
	// architecture-beta
	//     group cloud(cloud)[Cloud]
	//     junction j1 in cloud
}

// ExampleArchitecture_String returns the diagram without needing a writer,
// which is how it is handed to a markdown code block.
func ExampleArchitecture_String() {
	diagram := arch.NewArchitecture(io.Discard).
		Service("api", arch.IconServer, "API").
		String()

	_ = markdown.NewMarkdown(os.Stdout).
		CodeBlocks(markdown.SyntaxHighlightMermaid, diagram).
		Build()

	// Output:
	// ```mermaid
	// architecture-beta
	//     service api(server)[API]
	// ```
}

// ExampleArchitecture_Build writes the diagram and reports the error the chain
// recorded.
func ExampleArchitecture_Build() {
	err := arch.NewArchitecture(nil).Service("api", arch.IconServer, "API").Build()
	fmt.Println("error:", err)

	// Output:
	// error: output writer must not be nil
}

// ExampleArchitecture_Error reports the same error Build does, for code that
// wants to look before writing anything.
func ExampleArchitecture_Error() {
	a := arch.NewArchitecture(io.Discard).Service("api", arch.IconServer, "API")
	fmt.Println("error:", a.Error())

	// Output:
	// error: <nil>
}

// ExampleArchitecture_LF adds a blank line to the diagram body.
func ExampleArchitecture_LF() {
	_ = arch.NewArchitecture(os.Stdout).
		Service("api", arch.IconServer, "API").
		LF().
		Service("db", arch.IconDatabase, "Database").
		Build()

	// Output:
	// architecture-beta
	//     service api(server)[API]
	//
	//     service db(database)[Database]
}

// ExampleIcon shows the icons a service or a group can be drawn as.
func ExampleIcon() {
	_ = arch.NewArchitecture(os.Stdout).
		Service("a", arch.IconCloud, "Cloud").
		Service("b", arch.IconDatabase, "Database").
		Service("c", arch.IconDisk, "Disk").
		Service("d", arch.IconInternet, "Internet").
		Service("e", arch.IconServer, "Server").
		Build()

	// Output:
	// architecture-beta
	//     service a(cloud)[Cloud]
	//     service b(database)[Database]
	//     service c(disk)[Disk]
	//     service d(internet)[Internet]
	//     service e(server)[Server]
}

// ExampleEdge shows what one end of a line is: the service it leaves, the side
// it leaves from, and whether it carries an arrowhead.
func ExampleEdge() {
	from := arch.Edge{ServiceID: "api", Position: arch.PositionRight, Arrow: arch.ArrowNone}
	to := arch.Edge{ServiceID: "db", Position: arch.PositionLeft, Arrow: arch.ArrowRight}

	_ = arch.NewArchitecture(os.Stdout).
		Service("api", arch.IconServer, "API").
		Service("db", arch.IconDatabase, "Database").
		Edges(from, to).
		Build()

	// Output:
	// architecture-beta
	//     service api(server)[API]
	//     service db(database)[Database]
	//     api:R --> L:db
}

// ExamplePosition shows the four sides of a service a line can leave from.
func ExamplePosition() {
	for _, position := range []arch.Position{
		arch.PositionTop, arch.PositionBottom, arch.PositionLeft, arch.PositionRight,
	} {
		_ = arch.NewArchitecture(os.Stdout).
			Service("api", arch.IconServer, "API").
			Service("db", arch.IconDatabase, "Database").
			Edges(
				arch.Edge{ServiceID: "api", Position: position, Arrow: arch.ArrowNone},
				arch.Edge{ServiceID: "db", Position: arch.PositionLeft, Arrow: arch.ArrowNone},
			).
			Build()
		fmt.Println()
	}

	// Output:
	// architecture-beta
	//     service api(server)[API]
	//     service db(database)[Database]
	//     api:T -- L:db
	// architecture-beta
	//     service api(server)[API]
	//     service db(database)[Database]
	//     api:B -- L:db
	// architecture-beta
	//     service api(server)[API]
	//     service db(database)[Database]
	//     api:L -- L:db
	// architecture-beta
	//     service api(server)[API]
	//     service db(database)[Database]
	//     api:R -- L:db
}

// ExampleArrow shows what an end of a line can carry.
func ExampleArrow() {
	_ = arch.NewArchitecture(os.Stdout).
		Service("api", arch.IconServer, "API").
		Service("db", arch.IconDatabase, "Database").
		Edges(
			arch.Edge{ServiceID: "api", Position: arch.PositionRight, Arrow: arch.ArrowLeft},
			arch.Edge{ServiceID: "db", Position: arch.PositionLeft, Arrow: arch.ArrowRight},
		).
		Build()

	// Output:
	// architecture-beta
	//     service api(server)[API]
	//     service db(database)[Database]
	//     api:R <--> L:db
}

// ExampleOption shows what an Option is: a function that changes how the
// diagram is written, passed to NewArchitecture.
func ExampleOption() {
	options := []arch.Option{}

	_ = arch.NewArchitecture(os.Stdout, options...).
		Service("api", arch.IconServer, "API").
		Build()

	// Output:
	// architecture-beta
	//     service api(server)[API]
}
