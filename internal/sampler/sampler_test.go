package sampler

import (
	"testing"
)

func TestNewInvalidRate(t *testing.T) {
	_, err := New("nth", 0, 0)
	if err == nil {
		t.Fatal("expected error for rate=0")
	}
}

func TestNewInvalidMode(t *testing.T) {
	_, err := New("unknown", 2, 0)
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNthMode(t *testing.T) {
	s, err := New("nth", 3, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// lines 1..9; kept: 3,6,9
	expected := []bool{false, false, true, false, false, true, false, false, true}
	for i, want := range expected {
		got := s.Keep()
		if got != want {
			t.Errorf("line %d: Keep()=%v, want %v", i+1, got, want)
		}
	}
}

func TestNthModeRate1KeepsAll(t *testing.T) {
	s, err := New("nth", 1, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 10; i++ {
		if !s.Keep() {
			t.Errorf("line %d: expected Keep()=true for rate=1", i+1)
		}
	}
}

func TestRandomModeApproximateRate(t *testing.T) {
	s, err := New("random", 10, 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const total = 10000
	kept := 0
	for i := 0; i < total; i++ {
		if s.Keep() {
			kept++
		}
	}
	// expect ~10% kept; allow ±5%
	pct := float64(kept) / float64(total)
	if pct < 0.05 || pct > 0.15 {
		t.Errorf("random rate=10: kept %.2f%%, want ~10%%", pct*100)
	}
}

func TestReset(t *testing.T) {
	s, _ := New("nth", 2, 0)
	s.Keep() // 1 – skip
	s.Keep() // 2 – keep
	s.Reset()
	// after reset counter is 0, so next call is line 1 – skip
	if s.Keep() {
		t.Error("expected first line after Reset to be skipped for rate=2")
	}
	if !s.Keep() {
		t.Error("expected second line after Reset to be kept for rate=2")
	}
}
