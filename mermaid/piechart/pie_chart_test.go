// Package piechart is mermaid pie chart builder.
package piechart

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPieChart_Build(t *testing.T) {
	t.Parallel()

	t.Run("Build a pie chart with title", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		p := NewPieChart(
			b,
			WithTitle("mermaid pie chart builder"),
			WithShowData(true),
		)
		p.LabelAndIntValue("A", 10)
		p.LabelAndFloatValue("B", 20.1)
		p.LabelAndIntValue("C", 30)

		if err := p.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := fmt.Sprintf(
			"%s\n%s",
			"%%{init: {\"pie\": {\"textPosition\": 0.75}, \"themeVariables\": {\"pieOuterStrokeWidth\": \"5px\"}} }%%",
			"pie showData\n    title mermaid pie chart builder\n    \"A\" : 10\n    \"B\" : 20.100000\n    \"C\" : 30")
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a pie chart with no title", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		p := NewPieChart(
			b,
			WithShowData(true),
			WithTextPosition(0.5),
		)
		p.LabelAndIntValue("A", 10)
		p.LabelAndFloatValue("B", 20.1)
		p.LabelAndIntValue("C", 30)

		if err := p.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := fmt.Sprintf(
			"%s\n%s",
			"%%{init: {\"pie\": {\"textPosition\": 0.50}, \"themeVariables\": {\"pieOuterStrokeWidth\": \"5px\"}} }%%",
			"pie showData\n    \"A\" : 10\n    \"B\" : 20.100000\n    \"C\" : 30",
		)
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a pie chart with bad text position value(less than 0)", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		p := NewPieChart(
			b,
			WithTitle("mermaid pie chart builder"),
			WithShowData(true),
			WithTextPosition(-0.1),
		)
		p.LabelAndIntValue("A", 10)
		p.LabelAndFloatValue("B", 20.1)
		p.LabelAndIntValue("C", 30)

		if err := p.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := fmt.Sprintf(
			"%s\n%s",
			"%%{init: {\"pie\": {\"textPosition\": 0.75}, \"themeVariables\": {\"pieOuterStrokeWidth\": \"5px\"}} }%%",
			"pie showData\n    title mermaid pie chart builder\n    \"A\" : 10\n    \"B\" : 20.100000\n    \"C\" : 30")
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})

	t.Run("Build a pie chart with bad text position value(more than 1)", func(t *testing.T) {
		t.Parallel()

		b := new(bytes.Buffer)

		p := NewPieChart(
			b,
			WithTitle("mermaid pie chart builder"),
			WithShowData(true),
			WithTextPosition(1.1),
		)
		p.LabelAndIntValue("A", 10)
		p.LabelAndFloatValue("B", 20.1)
		p.LabelAndIntValue("C", 30)

		if err := p.Build(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		want := fmt.Sprintf(
			"%s\n%s",
			"%%{init: {\"pie\": {\"textPosition\": 0.75}, \"themeVariables\": {\"pieOuterStrokeWidth\": \"5px\"}} }%%",
			"pie showData\n    title mermaid pie chart builder\n    \"A\" : 10\n    \"B\" : 20.100000\n    \"C\" : 30")
		got := strings.ReplaceAll(b.String(), "\r\n", "\n")

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("value is mismatch (-want +got):%s", diff)
		}
	})
}
