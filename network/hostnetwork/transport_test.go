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
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const (
	InvalidPacket types.PacketType = 1024

	ID1 string = "123"
	ID2 string = "234"
	ID3 string = "345"
)

type MockResolver struct {
	mapping  map[core.RecordRef]*host.Host
	smapping map[core.ShortNodeID]*host.Host
}

func (m *MockResolver) Resolve(nodeID core.RecordRef) (*host.Host, error) {
	result, exist := m.mapping[nodeID]
	if !exist {
		return nil, errors.New("failed to resolve")
	}
	return result, nil
}

func (m *MockResolver) ResolveS(id core.ShortNodeID) (*host.Host, error) {
	result, exist := m.smapping[id]
	if !exist {
		return nil, errors.New("failed to resolve")
	}
	return result, nil
}

func (m *MockResolver) Inject(nodeKeeper network.NodeKeeper) {}
func (m *MockResolver) AddToKnownHosts(h *host.Host)         {}
func (m *MockResolver) Rebalance(network.PartitionPolicy)    {}
func (m *MockResolver) GetRandomNodes(int) []host.Host       { return nil }

func (m *MockResolver) addMapping(key, value string) error {
	k := core.NewRefFromBase58(key)
	h, err := host.NewHostN(value, k)
	if err != nil {
		return err
	}
	m.mapping[k] = h
	return nil
}

func (m *MockResolver) addMappingHost(h *host.Host) {
	m.mapping[h.NodeID] = h
	m.smapping[h.ShortID] = h
}

func newMockResolver() *MockResolver {
	return &MockResolver{
		mapping:  make(map[core.RecordRef]*host.Host),
		smapping: make(map[core.ShortNodeID]*host.Host),
	}
}

func mockConfiguration(address string) configuration.Configuration {
	result := configuration.Configuration{}
	result.Host.Transport = configuration.Transport{Protocol: "UTP", Address: address, BehindNAT: false}
	return result
}

func TestNewInternalTransport(t *testing.T) {
	// broken address
	_, err := NewInternalTransport(mockConfiguration("abirvalg"), ID1)
	require.Error(t, err)
	address := "127.0.0.1:0"
	tp, err := NewInternalTransport(mockConfiguration(address), ID1)
	require.NoError(t, err)
	defer tp.Stop()
	// require that new address with correct port has been assigned
	require.NotEqual(t, address, tp.PublicAddress())
	require.Equal(t, core.NewRefFromBase58(ID1), tp.GetNodeID())
}

func TestNewInternalTransport2(t *testing.T) {
	ctx := context.Background()
	tp, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), ID1)
	require.NoError(t, err)
	go tp.Start(ctx)
	// no assertion, check that Stop does not block
	defer func(t *testing.T) {
		tp.Stop()
		require.True(t, true)
	}(t)
}

func createTwoHostNetworks(id1, id2 string) (t1, t2 network.HostNetwork, err error) {
	m := newMockResolver()

	i1, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), ID1)
	if err != nil {
		return nil, nil, err
	}
	tr1 := NewHostTransport(i1, m)
	i2, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), ID2)
	if err != nil {
		return nil, nil, err
	}
	tr2 := NewHostTransport(i2, m)

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
	_, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), "")
	require.Error(t, err)
}

func TestNewHostTransport(t *testing.T) {
	ctx := context.Background()
	ctx2 := context.Background()
	t1, t2, err := createTwoHostNetworks(ID1, ID2)
	require.Equal(t, core.NewRefFromBase58(ID1), t1.GetNodeID())
	require.Equal(t, core.NewRefFromBase58(ID2), t2.GetNodeID())
	require.NoError(t, err)

	count := 10
	wg := sync.WaitGroup{}
	wg.Add(count)

	handler := func(request network.Request) (network.Response, error) {
		log.Info("handler triggered")
		wg.Done()
		return t2.BuildResponse(request, nil), nil
	}
	t2.RegisterRequestHandler(types.Ping, handler)

	t2.Start(ctx)
	t1.Start(ctx2)

	defer func() {
		t1.Stop()
		t2.Stop()
	}()

	for i := 0; i < count; i++ {
		request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
		_, err := t1.SendRequest(request, core.NewRefFromBase58(ID2))
		require.NoError(t, err)
	}
	success := utils.WaitTimeout(&wg, time.Second)
	require.True(t, success)
}

func TestHostTransport_SendRequestPacket(t *testing.T) {
	m := newMockResolver()
	ctx := context.Background()

	i1, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), ID1)
	require.NoError(t, err)
	t1 := NewHostTransport(i1, m)
	t1.Start(ctx)
	defer t1.Stop()

	unknownID := core.NewRefFromBase58("unknown")

	// should return error because cannot resolve NodeID -> Address
	request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
	_, err = t1.SendRequest(request, unknownID)
	require.Error(t, err)

	err = m.addMapping(ID2, "abirvalg")
	require.Error(t, err)
	err = m.addMapping(ID3, "127.0.0.1:7654")
	require.NoError(t, err)

	// should return error because resolved address is invalid
	_, err = t1.SendRequest(request, core.NewRefFromBase58(ID2))
	require.Error(t, err)
}

func TestHostTransport_SendRequestPacket2(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1, ID2)
	require.NoError(t, err)
	ctx := context.Background()
	ctx2 := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(r network.Request) (network.Response, error) {
		log.Info("handler triggered")
		require.Equal(t, core.NewRefFromBase58(ID1), r.GetSender())
		require.Equal(t, t1.PublicAddress(), r.GetSenderHost().Address.String())
		wg.Done()
		return t2.BuildResponse(r, nil), nil
	}

	t2.RegisterRequestHandler(types.Ping, handler)

	t2.Start(ctx)
	t1.Start(ctx2)
	defer func() {
		t1.Stop()
		t2.Stop()
	}()

	request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
	require.Equal(t, core.NewRefFromBase58(ID1), request.GetSender())
	require.Equal(t, t1.PublicAddress(), request.GetSenderHost().Address.String())

	_, err = t1.SendRequest(request, core.NewRefFromBase58(ID2))
	require.NoError(t, err)
	success := utils.WaitTimeout(&wg, time.Second)
	require.True(t, success)
}

func TestHostTransport_SendRequestPacket3(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1, ID2)
	require.NoError(t, err)
	ctx := context.Background()
	ctx2 := context.Background()

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

	t2.Start(ctx)
	t1.Start(ctx2)
	defer func() {
		t1.Stop()
		t2.Stop()
	}()

	magicNumber := 42
	request := t1.NewRequestBuilder().Type(types.Ping).Data(&Data{Number: magicNumber}).Build()
	f, err := t1.SendRequest(request, core.NewRefFromBase58(ID2))
	require.NoError(t, err)
	require.Equal(t, f.GetRequest().GetSender(), request.GetSender())

	r, err := f.GetResponse(time.Second)
	require.NoError(t, err)

	d := r.GetData().(*Data)
	require.Equal(t, magicNumber+1, d.Number)

	magicNumber = 666
	request = t1.NewRequestBuilder().Type(types.Ping).Data(&Data{Number: magicNumber}).Build()
	f, err = t1.SendRequest(request, core.NewRefFromBase58(ID2))
	require.NoError(t, err)

	r = <-f.Response()
	d = r.GetData().(*Data)
	require.Equal(t, magicNumber+1, d.Number)
}

func TestHostTransport_SendRequestPacket_errors(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1, ID2)
	require.NoError(t, err)
	ctx := context.Background()
	ctx2 := context.Background()

	handler := func(r network.Request) (network.Response, error) {
		log.Info("handler triggered")
		time.Sleep(time.Second)
		return t2.BuildResponse(r, nil), nil
	}
	t2.RegisterRequestHandler(types.Ping, handler)

	t2.Start(ctx)
	defer t2.Stop()
	t1.Start(ctx2)

	request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
	f, err := t1.SendRequest(request, core.NewRefFromBase58(ID2))
	require.NoError(t, err)

	_, err = f.GetResponse(time.Millisecond)
	require.Error(t, err)

	f, err = t1.SendRequest(request, core.NewRefFromBase58(ID2))
	require.NoError(t, err)
	t1.Stop()

	_, err = f.GetResponse(time.Second)
	require.Error(t, err)
}

func TestHostTransport_WrongHandler(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1, ID2)
	require.NoError(t, err)
	ctx := context.Background()
	ctx2 := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(r network.Request) (network.Response, error) {
		log.Info("handler triggered")
		wg.Done()
		return t2.BuildResponse(r, nil), nil
	}
	t2.RegisterRequestHandler(InvalidPacket, handler)

	t2.Start(ctx)
	t1.Start(ctx2)
	defer func() {
		t1.Stop()
		t2.Stop()
	}()

	request := t1.NewRequestBuilder().Type(types.Ping).Build()
	_, err = t1.SendRequest(request, core.NewRefFromBase58(ID2))
	require.NoError(t, err)

	// should timeout because there is no handler set for Ping packet
	result := utils.WaitTimeout(&wg, time.Millisecond*10)
	require.False(t, result)
}

func TestDoubleStart(t *testing.T) {
	ctx := context.Background()
	tp, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), ID1)
	require.NoError(t, err)
	wg := sync.WaitGroup{}
	wg.Add(2)

	f := func(group *sync.WaitGroup, t network.InternalTransport) {
		wg.Done()
		t.Start(ctx)
	}
	go f(&wg, tp)
	go f(&wg, tp)
	wg.Wait()
	defer tp.Stop()
}

func TestHostTransport_RegisterPacketHandler(t *testing.T) {
	m := newMockResolver()

	i1, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), ID1)
	require.NoError(t, err)
	tr1 := NewHostTransport(i1, m)
	defer tr1.Stop()
	handler := func(request network.Request) (network.Response, error) {
		return tr1.BuildResponse(request, nil), nil
	}
	f := func() {
		tr1.RegisterRequestHandler(types.Ping, handler)
	}
	require.NotPanics(t, f, "first request handler register should not panic")
	require.Panics(t, f, "second request handler register should panic because it is already registered")
}
