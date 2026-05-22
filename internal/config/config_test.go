package config

import (
	"testing"
	"time"
)

func baseConfig() *Config {
	return &Config{
		InputFile:    "test.log",
		From:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		To:           time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
		OutputFormat: FormatRaw,
	}
}

func TestValidateOK(t *testing.T) {
	if err := baseConfig().Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateMissingInput(t *testing.T) {
	c := baseConfig()
	c.InputFile = ""
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing input file")
	}
}

func TestValidateMissingFrom(t *testing.T) {
	c := baseConfig()
	c.From = time.Time{}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing --from")
	}
}

func TestValidateToBeforeFrom(t *testing.T) {
	c := baseConfig()
	c.To = c.From.Add(-time.Hour)
	if err := c.Validate(); err == nil {
		t.Fatal("expected error when --to is before --from")
	}
}

func TestValidateToEqualFrom(t *testing.T) {
	c := baseConfig()
	c.To = c.From
	if err := c.Validate(); err == nil {
		t.Fatal("expected error when --to equals --from")
	}
}

func TestValidateDefaultFormat(t *testing.T) {
	c := baseConfig()
	c.OutputFormat = ""
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.OutputFormat != FormatRaw {
		t.Errorf("expected default format %q, got %q", FormatRaw, c.OutputFormat)
	}
}

func TestValidateUnknownFormat(t *testing.T) {
	c := baseConfig()
	c.OutputFormat = "xml"
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for unknown format")
	}
}
