// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package gateway

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/pulse"
	mock "github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
)

func TestWaitMinroles_MinrolesNotHappenedInETA(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	nodeKeeper := mock.NewNodeKeeperMock(mc)
	nodeKeeper.GetAccessorMock.Set(func(p1 insolar.PulseNumber) (a1 network.Accessor) {
		accessor := mock.NewAccessorMock(mc)
		accessor.GetWorkingNodesMock.Set(func() (na1 []insolar.NetworkNode) {
			return []insolar.NetworkNode{}
		})
		return accessor
	})

	cert := &certificate.Certificate{}
	cert.MinRoles.HeavyMaterial = 1
	b := createBase(mc)
	b.CertificateManager = certificate.NewCertificateManager(cert)
	b.NodeKeeper = nodeKeeper
	waitMinRoles := newWaitMinRoles(b)

	gatewayer := mock.NewGatewayerMock(mc)
	gatewayer.GatewayMock.Set(func() network.Gateway {
		return waitMinRoles
	})

	assert.Equal(t, insolar.WaitMinRoles, waitMinRoles.GetState())
	waitMinRoles.Gatewayer = gatewayer
	waitMinRoles.bootstrapETA = time.Millisecond
	waitMinRoles.bootstrapTimer = time.NewTimer(waitMinRoles.bootstrapETA)

	waitMinRoles.Run(context.Background(), *insolar.EphemeralPulse)
}

func TestWaitMinroles_MinrolesHappenedInETA(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	gatewayer := mock.NewGatewayerMock(mc)
	gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		assert.Equal(t, insolar.WaitPulsar, state)
	})

	ref := gen.Reference()
	nodeKeeper := mock.NewNodeKeeperMock(mc)

	accessor1 := mock.NewAccessorMock(mc)
	accessor1.GetWorkingNodesMock.Set(func() (na1 []insolar.NetworkNode) {
		return []insolar.NetworkNode{}
	})
	accessor2 := mock.NewAccessorMock(mc)
	accessor2.GetWorkingNodesMock.Set(func() (na1 []insolar.NetworkNode) {
		n := node.NewNode(ref, insolar.StaticRoleLightMaterial, nil, "127.0.0.1:123", "")
		return []insolar.NetworkNode{n}
	})
	nodeKeeper.GetAccessorMock.Set(func(p insolar.PulseNumber) (a1 network.Accessor) {
		if p == pulse.MinTimePulse {
			return accessor1
		}
		return accessor2
	})

	discoveryNode := certificate.BootstrapNode{NodeRef: ref.String()}
	cert := &certificate.Certificate{MajorityRule: 1, BootstrapNodes: []certificate.BootstrapNode{discoveryNode}}
	cert.MinRoles.LightMaterial = 1
	pulseAccessor := mock.NewPulseAccessorMock(mc)
	pulseAccessor.GetPulseMock.Set(func(ctx context.Context, p1 insolar.PulseNumber) (p2 insolar.Pulse, err error) {
		p := *insolar.GenesisPulse
		p.PulseNumber += 10
		return p, nil
	})
	waitMinRoles := newWaitMinRoles(&Base{
		CertificateManager: certificate.NewCertificateManager(cert),
		NodeKeeper:         nodeKeeper,
		PulseAccessor:      pulseAccessor,
	})
	waitMinRoles.Gatewayer = gatewayer
	waitMinRoles.bootstrapETA = time.Second * 2
	waitMinRoles.bootstrapTimer = time.NewTimer(waitMinRoles.bootstrapETA)

	go waitMinRoles.Run(context.Background(), *insolar.EphemeralPulse)
	time.Sleep(100 * time.Millisecond)

	waitMinRoles.OnConsensusFinished(context.Background(), network.Report{PulseNumber: pulse.MinTimePulse + 10})
}
