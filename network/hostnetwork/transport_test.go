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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	InvalidPacket types.PacketType = 1024

	ID1       = "4K2V1kpVycZ6qSFsNdz2FtpNxnJs17eBNzf9rdCMcKoe"
	ID2       = "4NwnA4HWZurKyXWNowJwYmb9CwX4gBKzwQKov1ExMf8M"
	ID3       = "4Ss5JMkXAD9Z7cktFEdrqeMuT6jGMF1pVozTyPHZ6zT4"
	IDUNKNOWN = "4K3Mi2hyZ6QKgynGv33sR5n3zWmSzdo8zv5Em7X26r1w"
	DOMAIN    = ".4F7BsTMVPKFshM1MwLf6y23cid6fL3xMpazVoF9krzUw"
)

type MockResolver struct {
	mapping  map[core.RecordRef]*host.Host
	smapping map[core.ShortNodeID]*host.Host
}

func (m *MockResolver) ResolveConsensus(id core.ShortNodeID) (*host.Host, error) {
	result, exist := m.smapping[id]
	if !exist {
		return nil, errors.New("failed to resolve")
	}
	return result, nil
}

func (m *MockResolver) ResolveConsensusRef(nodeID core.RecordRef) (*host.Host, error) {
	return m.Resolve(nodeID)
}

func (m *MockResolver) Resolve(nodeID core.RecordRef) (*host.Host, error) {
	result, exist := m.mapping[nodeID]
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
	k, err := core.NewRefFromBase58(key)
	if err != nil {
		return err
	}
	h, err := host.NewHostN(value, *k)
	if err != nil {
		return err
	}
	m.mapping[*k] = h
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
	result.Host.Transport = configuration.Transport{Protocol: "TCP", Address: address, BehindNAT: false}
	return result
}

func TestNewInternalTransport(t *testing.T) {
	// broken address
	ctx := context.Background()
	_, err := NewInternalTransport(mockConfiguration("abirvalg"), ID1+DOMAIN)
	require.Error(t, err)
	address := "127.0.0.1:0"
	tp, err := NewInternalTransport(mockConfiguration(address), ID1+DOMAIN)
	require.NoError(t, err)
	defer tp.Stop(ctx)
	// require that new address with correct port has been assigned
	require.NotEqual(t, address, tp.PublicAddress())
	ref, err := core.NewRefFromBase58(ID1 + DOMAIN)
	require.NoError(t, err)
	require.Equal(t, *ref, tp.GetNodeID())
}

func TestNewInternalTransport2(t *testing.T) {
	ctx := context.Background()
	tp, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), ID1+DOMAIN)
	require.NoError(t, err)
	go tp.Start(ctx)
	time.Sleep(time.Millisecond)
	// no assertion, check that Stop does not block
	defer func() {
		tp.Stop(ctx)
		require.True(t, true)
	}()
}

func createTwoHostNetworks(id1, id2 string) (t1, t2 network.HostNetwork, err error) {
	m := newMockResolver()

	i1, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), ID1+DOMAIN)
	if err != nil {
		return nil, nil, err
	}
	tr1 := NewHostTransport(i1, m)
	i2, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), ID2+DOMAIN)
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
	t1, t2, err := createTwoHostNetworks(ID1+DOMAIN, ID2+DOMAIN)
	ref1, err := core.NewRefFromBase58(ID1 + DOMAIN)
	require.NoError(t, err)
	require.Equal(t, *ref1, t1.GetNodeID())
	ref2, err := core.NewRefFromBase58(ID2 + DOMAIN)
	require.Equal(t, *ref2, t2.GetNodeID())
	require.NoError(t, err)

	count := 10
	wg := sync.WaitGroup{}
	wg.Add(count)

	handler := func(ctx context.Context, request network.Request) (network.Response, error) {
		log.Info("handler triggered")
		wg.Done()
		return t2.BuildResponse(ctx, request, nil), nil
	}
	t2.RegisterRequestHandler(types.Ping, handler)

	t2.Start(ctx2)
	t1.Start(ctx)

	defer func() {
		t1.Stop(ctx)
		t2.Stop(ctx2)
	}()

	for i := 0; i < count; i++ {
		request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
		ref, err := core.NewRefFromBase58(ID2 + DOMAIN)
		require.NoError(t, err)
		_, err = t1.SendRequest(ctx, request, *ref)
		require.NoError(t, err)
	}
	success := utils.WaitTimeout(&wg, time.Second)
	require.True(t, success)
}

func TestHostTransport_SendRequestPacket(t *testing.T) {
	m := newMockResolver()
	ctx := context.Background()

	i1, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), ID1+DOMAIN)
	require.NoError(t, err)
	t1 := NewHostTransport(i1, m)
	t1.Start(ctx)
	defer t1.Stop(ctx)

	unknownID, err := core.NewRefFromBase58(IDUNKNOWN + DOMAIN)
	require.NoError(t, err)

	// should return error because cannot resolve NodeID -> Address
	request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
	_, err = t1.SendRequest(ctx, request, *unknownID)
	require.Error(t, err)

	err = m.addMapping(ID2+DOMAIN, "abirvalg")
	require.Error(t, err)
	err = m.addMapping(ID3+DOMAIN, "127.0.0.1:7654")
	require.NoError(t, err)

	ref, err := core.NewRefFromBase58(ID2 + DOMAIN)
	require.NoError(t, err)
	// should return error because resolved address is invalid
	_, err = t1.SendRequest(ctx, request, *ref)
	require.Error(t, err)
}

func TestHostTransport_SendRequestPacket2(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1+DOMAIN, ID2+DOMAIN)
	require.NoError(t, err)
	ctx := context.Background()
	ctx2 := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(ctx context.Context, r network.Request) (network.Response, error) {
		log.Info("handler triggered")
		ref, err := core.NewRefFromBase58(ID1 + DOMAIN)
		require.NoError(t, err)
		require.Equal(t, *ref, r.GetSender())
		require.Equal(t, t1.PublicAddress(), r.GetSenderHost().Address.String())
		wg.Done()
		return t2.BuildResponse(ctx, r, nil), nil
	}

	t2.RegisterRequestHandler(types.Ping, handler)

	t2.Start(ctx2)
	t1.Start(ctx)
	defer func() {
		t1.Stop(ctx)
		t2.Stop(ctx2)
	}()

	request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
	ref, err := core.NewRefFromBase58(ID1 + DOMAIN)
	require.NoError(t, err)
	require.Equal(t, *ref, request.GetSender())
	require.Equal(t, t1.PublicAddress(), request.GetSenderHost().Address.String())

	ref, err = core.NewRefFromBase58(ID2 + DOMAIN)
	require.NoError(t, err)
	_, err = t1.SendRequest(ctx, request, *ref)
	require.NoError(t, err)
	success := utils.WaitTimeout(&wg, time.Second)
	require.True(t, success)
}

func TestHostTransport_SendRequestPacket3(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1+DOMAIN, ID2+DOMAIN)
	require.NoError(t, err)
	ctx := context.Background()
	ctx2 := context.Background()

	type Data struct {
		Number int
	}
	gob.Register(&Data{})

	handler := func(ctx context.Context, r network.Request) (network.Response, error) {
		log.Info("handler triggered")
		d := r.GetData().(*Data)
		return t2.BuildResponse(ctx, r, &Data{Number: d.Number + 1}), nil
	}
	t2.RegisterRequestHandler(types.Ping, handler)

	t2.Start(ctx2)
	t1.Start(ctx)
	defer func() {
		t1.Stop(ctx)
		t2.Stop(ctx2)
	}()

	magicNumber := 42
	request := t1.NewRequestBuilder().Type(types.Ping).Data(&Data{Number: magicNumber}).Build()
	ref, err := core.NewRefFromBase58(ID2 + DOMAIN)
	require.NoError(t, err)
	f, err := t1.SendRequest(ctx, request, *ref)
	require.NoError(t, err)
	require.Equal(t, f.GetRequest().GetSender(), request.GetSender())

	r, err := f.GetResponse(time.Second)
	require.NoError(t, err)

	d := r.GetData().(*Data)
	require.Equal(t, magicNumber+1, d.Number)

	magicNumber = 666
	request = t1.NewRequestBuilder().Type(types.Ping).Data(&Data{Number: magicNumber}).Build()
	f, err = t1.SendRequest(ctx, request, *ref)
	require.NoError(t, err)

	r = <-f.Response()
	d = r.GetData().(*Data)
	require.Equal(t, magicNumber+1, d.Number)
}

func TestHostTransport_SendRequestPacket_errors(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1+DOMAIN, ID2+DOMAIN)
	require.NoError(t, err)
	ctx := context.Background()
	ctx2 := context.Background()

	handler := func(ctx context.Context, r network.Request) (network.Response, error) {
		log.Info("handler triggered")
		time.Sleep(time.Second)
		return t2.BuildResponse(ctx, r, nil), nil
	}
	t2.RegisterRequestHandler(types.Ping, handler)

	t2.Start(ctx2)
	defer t2.Stop(ctx2)
	t1.Start(ctx)

	request := t1.NewRequestBuilder().Type(types.Ping).Data(nil).Build()
	ref, err := core.NewRefFromBase58(ID2 + DOMAIN)
	require.NoError(t, err)
	f, err := t1.SendRequest(ctx, request, *ref)
	require.NoError(t, err)

	_, err = f.GetResponse(time.Millisecond)
	require.Error(t, err)

	f, err = t1.SendRequest(ctx, request, *ref)
	require.NoError(t, err)
	t1.Stop(ctx)

	_, err = f.GetResponse(time.Second)
	require.Error(t, err)
}

func TestHostTransport_WrongHandler(t *testing.T) {
	t1, t2, err := createTwoHostNetworks(ID1+DOMAIN, ID2+DOMAIN)
	require.NoError(t, err)
	ctx := context.Background()
	ctx2 := context.Background()

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(ctx context.Context, r network.Request) (network.Response, error) {
		log.Info("handler triggered")
		wg.Done()
		return t2.BuildResponse(ctx, r, nil), nil
	}
	t2.RegisterRequestHandler(InvalidPacket, handler)

	t2.Start(ctx2)
	t1.Start(ctx)
	defer func() {
		t1.Stop(ctx)
		t2.Stop(ctx2)
	}()

	request := t1.NewRequestBuilder().Type(types.Ping).Build()
	ref, err := core.NewRefFromBase58(ID2 + DOMAIN)
	require.NoError(t, err)
	_, err = t1.SendRequest(ctx, request, *ref)
	require.NoError(t, err)

	// should timeout because there is no handler set for Ping packet
	result := utils.WaitTimeout(&wg, time.Millisecond*10)
	require.False(t, result)
}

func TestDoubleStart(t *testing.T) {
	ctx := context.Background()
	tp, err := NewInternalTransport(mockConfiguration("127.0.0.1:0"), ID1+DOMAIN)
	require.NoError(t, err)

	err = tp.Start(ctx)
	assert.NoError(t, err)
	err = tp.Start(ctx)
	assert.Error(t, err)

	tp.Stop(ctx)
}
