package structure

import (
	"testing"
	"unicode/utf8"

	"github.com/verdoga/dslparser/internal/model"
	"github.com/verdoga/dslparser/internal/registry"
)

// FuzzAnalyzeBlocks проверяет устойчивость поиска фигурных границ и корректность их диапазонов.
func FuzzAnalyzeBlocks(f *testing.F) {
	for _, seed := range []string{"", "@text {\nbody\n}", "@variants {\n@variant a\n}", "}\n@text {"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, source string) {
		if !utf8.ValidString(source) {
			t.Skip()
		}
		rules, ok := registry.ForVersion("1.1")
		if !ok {
			t.Fatal("registry 1.1 is unavailable")
		}
		lines := fuzzLines(source)
		result := Analyze(lines, rules)
		for open, block := range result.Blocks {
			if open < 0 || open >= len(lines) || block.OpenLine != open || block.CloseLine >= len(lines) || block.HasClose && block.CloseLine <= open {
				t.Fatalf("invalid block at %d: %+v for %d lines", open, block, len(lines))
			}
		}
	})
}

// fuzzLines разбивает строку на физические строки без изменения их содержимого.
func fuzzLines(source string) []model.Line {
	var lines []model.Line
	start := 0
	for i := 0; i <= len(source); i++ {
		if i == len(source) || source[i] == '\n' {
			lines = append(lines, model.Line{Number: len(lines) + 1, Text: source[start:i]})
			start = i + 1
		}
	}
	return lines
}
