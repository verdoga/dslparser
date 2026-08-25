package dslparser

import "testing"

// TestPositionAndSpanUseOneBasedHalfOpenCoordinates проверяет наблюдаемый контракт.
func TestPositionAndSpanUseOneBasedHalfOpenCoordinates(t *testing.T) {
	start, end := (Position{line: 2, column: 3}), (Position{line: 2, column: 6})
	span := Span{start: start, end: end}
	if start.Line() != 2 || start.Column() != 3 || span.Start() != start || span.End() != end {
		t.Fatalf("unexpected coordinates: %#v %#v", start, span)
	}
}

// TestUnicodeColumnCountsRunes проверяет наблюдаемый контракт.
func TestUnicodeColumnCountsRunes(t *testing.T) {
	doc, err := ParseString("@dsl-version 1.1")
	if err != nil {
		t.Fatal(err)
	}
	if got := doc.VersionSpan().End().Column(); got != 17 {
		t.Fatalf("end column = %d, want 17", got)
	}
}
