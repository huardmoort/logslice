package transform

import (
	"strings"
	"testing"
)

func TestIsNoop(t *testing.T) {
	c := New()
	if !c.IsNoop() {
		t.Fatal("expected noop for empty chain")
	}
	c2 := New(TrimSpace())
	if c2.IsNoop() {
		t.Fatal("expected non-noop with funcs")
	}
}

func TestApplyEmpty(t *testing.T) {
	c := New()
	out, ok := c.Apply("hello")
	if !ok || out != "hello" {
		t.Fatalf("expected (hello, true), got (%q, %v)", out, ok)
	}
}

func TestTrimSpace(t *testing.T) {
	c := New(TrimSpace())
	out, ok := c.Apply("  hello world  ")
	if !ok || out != "hello world" {
		t.Fatalf("unexpected result: %q %v", out, ok)
	}
}

func TestDropEmpty(t *testing.T) {
	c := New(DropEmpty())
	_, ok := c.Apply("   ")
	if ok {
		t.Fatal("expected blank line to be dropped")
	}
	out, ok2 := c.Apply("data")
	if !ok2 || out != "data" {
		t.Fatalf("expected data to pass through, got %q %v", out, ok2)
	}
}

func TestReplaceAll(t *testing.T) {
	c := New(ReplaceAll("foo", "bar"))
	out, ok := c.Apply("foo and foo")
	if !ok || out != "bar and bar" {
		t.Fatalf("unexpected result: %q %v", out, ok)
	}
}

func TestMaxLength(t *testing.T) {
	c := New(MaxLength(5))
	out, ok := c.Apply("hello world")
	if !ok || out != "hello" {
		t.Fatalf("expected truncation to 5, got %q", out)
	}
	out2, ok2 := c.Apply("hi")
	if !ok2 || out2 != "hi" {
		t.Fatalf("short line should be unchanged, got %q", out2)
	}
}

func TestChainDropsOnFalse(t *testing.T) {
	called := false
	guard := Func(func(line string) (string, bool) {
		called = true
		return line, true
	})
	c := New(DropEmpty(), guard)
	_, ok := c.Apply("")
	if ok {
		t.Fatal("expected chain to stop on drop")
	}
	if called {
		t.Fatal("subsequent func should not be called after drop")
	}
}

func TestChainOrdering(t *testing.T) {
	c := New(TrimSpace(), ReplaceAll("x", "y"))
	out, ok := c.Apply("  xox  ")
	if !ok || out != "yoy" {
		t.Fatalf("unexpected: %q", out)
	}
	_ = strings.Contains // ensure import used
}
