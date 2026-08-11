package gantt_test

import (
	"strings"
	"testing"

	"github.com/nao1215/markdown/mermaid/gantt"
)

// TestChartError covers the accessor callers use to find out why a chart came
// out wrong, which nothing exercised.
func TestChartError(t *testing.T) {
	t.Parallel()

	t.Run("a well formed chart reports no error", func(t *testing.T) {
		t.Parallel()

		c := gantt.NewChart(nil).Section("build").Task("compile", "2024-01-01", "1d")
		if err := c.Error(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("a nil writer is reported through Error after Build", func(t *testing.T) {
		t.Parallel()

		c := gantt.NewChart(nil)
		if err := c.Build(); err == nil {
			t.Fatal("Build with a nil writer must fail")
		}
		if c.Error() == nil {
			t.Error("Error must report the failure Build returned")
		}
		if !strings.Contains(c.Error().Error(), "nil") {
			t.Errorf("unexpected error: %v", c.Error())
		}
	})
}
