// Package checkpoint implements resume support for logslice by persisting
// the last successfully processed position within a log file.
//
// When processing large log files that may be rotated or appended to over time,
// a checkpoint allows logslice to continue from where it left off rather than
// re-scanning from the beginning.
//
// Usage:
//
//	store, err := checkpoint.NewStore("/tmp/logslice/checkpoints")
//	if err != nil { ... }
//
//	// Resume from last position.
//	rec, err := store.Load("/var/log/app.log")
//	if errors.Is(err, checkpoint.ErrNotFound) {
//		// No prior checkpoint; start from beginning.
//	}
//
//	// After processing, persist the new position.
//	_ = store.Save(checkpoint.Record{
//		FilePath:  "/var/log/app.log",
//		Offset:    bytesRead,
//		LineCount: linesProcessed,
//	})
package checkpoint
