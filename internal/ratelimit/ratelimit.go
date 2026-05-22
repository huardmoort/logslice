// Package ratelimit provides a token-bucket rate limiter for controlling
// how many log lines are emitted per second during tail or pipeline operations.
package ratelimit

import (
	"errors"
	"time"
)

// Limiter controls the rate at which lines are processed.
type Limiter struct {
	tokens   float64
	max      float64
	rate     float64 // tokens per second
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows up to ratePerSec lines per second.
// ratePerSec must be greater than zero.
func New(ratePerSec float64) (*Limiter, error) {
	if ratePerSec <= 0 {
		return nil, errors.New("ratelimit: ratePerSec must be > 0")
	}
	now := time.Now()
	return &Limiter{
		tokens:   ratePerSec,
		max:      ratePerSec,
		rate:     ratePerSec,
		lastTick: now,
		clock:    time.Now,
	}, nil
}

// Allow returns true if a line may be processed right now, consuming one token.
// It refills tokens based on elapsed time since the last call.
func (l *Limiter) Allow() bool {
	now := l.clock()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.lastTick = now

	l.tokens += elapsed * l.rate
	if l.tokens > l.max {
		l.tokens = l.max
	}

	if l.tokens >= 1.0 {
		l.tokens -= 1.0
		return true
	}
	return false
}

// Wait blocks until a token is available, then consumes it.
func (l *Limiter) Wait() {
	for {
		if l.Allow() {
			return
		}
		// Sleep for roughly the time needed to accumulate one token.
		sleepDur := time.Duration(float64(time.Second) / l.rate)
		time.Sleep(sleepDur / 2)
	}
}

// SetRate updates the rate limit in tokens per second.
func (l *Limiter) SetRate(ratePerSec float64) error {
	if ratePerSec <= 0 {
		return errors.New("ratelimit: ratePerSec must be > 0")
	}
	l.rate = ratePerSec
	l.max = ratePerSec
	return nil
}
