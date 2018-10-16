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
	handler hosthandler.HostHandler
	keeper  nodekeeper.NodeKeeper
}

type communicatorSender struct {
	handler hosthandler.HostHandler
	keeper  nodekeeper.NodeKeeper
}

func (c *communicatorReceiver) ExchangeData(ctx context.Context, pulse core.PulseNumber,
	from core.RecordRef, data []*core.ActiveNode) ([]*core.ActiveNode, error) {

	// TODO: pass appropriate timeout
	unsyncHolder, err := c.keeper.GetUnsyncHolder(pulse, c.handler.GetPacketTimeout())
	if err != nil {
		return nil, errors.Wrap(err, "Error getting unsync holder on receiving side")
	}
	if unsyncHolder.GetPulse() > pulse {
		return nil, errors.Errorf("Received consensus unsync list exchange request with pulse %d but current is %d",
			pulse, unsyncHolder.GetPulse())
	}
	result := unsyncHolder.GetUnsync()
	unsyncHolder.AddUnsyncList(from, data)
	return result, nil
}

func (c *communicatorReceiver) ExchangeHash(ctx context.Context, pulse core.PulseNumber,
	from core.RecordRef, data []*consensus.NodeUnsyncHash) ([]*consensus.NodeUnsyncHash, error) {

	// TODO: pass appropriate timeout
	unsyncHolder, err := c.keeper.GetUnsyncHolder(pulse, c.handler.GetPacketTimeout())
	if err != nil {
		return nil, errors.Wrap(err, "Error getting unsync holder on receiving side")
	}
	if unsyncHolder.GetPulse() > pulse {
		return nil, errors.Errorf("Received consensus unsync hash exchange request with pulse %d but current is %d",
			pulse, unsyncHolder.GetPulse())
	}
	unsyncHolder.AddUnsyncHash(from, data)
	// TODO: pass appropriate timeout
	return unsyncHolder.GetHash(c.handler.GetPacketTimeout())
}

func (c *communicatorSender) ExchangeData(ctx context.Context, pulse core.PulseNumber,
	p consensus.Participant, data []*core.ActiveNode) ([]*core.ActiveNode, error) {

	unsyncHolder, err := c.keeper.GetUnsyncHolder(pulse, -1)
	if err != nil {
		log.Debugf("ExchangeData: error getting cache in consensus: " + err.Error())
	} else if result, ok := unsyncHolder.GetUnsyncList(p.GetActiveNode().NodeID); ok {
		log.Debugf("ExchangeData: got unsync list of remote party %s from cache", p.GetActiveNode().NodeID)
		return result, nil
	}
	log.Debugf("Sending consensus unsync list exchange request to %s", p.GetActiveNode().NodeID)
	sender, receiver, err := c.getSenderAndReceiver(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeData: error sending data to remote party")
	}
	log.Warn("recId: %s receiver != nil: %s", p.GetActiveNode().NodeID.String(), receiver != nil)
	request := packet.NewBuilder().Type(packet.TypeExchangeUnsyncLists).Sender(sender).Receiver(receiver).
		Request(&packet.RequestExchangeUnsyncLists{SenderID: c.keeper.GetID(), Pulse: pulse, UnsyncList: data}).Build()
	f, err := c.handler.SendRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeData: error sending data to remote party")
	}
	// TODO: fail fast if we don't have enough time to execute network request
	// TODO: count delay between two pulses and measure appropriate timeout. Maybe we also should pass time label
	// to instruct the receiving side how much time it has to process the operation
	response, err := f.GetResult(c.handler.GetPacketTimeout())
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeData: error sending data to remote party")
	}
	responseData := response.Data.(*packet.ResponseExchangeUnsyncLists)
	if responseData.Error != "" {
		return nil, errors.New("ExchangeData: got error from remote party: " + responseData.Error)
	}
	log.Debugf("ExchangeData: got unsync list of remote party %s from network", p.GetActiveNode().NodeID)
	return responseData.UnsyncList, nil
}

func (c *communicatorSender) ExchangeHash(ctx context.Context, pulse core.PulseNumber,
	p consensus.Participant, data []*consensus.NodeUnsyncHash) ([]*consensus.NodeUnsyncHash, error) {

	unsyncHolder, err := c.keeper.GetUnsyncHolder(pulse, -1)
	if err != nil {
		log.Debugf("ExchangeHash: error getting cache in consensus: " + err.Error())
	} else if result, ok := unsyncHolder.GetUnsyncHash(p.GetActiveNode().NodeID); ok {
		log.Debugf("ExchangeHash: got unsync hash of remote party %s from cache", p.GetActiveNode().NodeID)
		return result, nil
	}
	log.Debugf("Sending consensus unsync hash exchange request to %s", p.GetActiveNode().NodeID)
	sender, receiver, err := c.getSenderAndReceiver(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeHash: error sending data to remote party")
	}
	log.Warn("recId: %s, receiver != nil: %s", p.GetActiveNode().NodeID.String(), receiver != nil)
	request := packet.NewBuilder().Type(packet.TypeExchangeUnsyncHash).Sender(sender).Receiver(receiver).
		Request(&packet.RequestExchangeUnsyncHash{SenderID: c.keeper.GetID(), Pulse: pulse, UnsyncHash: data}).Build()
	f, err := c.handler.SendRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeHash: error sending data to remote party")
	}
	// TODO: fail fast if we don't have enough time to execute network request
	// TODO: count delay between two pulses and measure appropriate timeout. Maybe we also should pass time label
	// to instruct the receiving side how much time it has to process the operation
	response, err := f.GetResult(c.handler.GetPacketTimeout())
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeHash: error sending data to remote party")
	}
	responseData := response.Data.(*packet.ResponseExchangeUnsyncHash)
	if responseData.Error != "" {
		return nil, errors.New("ExchangeHash: got error from remote party: " + responseData.Error)
	}
	log.Debugf("ExchangeHash: got unsync hash of remote party %s from network", p.GetActiveNode().NodeID)
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
