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
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/dhtnetwork/hosthandler"
	"github.com/insolar/insolar/network/dhtnetwork/resolver"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

type communicatorReceiver struct {
	handler hosthandler.HostHandler
	keeper  network.NodeKeeper
}

type communicatorSender struct {
	handler hosthandler.HostHandler
	keeper  network.NodeKeeper
}

func (c *communicatorReceiver) ExchangeData(ctx context.Context, pulse core.PulseNumber,
	from core.RecordRef, data []core.Node) ([]core.Node, error) {

	// TODO: pass appropriate timeout
	unsyncHolder, err := c.keeper.GetUnsyncHolder(pulse, time.Second*5)
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
	from core.RecordRef, data []*network.NodeUnsyncHash) ([]*network.NodeUnsyncHash, error) {

	// TODO: pass appropriate timeout
	unsyncHolder, err := c.keeper.GetUnsyncHolder(pulse, time.Second*5)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting unsync holder on receiving side")
	}
	if unsyncHolder.GetPulse() > pulse {
		return nil, errors.Errorf("Received consensus unsync hash exchange request with pulse %d but current is %d",
			pulse, unsyncHolder.GetPulse())
	}
	unsyncHolder.AddUnsyncHash(from, data)
	// TODO: pass appropriate timeout
	return unsyncHolder.GetHash(time.Second * 5)
}

func (c *communicatorSender) ExchangeData(ctx context.Context, pulse core.PulseNumber,
	p consensus.Participant, data []core.Node) ([]core.Node, error) {

	unsyncHolder, err := c.keeper.GetUnsyncHolder(pulse, -1)
	if err != nil {
		log.Debugf("ExchangeData: error getting cache in consensus: " + err.Error())
	} else if result, ok := unsyncHolder.GetUnsyncList(p.GetActiveNode().ID()); ok {
		log.Debugf("ExchangeData: got unsync list of remote party %s from cache", p.GetActiveNode().ID())
		return result, nil
	}
	log.Debugf("Sending consensus unsync list exchange request to %s", p.GetActiveNode().ID())
	sender, receiver, err := c.getSenderAndReceiver(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeData: error sending data to remote party")
	}
	request := packet.NewBuilder(sender).
		Type(types.TypeExchangeUnsyncLists).
		Receiver(receiver).
		Request(&packet.RequestExchangeUnsyncLists{
			SenderID:   c.keeper.GetOrigin().ID(),
			Pulse:      pulse,
			UnsyncList: data,
		}).
		Build()
	f, err := c.handler.SendRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeData: error sending data to remote party")
	}
	// TODO: fail fast if we don't have enough time to execute network request
	// TODO: count delay between two pulses and measure appropriate timeout. Maybe we also should pass time label
	// to instruct the receiving side how much time it has to process the operation
	response, err := f.GetResult(time.Second * 5)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeData: error sending data to remote party")
	}
	responseData := response.Data.(*packet.ResponseExchangeUnsyncLists)
	if responseData.Error != "" {
		return nil, errors.New("ExchangeData: got error from remote party: " + responseData.Error)
	}
	log.Debugf("ExchangeData: got unsync list of remote party %s from network", p.GetActiveNode().ID())
	return responseData.UnsyncList, nil
}

func (c *communicatorSender) ExchangeHash(ctx context.Context, pulse core.PulseNumber,
	p consensus.Participant, data []*network.NodeUnsyncHash) ([]*network.NodeUnsyncHash, error) {

	unsyncHolder, err := c.keeper.GetUnsyncHolder(pulse, -1)
	if err != nil {
		log.Debugf("ExchangeHash: error getting cache in consensus: " + err.Error())
	} else if result, ok := unsyncHolder.GetUnsyncHash(p.GetActiveNode().ID()); ok {
		log.Debugf("ExchangeHash: got unsync hash of remote party %s from cache", p.GetActiveNode().ID())
		return result, nil
	}
	log.Debugf("Sending consensus unsync hash exchange request to %s", p.GetActiveNode().ID())
	sender, receiver, err := c.getSenderAndReceiver(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeHash: error sending data to remote party")
	}
	request := packet.NewBuilder(sender).
		Type(types.TypeExchangeUnsyncHash).
		Receiver(receiver).
		Request(&packet.RequestExchangeUnsyncHash{
			SenderID:   c.keeper.GetOrigin().ID(),
			Pulse:      pulse,
			UnsyncHash: data,
		}).
		Build()
	f, err := c.handler.SendRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeHash: error sending data to remote party")
	}
	// TODO: fail fast if we don't have enough time to execute network request
	// TODO: count delay between two pulses and measure appropriate timeout. Maybe we also should pass time label
	// to instruct the receiving side how much time it has to process the operation
	response, err := f.GetResult(time.Second * 5)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeHash: error sending data to remote party")
	}
	responseData := response.Data.(*packet.ResponseExchangeUnsyncHash)
	if responseData.Error != "" {
		return nil, errors.New("ExchangeHash: got error from remote party: " + responseData.Error)
	}
	log.Debugf("ExchangeHash: got unsync hash of remote party %s from network", p.GetActiveNode().ID())
	return responseData.UnsyncHash, nil
}

func (c *communicatorSender) getSenderAndReceiver(ctx context.Context, p consensus.Participant) (*host.Host, *host.Host, error) {
	ht := c.handler.HtFromCtx(ctx)
	sender := ht.Origin
	receiverID := resolver.ResolveHostID(p.GetActiveNode().NodeID)
	receiver, exists, err := c.handler.FindHost(ctx, receiverID)
	if err != nil || !exists {
		return nil, nil, errors.Wrap(err, "Error resolving receiver HostID -> Address")
	}
	return sender, receiver, nil
}
