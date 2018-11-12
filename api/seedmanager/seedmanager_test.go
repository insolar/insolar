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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSeedFromBytes_BadInputSize(t *testing.T) {
	var badSeedBytes []byte
	res := SeedFromBytes(badSeedBytes)
	require.Nil(t, res)
}

func TestSeedFromBytes(t *testing.T) {
	var seedBytes []byte
	for i := byte(0); i < byte(SeedSize); i++ {
		seedBytes = append(seedBytes, 'A'+i)
	}
	res := SeedFromBytes(seedBytes)
	require.NotNil(t, res)
	require.Equal(t, seedBytes, res[:])
}

func TestNew(t *testing.T) {
	sm := New()
	require.Empty(t, sm.seedPool)
}

func getSeed(t *testing.T) Seed {
	sg := SeedGenerator{}
	seed, err := sg.Next()
	require.NoError(t, err)
	return *seed
}

func TestSeedManager_Add(t *testing.T) {
	sm := NewSpecified(time.Duration(5*time.Millisecond), DefaultCleanPeriod)
	seed := getSeed(t)
	sm.Add(seed)
	require.True(t, sm.Exists(seed))
}

func TestSeedManager_ExpiredSeed(t *testing.T) {
	expTime := time.Duration(5 * time.Millisecond)
	sm := NewSpecified(expTime, DefaultCleanPeriod)
	seed := getSeed(t)
	sm.Add(seed)
	<-time.After(expTime * 2)
	require.False(t, sm.Exists(seed))
}

func TestSeedManager_ExistsThanExpiredSeed(t *testing.T) {
	seed := getSeed(t)
	ttl := time.Duration(8 * time.Millisecond)
	sm := NewSpecified(ttl, DefaultCleanPeriod)
	sm.Add(seed)
	require.True(t, sm.Exists(seed))
	<-time.After(ttl * 2)
	require.False(t, sm.Exists(seed))
}

func TestSeedManager_ExpiredSeedAfterCleaning(t *testing.T) {
	expTime := time.Duration(2 * time.Millisecond)
	sm := NewSpecified(expTime, 2*time.Millisecond)
	seed := getSeed(t)
	sm.Add(seed)
	<-time.After(8 * time.Millisecond)
	require.False(t, sm.Exists(seed))
}

func TestRace(t *testing.T) {
	const numConcurrent = 15

	expTime := time.Duration(2 * time.Millisecond)
	cleanPeriod := time.Duration(1 * time.Millisecond)
	sm := NewSpecified(expTime, cleanPeriod)

	wg := sync.WaitGroup{}
	wg.Add(numConcurrent)
	for i := 0; i < numConcurrent; i++ {
		go func() {
			defer wg.Done()
			var seeds []Seed
			for j := 0; j < 500; j++ {
				seeds = append(seeds, getSeed(t))
				sm.Add(seeds[len(seeds)-1])
			}
			<-time.After(cleanPeriod)
			for j := 0; j < 500; j++ {
				sm.Exists(seeds[j])
			}
		}()
	}
	wg.Wait()
}
