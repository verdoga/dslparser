package dslparser

import (
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"dslparser/internal/normalize"
	"dslparser/internal/registry"
)

// Parse разбирает байтовое представление документа или возвращает типизированную фатальную ошибку.
func Parse(data []byte) (*Document, error) { return parse(data) }

// ParseString разбирает строковое представление документа или возвращает типизированную фатальную ошибку.
func ParseString(source string) (*Document, error) { return parse([]byte(source)) }

// ParseReader читает весь документ и передаёт его общему конвейеру разбора.
func ParseReader(r io.Reader) (*Document, error) {
	if r == nil {
		return nil, newFatalError(ReadFailure, "не удалось прочитать документ: nil Reader", nil)
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, newFatalError(ReadFailure, "не удалось прочитать документ", err)
	}
	return parse(data)
}

// parse выполняет общий начальный конвейер разбора.
func parse(data []byte) (*Document, error) {
	lines, valid := normalize.Input(data)
	if !valid {
		return nil, newFatalError(InvalidUTF8, "вход содержит некорректный UTF-8", nil)
	}
	if len(lines) == 0 {
		return nil, newFatalError(MissingVersion, "отсутствует директива версии", nil)
	}
	const prefix = "@dsl-version"
	first := lines[0].Text
	if first == prefix || !strings.HasPrefix(first, prefix) || (len(first) > len(prefix) && !isSpaceAfterPrefix(first[len(prefix):])) {
		return nil, newFatalError(MissingVersion, "первая физическая строка не содержит корректную директиву версии", nil)
	}
	raw := strings.TrimSpace(first[len(prefix):])
	if raw == "" {
		return nil, newFatalError(MissingVersion, "значение версии отсутствует", nil)
	}
	if strings.IndexFunc(raw, unicode.IsSpace) >= 0 {
		return nil, newFatalError(MissingVersion, "директива версии повреждена", nil)
	}
	canonical := strings.ToLower(raw)
	if _, ok := registry.ForVersion(canonical); !ok {
		return nil, newFatalError(UnsupportedVersion, fmt.Sprintf("версия %q не поддерживается", raw), nil)
	}
	end := utf8.RuneCountInString(first) + 1
	return &Document{versionRaw: raw, version: canonical, versionSpan: Span{start: Position{line: 1, column: 1}, end: Position{line: 1, column: end}}}, nil
}

// isSpaceAfterPrefix сообщает, является ли первый символ хвоста Unicode-пробелом: true — является, false — нет.
func isSpaceAfterPrefix(tail string) bool {
	r, _ := utf8.DecodeRuneInString(tail)
	return unicode.IsSpace(r)
}
