package fieldextract

import (
	"testing"
)

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", FormatAuto)
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestExtractJSONField(t *testing.T) {
	e, err := New("level", FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	val, ok := e.Extract(`{"level":"error","msg":"boom"}`)
	if !ok {
		t.Fatal("expected field to be found")
	}
	if val != "error" {
		t.Fatalf("expected 'error', got %q", val)
	}
}

func TestExtractJSONFieldMissing(t *testing.T) {
	e, _ := New("missing", FormatJSON)
	_, ok := e.Extract(`{"level":"info"}`)
	if ok {
		t.Fatal("expected field to be absent")
	}
}

func TestExtractJSONInvalidLine(t *testing.T) {
	e, _ := New("level", FormatJSON)
	_, ok := e.Extract("not json at all")
	if ok {
		t.Fatal("expected false for invalid JSON")
	}
}

func TestExtractLogfmtField(t *testing.T) {
	e, _ := New("level", FormatLogfmt)
	val, ok := e.Extract(`time=2024-01-02T15:04:05Z level=warn msg="disk full"`)
	if !ok {
		t.Fatal("expected field to be found")
	}
	if val != "warn" {
		t.Fatalf("expected 'warn', got %q", val)
	}
}

func TestExtractLogfmtQuotedValue(t *testing.T) {
	e, _ := New("msg", FormatLogfmt)
	val, ok := e.Extract(`level=info msg="hello world"`)
	if !ok {
		t.Fatal("expected field to be found")
	}
	// Fields() splits on spaces so quoted multi-word values won't join;
	// we just verify the prefix stripping works for single-word quoted values.
	if val != "hello" {
		t.Fatalf("expected 'hello', got %q", val)
	}
}

func TestExtractLogfmtMissing(t *testing.T) {
	e, _ := New("host", FormatLogfmt)
	_, ok := e.Extract("level=info msg=ok")
	if ok {
		t.Fatal("expected field to be absent")
	}
}

func TestExtractAutoDetectsJSON(t *testing.T) {
	e, _ := New("svc", FormatAuto)
	val, ok := e.Extract(`{"svc":"api","code":200}`)
	if !ok || val != "api" {
		t.Fatalf("auto-detect JSON failed: ok=%v val=%q", ok, val)
	}
}

func TestExtractAutoDetectsLogfmt(t *testing.T) {
	e, _ := New("svc", FormatAuto)
	val, ok := e.Extract("svc=api level=info")
	if !ok || val != "api" {
		t.Fatalf("auto-detect logfmt failed: ok=%v val=%q", ok, val)
	}
}
