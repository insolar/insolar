/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

// Package backoff provides an exponential-backoff implementation.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Backoff is a time.Duration counter, starting at Min. After every call to
// the Duration method the current timing is multiplied by Factor, but it
// never exceeds Max.
//
// Backoff is not generally concurrent-safe, but the ForAttempt method can
// be used concurrently.
type Backoff struct {
	attempt int
	//Factor is the multiplying factor for each increment step
	Factor float64
	//Jitter eases contention by randomizing backoff steps
	Jitter bool
	//Min and Max are the minimum and maximum values of the counter
	Min, Max time.Duration
}

// Duration returns the duration for the current attempt before incrementing
// the attempt counter. See ForAttempt.
func (b *Backoff) Duration() time.Duration {
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
func (b *Backoff) ForAttempt(attempt int) time.Duration {
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
	//calculate this duration
	minf := float64(min)
	durf := minf * math.Pow(factor, float64(attempt))
	if b.Jitter {
		durf = rand.Float64()*(durf-minf) + minf
	}
	//ensure float64 wont overflow int64
	if durf > maxInt64 {
		return max
	}
	dur := time.Duration(durf)
	//keep within bounds
	if dur < min {
		return min
	}
	if dur > max {
		return max
	}
	return dur
}

// Reset restarts the current attempt counter at zero.
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Attempt returns the current attempt counter value.
func (b *Backoff) Attempt() int {
	return b.attempt
}

// Copy returns a backoff with equals constraints as the original
func (b *Backoff) Copy() *Backoff {
	return &Backoff{
		Factor: b.Factor,
		Jitter: b.Jitter,
		Min:    b.Min,
		Max:    b.Max,
	}
}
