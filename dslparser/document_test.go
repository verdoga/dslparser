package dslparser

import "testing"

// TestDocumentAccessorsDefensivelyCopySlices проверяет наблюдаемый контракт.
func TestDocumentAccessorsDefensivelyCopySlices(t *testing.T) {
	d := &Document{versionRaw: "1.1", version: "1.1", roots: []Node{{kind: NodeText, raw: "one", parsed: true}}, diagnostics: []Diagnostic{{code: EmptyTagName, message: "пустое имя"}}}
	roots := d.Roots()
	diagnostics := d.Diagnostics()
	roots[0].raw = "changed"
	diagnostics[0].message = "changed"
	if d.Roots()[0].Raw() != "one" || d.Diagnostics()[0].Message() != "пустое имя" {
		t.Fatal("accessor exposed backing storage")
	}
}

// TestNilDocumentAccessorsAreSafe проверяет наблюдаемый контракт.
func TestNilDocumentAccessorsAreSafe(t *testing.T) {
	var d *Document
	if d.Version() != "" || d.Roots() != nil || d.Diagnostics() != nil {
		t.Fatal("unexpected nil document data")
	}
}
