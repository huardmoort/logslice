package window

import (
	"reflect"
	"testing"
)

func TestNewZeroSizeIsNoop(t *testing.T) {
	w := New(0)
	w.Add("line")
	if w.Len() != 0 {
		t.Fatalf("expected Len 0, got %d", w.Len())
	}
	if w.Lines() != nil {
		t.Fatal("expected nil Lines for noop window")
	}
}

func TestNewNegativeSizeIsNoop(t *testing.T) {
	w := New(-5)
	w.Add("x")
	if w.Len() != 0 {
		t.Fatalf("expected Len 0, got %d", w.Len())
	}
}

func TestAddAndLines(t *testing.T) {
	w := New(3)
	w.Add("a")
	w.Add("b")
	w.Add("c")
	got := w.Lines()
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestEvictsOldestWhenFull(t *testing.T) {
	w := New(3)
	w.Add("a")
	w.Add("b")
	w.Add("c")
	w.Add("d") // evicts "a"
	got := w.Lines()
	want := []string{"b", "c", "d"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestLenGrowsUpToCapacity(t *testing.T) {
	w := New(4)
	for i, line := range []string{"x", "y", "z"} {
		w.Add(line)
		if w.Len() != i+1 {
			t.Fatalf("step %d: expected Len %d, got %d", i, i+1, w.Len())
		}
	}
}

func TestLenCapsAtSize(t *testing.T) {
	w := New(2)
	w.Add("a")
	w.Add("b")
	w.Add("c")
	if w.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", w.Len())
	}
}

func TestReset(t *testing.T) {
	w := New(3)
	w.Add("a")
	w.Add("b")
	w.Reset()
	if w.Len() != 0 {
		t.Fatalf("expected Len 0 after Reset, got %d", w.Len())
	}
	if w.Lines() != nil {
		t.Fatal("expected nil Lines after Reset")
	}
	// window should still be usable after reset
	w.Add("c")
	if w.Len() != 1 {
		t.Fatalf("expected Len 1 after re-add, got %d", w.Len())
	}
}

func TestLinesReturnsCopy(t *testing.T) {
	w := New(2)
	w.Add("a")
	w.Add("b")
	lines := w.Lines()
	lines[0] = "mutated"
	got := w.Lines()
	if got[0] != "a" {
		t.Fatalf("Lines should return a copy; got %v", got)
	}
}
