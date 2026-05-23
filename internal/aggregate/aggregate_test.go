package aggregate

import (
	"strings"
	"testing"
)

func TestNewInvalidMode(t *testing.T) {
	_, err := New("histogram")
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNewValidModes(t *testing.T) {
	for _, m := range []Mode{ModeCount, ModeUnique} {
		_, err := New(m)
		if err != nil {
			t.Fatalf("unexpected error for mode %q: %v", m, err)
		}
	}
}

func TestCountMode(t *testing.T) {
	a, _ := New(ModeCount)
	a.Add("error")
	a.Add("warn")
	a.Add("error")
	a.Add("error")

	results := a.Results()
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// sorted: error then warn
	if !strings.HasPrefix(results[0], "error\t3") {
		t.Errorf("expected error\t3, got %q", results[0])
	}
	if !strings.HasPrefix(results[1], "warn\t1") {
		t.Errorf("expected warn\t1, got %q", results[1])
	}
}

func TestUniqueMode(t *testing.T) {
	a, _ := New(ModeUnique)
	a.Add("alpha")
	a.Add("beta")
	a.Add("alpha")

	results := a.Results()
	if len(results) != 2 {
		t.Fatalf("expected 2 unique results, got %d", len(results))
	}
	if results[0] != "alpha" || results[1] != "beta" {
		t.Errorf("unexpected results: %v", results)
	}
}

func TestAddTrimsWhitespace(t *testing.T) {
	a, _ := New(ModeCount)
	a.Add("  info  ")
	a.Add("info")
	if a.Total() != 2 {
		t.Errorf("expected total 2, got %d", a.Total())
	}
	if len(a.Results()) != 1 {
		t.Errorf("expected 1 unique key after trim, got %d", len(a.Results()))
	}
}

func TestAddEmptyIgnored(t *testing.T) {
	a, _ := New(ModeCount)
	a.Add("")
	a.Add("   ")
	if a.Total() != 0 {
		t.Errorf("expected total 0, got %d", a.Total())
	}
}

func TestReset(t *testing.T) {
	a, _ := New(ModeCount)
	a.Add("x")
	a.Add("x")
	a.Reset()
	if a.Total() != 0 {
		t.Errorf("expected 0 after reset, got %d", a.Total())
	}
	if len(a.Results()) != 0 {
		t.Errorf("expected empty results after reset")
	}
}

func TestTotal(t *testing.T) {
	a, _ := New(ModeCount)
	a.Add("a")
	a.Add("b")
	a.Add("a")
	if got := a.Total(); got != 3 {
		t.Errorf("expected total 3, got %d", got)
	}
}
