package seedmanager

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
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
	sm.seedPool[seed] = storedSeed{time.Now().Add(-expTime).Add(-1 * time.Second), 0}
	_, exists := sm.Pop(seed)
	require.False(t, exists)
	sm.Stop()
}

func TestSeedManager_ExpiredSeedAfterCleaning(t *testing.T) {
	expTime := time.Duration(1 * time.Minute)
	sm := NewSpecified(expTime, 1*time.Minute)
	seed := getSeed(t)
	sm.Add(seed, 0)
	sm.seedPool[seed] = storedSeed{time.Now().Add(-expTime).Add(-1 * time.Second), 0}
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
