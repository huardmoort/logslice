// Package dedupe implements log line deduplication strategies for logslice.
//
// Two modes are supported:
//
//   - ModeConsecutive: suppresses a line only when it is identical to the
//     immediately preceding line. This is cheap (O(1) memory) and useful for
//     logs where repeated entries appear in bursts.
//
//   - ModeWindow: suppresses a line if the same content was seen anywhere
//     within the last N lines. The window is maintained as a circular buffer
//     of SHA-256 prefix hashes, keeping memory bounded regardless of line
//     length.
//
// Usage:
//
//	d, err := dedupe.New(dedupe.ModeConsecutive, 0)
//	if err != nil { ... }
//	for _, line := range lines {
//	    if d.Allow(line) {
//	        fmt.Println(line)
//	    }
//	}
//	fmt.Printf("suppressed %d duplicate lines\n", d.Suppressed)
package dedupe
