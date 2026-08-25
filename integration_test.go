package dslparser

import (
	"os"
	"strings"
	"testing"
	"unicode/utf8"
)

// TestCompleteLessons проверяет полноценные проверенные уроки как единый контракт конвейера.
func TestCompleteLessons(t *testing.T) {
	files := []string{"01.txt", "02.txt", "03.txt"}
	allTags := make(map[string]int)
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			source := readLesson(t, file)
			document, err := Parse(source)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got := document.Version(); got != "1.1" {
				t.Errorf("Version() = %q, want 1.1", got)
			}
			if diagnostics := document.Diagnostics(); len(diagnostics) != 0 {
				t.Errorf("Diagnostics() = %v, want none", diagnosticCodes(diagnostics))
			}
			assertCompleteTree(t, document, source)
			Walk(document, func(cursor Cursor) bool {
				allTags[cursor.Node().CanonicalName()]++
				return true
			})
		})
	}
	for _, name := range []string{"task", "endtask", "header", "newpage", "step", "editor", "speaking", "media", "example", "wordlist", "table", "script", "text", "key", "instr", "note", "alt", "question", "multifill", "choice", "multichoice", "matching", "ordering", "variants", "variant"} {
		if allTags[name] == 0 {
			t.Errorf("prepared lessons do not cover registry tag @%s", name)
		}
	}
}

// readLesson читает подготовленный пользователем урок и проверяет требования к fixture.
func readLesson(t *testing.T, file string) []byte {
	t.Helper()
	source, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("ReadFile(%q): %v", file, err)
	}
	if !utf8.Valid(source) {
		t.Fatalf("%s is not valid UTF-8", file)
	}
	if first, _, _ := strings.Cut(string(source), "\n"); strings.TrimSuffix(first, "\r") != "@dsl-version 1.1" {
		t.Fatalf("first physical line = %q", first)
	}
	return source
}

// diagnosticCodes возвращает коды диагностик для краткого сообщения теста.
func diagnosticCodes(diagnostics []Diagnostic) []DiagnosticCode {
	codes := make([]DiagnosticCode, len(diagnostics))
	for i, diagnostic := range diagnostics {
		codes[i] = diagnostic.Code()
	}
	return codes
}

// assertCompleteTree проверяет диапазоны, разбор тегов, логические области и порядок обхода.
func assertCompleteTree(t *testing.T, document *Document, source []byte) {
	t.Helper()
	lineCount := strings.Count(string(source), "\n") + 1
	seen := make(map[string]int)
	var cursors []Cursor
	Walk(document, func(cursor Cursor) bool {
		cursors = append(cursors, cursor)
		node := cursor.Node()
		assertSpanInDocument(t, node.Span(), lineCount, node.Raw() == "")
		if node.Kind() == NodeTag {
			seen[node.CanonicalName()]++
			if !node.Parsed() {
				t.Errorf("tag @%s at %v is not parsed", node.OriginalName(), node.Span())
			}
		}
		return true
	})
	for _, name := range []string{"task", "step", "instr", "media", "multifill", "key", "example", "wordlist", "table", "script", "choice", "multichoice", "ordering", "question", "note", "alt"} {
		if seen[name] == 0 {
			t.Errorf("full lesson does not exercise @%s", name)
		}
	}
	assertCursorPaths(t, document.Roots(), cursors)
	assertTaskRegions(t, document)
}

// assertSpanInDocument проверяет ненулевой полуоткрытый диапазон внутри исходного документа.
func assertSpanInDocument(t *testing.T, span Span, lineCount int, empty bool) {
	t.Helper()
	start, end := span.Start(), span.End()
	if start.Line() < 1 || start.Line() > lineCount || start.Column() < 1 || end.Line() < start.Line() || end.Line() > lineCount+1 || !empty && end.Line() == start.Line() && end.Column() <= start.Column() {
		t.Errorf("invalid document span %v for %d lines", span, lineCount)
	}
}

// assertCursorPaths независимо строит ожидаемый прямой обход и сравнивает пути курсоров.
func assertCursorPaths(t *testing.T, roots []Node, got []Cursor) {
	t.Helper()
	var wantPaths [][]int
	var visit func([]Node, []int)
	visit = func(nodes []Node, parent []int) {
		for i, node := range nodes {
			path := append(append([]int(nil), parent...), i)
			wantPaths = append(wantPaths, path)
			visit(node.Children(), path)
		}
	}
	visit(roots, nil)
	if len(got) != len(wantPaths) {
		t.Fatalf("Walk visited %d nodes, want %d", len(got), len(wantPaths))
	}
	for i, cursor := range got {
		if !equalInts(cursor.Path(), wantPaths[i]) || cursor.Depth() != len(wantPaths[i])-1 {
			t.Errorf("cursor %d: path/depth = %v/%d, want %v/%d", i, cursor.Path(), cursor.Depth(), wantPaths[i], len(wantPaths[i])-1)
		}
	}
}

// equalInts сообщает равенство целочисленных срезов: true — равны, false — различаются.
func equalInts(left, right []int) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

// assertTaskRegions проверяет принадлежность этапов и явных завершений заданиям.
func assertTaskRegions(t *testing.T, document *Document) {
	t.Helper()
	Walk(document, func(cursor Cursor) bool {
		node := cursor.Node()
		if node.CanonicalName() != "task" {
			return true
		}
		for _, child := range node.Children() {
			if child.CanonicalName() == "step" && len(child.Diagnostics()) != 0 {
				t.Errorf("@step inside @task has diagnostics at %v", child.Span())
			}
			if child.CanonicalName() == "endtask" && len(child.Diagnostics()) != 0 {
				t.Errorf("@endtask inside @task has diagnostics at %v", child.Span())
			}
		}
		return true
	})
}
