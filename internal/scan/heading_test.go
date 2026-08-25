package scan

import "testing"

// TestParseHeadingAcceptsAnyPositiveLevelAndUnicodeText проверяет наблюдаемый контракт.
func TestParseHeadingAcceptsAnyPositiveLevelAndUnicodeText(t *testing.T) {
	h := ParseHeading("#### Привет мир")
	if !h.Parsed || h.Level != 4 || h.Text != "Привет мир" || h.LevelStart != 1 || h.LevelEnd != 5 || h.TextStart != 6 || h.TextEnd != 16 {
		t.Fatalf("heading=%#v", h)
	}
}

// TestParseHeadingRejectsMissingSpaceOrText проверяет наблюдаемый контракт.
func TestParseHeadingRejectsMissingSpaceOrText(t *testing.T) {
	for _, line := range []string{"#bad", "#", "###   ", "plain"} {
		h := ParseHeading(line)
		if h.Parsed || len(h.Diagnostics) != 1 || h.Diagnostics[0].Code != "MALFORMED_HEADING" {
			t.Errorf("%q heading=%#v", line, h)
		}
	}
}
