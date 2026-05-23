// Package dedupe provides line deduplication for log output,
// suppressing consecutive or windowed repeated lines.
package dedupe

import (
	"crypto/sha256"
	"encoding/hex"
)

// Mode controls how deduplication is applied.
type Mode int

const (
	// ModeConsecutive suppresses a line only if it is identical to the immediately preceding line.
	ModeConsecutive Mode = iota
	// ModeWindow suppresses a line if it appeared anywhere within the last N lines.
	ModeWindow
)

// Deduper filters duplicate log lines.
type Deduper struct {
	mode       Mode
	windowSize int
	lastHash   string
	window     []string // circular buffer of hashes
	pos        int
	Suppressed int
}

// New creates a Deduper. For ModeConsecutive windowSize is ignored.
// For ModeWindow, windowSize must be >= 1.
func New(mode Mode, windowSize int) (*Deduper, error) {
	if mode == ModeWindow && windowSize < 1 {
		return nil, fmt.Errorf("dedupe: window size must be >= 1, got %d", windowSize)
	}
	d := &Deduper{mode: mode, windowSize: windowSize}
	if mode == ModeWindow {
		d.window = make([]string, windowSize)
	}
	return d, nil
}

// Allow returns true if the line should be forwarded, false if it is a duplicate.
func (d *Deduper) Allow(line string) bool {
	h := hash(line)
	switch d.mode {
	case ModeConsecutive:
		if h == d.lastHash {
			d.Suppressed++
			return false
		}
		d.lastHash = h
		return true
	case ModeWindow:
		for _, wh := range d.window {
			if wh == h {
				d.Suppressed++
				return false
			}
		}
		d.window[d.pos%d.windowSize] = h
		d.pos++
		return true
	}
	return true
}

func hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:8])
}
