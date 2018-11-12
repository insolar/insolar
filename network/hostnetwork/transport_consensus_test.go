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
	"sync"
	"testing"
	"time"

	consensus "github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	Number uint32
	Header consensus.RoutingHeader
}

func (d *Data) SetPacketHeader(header *consensus.RoutingHeader) error {
	d.Header = *header
	return nil
}

func (d *Data) GetPacketHeader() (*consensus.RoutingHeader, error) {
	return &d.Header, nil
}

func createTwoConsensusNetworks(id1, id2 core.ShortNodeID) (t1, t2 network.ConsensusNetwork, err error) {
	m := newMockResolver()

	origin1, err := host.NewHostNS("127.0.0.1:0", testutils.RandomRef(), id1)
	if err != nil {
		return nil, nil, err
	}
	m.addMappingHost(origin1)
	origin2, err := host.NewHostNS("127.0.0.1:0", testutils.RandomRef(), id2)
	if err != nil {
		return nil, nil, err
	}
	m.addMappingHost(origin2)

	cn1, err := NewConsensusNetwork(origin1, m)
	if err != nil {
		return nil, nil, err
	}
	cn2, err := NewConsensusNetwork(origin2, m)
	if err != nil {
		return nil, nil, err
	}

	return cn1, cn2, nil
}

func TestTransportConsensus_SendRequest(t *testing.T) {
	t.Skip("not completed yet")
	cn1, cn2, err := createTwoConsensusNetworks(0, 1)
	assert.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(r network.Request) {
		log.Info("handler triggered")
		wg.Done()
	}
	cn2.RegisterRequestHandler(types.Ping, handler)

	cn2.Start()
	cn1.Start()
	defer func() {
		cn1.Stop()
		cn2.Stop()
	}()

	request := cn1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
	err = cn1.SendRequest(request, cn2.GetNodeID())
	assert.NoError(t, err)
	success := network.WaitTimeout(&wg, time.Second)
	assert.True(t, success)
}

func TestTransportConsensus_RegisterPacketHandler(t *testing.T) {
	m := newMockResolver()

	origin, err := host.NewHostNS("127.0.0.1:0", testutils.RandomRef(), 0)
	assert.NoError(t, err)
	cn, err := NewConsensusNetwork(origin, m)
	assert.NoError(t, err)
	defer cn.Stop()
	handler := func(request network.Request) {
		// do nothing
	}
	f := func() {
		cn.RegisterRequestHandler(types.Ping, handler)
	}
	assert.NotPanics(t, f, "first request handler register should not panic")
	assert.Panics(t, f, "second request handler register should panic because it is already registered")
}
