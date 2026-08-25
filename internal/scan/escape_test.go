package scan

import "testing"

// TestEscapedUsesBackslashParity проверяет наблюдаемый контракт.
func TestEscapedUsesBackslashParity(t *testing.T) {
	for _, tt := range []struct {
		s     string
		index int
		want  bool
	}{{"@", 0, false}, {`\@`, 1, true}, {`\\@`, 2, false}, {`\\\@`, 3, true}} {
		if got := Escaped(tt.s, tt.index); got != tt.want {
			t.Errorf("Escaped(%q)=%v", tt.s, got)
		}
	}
}

// TestDecodeEscapesRemovesOneOnlyForSpecialCharacters проверяет наблюдаемый контракт.
func TestDecodeEscapesRemovesOneOnlyForSpecialCharacters(t *testing.T) {
	if got := DecodeEscapes(`\@ \# \{ \} \--- \x \\`, `@#{}-\`); got != `@ # { } --- \x \` {
		t.Fatalf("decoded = %q", got)
	}
}
