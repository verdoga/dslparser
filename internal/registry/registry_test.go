package registry

import "testing"

// TestVersion11ContainsCompleteTagSet проверяет наблюдаемый контракт.
func TestVersion11ContainsCompleteTagSet(t *testing.T) {
	r, ok := ForVersion("1.1")
	if !ok || r.Version() != "1.1" {
		t.Fatal("version 1.1 unavailable")
	}
	want := []string{"dsl-version", "task", "endtask", "header", "newpage", "step", "editor", "speaking", "media", "example", "wordlist", "table", "script", "text", "key", "instr", "note", "alt", "question", "multifill", "choice", "multichoice", "matching", "ordering", "variants", "variant"}
	entries := r.Entries()
	if len(entries) != len(want) {
		t.Fatalf("got %d entries, want %d", len(entries), len(want))
	}
	for i, name := range want {
		if entries[i].Name != name {
			t.Fatalf("entry %d = %q", i, entries[i].Name)
		}
		if _, ok := r.Lookup(name); !ok {
			t.Fatalf("missing %s", name)
		}
	}
}

// TestRegistryDescribesFormsBodiesAndSchemas проверяет наблюдаемый контракт.
func TestRegistryDescribesFormsBodiesAndSchemas(t *testing.T) {
	tests := []struct {
		name                            string
		forms                           Form
		args                            ArgumentSchema
		body                            BodyMode
		universal, special, bodyHandler HandlerID
		context                         ContextEffect
		empty                           bool
	}{
		{"dsl-version", FormInline, ArgsVersion, BodyNone, HandlerTokens, HandlerVersion, HandlerNone, ContextNone, false},
		{"editor", FormInline | FormBlock, ArgsContent, BodyOpaqueRaw, HandlerFreeText, HandlerNone, HandlerRawBody, ContextNone, false},
		{"media", FormInline, ArgsMediaTokens, BodyNone, HandlerTokens, HandlerMedia, HandlerNone, ContextNone, false},
		{"example", FormInline | FormBlock, ArgsContent, BodyOpaqueShaped, HandlerFreeText, HandlerNone, HandlerExample, ContextNone, false},
		{"text", FormBlock, ArgsOptionalTitle, BodyOpaqueRaw, HandlerFreeText, HandlerNone, HandlerRawBody, ContextNone, true},
		{"multifill", FormInline | FormBlock, ArgsOptionalInstruction, BodyOpaqueRaw, HandlerFreeText, HandlerNone, HandlerRawBody, ContextNone, true},
		{"choice", FormBlock, ArgsOptionalInstruction, BodyOpaqueShaped, HandlerFreeText, HandlerNone, HandlerItems, ContextNone, false},
		{"matching", FormBlock, ArgsOptionalInstruction, BodyOpaqueShaped, HandlerFreeText, HandlerNone, HandlerMatching, ContextNone, false},
		{"ordering", FormBlock, ArgsOptionalInstruction, BodyOpaqueShaped, HandlerFreeText, HandlerNone, HandlerOrdering, ContextNone, false},
		{"variants", FormBlock, ArgsNone, BodyStructural, HandlerNone, HandlerNone, HandlerStructural, ContextVariants, false},
		{"variant", FormInline, ArgsName, BodyLogical, HandlerFreeText, HandlerNone, HandlerNone, ContextVariant, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, _ := mustLookup(t, tt.name)
			if e.Forms != tt.forms || e.Arguments != tt.args || e.Body != tt.body || e.Universal != tt.universal || e.Special != tt.special || e.BodyHandler != tt.bodyHandler || e.Context != tt.context || e.PreserveEmpty != tt.empty {
				t.Fatalf("entry = %#v", e)
			}
		})
	}
}

// TestLookupIsCaseInsensitiveAndReturnedDataIsIndependent проверяет наблюдаемый контракт.
func TestLookupIsCaseInsensitiveAndReturnedDataIsIndependent(t *testing.T) {
	r, _ := ForVersion("1.1")
	entry, ok := r.Lookup("MeDiA")
	if !ok || entry.Name != "media" {
		t.Fatalf("entry=%#v", entry)
	}
	entries := r.Entries()
	entries[0].Name = "changed"
	again, _ := r.Lookup("dsl-version")
	if again.Name != "dsl-version" {
		t.Fatal("Entries exposed registry")
	}
	r2, _ := ForVersion("1.1")
	if r2.Entries()[0].Name != "dsl-version" {
		t.Fatal("registry instances share mutable data")
	}
}

// TestUnsupportedVersionsAreAbsent проверяет наблюдаемый контракт.
func TestUnsupportedVersionsAreAbsent(t *testing.T) {
	for _, v := range []string{"", "1.0", "1.10", "2.0", "1.1 "} {
		if _, ok := ForVersion(v); ok {
			t.Fatalf("version %q supported", v)
		}
	}
}

// mustLookup возвращает обязательную запись тестового реестра.
func mustLookup(t *testing.T, name string) (Entry, bool) {
	t.Helper()
	r, _ := ForVersion("1.1")
	e, ok := r.Lookup(name)
	if !ok {
		t.Fatalf("missing %s", name)
	}
	return e, ok
}
