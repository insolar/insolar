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

package drop

import (
	"math/rand"
	"sync"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func TestNewStorageMemory(t *testing.T) {
	ms := NewStorageMemory()

	require.NotNil(t, ms.drops)
}

func TestDropStorageMemory_Set(t *testing.T) {
	ms := NewStorageMemory()

	var drops []Drop
	genPulses := map[insolar.PulseNumber]struct{}{}
	genJets := map[insolar.JetID]struct{}{}

	f := fuzz.New().Funcs(func(jd *Drop, c fuzz.Continue) {
		pn := gen.PulseNumber()
		genPulses[pn] = struct{}{}
		jd.Pulse = pn

		j := gen.JetID()
		genJets[j] = struct{}{}
		jd.JetID = j
	}).NumElements(5, 1000)
	f.Fuzz(&drops)

	for _, jd := range drops {
		err := ms.Set(inslogger.TestContext(t), jd)
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

	jets := gen.UniqueJetIDs(2)

	fPn := gen.PulseNumber()
	_ = ms.Set(ctx, Drop{JetID: jets[0], Pulse: fPn})
	sPn := gen.PulseNumber()
	_ = ms.Set(ctx, Drop{JetID: jets[1], Pulse: sPn})

	drop, err := ms.ForPulse(ctx, jets[1], sPn)

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

	err := ms.Set(ctx, Drop{JetID: fJet, Pulse: fPn, Size: fSize})
	require.NoError(t, err)
	err = ms.Set(ctx, Drop{JetID: fJet, Pulse: fPn, Size: sSize})
	require.Error(t, err, ErrOverride)
}

func TestDropStorageMemory_Set_Concurrent(t *testing.T) {
	ctx := inslogger.TestContext(t)
	var ms Modifier = NewStorageMemory()

	gonum := 50
	startChannel := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(gonum)

	for i := 0; i < gonum; i++ {
		go func() {
			<-startChannel

			err := ms.Set(ctx, Drop{JetID: gen.JetID(), Pulse: gen.PulseNumber(), Size: rand.Uint64()})
			if err != nil {
				require.Error(t, err, ErrOverride)
			}

			wg.Done()
		}()
	}

	close(startChannel)
	wg.Wait()
}

func TestDropStorageMemory_Delete(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ms := NewStorageMemory()

	jets := gen.UniqueJetIDs(2)

	fPn := gen.PulseNumber()
	sPn := gen.PulseNumber()
	fSize := rand.Uint64()
	sSize := rand.Uint64()
	tSize := rand.Uint64()

	_ = ms.Set(ctx, Drop{JetID: jets[0], Pulse: fPn, Size: fSize})
	_ = ms.Set(ctx, Drop{JetID: jets[0], Pulse: sPn, Size: sSize})
	_ = ms.Set(ctx, Drop{JetID: jets[1], Pulse: fPn, Size: tSize})

	ms.DeleteForPN(ctx, fPn)

	drop, err := ms.ForPulse(ctx, jets[0], sPn)
	require.NoError(t, err)
	require.Equal(t, drop.Pulse, sPn)
	require.Equal(t, drop.Size, sSize)

	drop, err = ms.ForPulse(ctx, jets[0], fPn)
	require.Error(t, err, ErrNotFound)
	drop, err = ms.ForPulse(ctx, jets[1], sPn)
	require.Error(t, err, ErrNotFound)
}
