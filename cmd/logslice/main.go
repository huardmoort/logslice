// Command logslice extracts time-range segments from large structured log files.
package main

import (
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/pipeline"
	"github.com/yourorg/logslice/internal/stats"
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(2)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "logslice: invalid config: %v\n", err)
		os.Exit(2)
	}

	f, err := os.Open(cfg.Input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logslice: cannot open input: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	result, err := pipeline.Run(cfg, f, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logslice: %v\n", err)
		os.Exit(1)
	}

	if cfg.ShowStats {
		reporter := stats.NewReporter(os.Stderr)
		reporter.Report(result.Stats)
	}
}
