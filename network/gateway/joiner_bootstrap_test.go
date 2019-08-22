//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package gateway

import (
	"context"
	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/gateway/bootstrap"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	mock "github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
