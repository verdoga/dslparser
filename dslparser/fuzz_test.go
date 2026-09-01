package dslparser

import (
	"errors"
	"testing"
	"unicode/utf8"
)

// FuzzParseBytes проверяет, что произвольные байты дают документ либо типизированную фатальную ошибку.
func FuzzParseBytes(f *testing.F) {
	for _, seed := range [][]byte{nil, {0xff}, []byte("@dsl-version 1.1\n@task id"), []byte("\xef\xbb\xbf@dsl-version 1.1\nтекст")} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, source []byte) {
		document, err := Parse(source)
		if err != nil {
			var fatal *FatalError
			if document != nil || !errors.As(err, &fatal) {
				t.Fatalf("Parse() = (%v, %T), want nil document and *FatalError", document, err)
			}
			return
		}
		if document == nil {
			t.Fatal("Parse() returned nil document and nil error")
		}
	})
}

// FuzzParseUTF8 проверяет детерминированность разбора корректного UTF-8 после допустимой версии.
func FuzzParseUTF8(f *testing.F) {
	for _, seed := range []string{"", "текст", "@task ID\n@step Этап\n@endtask", "@text {\n}\nнепрозрачно"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, body string) {
		if !utf8.ValidString(body) {
			t.Skip()
		}
		source := "@dsl-version 1.1\n" + body
		first, firstErr := ParseString(source)
		second, secondErr := ParseString(source)
		if firstErr != nil || secondErr != nil {
			t.Fatalf("valid version produced errors %v and %v", firstErr, secondErr)
		}
		if documentFingerprint(first) != documentFingerprint(second) {
			t.Fatal("repeated parsing is not deterministic")
		}
	})
}

// documentFingerprint строит сравнимый снимок наблюдаемой структуры документа.
func documentFingerprint(document *Document) string {
	result := document.Version() + "|"
	Walk(document, func(cursor Cursor) bool {
		node := cursor.Node()
		result += node.Raw() + ":" + node.CanonicalName() + ";"
		return true
	})
	for _, diagnostic := range document.Diagnostics() {
		result += string(diagnostic.Code()) + ";"
	}
	return result
}
