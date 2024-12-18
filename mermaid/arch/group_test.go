package arch

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
