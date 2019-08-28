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

package hostnetwork

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/insolar/insolar/network"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/transport"
)

const (
	ID1       = "4K2V1kpVycZ6qSFsNdz2FtpNxnJs17eBNzf9rdCMcKoe"
	ID2       = "4NwnA4HWZurKyXWNowJwYmb9CwX4gBKzwQKov1ExMf8M"
	ID3       = "4Ss5JMkXAD9Z7cktFEdrqeMuT6jGMF1pVozTyPHZ6zT4"
	IDUNKNOWN = "4K3Mi2hyZ6QKgynGv33sR5n3zWmSzdo8zv5Em7X26r1w"
	DOMAIN    = ".4F7BsTMVPKFshM1MwLf6y23cid6fL3xMpazVoF9krzUw"
)

type MockResolver struct {
	mu       sync.RWMutex
	mapping  map[insolar.Reference]*host.Host
	smapping map[insolar.ShortNodeID]*host.Host
}

func (m *MockResolver) ResolveConsensus(id insolar.ShortNodeID) (*host.Host, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result, exist := m.smapping[id]
	if !exist {
		return nil, errors.New("failed to resolve")
	}
	return result, nil
}

func (m *MockResolver) ResolveConsensusRef(nodeID insolar.Reference) (*host.Host, error) {
	return m.Resolve(nodeID)
}

func (m *MockResolver) Resolve(nodeID insolar.Reference) (*host.Host, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result, exist := m.mapping[nodeID]
	if !exist {
		return nil, errors.New("failed to resolve")
	}
	return result, nil
}

func (m *MockResolver) AddToKnownHosts(h *host.Host)      {}
func (m *MockResolver) Rebalance(network.PartitionPolicy) {}

func (m *MockResolver) addMapping(key, value string) error {
	k, err := insolar.NewReferenceFromBase58(key)
	if err != nil {
		return err
	}
	h, err := host.NewHostN(value, *k)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.mapping[*k] = h
	return nil
}

func (m *MockResolver) addMappingHost(h *host.Host) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.mapping[h.NodeID] = h
	m.smapping[h.ShortID] = h
}

func newMockResolver() *MockResolver {
	return &MockResolver{
		mapping:  make(map[insolar.Reference]*host.Host),
		smapping: make(map[insolar.ShortNodeID]*host.Host),
	}
}

func TestNewHostNetwork_InvalidReference(t *testing.T) {
	n, err := NewHostNetwork("invalid reference")
	require.Error(t, err)
	require.Nil(t, n)
}

type hostSuite struct {
	t          *testing.T
	ctx1, ctx2 context.Context
	id1, id2   string
	n1, n2     network.HostNetwork
	resolver   *MockResolver
	cm1, cm2   *component.Manager
}

func newHostSuite(t *testing.T) *hostSuite {
	ctx1 := inslogger.ContextWithTrace(context.Background(), "AAA")
	ctx2 := inslogger.ContextWithTrace(context.Background(), "BBB")
	resolver := newMockResolver()
	id1 := ID1 + DOMAIN
	id2 := ID2 + DOMAIN

	cm1 := component.NewManager(nil)
	f1 := transport.NewFactory(configuration.NewHostNetwork().Transport)
	n1, err := NewHostNetwork(id1)
	require.NoError(t, err)
	cm1.Inject(f1, n1, resolver)

	cm2 := component.NewManager(nil)
	cfg2 := configuration.NewHostNetwork().Transport
	// cfg2.Address = "127.0.0.1:8087"
	f2 := transport.NewFactory(cfg2)
	n2, err := NewHostNetwork(id2)
	require.NoError(t, err)
	cm2.Inject(f2, n2, resolver)

	err = cm1.Init(ctx1)
	require.NoError(t, err)
	err = cm2.Init(ctx2)
	require.NoError(t, err)

	return &hostSuite{
		t: t, ctx1: ctx1, ctx2: ctx2, id1: id1, id2: id2, n1: n1, n2: n2, resolver: resolver, cm1: cm1, cm2: cm2,
	}
}

func (s *hostSuite) Start() {
	// start the second hostNetwork before the first because most test cases perform sending packets first -> second,
	// so the second hostNetwork should be ready to receive packets when the first starts to send
	err := s.cm1.Start(s.ctx1)
	require.NoError(s.t, err)
	err = s.cm2.Start(s.ctx2)
	require.NoError(s.t, err)

	err = s.resolver.addMapping(s.id1, s.n1.PublicAddress())
	require.NoError(s.t, err, "failed to add mapping %s -> %s: %s", s.id1, s.n1.PublicAddress(), err)
	err = s.resolver.addMapping(s.id2, s.n2.PublicAddress())
	require.NoError(s.t, err, "failed to add mapping %s -> %s: %s", s.id2, s.n2.PublicAddress(), err)
}

func (s *hostSuite) Stop() {
	// stop hostNetworks in the reverse order of their start
	err := s.cm1.Stop(s.ctx1)
	assert.NoError(s.t, err)
	err = s.cm2.Stop(s.ctx2)
	assert.NoError(s.t, err)
}

func TestNewHostNetwork(t *testing.T) {
	defer leaktest.Check(t)()

	s := newHostSuite(t)
	defer s.Stop()

	count := 10
	wg := sync.WaitGroup{}
	wg.Add(count)

	handler := func(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
		inslogger.FromContext(ctx).Info("handler triggered")
		wg.Done()
		return s.n2.BuildResponse(ctx, request, &packet.RPCResponse{}), nil
	}
	s.n2.RegisterRequestHandler(types.RPC, handler)

	s.Start()

	for i := 0; i < count; i++ {
		ref, err := insolar.NewReferenceFromBase58(ID2 + DOMAIN)
		require.NoError(t, err)
		f, err := s.n1.SendRequest(s.ctx1, types.RPC, &packet.RPCRequest{}, *ref)
		require.NoError(t, err)
		f.Cancel()
	}

	wg.Wait()
}

func TestHostNetwork_SendRequestPacket(t *testing.T) {
	m := newMockResolver()
	ctx := context.Background()

	n1, err := NewHostNetwork(ID1 + DOMAIN)
	require.NoError(t, err)

	cm := component.NewManager(nil)
	cm.Register(m, n1, transport.NewFactory(configuration.NewHostNetwork().Transport))
	cm.Inject()
	err = cm.Init(ctx)
	require.NoError(t, err)
	err = cm.Start(ctx)
	require.NoError(t, err)

	defer func() {
		err = cm.Stop(ctx)
		assert.NoError(t, err)
	}()

	unknownID, err := insolar.NewReferenceFromBase58(IDUNKNOWN + DOMAIN)
	require.NoError(t, err)

	// should return error because cannot resolve NodeID -> Address
	f, err := n1.SendRequest(ctx, types.Pulse, &packet.PulseRequest{}, *unknownID)
	require.Error(t, err)
	assert.Nil(t, f)

	err = m.addMapping(ID2+DOMAIN, "abirvalg")
	require.Error(t, err)
	err = m.addMapping(ID3+DOMAIN, "127.0.0.1:7654")
	require.NoError(t, err)

	ref, err := insolar.NewReferenceFromBase58(ID2 + DOMAIN)
	require.NoError(t, err)
	// should return error because resolved address is invalid
	f, err = n1.SendRequest(ctx, types.Pulse, &packet.PulseRequest{}, *ref)
	require.Error(t, err)
	assert.Nil(t, f)
}

func TestHostNetwork_SendRequestPacket2(t *testing.T) {
	defer leaktest.Check(t)()
	s := newHostSuite(t)
	defer s.Stop()

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(ctx context.Context, r network.ReceivedPacket) (network.Packet, error) {
		inslogger.FromContext(ctx).Info("handler triggered")
		ref, err := insolar.NewReferenceFromBase58(ID1 + DOMAIN)
		require.NoError(t, err)
		require.Equal(t, *ref, r.GetSender())
		require.Equal(t, s.n1.PublicAddress(), r.GetSenderHost().Address.String())
		wg.Done()
		return s.n2.BuildResponse(ctx, r, &packet.RPCResponse{}), nil
	}

	s.n2.RegisterRequestHandler(types.RPC, handler)

	s.Start()

	ref, err := insolar.NewReferenceFromBase58(ID2 + DOMAIN)
	require.NoError(t, err)
	f, err := s.n1.SendRequest(s.ctx1, types.RPC, &packet.RPCRequest{}, *ref)
	require.NoError(t, err)
	f.Cancel()

	wg.Wait()
}

func TestHostNetwork_SendRequestPacket3(t *testing.T) {
	s := newHostSuite(t)
	defer s.Stop()

	handler := func(ctx context.Context, r network.ReceivedPacket) (network.Packet, error) {
		inslogger.FromContext(ctx).Info("handler triggered")
		return s.n2.BuildResponse(ctx, r, &packet.BasicResponse{Error: "Error"}), nil
	}
	s.n2.RegisterRequestHandler(types.Pulse, handler)

	s.Start()

	request := &packet.PulseRequest{}
	ref, err := insolar.NewReferenceFromBase58(ID2 + DOMAIN)
	require.NoError(t, err)
	f, err := s.n1.SendRequest(s.ctx1, types.Pulse, request, *ref)
	require.NoError(t, err)

	r, err := f.WaitResponse(time.Minute)
	require.NoError(t, err)

	d := r.GetResponse().GetBasic().Error
	require.Equal(t, "Error", d)

	request = &packet.PulseRequest{}
	f, err = s.n1.SendRequest(s.ctx1, types.Pulse, request, *ref)
	require.NoError(t, err)

	r, err = f.WaitResponse(time.Second)
	assert.NoError(t, err)
	d = r.GetResponse().GetBasic().Error
	require.Equal(t, d, "Error")
}

func TestHostNetwork_SendRequestPacket_errors(t *testing.T) {
	s := newHostSuite(t)
	defer s.Stop()

	handler := func(ctx context.Context, r network.ReceivedPacket) (network.Packet, error) {
		inslogger.FromContext(ctx).Info("handler triggered")
		time.Sleep(time.Millisecond * 100)
		return s.n2.BuildResponse(ctx, r, &packet.RPCResponse{}), nil
	}
	s.n2.RegisterRequestHandler(types.RPC, handler)

	s.Start()

	ref, err := insolar.NewReferenceFromBase58(ID2 + DOMAIN)
	require.NoError(t, err)
	f, err := s.n1.SendRequest(s.ctx1, types.RPC, &packet.RPCRequest{}, *ref)
	require.NoError(t, err)

	_, err = f.WaitResponse(time.Microsecond * 10)
	require.Error(t, err)

	f, err = s.n1.SendRequest(s.ctx1, types.RPC, &packet.RPCRequest{}, *ref)
	require.NoError(t, err)

	_, err = f.WaitResponse(time.Minute)
	require.NoError(t, err)
}

func TestHostNetwork_WrongHandler(t *testing.T) {
	defer leaktest.Check(t)()
	s := newHostSuite(t)
	defer s.Stop()

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(ctx context.Context, r network.ReceivedPacket) (network.Packet, error) {
		inslogger.FromContext(ctx).Info("handler triggered")
		wg.Done()
		return s.n2.BuildResponse(ctx, r, nil), nil
	}
	s.n2.RegisterRequestHandler(types.Unknown, handler)

	s.Start()

	ref, err := insolar.NewReferenceFromBase58(ID2 + DOMAIN)
	require.NoError(t, err)
	f, err := s.n1.SendRequest(s.ctx1, types.Pulse, &packet.PulseRequest{}, *ref)
	require.NoError(t, err)
	f.Cancel()

	// should timeout because there is no handler set for Ping packet
	result := network.WaitTimeout(&wg, time.Millisecond*100)
	require.False(t, result)
	wg.Done()
}

func TestStartStopSend(t *testing.T) {
	defer leaktest.Check(t)()
	s := newHostSuite(t)
	defer s.Stop()

	wg := sync.WaitGroup{}
	wg.Add(2)

	handler := func(ctx context.Context, r network.ReceivedPacket) (network.Packet, error) {
		inslogger.FromContext(ctx).Info("handler triggered")
		wg.Done()
		return s.n2.BuildResponse(ctx, r, &packet.RPCResponse{}), nil
	}
	s.n2.RegisterRequestHandler(types.RPC, handler)

	s.Start()

	send := func() {
		ref, err := insolar.NewReferenceFromBase58(ID2 + DOMAIN)
		require.NoError(t, err)
		f, err := s.n1.SendRequest(s.ctx1, types.RPC, &packet.RPCRequest{}, *ref)
		require.NoError(t, err)
		_, err = f.WaitResponse(time.Second)
		assert.NoError(t, err)
	}

	send()

	err := s.cm1.Stop(s.ctx1)
	require.NoError(t, err)
	<-time.After(time.Millisecond * 10)

	s.ctx1 = context.Background()
	err = s.cm1.Start(s.ctx1)
	require.NoError(t, err)

	send()
	wg.Wait()
}

func TestHostNetwork_SendRequestToHost_NotStarted(t *testing.T) {
	defer leaktest.Check(t)()

	hn, err := NewHostNetwork(ID1 + DOMAIN)
	require.NoError(t, err)

	f, err := hn.SendRequestToHost(context.Background(), types.Unknown, nil, nil)
	require.EqualError(t, err, "host network is not started")
	assert.Nil(t, f)
}
