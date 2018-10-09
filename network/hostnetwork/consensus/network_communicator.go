/*
 *
 *  *    Copyright 2018 Insolar
 *  *
 *  *    Licensed under the Apache License, Version 2.0 (the "License");
 *  *    you may not use this file except in compliance with the License.
 *  *    You may obtain a copy of the License at
 *  *
 *  *        http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  *    Unless required by applicable law or agreed to in writing, software
 *  *    distributed under the License is distributed on an "AS IS" BASIS,
 *  *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  *    See the License for the specific language governing permissions and
 *  *    limitations under the License.
 *
 */

package consensus

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/nodekeeper"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/pkg/errors"
)

type communicatorReceiver struct {
	keeper nodekeeper.NodeKeeper
}

type communicatorSender struct {
	handler hosthandler.HostHandler
	network *nodenetwork.NodeNetwork
}

func (c *communicatorReceiver) ExchangeData(number core.PulseNumber, ctx context.Context,
	p consensus.Participant, data []*core.ActiveNode) ([]*core.ActiveNode, error) {
	return nil, errors.New("not implemented")
}

func (c *communicatorReceiver) ExchangeHash(number core.PulseNumber, ctx context.Context,
	p consensus.Participant, data []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (c *communicatorSender) ExchangeData(number core.PulseNumber, ctx context.Context,
	p consensus.Participant, data []*core.ActiveNode) ([]*core.ActiveNode, error) {

	sender, receiver, err := c.getSenderAndReceiver(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeData: error sending data to remote party")
	}
	request := packet.NewBuilder().Type(packet.TypeExchangeUnsyncLists).Sender(sender).Receiver(receiver).
		Request(&packet.RequestExchangeUnsyncLists{Pulse: number, UnsyncList: data}).Build()
	f, err := c.handler.SendRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeData: error sending data to remote party")
	}
	response, err := f.GetResult(c.handler.GetPacketTimeout())
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeData: error sending data to remote party")
	}
	responseData := response.Data.(*packet.ResponseExchangeUnsyncLists)
	if responseData.Error != "" {
		return nil, errors.New("ExchangeData: got error from remote party: " + responseData.Error)
	}
	return responseData.UnsyncList, nil
}

func (c *communicatorSender) ExchangeHash(number core.PulseNumber, ctx context.Context,
	p consensus.Participant, data []byte) ([]byte, error) {

	sender, receiver, err := c.getSenderAndReceiver(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeHash: error sending data to remote party")
	}
	request := packet.NewBuilder().Type(packet.TypeExchangeUnsyncHash).Sender(sender).Receiver(receiver).
		Request(&packet.RequestExchangeUnsyncHash{Pulse: number, UnsyncHash: data}).Build()
	f, err := c.handler.SendRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeHash: error sending data to remote party")
	}
	response, err := f.GetResult(c.handler.GetPacketTimeout())
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeHash: error sending data to remote party")
	}
	responseData := response.Data.(*packet.ResponseExchangeUnsyncHash)
	if responseData.Error != "" {
		return nil, errors.New("ExchangeHash: got error from remote party: " + responseData.Error)
	}
	return responseData.UnsyncHash, nil
}

func (c *communicatorSender) getSenderAndReceiver(ctx context.Context, p consensus.Participant) (*host.Host, *host.Host, error) {
	ht := c.handler.HtFromCtx(ctx)
	sender := ht.Origin
	receiverID := c.network.ResolveHostID(p.GetActiveNode().NodeID)
	receiver, exists, err := c.handler.FindHost(ctx, receiverID)
	if err != nil || !exists {
		return nil, nil, errors.Wrap(err, "Error resolving receiver HostID -> Address")
	}
	return sender, receiver, nil
}
