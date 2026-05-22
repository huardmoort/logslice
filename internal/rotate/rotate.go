// Package rotate provides log rotation detection and handling for logslice.
// It monitors log files for rotation events (truncation or inode change)
// and allows consumers to react accordingly.
package rotate

import (
	"errors"
	"os"
	"time"
)

// Event describes why a rotation was detected.
type Event int

const (
	EventTruncated Event = iota // file size decreased
	EventReplaced               // inode changed (file replaced)
)

func (e Event) String() string {
	switch e {
	case EventTruncated:
		return "truncated"
	case EventReplaced:
		return "replaced"
	default:
		return "unknown"
	}
}

// Detector watches a file for rotation.
type Detector struct {
	path    string
	size    int64
	inode   uint64
	pollInt time.Duration
}

// New creates a Detector for the given path with the provided poll interval.
// Returns an error if the file cannot be stat'd initially.
func New(path string, pollInterval time.Duration) (*Detector, error) {
	if pollInterval <= 0 {
		return nil, errors.New("rotate: poll interval must be positive")
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	d := &Detector{
		path:    path,
		size:    info.Size(),
		inode:   inode(info),
		pollInt: pollInterval,
	}
	return d, nil
}

// Check tests whether the file has been rotated since the last call to Check
// or since the Detector was created. Returns (event, true) on rotation.
func (d *Detector) Check() (Event, bool) {
	info, err := os.Stat(d.path)
	if err != nil {
		// File gone — treat as replaced.
		return EventReplaced, true
	}
	newInode := inode(info)
	newSize := info.Size()

	if newInode != d.inode {
		d.inode = newInode
		d.size = newSize
		return EventReplaced, true
	}
	if newSize < d.size {
		d.size = newSize
		return EventTruncated, true
	}
	d.size = newSize
	return 0, false
}

// Watch blocks and sends rotation Events on the returned channel until ctx is
// done or the channel is drained and closed by the caller.
func (d *Detector) Watch(stop <-chan struct{}) <-chan Event {
	ch := make(chan Event, 1)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(d.pollInt)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				if ev, rotated := d.Check(); rotated {
					select {
					case ch <- ev:
					case <-stop:
						return
					}
				}
			}
		}
	}()
	return ch
}
