// Package rotate detects log file rotation events for logslice.
//
// A Detector watches a single file path and reports when the file has been
// rotated, either by truncation (the file shrank in size) or replacement
// (the underlying inode changed, as happens with tools like logrotate).
//
// Usage:
//
//	d, err := rotate.New("/var/log/app.log", 500*time.Millisecond)
//	if err != nil { ... }
//
//	stop := make(chan struct{})
//	for ev := range d.Watch(stop) {
//		switch ev {
//		case rotate.EventTruncated:
//			// seek to beginning
//		case rotate.EventReplaced:
//			// reopen file
//		}
//	}
//
// The Watch method polls the file at the configured interval. For single
// one-shot checks, call Check directly.
package rotate
