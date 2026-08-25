package dslparser

import "testing"

// TestElementExposesDataAndDefensivelyCopiesSlices проверяет наблюдаемый контракт.
func TestElementExposesDataAndDefensivelyCopiesSlices(t *testing.T) {
	e := Element{kind: ElementGroup, name: ElementItems, raw: `\@raw`, value: "@raw", parsed: false, children: []Element{{raw: "child"}}, diagnostics: []Diagnostic{{code: UnclosedQuote, message: "ошибка"}}}
	if e.Kind() != ElementGroup || e.Name() != ElementItems || e.Raw() != `\@raw` || e.Value() != "@raw" || e.Parsed() {
		t.Fatal("element lost data")
	}
	c := e.Children()
	d := e.Diagnostics()
	c[0].raw = "bad"
	d[0].message = "bad"
	if e.Children()[0].Raw() != "child" || e.Diagnostics()[0].Message() != "ошибка" {
		t.Fatal("element exposed slices")
	}
}

// TestTechnicalNamesAreStable проверяет наблюдаемый контракт.
func TestTechnicalNamesAreStable(t *testing.T) {
	want := []string{"version", "level", "id", "name", "title", "content", "instruction", "media_type", "source", "source_kind", "body", "items", "before", "after", "upper", "lower", "positions", "separator"}
	got := []string{ElementVersion, ElementLevel, ElementID, ElementName, ElementTitle, ElementContent, ElementInstruction, ElementMediaType, ElementSource, ElementSourceKind, ElementBody, ElementItems, ElementBefore, ElementAfter, ElementUpper, ElementLower, ElementPositions, ElementSeparatorName}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("technical name %d = %q", i, got[i])
		}
	}
}
