package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const defaultTimeFormat = "2006-01-02T15:04:05"

// ParseFlags parses command-line arguments into a Config.
// It returns the populated Config or an error.
func ParseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		fromStr   = fs.String("from", "", "start of time range (e.g. 2006-01-02T15:04:05)")
		toStr     = fs.String("to", "", "end of time range")
		timeFmt   = fs.String("time-format", defaultTimeFormat, "Go time layout for parsing timestamps")
		pattern   = fs.String("pattern", "", "include only lines matching this regex")
		exclude   = fs.String("exclude", "", "exclude lines matching this regex")
		outFile   = fs.String("output", "", "write output to file (default: stdout)")
		outFmt    = fs.String("format", "raw", "output format: raw | numbered | json")
		showStats = fs.Bool("stats", false, "print statistics after processing")
		skipUnp   = fs.Bool("skip-unparsed", false, "omit lines with no parseable timestamp")
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() < 1 {
		return nil, fmt.Errorf("usage: logslice [flags] <input-file>")
	}

	cfg := &Config{
		InputFile:    fs.Arg(0),
		TimeFormat:   *timeFmt,
		Pattern:      *pattern,
		Exclude:      *exclude,
		OutputFile:   *outFile,
		OutputFormat: Format(*outFmt),
		ShowStats:    *showStats,
		SkipUnparsed: *skipUnp,
	}

	var err error
	if *fromStr == "" {
		return nil, fmt.Errorf("--from is required")
	}
	if cfg.From, err = time.Parse(*timeFmt, *fromStr); err != nil {
		return nil, fmt.Errorf("invalid --from: %w", err)
	}
	if *toStr == "" {
		return nil, fmt.Errorf("--to is required")
	}
	if cfg.To, err = time.Parse(*timeFmt, *toStr); err != nil {
		return nil, fmt.Errorf("invalid --to: %w", err)
	}

	return cfg, cfg.Validate()
}
