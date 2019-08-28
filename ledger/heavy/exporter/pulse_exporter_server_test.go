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

package exporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils/network"
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
		server := NewPulseServer(nil, nil)

		err := server.Export(&GetPulses{Count: 0}, &pulseStreamMock{})

		require.Error(t, err)
	})

	t.Run("exporter works well. passed pulse is 0.", func(t *testing.T) {
		var pulses []insolar.PulseNumber
		pulseGatherer := func(p *Pulse) error {
			pulses = append(pulses, p.PulseNumber)
			return nil
		}
		stream := pulseStreamMock{checker: pulseGatherer}

		pulseCalculator := network.NewPulseCalculatorMock(t)
		pulseCalculator.ForwardsMock.When(context.TODO(), pulse.MinTimePulse, 1).Then(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 1}, nil)
		pulseCalculator.ForwardsMock.When(context.TODO(), pulse.MinTimePulse+1, 1).Then(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 2}, nil)

		jetKeeper := executor.NewJetKeeperMock(t)
		jetKeeper.TopSyncPulseMock.Return(pulse.MinTimePulse + 2)

		server := NewPulseServer(pulseCalculator, jetKeeper)

		err := server.Export(&GetPulses{PulseNumber: 0, Count: 10}, &stream)
		require.NoError(t, err)

		require.Equal(t, 3, len(pulses))
		require.Equal(t, pulse.MinTimePulse, int(pulses[0]))
		require.Equal(t, pulse.MinTimePulse+1, int(pulses[1]))
		require.Equal(t, pulse.MinTimePulse+2, int(pulses[2]))
	})

	t.Run("exporter works well. passed pulse is 0. read until top-sync", func(t *testing.T) {
		var pulses []insolar.PulseNumber
		pulseGatherer := func(p *Pulse) error {
			pulses = append(pulses, p.PulseNumber)
			return nil
		}
		stream := pulseStreamMock{checker: pulseGatherer}

		pulseCalculator := network.NewPulseCalculatorMock(t)
		pulseCalculator.ForwardsMock.When(context.TODO(), pulse.MinTimePulse, 1).Then(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 1}, nil)

		jetKeeper := executor.NewJetKeeperMock(t)
		jetKeeper.TopSyncPulseMock.Return(pulse.MinTimePulse + 1)

		server := NewPulseServer(pulseCalculator, jetKeeper)

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

		pulseCalculator := network.NewPulseCalculatorMock(t)
		pulseCalculator.ForwardsMock.When(context.TODO(), pulse.MinTimePulse, 1).Then(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 1}, nil)

		jetKeeper := executor.NewJetKeeperMock(t)
		jetKeeper.TopSyncPulseMock.Return(pulse.MinTimePulse + 1)

		server := NewPulseServer(pulseCalculator, jetKeeper)

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

		server := NewPulseServer(nil, jetKeeper)

		err := server.Export(&GetPulses{PulseNumber: 0, Count: 1}, &stream)
		require.NoError(t, err)

		require.Equal(t, 1, len(pulses))
		require.Equal(t, pulse.MinTimePulse, int(pulses[0]))
	})
}
