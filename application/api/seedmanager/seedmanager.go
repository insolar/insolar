//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package seedmanager

import (
	"sync"
	"time"

	"github.com/insolar/insolar/insolar"
)

// DefaultTTL is default time period for deleting expired seeds
const DefaultTTL = 600 * time.Second

// DefaultCleanPeriod default time period for launching cleaning goroutine
const DefaultCleanPeriod = 60 * time.Second

type storedSeed struct {
	ts    time.Time
	pulse insolar.PulseNumber
}

// SeedManager manages working with seed pool
// It's thread safe
type SeedManager struct {
	mutex    sync.Mutex
	seedPool map[Seed]storedSeed
	ttl      time.Duration
	stopped  chan struct{}
}

// New creates new seed manager with default params
func New() *SeedManager {
	return NewSpecified(DefaultTTL, DefaultCleanPeriod)
}

// NewSpecified creates new seed manager with custom params
func NewSpecified(ttl time.Duration, cleanPeriod time.Duration) *SeedManager {
	sm := SeedManager{
		seedPool: make(map[Seed]storedSeed),
		ttl:      ttl,
		stopped:  make(chan struct{}),
	}

	ticker := time.NewTicker(cleanPeriod)
	defer ticker.Stop()

	go func() {
		var stop = false
		for !stop {
			select {
			case <-ticker.C:
				sm.deleteExpired()
			case <-sm.stopped:
				stop = true
			}
		}
		sm.stopped <- struct{}{}
	}()

	return &sm
}

func (sm *SeedManager) Stop() {
	sm.stopped <- struct{}{}
	<-sm.stopped
}

// Add adds seed to pool
func (sm *SeedManager) Add(seed Seed, pulse insolar.PulseNumber) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	sm.seedPool[seed] = storedSeed{
		ts:    time.Now(),
		pulse: pulse,
	}

}

func (sm *SeedManager) isExpired(ts time.Time) bool {
	return time.Now().Sub(ts) > sm.ttl
}

// Pop deletes and returns seed from the pool
func (sm *SeedManager) Pop(seed Seed) (insolar.PulseNumber, bool) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	stored, ok := sm.seedPool[seed]

	if ok && !sm.isExpired(stored.ts) {
		delete(sm.seedPool, seed)
		return stored.pulse, true
	}

	return 0, false
}

func (sm *SeedManager) deleteExpired() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for seed, stored := range sm.seedPool {
		if sm.isExpired(stored.ts) {
			delete(sm.seedPool, seed)
		}
	}
}

// SeedFromBytes converts slice of bytes to Seed. Returns nil if slice's size is not equal to SeedSize
func SeedFromBytes(slice []byte) *Seed {
	if len(slice) != int(SeedSize) {
		return nil
	}
	var result Seed
	copy(result[:], slice[:SeedSize])
	return &result
}
