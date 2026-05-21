// Package index provides log file indexing for fast time-range lookups.
//
// An index maps timestamps to byte offsets within a log file, allowing
// the slicer to seek directly to the relevant portion of a large file
// rather than scanning from the beginning.
//
// # Building an Index
//
// Use Build to scan a log file and produce an Index:
//
//	idx, err := index.Build(filename, format)
//
// # Finding a Range
//
// Use FindRange to locate the byte offsets for a time range:
//
//	start, end, err := idx.FindRange(from, to)
//
// # Caching
//
// NewCache wraps Build with an in-memory TTL cache keyed by filename,
// so repeated queries against the same file avoid redundant I/O.
package index
