package scan

import (
	"testing"
	"unicode/utf8"
)

// FuzzTokenize проверяет устойчивость и диапазоны токенизации произвольного UTF-8.
func FuzzTokenize(f *testing.F) {
	for _, seed := range []string{"", `one "two three"`, `"escaped \\" quote"`, `"незакрыто`} {
		f.Add(seed, 1)
	}
	f.Fuzz(func(t *testing.T, raw string, column int) {
		if !utf8.ValidString(raw) || column < 1 || column > 1_000_000 {
			t.Skip()
		}
		tokens := Tokenize(raw, column)
		last := column
		for _, token := range tokens.Values {
			if token.Start < last || token.End < token.Start || token.End > column+utf8.RuneCountInString(raw) {
				t.Fatalf("invalid token range [%d,%d) in %q", token.Start, token.End, raw)
			}
			last = token.End
		}
	})
}

// FuzzDecodeEscapes проверяет снятие экранирования и согласованность определения escape.
func FuzzDecodeEscapes(f *testing.F) {
	for _, seed := range []string{"", `\\@`, `\\\#`, `текст \\{`} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, raw string) {
		if !utf8.ValidString(raw) {
			t.Skip()
		}
		decoded := DecodeEscapes(raw, "@#{}")
		if !utf8.ValidString(decoded) || len(decoded) > len(raw) {
			t.Fatalf("DecodeEscapes(%q) = %q", raw, decoded)
		}
		for i := range raw {
			_ = Escaped(raw, i)
		}
	})
}
