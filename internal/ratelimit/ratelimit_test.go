package ratelimit

import (
	"testing"
	"time"
)

func TestNewInvalidRate(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero rate")
	}
	_, err = New(-5)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestNewValidRate(t *testing.T) {
	l, err := New(100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
}

func TestAllowConsumesTokens(t *testing.T) {
	l, _ := New(3)
	// Freeze time so no refill happens.
	fixed := time.Now()
	l.clock = func() time.Time { return fixed }
	l.lastTick = fixed
	l.tokens = 3

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected Allow()=false after tokens exhausted")
	}
}

func TestAllowRefillsOverTime(t *testing.T) {
	l, _ := New(10)
	base := time.Now()
	l.clock = func() time.Time { return base }
	l.lastTick = base
	l.tokens = 0

	// Advance clock by 1 second — should refill 10 tokens.
	base = base.Add(time.Second)
	for i := 0; i < 10; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d after refill", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("expected Allow()=false after refilled tokens exhausted")
	}
}

func TestTokensCapAtMax(t *testing.T) {
	l, _ := New(5)
	base := time.Now()
	l.clock = func() time.Time { return base }
	l.lastTick = base
	l.tokens = 0

	// Advance by 10 seconds — should cap at max (5), not 50.
	base = base.Add(10 * time.Second)
	l.Allow() // trigger refill
	if l.tokens > l.max {
		t.Fatalf("tokens %f exceeded max %f", l.tokens, l.max)
	}
}

func TestSetRate(t *testing.T) {
	l, _ := New(1)
	if err := l.SetRate(50); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.rate != 50 {
		t.Fatalf("expected rate 50, got %f", l.rate)
	}
	if err := l.SetRate(0); err == nil {
		t.Fatal("expected error for zero rate in SetRate")
	}
}
