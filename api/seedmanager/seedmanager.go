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

// Expiration represents time of expiration
type Expiration = int64

// DefaultTTL is default time period for deleting expired seeds
const DefaultTTL = time.Duration(1 * time.Second)

// DefaultCleanPeriod default time period for launching cleaning goroutine
const DefaultCleanPeriod = time.Duration(1 * time.Second)

// SeedManager manages working with seed pool
// It's thread safe
type SeedManager struct {
	mutex    sync.RWMutex
	seedPool map[Seed]Expiration
	ttl      time.Duration
}

// New creates new seed manager with default params
func New() *SeedManager {
	return NewSpecified(DefaultTTL, DefaultCleanPeriod)
}

// NewSpecified creates new seed manager with custom params
func NewSpecified(TTL time.Duration, cleanPeriod time.Duration) *SeedManager {
	sm := SeedManager{seedPool: make(map[Seed]Expiration), ttl: TTL}
	go func() {
		for range time.Tick(cleanPeriod) {
			sm.deleteExpired()
		}
	}()

	return &sm
}

// Add adds seed to pool
func (sm *SeedManager) Add(seed Seed) {
	expTime := time.Now().Add(sm.ttl).UnixNano()

	sm.mutex.Lock()
	sm.seedPool[seed] = expTime
	sm.mutex.Unlock()

}

func (sm *SeedManager) isExpired(expTime Expiration) bool {
	return expTime < time.Now().UnixNano()
}

// Exists checks whether seed in the pool
func (sm *SeedManager) Exists(seed Seed) bool {
	sm.mutex.RLock()
	expTime, ok := sm.seedPool[seed]
	sm.mutex.RUnlock()

	isSeedOk := ok && !sm.isExpired(expTime)
	if isSeedOk {
		sm.mutex.Lock()
		delete(sm.seedPool, seed)
		sm.mutex.Unlock()
	}

	return isSeedOk
}

func (sm *SeedManager) deleteExpired() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	for seed, expTime := range sm.seedPool {
		if sm.isExpired(expTime) {
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
