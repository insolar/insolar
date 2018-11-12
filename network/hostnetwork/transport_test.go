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
	"encoding/gob"
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	InvalidPacket types.PacketType = 1024

	ID1 string = "123"
	ID2 string = "234"
	ID3 string = "345"
)

type MockResolver struct {
	mapping map[core.RecordRef]*host.Host
}

func (m *MockResolver) Resolve(nodeID core.RecordRef) (*host.Host, error) {
	result, exist := m.mapping[nodeID]
	if !exist {
		return nil, errors.New("failed to resolve")
	}
	return result, nil
}

func (m *MockResolver) ResolveS(core.ShortNodeID) (*host.Host, error) {
	return nil, errors.New("not needed")
}
func (m *MockResolver) Start(components core.Components)  {}
func (m *MockResolver) AddToKnownHosts(h *host.Host)      {}
func (m *MockResolver) Rebalance(network.PartitionPolicy) {}
func (m *MockResolver) GetLocalNodes() []core.RecordRef   { return nil }
func (m *MockResolver) GetRandomNodes(int) []host.Host    { return nil }

func (m *MockResolver) addMapping(key, value string) error {
	k := core.NewRefFromBase58(key)
	h, err := host.NewHostN(value, k)
	if err != nil {
		return err
	}
	m.mapping[k] = h
	return nil
}

func mockConfiguration(nodeID string, address string) configuration.Configuration {
	result := configuration.Configuration{}
	result.Host.Transport = configuration.Transport{Protocol: "UTP", Address: address, BehindNAT: false}
	result.Node.Node = &configuration.Node{nodeID}
	return result
}

func TestNewInternalTransport(t *testing.T) {
	// broken address
	_, err := NewInternalTransport(mockConfiguration(ID1, "abirvalg"))
	assert.Error(t, err)
	address := "127.0.0.1:0"
	tp, err := NewInternalTransport(mockConfiguration(ID1, address))
	assert.NoError(t, err)
	defer tp.Stop()
	// assert that new address with correct port has been assigned
	assert.NotEqual(t, address, tp.PublicAddress())
	assert.Equal(t, core.NewRefFromBase58(ID1), tp.GetNodeID())
}

func TestNewInternalTransport2(t *testing.T) {
	tp, err := NewInternalTransport(mockConfiguration(ID1, "127.0.0.1:0"))
	assert.NoError(t, err)
	go tp.Start()
	// no assertion, check that Stop does not block
	defer func(t *testing.T) {
		tp.Stop()
		assert.True(t, true)
	}(t)
}

func createTwoHostNetworks(id1, id2 string) (t1, t2 network.HostNetwork, err error) {
	m := MockResolver{
		mapping: make(map[core.RecordRef]*host.Host),
	}

	i1, err := NewInternalTransport(mockConfiguration(ID1, "127.0.0.1:0"))
	if err != nil {
		return nil, nil, err
	}
	tr1 := NewHostTransport(i1, &m)
	i2, err := NewInternalTransport(mockConfiguration(ID2, "127.0.0.1:0"))
	if err != nil {
		return nil, nil, err
	}
	tr2 := NewHostTransport(i2, &m)

	err = m.addMapping(id1, tr1.PublicAddress())
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to add mapping %s -> %s", id1, tr1.PublicAddress())
	}
	err = m.addMapping(id2, tr2.PublicAddress())
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to add mapping %s -> %s", id2, tr2.PublicAddress())
	}

	return tr1, tr2, nil
}

func TestNewInternalTransport3(t *testing.T) {
	_, err := NewInternalTransport(mockConfiguration("", "127.0.0.1:0"))
	assert.Error(t, err)
}

func TestNewHostTransport(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1, ID2)
	assert.Equal(t, core.NewRefFromBase58(ID1), t1.GetNodeID())
	assert.Equal(t, core.NewRefFromBase58(ID2), t2.GetNodeID())
	assert.NoError(t, err)

	count := 10
	wg := sync.WaitGroup{}
	wg.Add(count)

	handler := func(request network.Request) (network.Response, error) {
		log.Info("handler triggered")
		wg.Done()
		return t2.BuildResponse(request, nil), nil
	}
	t2.RegisterRequestHandler(types.Ping, handler)

	t2.Start()
	t1.Start()

	defer func() {
		t1.Stop()
		t2.Stop()
	}()

	for i := 0; i < count; i++ {
		request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
		_, err := t1.SendRequest(request, core.NewRefFromBase58(ID2))
		assert.NoError(t, err)
	}
	wg.Wait()
}

func TestHostTransport_SendRequestPacket(t *testing.T) {
	m := MockResolver{
		mapping: make(map[core.RecordRef]*host.Host),
	}

	i1, err := NewInternalTransport(mockConfiguration(ID1, "127.0.0.1:0"))
	assert.NoError(t, err)
	t1 := NewHostTransport(i1, &m)
	t1.Start()
	defer t1.Stop()

	unknownID := core.NewRefFromBase58("unknown")

	// should return error because cannot resolve NodeID -> Address
	request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
	_, err = t1.SendRequest(request, unknownID)
	assert.Error(t, err)

	err = m.addMapping(ID2, "abirvalg")
	assert.Error(t, err)
	err = m.addMapping(ID3, "127.0.0.1:9090")
	assert.NoError(t, err)

	// should return error because resolved address is invalid
	_, err = t1.SendRequest(request, core.NewRefFromBase58(ID2))
	assert.Error(t, err)

	// request = t1.NewRequestBuilder().Type(InvalidPacket).Data(nil).Build()
	// should return error because packet type is invalid
	// _, err = t1.SendRequest(request, core.NewRefFromBase58(ID3))
	// assert.Error(t, err)
}

func TestHostTransport_SendRequestPacket2(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1, ID2)
	assert.NoError(t, err)

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(r network.Request) (network.Response, error) {
		log.Info("handler triggered")
		assert.Equal(t, core.NewRefFromBase58(ID1), r.GetSender())
		assert.Equal(t, t1.PublicAddress(), r.GetSenderHost().Address.String())
		wg.Done()
		return t2.BuildResponse(r, nil), nil
	}

	t2.RegisterRequestHandler(types.Ping, handler)

	t2.Start()
	t1.Start()
	defer func() {
		t1.Stop()
		t2.Stop()
	}()

	request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
	assert.Equal(t, core.NewRefFromBase58(ID1), request.GetSender())
	assert.Equal(t, t1.PublicAddress(), request.GetSenderHost().Address.String())

	_, err = t1.SendRequest(request, core.NewRefFromBase58(ID2))
	assert.NoError(t, err)
	wg.Wait()
	assert.True(t, true)
}

func TestHostTransport_SendRequestPacket3(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1, ID2)
	assert.NoError(t, err)

	type Data struct {
		Number int
	}
	gob.Register(&Data{})

	handler := func(r network.Request) (network.Response, error) {
		log.Info("handler triggered")
		d := r.GetData().(*Data)
		return t2.BuildResponse(r, &Data{Number: d.Number + 1}), nil
	}
	t2.RegisterRequestHandler(types.Ping, handler)

	t2.Start()
	t1.Start()
	defer func() {
		t1.Stop()
		t2.Stop()
	}()

	magicNumber := 42
	request := t1.NewRequestBuilder().Type(types.Ping).Data(&Data{Number: magicNumber}).Build()
	f, err := t1.SendRequest(request, core.NewRefFromBase58(ID2))
	assert.NoError(t, err)
	assert.Equal(t, f.GetRequest().GetSender(), request.GetSender())

	r, err := f.GetResponse(time.Second)
	assert.NoError(t, err)

	d := r.GetData().(*Data)
	assert.Equal(t, magicNumber+1, d.Number)

	magicNumber = 666
	request = t1.NewRequestBuilder().Type(types.Ping).Data(&Data{Number: magicNumber}).Build()
	f, err = t1.SendRequest(request, core.NewRefFromBase58(ID2))
	assert.NoError(t, err)

	r = <-f.Response()
	d = r.GetData().(*Data)
	assert.Equal(t, magicNumber+1, d.Number)
}

func TestHostTransport_SendRequestPacket_errors(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1, ID2)
	assert.NoError(t, err)

	handler := func(r network.Request) (network.Response, error) {
		log.Info("handler triggered")
		time.Sleep(time.Second)
		return t2.BuildResponse(r, nil), nil
	}
	t2.RegisterRequestHandler(types.Ping, handler)

	t2.Start()
	defer t2.Stop()
	t1.Start()

	request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
	f, err := t1.SendRequest(request, core.NewRefFromBase58(ID2))
	assert.NoError(t, err)

	_, err = f.GetResponse(time.Millisecond)
	assert.Error(t, err)

	f, err = t1.SendRequest(request, core.NewRefFromBase58(ID2))
	assert.NoError(t, err)
	t1.Stop()

	_, err = f.GetResponse(time.Second)
	assert.Error(t, err)
}
