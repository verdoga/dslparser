package scan

import "testing"

// TestTokenizeUnicodeWhitespaceQuotesEscapesAndRanges проверяет наблюдаемый контракт.
func TestTokenizeUnicodeWhitespaceQuotesEscapesAndRanges(t *testing.T) {
	r := Tokenize("one\u2003\"two words\" \"a\\\"b\\\\c\" \"\"", 5)
	wantRaw := []string{"one", `"two words"`, `"a\"b\\c"`, `""`}
	wantValue := []string{"one", "two words", `a"b\c`, ""}
	wantStart := []int{5, 9, 21, 31}
	if !r.Parsed || len(r.Values) != 4 {
		t.Fatalf("tokens=%#v", r)
	}
	for i := range wantRaw {
		if r.Values[i].Raw != wantRaw[i] || r.Values[i].Value != wantValue[i] || r.Values[i].Start != wantStart[i] || !r.Values[i].Parsed {
			t.Errorf("token %d=%#v", i, r.Values[i])
		}
	}
}

// TestTokenizeReturnsPartialTokenAndDiagnosticForUnclosedQuote проверяет наблюдаемый контракт.
func TestTokenizeReturnsPartialTokenAndDiagnosticForUnclosedQuote(t *testing.T) {
	r := Tokenize(`first "not closed`, 1)
	if r.Parsed || len(r.Values) != 2 || r.Values[0].Value != "first" || r.Values[1].Value != "not closed" || r.Values[1].Parsed || len(r.Diagnostics) != 1 || r.Diagnostics[0].Code != "UNCLOSED_QUOTE" {
		t.Fatalf("tokens=%#v", r)
	}
}

// TestTokenizeDoesNotInventSingleQuoteOrKeyValueSyntax проверяет наблюдаемый контракт.
func TestTokenizeDoesNotInventSingleQuoteOrKeyValueSyntax(t *testing.T) {
	r := Tokenize(`'two words' key=value`, 1)
	if len(r.Values) != 3 || r.Values[0].Value != "'two" || r.Values[1].Value != "words'" || r.Values[2].Value != "key=value" {
		t.Fatalf("tokens=%#v", r)
	}
}

// TestFreeTextDoesNotInterpretQuotes проверяет наблюдаемый контракт.
func TestFreeTextDoesNotInterpretQuotes(t *testing.T) {
	raw := `words "remain syntax"`
	if FreeText(raw) != raw {
		t.Fatal("free text changed")
	}
}

// TestUnknownTagCanUseUniversalTokens проверяет наблюдаемый контракт.
func TestUnknownTagCanUseUniversalTokens(t *testing.T) {
	c := Classify(`@unknown one "two words"`, Structural)
	tokens := Tokenize(c.Tail, c.NameEnd+1)
	if c.Kind != Tag || len(tokens.Values) != 2 || tokens.Values[1].Value != "two words" {
		t.Fatalf("candidate=%#v tokens=%#v", c, tokens)
	}
}
