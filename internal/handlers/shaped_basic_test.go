package handlers

import (
	"testing"

	"github.com/verdoga/dslparser/internal/model"
	"github.com/verdoga/dslparser/internal/registry"
	"github.com/verdoga/dslparser/internal/structure"
)

// TestApplyBuildsShapedGroupsAndExactPlaceholders проверяет группы, порядок и точные Unicode-диапазоны.
func TestApplyBuildsShapedGroupsAndExactPlaceholders(t *testing.T) {
	texts := []string{"@dsl-version 1.1", "@example {", "до _____", "---", "после \\_____ _____", "}", "@wordlist один;два\\;три;_____"}
	lines := make([]model.Line, len(texts))
	for i, s := range texts {
		lines[i] = model.Line{Number: i + 1, Text: s}
	}
	rules, _ := registry.ForVersion("1.1")
	blocks := structure.Analyze(lines, rules)
	roots := structure.Build(lines, blocks)
	Apply(roots, lines, rules)
	if len(roots[0].Elements) != 4 || roots[0].Elements[1].Name != "before" || roots[0].Elements[2].Name != "separator" || roots[0].Elements[3].Name != "after" {
		t.Fatalf("example=%#v", roots[0].Elements)
	}
	before := roots[0].Elements[1].Children[0].Children
	if len(before) != 1 || before[0].Span.StartColumn != 4 || before[0].Span.EndColumn != 9 {
		t.Fatalf("before placeholders=%#v", before)
	}
	after := roots[0].Elements[3].Children[0].Children
	if len(after) != 1 || after[0].Span.StartColumn != 14 {
		t.Fatalf("after placeholders=%#v", after)
	}
	items := roots[1].Elements[len(roots[1].Elements)-1]
	if len(items.Children) != 3 || items.Children[1].Raw != `два\;три` || len(items.Children[2].Children) != 1 {
		t.Fatalf("wordlist=%#v", items)
	}
}

// TestApplyRawBodiesPreserveOnlyConfiguredEmptyLines проверяет непрозрачность и правила пустых строк.
func TestApplyRawBodiesPreserveOnlyConfiguredEmptyLines(t *testing.T) {
	for _, tt := range []struct {
		name string
		want int
	}{{"text", 4}, {"table", 3}} {
		t.Run(tt.name, func(t *testing.T) {
			texts := []string{"@dsl-version 1.1", "@" + tt.name + " {", "@nested", "", "---", "}"}
			lines := make([]model.Line, len(texts))
			for i, s := range texts {
				lines[i] = model.Line{Number: i + 1, Text: s}
			}
			rules, _ := registry.ForVersion("1.1")
			b := structure.Analyze(lines, rules)
			roots := structure.Build(lines, b)
			Apply(roots, lines, rules)
			if len(roots[0].Elements) != tt.want {
				t.Fatalf("elements=%#v", roots[0].Elements)
			}
		})
	}
}
