package config

import (
	"testing"
)

func TestParseFlagsMinimal(t *testing.T) {
	args := []string{
		"--from", "2024-01-01T00:00:00",
		"--to", "2024-01-01T01:00:00",
		"app.log",
	}
	cfg, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.InputFile != "app.log" {
		t.Errorf("expected InputFile=app.log, got %q", cfg.InputFile)
	}
	if cfg.OutputFormat != FormatRaw {
		t.Errorf("expected format raw, got %q", cfg.OutputFormat)
	}
}

func TestParseFlagsAllOptions(t *testing.T) {
	args := []string{
		"--from", "2024-03-01T08:00:00",
		"--to", "2024-03-01T09:00:00",
		"--format", "json",
		"--pattern", "ERROR",
		"--exclude", "DEBUG",
		"--stats",
		"--skip-unparsed",
		"service.log",
	}
	cfg, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.OutputFormat != FormatJSON {
		t.Errorf("expected json format, got %q", cfg.OutputFormat)
	}
	if cfg.Pattern != "ERROR" {
		t.Errorf("expected pattern ERROR, got %q", cfg.Pattern)
	}
	if !cfg.ShowStats {
		t.Error("expected ShowStats=true")
	}
	if !cfg.SkipUnparsed {
		t.Error("expected SkipUnparsed=true")
	}
}

func TestParseFlagsMissingInput(t *testing.T) {
	args := []string{"--from", "2024-01-01T00:00:00", "--to", "2024-01-01T01:00:00"}
	if _, err := ParseFlags(args); err == nil {
		t.Fatal("expected error for missing input file")
	}
}

func TestParseFlagsInvalidFrom(t *testing.T) {
	args := []string{"--from", "not-a-date", "--to", "2024-01-01T01:00:00", "a.log"}
	if _, err := ParseFlags(args); err == nil {
		t.Fatal("expected error for invalid --from")
	}
}

func TestParseFlagsToBeforeFrom(t *testing.T) {
	args := []string{
		"--from", "2024-01-01T02:00:00",
		"--to", "2024-01-01T01:00:00",
		"a.log",
	}
	if _, err := ParseFlags(args); err == nil {
		t.Fatal("expected error when --to is before --from")
	}
}
