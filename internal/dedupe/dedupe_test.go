package dedupe

import (
	"testing"
)

func TestNewInvalidWindowSize(t *testing.T) {
	_, err := New(ModeWindow, 0)
	if err == nil {
		t.Fatal("expected error for window size 0")
	}
}

func TestConsecutiveAllowsFirst(t *testing.T) {
	d, _ := New(ModeConsecutive, 0)
	if !d.Allow("hello") {
		t.Fatal("first line should be allowed")
	}
}

func TestConsecutiveSuppressDuplicate(t *testing.T) {
	d, _ := New(ModeConsecutive, 0)
	d.Allow("hello")
	if d.Allow("hello") {
		t.Fatal("consecutive duplicate should be suppressed")
	}
	if d.Suppressed != 1 {
		t.Fatalf("expected Suppressed=1, got %d", d.Suppressed)
	}
}

func TestConsecutiveAllowsAfterChange(t *testing.T) {
	d, _ := New(ModeConsecutive, 0)
	d.Allow("hello")
	d.Allow("hello") // suppressed
	if !d.Allow("world") {
		t.Fatal("different line should be allowed")
	}
	if !d.Allow("hello") {
		t.Fatal("hello should be allowed again after different line")
	}
}

func TestWindowSuppressesWithinWindow(t *testing.T) {
	d, _ := New(ModeWindow, 3)
	d.Allow("a")
	d.Allow("b")
	if d.Allow("a") {
		t.Fatal("'a' should be suppressed within window of 3")
	}
	if d.Suppressed != 1 {
		t.Fatalf("expected Suppressed=1, got %d", d.Suppressed)
	}
}

func TestWindowAllowsOutsideWindow(t *testing.T) {
	d, _ := New(ModeWindow, 2)
	d.Allow("a")
	d.Allow("b")
	d.Allow("c") // evicts "a" from window
	if !d.Allow("a") {
		t.Fatal("'a' should be allowed after falling out of window")
	}
}

func TestWindowAllowsNewLine(t *testing.T) {
	d, _ := New(ModeWindow, 5)
	lines := []string{"x", "y", "z"}
	for _, l := range lines {
		if !d.Allow(l) {
			t.Fatalf("expected %q to be allowed", l)
		}
	}
}

func TestHashConsistency(t *testing.T) {
	if hash("foo") != hash("foo") {
		t.Fatal("hash should be deterministic")
	}
	if hash("foo") == hash("bar") {
		t.Fatal("different strings should have different hashes")
	}
}
