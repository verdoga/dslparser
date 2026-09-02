package handlers

import (
	"testing"

	"github.com/verdoga/dslparser/internal/model"
)

// TestSeparatorBodies проверяет все формы тел с локальным разделителем.
func TestSeparatorBodies(t *testing.T) {
	tests := []struct {
		name, tag string
		body      []string
		groups    []string
		code      string
	}{
		{name: "matching", tag: "matching", body: []string{"one", "---", "two"}, groups: []string{"upper", "separator", "lower"}},
		{name: "matching missing", tag: "matching", body: []string{"one"}, groups: []string{"body"}, code: "MISSING_SEPARATOR"},
		{name: "matching repeated", tag: "matching", body: []string{"---", "---"}, groups: []string{"body"}, code: "MULTIPLE_SEPARATORS"},
		{name: "ordering short", tag: "ordering", body: []string{"one", "two"}, groups: []string{"items"}},
		{name: "ordering extended", tag: "ordering", body: []string{"one", "---", "first"}, groups: []string{"items", "separator", "positions"}},
		{name: "ordering repeated", tag: "ordering", body: []string{"---", "---"}, groups: []string{"body"}, code: "MULTIPLE_SEPARATORS"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lines := []model.Line{{Number: 1, Text: "@" + tc.tag + " {"}}
			for i, text := range tc.body {
				lines = append(lines, model.Line{Number: i + 2, Text: text})
			}
			n := &model.Node{Parsed: true, Span: model.Span{StartLine: 1, StartColumn: 1, EndLine: 1, EndColumn: 2}, Block: &model.Block{OpenLine: 0, CloseLine: len(lines), HasClose: true}}
			if tc.tag == "matching" {
				applyMatching(n, lines)
			} else {
				applyOrdering(n, lines)
			}
			if len(n.Elements) != len(tc.groups) {
				t.Fatalf("elements = %d, want %d", len(n.Elements), len(tc.groups))
			}
			for i, name := range tc.groups {
				if n.Elements[i].Name != name {
					t.Errorf("element %d = %q, want %q", i, n.Elements[i].Name, name)
				}
			}
			if tc.code != "" && (len(n.Diagnostics) != 1 || n.Diagnostics[0].Code != tc.code || n.Parsed) {
				t.Errorf("diagnostics = %+v, parsed = %v", n.Diagnostics, n.Parsed)
			}
		})
	}
}
