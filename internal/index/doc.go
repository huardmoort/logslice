// Package index provides byte-offset indexing for log files to enable
// fast seeking into large files without scanning from the beginning.
//
// # Overview
//
// An Index is built by scanning a log file and recording the byte offset
// and parsed timestamp of each line. Once built, FindRange can locate
// the byte offsets that bracket a requested time range, allowing the
// slicer to seek directly to the relevant portion of the file.
//
// # Caching
//
// The Cache type wraps Index construction with an in-memory LRU-style
// store keyed by file path (or any string identifier). A configurable
// TTL controls how long entries remain valid before being rebuilt on
// the next access.
//
// # Usage
//
//	cache := index.NewCache(5 * time.Minute)
//	idx, err := cache.GetOrBuild(filePath, readSeeker, timestampFormat)
//	start, end, ok := idx.FindRange(from, to)
package index
