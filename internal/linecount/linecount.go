// Package linecount provides utilities for counting and tracking line
// positions within log files, enabling efficient offset-based navigation.
package linecount

import (
	"bufio"
	"io"
	"sync"
)

// Counter tracks line counts and byte offsets for a stream.
type Counter struct {
	mu      sync.Mutex
	lines   int64
	bytes   int64
	offsets []int64 // byte offset at the start of each line
}

// New returns a new Counter.
func New() *Counter {
	return &Counter{
		offsets: []int64{0},
	}
}

// Count reads all bytes from r and records per-line offsets.
// It returns the total number of lines counted and any read error.
func (c *Counter) Count(r io.Reader) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	scanner := bufio.NewScanner(r)
	var offset int64
	for scanner.Scan() {
		line := scanner.Bytes()
		offset += int64(len(line)) + 1 // +1 for newline
		c.lines++
		c.offsets = append(c.offsets, offset)
	}
	c.bytes = offset
	return c.lines, scanner.Err()
}

// Lines returns the total number of lines recorded.
func (c *Counter) Lines() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lines
}

// Bytes returns the total number of bytes recorded.
func (c *Counter) Bytes() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.bytes
}

// OffsetOf returns the byte offset of the start of the given 1-based line
// number. Returns -1 if the line number is out of range.
func (c *Counter) OffsetOf(lineNum int64) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	if lineNum < 1 || int(lineNum) >= len(c.offsets) {
		return -1
	}
	return c.offsets[lineNum-1]
}

// Reset clears all recorded state.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lines = 0
	c.bytes = 0
	c.offsets = []int64{0}
}
