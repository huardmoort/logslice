// Package sampler provides line-rate sampling for large log streams.
// It supports both fixed-interval (every Nth line) and random sampling.
package sampler

import (
	"fmt"
	"math/rand"
)

// Mode describes the sampling strategy.
type Mode int

const (
	// ModeNth keeps every Nth line.
	ModeNth Mode = iota
	// ModeRandom keeps each line with probability 1/Rate.
	ModeRandom
)

// Sampler decides whether a given line should be kept.
type Sampler struct {
	mode    Mode
	rate    int
	counter int
	rng     *rand.Rand
}

// New creates a Sampler from a mode string and a rate.
// mode must be "nth" or "random"; rate must be >= 1.
func New(mode string, rate int, seed int64) (*Sampler, error) {
	if rate < 1 {
		return nil, fmt.Errorf("sampler: rate must be >= 1, got %d", rate)
	}
	var m Mode
	switch mode {
	case "nth":
		m = ModeNth
	case "random":
		m = ModeRandom
	default:
		return nil, fmt.Errorf("sampler: unknown mode %q, want \"nth\" or \"random\"", mode)
	}
	return &Sampler{
		mode: m,
		rate: rate,
		rng:  rand.New(rand.NewSource(seed)),
	}, nil
}

// Keep reports whether the current line should be included in output.
// It must be called exactly once per line in order.
func (s *Sampler) Keep() bool {
	s.counter++
	switch s.mode {
	case ModeNth:
		return s.counter%s.rate == 0
	case ModeRandom:
		return s.rng.Intn(s.rate) == 0
	}
	return true
}

// Reset resets the internal counter (useful between files).
func (s *Sampler) Reset() {
	s.counter = 0
}
