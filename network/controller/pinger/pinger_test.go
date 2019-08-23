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

package pinger

import (
	"context"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/future"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"

	testutils "github.com/insolar/insolar/testutils/network"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPing_Errors(t *testing.T) {
	cm := component.NewManager(nil)
	f := transport.NewFactory(configuration.NewHostNetwork().Transport)
	n, err := hostnetwork.NewHostNetwork(insolar.NewEmptyReference().String())
	require.NoError(t, err)
	cm.Inject(f, n, testutils.NewRoutingTableMock(t))

	pinger := NewPinger(n)
	_, err = pinger.Ping(context.Background(), "invalid", time.Second)
	assert.Error(t, err)
	_, err = pinger.Ping(context.Background(), "127.0.0.1:0", time.Second)
	assert.Error(t, err)
}

func TestPing_HappyPath(t *testing.T) {
	ctx := context.Background()
	refs := gen.UniqueReferences(2)

	cm2 := component.NewManager(nil)
	f2 := transport.NewFactory(configuration.NewHostNetwork().Transport)
	n2, err := hostnetwork.NewHostNetwork(refs[1].String())
	require.NoError(t, err)
	defer n2.Stop(ctx)
	n2.RegisterRequestHandler(types.Ping, func(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
		return n2.BuildResponse(ctx, request, &packet.Ping{}), nil
	})
	resolver2 := testutils.NewRoutingTableMock(t)
	cm2.Inject(f2, n2, resolver2)
	err = cm2.Init(ctx)
	require.NoError(t, err)
	err = cm2.Start(ctx)
	require.NoError(t, err)

	cm := component.NewManager(nil)
	f := transport.NewFactory(configuration.NewHostNetwork().Transport)
	n, err := hostnetwork.NewHostNetwork(refs[0].String())
	defer n.Stop(ctx)
	require.NoError(t, err)
	resolver := testutils.NewRoutingTableMock(t)
	cm.Inject(f, n, resolver)
	err = cm.Init(ctx)
	require.NoError(t, err)
	err = cm.Start(ctx)
	require.NoError(t, err)

	pinger := NewPinger(n)

	_, err = pinger.Ping(ctx, n2.PublicAddress(), time.Minute)
	assert.NoError(t, err)
}

func TestPing_Timeout(t *testing.T) {
	ctx := context.Background()
	refs := gen.UniqueReferences(2)

	cm2 := component.NewManager(nil)
	f2 := transport.NewFactory(configuration.NewHostNetwork().Transport)
	n2, err := hostnetwork.NewHostNetwork(refs[1].String())
	require.NoError(t, err)
	defer n2.Stop(ctx)

	startRespondig := make(chan struct{})
	responded := make(chan struct{})
	n2.RegisterRequestHandler(types.Ping, func(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
		defer func() {
			close(responded)
		}()
		<-startRespondig
		return n2.BuildResponse(ctx, request, &packet.Ping{}), nil
	})
	resolver2 := testutils.NewRoutingTableMock(t)
	cm2.Inject(f2, n2, resolver2)
	err = cm2.Init(ctx)
	require.NoError(t, err)
	err = cm2.Start(ctx)
	require.NoError(t, err)

	cm := component.NewManager(nil)
	f := transport.NewFactory(configuration.NewHostNetwork().Transport)
	n, err := hostnetwork.NewHostNetwork(refs[0].String())
	defer n.Stop(ctx)
	require.NoError(t, err)
	resolver := testutils.NewRoutingTableMock(t)
	cm.Inject(f, n, resolver)
	err = cm.Init(ctx)
	require.NoError(t, err)
	err = cm.Start(ctx)
	require.NoError(t, err)

	pinger := NewPinger(n)

	_, err = pinger.Ping(ctx, n2.PublicAddress(), time.Nanosecond)

	require.Error(t, err)
	assert.Contains(t, err.Error(), future.ErrTimeout.Error())
	close(startRespondig)
	<-responded
}
