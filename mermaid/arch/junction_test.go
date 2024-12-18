package arch

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
