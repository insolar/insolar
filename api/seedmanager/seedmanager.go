/*
 *    Copyright 2018 Insolar
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

package seedmanager

import (
	"sync"
	"time"
)

type Expiration = int64

const TTL = time.Duration(1 * time.Second)

// SeedManager manages working with seed pool
// It's thread safe
type SeedManager struct {
	mu       sync.RWMutex
	seedPool map[Seed]Expiration
}

// New creates new seed manager
func New() *SeedManager {
	sm := SeedManager{seedPool: make(map[Seed]Expiration)}
	go func() {
		for range time.Tick(time.Second) {
			sm.deleteExpired()
		}
	}()

	return &sm
}

// Add adds seed to pool
func (sm *SeedManager) Add(seed Seed) {
	expTime := time.Now().Add(TTL).UnixNano()

	sm.mu.Lock()
	sm.seedPool[seed] = expTime
	sm.mu.Unlock()

}

func (sm *SeedManager) isExpired(expTime Expiration) bool {
	return expTime < time.Now().UnixNano()
}

// Exists checks whether seed in the pool
func (sm *SeedManager) Exists(seed Seed) bool {
	sm.mu.RLock()
	expTime, ok := sm.seedPool[seed]
	sm.mu.RUnlock()

	return ok && !sm.isExpired(expTime)
}

func (sm *SeedManager) deleteExpired() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for seed, expTime := range sm.seedPool {
		if sm.isExpired(expTime) {
			delete(sm.seedPool, seed)
		}
	}
}
