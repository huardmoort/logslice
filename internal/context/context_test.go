package context_test

import (
	"testing"
	"time"

	lsctx "github.com/yourorg/logslice/internal/context"
)

func TestWithInputFile(t *testing.T) {
	ctx := lsctx.Background()
	ctx = lsctx.WithInputFile(ctx, "app.log")

	got := lsctx.InputFile(ctx)
	if got != "app.log" {
		t.Errorf("InputFile = %q; want %q", got, "app.log")
	}
}

func TestInputFileMissing(t *testing.T) {
	ctx := lsctx.Background()
	got := lsctx.InputFile(ctx)
	if got != "" {
		t.Errorf("InputFile = %q; want empty string", got)
	}
}

func TestWithJobID(t *testing.T) {
	ctx := lsctx.Background()
	ctx = lsctx.WithJobID(ctx, "job-42")

	got := lsctx.JobID(ctx)
	if got != "job-42" {
		t.Errorf("JobID = %q; want %q", got, "job-42")
	}
}

func TestJobIDMissing(t *testing.T) {
	ctx := lsctx.Background()
	got := lsctx.JobID(ctx)
	if got != "" {
		t.Errorf("JobID = %q; want empty string", got)
	}
}

func TestValuesAreIndependent(t *testing.T) {
	ctx := lsctx.Background()
	ctx = lsctx.WithInputFile(ctx, "server.log")
	ctx = lsctx.WithJobID(ctx, "job-99")

	if lsctx.InputFile(ctx) != "server.log" {
		t.Errorf("InputFile clobbered by WithJobID")
	}
	if lsctx.JobID(ctx) != "job-99" {
		t.Errorf("JobID clobbered by WithInputFile")
	}
}

func TestWithTimeout(t *testing.T) {
	ctx, cancel := lsctx.WithTimeout(lsctx.Background(), 100*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		// expected after timeout
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context did not expire within expected window")
	}
}

func TestWithDeadline(t *testing.T) {
	deadline := time.Now().Add(100 * time.Millisecond)
	ctx, cancel := lsctx.WithDeadline(lsctx.Background(), deadline)
	defer cancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context did not reach deadline within expected window")
	}
}
