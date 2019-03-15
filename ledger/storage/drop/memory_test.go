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

	"github.com/google/gofuzz"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/stretchr/testify/require"
)

func TestNewStorageMemory(t *testing.T) {
	ms := NewStorageMemory()

	require.NotNil(t, ms.drops)
}

func TestDropStorageMemory_Set(t *testing.T) {
	ms := NewStorageMemory()

	var drops []jet.Drop
	genPulses := map[core.PulseNumber]struct{}{}
	f := fuzz.New().Funcs(func(jd *jet.Drop, c fuzz.Continue) {
		pn := gen.PulseNumber()
		genPulses[pn] = struct{}{}
		jd.Pulse = pn
	}).NumElements(5, 1000)
	f.Fuzz(&drops)

	genJets := map[core.JetID]struct{}{}
	for _, jd := range drops {
		j := gen.JetID()
		genJets[j] = struct{}{}
		err := ms.Set(inslogger.TestContext(t), j, jd)
		require.NoError(t, err)
	}

	require.Equal(t, len(drops), len(ms.drops))
	for k, jd := range ms.drops {
		_, ok := genPulses[jd.Pulse]
		require.Equal(t, true, ok)
		require.Equal(t, k.pulse, jd.Pulse)

		_, ok = genJets[k.jetID]
		require.Equal(t, true, ok)
	}
}

func TestDropStorageMemory_ForPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ms := NewStorageMemory()

	fJet := gen.JetID()
	fPn := gen.PulseNumber()
	_ = ms.Set(ctx, fJet, jet.Drop{Pulse: fPn})
	sJet := gen.JetID()
	sPn := gen.PulseNumber()
	_ = ms.Set(ctx, sJet, jet.Drop{Pulse: sPn})

	drop, err := ms.ForPulse(ctx, sJet, sPn)

	require.NoError(t, err)
	require.Equal(t, sPn, drop.Pulse)
	require.Equal(t, 2, len(ms.drops))
}

func TestDropStorageMemory_DoubleSet(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ms := NewStorageMemory()

	fJet := gen.JetID()
	fPn := gen.PulseNumber()
	fSize := rand.Uint64()
	sSize := rand.Uint64()

	_ = ms.Set(ctx, fJet, jet.Drop{Pulse: fPn, Size: fSize})
	_ = ms.Set(ctx, fJet, jet.Drop{Pulse: fPn, Size: sSize})

	drop, err := ms.ForPulse(ctx, fJet, fPn)

	require.NoError(t, err)
	require.Equal(t, fPn, drop.Pulse)
	require.Equal(t, sSize, drop.Size)
	require.Equal(t, 1, len(ms.drops))
}

func TestDropStorageDB_Set_Concurrent(t *testing.T) {
	ctx := inslogger.TestContext(t)
	var ms Modifier = NewStorageMemory()

	gonum := 10000
	startChannel := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(gonum)

	for i := 0; i < gonum; i++ {
		go func() {
			<-startChannel

			err := ms.Set(ctx, gen.JetID(), jet.Drop{Pulse: gen.PulseNumber(), Size: rand.Uint64()})
			require.NoError(t, err)

			wg.Done()
		}()
	}

	close(startChannel)
	wg.Wait()
}

func TestDropStorageMemory_Delete(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ms := NewStorageMemory()

	fJet := gen.JetID()
	sJet := gen.JetID()
	fPn := gen.PulseNumber()
	sPn := gen.PulseNumber()
	fSize := rand.Uint64()
	sSize := rand.Uint64()
	tSize := rand.Uint64()

	_ = ms.Set(ctx, fJet, jet.Drop{Pulse: fPn, Size: fSize})
	_ = ms.Set(ctx, sJet, jet.Drop{Pulse: fPn, Size: sSize})
	_ = ms.Set(ctx, fJet, jet.Drop{Pulse: sPn, Size: tSize})

	ms.Delete(fPn)

	drop, err := ms.ForPulse(ctx, fJet, sPn)
	require.NoError(t, err)
	require.Equal(t, drop.Pulse, sPn)
	require.Equal(t, drop.Size, tSize)

	drop, err = ms.ForPulse(ctx, fJet, fPn)
	require.Error(t, err, ErrNotFound)
	drop, err = ms.ForPulse(ctx, sJet, sPn)
	require.Error(t, err, ErrNotFound)
}
