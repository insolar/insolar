// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gateway

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/gateway/bootstrap"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	mock "github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
)

type fixture struct {
	mc              *minimock.Controller
	joinerBootstrap *JoinerBootstrap
	gatewayer       *mock.GatewayerMock
	requester       *bootstrap.RequesterMock
}

func createFixture(t *testing.T) fixture {
	mc := minimock.NewController(t)
	cert := &certificate.Certificate{}
	gatewayer := mock.NewGatewayerMock(mc)
	requester := bootstrap.NewRequesterMock(mc)

	joinerBootstrap := newJoinerBootstrap(&Base{
		CertificateManager: certificate.NewCertificateManager(cert),
		BootstrapRequester: requester,
		Gatewayer:          gatewayer,
		originCandidate:    &adapters.Candidate{},
	})

	return fixture{
		mc:              mc,
		joinerBootstrap: joinerBootstrap,
		gatewayer:       gatewayer,
		requester:       requester,
	}
}

func TestJoinerBootstrap_Run_AuthorizeRequestFailed(t *testing.T) {
	f := createFixture(t)
	defer f.mc.Finish()
	defer f.mc.Wait(time.Minute)

	f.gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		assert.Equal(t, insolar.NoNetworkState, state)
	})

	f.requester.AuthorizeMock.Set(func(ctx context.Context, c2 insolar.Certificate) (pp1 *packet.Permit, err error) {
		return nil, insolar.ErrUnknown
	})

	assert.Equal(t, insolar.JoinerBootstrap, f.joinerBootstrap.GetState())
	f.joinerBootstrap.Run(context.Background(), *insolar.EphemeralPulse)
}

func TestJoinerBootstrap_Run_BootstrapRequestFailed(t *testing.T) {
	f := createFixture(t)
	defer f.mc.Finish()
	defer f.mc.Wait(time.Minute)

	f.gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		assert.Equal(t, insolar.NoNetworkState, state)
	})

	f.requester.AuthorizeMock.Set(func(ctx context.Context, c2 insolar.Certificate) (pp1 *packet.Permit, err error) {
		return &packet.Permit{}, nil
	})

	f.requester.BootstrapMock.Set(func(ctx context.Context, pp1 *packet.Permit, c2 adapters.Candidate, pp2 *insolar.Pulse) (bp1 *packet.BootstrapResponse, err error) {
		return nil, insolar.ErrUnknown
	})

	f.joinerBootstrap.Run(context.Background(), *insolar.EphemeralPulse)
}

func TestJoinerBootstrap_Run_BootstrapSucceeded(t *testing.T) {
	f := createFixture(t)
	defer f.mc.Finish()
	defer f.mc.Wait(time.Minute)

	f.gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		assert.Equal(t, insolar.PulseNumber(123), pulse.PulseNumber)
		assert.Equal(t, insolar.WaitConsensus, state)
	})

	f.requester.AuthorizeMock.Set(func(ctx context.Context, c2 insolar.Certificate) (pp1 *packet.Permit, err error) {
		return &packet.Permit{}, nil
	})

	f.requester.BootstrapMock.Set(func(ctx context.Context, pp1 *packet.Permit, c2 adapters.Candidate, pp2 *insolar.Pulse) (bp1 *packet.BootstrapResponse, err error) {
		p := pulse.PulseProto{PulseNumber: 123}
		return &packet.BootstrapResponse{
			ETASeconds: 90,
			Pulse:      p,
		}, nil
	})

	f.joinerBootstrap.Run(context.Background(), *insolar.EphemeralPulse)

	assert.Equal(t, true, f.joinerBootstrap.bootstrapTimer.Stop())
	assert.Equal(t, time.Duration(0), f.joinerBootstrap.backoff)
	assert.Equal(t, time.Duration(time.Second*90), f.joinerBootstrap.bootstrapETA)
}
