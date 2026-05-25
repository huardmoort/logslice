// Package window provides a sliding window buffer that retains the last N lines
// seen during log processing. It is useful for capturing context lines around
// matched entries (e.g. --before-context / --after-context style behaviour).
package window

// Window is a fixed-capacity circular buffer of strings.
type Window struct {
	buf  []string
	size int
	head int
	count int
}

// New creates a Window that retains up to size lines.
// If size is <= 0 the window is a no-op (Len always returns 0).
func New(size int) *Window {
	if size <= 0 {
		return &Window{}
	}
	return &Window{
		buf:  make([]string, size),
		size: size,
	}
}

// Add appends a line to the window, evicting the oldest entry when full.
func (w *Window) Add(line string) {
	if w.size == 0 {
		return
	}
	w.buf[w.head] = line
	w.head = (w.head + 1) % w.size
	if w.count < w.size {
		w.count++
	}
}

// Lines returns the buffered lines in insertion order (oldest first).
// The returned slice is a copy and is safe to modify.
func (w *Window) Lines() []string {
	if w.count == 0 {
		return nil
	}
	out := make([]string, w.count)
	start := (w.head - w.count + w.size) % w.size
	for i := 0; i < w.count; i++ {
		out[i] = w.buf[(start+i)%w.size]
	}
	return out
}

// Len returns the number of lines currently held in the window.
func (w *Window) Len() int { return w.count }

// Reset clears all buffered lines without changing the capacity.
func (w *Window) Reset() {
	w.head = 0
	w.count = 0
}
