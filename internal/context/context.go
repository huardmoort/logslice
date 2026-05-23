// Package context provides cancellation and timeout helpers for logslice
// pipeline operations, wrapping standard context with log-specific metadata.
package context

import (
	"context"
	"time"
)

// Key is an unexported type for context keys in this package.
type Key int

const (
	// KeyInputFile holds the name of the input file being processed.
	KeyInputFile Key = iota
	// KeyJobID holds an optional job identifier for tracking.
	KeyJobID
)

// WithInputFile returns a new context with the given input file name attached.
func WithInputFile(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, KeyInputFile, name)
}

// InputFile retrieves the input file name from the context.
// Returns an empty string if not set.
func InputFile(ctx context.Context) string {
	v, _ := ctx.Value(KeyInputFile).(string)
	return v
}

// WithJobID returns a new context with the given job ID attached.
func WithJobID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, KeyJobID, id)
}

// JobID retrieves the job ID from the context.
// Returns an empty string if not set.
func JobID(ctx context.Context) string {
	v, _ := ctx.Value(KeyJobID).(string)
	return v
}

// WithTimeout wraps context.WithTimeout for convenience.
func WithTimeout(parent context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, d)
}

// WithDeadline wraps context.WithDeadline for convenience.
func WithDeadline(parent context.Context, t time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(parent, t)
}

// Background returns context.Background for use across the package.
func Background() context.Context {
	return context.Background()
}
