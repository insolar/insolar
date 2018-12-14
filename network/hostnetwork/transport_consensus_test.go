/*
 *    Copyright 2018 Insolar
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
