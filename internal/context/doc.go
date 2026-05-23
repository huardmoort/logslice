// Package context provides lightweight helpers for propagating logslice-specific
// metadata through standard Go contexts.
//
// It wraps the standard library "context" package and adds typed accessors for
// values that are commonly threaded through the logslice processing pipeline:
//
//   - InputFile: the path of the log file currently being processed.
//   - JobID: an optional string identifier for the current slicing job,
//     useful when multiple jobs run concurrently or results are aggregated.
//
// Usage:
//
//	ctx := context.WithInputFile(context.Background(), "/var/log/app.log")
//	ctx  = context.WithJobID(ctx, "run-20240101-001")
//
//	// Later, in a deeply nested function:
//	file := context.InputFile(ctx)  // "/var/log/app.log"
//	job  := context.JobID(ctx)      // "run-20240101-001"
//
// Timeout and deadline helpers are also re-exported for convenience so callers
// need only import this package rather than both it and "context".
package context
