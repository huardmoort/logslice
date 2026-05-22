// Package multifile provides utilities for treating multiple log files as a
// single sequential stream.
//
// # Overview
//
// Large deployments often rotate log files (app.log, app.log.1, app.log.2, …)
// or write to per-host shards. multifile lets logslice ingest all of them
// transparently:
//
//	mr, err := multifile.New([]string{"app.log.2", "app.log.1", "app.log"})
//
// Files are read in the order supplied. Use NewSorted with explicit Weight
// values to impose a custom ordering, or GlobFiles to expand shell-style
// wildcard patterns:
//
//	files, err := multifile.GlobFiles([]string{"/var/log/app/*.log"})
//	mr, err := multifile.NewSorted(files)
//
// MultiReader implements io.ReadCloser and can be passed directly to
// reader.NewLineReaderFromReader or any other component that accepts an
// io.Reader.
package multifile
