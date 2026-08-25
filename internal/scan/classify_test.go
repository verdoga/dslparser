package scan

import "testing"

// TestClassifyStructuralKinds проверяет наблюдаемый контракт.
func TestClassifyStructuralKinds(t *testing.T) {
	tests := []struct {
		line string
		want Kind
	}{{"", Empty}, {"@dsl-version 1.1", Version}, {"# title", Heading}, {"@tag x", Tag}, {"@tag x {", BlockOpen}, {"}", BlockClose}, {"} extra", InvalidBoundary}, {"plain", Text}, {`\@tag`, Text}, {`\# title`, Text}, {`@tag \{`, Tag}}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := Classify(tt.line, Structural); got.Kind != tt.want {
				t.Fatalf("kind=%v candidate=%#v", got.Kind, got)
			}
		})
	}
}

// TestClassifyTagPreservesNamesAndFullTail проверяет наблюдаемый контракт.
func TestClassifyTagPreservesNamesAndFullTail(t *testing.T) {
	c := Classify("@MeDiA audio path/file {", Structural)
	if c.Name != "MeDiA" || c.CanonicalName != "media" || c.Tail != "audio path/file {" || c.NameStart != 2 || c.NameEnd != 7 {
		t.Fatalf("candidate=%#v", c)
	}
}

// TestOpaqueModeRecognizesOnlyExactUnescapedClose проверяет наблюдаемый контракт.
func TestOpaqueModeRecognizesOnlyExactUnescapedClose(t *testing.T) {
	for _, line := range []string{"@tag", "# heading", "} extra", `\}`, "---", "{"} {
		if got := Classify(line, Opaque).Kind; got != OpaqueLine {
			t.Errorf("%q kind=%v", line, got)
		}
	}
	if got := Classify("}", Opaque).Kind; got != BlockClose {
		t.Fatalf("close kind=%v", got)
	}
}
