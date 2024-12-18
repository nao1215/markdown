// Package arch is mermaid architecture diagram builder.
// The building blocks of an architecture are groups, services, edges, and junctions.
// The arch package incorporates beta features of Mermaid, so the specifications are subject to significant changes.
package arch

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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
