// Package merge provides utilities for merging multiple sorted log streams
// into a single chronologically ordered output.
package merge

import (
	"container/heap"
	"time"
)

// Entry represents a single log line with its parsed timestamp and source index.
type Entry struct {
	Line      string
	Timestamp time.Time
	Source    int
}

// entryHeap implements heap.Interface for min-heap ordering by timestamp.
type entryHeap []Entry

func (h entryHeap) Len() int            { return len(h) }
func (h entryHeap) Less(i, j int) bool { return h[i].Timestamp.Before(h[j].Timestamp) }
func (h entryHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *entryHeap) Push(x any) {
	*h = append(*h, x.(Entry))
}

func (h *entryHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Merger merges pre-sorted entry channels into a single ordered stream.
type Merger struct {
	sources []<-chan Entry
}

// New creates a Merger from the provided source channels.
// Each source must emit entries in ascending timestamp order.
func New(sources []<-chan Entry) *Merger {
	return &Merger{sources: sources}
}

// Merge reads from all sources and emits entries in timestamp order.
// It closes the returned channel when all sources are exhausted.
func (m *Merger) Merge() <-chan Entry {
	out := make(chan Entry, 64)
	go func() {
		defer close(out)
		h := &entryHeap{}
		heap.Init(h)

		channels := make([]<-chan Entry, len(m.sources))
		copy(channels, m.sources)

		// Seed the heap with the first entry from each source.
		for i, ch := range channels {
			if e, ok := <-ch; ok {
				e.Source = i
				heap.Push(h, e)
			}
		}

		for h.Len() > 0 {
			smallest := heap.Pop(h).(Entry)
			out <- smallest
			if next, ok := <-channels[smallest.Source]; ok {
				next.Source = smallest.Source
				heap.Push(h, next)
			}
		}
	}()
	return out
}
