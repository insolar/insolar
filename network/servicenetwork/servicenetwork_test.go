// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package servicenetwork

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/network/controller"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/testutils"
	networkUtils "github.com/insolar/insolar/testutils/network"
)

type PublisherMock struct{}

func (p *PublisherMock) Publish(topic string, messages ...*message.Message) error {
	return nil
}

func (p *PublisherMock) Close() error {
	return nil
}

func prepareNetwork(t *testing.T, cfg configuration.GenericConfiguration) *ServiceNetwork {
	serviceNetwork, err := NewServiceNetwork(cfg.Host, component.NewManager(nil))
	require.NoError(t, err)

	nodeKeeper := networkUtils.NewNodeKeeperMock(t)
	nodeMock := networkUtils.NewNetworkNodeMock(t)
	nodeMock.IDMock.Return(gen.Reference())
	nodeKeeper.GetOriginMock.Return(nodeMock)
	serviceNetwork.NodeKeeper = nodeKeeper

	return serviceNetwork
}

func TestSendMessageHandler_ReceiverNotSet(t *testing.T) {
	cfg := configuration.NewGenericConfiguration()

	serviceNetwork := prepareNetwork(t, cfg)

	p := []byte{1, 2, 3, 4, 5}
	meta := payload.Meta{
		Payload: p,
	}
	data, err := meta.Marshal()
	require.NoError(t, err)

	inMsg := message.NewMessage(watermill.NewUUID(), data)

	err = serviceNetwork.SendMessageHandler(inMsg)
	require.NoError(t, err)
}

func TestSendMessageHandler_SameNode(t *testing.T) {
	cfg := configuration.NewGenericConfiguration()
	serviceNetwork, err := NewServiceNetwork(cfg.Host, component.NewManager(nil))
	nodeRef := gen.Reference()
	nodeN := networkUtils.NewNodeKeeperMock(t)
	nodeN.GetOriginMock.Set(func() (r insolar.NetworkNode) {
		n := networkUtils.NewNetworkNodeMock(t)
		n.IDMock.Set(func() (r insolar.Reference) {
			return nodeRef
		})
		return n
	})
	pubMock := &PublisherMock{}
	pulseMock := networkUtils.NewPulseAccessorMock(t)
	pulseMock.GetLatestPulseMock.Return(*insolar.GenesisPulse, nil)
	serviceNetwork.PulseAccessor = pulseMock
	serviceNetwork.NodeKeeper = nodeN
	serviceNetwork.Pub = pubMock

	p := []byte{1, 2, 3, 4, 5}
	meta := payload.Meta{
		Payload:  p,
		Receiver: nodeRef,
	}
	data, err := meta.Marshal()
	require.NoError(t, err)

	inMsg := message.NewMessage(watermill.NewUUID(), data)

	err = serviceNetwork.SendMessageHandler(inMsg)
	require.NoError(t, err)
}

func TestSendMessageHandler_SendError(t *testing.T) {
	cfg := configuration.NewGenericConfiguration()
	pubMock := &PublisherMock{}
	serviceNetwork, err := NewServiceNetwork(cfg.Host, component.NewManager(nil))
	serviceNetwork.Pub = pubMock
	nodeN := networkUtils.NewNodeKeeperMock(t)
	nodeN.GetOriginMock.Set(func() (r insolar.NetworkNode) {
		n := networkUtils.NewNetworkNodeMock(t)
		n.IDMock.Set(func() (r insolar.Reference) {
			return gen.Reference()
		})
		return n
	})
	rpc := controller.NewRPCControllerMock(t)
	rpc.SendBytesMock.Set(func(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error) {
		return nil, errors.New("test error")
	})
	pulseMock := networkUtils.NewPulseAccessorMock(t)
	pulseMock.GetLatestPulseMock.Return(*insolar.GenesisPulse, nil)
	serviceNetwork.PulseAccessor = pulseMock
	serviceNetwork.RPC = rpc
	serviceNetwork.NodeKeeper = nodeN

	p := []byte{1, 2, 3, 4, 5}
	meta := payload.Meta{
		Payload:  p,
		Receiver: gen.Reference(),
	}
	data, err := meta.Marshal()
	require.NoError(t, err)

	inMsg := message.NewMessage(watermill.NewUUID(), data)

	err = serviceNetwork.SendMessageHandler(inMsg)
	require.NoError(t, err)
}

func TestSendMessageHandler_WrongReply(t *testing.T) {
	cfg := configuration.NewGenericConfiguration()
	pubMock := &PublisherMock{}
	serviceNetwork, err := NewServiceNetwork(cfg.Host, component.NewManager(nil))
	serviceNetwork.Pub = pubMock
	nodeN := networkUtils.NewNodeKeeperMock(t)
	nodeN.GetOriginMock.Set(func() (r insolar.NetworkNode) {
		n := networkUtils.NewNetworkNodeMock(t)
		n.IDMock.Set(func() (r insolar.Reference) {
			return gen.Reference()
		})
		return n
	})
	rpc := controller.NewRPCControllerMock(t)
	rpc.SendBytesMock.Set(func(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error) {
		return nil, nil
	})
	pulseMock := networkUtils.NewPulseAccessorMock(t)
	pulseMock.GetLatestPulseMock.Return(*insolar.GenesisPulse, nil)
	serviceNetwork.PulseAccessor = pulseMock
	serviceNetwork.RPC = rpc
	serviceNetwork.NodeKeeper = nodeN

	p := []byte{1, 2, 3, 4, 5}
	meta := payload.Meta{
		Payload:  p,
		Receiver: gen.Reference(),
	}
	data, err := meta.Marshal()
	require.NoError(t, err)

	inMsg := message.NewMessage(watermill.NewUUID(), data)

	err = serviceNetwork.SendMessageHandler(inMsg)
	require.NoError(t, err)
}

func TestSendMessageHandler(t *testing.T) {
	cfg := configuration.NewGenericConfiguration()
	serviceNetwork, err := NewServiceNetwork(cfg.Host, component.NewManager(nil))
	nodeN := networkUtils.NewNodeKeeperMock(t)
	nodeN.GetOriginMock.Set(func() (r insolar.NetworkNode) {
		n := networkUtils.NewNetworkNodeMock(t)
		n.IDMock.Set(func() (r insolar.Reference) {
			return gen.Reference()
		})
		return n
	})
	rpc := controller.NewRPCControllerMock(t)
	rpc.SendBytesMock.Set(func(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error) {
		return ack, nil
	})
	pulseMock := networkUtils.NewPulseAccessorMock(t)
	pulseMock.GetLatestPulseMock.Return(*insolar.GenesisPulse, nil)
	serviceNetwork.PulseAccessor = pulseMock
	serviceNetwork.RPC = rpc
	serviceNetwork.NodeKeeper = nodeN

	p := []byte{1, 2, 3, 4, 5}
	meta := payload.Meta{
		Payload:  p,
		Receiver: gen.Reference(),
	}
	data, err := meta.Marshal()
	require.NoError(t, err)

	inMsg := message.NewMessage(watermill.NewUUID(), data)

	err = serviceNetwork.SendMessageHandler(inMsg)
	require.NoError(t, err)
}

type stater struct{}

func (s *stater) State() []byte {
	return []byte("123")
}

func TestServiceNetwork_StartStop(t *testing.T) {
	t.Skip("fixme")
	cm := component.NewManager(nil)
	origin := gen.Reference()
	nk := nodenetwork.NewNodeKeeper(node.NewNode(origin, insolar.StaticRoleUnknown, nil, "127.0.0.1:0", ""))
	cert := &certificate.Certificate{}
	cert.Reference = origin.String()
	certManager := certificate.NewCertificateManager(cert)
	serviceNetwork, err := NewServiceNetwork(configuration.NewGenericConfiguration().Host, cm)
	require.NoError(t, err)
	ctx := context.Background()
	defer serviceNetwork.Stop(ctx)

	cm.Inject(serviceNetwork, nk, certManager, testutils.NewCryptographyServiceMock(t), pulse.NewAccessorMock(t),
		testutils.NewTerminationHandlerMock(t), testutils.NewPulseManagerMock(t), &PublisherMock{},
		testutils.NewContractRequesterMock(t), bus.NewSenderMock(t), &stater{},
		testutils.NewPlatformCryptographyScheme(), testutils.NewKeyProcessorMock(t))
	err = serviceNetwork.Init(ctx)
	require.NoError(t, err)
	err = serviceNetwork.Start(ctx)
	require.NoError(t, err)
}

type publisherMock struct {
	Error error
}

func (pm *publisherMock) Publish(topic string, messages ...*message.Message) error { return pm.Error }
func (pm *publisherMock) Close() error                                             { return nil }

func TestServiceNetwork_processIncoming(t *testing.T) {
	serviceNetwork, err := NewServiceNetwork(configuration.NewGenericConfiguration().Host, component.NewManager(nil))
	require.NoError(t, err)
	pub := &publisherMock{}
	serviceNetwork.Pub = pub
	ctx := context.Background()
	_, err = serviceNetwork.processIncoming(ctx, []byte("ololo"))
	assert.Error(t, err)
	msg := message.NewMessage("1", nil)
	data, err := serializeMessage(msg)
	require.NoError(t, err)
	_, err = serviceNetwork.processIncoming(ctx, data)
	assert.NoError(t, err)
	pub.Error = errors.New("Failed to publish message")
	_, err = serviceNetwork.processIncoming(ctx, data)
	assert.Error(t, err)
}
