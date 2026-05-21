// Package tail provides functionality to follow a log file as new lines
// are appended, similar to `tail -f`.
package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Options configures the tail follower.
type Options struct {
	// PollInterval is how often to check for new data when the reader
	// reaches EOF. Defaults to 200ms if zero.
	PollInterval time.Duration
}

// Follower reads new lines appended to a file, sending each line to Lines.
type Follower struct {
	Lines  <-chan string
	Errors <-chan error
	close  func()
}

// Close stops the follower.
func (f *Follower) Close() { f.close() }

// Follow opens the named file and emits new lines as they are appended.
// The follower runs until ctx is cancelled or Close is called.
func Follow(ctx context.Context, filename string, opts Options) (*Follower, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	if opts.PollInterval == 0 {
		opts.PollInterval = 200 * time.Millisecond
	}

	lines := make(chan string, 64)
	errors := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer close(lines)
		defer close(errors)
		defer f.Close()
		r := bufio.NewReader(f)
		for {
			line, err := r.ReadString('\n')
			if err == nil {
				select {
				case lines <- line:
				case <-ctx.Done():
					return
				}
				continue
			}
			if err != io.EOF {
				errors <- err
				return
			}
			// EOF: wait for more data
			select {
			case <-ctx.Done():
				return
			case <-time.After(opts.PollInterval):
			}
		}
	}()

	return &Follower{
		Lines:  lines,
		Errors: errors,
		close:  cancel,
	}, nil
}
