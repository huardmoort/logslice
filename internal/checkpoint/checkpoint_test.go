package checkpoint_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/checkpoint"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "checkpoint-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestSaveAndLoad(t *testing.T) {
	store, err := checkpoint.NewStore(tempDir(t))
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}

	rec := checkpoint.Record{
		FilePath:  "/var/log/app.log",
		Offset:    4096,
		LineCount: 128,
	}
	if err := store.Save(rec); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := store.Load("/var/log/app.log")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Offset != 4096 {
		t.Errorf("Offset: got %d, want 4096", got.Offset)
	}
	if got.LineCount != 128 {
		t.Errorf("LineCount: got %d, want 128", got.LineCount)
	}
	if got.SavedAt.IsZero() {
		t.Error("SavedAt should be set")
	}
	if time.Since(got.SavedAt) > 5*time.Second {
		t.Errorf("SavedAt too old: %v", got.SavedAt)
	}
}

func TestLoadNotFound(t *testing.T) {
	store, _ := checkpoint.NewStore(tempDir(t))
	_, err := store.Load("/nonexistent/file.log")
	if !errors.Is(err, checkpoint.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	store, _ := checkpoint.NewStore(tempDir(t))
	rec := checkpoint.Record{FilePath: "/tmp/test.log", Offset: 100}
	_ = store.Save(rec)

	if err := store.Delete("/tmp/test.log"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := store.Load("/tmp/test.log")
	if !errors.Is(err, checkpoint.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNonExistent(t *testing.T) {
	store, _ := checkpoint.NewStore(tempDir(t))
	// Should not return an error for missing file.
	if err := store.Delete("/no/such/file.log"); err != nil {
		t.Errorf("Delete non-existent: unexpected error %v", err)
	}
}

func TestOverwriteCheckpoint(t *testing.T) {
	store, _ := checkpoint.NewStore(tempDir(t))
	path := "/logs/service.log"
	_ = store.Save(checkpoint.Record{FilePath: path, Offset: 1000})
	_ = store.Save(checkpoint.Record{FilePath: path, Offset: 2000})

	got, err := store.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Offset != 2000 {
		t.Errorf("Offset: got %d, want 2000", got.Offset)
	}
}
