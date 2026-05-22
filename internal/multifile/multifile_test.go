package multifile

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "mf-*.log")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestNewEmpty(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestNewSortedEmpty(t *testing.T) {
	_, err := NewSorted(nil)
	if err == nil {
		t.Fatal("expected error for empty files")
	}
}

func TestReadSequential(t *testing.T) {
	a := writeTmp(t, "line1\nline2\n")
	b := writeTmp(t, "line3\nline4\n")

	mr, err := New([]string{a, b})
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	data, err := io.ReadAll(mr)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)
	if !strings.Contains(got, "line1") || !strings.Contains(got, "line3") {
		t.Errorf("unexpected content: %q", got)
	}
}

func TestNewSortedByWeight(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.log")
	b := filepath.Join(dir, "b.log")
	os.WriteFile(a, []byte("AAA"), 0644)
	os.WriteFile(b, []byte("BBB"), 0644)

	mr, err := NewSorted([]File{
		{Path: b, Weight: 0},
		{Path: a, Weight: 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	data, err := io.ReadAll(mr)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "BBBAAA" {
		t.Errorf("expected BBBAAA, got %q", string(data))
	}
}

func TestCurrentPath(t *testing.T) {
	a := writeTmp(t, "hello")
	mr, err := New([]string{a})
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	buf := make([]byte, 1)
	mr.Read(buf)
	if mr.CurrentPath() != a {
		t.Errorf("expected %q, got %q", a, mr.CurrentPath())
	}
}

func TestReadMissingFile(t *testing.T) {
	mr, err := New([]string{"/nonexistent/file.log"})
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()
	buf := make([]byte, 4)
	_, err = mr.Read(buf)
	if err == nil {
		t.Fatal("expected error reading missing file")
	}
}
