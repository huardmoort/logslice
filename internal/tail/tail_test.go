package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/tail"
)

func TestFollowExistingLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, _ = f.WriteString("line one\nline two\n")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	follow, err := tail.Follow(ctx, f.Name(), tail.Options{PollInterval: 50 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	defer follow.Close()

	var got []string
	for line := range follow.Lines {
		got = append(got, line)
		if len(got) == 2 {
			follow.Close()
			break
		}
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestFollowNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-new-*.log")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	follow, err := tail.Follow(ctx, f.Name(), tail.Options{PollInterval: 30 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	defer follow.Close()

	go func() {
		time.Sleep(80 * time.Millisecond)
		_, _ = f.WriteString("appended line\n")
		f.Close()
	}()

	select {
	case line, ok := <-follow.Lines:
		if !ok {
			t.Fatal("channel closed before receiving line")
		}
		if line != "appended line\n" {
			t.Fatalf("unexpected line: %q", line)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for appended line")
	}
}

func TestFollowInvalidFile(t *testing.T) {
	_, err := tail.Follow(context.Background(), "/nonexistent/path/file.log", tail.Options{})
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}
