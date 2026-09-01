package normalize

import "testing"

// TestInputRejectsInvalidUTF8 проверяет наблюдаемый контракт.
func TestInputRejectsInvalidUTF8(t *testing.T) {
	if _, ok := Input([]byte{0xff}); ok {
		t.Fatal("invalid UTF-8 accepted")
	}
}

// TestInputRemovesOnlyInitialBOMAndTrimsUnicodeSpace проверяет наблюдаемый контракт.
func TestInputRemovesOnlyInitialBOMAndTrimsUnicodeSpace(t *testing.T) {
	lines, ok := Input([]byte("\ufeff\u2003first\u00a0\r\n \ufeffsecond \n"))
	if !ok {
		t.Fatal("valid input rejected")
	}
	want := []Line{{1, "first"}, {2, "\ufeffsecond"}, {3, ""}}
	if len(lines) != len(want) {
		t.Fatalf("lines = %#v", lines)
	}
	for i := range want {
		if lines[i] != want[i] {
			t.Fatalf("line %d = %#v, want %#v", i, lines[i], want[i])
		}
	}
}

// TestInputPreservesPhysicalEmptyLines проверяет наблюдаемый контракт.
func TestInputPreservesPhysicalEmptyLines(t *testing.T) {
	lines, _ := Input([]byte("a\n\nb"))
	if len(lines) != 3 || lines[1].Number != 2 || lines[1].Text != "" {
		t.Fatalf("lines = %#v", lines)
	}
}
