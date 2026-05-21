package pipeline_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/pipeline"
)

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

const sampleLog = `2024-01-01T10:00:00Z INFO  service started
2024-01-01T10:01:00Z DEBUG request received path=/health
2024-01-01T10:02:00Z ERROR disk full
2024-01-01T10:03:00Z INFO  request received path=/api
2024-01-01T10:04:00Z INFO  service stopped
`

func TestRunBasic(t *testing.T) {
	cfg := &config.Config{
		From:            mustTime("2024-01-01T10:01:00Z"),
		To:              mustTime("2024-01-01T10:03:00Z"),
		Format:          "raw",
		TimestampFormat: time.RFC3339,
	}
	var out bytes.Buffer
	res, err := pipeline.Run(cfg, strings.NewReader(sampleLog), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res == nil || res.Stats == nil {
		t.Fatal("expected non-nil result")
	}
	if !strings.Contains(out.String(), "disk full") {
		t.Errorf("expected 'disk full' in output, got: %s", out.String())
	}
}

func TestRunWithPattern(t *testing.T) {
	cfg := &config.Config{
		From:            mustTime("2024-01-01T10:00:00Z"),
		To:              mustTime("2024-01-01T10:04:00Z"),
		Pattern:         "path=",
		Format:          "raw",
		TimestampFormat: time.RFC3339,
	}
	var out bytes.Buffer
	_, err := pipeline.Run(cfg, strings.NewReader(sampleLog), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 matched lines, got %d: %v", len(lines), lines)
	}
}

func TestRunNoMatch(t *testing.T) {
	cfg := &config.Config{
		From:            mustTime("2024-01-01T12:00:00Z"),
		To:              mustTime("2024-01-01T13:00:00Z"),
		Format:          "raw",
		TimestampFormat: time.RFC3339,
	}
	var out bytes.Buffer
	res, err := pipeline.Run(cfg, strings.NewReader(sampleLog), &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Len() != 0 {
		t.Errorf("expected empty output, got: %s", out.String())
	}
	if res.Stats.Matched() != 0 {
		t.Errorf("expected 0 matched, got %d", res.Stats.Matched())
	}
}
