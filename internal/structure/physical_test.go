package structure

import (
	"testing"

	"github.com/verdoga/dslparser/internal/registry"
)

// TestBuildPreservesPhysicalOrderWithoutBoundaryNodes проверяет вложенность и отсутствие дублирования границ.
func TestBuildPreservesPhysicalOrderWithoutBoundaryNodes(t *testing.T) {
	lines := testLines("@dsl-version 1.1", "@variants {", "# title", "plain", "}", "@Unknown one")
	rules, _ := registry.ForVersion("1.1")
	blocks := Analyze(lines, rules)
	roots := Build(lines, blocks)
	if len(roots) != 2 || len(roots[0].Children) != 2 || roots[0].Children[0].Kind != 2 || roots[0].Children[1].Raw != "plain" {
		t.Fatalf("roots=%#v", roots)
	}
	if roots[1].OriginalName != "Unknown" || roots[1].CanonicalName != "unknown" || roots[1].Tail != "one" {
		t.Fatalf("unknown=%#v", roots[1])
	}
}
