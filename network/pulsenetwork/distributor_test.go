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

package pulsenetwork

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/transport"
	mock "github.com/insolar/insolar/testutils/network"
)

const (
	PULSENUMBER = 155
	ID1         = "14K2V1kpVycZ6qSFsNdz2FtpNxnJs17eBNzf9rdCMcKoe"
	DOMAIN      = ".14F7BsTMVPKFshM1MwLf6y23cid6fL3xMpazVoF9krzUw"
)

func createHostNetwork(t *testing.T) (network.HostNetwork, error) {
	m := mock.NewRoutingTableMock(t)

	cm1 := component.NewManager(nil)
	f1 := transport.NewFactory(configuration.NewHostNetwork().Transport)
	n1, err := hostnetwork.NewHostNetwork(ID1 + DOMAIN)
	if err != nil {
		return nil, err
	}
	cm1.Inject(f1, n1, m)

	ctx := context.Background()

	err = n1.Start(ctx)
	if err != nil {
		return nil, err
	}

	return n1, nil
}

func TestDistributor_Distribute(t *testing.T) {
	n1, err := createHostNetwork(t)
	require.NoError(t, err)
	ctx := context.Background()

	handler := func(ctx context.Context, r network.ReceivedPacket) (network.Packet, error) {
		log.Info("handle Pulse")
		pulse := r.GetRequest().GetPulse()
		assert.EqualValues(t, PULSENUMBER, pulse.Pulse.PulseNumber)
		return n1.BuildResponse(ctx, r, &packet.BasicResponse{Success: true}), nil
	}
	n1.RegisterRequestHandler(types.Pulse, handler)

	err = n1.Start(ctx)
	require.NoError(t, err)
	defer func() {
		err = n1.Stop(ctx)
		require.NoError(t, err)
	}()

	pulsarCfg := configuration.NewPulsar()
	pulsarCfg.DistributionTransport.Address = "127.0.0.1:0"
	pulsarCfg.PulseDistributor.BootstrapHosts = []string{n1.PublicAddress()}

	d, err := NewDistributor(pulsarCfg.PulseDistributor)
	require.NoError(t, err)
	assert.NotNil(t, d)

	cm := component.NewManager(nil)
	cm.Inject(d, transport.NewFactory(pulsarCfg.DistributionTransport))
	err = cm.Init(ctx)
	require.NoError(t, err)
	err = cm.Start(ctx)
	require.NoError(t, err)
	defer func() {
		err = cm.Stop(ctx)
		require.NoError(t, err)
	}()

	d.Distribute(ctx, insolar.Pulse{PulseNumber: PULSENUMBER})
	d.Distribute(ctx, insolar.Pulse{PulseNumber: PULSENUMBER})
	d.Distribute(ctx, insolar.Pulse{PulseNumber: PULSENUMBER})
}
