// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package hostnetwork

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/testutils"
)

var id1, id2, id3, idunknown string

func init() {
	id1 = gen.Reference().String()
	id2 = gen.Reference().String()
	id3 = gen.Reference().String()
	idunknown = gen.Reference().String()
}

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

func (m *MockResolver) addMapping(key, value string) error {
	k, err := insolar.NewReferenceFromString(key)
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
	defer testutils.LeakTester(t)

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
		ref, err := insolar.NewReferenceFromString(id2)
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

	n1, err := NewHostNetwork(id1)
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

	unknownID, err := insolar.NewReferenceFromString(idunknown)
	require.NoError(t, err)

	// should return error because cannot resolve NodeID -> Address
	f, err := n1.SendRequest(ctx, types.Pulse, &packet.PulseRequest{}, *unknownID)
	require.Error(t, err)
	assert.Nil(t, f)

	err = m.addMapping(id2, "abirvalg")
	require.Error(t, err)
	err = m.addMapping(id3, "127.0.0.1:7654")
	require.NoError(t, err)

	ref, err := insolar.NewReferenceFromString(id3)
	require.NoError(t, err)
	// should return error because resolved address is invalid
	f, err = n1.SendRequest(ctx, types.Pulse, &packet.PulseRequest{}, *ref)
	require.Error(t, err)
	assert.Nil(t, f)
}

func TestHostNetwork_SendRequestPacket2(t *testing.T) {
	defer testutils.LeakTester(t)
	s := newHostSuite(t)
	defer s.Stop()

	wg := sync.WaitGroup{}
	wg.Add(1)

	handler := func(ctx context.Context, r network.ReceivedPacket) (network.Packet, error) {
		inslogger.FromContext(ctx).Info("handler triggered")
		ref, err := insolar.NewReferenceFromString(id1)
		require.NoError(t, err)
		require.Equal(t, *ref, r.GetSender())
		require.Equal(t, s.n1.PublicAddress(), r.GetSenderHost().Address.String())
		wg.Done()
		return s.n2.BuildResponse(ctx, r, &packet.RPCResponse{}), nil
	}

	s.n2.RegisterRequestHandler(types.RPC, handler)

	s.Start()

	ref, err := insolar.NewReferenceFromString(id2)
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
	ref, err := insolar.NewReferenceFromString(id2)
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

	ref, err := insolar.NewReferenceFromString(id2)
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
	defer testutils.LeakTester(t)
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

	ref, err := insolar.NewReferenceFromString(id2)
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
	defer testutils.LeakTester(t)
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
		ref, err := insolar.NewReferenceFromString(id2)
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
	defer testutils.LeakTester(t)

	hn, err := NewHostNetwork(id1)
	require.NoError(t, err)

	f, err := hn.SendRequestToHost(context.Background(), types.Unknown, nil, nil)
	require.EqualError(t, err, "host network is not started")
	assert.Nil(t, f)
}
