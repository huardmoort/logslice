package redact

import (
	"testing"
)

func TestNewInvalidPattern(t *testing.T) {
	_, err := New([]string{"[invalid"}, "")
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestIsNoop(t *testing.T) {
	r, err := New(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.IsNoop() {
		t.Error("expected IsNoop true for empty pattern list")
	}
}

func TestIsNoopFalseWithPatterns(t *testing.T) {
	r, err := New([]string{`\d+`}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.IsNoop() {
		t.Error("expected IsNoop false when patterns present")
	}
}

func TestApplyNoPatterns(t *testing.T) {
	r, _ := New(nil, "")
	got := r.Apply("hello world")
	if got != "hello world" {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApplySinglePattern(t *testing.T) {
	r, err := New([]string{`password=\S+`}, "[REDACTED]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "user=alice password=secret123 action=login"
	want := "user=alice [REDACTED] action=login"
	got := r.Apply(input)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApplyMultiplePatterns(t *testing.T) {
	r, err := New([]string{
		`token=\S+`,
		`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`,
	}, "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := "token=abc123 src=192.168.1.1 msg=ok"
	got := r.Apply(input)
	if got == input {
		t.Error("expected line to be modified")
	}
	for _, sub := range []string{"abc123", "192.168.1.1"} {
		for _, c := range []byte(got) {
			_ = c
		}
		// ensure sensitive values are gone
		if contains(got, sub) {
			t.Errorf("sensitive value %q still present in %q", sub, got)
		}
	}
}

func TestApplyCustomPlaceholder(t *testing.T) {
	r, err := New([]string{`secret`}, "<hidden>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Apply("this is a secret value")
	want := "this is a <hidden> value"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApplyDefaultPlaceholder(t *testing.T) {
	r, err := New([]string{`secret`}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Apply("secret")
	if got != defaultPlaceholder {
		t.Errorf("got %q, want %q", got, defaultPlaceholder)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub ||
		len(s) > 0 && containsHelper(s, sub))
}

func containsHelper(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
