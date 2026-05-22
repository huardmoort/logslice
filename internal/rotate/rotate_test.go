package rotate_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/rotate"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "rotatetest-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestNewInvalidPollInterval(t *testing.T) {
	path := writeTmp(t, "hello\n")
	_, err := rotate.New(path, 0)
	if err == nil {
		t.Fatal("expected error for zero poll interval")
	}
}

func TestNewMissingFile(t *testing.T) {
	_, err := rotate.New("/nonexistent/path/log.txt", time.Millisecond)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestCheckNoRotation(t *testing.T) {
	path := writeTmp(t, "line1\n")
	d, err := rotate.New(path, time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	_, rotated := d.Check()
	if rotated {
		t.Fatal("expected no rotation on stable file")
	}
}

func TestCheckTruncation(t *testing.T) {
	path := writeTmp(t, "some longer content\n")
	d, err := rotate.New(path, time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	// Truncate the file.
	if err := os.WriteFile(path, []byte("x\n"), 0644); err != nil {
		t.Fatal(err)
	}
	ev, rotated := d.Check()
	if !rotated {
		t.Fatal("expected rotation to be detected")
	}
	if ev != rotate.EventTruncated {
		t.Fatalf("expected EventTruncated, got %v", ev)
	}
}

func TestCheckFileGone(t *testing.T) {
	path := writeTmp(t, "data\n")
	d, err := rotate.New(path, time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	ev, rotated := d.Check()
	if !rotated {
		t.Fatal("expected rotation when file is gone")
	}
	if ev != rotate.EventReplaced {
		t.Fatalf("expected EventReplaced, got %v", ev)
	}
}

func TestWatchReceivesEvent(t *testing.T) {
	path := writeTmp(t, "initial content that is long\n")
	d, err := rotate.New(path, 5*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	stop := make(chan struct{})
	ch := d.Watch(stop)

	time.AfterFunc(20*time.Millisecond, func() {
		os.WriteFile(path, []byte("x\n"), 0644) //nolint:errcheck
	})

	select {
	case ev := <-ch:
		if ev != rotate.EventTruncated {
			t.Fatalf("expected EventTruncated, got %v", ev)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for rotation event")
	}
	close(stop)
}

func TestEventString(t *testing.T) {
	if rotate.EventTruncated.String() != "truncated" {
		t.Errorf("unexpected string for EventTruncated")
	}
	if rotate.EventReplaced.String() != "replaced" {
		t.Errorf("unexpected string for EventReplaced")
	}
}
