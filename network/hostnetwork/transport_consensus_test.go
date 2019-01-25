/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package hostnetwork

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/network/utils"
	"github.com/stretchr/testify/require"
)

func createTwoConsensusNetworks(id1, id2 core.ShortNodeID) (t1, t2 network.ConsensusNetwork, err error) {
	m := newMockResolver()

	cn1, err := NewConsensusNetwork("127.0.0.1:0", ID1+DOMAIN, id1, m)
	if err != nil {
		return nil, nil, err
	}
	cn2, err := NewConsensusNetwork("127.0.0.1:0", ID2+DOMAIN, id2, m)
	if err != nil {
		return nil, nil, err
	}

	ref1, err := core.NewRefFromBase58(ID2 + DOMAIN)
	if err != nil {
		return nil, nil, err
	}
	routing1, err := host.NewHostNS(cn1.PublicAddress(), *ref1, id1)
	if err != nil {
		return nil, nil, err
	}
	ref2, err := core.NewRefFromBase58(ID2 + DOMAIN)
	if err != nil {
		return nil, nil, err
	}
	routing2, err := host.NewHostNS(cn2.PublicAddress(), *ref2, id2)
	if err != nil {
		return nil, nil, err
	}
	m.addMappingHost(routing1)
	m.addMappingHost(routing2)

	return cn1, cn2, nil
}

func TestTransportConsensus_SendRequest(t *testing.T) {
	cn1, cn2, err := createTwoConsensusNetworks(0, 1)
	require.NoError(t, err)
	ctx := context.Background()
	ctx2 := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(r network.Request) {
		log.Info("handler triggered")
		wg.Done()
	}
	cn2.RegisterRequestHandler(types.Phase1, handler)

	cn2.Start(ctx)
	cn1.Start(ctx2)
	defer func() {
		cn1.Stop()
		cn2.Stop()
	}()

	packet := packets.NewPhase1Packet()
	request := cn1.NewRequestBuilder().Type(types.Phase1).Data(packet).Build()
	err = cn1.SendRequest(request, cn2.GetNodeID())
	require.NoError(t, err)
	success := utils.WaitTimeout(&wg, time.Second)
	require.True(t, success)
}

func TestTransportConsensus_RegisterPacketHandler(t *testing.T) {
	m := newMockResolver()

	cn, err := NewConsensusNetwork("127.0.0.1:0", ID1+DOMAIN, 0, m)
	require.NoError(t, err)
	defer cn.Stop()
	handler := func(request network.Request) {
		// do nothing
	}
	f := func() {
		cn.RegisterRequestHandler(types.Phase1, handler)
	}
	require.NotPanics(t, f, "first request handler register should not panic")
	require.Panics(t, f, "second request handler register should panic because it is already registered")
}
