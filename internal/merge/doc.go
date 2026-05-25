// Package merge implements a k-way merge for sorted log entry streams.
//
// When logslice processes multiple input files simultaneously, each file
// produces a stream of timestamped entries in ascending order. The merge
// package combines these streams into a single chronologically ordered
// output using a min-heap, achieving O(N log k) time complexity where N
// is the total number of entries and k is the number of sources.
//
// Basic usage:
//
//	src1 := make(chan merge.Entry)
//	src2 := make(chan merge.Entry)
//	// ... populate sources from readers ...
//
//	m := merge.New([]<-chan merge.Entry{src1, src2})
//	for entry := range m.Merge() {
//		fmt.Println(entry.Line)
//	}
//
// Each source channel must emit entries in ascending timestamp order.
// The Merge method is safe to call once per Merger instance.
package merge
