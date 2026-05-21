package filter

import (
	"testing"
)

func TestNewInvalidPattern(t *testing.T) {
	_, err := New(Options{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern regex")
	}
}

func TestNewInvalidExclude(t *testing.T) {
	_, err := New(Options{Exclude: "[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid exclude regex")
	}
}

func TestIsEmpty(t *testing.T) {
	f, err := New(Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !f.IsEmpty() {
		t.Error("expected empty filter")
	}
}

func TestMatchPattern(t *testing.T) {
	f, err := New(Options{Pattern: `ERROR`})
	if err != nil {
		t.Fatal(err)
	}
	if !f.Match("2024-01-01 ERROR something went wrong") {
		t.Error("expected match for ERROR line")
	}
	if f.Match("2024-01-01 INFO all good") {
		t.Error("expected no match for INFO line")
	}
}

func TestMatchExclude(t *testing.T) {
	f, err := New(Options{Exclude: `DEBUG`})
	if err != nil {
		t.Fatal(err)
	}
	if f.Match("2024-01-01 DEBUG verbose output") {
		t.Error("expected DEBUG line to be excluded")
	}
	if !f.Match("2024-01-01 INFO important message") {
		t.Error("expected INFO line to pass exclude filter")
	}
}

func TestMatchKeywords(t *testing.T) {
	f, err := New(Options{Keywords: []string{"timeout", "db"}})
	if err != nil {
		t.Fatal(err)
	}
	if !f.Match("connection db timeout exceeded") {
		t.Error("expected match when all keywords present")
	}
	if f.Match("connection db established") {
		t.Error("expected no match when keyword 'timeout' missing")
	}
}

func TestMatchCombined(t *testing.T) {
	f, err := New(Options{
		Pattern:  `ERROR|WARN`,
		Keywords: []string{"disk"},
		Exclude:  `test`,
	})
	if err != nil {
		t.Fatal(err)
	}

	// passes all criteria
	if !f.Match("ERROR disk full") {
		t.Error("expected match")
	}
	// excluded by exclude pattern
	if f.Match("ERROR disk test failure") {
		t.Error("expected exclusion due to 'test'")
	}
	// missing keyword
	if f.Match("ERROR memory full") {
		t.Error("expected no match: missing keyword 'disk'")
	}
}
