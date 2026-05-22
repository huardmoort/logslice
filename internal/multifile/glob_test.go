package multifile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGlobFilesNoPatterns(t *testing.T) {
	_, err := GlobFiles(nil)
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestGlobFilesNoMatch(t *testing.T) {
	_, err := GlobFiles([]string{"/no/such/path/*.log"})
	if err == nil {
		t.Fatal("expected error when no files match")
	}
}

func TestGlobFilesMatchesSorted(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"c.log", "a.log", "b.log"} {
		os.WriteFile(filepath.Join(dir, name), []byte(name), 0644)
	}

	files, err := GlobFiles([]string{filepath.Join(dir, "*.log")})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(files))
	}
	expected := []string{"a.log", "b.log", "c.log"}
	for i, f := range files {
		if filepath.Base(f.Path) != expected[i] {
			t.Errorf("files[%d] = %q, want %q", i, f.Path, expected[i])
		}
	}
}

func TestGlobFilesDeduplicates(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "only.log")
	os.WriteFile(p, []byte("x"), 0644)

	files, err := GlobFiles([]string{p, p})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Errorf("expected 1 file after dedup, got %d", len(files))
	}
}

func TestGlobFilesMultiplePatterns(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "app.log"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir, "sys.log"), []byte("s"), 0644)

	files, err := GlobFiles([]string{
		filepath.Join(dir, "app.log"),
		filepath.Join(dir, "sys.log"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
}
