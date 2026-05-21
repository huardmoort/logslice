package config

import (
	"errors"
	"time"
)

// Format represents the output format type.
type Format string

const (
	FormatRaw      Format = "raw"
	FormatNumbered Format = "numbered"
	FormatJSON     Format = "json"
)

// Config holds all runtime configuration for a logslice run.
type Config struct {
	// Input
	InputFile string

	// Time range
	From      time.Time
	To        time.Time
	TimeFormat string

	// Filtering
	Pattern string
	Exclude string

	// Output
	OutputFile string
	OutputFormat Format
	ShowStats  bool

	// Behaviour
	SkipUnparsed bool
}

// Validate checks that the config is internally consistent.
func (c *Config) Validate() error {
	if c.InputFile == "" {
		return errors.New("input file must be specified")
	}
	if c.From.IsZero() {
		return errors.New("--from timestamp is required")
	}
	if c.To.IsZero() {
		return errors.New("--to timestamp is required")
	}
	if !c.To.After(c.From) {
		return errors.New("--to must be after --from")
	}
	switch c.OutputFormat {
	case FormatRaw, FormatNumbered, FormatJSON:
		// valid
	case "":
		c.OutputFormat = FormatRaw
	default:
		return errors.New("unknown output format: " + string(c.OutputFormat))
	}
	return nil
}
