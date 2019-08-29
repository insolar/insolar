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
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/pulse"
	mock "github.com/insolar/insolar/testutils/network"
)

func TestWaitMajority_MajorityNotHappenedInETA(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	gatewayer := mock.NewGatewayerMock(mc)
	gatewayer.FailStateMock.Set(func(ctx context.Context, reason string) {
		require.Equal(t, "Bootstrap timeout exceeded", reason)
	})

	nodeKeeper := mock.NewNodeKeeperMock(mc)
	nodeKeeper.GetAccessorMock.Set(func(p1 insolar.PulseNumber) (a1 network.Accessor) {
		accessor := mock.NewAccessorMock(mc)
		accessor.GetWorkingNodesMock.Set(func() (na1 []insolar.NetworkNode) {
			return []insolar.NetworkNode{}
		})
		return accessor
	})

	cert := &certificate.Certificate{MajorityRule: 4}
	waitMajority := newWaitMajority(&Base{
		CertificateManager: certificate.NewCertificateManager(cert),
		NodeKeeper:         nodeKeeper,
	})
	assert.Equal(t, insolar.WaitMajority, waitMajority.GetState())
	waitMajority.Gatewayer = gatewayer
	waitMajority.bootstrapETA = time.Millisecond
	waitMajority.bootstrapTimer = time.NewTimer(waitMajority.bootstrapETA)

	waitMajority.Run(context.Background(), *insolar.EphemeralPulse)
}

func TestWaitMajority_MajorityHappenedInETA(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	gatewayer := mock.NewGatewayerMock(mc)
	gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		assert.Equal(t, insolar.WaitMinRoles, state)
	})

	ref := gen.Reference()
	nodeKeeper := mock.NewNodeKeeperMock(mc)
	nodeKeeper.GetAccessorMock.Set(func(p insolar.PulseNumber) (a1 network.Accessor) {
		accessor := mock.NewAccessorMock(mc)
		if p == pulse.MinTimePulse {
			accessor.GetWorkingNodesMock.Set(func() (na1 []insolar.NetworkNode) {
				return []insolar.NetworkNode{}
			})
		} else {
			accessor.GetWorkingNodesMock.Set(func() (na1 []insolar.NetworkNode) {
				n := node.NewNode(ref, insolar.StaticRoleHeavyMaterial, nil, "127.0.0.1:123", "")
				return []insolar.NetworkNode{n}
			})
		}

		return accessor
	})

	discoveryNode := certificate.BootstrapNode{NodeRef: ref.String()}
	cert := &certificate.Certificate{MajorityRule: 1, BootstrapNodes: []certificate.BootstrapNode{discoveryNode}}
	pulseAccessor := mock.NewPulseAccessorMock(mc)
	pulseAccessor.GetPulseMock.Set(func(ctx context.Context, p1 insolar.PulseNumber) (p2 insolar.Pulse, err error) {
		p := *insolar.GenesisPulse
		p.PulseNumber += 10
		return p, nil
	})
	waitMajority := newWaitMajority(&Base{
		CertificateManager: certificate.NewCertificateManager(cert),
		NodeKeeper:         nodeKeeper,
		PulseAccessor:      pulseAccessor,
	})
	waitMajority.Gatewayer = gatewayer
	waitMajority.bootstrapETA = time.Second * 2
	waitMajority.bootstrapTimer = time.NewTimer(waitMajority.bootstrapETA)

	go waitMajority.Run(context.Background(), *insolar.EphemeralPulse)
	time.Sleep(100 * time.Millisecond)

	waitMajority.OnConsensusFinished(context.Background(), network.Report{PulseNumber: pulse.MinTimePulse + 10})
}
