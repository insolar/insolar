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
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
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
	sm.Stop()
}

func getSeed(t *testing.T) Seed {
	sg := SeedGenerator{}
	seed, err := sg.Next()
	require.NoError(t, err)
	return *seed
}

func TestSeedManager_Add(t *testing.T) {
	sm := NewSpecified(time.Duration(1*time.Minute), DefaultCleanPeriod)
	seed := getSeed(t)
	sm.Add(seed, 5)
	pulse, exists := sm.Pop(seed)
	require.True(t, exists)
	require.Equal(t, insolar.PulseNumber(5), pulse)
	sm.Stop()
}

func TestSeedManager_ExpiredSeed(t *testing.T) {
	expTime := time.Duration(1 * time.Minute)
	sm := NewSpecified(expTime, DefaultCleanPeriod)
	seed := getSeed(t)
	sm.Add(seed, 0)
	sm.seedPool[seed] = storedSeed{time.Now().UnixNano() - 1000, 0}
	_, exists := sm.Pop(seed)
	require.False(t, exists)
	sm.Stop()
}

func TestSeedManager_ExpiredSeedAfterCleaning(t *testing.T) {
	expTime := time.Duration(1 * time.Minute)
	sm := NewSpecified(expTime, 1*time.Minute)
	seed := getSeed(t)
	sm.Add(seed, 0)
	sm.seedPool[seed] = storedSeed{time.Now().UnixNano() - 1000, 0}
	_, exists := sm.Pop(seed)
	require.False(t, exists)
	sm.Stop()
}

func TestRace(t *testing.T) {
	const numConcurrent = 10

	expTime := time.Duration(2 * time.Millisecond)
	cleanPeriod := time.Duration(1 * time.Millisecond)
	sm := NewSpecified(expTime, cleanPeriod)

	wg := sync.WaitGroup{}
	wg.Add(numConcurrent)
	for i := 0; i < numConcurrent; i++ {
		go func() {
			defer wg.Done()
			var seeds []Seed
			numIterations := 300
			for j := 0; j < numIterations; j++ {
				seeds = append(seeds, getSeed(t))
				sm.Add(seeds[len(seeds)-1], 0)
			}
			<-time.After(cleanPeriod)
			for j := 0; j < numIterations; j++ {
				sm.Pop(seeds[j])
			}
		}()
	}
	wg.Wait()
	sm.Stop()
}
