package dslparser

import (
	"errors"
	"io"
	"strings"
	"testing"
)

// failingReader всегда завершает чтение ошибкой.
type failingReader struct{}

// Read возвращает воспроизводимую ошибку чтения без данных.
func (failingReader) Read([]byte) (int, error) { return 0, errors.New("boom") }

// TestParseEntryPointsShareVersionContract проверяет наблюдаемый контракт.
func TestParseEntryPointsShareVersionContract(t *testing.T) {
	for name, parse := range map[string]func() (*Document, error){"bytes": func() (*Document, error) { return Parse([]byte("@dsl-version 1.1")) }, "string": func() (*Document, error) { return ParseString("@dsl-version 1.1") }, "reader": func() (*Document, error) { return ParseReader(strings.NewReader("@dsl-version 1.1")) }} {
		t.Run(name, func(t *testing.T) {
			d, err := parse()
			if err != nil || d == nil || d.Version() != "1.1" || d.VersionRaw() != "1.1" || len(d.Roots()) != 0 {
				t.Fatalf("document=%#v error=%v", d, err)
			}
		})
	}
}

// TestParseFatalErrorsHaveStableCodes проверяет наблюдаемый контракт.
func TestParseFatalErrorsHaveStableCodes(t *testing.T) {
	cases := []struct {
		name string
		run  func() (*Document, error)
		code DiagnosticCode
	}{{"invalid UTF-8", func() (*Document, error) { return Parse([]byte{0xff}) }, InvalidUTF8}, {"empty", func() (*Document, error) { return Parse(nil) }, MissingVersion}, {"leading blank", func() (*Document, error) { return ParseString("\n@dsl-version 1.1") }, MissingVersion}, {"empty value", func() (*Document, error) { return ParseString("@dsl-version") }, MissingVersion}, {"malformed", func() (*Document, error) { return ParseString("@dsl-versionx 1.1") }, MissingVersion}, {"unsupported", func() (*Document, error) { return ParseString("@dsl-version 2.0") }, UnsupportedVersion}, {"reader", func() (*Document, error) { return ParseReader(failingReader{}) }, ReadFailure}, {"nil reader", func() (*Document, error) { return ParseReader(nil) }, ReadFailure}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := tc.run()
			if d != nil || err == nil {
				t.Fatalf("document=%v error=%v", d, err)
			}
			var fatal *FatalError
			if !errors.As(err, &fatal) || fatal.Code() != tc.code {
				t.Fatalf("error=%T %v", err, err)
			}
		})
	}
}

var _ io.Reader = failingReader{}
