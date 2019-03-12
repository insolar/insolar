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

package drop

import (
	"math/rand"
	"sync"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/stretchr/testify/require"
)

func TestNewStorageMemory(t *testing.T) {
	ms := NewStorageMemory()

	require.NotNil(t, ms.jets)
}

func TestDropStorageMemory_Set(t *testing.T) {
	ms := NewStorageMemory()

	var drops []jet.Drop
	dropsIndex := 0
	f := fuzz.New().Funcs(func(jd *jet.Drop, c fuzz.Continue) {
		dropsIndex++
		jd.Pulse = core.PulseNumber(dropsIndex)
	}).NumElements(5, 10)
	f.Fuzz(&drops)

	for _, jd := range drops {
		err := ms.Set(inslogger.TestContext(t), core.ZeroJetID, jd)
		require.NoError(t, err)
	}

	require.Equal(t, dropsIndex, len(ms.jets))
}

func TestDropStorageMemory_ForPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ms := NewStorageMemory()

	err := ms.Set(ctx, core.ZeroJetID, jet.Drop{Pulse: 123})
	require.NoError(t, err)
	err = ms.Set(ctx, *core.NewJetID(1, nil), jet.Drop{Pulse: 2})
	require.NoError(t, err)

	drop, err := ms.ForPulse(ctx, core.ZeroJetID, core.PulseNumber(123))

	require.Equal(t, core.PulseNumber(123), drop.Pulse)
	require.Equal(t, 2, len(ms.jets))
}

func TestDropStorageMemory_DoubleSet(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ms := NewStorageMemory()

	err := ms.Set(ctx, core.ZeroJetID, jet.Drop{Pulse: 123, Size: 123})
	require.NoError(t, err)
	err = ms.Set(ctx, core.ZeroJetID, jet.Drop{Pulse: 123, Size: 999})
	require.NoError(t, err)

	drop, err := ms.ForPulse(ctx, core.ZeroJetID, core.PulseNumber(123))

	require.Equal(t, core.PulseNumber(123), drop.Pulse)
	require.Equal(t, uint64(999), drop.Size)
	require.Equal(t, 1, len(ms.jets))
}

func TestDropStorageDB_Set_Concurrent(t *testing.T) {
	ctx := inslogger.TestContext(t)
	var ms Modifier = NewStorageMemory()

	gonum := 10000
	startChannel := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(gonum)

	for i := 0; i < gonum; i++ {
		go func(pulseNumber int) {
			<-startChannel

			err := ms.Set(ctx, gen.JetID(), jet.Drop{Pulse: core.PulseNumber(pulseNumber), Size: rand.Uint64()})
			require.NoError(t, err)

			wg.Done()
		}(i + 1)
	}

	close(startChannel)
	wg.Wait()
}

func TestDropStorageMemory_Delete(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ms := NewStorageMemory()

	err := ms.Set(ctx, core.ZeroJetID, jet.Drop{Pulse: 123, Size: 123})
	require.NoError(t, err)
	err = ms.Set(ctx, core.ZeroJetID, jet.Drop{Pulse: 222, Size: 999})
	require.NoError(t, err)

	ms.Delete(core.PulseNumber(123))

	drop, err := ms.ForPulse(ctx, core.ZeroJetID, core.PulseNumber(222))
	require.NoError(t, err)
	require.Equal(t, drop.Pulse, core.PulseNumber(222))
	require.Equal(t, drop.Size, uint64(999))

	drop, err = ms.ForPulse(ctx, core.ZeroJetID, core.PulseNumber(123))
	require.Error(t, err, db.ErrNotFound)
}
