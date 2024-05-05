package sequence

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewConfig(t *testing.T) {
	t.Parallel()

	t.Run("should return a new config", func(t *testing.T) {
		t.Parallel()

		got := NewConfig()
		want := &Config{
			MirrorActors:            false,
			BottomMariginAdjustment: 1,
			ActorFontSize:           14,
			ActorFontFamily:         "Open Sans, sans-serif",
			ActorFontWeight:         "Open Sans, sans-serif",
			NoteFontSize:            14,
			NoteFontFamily:          "trebuchet ms, verdana, arial",
			NoteFontWeight:          "trebuchet ms, verdana, arial",
			NoteAlign:               "center",
			MessageFontSize:         16,
			MessageFontFamily:       "trebuchet ms, verdana, arial",
			MessageFontWeight:       "trebuchet ms, verdana, arial",
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}
