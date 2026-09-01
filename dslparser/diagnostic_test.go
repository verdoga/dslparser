package dslparser

import (
	"errors"
	"testing"
)

// TestFatalErrorIsTypedAndKeepsCause проверяет наблюдаемый контракт.
func TestFatalErrorIsTypedAndKeepsCause(t *testing.T) {
	cause := errors.New("read")
	err := newFatalError(ReadFailure, "ошибка чтения", cause)
	var fatal *FatalError
	if !errors.As(err, &fatal) || fatal.Code() != ReadFailure || !errors.Is(err, cause) {
		t.Fatalf("unexpected fatal error: %v", err)
	}
}

// TestDiagnosticAccessorsAndOptionalRelatedSpan проверяет наблюдаемый контракт.
func TestDiagnosticAccessorsAndOptionalRelatedSpan(t *testing.T) {
	s := Span{start: Position{line: 1, column: 1}, end: Position{line: 1, column: 2}}
	d := Diagnostic{code: UnclosedQuote, message: "не закрыта кавычка", span: s}
	if d.Code() != UnclosedQuote || d.Message() == "" || d.Span() != s {
		t.Fatal("diagnostic lost data")
	}
	if _, ok := d.RelatedSpan(); ok {
		t.Fatal("unexpected related span")
	}
}
