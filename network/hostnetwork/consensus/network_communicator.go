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

package consensus

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
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
}

func (c *communicatorReceiver) ExchangeData(pulse core.PulseNumber, ctx context.Context,
	p consensus.Participant, data []*core.ActiveNode) ([]*core.ActiveNode, error) {

	currentPulse := c.keeper.GetPulse()
	if currentPulse > pulse {
		return nil, errors.Errorf("Received consensus unsync list exchange request with pulse %d but current is %d",
			pulse, currentPulse)
	}
	// TODO: block on getting unsync if currentPulse < number
	// TODO: write to communicatorSender map to decrease network requests
	// return c.keeper.GetUnsync(), nil
	return nil, errors.New("not implemented")
}

func (c *communicatorReceiver) ExchangeHash(pulse core.PulseNumber, ctx context.Context,
	p consensus.Participant, data []*consensus.NodeUnsyncHash) ([]byte, error) {

	currentPulse := c.keeper.GetPulse()
	if currentPulse > pulse {
		return nil, errors.Errorf("Received consensus unsync hash exchange request with pulse %d but current is %d",
			pulse, currentPulse)
	}
	// TODO: block on getting unsync hash if currentPulse < number
	// hash, _, err := c.keeper.GetUnsyncHash()
	// if err != nil {
	// 	return nil, errors.Wrap(err, "Failed to calculate unsync hash")
	// }
	// TODO: write to communicatorSender map to decrease network requests
	// return hash, nil
	return nil, errors.New("not implemented")
}

func (c *communicatorSender) ExchangeData(pulse core.PulseNumber, ctx context.Context,
	p consensus.Participant, data []*core.ActiveNode) ([]*core.ActiveNode, error) {

	log.Infof("Sending consensus unsync list exchange request to %s", p.GetActiveNode().NodeID)
	sender, receiver, err := c.getSenderAndReceiver(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeData: error sending data to remote party")
	}
	request := packet.NewBuilder().Type(packet.TypeExchangeUnsyncLists).Sender(sender).Receiver(receiver).
		Request(&packet.RequestExchangeUnsyncLists{Pulse: pulse, UnsyncList: data}).Build()
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

func (c *communicatorSender) ExchangeHash(pulse core.PulseNumber, ctx context.Context,
	p consensus.Participant, data []*consensus.NodeUnsyncHash) ([]byte, error) {

	log.Infof("Sending consensus unsync hash exchange request to %s", p.GetActiveNode().NodeID)
	sender, receiver, err := c.getSenderAndReceiver(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeHash: error sending data to remote party")
	}
	request := packet.NewBuilder().Type(packet.TypeExchangeUnsyncHash).Sender(sender).Receiver(receiver).
		Request(&packet.RequestExchangeUnsyncHash{Pulse: pulse, UnsyncHash: data}).Build()
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
	receiverID := nodenetwork.ResolveHostID(p.GetActiveNode().NodeID)
	receiver, exists, err := c.handler.FindHost(ctx, receiverID)
	if err != nil || !exists {
		return nil, nil, errors.Wrap(err, "Error resolving receiver HostID -> Address")
	}
	return sender, receiver, nil
}
