package dslparser

import "testing"

// TestWalkReportsPreorderPathAndPrunesDescendants проверяет порядок, путь и локальное отсечение обхода.
func TestWalkReportsPreorderPathAndPrunesDescendants(t *testing.T) {
	doc, err := ParseString("@dsl-version 1.1\n# Unit\n### First\n@task ID\n@step Part\ntext\n@endtask\n### Second\ntext")
	if err != nil {
		t.Fatal(err)
	}
	var paths [][]int
	Walk(doc, func(c Cursor) bool {
		paths = append(paths, c.Path())
		return !(c.Node().CanonicalName() == "task")
	})
	want := [][]int{{0}, {0, 0}, {0, 1}, {0, 2}, {1}, {1, 0}, {1, 1}}
	if len(paths) != len(want) {
		t.Fatalf("paths = %v, want %v", paths, want)
	}
	for i := range want {
		if len(paths[i]) != len(want[i]) {
			t.Fatalf("path %d = %v, want %v", i, paths[i], want[i])
		}
		for j := range want[i] {
			if paths[i][j] != want[i][j] {
				t.Fatalf("path %d = %v, want %v", i, paths[i], want[i])
			}
		}
	}
}

// TestWalkNilInputs проверяет безопасные пустые вызовы.
func TestWalkNilInputs(t *testing.T) {
	Walk(nil, func(Cursor) bool { t.Fatal("unexpected visit"); return true })
	doc, err := ParseString("@dsl-version 1.1")
	if err != nil {
		t.Fatal(err)
	}
	Walk(doc, nil)
}
