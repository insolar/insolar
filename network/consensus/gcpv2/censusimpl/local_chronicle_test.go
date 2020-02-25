// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package censusimpl

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/pulse"
)

func TestNewLocalChronicles(t *testing.T) {
	pf := profiles.NewFactoryMock(t)
	lcc := NewLocalChronicles(pf)
	ksf := cryptkit.NewKeyStoreFactoryMock(t)
	require.Equal(t, pf, lcc.GetProfileFactory(ksf))
}

func TestGetLatestCensus(t *testing.T) {
	lc := localChronicles{}
	latestCensus, _ := lc.GetLatestCensus()
	require.Nil(t, latestCensus)

	exp := census.NewExpectedMock(t)
	lc.expected = exp
	latestCensus, _ = lc.GetLatestCensus()
	require.Equal(t, exp, latestCensus)

	lc.expected = nil
	pct := &PrimingCensusTemplate{CensusTemplate{pd: pulse.Data{PulseNumber: 1}}}
	lc.active = pct
	latestCensus, _ = lc.GetLatestCensus()
	require.Equal(t, pct, latestCensus)
}

func TestGetRecentCensus(t *testing.T) {
	lc := localChronicles{}
	exp := census.NewExpectedMock(t)
	pn := pulse.Number(pulse.MinTimePulse)
	exp.GetPulseNumberMock.Set(func() pulse.Number { return pn })
	lc.expected = exp

	require.Equal(t, exp, lc.GetRecentCensus(pn))

	require.Panics(t, func() { lc.GetRecentCensus(pn + 1) })

	lc.expected = nil
	require.Panics(t, func() { lc.GetRecentCensus(pn) })

	pct := &PrimingCensusTemplate{CensusTemplate{pd: pulse.Data{PulseNumber: pn}}}
	lc.active = pct
	require.Equal(t, pct, lc.GetRecentCensus(pn))

	require.Panics(t, func() { lc.GetRecentCensus(pn + 1) })
}

func TestGetActiveCensus(t *testing.T) {
	lc := localChronicles{}
	pct := &PrimingCensusTemplate{CensusTemplate{pd: pulse.Data{PulseNumber: 1}}}
	lc.active = pct
	require.Equal(t, pct, lc.GetActiveCensus())
}

func TestGetExpectedCensus(t *testing.T) {
	lc := localChronicles{}
	exp := census.NewExpectedMock(t)
	lc.expected = exp
	require.Equal(t, exp, lc.GetExpectedCensus())
}

func TestLCMakeActive(t *testing.T) {
	t.Skip("merge")
	// lc := localChronicles{}
	// exp1 := census.NewExpectedMock(t)
	// lc.expected = exp1
	// pn := pulse.Number(pulse.MinTimePulse)
	// pct := &PrimingCensusTemplate{CensusTemplate{pd: pulse.Data{PulseNumber: pn}}}
	// lc.active = pct
	// exp2 := census.NewExpectedMock(t)
	// require.Panics(t, func() { lc.makeActive(exp2, pct) })
	//
	// lc.expectedPulseNumber = pn + 1
	// require.Panics(t, func() { lc.makeActive(exp1, pct) })
	//
	// lc.expectedPulseNumber = pn
	// pct.pd.PulseNumber = 1
	// require.Panics(t, func() { lc.makeActive(exp1, pct) })
	//
	// lc.expectedPulseNumber = pulse.Unknown
	// require.Panics(t, func() { lc.makeActive(exp1, pct) })
	//
	// lc.expectedPulseNumber = pn
	// pct.pd.PulseNumber = pn
	// require.Panics(t, func() { lc.makeActive(exp1, pct) })
	//
	// pct.pd.PulseEpoch = pulse.MaxTimePulse + 1
	// require.Panics(t, func() { lc.makeActive(exp1, pct) })
	//
	// pct.pd.PulseEpoch = pulse.MaxTimePulse
	// pct.pd.NextPulseDelta = 1
	// registries := census.NewVersionedRegistriesMock(t)
	// vr := census.NewVersionedRegistriesMock(t)
	// registries.CommitNextPulseMock.Set(
	// 	func(pulse.Data, census.OnlinePopulation) census.VersionedRegistries { return vr })
	// pct.registries = registries
	// lc.makeActive(exp1, pct)
	// require.Equal(t, pn+pulse.Number(pct.pd.NextPulseDelta), lc.expectedPulseNumber)
	//
	// require.Equal(t, pct, lc.active)
	//
	// require.Nil(t, lc.expected)
	//
	// lc = localChronicles{expected: exp1}
	// pct = &PrimingCensusTemplate{}
	// require.Panics(t, func() { lc.makeActive(exp1, pct) })
	//
	// pct.registries = registries
	// pct.pd = pulse.Data{PulseNumber: pn, DataExt: pulse.DataExt{NextPulseDelta: 0}}
	// lc.makeActive(exp1, pct)
	// require.Equal(t, pn, lc.expectedPulseNumber)
	//
	// require.Equal(t, pct, lc.active)
	//
	// require.Nil(t, lc.expected)
	//
	// lc = localChronicles{expected: exp1}
	// pct = &PrimingCensusTemplate{CensusTemplate{registries: registries, pd: pulse.Data{PulseNumber: pn,
	// 	DataExt: pulse.DataExt{NextPulseDelta: 1, PulseEpoch: pulse.MaxTimePulse}}}}
	// lc.makeActive(exp1, pct)
	// require.Equal(t, pn+pulse.Number(pct.pd.NextPulseDelta), lc.expectedPulseNumber)
	//
	// require.Equal(t, pct, lc.active)
	//
	// require.Nil(t, lc.expected)
	//
	// lc = localChronicles{expected: exp1}
	// pct = &PrimingCensusTemplate{CensusTemplate{registries: registries, pd: pulse.Data{PulseNumber: pn,
	// 	DataExt: pulse.DataExt{NextPulseDelta: 1, PulseEpoch: pulse.MaxTimePulse + 1}}}}
	// lc.makeActive(exp1, pct)
	// require.Equal(t, pulse.Unknown, lc.expectedPulseNumber)
	//
	// require.Equal(t, pct, lc.active)
	//
	// require.Nil(t, lc.expected)
}

func TestLCMakeExpected(t *testing.T) {
	lc := localChronicles{}
	exp := census.NewExpectedMock(t)
	exp.GetPreviousMock.Set(func() (a1 census.Active) {
		return census.NewActiveMock(t)
	})
	exp.GetOnlinePopulationMock.Set(func() (o1 census.OnlinePopulation) {
		return nil
	})

	lc.expected = exp
	require.Panics(t, func() { lc.makeExpected(exp) })

	lc.expected = nil
	pn := pulse.Number(pulse.MinTimePulse)
	exp.GetPulseNumberMock.Set(func() pulse.Number { return pn })
	exp.GetExpectedPulseNumberMock.Set(func() (n1 pulse.Number) {
		return pulse.MaxTimePulse
	})

	require.Panics(t, func() { lc.makeExpected(exp) })

	exp.GetExpectedPulseNumberMock.Set(func() (n1 pulse.Number) {
		return pulse.Unknown
	})
	exp.GetPreviousMock.Set(func() (a1 census.Active) {
		return nil
	})
	lc.makeExpected(exp)
	require.Equal(t, exp, lc.expected)
}

func TestGetProfileFactory(t *testing.T) {
	pf := profiles.NewFactoryMock(t)
	lc := localChronicles{profileFactory: pf}
	require.Equal(t, pf, lc.GetProfileFactory(cryptkit.NewKeyStoreFactoryMock(t)))
}
