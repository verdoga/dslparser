package structure

import (
	"testing"

	"github.com/verdoga/dslparser/internal/model"
	"github.com/verdoga/dslparser/internal/registry"
)

// TestAnalyzeDistinguishesStructuralAndOpaqueBlocks проверяет режимы контейнеров и единую карту границ.
func TestAnalyzeDistinguishesStructuralAndOpaqueBlocks(t *testing.T) {
	lines := testLines("@dsl-version 1.1", "@variants {", "@text {", "@inside {", "}", "}")
	rules, _ := registry.ForVersion("1.1")
	result := Analyze(lines, rules)
	outer, ok := result.Blocks[1]
	if !ok || !outer.Structural || !outer.HasClose || outer.CloseLine != 5 {
		t.Fatalf("outer=%#v", outer)
	}
	inner, ok := result.Blocks[2]
	if !ok || inner.Structural || !inner.HasClose || inner.CloseLine != 4 {
		t.Fatalf("inner=%#v", inner)
	}
}

// TestAnalyzeRecoversFromEveryInvalidBoundary проверяет независимые локальные ошибки и незакрытый блок.
func TestAnalyzeRecoversFromEveryInvalidBoundary(t *testing.T) {
	lines := testLines("@dsl-version 1.1", "{", "}", "} хвост", "@choice { tail", "@choice {")
	rules, _ := registry.ForVersion("1.1")
	result := Analyze(lines, rules)
	want := map[int]string{1: "INVALID_BLOCK_OPEN", 2: "ORPHAN_BLOCK_CLOSE", 3: "INVALID_BLOCK_CLOSE", 4: "INVALID_BLOCK_OPEN", 5: "UNCLOSED_BLOCK"}
	for line, code := range want {
		if len(result.Diagnostics[line]) == 0 || result.Diagnostics[line][0].Code != code {
			t.Errorf("line %d diagnostics=%#v", line, result.Diagnostics[line])
		}
	}
}

// testLines создаёт последовательность физических строк для структурных тестов.
func testLines(texts ...string) []model.Line {
	lines := make([]model.Line, len(texts))
	for i, text := range texts {
		lines[i] = model.Line{Number: i + 1, Text: text}
	}
	return lines
}
