// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package exporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/pulse"
)

type pulseStreamMock struct {
	checker func(*Pulse) error
}

func (p *pulseStreamMock) Send(pulse *Pulse) error {
	return p.checker(pulse)
}

func (p *pulseStreamMock) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (p *pulseStreamMock) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (p *pulseStreamMock) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (p *pulseStreamMock) Context() context.Context {
	return context.TODO()
}

func (p *pulseStreamMock) SendMsg(m interface{}) error {
	panic("implement me")
}

func (p *pulseStreamMock) RecvMsg(m interface{}) error {
	panic("implement me")
}

func TestPulseServer_Export(t *testing.T) {
	t.Run("fails if count is 0", func(t *testing.T) {
		server := NewPulseServer(nil, nil, nil)

		err := server.Export(&GetPulses{Count: 0}, &pulseStreamMock{})

		require.Equal(t, err, ErrNilCount)
	})

	t.Run("exporter works well. passed pulse is 0.", func(t *testing.T) {
		var pulses []insolar.PulseNumber
		pulseGatherer := func(p *Pulse) error {
			pulses = append(pulses, p.PulseNumber)
			return nil
		}
		stream := pulseStreamMock{checker: pulseGatherer}

		pulseCalculator := insolarPulse.NewCalculatorMock(t)
		pulseCalculator.ForwardsMock.When(context.TODO(), pulse.MinTimePulse, 1).Then(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 1}, nil)
		pulseCalculator.ForwardsMock.When(context.TODO(), pulse.MinTimePulse+1, 1).Then(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 2}, nil)

		jetKeeper := executor.NewJetKeeperMock(t)
		jetKeeper.TopSyncPulseMock.Return(pulse.MinTimePulse + 2)

		nodeList := []insolar.Node{{Role: insolar.StaticRoleLightMaterial}}
		nodeAccessor := node.NewAccessorMock(t)
		nodeAccessor.AllMock.When(pulse.MinTimePulse+1).Then(nodeList, nil)
		nodeAccessor.AllMock.When(pulse.MinTimePulse+2).Then(nodeList, nil)

		server := NewPulseServer(pulseCalculator, jetKeeper, nodeAccessor)

		err := server.Export(&GetPulses{PulseNumber: 0, Count: 10}, &stream)
		require.NoError(t, err)

		require.Equal(t, 3, len(pulses))
		require.Equal(t, pulse.MinTimePulse, int(pulses[0]))
		require.Equal(t, pulse.MinTimePulse+1, int(pulses[1]))
		require.Equal(t, pulse.MinTimePulse+2, int(pulses[2]))

		pulseCalculator.ForwardsMock.Return(insolar.Pulse{PulseNumber: 100}, nil)
		jetStorage := jet.NewStorageMock(t)
		jetKeeper.StorageMock.Return(jetStorage)
		jetStorage.AllMock.Return(nil)
		res2, err := server.NextFinalizedPulse(context.Background(), &GetNextFinalizedPulse{0})
		require.NoError(t, err)
		require.NotNil(t, res2)
		require.Equal(t, pulse.Number(100), res2.PulseNumber)

		res3, err := server.NextFinalizedPulse(context.Background(), &GetNextFinalizedPulse{pulse.MinTimePulse})
		require.NoError(t, err)
		require.NotNil(t, res3)
		require.Equal(t, pulse.Number(pulse.MinTimePulse+1), res3.PulseNumber)

	})

	t.Run("exporter works well. passed pulse is 0. read until top-sync", func(t *testing.T) {
		var pulses []insolar.PulseNumber
		pulseGatherer := func(p *Pulse) error {
			pulses = append(pulses, p.PulseNumber)
			return nil
		}
		stream := pulseStreamMock{checker: pulseGatherer}

		pulseCalculator := insolarPulse.NewCalculatorMock(t)
		pulseCalculator.ForwardsMock.When(context.TODO(), pulse.MinTimePulse, 1).Then(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 1}, nil)

		jetKeeper := executor.NewJetKeeperMock(t)
		jetKeeper.TopSyncPulseMock.Return(pulse.MinTimePulse + 1)

		nodeList := []insolar.Node{{Role: insolar.StaticRoleLightMaterial}}
		nodeAccessor := node.NewAccessorMock(t)
		nodeAccessor.AllMock.When(pulse.MinTimePulse+1).Then(nodeList, nil)

		server := NewPulseServer(pulseCalculator, jetKeeper, nodeAccessor)

		err := server.Export(&GetPulses{PulseNumber: 0, Count: 10}, &stream)
		require.NoError(t, err)

		require.Equal(t, 2, len(pulses))
		require.Equal(t, pulse.MinTimePulse, int(pulses[0]))
		require.Equal(t, pulse.MinTimePulse+1, int(pulses[1]))
	})

	t.Run("exporter works well. passed pulse is firstPulseNumber. read until top-sync", func(t *testing.T) {
		var pulses []insolar.PulseNumber
		pulseGatherer := func(p *Pulse) error {
			pulses = append(pulses, p.PulseNumber)
			return nil
		}
		stream := pulseStreamMock{checker: pulseGatherer}

		pulseCalculator := insolarPulse.NewCalculatorMock(t)
		pulseCalculator.ForwardsMock.When(context.TODO(), pulse.MinTimePulse, 1).Then(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 1}, nil)

		jetKeeper := executor.NewJetKeeperMock(t)
		jetKeeper.TopSyncPulseMock.Return(pulse.MinTimePulse + 1)

		nodeList := []insolar.Node{{Role: insolar.StaticRoleLightMaterial}}
		nodeAccessor := node.NewAccessorMock(t)
		nodeAccessor.AllMock.When(pulse.MinTimePulse+1).Then(nodeList, nil)

		server := NewPulseServer(pulseCalculator, jetKeeper, nodeAccessor)

		err := server.Export(&GetPulses{PulseNumber: pulse.MinTimePulse, Count: 10}, &stream)
		require.NoError(t, err)

		require.Equal(t, 1, len(pulses))
		require.Equal(t, pulse.MinTimePulse+1, int(pulses[0]))
	})

	t.Run("exporter works well. passed pulse is 0. read only 1", func(t *testing.T) {
		var pulses []insolar.PulseNumber
		pulseGatherer := func(p *Pulse) error {
			pulses = append(pulses, p.PulseNumber)
			return nil
		}
		stream := pulseStreamMock{checker: pulseGatherer}

		jetKeeper := executor.NewJetKeeperMock(t)
		jetKeeper.TopSyncPulseMock.Return(pulse.MinTimePulse + 2)

		nodeAccessor := node.NewAccessorMock(t)

		server := NewPulseServer(nil, jetKeeper, nodeAccessor)

		err := server.Export(&GetPulses{PulseNumber: 0, Count: 1}, &stream)
		require.NoError(t, err)

		require.Equal(t, 1, len(pulses))
		require.Equal(t, pulse.MinTimePulse, int(pulses[0]))
	})
}

func TestPulseServer_TopSyncPulse(t *testing.T) {
	pn := pulse.MinTimePulse + 2
	jetKeeper := executor.NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseMock.Return(insolar.PulseNumber(pn))
	nodeAccessor := node.NewAccessorMock(t)
	server := NewPulseServer(nil, jetKeeper, nodeAccessor)

	res, err := server.TopSyncPulse(context.Background(), nil)

	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, uint32(pn), res.PulseNumber)
}
