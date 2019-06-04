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

package servicenetwork

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils"
	networkUtils "github.com/insolar/insolar/testutils/network"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type PublisherMock struct {
	msg *message.Message
}

func (p *PublisherMock) Publish(topic string, messages ...*message.Message) error {
	return nil
}

func (p *PublisherMock) Close() error {
	return nil
}

func checkRepliedMsg(t *testing.T, expectedMsg *message.Message, errText string) {
	require.NotNil(t, expectedMsg)
	replyPayload, err := payload.Unmarshal(expectedMsg.Payload)
	require.NoError(t, err)
	require.Contains(t, replyPayload.(*payload.Error).Text, errText)
}

func prepareNetwork(t *testing.T, cfg configuration.Configuration) *ServiceNetwork {
	serviceNetwork, err := NewServiceNetwork(cfg, &component.Manager{}, false)
	require.NoError(t, err)

	nodeKeeper := networkUtils.NewNodeKeeperMock(t)
	nodeMock := networkUtils.NewNetworkNodeMock(t)
	nodeMock.IDMock.Return(testutils.RandomRef())
	nodeKeeper.GetOriginMock.Return(nodeMock)
	serviceNetwork.NodeKeeper = nodeKeeper

	return serviceNetwork
}

func TestSendMessageHandler_ReceiverNotSet(t *testing.T) {
	cfg := configuration.NewConfiguration()
	var expectedMsg *message.Message

	serviceNetwork := prepareNetwork(t, cfg)
	sender := bus.NewSenderMock(t)
	serviceNetwork.Sender = sender
	sender.ReplyFunc = func(p context.Context, p1 *message.Message, p2 *message.Message) {
		expectedMsg = p2
	}

	p := []byte{1, 2, 3, 4, 5}
	inMsg := message.NewMessage(watermill.NewUUID(), p)

	outMsgs, err := serviceNetwork.SendMessageHandler(inMsg)
	require.NoError(t, err)
	checkRepliedMsg(t, expectedMsg, "failed to send message: Receiver in msg.Metadata not set")
	require.Nil(t, outMsgs)
}

func TestSendMessageHandler_IncorrectReceiver(t *testing.T) {
	cfg := configuration.NewConfiguration()
	cfg.Service.Skip = 5

	var expectedMsg *message.Message

	serviceNetwork := prepareNetwork(t, cfg)
	sender := bus.NewSenderMock(t)
	serviceNetwork.Sender = sender
	sender.ReplyFunc = func(p context.Context, p1 *message.Message, p2 *message.Message) {
		expectedMsg = p2
	}

	p := []byte{1, 2, 3, 4, 5}
	inMsg := message.NewMessage(watermill.NewUUID(), p)
	inMsg.Metadata.Set(bus.MetaReceiver, "someBadValue")

	outMsgs, err := serviceNetwork.SendMessageHandler(inMsg)
	require.NoError(t, err)
	checkRepliedMsg(t, expectedMsg, "incorrect Receiver in msg.Metadata")
	require.Nil(t, outMsgs)
}

func TestSendMessageHandler_SameNode(t *testing.T) {
	cfg := configuration.NewConfiguration()
	cfg.Service.Skip = 5
	serviceNetwork, err := NewServiceNetwork(cfg, &component.Manager{}, false)
	nodeRef := testutils.RandomRef()
	nodeN := networkUtils.NewNodeKeeperMock(t)
	nodeN.GetOriginFunc = func() (r insolar.NetworkNode) {
		n := networkUtils.NewNetworkNodeMock(t)
		n.IDFunc = func() (r insolar.Reference) {
			return nodeRef
		}
		return n
	}
	pubMock := &PublisherMock{}
	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)
	serviceNetwork.PulseAccessor = pulseMock
	serviceNetwork.NodeKeeper = nodeN
	serviceNetwork.Pub = pubMock

	payload := []byte{1, 2, 3, 4, 5}
	inMsg := message.NewMessage(watermill.NewUUID(), payload)
	inMsg.Metadata.Set(bus.MetaReceiver, nodeRef.String())

	outMsgs, err := serviceNetwork.SendMessageHandler(inMsg)
	require.NoError(t, err)
	require.Nil(t, outMsgs)
}

func TestSendMessageHandler_SendError(t *testing.T) {
	cfg := configuration.NewConfiguration()
	cfg.Service.Skip = 5
	pubMock := &PublisherMock{}
	serviceNetwork, err := NewServiceNetwork(cfg, &component.Manager{}, false)
	serviceNetwork.Pub = pubMock
	nodeN := networkUtils.NewNodeKeeperMock(t)
	nodeN.GetOriginFunc = func() (r insolar.NetworkNode) {
		n := networkUtils.NewNetworkNodeMock(t)
		n.IDFunc = func() (r insolar.Reference) {
			return testutils.RandomRef()
		}
		return n
	}
	controller := networkUtils.NewControllerMock(t)
	controller.SendBytesFunc = func(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error) {
		return nil, errors.New("test error")
	}
	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)
	serviceNetwork.PulseAccessor = pulseMock
	serviceNetwork.Controller = controller
	serviceNetwork.NodeKeeper = nodeN
	var expectedMsg *message.Message
	sender := bus.NewSenderMock(t)
	serviceNetwork.Sender = sender
	sender.ReplyFunc = func(p context.Context, p1 *message.Message, p2 *message.Message) {
		expectedMsg = p2
	}

	p := []byte{1, 2, 3, 4, 5}
	inMsg := message.NewMessage(watermill.NewUUID(), p)
	inMsg.Metadata.Set(bus.MetaReceiver, testutils.RandomRef().String())

	outMsgs, err := serviceNetwork.SendMessageHandler(inMsg)
	require.NoError(t, err)
	checkRepliedMsg(t, expectedMsg, "error while sending watermillMsg to controller")
	require.Nil(t, outMsgs)
}

func TestSendMessageHandler_WrongReply(t *testing.T) {
	cfg := configuration.NewConfiguration()
	cfg.Service.Skip = 5
	pubMock := &PublisherMock{}
	serviceNetwork, err := NewServiceNetwork(cfg, &component.Manager{}, false)
	serviceNetwork.Pub = pubMock
	nodeN := networkUtils.NewNodeKeeperMock(t)
	nodeN.GetOriginFunc = func() (r insolar.NetworkNode) {
		n := networkUtils.NewNetworkNodeMock(t)
		n.IDFunc = func() (r insolar.Reference) {
			return testutils.RandomRef()
		}
		return n
	}
	controller := networkUtils.NewControllerMock(t)
	controller.SendBytesFunc = func(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error) {
		return nil, nil
	}
	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)
	serviceNetwork.PulseAccessor = pulseMock
	serviceNetwork.Controller = controller
	serviceNetwork.NodeKeeper = nodeN
	var expectedMsg *message.Message
	sender := bus.NewSenderMock(t)
	serviceNetwork.Sender = sender
	sender.ReplyFunc = func(p context.Context, p1 *message.Message, p2 *message.Message) {
		expectedMsg = p2
	}

	payload := []byte{1, 2, 3, 4, 5}
	inMsg := message.NewMessage(watermill.NewUUID(), payload)
	inMsg.Metadata.Set(bus.MetaReceiver, testutils.RandomRef().String())

	outMsgs, err := serviceNetwork.SendMessageHandler(inMsg)
	require.NoError(t, err)
	checkRepliedMsg(t, expectedMsg, "reply is not ack")
	require.Nil(t, outMsgs)
}

func TestSendMessageHandler(t *testing.T) {
	cfg := configuration.NewConfiguration()
	cfg.Service.Skip = 5
	serviceNetwork, err := NewServiceNetwork(cfg, &component.Manager{}, false)
	nodeN := networkUtils.NewNodeKeeperMock(t)
	nodeN.GetOriginFunc = func() (r insolar.NetworkNode) {
		n := networkUtils.NewNetworkNodeMock(t)
		n.IDFunc = func() (r insolar.Reference) {
			return testutils.RandomRef()
		}
		return n
	}
	controller := networkUtils.NewControllerMock(t)
	controller.SendBytesFunc = func(p context.Context, p1 insolar.Reference, p2 string, p3 []byte) (r []byte, r1 error) {
		return ack, nil
	}
	pulseMock := pulse.NewAccessorMock(t)
	pulseMock.LatestMock.Return(*insolar.GenesisPulse, nil)
	serviceNetwork.PulseAccessor = pulseMock
	serviceNetwork.Controller = controller
	serviceNetwork.NodeKeeper = nodeN

	payload := []byte{1, 2, 3, 4, 5}
	inMsg := message.NewMessage(watermill.NewUUID(), payload)
	inMsg.Metadata.Set(bus.MetaReceiver, testutils.RandomRef().String())

	outMsgs, err := serviceNetwork.SendMessageHandler(inMsg)
	require.NoError(t, err)
	require.Nil(t, outMsgs)
}
