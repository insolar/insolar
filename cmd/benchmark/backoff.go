package main

import (
	"math"
	"time"
)

// backoff is a time.Duration counter, starting at Min. After every call to
// the Duration method the current timing is multiplied by Factor, but it
// never exceeds Max.
type backoff struct {
	// Factor is the multiplying factor for each increment step
	attempt, Factor float64
	// Min and Max are the minimum and maximum values of the counter
	Min, Max time.Duration
}

// Duration returns the duration for the current attempt before incrementing
// the attempt counter. See ForAttempt.
func (b *backoff) Duration() time.Duration {
	d := b.ForAttempt(b.attempt)
	b.attempt++
	return d
}

const maxInt64 = float64(math.MaxInt64 - 512)

// ForAttempt returns the duration for a specific attempt. This is useful if
// you have a large number of independent Backoffs, but don't want use
// unnecessary memory storing the Backoff parameters per Backoff. The first
// attempt should be 0.
//
// ForAttempt is concurrent-safe.
func (b *backoff) ForAttempt(attempt float64) time.Duration {
	// Zero-values are nonsensical, so we use
	// them to apply defaults
	min := b.Min
	if min <= 0 {
		min = 100 * time.Millisecond
	}
	max := b.Max
	if max <= 0 {
		max = 10 * time.Second
	}
	if min >= max {
		// short-circuit
		return max
	}
	factor := b.Factor
	if factor <= 0 {
		factor = 2
	}
	// calculate this duration
	minf := float64(min)
	durf := minf * math.Pow(factor, attempt)
	// ensure float64 wont overflow int64
	if durf > maxInt64 {
		return max
	}
	dur := time.Duration(durf)
	// keep within bounds
	if dur < min {
		return min
	}
	if dur > max {
		return max
	}
	return dur
}
