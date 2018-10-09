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
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/nodekeeper"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/pkg/errors"
)

type communicatorReceiver struct {
	keeper nodekeeper.NodeKeeper
}

type communicatorSender struct {
	handler hosthandler.HostHandler
	network nodenetwork.NodeNetwork
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

	// nodeID := p.GetActiveNode().NodeID
	// hostID := c.network.ResolveHostID(nodeID)

	return nil, errors.New("not implemented")
}

func (c *communicatorSender) ExchangeHash(number core.PulseNumber, ctx context.Context,
	p consensus.Participant, data []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}
