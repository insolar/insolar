// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package drop

import (
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

	err := ms.Set(ctx, Drop{JetID: fJet, Pulse: fPn})
	require.NoError(t, err)
	err = ms.Set(ctx, Drop{JetID: fJet, Pulse: fPn})
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

			err := ms.Set(ctx, Drop{JetID: gen.JetID(), Pulse: gen.PulseNumber()})
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

	_ = ms.Set(ctx, Drop{JetID: jets[0], Pulse: fPn})
	_ = ms.Set(ctx, Drop{JetID: jets[0], Pulse: sPn})
	_ = ms.Set(ctx, Drop{JetID: jets[1], Pulse: fPn})

	ms.DeleteForPN(ctx, fPn)

	drop, err := ms.ForPulse(ctx, jets[0], sPn)
	require.NoError(t, err)
	require.Equal(t, drop.Pulse, sPn)

	drop, err = ms.ForPulse(ctx, jets[0], fPn)
	require.Error(t, err, ErrNotFound)
	drop, err = ms.ForPulse(ctx, jets[1], sPn)
	require.Error(t, err, ErrNotFound)
}
