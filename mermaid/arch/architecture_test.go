// Package arch is mermaid architecture diagram builder.
// The building blocks of an architecture are groups, services, edges, and junctions.
// The arch package incorporates beta features of Mermaid, so the specifications are subject to significant changes.
package arch

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nao1215/markdown/internal/buildertest"
	"github.com/nao1215/markdown/internal/golden"
)

func TestArchitecture_Build(t *testing.T) {
	t.Parallel()

	t.Run("build architecture-beta sample code", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)
		a := NewArchitecture(b)
		a.Service("left_disk", IconDisk, "Disk").
			Service("top_disk", IconDisk, "Disk").
			Service("bottom_disk", IconDisk, "Disk").
			Service("top_gateway", IconInternet, "Gateway").
			Service("bottom_gateway", IconInternet, "Gateway").
			Junction("junctionCenter").
			Junction("junctionRight").
			LF().
			Edges(
				Edge{
					ServiceID: "left_disk",
					Position:  PositionRight,
					Arrow:     ArrowNone,
				},
				Edge{
					ServiceID: "junctionCenter",
					Position:  PositionLeft,
					Arrow:     ArrowNone,
				}).
			Edges(
				Edge{
					ServiceID: "top_disk",
					Position:  PositionBottom,
					Arrow:     ArrowNone,
				},
				Edge{
					ServiceID: "junctionCenter",
					Position:  PositionTop,
					Arrow:     ArrowNone,
				}).
			Edges(
				Edge{
					ServiceID: "bottom_disk",
					Position:  PositionTop,
					Arrow:     ArrowNone,
				},
				Edge{
					ServiceID: "junctionCenter",
					Position:  PositionBottom,
					Arrow:     ArrowNone,
				}).
			Edges(
				Edge{
					ServiceID: "junctionCenter",
					Position:  PositionRight,
					Arrow:     ArrowNone,
				},
				Edge{
					ServiceID: "junctionRight",
					Position:  PositionLeft,
					Arrow:     ArrowNone,
				}).
			Edges(
				Edge{
					ServiceID: "top_gateway",
					Position:  PositionBottom,
					Arrow:     ArrowNone,
				},
				Edge{
					ServiceID: "junctionRight",
					Position:  PositionTop,
					Arrow:     ArrowNone,
				}).
			Edges(
				Edge{
					ServiceID: "bottom_gateway",
					Position:  PositionTop,
					Arrow:     ArrowNone,
				},
				Edge{
					ServiceID: "junctionRight",
					Position:  PositionBottom,
					Arrow:     ArrowNone,
				})
		if err := a.Build(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if err := a.Error(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := `architecture-beta
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
    bottom_gateway:T -- B:junctionRight`

		want = strings.ReplaceAll(want, "\r\n", "\n")
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("build architecture-beta with Group", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)
		a := NewArchitecture(b)
		a.Group("group1", IconCloud, "Group1").
			Group("group2", IconCloud, "Group2").
			ServiceInGroup("left_disk", IconDisk, "Disk", "group1").
			ServiceInGroup("right_disk", IconDisk, "Disk", "group2").
			EdgesInAnothorGroup(
				Edge{
					ServiceID: "left_disk",
					Position:  PositionRight,
					Arrow:     ArrowNone,
				},
				Edge{
					ServiceID: "right_disk",
					Position:  PositionLeft,
					Arrow:     ArrowNone,
				})
		if err := a.Build(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if err := a.Error(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := `architecture-beta
    group group1(cloud)[Group1]
    group group2(cloud)[Group2]
    service left_disk(disk)[Disk] in group1
    service right_disk(disk)[Disk] in group2
    left_disk{group}:R -- L:right_disk{group}`

		want = strings.ReplaceAll(want, "\r\n", "\n")
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

// TestBuildContract asserts the error handling every builder in this module
// shares. The contract itself lives in internal/buildertest.
func TestBuildContract(t *testing.T) {
	t.Parallel()

	buildertest.RunBuildContract(t, func(w io.Writer) buildertest.Builder {
		return NewArchitecture(w).Service("api", IconServer, "API")
	})
}

// TestGoldenArchitecture pins the rendered diagram of every builder method of
// this package.
func TestGoldenArchitecture(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	err := NewArchitecture(buf).
		Group("api", IconCloud, "API").
		GroupInParentGroup("storage", IconDatabase, "Storage", "api").
		Service("gateway", IconInternet, "Gateway").
		ServiceInGroup("db", IconDatabase, "Database", "storage").
		ServiceInGroup("disk", IconDisk, "Disk", "storage").
		ServiceInGroup("worker", IconServer, "Worker", "api").
		Junction("hub").
		JunctionsInParent("inner", "storage").
		LF().
		Edges(
			Edge{ServiceID: "gateway", Position: PositionRight, Arrow: ArrowRight},
			Edge{ServiceID: "hub", Position: PositionLeft, Arrow: ArrowNone},
		).
		Edges(
			Edge{ServiceID: "hub", Position: PositionBottom, Arrow: ArrowNone},
			Edge{ServiceID: "worker", Position: PositionTop, Arrow: ArrowLeft},
		).
		EdgesInAnothorGroup(
			Edge{ServiceID: "worker", Position: PositionRight, Arrow: ArrowNone},
			Edge{ServiceID: "db", Position: PositionLeft, Arrow: ArrowRight},
		).
		Build()
	if err != nil {
		t.Fatalf("Build() = %v, want nil", err)
	}

	if err := golden.Assert("architecture.md", buf.String()); err != nil {
		t.Error(err)
	}
}

func TestArchitecture_GroupInParentGroup(t *testing.T) {
	t.Parallel()

	t.Run("set group in parent group", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)
		a := NewArchitecture(b)

		if err := a.GroupInParentGroup("group1", "icon", "title", "parentGroup").Build(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := `architecture-beta
    group group1(icon)[title] in parentGroup`
		want = strings.ReplaceAll(want, "\r\n", "\n")
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

func TestArchitecture_JunctionsInParent(t *testing.T) {
	t.Parallel()

	t.Run("set junctions in parent", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)
		a := NewArchitecture(b)

		if err := a.JunctionsInParent("junction1", "parentGroup").Build(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		want := `architecture-beta
    junction junction1 in parentGroup`
		want = strings.ReplaceAll(want, "\r\n", "\n")
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}

// TestBuildWithNilWriter covers the case where a architecture is built for String()
// only and Build() is called by mistake. Build() used to dereference the nil
// writer and take the process down; it has to return an error instead.
func TestBuildWithNilWriter(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Build() panicked with a nil writer: %v", r)
		}
	}()

	d := NewArchitecture(nil)

	// String() has always worked without a writer, and callers rely on it.
	_ = d.String()

	err := d.Build()
	if err == nil {
		t.Fatal("Build() with a nil writer must return an error")
	}
	if err.Error() != "output writer must not be nil" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// errWrite is the failure the writer below reports, so the test can assert that
// Build passed it through rather than inventing an error of its own.
var errWrite = errors.New("write failed")

// errWriter fails every write, which is what a full disk or a closed pipe looks
// like to Build.
type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errWrite
}

// TestBuildReportsWriteFailure covers the branch where the destination accepts
// the diagram and then fails. Silently returning nil there would hand the caller
// a document that was never written.
func TestBuildReportsWriteFailure(t *testing.T) {
	t.Parallel()

	err := NewArchitecture(errWriter{}).Build()
	if err == nil {
		t.Fatal("Build must report a failing writer")
	}
	if !errors.Is(err, errWrite) {
		t.Errorf("Build lost the destination error: %v", err)
	}
}

// TestTitleIsPassedThroughUnchanged pins the decision this package's
// documentation records: mermaid's architecture-beta grammar accepts
// only [A-Za-z0-9_ ] in a title and refuses even its own "#name;" escape there,
// so there is nothing to encode to.
//
// A title outside that set is passed through as it was given rather than
// mangled into something that renders but says something else, and rather than
// rejected, because a label mermaid cannot take is not a caller error. What
// this test guards is that no later sweep quietly starts escaping here: the
// escape would not work, and the output would change for nothing.
func TestTitleIsPassedThroughUnchanged(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"the set mermaid accepts": "plain_label 1",
		"a quotation mark":        `the "core"`,
		"a hyphen":                "Order-Service",
		"Japanese":                "注文サービス",
		"an entity escape":        "a#quot;b",
	}

	for name, title := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := NewArchitecture(io.Discard).Service("api", IconServer, title).String()
			if want := "    service api(server)[" + title + "]"; !strings.Contains(got, want) {
				t.Errorf("diagram =\n%s\nwant it to contain\n%s", got, want)
			}
		})
	}
}
