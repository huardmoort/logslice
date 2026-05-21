// Package pipeline provides the top-level orchestration layer for logslice.
//
// It connects the individual processing stages — line reading, time-range
// slicing, pattern filtering, and output formatting — into a single
// streaming pass over a log file.
//
// Typical usage:
//
//	cfg := &config.Config{ ... }
//	result, err := pipeline.Run(cfg, inputReader, outputWriter)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("matched %d lines in %s\n", result.Stats.Matched(), result.Stats.Duration())
package pipeline
