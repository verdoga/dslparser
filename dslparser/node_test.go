package dslparser

import "testing"

// TestNodeExposesUniversalDataWithoutMutation проверяет наблюдаемый контракт.
func TestNodeExposesUniversalDataWithoutMutation(t *testing.T) {
	b := BlockInfo{mode: BodyStructural, hasClose: true}
	n := Node{kind: NodeTag, raw: "@Tag x", value: "x", originalName: "Tag", canonicalName: "tag", headingLevel: 0, parsed: true, synthetic: false, elements: []Element{{raw: "x"}}, children: []Node{{raw: "child"}}, block: &b}
	if n.Kind() != NodeTag || n.Raw() != "@Tag x" || n.Value() != "x" || n.OriginalName() != "Tag" || n.CanonicalName() != "tag" || !n.Parsed() || n.Synthetic() {
		t.Fatal("node lost data")
	}
	e := n.Elements()
	c := n.Children()
	e[0].raw = "bad"
	c[0].raw = "bad"
	if n.Elements()[0].Raw() != "x" || n.Children()[0].Raw() != "child" {
		t.Fatal("node exposed slices")
	}
	block, ok := n.Block()
	if !ok || !block.Closed() || block.Mode() != BodyStructural {
		t.Fatal("node lost block")
	}
}

// TestNodeKindsAreDistinct проверяет наблюдаемый контракт.
func TestNodeKindsAreDistinct(t *testing.T) {
	seen := map[NodeKind]bool{}
	for _, kind := range []NodeKind{NodeStep, NodeHeading, NodeTag, NodeText, NodeBlockBoundary} {
		if seen[kind] {
			t.Fatal("duplicate node kind")
		}
		seen[kind] = true
	}
}
